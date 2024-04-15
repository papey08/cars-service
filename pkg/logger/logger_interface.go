package logger

type Fields map[string]any

type Logger interface {
	Debug(fields Fields, message string)
	Info(fields Fields, message string)
	Error(fields Fields, message string)
	Fatal(fields Fields, message string)
}
