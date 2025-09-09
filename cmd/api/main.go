package main

import (
	"github.com/ilindan-dev/shortener/internal/app"
	"go.uber.org/fx"
)

// main is the entry point for the URL shortener API application.
func main() {
	fx.New(app.Module).Run()
}
