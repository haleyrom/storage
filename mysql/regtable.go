package mysql

import (
	"fmt"
	"reflect"
	"sync"
)

// tableTypeRegistry tableTypeRegistry
var tableTypeRegistry sync.Map

// RegisterTableObj RegisterTableObj
func RegisterTableObj(tableObj interface{}) {
	// 得到类型
	tableObjType := reflect.Indirect(reflect.ValueOf(tableObj)).Type()
	tableObjFullName := tableObjType.String()
	tableObjName := tableObjType.Name()

	fmt.Println("tableObjFullName:", tableObjFullName)
	fmt.Println("tableObjName:", tableObjName)

	tableTypeRegistry.Store(tableObjFullName, tableObjType)
	tableTypeRegistry.Store(tableObjName, tableObjType)
}

// NewTableObj NewTableObj
func NewTableObj(tblStructName string) (interface{}, error) {
	typeV, exist := tableTypeRegistry.Load(tblStructName)
	if exist {
		if tableType, ok := typeV.(reflect.Type); ok {
			newTableObj := reflect.New(tableType)
			return newTableObj.Interface(), nil
		}
	}

	return nil, fmt.Errorf("tblStructName[%s] no registration", tblStructName)
}
