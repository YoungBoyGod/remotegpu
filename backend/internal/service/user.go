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
	customerDao *dao.CustomerDao
}

func NewUserService() *UserService {
	return &UserService{
		customerDao: dao.NewCustomerDao(),
	}
}

// Register 用户注册
func (s *UserService) Register(req *v1.RegisterRequest) error {
	// 检查用户名是否存在
	if _, err := s.customerDao.GetByUsername(req.Username); err == nil {
		return errors.New("用户名已存在")
	}

	// 检查邮箱是否存在
	if _, err := s.customerDao.GetByEmail(req.Email); err == nil {
		return errors.New("邮箱已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 创建用户
	user := &entity.Customer{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		DisplayName:  req.Nickname,
		UserType:     "external", // 默认为外部用户
		AccountType:  "individual", // 默认为个人账户
		Status:       "active",
	}

	return s.customerDao.Create(user)
}

// Login 用户登录
func (s *UserService) Login(req *v1.LoginRequest) (*v1.LoginResponse, error) {
	// 获取用户
	user, err := s.customerDao.GetByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 检查用户状态
	if user.Status != "active" {
		return nil, errors.New("用户已被禁用")
	}

	// 生成 token
	token, err := auth.GenerateToken(user.ID, user.Username, user.UserType)
	if err != nil {
		return nil, err
	}

	return &v1.LoginResponse{
		Token: token,
		User: &v1.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Nickname: user.DisplayName,
			Avatar:   user.AvatarURL,
			Role:     user.UserType,
			Status:   1, // 临时保持兼容
		},
	}, nil
}

// GetUserInfo 获取用户信息
func (s *UserService) GetUserInfo(userID uint) (*v1.UserInfo, error) {
	user, err := s.customerDao.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return &v1.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Nickname: user.DisplayName,
		Avatar:   user.AvatarURL,
		Role:     user.UserType,
		Status:   1, // 临时保持兼容
	}, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(userID uint, req *v1.UpdateUserRequest) error {
	user, err := s.customerDao.GetByID(userID)
	if err != nil {
		return err
	}

	if req.Nickname != "" {
		user.DisplayName = req.Nickname
	}
	if req.Avatar != "" {
		user.AvatarURL = req.Avatar
	}

	return s.customerDao.Update(user)
}
