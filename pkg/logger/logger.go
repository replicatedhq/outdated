package logger

import (
	"fmt"

	"github.com/fatih/color"
)

type Logger struct {
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Info(msg string, args ...interface{}) {
	if msg == "" {
		fmt.Println("")
		return
	}

	c := color.New(color.FgHiCyan)
	c.Println(fmt.Sprintf(msg, args...))
}

func (l *Logger) Error(err error) {
	c := color.New(color.FgHiRed)
	c.Println(fmt.Sprintf("%#v", err))
}

func (l *Logger) Header(msg string, args ...interface{}) {
	c := color.New(color.FgHiWhite)
	c.Println(fmt.Sprintf(msg, args...))
}

func (l *Logger) StartImageLine(msg string, args ...interface{}) {
	c := color.New(color.FgHiYellow)
	c.Printf(fmt.Sprintf(msg, args...))
}

func (l *Logger) FinalizeImageLine(behind int64, msg string, args ...interface{}) {
	var c *color.Color

	if behind == 0 {
		c = color.New(color.FgHiGreen)
	} else if behind < 3 {
		c = color.New(color.FgHiYellow)
	} else {
		c = color.New(color.FgHiRed)
	}
	c.Println(fmt.Sprintf("\r"+msg, args...))
}

func (l *Logger) FinalizeImageLineWithError(msg string, args ...interface{}) {
	c := color.New(color.FgHiMagenta)

	c.Println(fmt.Sprintf("\r"+msg, args...))
}
