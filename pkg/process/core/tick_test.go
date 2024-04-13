package core_test

import (
	"testing"
	"time"

	"github.com/dgf/go-fbp-x/pkg/process/core"
)

func TestTick(t *testing.T) {
	tick := core.Tick()
	out := tick.Outputs()["out"]

	intv := tick.Inputs()["intv"]
	defer close(intv.Channel)
	data := tick.Inputs()["data"]
	defer close(data.Channel)

	packet := "test"
	intv.Channel <- "0.001S"
	data.Channel <- packet

	select {
	case act := <-out.Channel:
		if act != packet {
			t.Errorf("Invalid output, got: %q, want: %q", act, packet)
		}
	case <-time.After(73 * time.Millisecond):
		t.Error("test timed out, no output")
	}
}

func BenchmarkTick(b *testing.B) {
	tick := core.Tick()

	intv := tick.Inputs()["intv"]
	defer close(intv.Channel)
	data := tick.Inputs()["data"]
	defer close(data.Channel)

	intv.Channel <- "1H2M3S"
	data.Channel <- "test"

	b.Run("interval", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			intv.Channel <- "1H2M3S"
		}
	})

	b.Run("data", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			data.Channel <- "test"
		}
	})
}
