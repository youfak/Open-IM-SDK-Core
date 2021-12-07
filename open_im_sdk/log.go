package open_im_sdk

import (
	"bufio"
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var logger *Logger

type Logger struct {
	*logrus.Logger
	Pid int
}

//func init() {
//	logger = loggerInit("")
//
//}
func NewPrivateLog(moduleName string) {
	logger = loggerInit(moduleName)
}
func IsNil() bool {
	if logger != nil {
		return false
	}
	return true
}

func loggerInit(moduleName string) *Logger {
	var logger = logrus.New()
	//All logs will be printed
	logger.SetLevel(logrus.Level(6))
	//Close std console output
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(err.Error())
	}
	writer := bufio.NewWriter(src)
	logger.SetOutput(writer)
	//Log Console Print Style Setting
	logger.SetFormatter(&nested.Formatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		HideKeys:        false,
		FieldsOrder:     []string{"PID", "FilePath", "OperationID"},
	})
	//File name and line number display hook
	logger.AddHook(newFileHook())

	//Log file segmentation hook
	hook := NewLfsHook(time.Duration(24)*time.Hour, 5, moduleName)
	logger.AddHook(hook)
	return &Logger{
		logger,
		os.Getpid(),
	}
}
func NewLfsHook(rotationTime time.Duration, maxRemainNum uint, moduleName string) logrus.Hook {
	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: initRotateLogs(rotationTime, maxRemainNum, "debug", moduleName),
		logrus.InfoLevel:  initRotateLogs(rotationTime, maxRemainNum, "info", moduleName),
		logrus.WarnLevel:  initRotateLogs(rotationTime, maxRemainNum, "warn", moduleName),
		logrus.ErrorLevel: initRotateLogs(rotationTime, maxRemainNum, "error", moduleName),
	}, &nested.Formatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		HideKeys:        false,
		FieldsOrder:     []string{"PID", "FilePath", "OperationID"},
	})
	return lfsHook
}
func initRotateLogs(rotationTime time.Duration, maxRemainNum uint, level string, moduleName string) *rotatelogs.RotateLogs {
	if moduleName != "" {
		moduleName = moduleName + "."
	}
	writer, err := rotatelogs.New(
		"../logs/"+moduleName+level+"."+"%Y-%m-%d",
		rotatelogs.WithRotationTime(rotationTime),
		rotatelogs.WithRotationCount(maxRemainNum),
	)
	if err != nil {
		panic(err.Error())
	} else {
		return writer
	}
}

//Deprecated
func Info(token, OperationID, format string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"PID":         logger.Pid,
		"OperationID": OperationID,
	}).Infof(format, args...)

}

//Deprecated
func Error(token, OperationID, format string, args ...interface{}) {

	logger.WithFields(logrus.Fields{
		"PID":         logger.Pid,
		"OperationID": OperationID,
	}).Errorf(format, args...)

}

//Deprecated
func Debug(token, OperationID, format string, args ...interface{}) {

	logger.WithFields(logrus.Fields{
		"PID":         logger.Pid,
		"OperationID": OperationID,
	}).Debugf(format, args...)

}

//Deprecated
func Warning(token, OperationID, format string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"PID":         logger.Pid,
		"OperationID": OperationID,
	}).Warningf(format, args...)

}

//Deprecated
func InfoByArgs(format string, args ...interface{}) {
	logger.WithFields(logrus.Fields{}).Infof(format, args)
}

//Deprecated
func ErrorByArgs(format string, args ...interface{}) {
	logger.WithFields(logrus.Fields{}).Errorf(format, args...)
}

//Print log information in k, v format,
//kv is best to appear in pairs. tipInfo is the log prompt information for printing,
//and kv is the key and value for printing.
//Deprecated
func InfoByKv(tipInfo, OperationID string, args ...interface{}) {
	fields := make(logrus.Fields)
	argsHandle(OperationID, fields, args)
	logger.WithFields(fields).Info(tipInfo)
}

//Deprecated
func ErrorByKv(tipInfo, OperationID string, args ...interface{}) {
	fields := make(logrus.Fields)
	argsHandle(OperationID, fields, args)
	logger.WithFields(fields).Error(tipInfo)
}

//Deprecated
func DebugByKv(tipInfo, OperationID string, args ...interface{}) {
	fields := make(logrus.Fields)
	argsHandle(OperationID, fields, args)
	logger.WithFields(fields).Debug(tipInfo)
}

//Deprecated
func WarnByKv(tipInfo, OperationID string, args ...interface{}) {
	fields := make(logrus.Fields)
	argsHandle(OperationID, fields, args)
	logger.WithFields(fields).Warn(tipInfo)
}

//internal method
func argsHandle(OperationID string, fields logrus.Fields, args []interface{}) {
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			fields[fmt.Sprintf("%v", args[i])] = args[i+1]
		} else {
			fields[fmt.Sprintf("%v", args[i])] = ""
		}
	}
	fields["OperationID"] = OperationID
	fields["PID"] = logger.Pid
}
func NewInfo(OperationID string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"OperationID": OperationID,
		"PID":         logger.Pid,
	}).Infoln(args)
}
func NewError(OperationID string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"OperationID": OperationID,
		"PID":         logger.Pid,
	}).Errorln(args)
}
func NewDebug(OperationID string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"OperationID": OperationID,
		"PID":         logger.Pid,
	}).Debugln(args)
}
func NewWarn(OperationID string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"OperationID": OperationID,
		"PID":         logger.Pid,
	}).Warnln(args)
}