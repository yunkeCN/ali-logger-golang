package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	zerologger "github.com/rs/zerolog/log"
	"log"
	"os"
	"path/filepath"
)

var (
	accessLogger   zerolog.Logger
	businessLogger zerolog.Logger
	errorLogger    zerolog.Logger
)

var defaultLogger = WithTag("")

type Options struct {
	IsDev       bool
	ProjectName string
}

var optionsInner Options

func getLogger(filePath string) zerolog.Logger {
	var basePath string
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
		os.Exit(-1)
	}
	basePath = dir + "/log/"

	if !optionsInner.IsDev {
		basePath = "/var/log/service/"
	}

	s := basePath + filePath

	err = os.MkdirAll(filepath.Dir(s), 0777)
	if err != nil {
		log.Fatalln(err)
		os.Exit(-1)
	}

	f, err := os.OpenFile(s, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)
		os.Exit(-1)
	}

	return zerologger.Output(f)
}

func Init(options Options) {
	optionsInner = options
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerologger.Logger = zerologger.With().Logger()
	if optionsInner.IsDev {
		accessLogger = zerologger.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		businessLogger = zerologger.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		errorLogger = zerologger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	} else {
		accessLogger = getLogger(fmt.Sprintf("access/%s/out.log", optionsInner.ProjectName))
		businessLogger = getLogger(fmt.Sprintf("business/%s/out.log", optionsInner.ProjectName))
		errorLogger = getLogger(fmt.Sprintf("err/%s/error.log", optionsInner.ProjectName))
	}

	zerolog.MessageFieldName = "msg"
}

func Access(msg string) {
	defaultLogger.Access(msg)
}

func Accessf(format string, v ...interface{}) {
	defaultLogger.Accessf(format, v...)
}

func Business(msg string) {
	defaultLogger.Business(msg)
}

func Businessf(format string, v ...interface{}) {
	defaultLogger.Businessf(format, v...)
}

func Error(msg string) {
	defaultLogger.Error(msg)
}

func Errorf(format string, v ...interface{}) {
	defaultLogger.Errorf(format, v...)
}

//Deprecated:请使用WithTag代替
func WithPrefix(prefix string) Logger {
	return WithTag(prefix)
}

func WithTag(tag string) Logger {
	return Logger{accessLogger: &accessLogger, businessLogger: &businessLogger, errorLogger: &errorLogger, tag: tag}
}
