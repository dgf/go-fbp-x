package web

import (
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
)

//go:embed static
var static embed.FS

func Serve(addr string) error {
	slog.Info("listen and serve", "addr", addr)
	if root, err := fs.Sub(static, "static"); err != nil {
		return fmt.Errorf("serve root failed: %w", err)
	} else {
		return http.ListenAndServe(addr, http.FileServer(http.FS(root)))
	}
}
