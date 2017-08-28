package log

import (
    l4g "github.com/ms2008/log4go"
)

var lg = make(l4g.Logger)

func init() {
    clw := l4g.NewConsoleLogWriter()
    clw.SetFormat("[%D %T] [%L] %M")
    //lg.AddFilter("stdout", l4g.DEBUG, l4g.NewConsoleLogWriter())
    lg.AddFilter("stdout", l4g.DEBUG, clw)

    flw := l4g.NewFileLogWriter("connect.log", false)
    flw.SetFormat("[%D %T] [%L] %M")
    lg.AddFilter("file", l4g.CRITICAL, flw)
}

func Close() {
    lg.Close()
}

func Finest(arg0 interface{}, args ...interface{}) {
    lg.Finest(arg0, args...)
}

func Fine(arg0 interface{}, args ...interface{}) {
    lg.Fine(arg0, args...)
}

func Debug(arg0 interface{}, args ...interface{}) {
    lg.Debug(arg0, args...)
}

func Trace(arg0 interface{}, args ...interface{}) {
    lg.Trace(arg0, args...)
}

func Info(arg0 interface{}, args ...interface{}) {
    lg.Info(arg0, args...)
}

func Error(arg0 interface{}, args ...interface{}) {
    lg.Error(arg0, args...)
}

func Critical(arg0 interface{}, args ...interface{}) {
    lg.Critical(arg0, args...)
}