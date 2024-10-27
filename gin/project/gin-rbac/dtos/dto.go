package dtos


import (
	"fmt"
	"reflect"
)

// ConvertDTOToModel 是一个将 DTO 转换成数据库模型的通用函数
//
// 描述:
// 该函数接收两个interface{}类型的参数，分别是DTO和模型对象。
// 注意：这里传参采用指针的形式，即DTO和模型对象都为指针，否则会报错。
// 它将DTO中的字段值复制到模型对象的相应字段中。
// 如果DTO或模型对象不是结构体类型，或者字段类型不匹配，则返回错误。
//
// 参数:
// - dto: 数据传输对象（&DTO）
// - model: 目标模型对象（&Model）
func ConvertDTOToModel(dto, model interface{}) error {
	dtoVal := reflect.ValueOf(dto)
	modelVal := reflect.ValueOf(model)

	if dtoVal.Elem().Kind() != reflect.Struct || modelVal.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("both arguments must be of type 'struct'")
	}

	dtoType := dtoVal.Elem().Type()
	modelType := modelVal.Elem().Type()

	for i := 0; i < modelVal.Elem().NumField(); i++ {
		modelField := modelVal.Elem().Field(i)
		modelFieldType := modelType.Field(i)

		// 检查 DTO 是否包含对应字段
		dtoField, ok := dtoType.FieldByName(modelFieldType.Name)
		if !ok {
			continue // 如果 DTO 中没有该字段，则跳过
		}

		dtoValue := dtoVal.Elem().FieldByIndex(dtoField.Index)

		if !modelField.CanSet() {
			continue // 如果字段不可设置，则跳过
		}

		if dtoValue.Type().AssignableTo(modelField.Type()) {
			modelField.Set(dtoValue)
		} else {
			return fmt.Errorf("field %s cannot be assigned from DTO to model: types do not match", modelFieldType.Name)
		}
	}

	return nil
}

// ConvertModelToDTO 是一个将数据库模型转换成 DTO 的通用函数
// 描述:
// 该函数接收两个interface{}类型的参数，分别是模型对象和DTO对象。
// 注意：这里传参采用指针的形式，即模型对象和DTO对象都为指针，否则会报错。
// 它将模型对象的字段值复制到DTO对象的相应字段中。
// 如果模型对象或DTO对象不是结构体类型，或者字段类型不匹配，则返回错误。
func ConvertModelToDTO(model, dto interface{}) error {
	modelVal := reflect.ValueOf(model)
	dtoVal := reflect.ValueOf(dto)

	if modelVal.Elem().Kind() != reflect.Struct || dtoVal.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("both arguments must be of type 'struct'")
	}

	modelType := modelVal.Elem().Type()
	dtoType := dtoVal.Elem().Type()

	for i := 0; i < dtoVal.Elem().NumField(); i++ {
		dtoField := dtoVal.Elem().Field(i)
		dtoFieldType := dtoType.Field(i)

		// 检查 model 是否包含对应字段
		modelField, ok := modelType.FieldByName(dtoFieldType.Name)
		if !ok {
			continue // 如果 model 中没有该字段，则跳过
		}

		modelValue := modelVal.Elem().FieldByIndex(modelField.Index)

		if !dtoField.CanSet() {
			continue // 如果字段不可设置，则跳过
		}
		// 检查 modelValue 的类型是否可以赋值给 dtoField
		if modelValue.Type().AssignableTo(dtoField.Type()) {
			dtoField.Set(modelValue)
		} else {
			return fmt.Errorf("field %s cannot be assigned from model to DTO: types do not match", dtoFieldType.Name)
		}
	}

	return nil
}
