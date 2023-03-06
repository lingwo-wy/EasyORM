//将子句拼接成完整的SQL语句，封装成可以被外部调用的函数。

package session

import (
	"EasyORM/clause"
	"reflect"
	"errors"
)


// 多次调用 clause.Set() 构造好每一个子句。
// 调用一次 clause.Build() 按照传入的顺序构造出最终的SQL语句。构造完成后，调用Raw().Exec()方法执行。
func (s *Session) Insert (values ...interface{}) (int64, error) {
	recordValues := make([]interface{},0)
	for _, value := range values {
		s.CallMethod(BeforeInsert, value)
		table := s.Model(value).RefTable()  //Model是解析入库，RefTable是判断映射表是否为空从而得知Model是否为空
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)  //实现子句
		recordValues = append(recordValues, table.RecordValues(value))  //将值转换再追加
	}
	s.clause.Set(clause.VALUES, recordValues...)   //构造子句
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES) //生成完整的SQL
	result, err := s.Raw(sql, vars...).Exec()  //用来执行sql语句
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterInsert, nil)
	return result.RowsAffected()  //返回插入行数
}


//Find的功能实现和Inser是反过来，Find要通过平铺字段值构造出对象

//应该设计成如下：
// s := geeorm.NewEngine("sqlite3", "gee.db").NewSession()
// var users []User
// s.Find(&users);

// destSlice.Type().Elem() 获取切片的单个元素的类型 destType，使用 reflect.New() 方法创建一个 destType 的实例，作为 Model() 的入参，映射出表结构 RefTable()。
// 根据表结构，使用 clause 构造出 SELECT 语句，查询到所有符合条件的记录 rows。
// 遍历每一行记录，利用反射创建 destType 的实例 dest，将 dest 的所有字段平铺开，构造切片 values。
// 调用 rows.Scan() 将该行记录每一列的值依次赋值给 values 中的每一个字段。
// 将 dest 添加到切片 destSlice 中。循环直到所有的记录都添加到切片 destSlice 中。

// destType 是作为参数传进来的，它表示一个用户定义的结构体类型。.Elem() 方法对指针进行解引用，返回指向的对象。因此，reflect.New(destType).Elem() 表示创建一个指向 destType 类型的新对象，
// 并返回这个新对象。这个操作得到的是一个新对象的 reflect.Value。Interface() 方法返回 reflect.Value 的接口表示形式。它用于将一个 reflect.Value 转换为普通的 Go 变量，为之后的操作做准备。
// 在这里，它将之前得到的新对象的 reflect.Value 转换为了一个普通的空接口类型，以便在之后可以作为 Model 的参数进行传递。
// 最后，调用 RefTable() 方法获取模型对应的映射表。
func (s *Session) Find (values interface{}) error {

	s.CallMethod(BeforeQuery, nil)
	destSlice := reflect.Indirect(reflect.ValueOf(values)) // 实际上ValueOf返回values的指针，再通过Indirect就取出Values的值
	destType := destSlice.Type().Elem()  //。在指向类型的指针变量的 reflect.Value 上调用 Type().Elem() 方法时，它会返回这个指针所指向的对象类型的 reflect.Type。
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()  
	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)

	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)  //生成SQL整句
	rows, err := s.Raw(sql, vars...).QueryRows()   //读取多行
	if err != nil {
		return err
	}

	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var values []interface{}
		for _, name := range table.FieldNames {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		if err := rows.Scan(values...); err != nil {
			return err
		}
		s.CallMethod(AfterQuery, dest.Addr().Interface())
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}


//Update
//支持 map[string]interface{},还支持kv列表："Name","Tom","Age",18,...

func (s *Session) Update (mapkv ...interface{}) (int64, error) {

	s.CallMethod(BeforeUpdate, nil)
	m, ok := mapkv[0].(map[string]interface{})
	if !ok {
		m = make(map[string]interface{})
		for i := 0; i < len(mapkv); i += 2 {
			m[mapkv[i].(string)] = mapkv[i+1]
		}
	}
	s.clause.Set(clause.UPDATE, s.RefTable().Name, m)  //生成子句
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE) //生成完整的SQL语句。
	result, err := s.Raw(sql, vars...).Exec()  //执行SQL
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterUpdate, nil)
	return result.RowsAffected()   //返回改变的数据行数
}

//Delete 删除带有 where 子句的记录

func (s *Session) Delete() (int64, error) {

	s.CallMethod(BeforeDelete, nil)
	s.clause.Set(clause.DELETE, s.RefTable().Name)  //生成删除子句
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)  
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterDelete, nil)
	return result.RowsAffected()
}

// Count 使用 where 子句对记录进行计数
func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vars...).QueryRow()  
	var tmp int64
	if err := row.Scan(&tmp); err != nil {
		return 0, err
	}
	return tmp, nil
}

// Limit 将限制条件添加到子句
func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)  
	return s
}

//Where 向子句添加Where条件
func (s *Session) Where(desc string, args ...interface{}) *Session {
	var vars []interface{}
	s.clause.Set(clause.WHERE, append(append(vars, desc), args...)...)
	return s
}

//OrderBy 将Order by条件添加到子句中
func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.ORDERBY, desc)
	return s
}

//First查询一条记录。实现原理：根据传入的类型，利用反射构造切片，调用 Limit(1) 限制返回的行数，调用 Find 方法获取到查询结果。
func (s *Session) First(value interface{}) error {
	dest := reflect.Indirect(reflect.ValueOf(value))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {  //复用Limit和Find
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("NOT FOUND")  //找不到记录
	}
	dest.Set(destSlice.Index(0))
	return nil
}


