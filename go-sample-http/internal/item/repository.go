package item

import (
	"database/sql"
	"errors"
)

var ErrNotFound = errors.New("item not found")

type Repository interface {
	FindAll() ([]Item, error)
	FindByID(id int64) (*Item, error)
	Create(req CreateRequest) (*Item, error)
	Update(id int64, req UpdateRequest) (*Item, error)
	Delete(id int64) error
}

type sqliteRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &sqliteRepository{db: db}
}

func (r *sqliteRepository) FindAll() ([]Item, error) {
	rows, err := r.db.Query(`SELECT id, name, description, created_at FROM items ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *sqliteRepository) FindByID(id int64) (*Item, error) {
	var item Item
	err := r.db.QueryRow(`SELECT id, name, description, created_at FROM items WHERE id = ?`, id).
		Scan(&item.ID, &item.Name, &item.Description, &item.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &item, err
}

func (r *sqliteRepository) Create(req CreateRequest) (*Item, error) {
	result, err := r.db.Exec(`INSERT INTO items (name, description) VALUES (?, ?)`, req.Name, req.Description)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.FindByID(id)
}

func (r *sqliteRepository) Update(id int64, req UpdateRequest) (*Item, error) {
	result, err := r.db.Exec(`UPDATE items SET name = ?, description = ? WHERE id = ?`, req.Name, req.Description, id)
	if err != nil {
		return nil, err
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		return nil, ErrNotFound
	}
	return r.FindByID(id)
}

func (r *sqliteRepository) Delete(id int64) error {
	result, err := r.db.Exec(`DELETE FROM items WHERE id = ?`, id)
	if err != nil {
		return err
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		return ErrNotFound
	}
	return nil
}
