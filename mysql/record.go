package mysql

import (
	"errors"
	"fmt"
	"github.com/haleyrom/storage/logging"
	"reflect"
	"strings"

	"github.com/go-xorm/builder"
	"github.com/go-xorm/core"
)

// InsertRecord InsertRecord
func InsertRecord(lc *logging.LogContext, engine *DBEngineInfoS, tableName string, object interface{}, omitColumns ...string) (int64, error) {
	if tableName == "" {
		return 0, fmt.Errorf("table name is empty")
	}
	if engine != nil {
		err := engine.RealEngine.Ping()
		if err == nil {
			affected, err := engine.RealEngine.Table(tableName).Omit(omitColumns...).Insert(object)
			if err != nil {
				logging.Warn(lc, "InsertMultiRecords Table[%s] failed! reason[%s]", tableName,
					err.Error())
				return -1, err
			}

			return affected, nil
		}

		return -1, err
	}

	return -1, fmt.Errorf("engine is nil")
}

// InsertSliceRecordsToSameTable InsertSliceRecordsToSameTable
func InsertSliceRecordsToSameTable(lc *logging.LogContext, engine *DBEngineInfoS, tableName string, sliceObjs interface{}) (int64, error) {
	return InsertRecord(lc, engine, tableName, sliceObjs)
}

// UpdateRecord UpdateRecord
func UpdateRecord(engine *DBEngineInfoS, tableName string, primaryKeys *core.PK, objPtr interface{}) (affectedRows int64, err error) {
	affectedRows = 0
	err = nil

	if engine != nil && len(tableName) > 0 {
		err = engine.RealEngine.Ping()
		if err == nil {
			affectedRows, err = engine.RealEngine.Table(tableName).Id(primaryKeys).AllCols().Update(objPtr)
			return
		}
	}
	return 0, fmt.Errorf("engine is nil or tableName is empty")
}

// UpdateRecordCols UpdateRecordCols
func UpdateRecordCols(engine *DBEngineInfoS, tableName string, primaryKeys *core.PK, objPtr interface{}, updateCols ...string) (affectedRows int64, err error) {
	affectedRows = 0
	err = nil

	if engine != nil && len(tableName) > 0 {
		err = engine.RealEngine.Ping()
		if err == nil {
			session := engine.RealEngine.Table(tableName)
			session = session.Id(primaryKeys)

			if len(updateCols) > 0 {
				tableInfo := engine.RealEngine.TableInfo(objPtr)
				omitCols := make([]string, 0)
				for _, col := range tableInfo.Columns() {
					if col.IsPrimaryKey {
						continue
					}

					update := false
					for _, colName := range updateCols {
						lowerColName := strings.ToLower(colName)
						if strings.ToLower(col.Name) == lowerColName {
							update = true
							break
						}
					}

					if update {
						continue
					}

					omitCols = append(omitCols, col.Name)
				}
				session = session.Omit(omitCols...)
			}

			affectedRows, err = session.AllCols().Update(objPtr)

			return
		}
	}
	return 0, fmt.Errorf("engine is nil or tableName is empty")
}

// UpdateRecordSpecifiedFieldsByCondEx UpdateRecordSpecifiedFieldsByCondEx
func UpdateRecordSpecifiedFieldsByCondEx(engine *DBEngineInfoS, tableName string, cond *builder.Cond, fieldMap map[string]interface{}) (affectedRows int64, err error) {
	affectedRows = 0
	err = nil
	if engine == nil {
		return 0, fmt.Errorf("UpdateRecordSpecifiedFieldsByCond,engine is nil")
	}
	if len(tableName) <= 0 {
		return 0, fmt.Errorf("UpdateRecordSpecifiedFieldsByCond,tableName is empty")
	}
	if cond == nil {
		return 0, fmt.Errorf("UpdateRecordSpecifiedFieldsByCond,cond is empty")
	}

	err = engine.RealEngine.Ping()
	if err == nil {
		affectedRows, err = engine.RealEngine.Table(tableName).Where(*cond).Update(fieldMap)
		return
	}
	return
}

// GetRecord 根据primary key 查询得到一条记录
func GetRecord(engine *DBEngineInfoS, tableName string, objPtr interface{}) (bool, error) {
	if engine != nil {
		objV := reflect.ValueOf(objPtr)
		if objV.Kind() == reflect.Ptr {
			err := engine.RealEngine.Ping()
			if err == nil {
				exist, err := engine.RealEngine.Table(tableName).Get(objPtr)
				if exist && err == nil {
					return true, nil
				} else if !exist && err == nil {
					return false, nil
				} else {
					return false, err
				}
			}
		}

		return false, fmt.Errorf("objPtr kind is not Ptr")
	}
	return false, fmt.Errorf("engine is nil")
}

// GetRecordByCond 根据查询条件来查询一条记录
func GetRecordByCond(engine *DBEngineInfoS, tableName string, cond string, objPtr interface{}) (bool, error) {
	if engine != nil {
		objV := reflect.ValueOf(objPtr)
		if objV.Kind() == reflect.Ptr {
			err := engine.RealEngine.Ping()
			if err == nil {
				exist, err := engine.RealEngine.Table(tableName).Where(cond).Get(objPtr)
				if exist && err == nil {
					return true, nil
				} else if !exist && err == nil {
					return false, nil
				} else {
					return false, err
				}
			}
		} else {
			return false, fmt.Errorf("objPtr kind is not Ptr")
		}
	}
	return false, fmt.Errorf("engine is nil")
}

// GetRecordByCond2 根据查询条件来查询一条记录
func GetRecordByCond2(engine *DBEngineInfoS, tableName string, cond *builder.Cond, objPtr interface{}) (bool, error) {
	if engine != nil {
		objV := reflect.ValueOf(objPtr)
		if objV.Kind() == reflect.Ptr {
			err := engine.RealEngine.Ping()
			if err == nil {
				exist, err := engine.RealEngine.Table(tableName).Where(*cond).Get(objPtr)
				if exist && err == nil {
					return true, nil
				} else if !exist && err == nil {
					return false, nil
				} else {
					return false, err
				}
			}
		}
		return false, fmt.Errorf("objPtr kind is not Ptr")
	}
	return false, fmt.Errorf("engine is nil")
}

// FindRecordsBySimpleCond cond：where条件，如果传入空字符串代表没有查询条件，如果查询条件为空，limit也是无效的
// 只支持单一条件
// tableName：查询的表名
// limit：限制大小，如果不适用填写0
// start：limit的偏移量
// FindRecordsBySimpleCond FindRecordsBySimpleCond
func FindRecordsBySimpleCond(engine *DBEngineInfoS, tableName string, cond string, limit int, start int, resultSlicePtr interface{}) error {
	if engine != nil {
		resultV := reflect.ValueOf(resultSlicePtr)
		if resultV.Kind() == reflect.Ptr {
			err := engine.RealEngine.Ping()
			if err == nil {
				if limit > 0 {
					err = engine.RealEngine.Table(tableName).Where(cond).Limit(limit, start).Find(resultSlicePtr)
				} else {
					err = engine.RealEngine.Table(tableName).Where(cond).Find(resultSlicePtr)
				}

				return err
			}
		}
		return fmt.Errorf("resultSlicePtr kind is not Ptr")
	}

	return fmt.Errorf("%s", "engine is nil or cond is empty!")
}

// FindRecordsByMultiConds FindRecordsByMultiConds
func FindRecordsByMultiConds(engine *DBEngineInfoS, tableName string, cond *builder.Cond, limit int, start int, resultSlicePtr interface{}) error {
	if engine != nil {
		resultV := reflect.ValueOf(resultSlicePtr)
		if resultV.Kind() == reflect.Ptr {
			err := engine.RealEngine.Ping()
			if err == nil {
				session := engine.RealEngine.Table(tableName)
				if cond != nil {
					session = session.Where(*cond)
				}
				if limit > 0 {
					session = session.Limit(limit, start)
				}
				err = session.Find(resultSlicePtr)

				return err
			}
		} else {
			return errors.New("resultSlicePtr kind is not Ptr")
		}
	}
	return fmt.Errorf("%s", "engine is nil or cond is empty!")
}

// FindRecordsByMultiCondsDesc FindRecordsByMultiCondsDesc
func FindRecordsByMultiCondsDesc(engine *DBEngineInfoS, tableName string, cond *builder.Cond, limit int, start int, descfield string, resultSlicePtr interface{}) error {
	if engine != nil {
		resultV := reflect.ValueOf(resultSlicePtr)
		if resultV.Kind() == reflect.Ptr {
			err := engine.RealEngine.Ping()
			if err == nil {
				session := engine.RealEngine.Table(tableName)
				if cond != nil {
					session = session.Where(*cond)
				}
				session = session.Desc(descfield)
				if limit > 0 {
					session = session.Limit(limit, start)
				}
				err = session.Find(resultSlicePtr)

				return err
			}
		}

		return fmt.Errorf("resultSlicePtr kind is not Ptr")
	}
	return fmt.Errorf("%s", "engine is nil or cond is empty!")
}

// FindBillAmountByConds FindBillAmountByConds
func FindBillAmountByConds(engine *DBEngineInfoS, tableName string, cond *builder.Cond, resultSlicePtr interface{}) error {
	if engine != nil {
		resultV := reflect.ValueOf(resultSlicePtr)
		if resultV.Kind() == reflect.Ptr {
			err := engine.RealEngine.Ping()
			if err == nil {
				session := engine.RealEngine.Table(tableName)
				session.Select("sum(changeQuan) as changeQuan, payFlag")
				if cond != nil {
					session = session.Where(*cond)
				}
				session.GroupBy("payFlag")
				err = session.Find(resultSlicePtr)

				return err
			}
		} else {
			return fmt.Errorf("resultSlicePtr kind is not Ptr")
		}
	}
	return fmt.Errorf("%s", "engine is nil or cond is empty!")
}

// FindRecordsByMultiCondsCount FindRecordsByMultiCondsCount
func FindRecordsByMultiCondsCount(engine *DBEngineInfoS, tableName string, cond *builder.Cond, resultSlicePtr interface{}) (int64, error) {
	if engine != nil {
		resultV := reflect.ValueOf(resultSlicePtr)
		if resultV.Kind() == reflect.Ptr {
			err := engine.RealEngine.Ping()
			if err == nil {
				session := engine.RealEngine.Table(tableName)
				if cond != nil {
					session = session.Where(*cond)
				}
				count, err := session.Count(resultSlicePtr)
				return count, err
			}
		} else {
			return 0, fmt.Errorf("resultSlicePtr kind is not Ptr")
		}
	}
	return 0, fmt.Errorf("%s", "engine is nil or cond is empty!")
}

// FindDistinctRecordsByMultiConds FindDistinctRecordsByMultiConds
func FindDistinctRecordsByMultiConds(engine *DBEngineInfoS, tableName string, distinctColumns []string, cond *builder.Cond, limit int, start int, resultSlicePtr interface{}) error {
	if engine != nil {
		resultV := reflect.ValueOf(resultSlicePtr)
		if resultV.Kind() == reflect.Ptr {
			err := engine.RealEngine.Ping()
			if err == nil {
				session := engine.RealEngine.Table(tableName)
				if len(distinctColumns) > 0 {
					session = session.Distinct(distinctColumns...)
				}

				if cond != nil {
					session = session.Where(*cond)
				}
				if limit > 0 {
					session = session.Limit(limit, start)
				}

				err = session.Find(resultSlicePtr)
				return err
			}
		} else {
			return fmt.Errorf("resultSlicePtr kind is not Ptr")
		}
	}
	return fmt.Errorf("%s", "engine is nil or cond is empty!")
}

// FindRecordsBySimpleCondWithOrderBy  FindRecordsBySimpleCondWithOrderBy OrderBy("name desc")
func FindRecordsBySimpleCondWithOrderBy(engine *DBEngineInfoS, tableName string, cond string, limit int, start int, orderbyLst []string, resultSlicePtr interface{}) error {
	if engine != nil {
		resultV := reflect.ValueOf(resultSlicePtr)
		if resultV.Kind() == reflect.Ptr {
			err := engine.RealEngine.Ping()
			if err == nil {
				session := engine.RealEngine.Table(tableName).Where(cond)
				// 设置orderby
				if len(orderbyLst) > 0 {
					for _, orderbyContent := range orderbyLst {
						session = session.OrderBy(orderbyContent)
					}
				}
				if limit > 0 {
					session = session.Limit(limit, start)
				}
				err = session.Find(resultSlicePtr)
				return err
			}
		} else {
			return fmt.Errorf("resultSlicePtr kind is not Ptr")
		}
	}
	return fmt.Errorf("%s", "engine is nil or cond is empty!")
}

// QueryMultiConds QueryMultiConds
type QueryMultiConds struct {
	Engine         *DBEngineInfoS
	TableName      string
	Cond           *builder.Cond
	QueryFields    []string
	Limit          int
	StartPos       int
	OrderBys       []string
	ResultSlicePtr interface{}
}

// FindRecordsByMultiCondsStruct FindRecordsByMultiCondsStruct
func FindRecordsByMultiCondsStruct(multiCond QueryMultiConds) error {
	if multiCond.Engine == nil {
		return fmt.Errorf("FindRecordsByMultiCondsStruct engine is nil")
	}
	if multiCond.TableName == "" {
		return fmt.Errorf("FindRecordsByMultiCondsStruct table is empty")
	}
	resultV := reflect.ValueOf(multiCond.ResultSlicePtr)
	if resultV.Kind() != reflect.Ptr {
		return fmt.Errorf("FindRecordsByMultiCondsStruct resultSlicePtr kind is not Ptr")
	}

	err := multiCond.Engine.RealEngine.Ping()
	if err != nil {
		return fmt.Errorf("FindRecordsByMultiCondsStruct Ping error:%s", err)
	}
	session := multiCond.Engine.RealEngine.Table(multiCond.TableName)
	if multiCond.Cond != nil {
		session = session.Where(*multiCond.Cond)
	}
	// 设置orderby
	if len(multiCond.OrderBys) > 0 {
		for _, orderbyContent := range multiCond.OrderBys {
			session = session.OrderBy(orderbyContent)
		}
	}
	if len(multiCond.QueryFields) > 0 {
		session.Cols(multiCond.QueryFields...)
	}
	if multiCond.Limit > 0 {
		session = session.Limit(multiCond.Limit, multiCond.StartPos)
	}
	err = session.Find(multiCond.ResultSlicePtr)
	return err
}

// DeleteRecordsByMultiConds DeleteRecordsByMultiConds
func DeleteRecordsByMultiConds(engine *DBEngineInfoS, tableName string, cond *builder.Cond, objPtr interface{}) (affectedRows int64, err error) {
	affectedRows = 0
	err = nil

	if engine != nil && len(tableName) > 0 {
		err = engine.RealEngine.Ping()
		if err == nil {
			affectedRows, err = engine.RealEngine.Table(tableName).Where(*cond).Delete(objPtr)
			return
		}
	}
	return 0, fmt.Errorf("engine is nil or tableName is empty")
}

// DeleteRecord DeleteRecord
func DeleteRecord(engine *DBEngineInfoS, tableName string, primaryKeys *core.PK, objPtr interface{}) (affectedRows int64, err error) {
	affectedRows = 0
	err = nil

	if engine != nil && len(tableName) > 0 {
		err = engine.RealEngine.Ping()
		if err == nil {
			affectedRows, err = engine.RealEngine.Table(tableName).Id(primaryKeys).Delete(objPtr)
			return
		}
	}
	return 0, fmt.Errorf("engine is nil or tableName is empty")
}

// SelectRecordsByCond SelectRecordsByCond
func SelectRecordsByCond(engine *DBEngineInfoS, tableName string, cond string, record interface{}) ([]interface{}, error) {
	if engine != nil && len(tableName) > 0 {
		recordT := reflect.Indirect(reflect.ValueOf(record)).Type()
		fmt.Println("recordT", recordT.String())
		// 这个不带package名字
		fmt.Println("recordT", recordT.Name())
		// // 动态创建对象
		rows, err := engine.RealEngine.Table(tableName).Where(cond).Rows(record)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		// 结果集
		result := make([]interface{}, 0)

		for rows.Next() {
			newRecord := reflect.New(recordT)
			err = rows.Scan(newRecord.Interface())
			if err != nil {
				return nil, err
			}

			result = append(result, newRecord.Interface())
		}
		return result, nil

	}
	return nil, fmt.Errorf("engine is nil or tableName is empty")
}

// SelectRecordsByCond2 SelectRecordsByCond2
func SelectRecordsByCond2(engine *DBEngineInfoS, tableName string, cond string, tblStructName string) ([]interface{}, error) {
	if engine != nil && len(tableName) > 0 {
		record, err := NewTableObj(tblStructName)
		if err != nil {
			return nil, err
		}

		rows, err := engine.RealEngine.Table(tableName).Where(cond).Rows(record)
		if err != nil {
			return nil, err
		}

		defer rows.Close()
		// 结果集
		result := make([]interface{}, 0)

		for rows.Next() {
			newRecord, _ := NewTableObj(tblStructName)
			err = rows.Scan(newRecord)
			if err != nil {
				return nil, err
			}

			result = append(result, newRecord)
		}
		return result, nil
	}
	return nil, fmt.Errorf("engine is nil or tableName is empty")
}

// FindRecordsByMultiCondsOrderBy 排序  可传参   sortField 排序的字段   sortStatus 排序类型  desc asc
func FindRecordsByMultiCondsOrderBy(lc *logging.LogContext, engine *DBEngineInfoS, tableName string, cond *builder.Cond, limit int, start int, sortField string, sortStatus string, resultSlicePtr interface{}) error {
	if engine != nil {
		resultV := reflect.ValueOf(resultSlicePtr)
		if resultV.Kind() == reflect.Ptr {
			err := engine.RealEngine.Ping()
			if err == nil {
				session := engine.RealEngine.Table(tableName)
				if cond != nil {
					session = session.Where(*cond)
				}
				if len(sortField) > 0 && sortStatus != "" {

					logging.Info(lc, "sort sortStatus: %s  sortField:  ", sortStatus, sortField)
					if strings.Compare("asc", sortStatus) == 0 {
						logging.Info(lc, "sort    asc")

						session.Asc(sortField)
					} else if strings.Compare("desc", sortStatus) == 0 {
						logging.Info(lc, "sort   desc")

						session.Desc(sortField)
					}
				} else {
					session.Desc("create_time")

				}
				if limit > 0 {
					session = session.Limit(limit, start)

				}
				err = session.Find(resultSlicePtr)

				return err
			}
		}
		return errors.New("resultSlicePtr kind is not Ptr")

	}
	return fmt.Errorf("%s", "engine is nil or cond is empty!")
}
