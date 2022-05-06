package schema

/*
设计 Schema，利用反射(reflect)完成结构体和数据库表结构的映射，包括表名、字段名、字段类型、字段 tag 等
*/

import (
	"go/ast"
	"orm/dialect"
	"reflect"
)
type Field struct {
	Name string	// 字段名
	Type string	// 类型
	Tag string	// 约束条件
}

type Schema struct {
	Model interface{}	// 被映射的对象
	Name string	// 表名
	Fields []*Field
	FieldNames []string
	fieldMap map[string]*Field	// 字段名和field的映射
}

func (s *Schema) GetField(name string) *Field {
	return s.fieldMap[name]
}

/*
reflect
TypeOf(): 返回入参的类型
ValueOf(): 返回入参的值
Indirect(): 因为设计的入参是一个对象的指针，Indirect() 获取指针指向的实例
*/

func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model: dest,
		Name: modelType.Name(),
		fieldMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}

			if v, ok := p.Tag.Lookup("orm"); ok {
				field.Tag = v
			}

			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}

	return schema
}

func (s *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))

	var fieldValues []interface{}

	for _, field := range s.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}

	return fieldValues
}