package service

import (
	"os"
	"testing"

	"github.com/YoungBoyGod/remotegpu/api/v1"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/auth"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	// 初始化JWT用于测试
	auth.InitJWT("test-secret-key-for-unit-testing-only", 3600)

	// 运行测试
	code := m.Run()

	os.Exit(code)
}

// MockCustomerDao 是 CustomerDao 的 mock 实现
type MockCustomerDao struct {
	mock.Mock
}

func (m *MockCustomerDao) Create(customer *entity.Customer) error {
	args := m.Called(customer)
	return args.Error(0)
}

func (m *MockCustomerDao) GetByID(id uint) (*entity.Customer, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Customer), args.Error(1)
}

func (m *MockCustomerDao) GetByUsername(username string) (*entity.Customer, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Customer), args.Error(1)
}

func (m *MockCustomerDao) GetByEmail(email string) (*entity.Customer, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Customer), args.Error(1)
}

func (m *MockCustomerDao) Update(customer *entity.Customer) error {
	args := m.Called(customer)
	return args.Error(0)
}

func (m *MockCustomerDao) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockCustomerDao) List(page, pageSize int) ([]*entity.Customer, int64, error) {
	args := m.Called(page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.Customer), args.Get(1).(int64), args.Error(2)
}

// TestUserService_Register_Success 测试注册成功
func TestUserService_Register_Success(t *testing.T) {
	mockDao := new(MockCustomerDao)
	service := &UserService{
		customerDao: mockDao,
	}

	req := &v1.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Nickname: "Test User",
	}

	// Mock GetByUsername - 用户名不存在
	mockDao.On("GetByUsername", req.Username).Return(nil, gorm.ErrRecordNotFound)
	// Mock GetByEmail - 邮箱不存在
	mockDao.On("GetByEmail", req.Email).Return(nil, gorm.ErrRecordNotFound)
	// Mock Create - 创建成功
	mockDao.On("Create", mock.AnythingOfType("*entity.Customer")).Return(nil)

	err := service.Register(req)
	assert.NoError(t, err)
	mockDao.AssertExpectations(t)
}

// TestUserService_Register_DuplicateUsername 测试用户名重复
func TestUserService_Register_DuplicateUsername(t *testing.T) {
	mockDao := new(MockCustomerDao)
	service := &UserService{
		customerDao: mockDao,
	}

	req := &v1.RegisterRequest{
		Username: "existinguser",
		Email:    "test@example.com",
		Password: "password123",
	}

	existingUser := &entity.Customer{
		ID:       1,
		Username: "existinguser",
	}

	// Mock GetByUsername - 用户名已存在
	mockDao.On("GetByUsername", req.Username).Return(existingUser, nil)

	err := service.Register(req)
	assert.Error(t, err)
	assert.Equal(t, "用户名已存在", err.Error())
	mockDao.AssertExpectations(t)
}

// TestUserService_Register_DuplicateEmail 测试邮箱重复
func TestUserService_Register_DuplicateEmail(t *testing.T) {
	mockDao := new(MockCustomerDao)
	service := &UserService{
		customerDao: mockDao,
	}

	req := &v1.RegisterRequest{
		Username: "newuser",
		Email:    "existing@example.com",
		Password: "password123",
	}

	existingUser := &entity.Customer{
		ID:    1,
		Email: "existing@example.com",
	}

	// Mock GetByUsername - 用户名不存在
	mockDao.On("GetByUsername", req.Username).Return(nil, gorm.ErrRecordNotFound)
	// Mock GetByEmail - 邮箱已存在
	mockDao.On("GetByEmail", req.Email).Return(existingUser, nil)

	err := service.Register(req)
	assert.Error(t, err)
	assert.Equal(t, "邮箱已存在", err.Error())
	mockDao.AssertExpectations(t)
}

// TestUserService_Login_Success 测试登录成功
func TestUserService_Login_Success(t *testing.T) {
	mockDao := new(MockCustomerDao)
	service := &UserService{
		customerDao: mockDao,
	}

	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	req := &v1.LoginRequest{
		Username: "testuser",
		Password: password,
	}

	existingUser := &entity.Customer{
		ID:           1,
		UUID:         uuid.New(),
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
		DisplayName:  "Test User",
		UserType:     "external",
		Status:       "active",
	}

	// Mock GetByUsername - 用户存在
	mockDao.On("GetByUsername", req.Username).Return(existingUser, nil)

	resp, err := service.Login(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, existingUser.ID, resp.User.ID)
	assert.Equal(t, existingUser.Username, resp.User.Username)
	mockDao.AssertExpectations(t)
}

// TestUserService_Login_WrongPassword 测试密码错误
func TestUserService_Login_WrongPassword(t *testing.T) {
	mockDao := new(MockCustomerDao)
	service := &UserService{
		customerDao: mockDao,
	}

	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	req := &v1.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	existingUser := &entity.Customer{
		ID:           1,
		Username:     "testuser",
		PasswordHash: string(hashedPassword),
		Status:       "active",
	}

	// Mock GetByUsername - 用户存在
	mockDao.On("GetByUsername", req.Username).Return(existingUser, nil)

	resp, err := service.Login(req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "用户名或密码错误", err.Error())
	mockDao.AssertExpectations(t)
}

// TestUserService_Login_UserNotFound 测试用户不存在
func TestUserService_Login_UserNotFound(t *testing.T) {
	mockDao := new(MockCustomerDao)
	service := &UserService{
		customerDao: mockDao,
	}

	req := &v1.LoginRequest{
		Username: "nonexistent",
		Password: "password123",
	}

	// Mock GetByUsername - 用户不存在
	mockDao.On("GetByUsername", req.Username).Return(nil, gorm.ErrRecordNotFound)

	resp, err := service.Login(req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "用户名或密码错误", err.Error())
	mockDao.AssertExpectations(t)
}

// TestUserService_Login_UserSuspended 测试用户被禁用
func TestUserService_Login_UserSuspended(t *testing.T) {
	mockDao := new(MockCustomerDao)
	service := &UserService{
		customerDao: mockDao,
	}

	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	req := &v1.LoginRequest{
		Username: "testuser",
		Password: password,
	}

	existingUser := &entity.Customer{
		ID:           1,
		Username:     "testuser",
		PasswordHash: string(hashedPassword),
		Status:       "suspended",
	}

	// Mock GetByUsername - 用户存在但被禁用
	mockDao.On("GetByUsername", req.Username).Return(existingUser, nil)

	resp, err := service.Login(req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "用户已被禁用", err.Error())
	mockDao.AssertExpectations(t)
}

// TestUserService_GetUserInfo 测试获取用户信息
func TestUserService_GetUserInfo(t *testing.T) {
	mockDao := new(MockCustomerDao)
	service := &UserService{
		customerDao: mockDao,
	}

	userID := uint(1)
	existingUser := &entity.Customer{
		ID:          userID,
		UUID:        uuid.New(),
		Username:    "testuser",
		Email:       "test@example.com",
		DisplayName: "Test User",
		AvatarURL:   "http://example.com/avatar.jpg",
		UserType:    "external",
		Status:      "active",
	}

	// Mock GetByID - 用户存在
	mockDao.On("GetByID", userID).Return(existingUser, nil)

	userInfo, err := service.GetUserInfo(userID)
	assert.NoError(t, err)
	assert.NotNil(t, userInfo)
	assert.Equal(t, existingUser.ID, userInfo.ID)
	assert.Equal(t, existingUser.Username, userInfo.Username)
	assert.Equal(t, existingUser.Email, userInfo.Email)
	mockDao.AssertExpectations(t)
}

// TestUserService_UpdateUser 测试更新用户信息
func TestUserService_UpdateUser(t *testing.T) {
	mockDao := new(MockCustomerDao)
	service := &UserService{
		customerDao: mockDao,
	}

	userID := uint(1)
	existingUser := &entity.Customer{
		ID:          userID,
		Username:    "testuser",
		Email:       "test@example.com",
		DisplayName: "Old Name",
		AvatarURL:   "http://example.com/old.jpg",
	}

	req := &v1.UpdateUserRequest{
		Nickname: "New Name",
		Avatar:   "http://example.com/new.jpg",
	}

	// Mock GetByID - 用户存在
	mockDao.On("GetByID", userID).Return(existingUser, nil)
	// Mock Update - 更新成功
	mockDao.On("Update", mock.AnythingOfType("*entity.Customer")).Return(nil)

	err := service.UpdateUser(userID, req)
	assert.NoError(t, err)
	mockDao.AssertExpectations(t)
}
