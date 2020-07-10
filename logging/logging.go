package logging

import "fmt"

// LogContext LogContext
type LogContext struct {
	//本次会话的追踪ID
	TraceID string
	//业务的ID,例如订单ID ,orderid （最大16位)
	Callid string
}

var (
	// StartupLogContext StartupLogContext
	StartupLogContext LogContext
)

// InitStartupLogContext InitStartupLogContext
func InitStartupLogContext(serviceName string) *LogContext {
	var lc LogContext
	lc.TraceID = fmt.Sprintf("%s_startup", serviceName)
	return &lc
}
