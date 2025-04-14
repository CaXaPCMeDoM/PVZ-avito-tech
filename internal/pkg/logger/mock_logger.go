package logger

import (
	"fmt"
	"strings"
)

// MockLogger реализует Interface для тестирования
type MockLogger struct {
	DebugLogs []string
	InfoLogs  []string
	WarnLogs  []string
	ErrorLogs []string
	FatalLogs []string
}

func NewMock() *MockLogger {
	return &MockLogger{
		DebugLogs: make([]string, 0),
		InfoLogs:  make([]string, 0),
		WarnLogs:  make([]string, 0),
		ErrorLogs: make([]string, 0),
		FatalLogs: make([]string, 0),
	}
}

func (m *MockLogger) Debug(message interface{}, args ...interface{}) {
	m.DebugLogs = append(m.DebugLogs, formatMessage(message, args...))
}

func (m *MockLogger) Info(message string, args ...interface{}) {
	m.InfoLogs = append(m.InfoLogs, fmt.Sprintf(message, args...))
}

func (m *MockLogger) Warn(message string, args ...interface{}) {
	m.WarnLogs = append(m.WarnLogs, fmt.Sprintf(message, args...))
}

func (m *MockLogger) Error(message interface{}, args ...interface{}) {
	m.ErrorLogs = append(m.ErrorLogs, formatMessage(message, args...))
}

func (m *MockLogger) Fatal(message interface{}, args ...interface{}) {
	m.FatalLogs = append(m.FatalLogs, formatMessage(message, args...))
}

// Вспомогательная функция для форматирования сообщений
func formatMessage(message interface{}, args ...interface{}) string {
	switch msg := message.(type) {
	case string:
		if len(args) > 0 {
			return fmt.Sprintf(msg, args...)
		}
		return msg
	case error:
		return msg.Error()
	case fmt.Stringer:
		return msg.String()
	default:
		return fmt.Sprintf("%v", message)
	}
}

func (m *MockLogger) Reset() {
	m.DebugLogs = nil
	m.InfoLogs = nil
	m.WarnLogs = nil
	m.ErrorLogs = nil
	m.FatalLogs = nil
}

func (m *MockLogger) HasDebug(msg string) bool {
	return hasMessage(m.DebugLogs, msg)
}

func (m *MockLogger) HasInfo(msg string) bool {
	return hasMessage(m.InfoLogs, msg)
}

func (m *MockLogger) HasWarn(msg string) bool {
	return hasMessage(m.WarnLogs, msg)
}

func (m *MockLogger) HasError(msg string) bool {
	return hasMessage(m.ErrorLogs, msg)
}

func (m *MockLogger) HasFatal(msg string) bool {
	return hasMessage(m.FatalLogs, msg)
}

func hasMessage(logs []string, search string) bool {
	for _, entry := range logs {
		if strings.Contains(entry, search) {
			return true
		}
	}
	return false
}
