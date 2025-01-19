package main

import (
	"flag"
	"fmt"
	"lab/feature"
	"strings"
)

func main() {
	var action = feature.MakeAction()

	feature.LogSetConfig()
	var logActive = feature.GetLogActivity()

	logActive.WriteLog("Iniciando aplicação...")
	println("Iniciando aplicação...")

	arg_action := flag.String("action", "", "Action to invoke (upload/download/remove/clear)")
	arg_file_target := flag.String("file-target", "", "File target to action")

	flag.Parse()

	if !validArgAction(arg_action) {
		fmt.Println("Invalid action:", *arg_action)
		return
	}

	if *arg_action != feature.ACTION_CLEAR && !validArgFileTarget(arg_file_target) {
		fmt.Println("Invalid target:", arg_file_target)
		return
	}

	action.Action = *arg_action
	action.FileTarget = *arg_file_target
	err := action.Execute()
	if err != nil {
		println(err.Error())
	}

	logActive.WriteLog("Finalizando aplicação...")
	println("Finalizando aplicação...")
}

func validArgAction(arg_action *string) bool {
	switch strings.ToLower(*arg_action) {
	case feature.ACTION_CLEAR,
		feature.ACTION_DOWNLOAD,
		feature.ACTION_UPLOAD,
		feature.ACTION_REMOVE:
		return true
	}

	return false
}

func validArgFileTarget(file_target *string) bool {
	return len(*file_target) > 1
}
