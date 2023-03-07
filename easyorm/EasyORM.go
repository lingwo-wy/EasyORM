// 引擎部分：负责与数据库交互后的工作，关闭连接，这里也是geeorm和用户交互的入口
package easyorm

import (
	"EasyORM/dialect"
	"EasyORM/log"
	"EasyORM/session"
	"database/sql"
	"fmt"
	"strings"
)

type Engine struct {
	db *sql.DB
	dialect dialect.Dialect //2增加依赖
}

//连接数据库并返回*sql.DB,调用ping检查数据库连接.创建Engine实例，获取driver对应的dialect
func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver,source)  // 数据库驱动，数据源（链接地址等）
	if err != nil {
		log.Error(err)
		return
	}
	//发送ping测试数据库连接
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}

	dial, ok := dialect.GetDialect(driver)  //加入数据库驱动取出Dialect{},或者说确保特定的dialect存在

	if !ok {
		log.Errorf("dialect %s Not Found", driver)
		return
	}

	e = &Engine{db:db, dialect: dial}
	log.Info("连接数据库成功")

	return

}

func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		log.Error("数据库关闭失败")
	}
	log.Info("数据库关闭成功")
}

//为后续操作创建新会话
func (e *Engine) NewSession() *session.Session {
	return session.New(e.db,e.dialect)
}


//事务的交互部分
type TxFunc func(*session.Session) (interface{}, error)  //模仿gorm，在使用时，应该是要用户自定义成这个函数的自定义操作。

func (engine *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := engine.NewSession()
	if err := s.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {  //捕获错误，捕获到的情况下就回滚，或者上面的err!=nil也回滚，如果都没有问题再提交
			_ = s.Rollback()
			panic(p) // 回滚后再panic
		} else if err != nil {
			_ = s.Rollback() // 
		} else {
			err = s.Commit() // 
			if err != nil {
				_ = s.Rollback()  //提交失败再回滚
			}
		}
	}()

	return f(s)
}


//支持数据库表字段迁移
//仅支持字段新增与删除，不支持字段类型变更

//difference 返回a - b:新表 - 旧表 = 新增字段， 旧表-新表 = 删除字段
func difference( a, b []string) (diff []string) {
	Bmap := make(map[string]bool)
	for _, v := range b {
		Bmap[v] = true
	}
	for _, v := range a {
		if _, ok := Bmap[v]; !ok { //当v存在于Bmap即b中就不管，不存在再追加到diff
			diff = append(diff, v)
		}
	}
	return
}


//迁移表:使用ALTER（alter）语句新增字段，使用创建表并重命名的方式删除字段.
func (engine *Engine) Migrate(value interface{}) error {

	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {  //使用事务
		if !s.Model(value).HashTable() {    // 判断结构体对应的表是否存在
			log.Infof("表 %s 不存在", s.RefTable().Name)  
			return nil, s.CreateTable()  //不存在直接建表
		}
		table := s.RefTable()  //判断Model内是否为空，即refTable是否空
		rows, _ := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", table.Name)).QueryRows()
		columns, _ := rows.Columns()
		addCols := difference(table.FieldNames, columns)  //增加列
		delCols := difference(columns, table.FieldNames)  //删除列
		log.Infof("增加列 %v, 删除列 %v", addCols, delCols)

		for _, col := range addCols {  
			f := table.GetField(col)
			sqlStr := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", table.Name, f.Name, f.Type)
			if _, err = s.Raw(sqlStr).Exec(); err != nil {
				return
			}
		}

		if len(delCols) == 0 {  
			return
		}
		tmp := "tmp_" + table.Name
		fieldStr := strings.Join(table.FieldNames, ", ")
		s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s from %s;", tmp, fieldStr, table.Name))
		s.Raw(fmt.Sprintf("DROP TABLE %s;", table.Name))
		s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", tmp, table.Name))
		_, err = s.Exec()
		return
	})
	return err
}