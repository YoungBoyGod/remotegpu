package service

import (
	"errors"

	"github.com/YoungBoyGod/remotegpu/api/v1"
	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/auth"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	userDao *dao.UserDao
}

func NewUserService() *UserService {
	return &UserService{
		userDao: dao.NewUserDao(),
	}
}

// Register 用户注册
func (s *UserService) Register(req *v1.RegisterRequest) error {
	// 检查用户名是否存在
	if _, err := s.userDao.GetByUsername(req.Username); err == nil {
		return errors.New("用户名已存在")
	}

	// 检查邮箱是否存在
	if _, err := s.userDao.GetByEmail(req.Email); err == nil {
		return errors.New("邮箱已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 创建用户
	user := &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
		Role:     auth.RoleUser, // 默认角色为普通用户
		Status:   1,
	}

	return s.userDao.Create(user)
}

// Login 用户登录
func (s *UserService) Login(req *v1.LoginRequest) (*v1.LoginResponse, error) {
	// 获取用户
	user, err := s.userDao.GetByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, errors.New("用户已被禁用")
	}

	// 生成 token
	token, err := auth.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	return &v1.LoginResponse{
		Token: token,
		User: &v1.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Role:     user.Role,
			Status:   user.Status,
		},
	}, nil
}

// GetUserInfo 获取用户信息
func (s *UserService) GetUserInfo(userID uint) (*v1.UserInfo, error) {
	user, err := s.userDao.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return &v1.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Role:     user.Role,
		Status:   user.Status,
	}, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(userID uint, req *v1.UpdateUserRequest) error {
	user, err := s.userDao.GetByID(userID)
	if err != nil {
		return err
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	return s.userDao.Update(user)
}
