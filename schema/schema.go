package schema

import (
	"EasyORM/dialect"
	"go/ast"  //抽象语法树
	"reflect"
)

// Field表示数据库的一列
type Field struct {
	Name string  //字段名
	Type string  //类型
	Tag  string  //约束条件Tag
}


//schema 表示数据库表
type Schema struct {
	Model  	    interface{}       //被映射对象
	Name   	    string            //表名
	Fields 	    []*Field          //字段
	FieldNames   []string		  //字段名：包含所有的表列名
	fieldMap    map[string]*Field //记录列名和Field的映射关系，方便使用，不用再在使用的时候遍历Fields
}


//获取field字段名字对应的Field信息
func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}


type ITableName interface {
	TableName() string
}


//Parse 解析函数，将任意的对象解析为Schema实例
func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type() //TypeOf() 和 ValueOf() 分别用来返回入参的类型和值。因为设计的入参是一个对象的指针，因此需要reflect.Indirect()获取指针指向的实例。
	schema := &Schema{
		Model: dest,
		Name: modelType.Name(),  //获取到结构体的名称，用来作为表名
		fieldMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {  //modelType.NumField()返回字段的数量,通过下标获取到特定字段 p := modelType.Field(i)。
		p := modelType.Field(i)  //返回结构类型的第 i 个字段。
		if !p.Anonymous && ast.IsExported(p.Name) {   //Anonymous判断p是否匿名，IsExported判断p.Name（即字段名称）是否以大写字母开头:即if 不匿名且非大写
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))), //转换成数据库字段类型
			}
			if v, ok := p.Tag.Lookup("geeorm"); ok {   //查找与标签字符串中的键相关联的值，如果标签中没有此类键，则返回一个空字符串，并确定是否将标签明确设置为空字符串。（标签关键字是geeorm，必须使用）
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field  //每一个字段名对应的键是一个完整的字段信息
		}

	}

	return schema
}

//希望通过这种方式调用函数
// s := geeorm.NewEngine("sqlite3", "gee.db").NewSession()
// u1 := &User{Name: "Tom", Age: 18}
// u2 := &User{Name: "Sam", Age: 25}
// s.Insert(u1, u2, ...)
//那么就要加一个转换步骤：将u1和u2转成（"Tom",18）（"Sam",25）这种形式

func (schema *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))  //取值
	var fieldValues []interface{}
	for _, field := range schema.Fields {  //遍历字段列表
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}