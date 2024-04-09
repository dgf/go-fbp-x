package http_test

import (
	"fmt"
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgf/go-fbp-x/pkg/process/http"
)

func TestGet(t *testing.T) {
	get := http.Get()
	payload := "test-payload"

	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, _ *nethttp.Request) {
		fmt.Fprint(w, payload)
	}))
	defer server.Close()

	inputs := get.Inputs()
	inputs["timeout"].Channel <- "1S"
	inputs["url"].Channel <- server.URL

	outputs := get.Outputs()
	act := <-outputs["out"].Channel

	if act != payload {
		t.Errorf("output not matches got: %s, want: %s", act, payload)
	}
}
