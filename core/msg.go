package core

import (
	"fmt"
)

var reset  = "\033[0m"
var red    = "\033[31m"
var green  = "\033[32m"
//var Yellow = "\033[33m"
//var Blue   = "\033[34m"


func MsgInfo(text string) {
	fmt.Println(green + "INFO: " + reset + text)
}

func MsgErr(text string) {
	fmt.Println(red + "ERROR: " + text + reset)
}
