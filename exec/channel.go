package exec

import (
	"context"
	"github.com/DailyC/sks-agent/transport"
)

type Channel interface {
	// Run command
	Run(ctx context.Context, script, args string) *transport.Response

	// GetScriptPath
	GetScriptPath() string
}
