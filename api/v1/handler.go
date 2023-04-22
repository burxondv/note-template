package v1

import (
	"github.com/burxondv/note-template/api/models"
	"github.com/burxondv/note-template/config"
	"github.com/burxondv/note-template/storage"
)

type handlerV1 struct {
	cfg      *config.Config
	storage  storage.StorageI
	inMemory storage.InMemoryStorageI
}

type HandlerV10Options struct {
	Cfg      *config.Config
	Storage  storage.StorageI
	InMemory storage.InMemoryStorageI
}

func New(options *HandlerV10Options) *handlerV1 {
	return &handlerV1{
		cfg:      options.Cfg,
		storage:  options.Storage,
		inMemory: options.InMemory,
	}
}

func errorResponse(err error) *models.ErrorResponse {
	return &models.ErrorResponse{
		Error: err.Error(),
	}
}
