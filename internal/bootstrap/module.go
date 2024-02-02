package bootstrap

import (
	"context"
)

type Module interface {
	Name() string
	Execute(context.Context) error
}
