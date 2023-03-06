//实现

package dialect

import (
	"reflect"  
)

type Dialect interface {
	DataTypeOf(Type reflect.Value) string  // 用于将 Go 语言的类型转换为该数据库的数据类型。因为不知道传入的值类型，只能用反射确定
	TableExistSQL(tableName string) (string, []interface{})  // TableExistSQL 返回某个表是否存在的 SQL 语句，参数是表名(table)。
}

var dialectsMap = map[string]Dialect{}

//注册dialect实例，新增对某个数据库的支持时，调用可以注册到全局
func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

//取出dialectsMap的name对应的值：Dialect{}
func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}