package network

import "fmt"

type Trace struct {
	Packet any
	Connection
}

func (t Trace) String() string {
	return fmt.Sprintf("%s: %q", t.Connection, t.Packet)
}
