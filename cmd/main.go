package main

import (
	"context"
	"scp/internal/app"
)

func main() {
	ctx := context.Background()

	app.Run(
		ctx,
	)
}
