package storage

import (
	"github.com/burxondv/note-template/storage/postgres"
	"github.com/burxondv/note-template/storage/repo"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type StorageI interface {
	User() repo.UserStorageI
	Note() repo.NoteStorageI
}

type storagePg struct {
	userRepo repo.UserStorageI
	noteRepo repo.NoteStorageI
}

func NewStoragePg(db *sqlx.DB) StorageI {
	return &storagePg{
		userRepo: postgres.NewUser(db),
		noteRepo: postgres.NewNote(db),
	}
}

func (s *storagePg) User() repo.UserStorageI {
	return s.userRepo
}

func (s *storagePg) Note() repo.NoteStorageI {
	return s.noteRepo
}
