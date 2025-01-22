package interfaces

type LogUnified interface {
	LogError
	LogActivity
	ClearLog
}

type LogError interface {
	WriteLogMessageError(msg string) error
}

type LogActivity interface {
	WriteLogMessageInfo(msg string) error
}

type ClearLog interface {
	ClearLog() error
}
