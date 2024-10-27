package errs


var (
	// system error
	ErrNoFieldsUpdated     = NewCustomError("ErrNoFieldsUpdated", "无字段更新")
	ErrInternalServerError = NewCustomError("ErrInternalServerError", "服务器错误")
	
	// user error
	ErrUserNotFound             = NewCustomError("ErrUserNotFound", "用户不存在")
	ErrUserAlreadyExists        = NewCustomError("ErrUserAlreadyExists", "用户已存在")
	ErrUserConflict             = NewCustomError("ErrUserConflict", "用户冲突")
	ErrInvalidPassword          = NewCustomError("ErrInvalidPassword", "密码错误")
	ErrInvalidPasswordSameAsOld = NewCustomError("ErrInvalidPasswordSameAsOld", "新密码与旧密码一致")
	ErrInvalidUsernameFormat    = NewCustomError("ErrInvalidUsernameFormat", "用户名格式错误")
	ErrInvalidPhoneNumFormat    = NewCustomError("ErrInvalidPhoneNumFormat", "电话号码名格式错误")
	ErrInvalidEmailFormat       = NewCustomError("ErrInvalidEmailFormat", "邮箱格式错误")
	ErrInvalidPasswordFormat    = NewCustomError("ErrInvalidPasswordFormat", "密码格式错误")
	ErrUsernameAlreadyExists    = NewCustomError("ErrUsernameAlreadyExists", "用户名已存在")
	ErrUsernameNotFound         = NewCustomError("ErrUsernameNotFound", "用户名不存在")
	ErrPhoneNumNotFound         = NewCustomError("ErrPhoneNumNotFound", "电话号码不存在")
	ErrPhoneNumAlreadyExists    = NewCustomError("ErrPhoneNumAlreadyExists", "电话号码已存在")
	ErrEmailNotFound            = NewCustomError("ErrEmailNotFound", "邮箱不存在")
	ErrEmailAlreadyExists       = NewCustomError("ErrEmailAlreadyExists", "邮箱已存在")
	ErrInvalidEmailOrPhoneNum   = NewCustomError("ErrInvalidEmailOrPhoneNum", "邮箱或电话号码不存在")
	
	// user-role error
	ErrUserRoleAlreadyExists = NewCustomError("", "用户角色已存在")
	ErrUserRoleNotFound      = NewCustomError("", "用户角色不存在")

	// role error
	ErrRoleAlreadyExists = NewCustomError("ErrRoleAlreadyExists", "角色已存在")
	ErrRoleNotFound      = NewCustomError("ErrRoleNotFound", "角色不存在")
	ErrRoleConflict      = NewCustomError("ErrRoleConflict", "角色冲突")
	
	// role_permisson error
	ErrRolePermissionNotFound   = NewCustomError("ErrRolePermissionNotFound", "角色权限不存在")
	ErrRoleOrPermissionNotFound = NewCustomError("ErrRoleOrPermissionNotFound", "角色或权限不存在")
	
	// permission error
	ErrPermissionAlreadyExists           = NewCustomError("ErrPermissionAlreadyExists", "权限已存在")
	ErrPermissionNameAlreadyExists       = NewCustomError("ErrPermissionNameAlreadyExists", "permission name already exists")
	ErrPermissionPathMethodAlreadyExists = NewCustomError("ErrPermissionPathMethodAlreadyExists", "permission path method already exists")
	ErrInvalidAPIPathMethodFormat        = NewCustomError("ErrInvalidAPIPathMethodFormat", "invalid api path or method format")
	ErrPermissionNotFound                = NewCustomError("ErrPermissionNotFound", "permission not found")
	ErrPermissionConflict                = NewCustomError("ErrPermissionConflict", "permission conflict")

	// image error
	ErrImageNotFound = NewCustomError("ErrImageNotFound", "Image not found")
)