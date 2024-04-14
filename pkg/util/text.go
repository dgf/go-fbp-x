package util

import (
	"fmt"
	"strings"
)

var rawUnescape = strings.NewReplacer(`\n`, "\n", `\r`, "\r", `\t`, "\t")

func CastAndUnescapeRaw(raw any) (string, error) {
	s, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("value %v", raw)
	}
	return rawUnescape.Replace(s), nil
}
