package queue

import (
	"github.com/deliergky/monitorink/data"
)

type ResponseDataIterator interface {
	data.Iterator
	Current() data.ResponseData
}

type Queue interface {
	ResponseDataIterator
	Push(data.ResponseData) error
	Pop() (data.ResponseData, error)
}
