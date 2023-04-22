package postgres

import (
	"database/sql"
	"fmt"

	"github.com/burxondv/note-template/storage/repo"
	"github.com/jmoiron/sqlx"
)

type noteRepo struct {
	db *sqlx.DB
}

func NewNote(db *sqlx.DB) repo.NoteStorageI {
	return &noteRepo{
		db: db,
	}
}

func (ur *noteRepo) Create(note *repo.Note) (*repo.Note, error) {
	query := `
		INSERT INTO notes(
			user_id,
			title,
			description
		) VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	row := ur.db.QueryRow(
		query,
		note.UserID,
		note.Title,
		note.Description,
	)

	err := row.Scan(
		&note.ID,
		&note.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return note, nil
}

func (ur *noteRepo) Get(id int64) (*repo.Note, error) {
	var result repo.Note

	query := `
		SELECT 
		    id,
            user_id,
            title,
			description,
            created_at,
			updated_at,
			deleted_at
		FROM notes
        WHERE id=$1
	`

	row := ur.db.QueryRow(query, id)

	err := row.Scan(
		&result.ID,
		&result.UserID,
		&result.Title,
		&result.Description,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ur *noteRepo) GetAll(params *repo.GetAllNotesParams) (*repo.GetAllNotesResult, error) {
	result := repo.GetAllNotesResult{
		Notes: make([]*repo.Note, 0),
	}

	offset := (params.Page - 1) * params.Limit

	limit := fmt.Sprintf(" LIMIT %d OFFSET %d", params.Limit, offset)

	filter := "WHERE true"
	if params.Search != "" {
		filter += " AND title ilike '%" + params.Search + "%' "
	}

	if params.UserID != 0 {
		filter += fmt.Sprintf(" AND user_id=%d ", params.UserID)
	}

	orderBy := " ORDER BY created_at desc "
	if params.SortByData != "" {
		orderBy = fmt.Sprintf(" ORDER BY created_at %s ", params.SortByData)
	}

	query := `
		SELECT
			id,
			user_id,
			title,
			description,
			created_at,
            updated_at,
            deleted_at
		FROM notes
		` + filter + orderBy + limit

	rows, err := ur.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var u repo.Note

		err := rows.Scan(
			&u.ID,
			&u.UserID,
			&u.Title,
			&u.Description,
			&u.CreatedAt,
			&u.UpdatedAt,
			&u.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		result.Notes = append(result.Notes, &u)
	}

	queryCount := `SELECT count(*) FROM notes ` + filter
	err = ur.db.QueryRow(queryCount).Scan(&result.Count)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ur *noteRepo) Update(note *repo.Note) (*repo.Note, error) {
	query := `
		UPDATE notes SET
			title=$1,
            description=$2
		WHERE id=$3
		RETURNING id, user_id, title, description, created_at, updated_at, deleted_at
	`

	row := ur.db.QueryRow(
		query,
		note.Title,
		note.Description,
		note.ID,
	)

	var result repo.Note
	err := row.Scan(
		&result.ID,
		&result.UserID,
		&result.Title,
		&result.Description,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return note, nil
}

func (ur *noteRepo) Delete(id int64) error {
	query := "DELETE FROM notes WHERE id=$1"

	result, err := ur.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsCount, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsCount == 0 {
		return sql.ErrNoRows
	}

	return nil
}
