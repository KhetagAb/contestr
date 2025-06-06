package wiresets

import (
	"contestr/pkg/config"
	"context"
	"github.com/google/wire"
)

func NewContextProvider() context.Context {
	return context.Background()
}

var CommonSet = wire.NewSet(
	NewContextProvider,
	config.NewConfig,
)
