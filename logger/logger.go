package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	zerologger "github.com/rs/zerolog/log"
	"io"
	"log"
	"os"
	"path/filepath"
)

var (
	accessLogger   zerolog.Logger
	businessLogger zerolog.Logger
	errorLogger    zerolog.Logger
	AccessWriter   io.Writer
	BusinessWriter io.Writer
	ErrorWriter    io.Writer
)

type Options struct {
	IsDev       bool
	ProjectName string
}

var optionsInner Options

var isDevInner bool

func Init(options Options) {
	optionsInner = options
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if optionsInner.IsDev {
		accessLogger = zerologger.Level(zerolog.GlobalLevel())
		businessLogger = zerologger.Level(zerolog.GlobalLevel())
		errorLogger = zerologger.Level(zerolog.GlobalLevel())
		AccessWriter = os.Stdout
		BusinessWriter = os.Stdout
		ErrorWriter = os.Stderr
	} else {
		accessLogger, AccessWriter = getLogger(fmt.Sprintf("access/%s/out.log", optionsInner.ProjectName))
		businessLogger, BusinessWriter = getLogger(fmt.Sprintf("business/%s/out.log", optionsInner.ProjectName))
		errorLogger, ErrorWriter = getLogger(fmt.Sprintf("err/%s/error.log", optionsInner.ProjectName))
	}
}

func getLogger(filePath string) (zerolog.Logger, io.Writer) {
	var basePath string
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
		os.Exit(-1)
	}
	basePath = dir + "/log/"

	if isDevInner {
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

	return zerologger.Output(f), io.MultiWriter(f)
}

func Access(msg string) {
	accessLogger.Info().Msg(msg)
}

func Accessf(msg string, v ...interface{}) {
	accessLogger.Info().Msgf(msg, v...)
}

func Business(msg string) {
	businessLogger.Trace().Msg(msg)
}

func Businessf(msg string, v ...interface{}) {
	businessLogger.Trace().Msgf(msg, v...)
}

func Error(msg string) {
	errorLogger.Error().Msg(msg)
}

func Errorf(msg string, v ...interface{}) {
	errorLogger.Error().Msgf(msg, v...)
}
