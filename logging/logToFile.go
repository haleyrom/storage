package logging

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// PathExists PathExists
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// checkFileIsExist checkFileIsExist
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// checkAndCreateNoColorPath checkAndCreateNoColorPath
func (log *Dylogger) checkAndCreateNoColorPath() error {

	if checkFileIsExist(log.logNoclorPath) {
		return nil
	}

	err := os.Mkdir(log.logNoclorPath, os.ModePerm)
	if err != nil {
		fmt.Printf("mkdir [%s] failed[%s]\n", log.logNoclorPath, err)
	}

	return err
}

// getLogfieName getLogfieName
func (log *Dylogger) getLogfieName(destLogFile string) string {
	strarray := strings.Split(destLogFile, ".")

	curTime := time.Now()

	strFullPath := fmt.Sprintf("%s/%s_%02d%02d.log", log.logBasePath, strarray[0], curTime.Month(), curTime.Day())
	return strFullPath
}
func (log *Dylogger) getLogfieNameNoColor(destLogFile string) string {
	strarray := strings.Split(destLogFile, ".")

	curTime := time.Now()

	strFullPath := fmt.Sprintf("%s/%s_%02d%02d.log", log.logNoclorPath, strarray[0], curTime.Month(), curTime.Day())
	return strFullPath
}

// writeMsgToFile writeMsgToFile
func (log *Dylogger) writeMsgToFile(wireteString string, destLogFile string) {
	if destLogFile == "" {
		return
	}
	var filename = log.getLogfieName(destLogFile)

	fd, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

	buf := []byte(wireteString)
	fd.Write(buf)
	fd.Close()

}

// writeMsgToFileNoColor writeMsgToFileNoColor
func (log *Dylogger) writeMsgToFileNoColor(wireteString string, destLogFile string) {

	var filename = log.getLogfieNameNoColor(destLogFile)

	fd, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

	buf := []byte(wireteString)
	fd.Write(buf)
	fd.Close()

}
