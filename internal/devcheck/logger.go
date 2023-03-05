package devcheck

import "github.com/fatih/color"

type Logger struct {
	indent string
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Indented() *Logger {
	return &Logger{
		indent: l.indent + "  ",
	}
}

func (l *Logger) writeMessage(prefixColor *color.Color, prefix string, format string, a ...interface{}) {
	reset := color.New(color.Reset)
	reset.Print(l.indent)
	prefixColor.Print(prefix)
	reset.Print(" ")
	reset.Printf(format, a...)
	reset.Println("")
}

func (l *Logger) Info(format string, a ...interface{}) {
	l.writeMessage(color.New(color.FgBlue), "i", format, a...)
}

func (l *Logger) Success(format string, a ...interface{}) {
	l.writeMessage(color.New(color.FgGreen), "✓", format, a...)
}

func (l *Logger) Warning(format string, a ...interface{}) {
	l.writeMessage(color.New(color.FgYellow), "!", format, a...)
}

func (l *Logger) Failure(format string, a ...interface{}) {
	l.writeMessage(color.New(color.FgRed), "✗", format, a...)
}

func (l *Logger) Bullet(format string, a ...interface{}) {
	l.writeMessage(color.New(color.Reset), "-", format, a...)
}

func (l *Logger) Error(e error) {
	// TODO
}

