package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/arifwahyu/petverse-be/internal/app"
)

func main() {
	if err := app.Run(context.Background()); err != nil {
		slog.Error("application stopped", "error", err)
		os.Exit(1)
	}
}
