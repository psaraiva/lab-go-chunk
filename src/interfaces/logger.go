package interfaces

type LogUnified interface {
	LogError
	LogActivity
	Clear
}

type LogError interface {
	WriteLogMessageError(msg string) error
}

type LogActivity interface {
	WriteLogMessageInfo(msg string) error
}

type Clear interface {
	Clear() error
}
