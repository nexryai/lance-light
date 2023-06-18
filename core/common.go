package core

import (
	"fmt"
	"os"
)

var reset  = "\033[0m"
var red    = "\033[31m"
var green  = "\033[32m"
var yellow = "\033[33m"
var gray   = "\033[37m"


func MsgInfo(text string) {
	fmt.Println(green + "INFO: " + reset + text)
}

func MsgErr(text string) {
	fmt.Println(red + "✘ ERROR: " + text + reset)
}

func MsgWarn(text string) {
	fmt.Println(yellow + "⚠ WARNING: " + reset + text)
}

func MsgDebug(text string) {
	fmt.Println(gray + "DEBUG: " + text + reset)
}

func MsgDetail(text string) {
	fmt.Println(gray + "  ↳ " + reset + text)
}

func ExitOnError(err error, message string) {
	if err != nil {
		errorInfo := fmt.Sprintf("Fatal error: %v", err)
		MsgErr(errorInfo)
		MsgDetail(message)
		os.Exit(1)
	}

	return
}
