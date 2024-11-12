package permission

var (
	setPermission = `
		DELETE FROM role_permissions
		WHERE role_id = @roleID;

		INSERT INTO role_permissions (id, name)
		VALUES (@roleID, @permission);
	`
	getPermissionByID = `
		SELECT * FROM role_permissions
		WHERE role_id = @roleID;
	`
)
