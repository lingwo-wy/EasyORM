package session

import (
	"EasyORM/log"
	"reflect"
)

//Hooks 常量
//钩子与结构体绑定，也就是每个结构体需要实现各自的钩子

const (
	BeforeQuery  = "BeforeQuery" //查询前
	AfterQuery   = "AfterQuery"  //查询后
	BeforeUpdate = "BeforeUpdate" //更新前
	AfterUpdate  = "AfterUpdate"  //更新后
	BeforeDelete = "BeforeDelete"  //删除前
	AfterDelete  = "AfterDelete"   //删除后
	BeforeInsert = "BeforeInsert"  //插入前
	AfterInsert  = "AfterInsert"   //插入后
)


//调用方法调用已注册的钩子
//s.RefTable().Model 或 value 即当前会话正在操作的对象，使用 MethodByName 方法反射得到该对象的方法(绑定结构体的方法)。
//将 s *Session 作为入参调用。每一个钩子的入参类型均是 *Session。
//将 CallMethod() 方法在 Find、Insert、Update、Delete 方法内部调用即可。
func (s *Session) CallMethod(method string, value interface{}) {
	fm := reflect.ValueOf(s.RefTable().Model).MethodByName(method)  //fm得出的是该对象的方法：钩子常量
	if value != nil {                    //如果value非空，说明对应的增删改查操作有新的值需要处理
		fm = reflect.ValueOf(value).MethodByName(method)
	}
	param := []reflect.Value{reflect.ValueOf(s)}
	if fm.IsValid() {                   //判断一个 reflect.Value 对象是否有效，即它是否包含一个非零值。
		if v := fm.Call(param); len(v) > 0 {    //用于调用一个函数类型值，其参数为一个 reflect.Value 类型的参数列表，并返回一个 reflect.Value 类型的结果列表。
			if err, ok := v[0].Interface().(error); ok {  //
				log.Error(err)
			}
		}
	}
}

//将上面的MethodByName(method)修改成interface实现

// type TBeforeQuery interface {
// 	BeforeQuery(s *Session) error
// }
// type TAfterQuery interface {
// 	AfterQuery(s *Session) error
// }
// type TBeforeUpdate interface {
// 	BeforeUpdate(s *Session) error
// }
// type TAfterUpdate interface {
// 	AfterUpdate(s *Session) error
// }
// type TBeforeDelete interface {
// 	BeforeDelete(s *Session) error
// }
// type TAfterDelete interface {
// 	AfterDelete(s *Session) error
// }
// type TBeforeInsert interface {
// 	BeforeInsert(s *Session) error
// }
// type TAfterInsert interface {
// 	AfterInsert(s *Session) error
// }

// func (s *Session) CallMethod(method string, value interface{}) {
// 	param := reflect.ValueOf(value)
//     switch method {
//     case BeforeQuery:
//         if i, ok := param.Interface().(TBeforeQuery); ok {
//             i.BeforeQuery(s)
//         }
//     case AfterQuery:
//         if i, ok := param.Interface().(TAfterQuery); ok {
//             i.AfterQuery(s)
//         }
// 	case BeforeUpdate:
//         if i, ok := param.Interface().(TBeforeUpdate); ok {
//             i.BeforeUpdate(s)
//         }
// 	case AfterUpdate:
//         if i, ok := param.Interface().(TAfterUpdate); ok {
//             i.AfterUpdate(s)
//         }
// 	case BeforeDelete:
//         if i, ok := param.Interface().(TBeforeDelete); ok {
//             i.BeforeDelete(s)
//         }
// 	case AfterDelete:
//         if i, ok := param.Interface().(TAfterDelete); ok {
//             i.AfterDelete(s)
//         }
// 	case BeforeInsert:
//         if i, ok := param.Interface().(TBeforeInsert); ok {
//             i.BeforeInsert(s)
//         }
// 	case AfterInsert:
//         if i, ok := param.Interface().(TAfterInsert); ok {
//             i.AfterInsert(s)
//         }
//     default:
//         panic("不是支持钩子的方法")
//     }
// }
