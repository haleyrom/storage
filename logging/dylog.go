package logging

import (
	"fmt"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"runtime"

	"sync"
	"time"
)

const (
	// LevelDebug LevelDebug
	LevelDebug = iota
	// LevelTrace LevelTrace
	LevelTrace
	// LevelInfo LevelInfo
	LevelInfo
	// LevelWarn LevelWarn
	LevelWarn
	// LevelError LevelError
	LevelError
	// LevelFatal LevelFatal
	LevelFatal
)

var (
	// LeverStr LeverStr
	LeverStr = map[int]string{
		0: "Debug",
		1: "Trace",
		2: "Info",
		3: "Warn",
		4: "Error",
		5: "Fatal",
	}

	// dylogger dylogger
	dylogger *Dylogger

	// LogErrorFuncCallback LogErrorFuncCallback
	LogErrorFuncCallback ErrorFuncCallback
)

// ErrorFuncCallback ErrorFuncCallback
type ErrorFuncCallback func(string)

// Debug Debug
func Debug(lc *LogContext, template string, v ...interface{}) {
	if LevelDebug < dylogger.level {
		return
	}
	dylogger.writeMsg(lc, LevelDebug, template, v...)
}

// Trace Trace
func Trace(lc *LogContext, template string, v ...interface{}) {
	if LevelTrace < dylogger.level {
		return
	}
	dylogger.writeMsg(lc, LevelTrace, template, v...)
}

// Info Info
func Info(lc *LogContext, template string, v ...interface{}) {

	if LevelInfo < dylogger.level {
		return
	}
	dylogger.writeMsg(lc, LevelInfo, template, v...)
}

// Notice Notice
func Notice(lc *LogContext, template string, v ...interface{}) {
	if LevelInfo < dylogger.level {
		return
	}
	dylogger.writeMsg(lc, LevelTrace, template, v...)
}

// Warn Warn
func Warn(lc *LogContext, template string, v ...interface{}) {

	if LevelWarn < dylogger.level {
		return
	}
	dylogger.writeMsg(lc, LevelWarn, template, v...)
}

// Error Error
func Error(lc *LogContext, template string, v ...interface{}) {

	if LevelError < dylogger.level {
		return
	}
	dylogger.writeMsg(lc, LevelError, template, v...)
}

// Fatal Fatal
func Fatal(lc *LogContext, template string, v ...interface{}) {

	if LevelFatal < dylogger.level {
		return
	}
	dylogger.writeMsg(lc, LevelFatal, template, v...)
}

// Fatal2 Fatal2
func Fatal2(lc *LogContext, template string, v ...interface{}) {

	if LevelFatal < dylogger.level {
		return
	}
	dylogger.writeMsg(lc, LevelFatal, template, v...)
}

// Dylogger Dylogger
type Dylogger struct {
	//锁
	lock sync.Mutex
	//是否已经初始化
	init bool
	//日志级别
	level int
	//调用栈的深度
	loggerFuncCallDepth int
	//是否打印調用棧
	printCallStack bool

	lockWhileMap sync.Mutex
	//打印日志的代码文件的白名单
	printLogWhiteMap map[string]int
	//输出日志的路径
	logBasePath string
	//非颜色目录日志路径
	logNoclorPath string
	// 是否已经设置设置日志目录
	logSetPath bool
	//日志文件名
	logFileName string
	//输出日志 到控制台的最大行数
	maxLogConsolLineLimitCount int
	//当前输出的日志到控制台的行数
	curLogConsolLine int
}

// Loginit Loginit
func Loginit(callDepth int, printStack bool) {
	if callDepth < 3 {
		os.Stderr.WriteString("callDepth can't <3")
	}

	dylogger = NewDylogger()
	dylogger.loggerFuncCallDepth = callDepth
	dylogger.printCallStack = printStack
	dylogger.logBasePath = "."
	dylogger.logSetPath = false
	dylogger.maxLogConsolLineLimitCount = 30000

	/*
		//建立个非颜色版本的目录
		dylogger.logNoclorPath = fmt.Sprintf("%s/nocolor", dylogger.logBasePath)
		dylogger.checkAndCreateNoColorPath()
	*/

}

// init init
func init() {
	Loginit(3, false)
}

// SetLogPath SetLogPath
func SetLogPath(path string, logfileName string) error {

	dylogger.logSetPath = true
	dylogger.logFileName = logfileName

	dylogger.logBasePath = path
	if path == "" {
		dylogger.logBasePath = "."
	}

	ok, err := PathExists(dylogger.logBasePath)
	if !ok {
		err = os.Mkdir(dylogger.logBasePath, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir [%s] failed[%s]\n", dylogger.logBasePath, err)
			return err
		}
	}

	//建立个非颜色版本的目录
	dylogger.logNoclorPath = fmt.Sprintf("%s/nocolor", dylogger.logBasePath)
	dylogger.checkAndCreateNoColorPath()

	return nil
}

// GetLogPath GetLogPath
func GetLogPath() string {
	return dylogger.logBasePath
}

// SetLogLevel SetLogLevel
func SetLogLevel(level int) {
	dylogger.level = level

}

// NewDylogger NewDylogger
func NewDylogger() *Dylogger {
	logger := new(Dylogger)
	logger.init = true
	logger.level = LevelDebug
	logger.loggerFuncCallDepth = 3
	logger.printLogWhiteMap = make(map[string]int)

	return logger
}

// Zlog 适配zlog 的日志输出使用这个日志打屏
func Zlog(lvl zapcore.Level, msg string) {

	if lvl == zapcore.DebugLevel {

		dylogger.writeMsg2(LevelDebug, msg)
	} else if lvl == zapcore.InfoLevel {
		dylogger.writeMsg2(LevelInfo, msg)
	} else if lvl == zapcore.WarnLevel {
		dylogger.writeMsg2(LevelWarn, msg)
	} else if lvl == zapcore.ErrorLevel {
		dylogger.writeMsg2(LevelError, msg)

	} else if lvl == zapcore.DPanicLevel {
		dylogger.writeMsg2(LevelFatal, msg)

	} else if lvl == zapcore.PanicLevel {
		dylogger.writeMsg2(LevelFatal, msg)

	} else if lvl == zapcore.FatalLevel {
		dylogger.writeMsg2(LevelFatal, msg)
	}
}

//获取时间的前缀
func (log *Dylogger) getTimePrefix(t *time.Time) string {
	strTime := fmt.Sprintf("%02d-%02d %02d:%02d:%02d", t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	return strTime
}

//获取日志级别的字符串
func (log *Dylogger) getStrLevel(logLevel int) string {
	if logLevel == LevelDebug {
		return "[DEBUG]"
	} else if logLevel == LevelTrace {
		return "[TRACE]"
	} else if logLevel == LevelInfo {
		return "[INFO ]"
	} else if logLevel == LevelWarn {
		return "[WARN ]"
	} else if logLevel == LevelError {
		return "[ERROR]"
	} else if logLevel == LevelFatal {
		return "[FATAL]"
	} else {
		return "[WHAT ]"
	}
}

/*
char COLOR_DEF[][32] =
{
"\033[00;32m",//GREEN  DEBUG
"\033[00;35m",//BLUE TRACE
"\033[01;36m",//CYAN INFO
"\033[01;33m",//BOLD YELLOW WARN
"\033[01;31m", // # BOLD RED
"\033[00;31m"  //# BOLD RED
};
#define  CLOSE_COLOR "\033[0m"
*/

// CLOSE_COLOR CLOSE_COLOR
var CLOSE_COLOR string = "\033[0m"

// getLevelColor 获取日志级别的颜色
func (log *Dylogger) getLevelColor(logLevel int) string {
	if logLevel == LevelDebug {
		return "\033[00;32m"
	} else if logLevel == LevelTrace {
		return "\033[00;35m"
	} else if logLevel == LevelInfo {
		return "\033[01;36m"
	} else if logLevel == LevelWarn {
		return "\033[01;33m"
	} else if logLevel == LevelError {
		return "\033[01;31m"
	} else if logLevel == LevelFatal {
		return "\033[01;31m"
	} else {
		return ""
	}
}

//获取调用栈：
func (log *Dylogger) getCaller() (file string, line int) {

	// fmt.Printf("loggerFuncCallDepth=%d\n",log.loggerFuncCallDepth)
	_, file, line, ok := runtime.Caller(log.loggerFuncCallDepth)
	if !ok {
		file = "???"
		line = 0
	}
	_, filename := path.Split(file)
	return filename, line

}

func (log *Dylogger) printStackStace() {

	//--------------------------
	log.lock.Lock()
	for i := 0; i < 10; i++ {

		_, file2, line2, ok2 := runtime.Caller(i)
		if !ok2 {
			file2 = "???"
			line2 = 0
		}
		fmt.Printf("i=%d %s  %d\n", i, file2, line2)
	}
	log.lock.Unlock()

	//-----------------

}

func (log *Dylogger) getCallerPrefix(filename string, line int) string {

	return fmt.Sprintf("F[%30s] L[%5d]", filename, line)

}

func (log *Dylogger) getLogContext(ctx *LogContext) string {
	if ctx == nil {
		//没传日志的上下文
		return fmt.Sprintf("T[%30s] C[%16s]", "", "")
	}
	return fmt.Sprintf("T[%30s] C[%16s]", ctx.TraceID, ctx.Callid)
}

// getLogPrefix 获取日志的前缀
func (log *Dylogger) getLogPrefix(logLevel int, t *time.Time) string {

	strPrefix := fmt.Sprintf("%s %s", log.getStrLevel(logLevel), log.getTimePrefix(t))

	return strPrefix
}

// AddPrintLogCodeFileWhiteList 设置代码文件名为本文件的白名单才打印
func AddPrintLogCodeFileWhiteList(codeFile string, level int) {

	dylogger.lockWhileMap.Lock()
	dylogger.printLogWhiteMap[codeFile] = level
	dylogger.lockWhileMap.Unlock()

}

//检查这个代码文件是否可打印日志
func (log *Dylogger) isCanPrintLogFilter(codeFile string) bool {

	log.lockWhileMap.Lock()
	defer log.lockWhileMap.Unlock()

	if len(log.printLogWhiteMap) == 0 {
		//没启用白名单的机制，全输出日志
		return true
	}

	_, ok := log.printLogWhiteMap[codeFile]

	return ok
}

// PrintLogstack 尝试打印调用堆栈，确定调用栈的深度
func PrintLogstack() {
	for i := 0; i < 10; i++ {

		_, file, line, ok := runtime.Caller(i)
		if !ok {
			file = "???"
			line = 0
		}
		fmt.Printf("i=%d %s  %d\n", i, file, line)
	}
}

func (log *Dylogger) writeMsg(lc *LogContext, logLevel int, template string, fmtArgs ...interface{}) {

	if logLevel < log.level {
		//小于了日志级别的限制，不输出日志
		return
	}

	//---------------------------
	msg := template
	if msg == "" && len(fmtArgs) > 0 {
		msg = fmt.Sprint(fmtArgs...)
	} else if msg != "" && len(fmtArgs) > 0 {
		msg = fmt.Sprintf(template, fmtArgs...)
	}
	//-----------------------------

	codefile, codeline := log.getCaller()
	//调用日志过滤器判断
	if !log.isCanPrintLogFilter(codefile) {
		return
	}
	if log.printCallStack {
		log.printStackStace()
	}

	curTime := time.Now()

	//存一份日志到kafka
	/*
		tcplog := new(TcpLog)
		tcplog.CodeFile = codefile
		tcplog.CodeLine = codeline
		if lc != nil {
			tcplog.LogCallId = lc.Callid
			tcplog.LogTraceId = lc.TraceID
		}
		tcplog.LogeLevel = logLevel
		tcplog.LogTime = curTime.Format("2006-01-02 15:04:05")
		tcplog.LogMsg = msg
		GlobalTcpLogMgr.Push(tcplog)
	*/
	toLogMsg := fmt.Sprintf("%s %s %s %s %s %s\n",
		log.getLevelColor(logLevel),
		log.getLogPrefix(logLevel, &curTime),
		//調用日誌的代碼
		log.getCallerPrefix(codefile, codeline),
		log.getLogContext(lc),

		msg,
		CLOSE_COLOR)

	/*
		noColorMsg := fmt.Sprintf("%s %s %s%s\n",
			log.getLogPrefix(logLevel, &curTime),
			//調用日誌的代碼
			log.getCallerPrefix(codefile, codeline),

			msg,
			CLOSE_COLOR)*/

	//将日志打印到屏幕
	log.lock.Lock()
	//fmt.Printf(toLogMsg)
	//os.Stderr.WriteString(toLogMsg)

	os.Stderr.WriteString(toLogMsg)

	if logLevel >= LevelError {
		if LogErrorFuncCallback != nil {

			noColorMsg1 := fmt.Sprintf("%s %s %s %s\n",
				log.getLogPrefix(logLevel, &curTime),
				//調用日誌的代碼
				log.getCallerPrefix(codefile, codeline),
				log.getLogContext(lc),
				msg)

			LogErrorFuncCallback(noColorMsg1)
		}
	}

	//if logLevel >= LevelError {
	//先暂时关下日志，降低IO压力  ternence
	if log.logSetPath {
		log.writeMsgToFile(toLogMsg, codefile)
		//	log.writeMsgToFileNoColor(noColorMsg, codefile)
		log.writeMsgToFile(toLogMsg, log.logFileName)
	}
	//	}

	log.lock.Unlock()

}

//沒有模板的直接
func (log *Dylogger) writeMsg2(logLevel int, msg string) {

	if logLevel < log.level {
		//小于了日志级别的限制，不输出日志
		return
	}

	codefile, codeline := log.getCaller()
	//调用日志过滤器判断
	if !log.isCanPrintLogFilter(codefile) {
		return
	}
	if log.printCallStack {
		log.printStackStace()
	}

	curTime := time.Now()

	toLogMsg := fmt.Sprintf("%s%s %s %s%s\n",
		log.getLevelColor(logLevel),
		log.getLogPrefix(logLevel, &curTime),
		//調用日誌的代碼
		log.getCallerPrefix(codefile, codeline),

		msg,
		CLOSE_COLOR)

	//将日志打印到屏幕
	log.lock.Lock()
	fmt.Printf(toLogMsg)

	log.writeMsgToFile(toLogMsg, codefile)
	log.writeMsgToFileNoColor(toLogMsg, codefile)
	log.writeMsgToFile(toLogMsg, log.logFileName)

	log.lock.Unlock()

}
