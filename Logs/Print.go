package Logs

import (
	"fmt"
	"github.com/fatih/color"
	"time"
)

const (
	Print_log = 1
	Error     = 2
	Warm      = 3
)

func Print(mode int, i interface{}) {
	//color.Set(color.FgRed, color.Bold)
	str := ""
	if mode == 1 {
		color.Set(color.FgGreen, color.Bold)
		str = "Log "
	} else if mode == 2 {
		color.Set(color.FgRed, color.Bold)
		str = "Error "
	} else if mode == 3 {
		color.Set(color.FgYellow, color.Bold)
		str = "Warming "
	} else {
		color.Set(color.FgWhite, color.Bold)
	}
	fmt.Println(time.Now().String()[:19]+"\n", str, i)
}
