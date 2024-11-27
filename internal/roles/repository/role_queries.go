package repository

var (
	createRole = `
		INSERT INTO roles (id, name, description, created_at)
		VALUES (@roleID, @name, @description, @createdAt)
		RETURNING id, name, description, created_at;
	`
	deleteRole = `
		DELETE FROM roles
		WHERE id = @roleID;
	`
	updateRole = `
		UPDATE roles
		SET name = @name, description = @description
		WHERE id = @roleID
		RETURNING id, name, description, created_at;
	`
	getRoleByID = `
		SELECT id, name, description, created_at FROM roles
		WHERE id = @roleID
	`

	getPermissionsByRole = `
		SELECT * FROM role_permissions
		WHERE id = @roleID;
	`
)
