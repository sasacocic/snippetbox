package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {

	// had to use a function here, because I don't think a prepared statement will evaluate interval $3 || 'days' - and things like this
	// but it will evaluate a function

	stmt := `INSERT INTO snippets(title, content, created, expires)
    VALUES($1, $2, CURRENT_TIMESTAMP, date_add(CURRENT_TIMESTAMP, concat_ws(' ', $3::int, 'days')::interval ));`

	_, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id := float64(420) // postgresql driver doesn't support LastInsertId
	// id, err := result.LastInsertId()
	//if err != nil {
	//	return 0, err
	//}

	// the ID returned has the type int64, so we convert it to an int type
	// before returning.
	return int(id), nil
}

func (m *SnippetModel) Get(id int) (Snippet, error) {

	stmt := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > NOW() AND id = $1`
	row := m.DB.QueryRow(stmt, id)

	var s Snippet

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]Snippet, error) {

	stmt := `SELECT id, title, content, created, expires
    FROM snippets
    WHERE expires > NOW()
    ORDER BY id
    DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var snippets []Snippet

	for rows.Next() {
		var s Snippet

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
