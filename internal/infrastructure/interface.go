package infrastructure

type ILogger interface {
	IsDebugMode() bool
	Debug(msg string, ctx ...interface{})
	Info(msg string, ctx ...interface{})
	Warn(msg string, ctx ...interface{})
	Error(msg string, ctx ...interface{})
	Crit(msg string, ctx ...interface{})
}

type IBroker interface {
	SendMessage(topic string, message []byte, retained bool) error
	Subscribe(topic string, handler func(message []byte)) error
	Close()
}
