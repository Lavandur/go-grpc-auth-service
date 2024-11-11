package repository

var (
	getUserByID = `
		SELECT id, login, visible_id, hashed_password, person, role_ids, 
			   created_date, updated_date, deleted_date, last_password_restore_date, search_index 
		FROM users 
		WHERE id = @userID AND deleted_date IS NULL;
	`

	insertUser = `
		INSERT INTO users (id, login, visible_id, hashed_password, person, role_ids, 
						   created_date, updated_date, deleted_date, last_password_restore_date) 
		VALUES (@userID, @login, @visibleID, @hashedPassword, @person, @roles, 
				@createdAt, @updatedAt, @deletedAt, @lastPasswordRestoreAt)
		RETURNING id, login, visible_id, hashed_password, person, role_ids, 
				  created_date, updated_date, deleted_date, last_password_restore_date, search_index;
	`

	updateUser = `
		UPDATE users 
		SET login = @login, visible_id = @visibleID, hashed_password = @hashedPassword, 
			person = @person, role_ids = @roles, updated_date = @updatedAt, 
			deleted_date = @deletedAt, last_password_restore_date = @lastPasswordRestoreAt 
		WHERE id = @userID
		RETURNING id, login, visible_id, hashed_password, person, role_ids, 
				  created_date, updated_date, deleted_date, last_password_restore_date, search_index;
	`

	deleteUser = `
		UPDATE users 
		SET deleted_date = @deletedAt
		WHERE id = @userID
		RETURNING id, login, visible_id, hashed_password, person, role_ids, 
				  created_date, updated_date, deleted_date, last_password_restore_date, search_index;
	`
)
