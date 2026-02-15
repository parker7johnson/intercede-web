package utils

import (
	"fmt"
	"log"
	"net/http"
)

// Logger provides structured logging methods
type Logger struct {
	prefix string
}

func NewLogger(prefix string) *Logger {
	return &Logger{prefix: prefix}
}

func (l *Logger) Info(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	log.Printf("[INFO] %s: %s", l.prefix, msg)
}

func (l *Logger) Error(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	log.Printf("[ERROR] %s: %s", l.prefix, msg)
}

func (l *Logger) Success(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	log.Printf("[SUCCESS] %s: %s", l.prefix, msg)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	log.Printf("[WARN] %s: %s", l.prefix, msg)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	log.Printf("[DEBUG] %s: %s", l.prefix, msg)
}

func (l *Logger) LogRequest(r *http.Request, churchCode string) {
	log.Printf("[INFO] %s: %s %s from %s, Church-Code: %s",
		l.prefix, r.Method, r.URL.Path, r.RemoteAddr, churchCode)
}

func (l *Logger) LogBadRequest(r *http.Request, churchCode, reason string, err error) {
	if err != nil {
		log.Printf("[ERROR] %s: %s - Request: %s %s from %s, Church-Code: %s, Error: %v",
			l.prefix, reason, r.Method, r.URL.Path, r.RemoteAddr, churchCode, err)
	} else {
		log.Printf("[ERROR] %s: %s - Request: %s %s from %s, Church-Code: %s",
			l.prefix, reason, r.Method, r.URL.Path, r.RemoteAddr, churchCode)
	}
}
