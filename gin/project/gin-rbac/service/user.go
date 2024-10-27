package service

import (
	"errors"

	"gin-rbac/common/errs"
	"gin-rbac/db/dao"
	"gin-rbac/db/model"
	"gin-rbac/dtos"
	"gin-rbac/global"
	"gin-rbac/utils"

	"gorm.io/gorm"
)



// UserService 用户服务
type UserService interface {
	// 登录用户
	LoginUser(loginUserReqDTO dtos.LoginUserReqDTO) (*dtos.TokenResDTO, error)
	// 创建用户
	CreateUser(createUserReqDTO dtos.CreateUserReqDTO) (*dtos.GetFullUserResDTO, error)
	// 根据id获取用户信息
	GetUserByID(getFullUserReqDTO dtos.GetFullUserReqDTO) (dtos.GetFullUserResDTO, error)
	// 获取公开用户列表
	GetPublicUserList(publicUserListReqDTO dtos.PublicUserListReqDTO) (dtos.PaginationResult[dtos.PublicUserResDTO], error)
	// 根据id获取公开用户信息
	GetPublicUserByID(publicUserReqDTO dtos.PublicUserReqDTO) (dtos.PublicUserResDTO, error)
	// 更新用户信息
	UpdateUser(getFullUserReqDTO dtos.GetFullUserReqDTO, updateUserReqDTO dtos.UpdateUserReqDTO) error
	// 更新用户头像
	UpdateUserAvatar(dtos.GetFullUserReqDTO, dtos.UpdateUserAvatarReqDTO) error
	// 更新密码
	UpdatePassword(getFullUserReqDTO dtos.GetFullUserReqDTO, updateUserPasswordReqDTO dtos.UpdatePasswordReqDTO) error
	// 重置密码
	ResetPassword(resetPasswordReqDTO dtos.ResetPasswordReqDTO) error
	// 删除用户
	DeleteUser(deleteUserReqDTO dtos.DeleteUserReqDTO) error
	// 恢复用户
	RecoverUser(recoverUserReqDTO dtos.RecoverUserReqDTO) error
}

// userService 用户服务实现
type userService struct {
	userDao dao.UserDao
}

// NewUserService 创建用户服务
func NewUserService(userDao dao.UserDao) UserService {
	return &userService{
		userDao: userDao,
	}
}

// LoginUser 用户登录
func (s *userService) LoginUser(loginUserReqDTO dtos.LoginUserReqDTO) (*dtos.TokenResDTO, error) {
	var user *model.UserModel
	var err error

	// 根据凭证尝试解析其类型，并根据其类型获取用户
	switch {
	case utils.IsValidEmail(loginUserReqDTO.Credential):
		user, err = s.userDao.GetUserByEmail(loginUserReqDTO.Credential)
	case utils.IsValidPhoneNumber(loginUserReqDTO.Credential):
		user, err = s.userDao.GetUserByPhoneNum(loginUserReqDTO.Credential)
	default:
		user, err = s.userDao.GetUserByUsername(loginUserReqDTO.Credential)
	}
	if err != nil {
		// 用户不存在
		if errors.Is(err, gorm.ErrRecordNotFound) {
			switch {
			case utils.IsValidEmail(loginUserReqDTO.Credential):
				// 404 邮箱不存在
				global.Log.Warnln("Failed to login user, Email not found: ", loginUserReqDTO.Credential)
				return nil, errs.ErrEmailNotFound
			case utils.IsValidPhoneNumber(loginUserReqDTO.Credential):
				// 404 手机号不存在
				global.Log.Warnln("Failed to login user, Phone number not found: ", loginUserReqDTO.Credential)
				return nil, errs.ErrPhoneNumNotFound
			default:
				// 404 用户名不存在
				global.Log.Warnln("Failed to login user, Username not found: ", loginUserReqDTO.Credential)
				return nil, errs.ErrUsernameNotFound
			}
		}
		// 500 服务器错误
		global.Log.Errorln("Failed to login user, Failed to get user: ", err)
		return nil, errs.ErrInternalServerError
	}

	// 确认密码匹配
	if !utils.ComparePasswordHash(user.Password, loginUserReqDTO.Password) {
		// 400 密码错误
		global.Log.Warnln("Failed to login user, Password does not match: ", err)
		return nil, errs.ErrInvalidPassword
	}

	// 生成JWT
	token, err := utils.GenerateJWT(user.ID, user.Username, user.IsSuperAdmin, global.Config.JWT.Secret, global.Config.JWT.Expiration)
	if err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to login user, Failed to generate JWT: ", err)
		return nil, errs.ErrInternalServerError
	}
	return &dtos.TokenResDTO{
		Token: token,
	}, nil
}

func (s *userService) CreateUser(createUserReqDTO dtos.CreateUserReqDTO) (*dtos.GetFullUserResDTO, error) {
	// 创建用户模型。其中的密码需要经过加密处理
	user := &model.UserModel{}
	err := dtos.ConvertDTOToModel(&createUserReqDTO, user)
	if err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to create user, Failed to convert DTO to model: ", err)
		return nil, errs.ErrInternalServerError
	}
	// 检查电话号码和邮箱是否非空
	if user.PhoneNum == "" && user.Email == "" {
		// 400 客户端错误
		global.Log.Warnln("Failed to create user, Phone number and email cannot be empty at the same time")
		return nil, errs.ErrInvalidEmailOrPhoneNum
	}
	// 检查用户名是否合法
	if !utils.IsValidUsername(user.Username) {
		// 400 客户端错误
		global.Log.Warnln("Failed to create user, Invalid username")
		return nil, errs.ErrInvalidUsernameFormat
	}
	// 检查用户名是否重复
	if err := s.checkDuplicateUsername(user.Username); err != nil {
		// 409 冲突
		global.Log.Warnln("Failed to create user, Username already exists: ", err)
		return nil, errs.ErrUsernameAlreadyExists
	}

	// 检查电话号码是否合法和重复
	if user.PhoneNum != "" {
		if !utils.IsValidPhoneNumber(user.PhoneNum) {
			// 400 客户端错误
			global.Log.Warnln("Failed to create user, Invalid phone number")
			return nil, errs.ErrInvalidPhoneNumFormat
		}
		if err := s.checkDuplicatePhoneNum(user.PhoneNum); err != nil {
			// 409 冲突
			global.Log.Warnln("Failed to create user, Phone number already exists: ", err)
			return nil, errs.ErrPhoneNumAlreadyExists
		}
	}
	// 检查邮箱是否合法和重复
	if user.Email != "" {
		if !utils.IsValidEmail(user.Email) {
			// 400 客户端错误
			global.Log.Warnln("Failed to create user, Invalid email")
			return nil, errs.ErrInvalidEmailFormat
		}
		if err := s.checkDuplicateEmail(user.Email); err != nil {
			// 409 冲突
			global.Log.Warnln("Failed to create user, Email already exists: ", err)
			return nil, errs.ErrEmailAlreadyExists
		}
	}
	// 验证密码格式
	if !utils.IsValidPassword(createUserReqDTO.Password) {
		// 400 客户端错误
		global.Log.Warnln("Failed to create user, Invalid password")
		return nil, errs.ErrInvalidPasswordFormat
	}

	// // 生成随机盐并散列密码
	// hashedPassword, err := utils.HashPassword(user.Password, passwordCost)
	// if err != nil {
	// 	// 500 服务器错误
	// 	global.Log.Errorln("Failed to create user, Failed to generate password hash: ", err)
	// 	return nil, errs.ErrInternalServerError
	// }
	// // 将密码和盐值转换为Base64编码
	// base64HashedPassword := base64.StdEncoding.EncodeToString([]byte(hashedPassword))
	// // 将Base64编码密码和盐值存入用户模型
	// user.Password = base64HashedPassword

	user, err = s.userDao.CreateUser(user)
	if err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to create user: ", err)
		return nil, errs.ErrInternalServerError
	}
	
	// // 生成JWT
	// token, err := utils.GenerateJWT(user.ID, user.Username, user.IsSuperAdmin, global.Config.JWT.Secret, global.Config.JWT.Expiration)
	// if err != nil {
	// 	// 500 服务器错误
	// 	global.Log.Errorln("Failed to login user, Failed to generate JWT: ", err)
	// 	return nil, errs.ErrInternalServerError
	// }
	
	return &dtos.GetFullUserResDTO{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		PhoneNum:   user.PhoneNum,
		Sex:        user.Sex,
		Intro:      user.Intro,
		CreatedAtFormat: user.CreatedAt.Format("2006-01-02 15:04:05"),
	},nil
}

func (s *userService) GetPublicUserList(publicUserListReqDTO dtos.PublicUserListReqDTO) (dtos.PaginationResult[dtos.PublicUserResDTO], error) {
	users, count, err := s.userDao.GetUserList(&publicUserListReqDTO)
	if err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to get user list: ", err)
		return dtos.PaginationResult[dtos.PublicUserResDTO]{}, errs.ErrInternalServerError
	}
	// 遍历用户列表，将每个用户转换为PublicUserInfoResDTO类型
	var usersData []dtos.PublicUserResDTO
	for _, user := range users {
		usersData = append(usersData, dtos.PublicUserDTOToPublicUserResDTO(*user))
	}
	// 创建分页结果
	paginatedResult := dtos.NewPaginationResult[dtos.PublicUserResDTO](
		count, usersData, publicUserListReqDTO.Page, publicUserListReqDTO.Size)

	return paginatedResult, nil

}

func (s *userService) GetPublicUserByID(publicUserReqDTO dtos.PublicUserReqDTO) (dtos.PublicUserResDTO, error) {
	id := publicUserReqDTO.ID
	user, err := s.userDao.GetPublicUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 404 用户不存在
			global.Log.Warnln("Failed to get public user, User not found: ", err)
			return dtos.PublicUserResDTO{}, errs.ErrUserNotFound
		}
		// 500 服务器错误
		global.Log.Errorln("Failed to get public user: ", err)
		return dtos.PublicUserResDTO{}, errs.ErrInternalServerError
	}
	return dtos.PublicUserDTOToPublicUserResDTO(*user), nil
}

func (s *userService) GetUserByID(getFullUserReqDTO dtos.GetFullUserReqDTO) (dtos.GetFullUserResDTO, error) {
	id := getFullUserReqDTO.ID
	user, err := s.userDao.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 404 用户不存在
			global.Log.Warnln("Failed to get user, User not found: ", err)
			return dtos.GetFullUserResDTO{}, errs.ErrUserNotFound
		}
		// 500 服务器错误
		global.Log.Errorln("Failed to get user: ", err)
		return dtos.GetFullUserResDTO{}, errs.ErrInternalServerError
	}
	return dtos.FullUserDTOToGetFullUserResDTO(*user), nil
}

func (s *userService) UpdateUser(getFullUserReqDTO dtos.GetFullUserReqDTO, updateUserReqDTO dtos.UpdateUserReqDTO) error {
	user := &model.UserModel{}
	user.ID = getFullUserReqDTO.ID
	// 根据ID获取当前用户
	existingUser, err := s.userDao.GetUserByID(getFullUserReqDTO.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 404 用户不存在
			global.Log.Warnln("Failed to update user, User not found: ", err)
			return errs.ErrUserNotFound
		}
		// 500 服务器错误
		global.Log.Errorln("Failed to update user, Failed to get user: ", err)
		return errs.ErrInternalServerError
	}
	// 创建一个布尔变量来追踪是否有字段需要更新
	hasChanges := false
	if updateUserReqDTO.Username != "" && updateUserReqDTO.Username != existingUser.Username {
		// 检查用户名是否重复
		if err := s.checkDuplicateUsername(updateUserReqDTO.Username); err != nil {
			// 409 冲突
			global.Log.Warnln("Failed to update user, Username already exists: ", err)
			return errs.ErrUsernameAlreadyExists
		}
		existingUser.Username = updateUserReqDTO.Username
		hasChanges = true
	}
	if updateUserReqDTO.Email != "" && updateUserReqDTO.Email != existingUser.Email {
		// 检查邮箱格式是否正确
		if !utils.IsValidEmail(updateUserReqDTO.Email) {
			// 400 客户端错误
			global.Log.Warnln("Failed to update user, Invalid email")
			return errs.ErrInvalidEmailFormat
		}
		// 检查邮箱是否重复
		if err := s.checkDuplicateEmail(updateUserReqDTO.Email); err != nil {
			// 409 冲突
			global.Log.Warnln("Failed to update user, Email already exists: ", err)
			return errs.ErrEmailAlreadyExists
		}
		existingUser.Email = updateUserReqDTO.Email
		hasChanges = true
	}
	if updateUserReqDTO.PhoneNum != "" && updateUserReqDTO.PhoneNum != existingUser.PhoneNum {
		// 检查电话号码格式是否正确
		if !utils.IsValidPhoneNumber(updateUserReqDTO.PhoneNum) {
			// 400 客户端错误
			global.Log.Warnln("Failed to update user, Invalid phone number")
			return errs.ErrInvalidPhoneNumFormat
		}
		// 检查电话号码是否重复
		if err := s.checkDuplicatePhoneNum(updateUserReqDTO.PhoneNum); err != nil {
			// 409 冲突
			global.Log.Warnln("Failed to update user, Phone number already exists: ", err)
			return errs.ErrPhoneNumAlreadyExists
		}
		existingUser.PhoneNum = updateUserReqDTO.PhoneNum
		hasChanges = true
	}
	if updateUserReqDTO.Sex != existingUser.Sex {
		existingUser.Sex = updateUserReqDTO.Sex
		hasChanges = true
	}
	if updateUserReqDTO.Intro != "" && updateUserReqDTO.Intro != existingUser.Intro {
		existingUser.Intro = updateUserReqDTO.Intro
		hasChanges = true
	}
	// 如果没有任何字段需要更新，则提前返回错误
	if !hasChanges {
		// 422 无字段更新
		global.Log.Warnln("Failed to update user, No fields to update")
		return errs.ErrNoFieldsUpdated
	}
	err = dtos.ConvertDTOToModel(&updateUserReqDTO, user)
	if err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to update user, Failed to convert DTO to model: ", err)
		return errs.ErrInternalServerError
	}
	if err = s.userDao.UpdateUser(user); err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to update user: ", err)
		return errs.ErrInternalServerError
	}
	return nil
}

func (s *userService) UpdateUserAvatar(getFullUserReqDTO dtos.GetFullUserReqDTO, updateUserAvatarReqDTO dtos.UpdateUserAvatarReqDTO) error {
	// 获取用户信息
	user, err := s.userDao.GetUserByID(getFullUserReqDTO.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 404 用户不存在
			global.Log.Warnln("Failed to update user avatar, User not found: ", err)
			return errs.ErrUserNotFound
		}
		// 500 服务器错误
		global.Log.Errorln("Failed to update user avatar, Failed to get user: ", err)
		return errs.ErrInternalServerError
	}
	if err = s.userDao.UpdateUserAvatar(user.ID, updateUserAvatarReqDTO.AvatarID); err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to update user avatar: ", err)
		return errs.ErrInternalServerError
	}
	return nil
}

func (s *userService) UpdatePassword(getFullUserReqDTO dtos.GetFullUserReqDTO, updatePasswordReqDTO dtos.UpdatePasswordReqDTO) error {
	// 检查新密码是否与旧密码相同
	if updatePasswordReqDTO.OldPassword == updatePasswordReqDTO.NewPassword {
		// 400 请求错误
		global.Log.Warnln("Failed to update user password, New password cannot be the same as old password")
		return errs.ErrInvalidPasswordSameAsOld
	}
	// 获取用户信息
	password, err := s.userDao.GetUserPasswordByID(getFullUserReqDTO.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 404 用户不存在
			global.Log.Warnln("Failed to update user password, User not found: ", err)
			return errs.ErrUserNotFound
		}
		// 500 服务器错误
		global.Log.Errorln("Failed to update user password, Failed to get user: ", err)
		return errs.ErrInternalServerError
	}
	
	// // 获取数据库密码的base64编码并将其转为哈希值
	// UserPassword, err := base64.StdEncoding.DecodeString(password)
	// if err != nil {
	// 	// 500 服务器错误
	// 	global.Log.Errorln("Failed to update user password, Failed to decode password: ", err)
	// 	return errs.ErrInternalServerError
	// }
	
	global.Log.Debugf("Password: %s", password)
	// 验证原密码
	if !utils.ComparePasswordHash("$2a$12$ZBII.37VJE4BE2gLm5NQAOWMJ6lloNB8.xmHtQHLKCBjaNr4.k6P2", updatePasswordReqDTO.OldPassword) {
		// 400 客户端错误
		global.Log.Warnln("Failed to update user password, Incorrect old password")
		return errs.ErrInvalidPassword
	}

	// // 生成新密码的哈希值
	// hashedNewPassword, err := utils.HashPassword(updatePasswordReqDTO.NewPassword, passwordCost)
	// if err != nil {
	// 	// 500 服务器错误
	// 	global.Log.Errorln("Failed to update user password, Failed to hash password: ", err)
	// 	return errs.ErrInternalServerError
	// }

	// // 将密码和盐值转换为Base64编码
	// base64HashedPassword := base64.StdEncoding.EncodeToString([]byte(hashedNewPassword))
	// // 将Base64编码密码和盐值存入用户模型
	// password = base64HashedPassword
	
	if err = s.userDao.UpdatePassword(getFullUserReqDTO.ID, updatePasswordReqDTO.NewPassword); err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to update user password: ", err)
		return errs.ErrInternalServerError
	}
	return nil
}

func (s *userService) ResetPassword(resetPasswordReqDTO dtos.ResetPasswordReqDTO) error {
	var err error
	var user *model.UserModel

	// 根据凭证尝试解析其类型
	switch {
	case utils.IsValidEmail(resetPasswordReqDTO.Credential):
		user, err = s.userDao.GetUserByEmail(resetPasswordReqDTO.Credential)
	case utils.IsValidPhoneNumber(resetPasswordReqDTO.Credential):
		user, err = s.userDao.GetUserByPhoneNum(resetPasswordReqDTO.Credential)
	default:
		user, err = s.userDao.GetUserByUsername(resetPasswordReqDTO.Credential)
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 404 用户不存在
			global.Log.Warnln("Failed to reset user password, User not found: ", err)
			return errs.ErrUserNotFound
		}
		// 500 服务器错误
		global.Log.Errorln("Failed to reset user password, Failed to get user: ", err)
		return errs.ErrInternalServerError
	}
	// // 生成新密码的哈希值
	// hashedNewPassword, err := utils.HashPassword(resetPasswordReqDTO.NewPassword, passwordCost)
	// if err != nil {
	// 	// 500 服务器错误
	// 	global.Log.Errorln("Failed to reset user password, Failed to hash password: ", err)
	// 	return errs.ErrInternalServerError
	// }
	// // 将密码和盐值转换为Base64编码
	// base64HashedNewPassword := base64.StdEncoding.EncodeToString([]byte(hashedNewPassword))
	// // 将Base64编码密码和盐值存入用户模型
	// password := base64HashedNewPassword
	
	if err = s.userDao.UpdatePassword(user.ID, resetPasswordReqDTO.NewPassword); err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to reset user password, Failed to update password: ", err)
		return errs.ErrInternalServerError
	}
	return nil
}

func (s *userService) DeleteUser(deleteUserReqDTO dtos.DeleteUserReqDTO) error {
	if _, err := s.userDao.GetUserByID(deleteUserReqDTO.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 404 资源不存在
			global.Log.Warnln("Failed to delete user, User not found: ", err)
			return errs.ErrUserNotFound
		}
		// 500 服务器错误
		global.Log.Errorln("Failed to delete user, Failed to get user: ", err)
		return errs.ErrInternalServerError
	}
	if err := s.userDao.DeleteUser(deleteUserReqDTO.ID); err != nil {
		// 500 服务器错误
		global.Log.Errorln("Failed to delete user: ", err)
		return errs.ErrInternalServerError
	}
	return nil
}

func (s *userService) RecoverUser(recoverUserReqDTO dtos.RecoverUserReqDTO) error {
	if user, err := s.userDao.GetDelUserByID(recoverUserReqDTO.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 404 用户不存在
			global.Log.Errorln("Failed to recover user, User not found: ", err)
			return errs.ErrUserNotFound
		}
		// 500 内部服务器错误
		global.Log.Errorln("Failed to recover user, Failed to get user: ", err)
		return errs.ErrInternalServerError
	} else if !user.DeletedAt.Valid {
		// 409 冲突
		global.Log.Errorln("Failed to recover user, User not deleted: ", err)
		return errs.ErrUserConflict
	}
	if err := s.userDao.RecoverUser(recoverUserReqDTO.ID); err != nil {
		// 500 内部服务器错误
		global.Log.Errorln("Failed to recover user: ", err)
		return errs.ErrInternalServerError
	}
	return nil
}

// checkDuplicateUsername 检查 Username 是否重复的函数
//
// 描述:
// 该函数检查提供的用户名是否已存在于数据库中。
// 如果用户名已存在，则返回错误。
// 如果查询数据库时发生其他错误，则也返回错误。
//
// 参数:
// - username: 需要检查的用户名
//
// 返回:
// - error: 如果用户名已存在或发生其他错误则返回
func (s *userService) checkDuplicateUsername(username string) error {
	// 查询数据库中是否存在该用户名
	user, err := s.userDao.GetUserByUsername(username)
	// 如果查询出错且不是记录不存在的错误，返回错误
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// 500 服务器错误
		global.Log.Errorln("Failed to check duplicate username:", err)
		return errs.ErrInternalServerError
	}
	// 如果查询结果不为空，说明用户名重复，返回错误(直接返回用户不存在即可)
	if user != nil {
		// 409 冲突
		global.Log.Warnln("Username already exists")
		return errs.ErrUserAlreadyExists
	}
	return nil
}

// checkDuplicatePhoneNum 检查 PhoneNum 是否重复的函数
//
// 描述:
// 该函数检查提供的电话号码是否已存在于数据库中。
// 如果电话号码已存在，则返回错误。
// 如果查询数据库时发生其他错误，则也返回错误。
//
// 参数:
// - phoneNum: 需要检查的电话号码
//
// 返回:
// - error: 如果电话号码已存在或发生其他错误则返回
func (s *userService) checkDuplicatePhoneNum(phoneNum string) error {
	// 查询数据库中是否存在该电话号码
	user, err := s.userDao.GetUserByPhoneNum(phoneNum)
	// 如果查询出错且不是记录不存在的错误，返回错误
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// 500 服务器错误
		global.Log.Errorln("Failed to check duplicate phone number: ", err)
		return errs.ErrInternalServerError
	}
	// 如果查询结果不为空，说明电话号码重复，返回错误
	if user != nil {
		// 409 冲突
		global.Log.Warnln("Phone number is already exists", err)
		return errs.ErrPhoneNumAlreadyExists
	}
	return nil
}

// checkDuplicateEmail 检查 Email 是否重复的函数
//
// 描述:
// 该函数检查提供的邮箱是否已存在于数据库中。
// 如果邮箱已存在，则返回错误。
// 如果查询数据库时发生其他错误，则也返回错误。
//
// 参数:
// - email: 需要检查的邮箱地址
//
// 返回:
// - error: 如果邮箱已存在或发生其他错误则返回
func (s *userService) checkDuplicateEmail(email string) error {
	// 查询数据库中是否存在该邮箱
	user, err := s.userDao.GetUserByEmail(email)
	// 如果查询出错且不是记录不存在的错误，返回错误
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// 500 服务器错误
		global.Log.Errorln("Failed to check duplicate email: ", err)
		return errs.ErrInternalServerError
	}
	// 如果查询结果不为空，说明邮箱重复，返回错误
	if user != nil {
		// 409 冲突
		global.Log.Warnln("Email is already exists")
		return errs.ErrEmailAlreadyExists
	}
	return nil
}
