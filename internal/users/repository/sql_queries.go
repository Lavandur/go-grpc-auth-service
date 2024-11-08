package repository

var (
	insertUsers = `INSERT INTO users (id, login, visible_id, hashed_password, person, role_ids, deleted_date, created_date, updated_date, last_password_restore_date_time, search_index) 
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`
	getListUsers = `SELECT * FROM users;`
	getUserByID  = `SELECT * FROM users WHERE id = $1;`
	deleteUser   = `DELETE FROM users WHERE id = $1;`
	// updateUser   = `INSERT INTO users (username, password) VALUES (?, ?)`
)
