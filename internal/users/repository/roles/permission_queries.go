package roles

var (
	clearPermission = `
		DELETE FROM role_permissions
		WHERE role_id = @roleID;
	`
	setPermissions = `
		INSERT INTO role_permissions (role_id, permission)
		VALUES (@roleID, @permission);
	`
	getPermissionsByID = `
		SELECT * FROM role_permissions
		WHERE role_id = @roleID;
	`
)
