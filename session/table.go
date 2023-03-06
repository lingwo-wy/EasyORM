//table.go是操作数据库表的相关代码
package session

import (
	"fmt"
	"EasyORM/log"
	"reflect"
	"strings"
	"EasyORM/schema"
)


//Model方法用于给refTable赋值，和gorm的一样，因为Parse函数的解析操作很耗时，将解析的结果保存到refTable中，当Model中的数据不发生变化时就不用再调用解析，直接复用refTable即可
func (s *Session) Model (value interface{}) *Session {
	//空值nil或者不同的model,更新refTable，即更新schema（存放表信息的）数据
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {  //这个model是结构体模型，定义为interface{}类型
		s.refTable = schema.Parse(value,s.dialect)  
	}
	return s
}

func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {   //不多加处理，如果有问题报日志即可，对表的操作不一定成功也没事
		log.Error("Model为空")
	}
	return s.refTable
}

//判断表是否存在
func (s *Session) HashTable() bool {
	sql, values := s.dialect.TableExistSQL(s.RefTable().Name)  //输入表名s.RefTable().Name的，返回对应的sql语句
	row := s.Raw(sql, values...).QueryRow() //单行查询
	var tmp string
	_ = row.Scan(&tmp)
	return tmp == s.RefTable().Name

}

//创建表函数，在这里设计只需要返回err就行了，提示日志在函数体实现
func (s *Session) CreateTable() error {  
	table := s.RefTable()  //先取出数据库表数据
	var columns []string
	for _, field := range table.Fields {  //这里遍历所有的字段
		columns = append(columns, fmt.Sprintf("%s %s %s",field.Name,field.Type,field.Tag))  //sql拼接追加到columns切片
	}
	desc := strings.Join(columns,",") //将columns切片所有的字符拼接，用,相隔
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", table.Name, desc)).Exec()   

	return err
}

//删除表

func (s *Session) DropTable() error {
	table := s.RefTable()
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", table.Name)).Exec()  
	return err
}

