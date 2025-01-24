package main

import (
	"flag"
	"fmt"
	"lab/src/internal/service"
	"lab/src/logger"
	"lab/src/repository"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var EngineRepositoryFile = ""
var EngineRepositoryChunk = ""

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	resp := isValidConfigRepositorty()
	if !resp {
		fmt.Println("Error loading engine repository")
		return
	}

	logger.LogSetConfig()
	var serviceAction = service.MakeAction()

	println("Iniciando aplicação...")
	err = logger.GetLogActivity().WriteLog("Iniciando aplicação...")
	if err != nil {
		logger.GetLogError().WriteLog(err.Error())
		fmt.Println("Não foi possível iniciar log de atividades, aplicação será encerrada")
		return
	}

	arg_action := flag.String("action", "", "Action to invoke (upload/download/remove/clear)")
	arg_file_target := flag.String("file-target", "", "File target to action")

	flag.Parse()

	if !isValidArgAction(arg_action) {
		logger.GetLogError().WriteLog(fmt.Errorf("invalid value of action: %s", *arg_action).Error())
		fmt.Println("Invalid value of action:", *arg_action)
		return
	}

	if *arg_action != service.ACTION_CLEAR && !isValidArgFileTarget(arg_file_target) {
		logger.GetLogError().WriteLog(fmt.Errorf("invalid value of file-target: %s", *arg_file_target).Error())
		fmt.Println("Invalid value of file-target:", *arg_file_target)
		return
	}

	serviceAction.Type = *arg_action
	serviceAction.FileTarget = *arg_file_target
	err = service.Execute(&serviceAction)
	if err != nil {
		println(err.Error())
	}

	logger.GetLogActivity().WriteLog("Finalizando aplicação...")
	println("Finalizando aplicação...")
}

func isValidArgAction(arg_action *string) bool {
	switch strings.ToLower(*arg_action) {
	case service.ACTION_CLEAR,
		service.ACTION_DOWNLOAD,
		service.ACTION_UPLOAD,
		service.ACTION_REMOVE:
		return true
	}

	return false
}

func isValidArgFileTarget(file_target *string) bool {
	return len(*file_target) > 1
}

func isValidConfigRepositorty() bool {
	resp_file, resp_chunk_item := false, false
	EngineRepositoryFile = os.Getenv("ENGINE_COLLECTION_FILE")
	EngineRepositoryChunk = os.Getenv("ENGINE_COLLECTION_CHUNK")

	switch EngineRepositoryFile {
	case repository.ENGINE_JSON:
		resp_file = true
	}

	switch EngineRepositoryChunk {
	case repository.ENGINE_JSON:
		resp_chunk_item = true
	}

	return resp_file && resp_chunk_item
}
