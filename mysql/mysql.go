package mysql

import (
	"errors"
	"fmt"
	"github.com/haleyrom/storage/logging"
	"time"
)

// MysqlManager MysqlManager
type MysqlManager struct {
	EngineMap map[string]*DBEngineInfoS
}

const (
	// MAX_DB_CONNECTION_COUNT MAX_DB_CONNECTION_COUNT
	MAX_DB_CONNECTION_COUNT int = 20
	// MYSQL_PING MYSQL_PING
	MYSQL_PING int = 300
)

var (
	// GMysqlManager GMysqlManager
	GMysqlManager *MysqlManager = new(MysqlManager)
)

// GetMysql GetMysql
func (mysqlMgr *MysqlManager) GetMysql(dbName string) (*DBEngineInfoS, error) {
	if len(dbName) == 0 {
		return nil, errors.New("parameters dbname is empty")
	}

	if dbEngine, ok := mysqlMgr.EngineMap[dbName]; ok {
		return dbEngine, nil
	}

	return nil, errors.New("database is not exist")
}

// Initialize Initialize
func (mysqlMgr *MysqlManager) Initialize(lc *logging.LogContext, key string, Host string, Port int, User string, Password string, dbName string) error {
	logging.Debug(lc, "MysqlManager Initialize enter")

	mysqlMgr.EngineMap = make(map[string]*DBEngineInfoS)

	addr := fmt.Sprintf("%s:%d", Host, Port)
	dbEngine, err := CreateDBEngnine(lc, "mysql", User, Password, addr, dbName)
	if err != nil {
		logging.Error(lc, "CreateDBEngnine failed reason[%s]", err.Error())
		return err
	}

	SetDBEngineConnectionParams(lc, dbEngine, MAX_DB_CONNECTION_COUNT, MAX_DB_CONNECTION_COUNT)

	DoDBKeepAlive(lc, dbEngine, time.Duration(MYSQL_PING)*time.Second)

	mysqlMgr.EngineMap[key] = dbEngine

	return nil
}
