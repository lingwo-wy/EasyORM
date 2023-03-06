// 实现字句生成规则,拼接SQL语句
package clause

import (
	"fmt"
	"strings"
)

type generator func(values ...interface{}) (string, []interface{})

var generators map[Type]generator

func init() {  

	generators = make(map[Type]generator)
	generators[INSERT] = _insert   //调用generators[INSERT]=(),就会调用_insert函数
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT]  = _limit
	generators[WHERE]  = _where
	generators[ORDERBY]= _orderBy
	generators[UPDATE] = _update
	generators[DELETE] = _delete
	generators[COUNT]  = _count
}

//占位符和,拼接
func genBindVars(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")  //将占位符存到vars
	}
	return strings.Join(vars,", ")  //拼接成 ?, ? 这样子
}


//Set(WHERE,"Name = ?","Tom") 要设计成这个样子的情况下，这里的规则函数要通过Set这个封装函数拿值，所以规则函数的设计就遵循传入可变参数列表就行了，然后再对这个值处理
//函数返回一个用于关键字语句的查询字符串(即SQL语句)和一个要替换到所有占位符中的值的切片。然后被clause.go的函数封装，就能使用了


func _insert(values ...interface{}) (string, []interface{}) {  //加_只能在本包访问
	// INSERT INTO tableName(fileds)  //fileds是字段名
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")  

	return fmt.Sprintf("INSERT INTO %s (%v)", tableName, fields), []interface{}{}
}

//VALUES关键字部分
func _values(values ...interface{}) (string, []interface{}) {
	// VALUES (v1), (v2),...
	var bindStr string
	var sql strings.Builder
	var vars []interface{}
	sql.WriteString("VALUES")  //将VALUES写入缓冲区
	for i, value := range values {
		v := value.([]interface{})
		if bindStr == "" {
			bindStr = genBindVars(len(v))  //传入values的长度，调用函数拼接占位符和，
		}
		sql.WriteString(fmt.Sprintf("(%v)", bindStr))
		if i + 1 != len(values) {
			sql.WriteString(", ")
		}
		vars = append(vars, v...)
	}

	return sql.String(), vars
}


//Set(SELECT, "User", []string{"*"})
func _select(values ...interface{}) (string, []interface{}) {
	// SELECT fields FROM tableName
	tableName := values[0]  //第一个元素是表名User
	fields := strings.Join(values[1].([]string), ",") //拼接 , 
	return fmt.Sprintf("SELECT %v FROM %s", fields, tableName), []interface{}{}
}


//Set(LIMIT, 3)
func _limit(values ...interface{}) (string, []interface{}) {
	// LIMIT num
	return "LIMIT ?", values
}

//Set(WHERE, "Name = ?", "Tom")
func _where(values ...interface{}) (string, []interface{}) {
	// WHERE desc
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %s", desc), vars
}

func _orderBy(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("ORDER BY %s", values[0]), []interface{}{}
}


//设计入参是2个，第一个参数是表名(table)，第二个参数是 map 类型，表示待更新的键值对。
//接受任意数量的接口参数作为输入参数，其中包含表名和一个包含列名及其更新值的 map，用于更新语句。
//然后迭代该映射，使用每个列名后跟“=？”作为新值的占位符，构建更新语句的SET子句。然后返回最终的查询字符串和一系列要替换到占位符中的值。
func _update(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	m := values[1].(map[string]interface{})
	var keys []string
	var vars []interface{}
	for k, v := range m {
		keys = append(keys, k+" = ?")
		vars = append(vars, v)
	}
	return fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(keys, ", ")), vars
}


//只有一个入参，即表名。
func _delete(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("DELETE FROM %s", values[0]), []interface{}{}  //直接操作
}

//只有一个入参，即表名，并复用 _select 生成器。
func _count(values ...interface{}) (string, []interface{}) {
	return _select(values[0], []string{"count(*)"})  //计数，values[1]是表中字段名
}