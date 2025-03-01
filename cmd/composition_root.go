package cmd

import (
	"context"
)

type CompositionRoot struct {
}

func NewCompositionRoot(ctx context.Context) CompositionRoot {
	app := CompositionRoot{}
	return app
}
