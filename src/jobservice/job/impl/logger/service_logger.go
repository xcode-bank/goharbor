package logger

import (
	"os"

	"github.com/vmware/harbor/src/common/utils/log"
)

//ServiceLogger is an implementation of logger.Interface.
//It used to log info in workerpool components.
type ServiceLogger struct {
	backendLogger *log.Logger
}

//NewServiceLogger to create new logger for job service
//nil might be returned
func NewServiceLogger(level string) *ServiceLogger {
	logLevel := parseLevel(level)
	backendLogger := log.New(os.Stdout, log.NewTextFormatter(), logLevel)

	return &ServiceLogger{
		backendLogger: backendLogger,
	}
}

//Debug ...
func (sl *ServiceLogger) Debug(v ...interface{}) {
	sl.backendLogger.Debug(v...)
}

//Debugf with format
func (sl *ServiceLogger) Debugf(format string, v ...interface{}) {
	sl.backendLogger.Debugf(format, v...)
}

//Info ...
func (sl *ServiceLogger) Info(v ...interface{}) {
	sl.backendLogger.Info(v...)
}

//Infof with format
func (sl *ServiceLogger) Infof(format string, v ...interface{}) {
	sl.backendLogger.Infof(format, v...)
}

//Warning ...
func (sl *ServiceLogger) Warning(v ...interface{}) {
	sl.backendLogger.Warning(v...)
}

//Warningf with format
func (sl *ServiceLogger) Warningf(format string, v ...interface{}) {
	sl.backendLogger.Warningf(format, v...)
}

//Error ...
func (sl *ServiceLogger) Error(v ...interface{}) {
	sl.backendLogger.Error(v...)
}

//Errorf with format
func (sl *ServiceLogger) Errorf(format string, v ...interface{}) {
	sl.backendLogger.Errorf(format, v...)
}

//Fatal error
func (sl *ServiceLogger) Fatal(v ...interface{}) {
	sl.backendLogger.Fatal(v...)
}

//Fatalf error
func (sl *ServiceLogger) Fatalf(format string, v ...interface{}) {
	sl.backendLogger.Fatalf(format, v...)
}
