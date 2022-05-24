package util

import (
	"encoding/json"
	"fmt"
	"runtime"
)

func GenerateCallerInfo(skip int) string {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	f := runtime.FuncForPC(pc)
	return fmt.Sprintf("%s:%d %s", file, line, f.Name())
}

type stack struct {
	Func   string `json:"func"`
	Line   int    `json:"line"`
	Source string `json:"source"`
}

func (s *stack) String() string {
	return fmt.Sprintf("%s:%d %s", s.Source, s.Line, s.Func)
}

func Stacks(skip int) string {
	var stacks []*stack
	pc := make([]uintptr, 32)
	n := runtime.Callers(skip, pc)
	frames := runtime.CallersFrames(pc[:n])
	for {
		frame, more := frames.Next()
		stacks = append(stacks, &stack{Source: frame.File, Line: frame.Line, Func: frame.Function})
		if !more {
			break
		}
	}

	s, _ := json.Marshal(stacks)
	return string(s)
}
