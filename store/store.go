package store

import (
	"context"

	"github.com/deliergky/monitorink/data"
)

type Store interface {
	Persist(context.Context, data.ResponseData) error
	Retrieve(context.Context) ([]data.ResponseData, error)
}
