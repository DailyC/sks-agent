package cmd

import (
	"fmt"
	"github.com/DailyC/sks-agent/transport"
	"github.com/spf13/cobra"
)

type ServerCommand struct {
	baseCommand
}

func (sc *ServerCommand) Init() {
	sc.command = &cobra.Command{
		Use:     "server",
		Short:   "Server mode starts, exposes web services",
		Long:    "Server mode starts, exposes web services. Under the mode, you can send http request to trigger experiments",
		Aliases: []string{"srv"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return transport.ReturnFail(transport.Code[transport.IllegalCommand],
				fmt.Sprintf("less start or stop command"))
		},
		Example: serverExample(),
	}
}

func serverExample() string {
	return `blade server start --port 8000`
}
