package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/yunkeCN/ali-logger-golang/util"
)

type Logger struct {
	accessLogger   *zerolog.Logger
	businessLogger *zerolog.Logger
	errorLogger    *zerolog.Logger

	tag            string
	isOutputCaller byte //这里使用3态逻辑，以区分是否做了设置   0未设置  1不输出 2输出
	isOutputStack  byte //这里使用3态逻辑，以区分是否做了设置   0未设置  1不输出 2输出
	commonFields   []map[string]interface{}
	fields         []map[string]interface{}
}

func (l Logger) WithCaller(isOutput bool) Logger {
	if isOutput {
		l.isOutputCaller = 2
	} else {
		l.isOutputCaller = 1
	}

	return l
}

func (l Logger) WithStack(isOutput bool) Logger {
	if isOutput {
		l.isOutputStack = 2
	} else {
		l.isOutputStack = 1
	}

	return l
}

//通用字段设置，输出为一级节点
func (l Logger) WithCommonField(key string, val interface{}) Logger {
	l.commonFields = append(l.commonFields, map[string]interface{}{key: val})

	return l
}

//通用字段设置，输出为一级节点
func (l Logger) WithCommonFields(fields map[string]interface{}) Logger {
	l.commonFields = append(l.commonFields, fields)

	return l
}

//自定义字段设置，作为一级节点attatch的内容输出
func (l Logger) WithField(key string, val interface{}) Logger {
	l.fields = append(l.fields, map[string]interface{}{key: val})

	return l
}

//自定义字段设置，作为一级节点attatch的内容输出
func (l Logger) WithFields(fields map[string]interface{}) Logger {
	l.fields = append(l.fields, fields)

	return l
}

func (l Logger) ClearFields() Logger {
	l.fields = l.fields[0:0]

	return l
}

func (l Logger) setDefaultField(zeroEvent *zerolog.Event) {
	zeroEvent.Str("app", optionsInner.ProjectName)
	zeroEvent.Str("tag", l.tag)
}

func (l Logger) msg(zeroEvent *zerolog.Event, msg string, withCaller bool, withStack bool) {
	l.setDefaultField(zeroEvent)

	if withCaller {
		zeroEvent.Str("caller", util.GenerateCallerInfo(3))
	}

	if withStack {
		zeroEvent.Str("stack", util.Stacks(4))
	}

	if l.commonFields != nil && len(l.commonFields) > 0 {
		for _, val := range l.commonFields {
			zeroEvent.Fields(val)
		}
	}

	if l.fields != nil && len(l.fields) > 0 {
		zeroEvent.Fields(map[string]interface{}{"attach": l.fields})
	}

	zeroEvent.Msg(msg)
}

func (l Logger) Access(msg string) {
	zeroEvent := l.getZeroLogger(l.accessLogger, zerolog.InfoLevel).Info()
	l.msg(zeroEvent, msg, l.getIsOutput(true, l.isOutputCaller), l.getIsOutput(false, l.isOutputStack))
}

func (l Logger) Accessf(format string, v ...interface{}) {
	zeroEvent := l.getZeroLogger(l.accessLogger, zerolog.InfoLevel).Info()
	l.msg(zeroEvent, fmt.Sprintf(format, v...), l.getIsOutput(true, l.isOutputCaller), l.getIsOutput(false, l.isOutputStack))
}

func (l Logger) Business(msg string) {
	zeroEvent := l.getZeroLogger(l.businessLogger, zerolog.TraceLevel).Trace()
	l.msg(zeroEvent, msg, l.getIsOutput(true, l.isOutputCaller), l.getIsOutput(false, l.isOutputStack))
}

func (l Logger) Businessf(format string, v ...interface{}) {
	zeroEvent := l.getZeroLogger(l.businessLogger, zerolog.TraceLevel).Trace()
	l.msg(zeroEvent, fmt.Sprintf(format, v...), l.getIsOutput(true, l.isOutputCaller), l.getIsOutput(false, l.isOutputStack))
}

func (l Logger) Error(msg string) {
	zeroEvent := l.getZeroLogger(l.errorLogger, zerolog.ErrorLevel).Error()
	l.msg(zeroEvent, msg, l.getIsOutput(false, l.isOutputCaller), l.getIsOutput(true, l.isOutputStack))
}

func (l Logger) Errorf(format string, v ...interface{}) {
	zeroEvent := l.getZeroLogger(l.errorLogger, zerolog.ErrorLevel).Error()
	l.msg(zeroEvent, fmt.Sprintf(format, v...), l.getIsOutput(false, l.isOutputCaller), l.getIsOutput(true, l.isOutputStack))
}

//获取最终是否要输出
//参数：defaultIsOutput  默认是否输出
//参数：isSetOutput  设置的是否输出
func (l Logger) getIsOutput(defaultIsOutput bool, isSetOutput byte) bool {
	if isSetOutput == 0 { //如果未设置，则直接用默认
		return defaultIsOutput
	} else {
		if isSetOutput == 1 { //如果设置不输出，则不输出
			return false
		} else {
			return true
		}
	}
}

func (l Logger) getZeroLogger(zeroLoggger *zerolog.Logger, loggerType zerolog.Level) *zerolog.Logger {
	if zeroLoggger != nil { //如果设置，则直接用
		return zeroLoggger
	} else {
		switch loggerType {
		case zerolog.InfoLevel:
			return &accessLogger
		case zerolog.TraceLevel:
			return &businessLogger
		case zerolog.ErrorLevel:
			return &errorLogger
		default:
			panic("其他类型不支持")
		}
	}
}
