package mysql

import "fmt"

// 直接执行sql命令

// SqlExec SqlExec
func SqlExec(engine *DBEngineInfoS, sqlContent string) (affectedRows, lastInsertID int64, err error) {
	affectedRows = 0
	lastInsertID = -1
	err = nil

	if len(sqlContent) > 0 && engine != nil {
		result, res := engine.RealEngine.Exec(sqlContent)
		if res != nil {
			return
		}
		affectedRows, _ = result.RowsAffected()
		lastInsertID, _ = result.LastInsertId()
	} else {
		err = fmt.Errorf("%s", "Params engine or sqlContent is invalid!")
	}

	return
}
