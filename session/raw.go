// 实现直接调用SQL进行原生交互(与数据库交互)
// 使用go原生的sql库
package session

import (
	"EasyORM/clause"
	"EasyORM/dialect"
	"EasyORM/log"
	"EasyORM/schema"
	"database/sql"
	"strings"
)

//Session保留指向sql.DB的指针，并提供所有类型数据库操作的所有执行。
type Session struct {
	db *sql.DB  //使用sql.Open连接数据库，成功以后返回的指针
	sql strings.Builder  //拼接SQL语句
	sqlVars []interface{}  //拼接SQL 语句中占位符的对应值（?）

	dialect dialect.Dialect   //2编写对象表结构映射模块时加的依赖
	refTable *schema.Schema

	clause clause.Clause  //3编写记录新增和查询

	tx *sql.Tx  //6增加事务
}




//以下将Session实例化，绑定对应的方法

//新建session实例
func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
	}
}

//初始化会话的状态
func (s *Session) Clear() {  //将Reset再封装成Clear以便于调用
	s.sql.Reset()  //Reset方法：清空Builder中存储的字符串内容，同时也会清空底层的[]byte缓冲区。
	s.sqlVars = nil
	s.clause = clause.Clause{}   //清空Clause结构体里面存储的信息，赋空。
}


//CommonDB 是 db 的最小函数集
type CommonDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

var _ CommonDB = (*sql.DB)(nil)
var _ CommonDB = (*sql.Tx)(nil)

// 如果 tx 开始，则 DB 返回 tx。否则返回 *SQL.DB
func (s *Session) DB() CommonDB {
	if s.tx != nil {  //增加了Tx事务
		return s.tx
	}
	return s.db
}


//Raw追加sql和sqlVars；且Raw后支持链式，因为返回了一个session对象
func (s *Session) Raw(sql string, values ...interface{}) *Session {  
	s.sql.WriteString(sql)  // strings.Builder里的方法：该方法会将字符串转换为 []byte 类型，然后通过调用 Write 方法将其追加到 buf 缓冲区中。
	s.sql.WriteString("")   //s.String()应该是sql+""
	s.sqlVars = append(s.sqlVars, values...)
	return s
}


//封装Exec、Query、QueryRow 原生方法

//Exec用于执行增删改等操作，后接QueryRow等方法
func (s *Session) Exec() (result sql.Result, err error) {  //使用sqlVars执行原生的sql
	defer s.Clear()  //延时清空缓冲区（存储的字符串内容）
	log.Info(s.sql.String(), s.sqlVars)  //这里的s.sql.String的值是raw中已经拼接好了的
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return

}

// QueryRow从数据库获取记录:原生的sql.QueryRow可以从数据库获取单条记录

func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)  
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...) 
}

//QueryRows获取数据库的记录列表(多行记录)

func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)  //日志输出对应的sql语句
	if rows, err =  s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)  //打印错误日志
	}
	return 
}

