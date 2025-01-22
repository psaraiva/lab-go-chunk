package logger

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
		l.WriteLogMessageInfo(msg)
		return nil
	}

	if l.logType == LOG_TYPE_ERROR {
		l.WriteLogMessageError(msg)
		return nil
	}

	return fmt.Errorf("nao foi possível escrever no log, tipo de Log inválido: %s", l.logType)
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

// IMPLEMENTAÇÃO DA INTERFACE DE LOG
func (l *Log) WriteLogMessageError(msg string) error {
	file, err := os.OpenFile(os.Getenv("LOG_FILE_ERROR"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("não foi possível escrever mensagem de erro, arquivo não encontrado")
	}
	log.SetOutput(file)
	log.Println(msg)
	return nil
}

func (l *Log) WriteLogMessageInfo(msg string) error {
	file, err := os.OpenFile(os.Getenv("LOG_FILE_ACTIVE"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("não foi possível escrever mensagem de informação, arquivo não encontrado")
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

	return fmt.Errorf("nao foi possível limpar log, tipo de Log inválido: %s", l.logType)
}
