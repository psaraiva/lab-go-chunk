package main

import (
	"flag"
	"fmt"
	"lab/src/internal/actions"
	"lab/src/logger"
	"lab/src/repository"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var EngineRepositoryFile = ""
var EngineRepositoryChunkItem = ""

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	resp := validConfigRepositorty()
	if !resp {
		panic("Error loading engine repository")
	}

	var action = actions.MakeAction()

	logger.LogSetConfig()
	var logActive = logger.GetLogActivity()

	println("Iniciando aplicação...")
	err = logActive.WriteLog("Iniciando aplicação...")
	if err != nil {
		panic("Não foi possível iniciar log de atividades, aplicação será encerrada")
	}

	arg_action := flag.String("action", "", "Action to invoke (upload/download/remove/clear)")
	arg_file_target := flag.String("file-target", "", "File target to action")

	flag.Parse()

	if !validArgAction(arg_action) {
		fmt.Println("Invalid action:", *arg_action)
		return
	}

	if *arg_action != actions.ACTION_CLEAR && !validArgFileTarget(arg_file_target) {
		fmt.Println("Invalid target:", arg_file_target)
		return
	}

	action.Type = *arg_action
	action.FileTarget = *arg_file_target
	err = actions.Execute(&action)
	if err != nil {
		println(err.Error())
	}

	logActive.WriteLog("Finalizando aplicação...")
	println("Finalizando aplicação...")
}

func validArgAction(arg_action *string) bool {
	switch strings.ToLower(*arg_action) {
	case actions.ACTION_CLEAR,
		actions.ACTION_DOWNLOAD,
		actions.ACTION_UPLOAD,
		actions.ACTION_REMOVE:
		return true
	}

	return false
}

func validArgFileTarget(file_target *string) bool {
	return len(*file_target) > 1
}

func validConfigRepositorty() bool {
	resp_file, resp_chunk_item := false, false
	EngineRepositoryFile = os.Getenv("ENGINE_COLLECTION_FILE")
	EngineRepositoryChunkItem = os.Getenv("ENGINE_COLLECTION_CHUNK")

	switch EngineRepositoryFile {
	case repository.ENGINE_JSON:
		resp_file = true
	}

	switch EngineRepositoryChunkItem {
	case repository.ENGINE_JSON:
		resp_chunk_item = true
	}

	return resp_file && resp_chunk_item
}
