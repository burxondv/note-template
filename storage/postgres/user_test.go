package postgres_test

import (
	"testing"

	"github.com/burxondv/note-template/storage/repo"
	"github.com/bxcodec/faker/v4"
	"github.com/stretchr/testify/require"
)

// har gal run qilganda argumentlarini ozgartirish kere, chunki hamasi unique

func createUser(t *testing.T) *repo.User {
	user, err := strg.User().Create(&repo.User{
		FirstName:   faker.Name(),
		LastName:    faker.Name(),
		PhoneNumber: faker.Phonenumber(),
		Email:       faker.Email(),
		ImageURL:    faker.URL(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, user)

	return user
}

func deleteUser(id int64, t *testing.T) {
	err := strg.User().Delete(id)
	require.NoError(t, err)
}

func TestCreateUser(t *testing.T) {
	createUser(t)
}

func TestGetUser(t *testing.T) {
	c := createUser(t)

	user, err := strg.User().Get(c.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user)
}

func TestGetAllUser(t *testing.T) {
	user := createUser(t)

	users, err := strg.User().GetAll(&repo.GetAllUsersParams{
		Limit: 10,
		Page:  1,
	})

	require.NoError(t, err)
	require.GreaterOrEqual(t, len(users.Users), 1)

	deleteUser(user.ID, t)
}

func TestUpdateUser(t *testing.T) {
	user := createUser(t)

	user.FirstName = faker.Name()
	user.LastName = faker.Name()
	user.PhoneNumber = faker.Phonenumber()
	user.Email = faker.Email()
	user.ImageURL = faker.URL()
}

func TestDeleteUser(t *testing.T) {
	user := createUser(t)

	deleteUser(user.ID, t)
}
