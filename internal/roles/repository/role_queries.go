package repository

var (
	createRole = `
		INSERT INTO roles (id, title, description, created_at)
		VALUES (@roleID, @title, @description, @createdAt)
		RETURNING id, title, description, created_at;
	`
	deleteRole = `
		DELETE FROM roles
		WHERE id = @roleID;
	`
	updateRole = `
		UPDATE roles
		SET title = @title, description = @description
		WHERE id = @roleID
		RETURNING id, title, description, created_at;
	`
	getRoleByID = `
		SELECT id, title, description, created_at FROM roles
		WHERE id = @roleID
	`

	getPermissionsByRole = `
		SELECT * FROM role_permissions
		WHERE id = @roleID;
	`
)
