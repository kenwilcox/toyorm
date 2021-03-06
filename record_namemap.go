/*
 * Copyright 2018. bigpigeon. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package toyorm

import (
	"fmt"
	"reflect"
)

type ModelNameMapRecords struct {
	model           *Model
	elemType        reflect.Type
	source          reflect.Value
	FieldValuesList []map[string]reflect.Value
}

func NewNameMapRecords(model *Model, v reflect.Value) *ModelNameMapRecords {
	records := &ModelNameMapRecords{
		model,
		v.Type().Elem(),
		v,
		[]map[string]reflect.Value{},
	}
	records.sync()
	return records
}

func NewNameMapRecord(model *Model, v reflect.Value) *ModelNameMapRecord {
	record := &ModelNameMapRecord{
		map[string]reflect.Value{},
		v,
		model,
	}

	for name := range model.GetNameFieldMap() {
		if fieldValue := v.MapIndex(reflect.ValueOf(name)); fieldValue.IsValid() {
			fieldValue = fieldValue.Elem()
			record.FieldValues[name] = fieldValue
		}
	}
	return record
}

func (m *ModelNameMapRecords) sync() {
	for i := len(m.FieldValuesList); i < m.source.Len(); i++ {
		// why need LoopIndirect? because m could be []*map[string]interface{}
		elem := LoopIndirect(m.source.Index(i))
		c := map[string]reflect.Value{}
		for name := range m.model.GetNameFieldMap() {
			if elemField := elem.MapIndex(reflect.ValueOf(name)); elemField.IsValid() {
				elemField = elemField.Elem()
				c[name] = elemField
			}
		}
		m.FieldValuesList = append(m.FieldValuesList, c)
	}
}

func (m *ModelNameMapRecords) GetRecords() []ModelRecord {
	var recordList []ModelRecord
	for i := 0; i < len(m.FieldValuesList); i++ {
		recordList = append(recordList, &ModelNameMapRecord{
			FieldValues: m.FieldValuesList[i],
			source:      LoopIndirect(m.source.Index(i)),
			model:       m.model,
		})
	}
	return recordList
}

func (m *ModelNameMapRecords) GroupBy(key string) ModelGroupBy {
	field := m.model.GetFieldWithName(key)
	if field.StructField().Type.Comparable() == false {
		panic(fmt.Sprintf("%s is not compareable field", field.Name()))
	}
	result := ModelGroupBy{}
	for i := 0; i < len(m.FieldValuesList); i++ {
		keyValue := m.FieldValuesList[i][key].Interface()
		result[keyValue] = append(result[keyValue], ModelIndexRecord{&ModelNameMapRecord{
			FieldValues: m.FieldValuesList[i],
			source:      LoopIndirect(m.source.Index(i)),
			model:       m.model,
		}, i})
	}
	return result
}

func (m *ModelNameMapRecords) GroupByFunc(key string, fn func(int, ModelRecord)) {
	field := m.model.GetFieldWithName(key)
	if field.StructField().Type.Comparable() == false {
		panic(fmt.Sprintf("%s is not compareable field", field.Name()))
	}
	//result := ModelGroupBy{}
	for i := 0; i < len(m.FieldValuesList); i++ {
		//keyValue := m.FieldValuesList[i][key].Interface()
		fn(i, &ModelNameMapRecord{
			FieldValues: m.FieldValuesList[i],
			source:      LoopIndirect(m.source.Index(i)),
			model:       m.model,
		})

	}
}

func (m *ModelNameMapRecords) GetRecord(i int) ModelRecord {
	return &ModelNameMapRecord{
		FieldValues: m.FieldValuesList[i],
		source:      LoopIndirect(m.source.Index(i)),
		model:       m.model,
	}
}

func (m *ModelNameMapRecords) Add(v reflect.Value) ModelRecord {
	if m.source.CanSet() == false {
		panic("Add need can set permission")
	}
	last := len(m.FieldValuesList)
	m.source.Set(reflect.Append(m.source, v))
	m.sync()
	return &ModelNameMapRecord{
		FieldValues: m.FieldValuesList[last],
		source:      LoopIndirect(m.source.Index(last)),
		model:       m.model,
	}
}

func (m *ModelNameMapRecords) GetFieldType(name string) reflect.Type {
	field := m.model.GetFieldWithName(name)
	fieldType := field.StructField().Type

	// TODO check the field is or not container field ?
	if field.Column() == "" || field.SqlType() == "" {
		t := LoopTypeIndirect(fieldType)
		switch t.Kind() {
		case reflect.Struct:
			return reflect.TypeOf(map[string]interface{}{})
		case reflect.Slice:
			return reflect.TypeOf([]map[string]interface{}{})
		}
	}
	return fieldType
}

func (m *ModelNameMapRecords) GetFieldAddressType(name string) reflect.Type {
	return m.GetFieldType(name)
}

func (m *ModelNameMapRecords) IsVariableContainer() bool {
	return true
}

func (m *ModelNameMapRecords) ElemType() reflect.Type {
	return m.elemType
}

func (m *ModelNameMapRecords) Len() int {
	return len(m.FieldValuesList)
}

func (m *ModelNameMapRecords) Source() reflect.Value {
	return m.source
}

type ModelNameMapRecord struct {
	FieldValues map[string]reflect.Value
	source      reflect.Value
	model       *Model
}

func (m *ModelNameMapRecord) SetField(name string, value reflect.Value) {
	if name == "" {
		return
	}
	if field := m.model.GetFieldWithName(name); field != nil {
		nameValue := reflect.ValueOf(name)
		elem := reflect.New(field.StructField().Type).Elem()
		safeSet(elem, value)
		m.source.SetMapIndex(nameValue, elem)
		m.FieldValues[name] = m.source.MapIndex(nameValue).Elem()
	}
}

func (m *ModelNameMapRecord) Field(name string) reflect.Value {
	return m.FieldValues[name]
}

func (m *ModelNameMapRecord) FieldAddress(name string) reflect.Value {
	return m.FieldValues[name]
}

func (m *ModelNameMapRecord) AllField() map[string]reflect.Value {
	return m.FieldValues
}

func (m *ModelNameMapRecord) IsVariableContainer() bool {
	return true
}

func (m *ModelNameMapRecord) Source() reflect.Value {
	return m.source
}

func (m *ModelNameMapRecord) GetFieldType(name string) reflect.Type {
	field := m.model.GetFieldWithName(name)
	fieldType := field.StructField().Type

	// TODO check the field is or not container field ?
	if field.Column() == "" || field.SqlType() == "" {
		t := LoopTypeIndirect(fieldType)
		switch t.Kind() {
		case reflect.Struct:
			return reflect.TypeOf(map[uintptr]interface{}{})
		case reflect.Slice:
			return reflect.TypeOf([]map[uintptr]interface{}{})
		}
	}
	return fieldType
}
