package main

import (
	"fmt"
	"github.com/Centny/gwf/util"
	"github.com/Centny/nms"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: nms <-c|-s> <configure file>")
		os.Exit(1)
		return
	}
	var fcfg = util.NewFcfg3()
	fcfg.InitWithFilePath2(os.Args[2], true)
	switch os.Args[1] {
	case "-s":
		fmt.Println(nms.RunNMS_S(fcfg))
	case "-c":
		fmt.Println(nms.RunNMS_C(fcfg))
	default:
		fmt.Println("unknow option %v", os.Args[1])
		os.Exit(1)
	}
}
