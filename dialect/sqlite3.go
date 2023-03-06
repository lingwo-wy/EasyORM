package dialect

import (
	"fmt"
	"reflect"
	"time"
)

type sqlite3 struct{}

var _ Dialect = (*sqlite3)(nil)  //判断*sqlite3中是否实现了Dialect这个方法，相当于类型断言

func init() {
	RegisterDialect("sqlite3",&sqlite3{})
}

//将go语言类型映射成sqlite3数据类型
func (s *sqlite3) DataTypeOf(typ reflect.Value) string {  
	switch typ.Kind() {   //取出typ的类型
	case reflect.Bool:  //当类型是bool时（kind返回的是reflect.Kind）
		return "bool"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,  //相应的整型对应integer
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		return "integer"
	case reflect.Int64, reflect.Uint64:  
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "real"
	case reflect.String:
		return "text"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))
}

// TableExistSQL 返回了在 SQLite 中判断表 tableName 是否存在的 SQL 语句。
func (s *sqlite3) TableExistSQL(tableName string) (string, []interface{}) {
	args := []interface{}{tableName}
	return "SELECT name FROM sqlite_master WHERE type='table' and name = ?", args
}
