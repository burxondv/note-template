package postgres_test

import (
	"testing"

	"github.com/burxondv/note-template/storage/repo"
	"github.com/bxcodec/faker/v4"
	"github.com/stretchr/testify/require"
)

func createNote(t *testing.T) *repo.Note {
	note, err := strg.Note().Create(&repo.Note{
		UserID:      3,
		Title:       faker.Name(),
		Description: faker.Sentence(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, note)

	return note
}

func deleteNote(id int64, t *testing.T) {
	err := strg.Note().Delete(id)
	require.NoError(t, err)
}

func TestCreateNote(t *testing.T) {
	createNote(t)
}

func TestGetNote(t *testing.T) {
	c := createNote(t)

	note, err := strg.Note().Get(c.ID)
	require.NoError(t, err)
	require.NotEmpty(t, note)
}

func TestGetAllNote(t *testing.T) {
	note := createNote(t)

	notes, err := strg.Note().GetAll(&repo.GetAllNotesParams{
		Limit: 10,
		Page:  1,
	})

	require.NoError(t, err)
	require.GreaterOrEqual(t, len(notes.Notes), 1)

	deleteUser(note.ID, t)
}

func TestUpdateNote(t *testing.T) {
	note := createNote(t)

	note.UserID = 1
	note.Title = faker.Name()
	note.Description = faker.Sentence()
}

func TestDeleteNote(t *testing.T) {
	user := createNote(t)

	deleteNote(user.ID, t)
}
