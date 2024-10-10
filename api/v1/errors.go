package v1

var (
	// common errors
	ErrSuccess             = newError(0, "ok")
	ErrBadRequest          = newError(400, "Bad Request")
	ErrUnauthorized        = newError(401, "Unauthorized")
	ErrNotFound            = newError(404, "Not Found")
	ErrInternalServerError = newError(500, "Internal Server Error")

	// more biz errors
	ErrEmailAlreadyUse = newError(1001, "The email is already in use.")

	ErrEmailOrPasswordIncorrect = newError(1002, "邮箱或密码不正确")
	ErrUserNotFound         = newError(1003, "用户不存在")
	ErrRoleNotFound         = newError(1004, "角色不存在")
)
