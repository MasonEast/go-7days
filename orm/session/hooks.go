package session

import (
	"orm/log"
	"reflect"
)

// Hooks constants
const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

// 通过反射找到注册的钩子并执行
func (s *Session) CallMethod(method string, value ...interface{}){
	fm := reflect.ValueOf(s.RefTable().Model).MethodByName(method)

	if value != nil {
		fm = reflect.ValueOf(value).MethodByName(method)
	}
	log.Info("fm.isva------",fm, fm.IsValid(), method)

	param := []reflect.Value{reflect.ValueOf(s)}
	if fm.IsValid() {
		// 每个钩子将session会话作为入参
		if v := fm.Call(param); len(v) > 0 {
			if err, ok := v[0].Interface().(error); ok {
				log.Error(err)
			}
		}
	}
}