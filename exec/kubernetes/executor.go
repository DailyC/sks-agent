package kubernetes

import (
	"context"
	"github.com/DailyC/sks-agent/exec"
	"github.com/DailyC/sks-agent/transport"
)

type Executor struct {
}

func (*Executor) Name() string {
	return "k8s"
}

func (e *Executor) SetChannel(channel exec.Channel) {
}

func (*Executor) Exec(uid string, ctx context.Context, model *exec.ExpModel) *transport.Response {
	return transport.ReturnSuccess("k8s command")
}
