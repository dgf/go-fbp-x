package network

import (
	"fmt"

	"github.com/dgf/go-fbp-x/pkg/dsl"
)

type Trace struct {
	Packet any
	dsl.Connection
}

func (t Trace) String() string {
	return fmt.Sprintf("%s: %q", t.Connection, t.Packet)
}
