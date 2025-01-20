package feature

import (
	"fmt"
	"log"
	"os"
)

const (
	LOG_TYPE_ERROR  = "error"
	LOG_TYPE_ACTIVE = "active"
)

type Log struct {
	logType string
	file    string
}

func LogSetConfig() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
}

func GetLogActivity() *Log {
	return &Log{logType: LOG_TYPE_ACTIVE, file: os.Getenv("LOG_FILE_ACTIVE")}
}

func GetLogError() *Log {
	return &Log{logType: LOG_TYPE_ERROR, file: os.Getenv("LOG_FILE_ERROR")}
}

func (l *Log) WriteLog(msg string) error {
	if l.logType == LOG_TYPE_ACTIVE {
		l.writeLogActivity(msg)
		return nil
	}

	if l.logType == LOG_TYPE_ERROR {
		l.writeLogError(msg)
		return nil
	}

	return fmt.Errorf("arquivo de LOG não está configurado")
}

func (l *Log) writeLogError(msg string) error {
	file, err := os.OpenFile(os.Getenv("LOG_FILE_ERROR"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("arquivo de LOG não encontrado")
	}
	log.SetOutput(file)
	log.Println(msg)
	return nil
}

func (l *Log) writeLogActivity(msg string) error {
	file, err := os.OpenFile(os.Getenv("LOG_FILE_ACTIVE"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("arquivo de LOG não encontrado")
	}
	log.SetOutput(file)
	log.Println(msg)
	return nil
}

func (l *Log) ClearLog() error {
	if l.logType == LOG_TYPE_ACTIVE {
		l.clearLogError()
		return nil
	}

	if l.logType == LOG_TYPE_ERROR {
		l.clearLogActivity()
		return nil
	}

	return fmt.Errorf("arquivo de LOG não está configurado")
}

func (l *Log) clearLogError() error {
	file, err := os.OpenFile(os.Getenv("LOG_FILE_ERROR"), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("falha ao abrir arquivo de log de erro")
	}
	defer file.Close()

	log.SetOutput(file)
	log.Println("Arquivo de log reiniciado")
	return nil
}

func (l *Log) clearLogActivity() error {
	file, err := os.OpenFile(os.Getenv("LOG_FILE_ACTIVE"), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("falha ao abrir arquivo de log de atividade")
	}
	defer file.Close()

	log.SetOutput(file)
	log.Println("Arquivo de log reiniciado")
	return nil
}
