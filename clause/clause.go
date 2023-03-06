//generator实现了字句生成规则，函数存储到了generators这个map里面，在clause里实现完整的独立子句，封装成函数使用

package clause

import "strings"

type Type int  //定义一个Type类型，用于存放sql关键字

type Clause struct {
	sql map[Type]string
	sqlVars map[Type][]interface{}
}

//定义sql语句的关键字
const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
	UPDATE
	DELETE
	COUNT
)

//Set 添加特定类型的子句:Set 方法根据 Type 调用对应的 generator，生成该子句对应的 SQL 语句。
func (c *Clause) Set (name Type, vars ...interface{}) {   //设计成Set(WHERE,"Name = ?","Tom") 这样子，Type类型就是const里面的都符合
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]interface{})
	}
	sql, vars := generators[name](vars...)  //传入vars...给规则函数处理
	c.sql[name] = sql  //调用相应的字段，返回sql字句，将字句赋值给Clause
	c.sqlVars[name] = vars  
}

//构建生成最终的 SQL 和 SQLVars:Build 方法根据传入的 Type 的顺序，构造出最终的 SQL 语句。
//例如：Build(SELECT, WHERE, ORDERBY, LIMIT) 构造出完整的sql语句
func (c *Clause) Build(orders ...Type) (string, []interface{}) {
	var sqls []string
	var vars []interface{}
	for _, order := range orders {
		if sql, ok := c.sql[order]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, c.sqlVars[order]...)
		}
	}
	return strings.Join(sqls," "),vars
}
