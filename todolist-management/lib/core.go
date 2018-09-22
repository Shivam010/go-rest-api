package todolist

import (
	"database/sql"
	"errors"
	"fmt"
)

// Generic error messages
var (
	ErrNotFound     = errors.New("list not found")
	ErrItemNotFound = errors.New("item not found")
)

// Core ...
type Core struct {
	db *sql.DB
}

// NewCore implements Todo List Management Core Logic
func NewCore(db *sql.DB) *Core {
	return &Core{db}
}

// Ex create user
func (c *Core) Ex(name string) error {
	if err := c.db.Ping(); err == nil {
		fmt.Println("pinging")
		return err
	}
	fmt.Println("yes", name)
	return nil
}

// TodoItem ...
type TodoItem struct {
	ID        int64  `json:"id"`
	Value     string `json:"value"`
	Completed bool   `json:"completed"`
}

// TodoList ...
type TodoList struct {
	ID    int64       `json:"id"`
	Items []*TodoItem `json:"items"`
	Name  string      `json:"name"`
}

// AddTodoList creates a todo list with it's items
func (c *Core) AddTodoList(list *TodoList) (*TodoList, error) {
	const listQuery = `INSERT INTO todolist_management.todo_lists (name) VALUES($1) returning id`
	const itemQuery = `INSERT INTO todolist_management.todo_items (value, list_id, completed) VALUES($1, $2, $3) returning id`

	tx, err := c.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := tx.QueryRow(listQuery, list.Name).Scan(&list.ID); err != nil {
		return nil, err
	}

	stmt, err := tx.Prepare(itemQuery)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	for _, item := range list.Items {
		if err := stmt.QueryRow(item.Value, list.ID, item.Completed).Scan(&item.ID); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return list, nil
}

// DeleteTodoList removes a todo list with it's items
func (c *Core) DeleteTodoList(id int64) error {
	const check = `SELECT id FROM todolist_management.todo_lists WHERE id = $1`
	cid := int64(0)
	if err := c.db.QueryRow(check, id).Scan(&cid); err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		return ErrNotFound
	}

	const listQuery = `DELETE FROM todolist_management.todo_lists WHERE id = $1`
	const itemQuery = `DELETE FROM todolist_management.todo_items WHERE list_id = $1`

	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(listQuery, id); err != nil {
		return err
	}
	if _, err := tx.Exec(itemQuery, id); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// EditTodoListName updates the name of the list
func (c *Core) EditTodoListName(id int64, name string) error {
	const check = `SELECT id FROM todolist_management.todo_lists WHERE id = $1`
	cid := int64(0)
	if err := c.db.QueryRow(check, id).Scan(&cid); err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		return ErrNotFound
	}

	const query = `UPDATE todolist_management.todo_lists SET name = $2 WHERE id = $1`
	if _, err := c.db.Exec(query, id, name); err != nil {
		return err
	}
	return nil
}

// AddTodoItem adds item to the list
func (c *Core) AddTodoItem(lid int64, item *TodoItem) (*TodoItem, error) {
	const check = `SELECT id FROM todolist_management.todo_lists WHERE id = $1`
	cid := int64(0)
	if err := c.db.QueryRow(check, lid).Scan(&cid); err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
		return nil, ErrNotFound
	}

	const query = `INSERT INTO todolist_management.todo_items (value, list_id, completed) VALUES($1, $2, $3) returning id`
	if err := c.db.QueryRow(query, item.Value, lid, item.Completed).Scan(&item.ID); err != nil {
		return nil, err
	}
	return item, nil
}

// DeleteTodoListItem removes items from the list
func (c *Core) DeleteTodoListItem(id int64) error {
	const query = `DELETE FROM todolist_management.todo_items WHERE id = $1`
	if _, err := c.db.Exec(query, id); err != nil {
		return err
	}
	return nil
}

// GetTodoListItem returns a todolist item
func (c *Core) GetTodoListItem(id int64) (*TodoItem, error) {
	const query = `SELECT id, value, completed FROM todolist_management.todo_items WHERE id = $1`
	item := &TodoItem{}
	if err := c.db.QueryRow(query, id).Scan(&item.ID, &item.Value, &item.Completed); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrItemNotFound
		}
		return nil, err
	}
	return item, nil
}

// UpdateTodoItem updates an item
func (c *Core) UpdateTodoItem(item *TodoItem) error {
	const query = `UPDATE todolist_management.todo_items SET value = $1, completed = $2 WHERE id = $3`
	if _, err := c.db.Exec(query, item.Value, item.Completed, item.ID); err != nil {
		return err
	}
	return nil
}

// GetTodoList returns whole todolist
func (c *Core) GetTodoList(id int64) (*TodoList, error) {
	const query = `SELECT * FROM todolist_management.todo_lists INNER JOIN todolist_management.todo_items ON todolist_management.todo_items.list_id = todolist_management.todo_lists.id AND todolist_management.todo_lists.id = $1`
	rows, err := c.db.Query(query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	defer rows.Close()
	list := &TodoList{
		Items: []*TodoItem{},
	}
	for rows.Next() {
		item := &TodoItem{}
		if err = rows.Scan(&list.ID, &list.Name, &item.ID, &item.Value, &list.ID, &item.Completed); err != nil {
			return nil, err
		}
		list.Items = append(list.Items, item)
	}
	return list, nil
}
