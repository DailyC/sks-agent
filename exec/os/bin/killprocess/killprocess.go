package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/DailyC/sks-agent/exec"
	"github.com/DailyC/sks-agent/exec/os/bin"
)

var killProcessName string
var killProcessInCmd string

func main() {
	flag.StringVar(&killProcessName, "process", "", "process name")
	flag.StringVar(&killProcessInCmd, "process-cmd", "", "process in command")

	flag.Parse()

	killProcess(killProcessName, killProcessInCmd)
}

func killProcess(process, processCmd string) {
	var pids []string
	var err error
	var ctx = context.WithValue(context.Background(), exec.ExcludeProcessKey, "blade")
	if process != "" {
		pids, err = exec.GetPidsByProcessName(process, ctx)
		if err != nil {
			bin.PrintErrAndExit(err.Error())
		}
		killProcessName = process
	} else if processCmd != "" {
		pids, err = exec.GetPidsByProcessCmdName(processCmd, ctx)
		if err != nil {
			bin.PrintErrAndExit(err.Error())
		}
		killProcessName = processCmd
	}

	if pids == nil || len(pids) == 0 {
		bin.PrintErrAndExit(fmt.Sprintf("%s process not found", killProcessName))
		return
	}
	response := exec.NewLocalChannel().Run(ctx, "kill", fmt.Sprintf("-9 %s", strings.Join(pids, " ")))
	if !response.Success {
		bin.PrintErrAndExit(response.Err)
		return
	}
	bin.PrintOutputAndExit(response.Result.(string))
}
