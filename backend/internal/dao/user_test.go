package dao

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setupMockDB 创建Mock数据库连接
func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *sql.DB) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("创建sqlmock失败: %v", err)
	}

	dialector := postgres.New(postgres.Config{
		Conn:       sqlDB,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("创建gorm DB失败: %v", err)
	}

	return gormDB, mock, sqlDB
}

// TestCustomerDao_Create 测试创建用户
func TestCustomerDao_Create(t *testing.T) {
	gormDB, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	dao := &CustomerDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		customer := &entity.Customer{
			Username:     "testuser",
			Email:        "test@example.com",
			PasswordHash: "hashedpassword",
			DisplayName:  "Test User",
			UserType:     "external",
			AccountType:  "individual",
			Status:       "active",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "customers"`)).
			WithArgs(
				customer.Username,
				customer.Email,
				customer.PasswordHash,
				customer.DisplayName,
				customer.AvatarURL,
				customer.Phone,
				customer.FullName,
				customer.Company,
				customer.UserType,
				customer.AccountType,
				customer.Status,
				customer.EmailVerified,
				customer.PhoneVerified,
				sqlmock.AnyArg(), // last_login_at
				sqlmock.AnyArg(), // created_at
				sqlmock.AnyArg(), // updated_at
				sqlmock.AnyArg(), // deleted_at
			).
			WillReturnRows(sqlmock.NewRows([]string{"uuid", "id"}).AddRow(uuid.New(), 1))
		mock.ExpectCommit()

		err := dao.Create(customer)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		customer := &entity.Customer{
			Username:     "testuser",
			Email:        "test@example.com",
			PasswordHash: "hashedpassword",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "customers"`)).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := dao.Create(customer)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestCustomerDao_GetByID 测试根据ID获取用户
func TestCustomerDao_GetByID(t *testing.T) {
	gormDB, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	dao := &CustomerDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		userID := uint(1)
		testUUID := uuid.New()
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "uuid", "username", "email", "password_hash", "display_name",
			"avatar_url", "phone", "full_name", "company", "user_type",
			"account_type", "status", "email_verified", "phone_verified",
			"last_login_at", "created_at", "updated_at", "deleted_at",
		}).AddRow(
			userID, testUUID, "testuser", "test@example.com", "hashedpassword",
			"Test User", "", "", "", "", "external", "individual", "active",
			false, false, nil, now, now, nil,
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers" WHERE id = $1`)).
			WithArgs(userID, 1).
			WillReturnRows(rows)

		customer, err := dao.GetByID(userID)
		assert.NoError(t, err)
		assert.NotNil(t, customer)
		assert.Equal(t, userID, customer.ID)
		assert.Equal(t, "testuser", customer.Username)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotFound", func(t *testing.T) {
		userID := uint(999)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers" WHERE id = $1`)).
			WithArgs(userID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		customer, err := dao.GetByID(userID)
		assert.Error(t, err)
		assert.Nil(t, customer)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		userID := uint(1)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers" WHERE id = $1`)).
			WithArgs(userID, 1).
			WillReturnError(sql.ErrConnDone)

		customer, err := dao.GetByID(userID)
		assert.Error(t, err)
		assert.Nil(t, customer)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestCustomerDao_GetByUsername 测试根据用户名获取用户
func TestCustomerDao_GetByUsername(t *testing.T) {
	gormDB, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	dao := &CustomerDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		username := "testuser"
		testUUID := uuid.New()
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "uuid", "username", "email", "password_hash", "display_name",
			"avatar_url", "phone", "full_name", "company", "user_type",
			"account_type", "status", "email_verified", "phone_verified",
			"last_login_at", "created_at", "updated_at", "deleted_at",
		}).AddRow(
			1, testUUID, username, "test@example.com", "hashedpassword",
			"Test User", "", "", "", "", "external", "individual", "active",
			false, false, nil, now, now, nil,
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers" WHERE username = $1`)).
			WithArgs(username, 1).
			WillReturnRows(rows)

		customer, err := dao.GetByUsername(username)
		assert.NoError(t, err)
		assert.NotNil(t, customer)
		assert.Equal(t, username, customer.Username)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotFound", func(t *testing.T) {
		username := "nonexistent"

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers" WHERE username = $1`)).
			WithArgs(username, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		customer, err := dao.GetByUsername(username)
		assert.Error(t, err)
		assert.Nil(t, customer)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestCustomerDao_GetByEmail 测试根据邮箱获取用户
func TestCustomerDao_GetByEmail(t *testing.T) {
	gormDB, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	dao := &CustomerDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		email := "test@example.com"
		testUUID := uuid.New()
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "uuid", "username", "email", "password_hash", "display_name",
			"avatar_url", "phone", "full_name", "company", "user_type",
			"account_type", "status", "email_verified", "phone_verified",
			"last_login_at", "created_at", "updated_at", "deleted_at",
		}).AddRow(
			1, testUUID, "testuser", email, "hashedpassword",
			"Test User", "", "", "", "", "external", "individual", "active",
			false, false, nil, now, now, nil,
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers" WHERE email = $1`)).
			WithArgs(email, 1).
			WillReturnRows(rows)

		customer, err := dao.GetByEmail(email)
		assert.NoError(t, err)
		assert.NotNil(t, customer)
		assert.Equal(t, email, customer.Email)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotFound", func(t *testing.T) {
		email := "nonexistent@example.com"

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers" WHERE email = $1`)).
			WithArgs(email, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		customer, err := dao.GetByEmail(email)
		assert.Error(t, err)
		assert.Nil(t, customer)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestCustomerDao_Update 测试更新用户
func TestCustomerDao_Update(t *testing.T) {
	gormDB, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	dao := &CustomerDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		customer := &entity.Customer{
			ID:           1,
			UUID:         uuid.New(),
			Username:     "testuser",
			Email:        "test@example.com",
			PasswordHash: "hashedpassword",
			DisplayName:  "Updated User",
			UserType:     "external",
			AccountType:  "individual",
			Status:       "active",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "customers"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := dao.Update(customer)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		customer := &entity.Customer{
			ID:           1,
			Username:     "testuser",
			Email:        "test@example.com",
			PasswordHash: "hashedpassword",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "customers"`)).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := dao.Update(customer)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestCustomerDao_Delete 测试删除用户
func TestCustomerDao_Delete(t *testing.T) {
	gormDB, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	dao := &CustomerDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		userID := uint(1)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "customers" SET "deleted_at"=$1 WHERE "customers"."id" = $2`)).
			WithArgs(sqlmock.AnyArg(), userID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := dao.Delete(userID)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		userID := uint(1)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "customers" SET "deleted_at"=$1 WHERE "customers"."id" = $2`)).
			WithArgs(sqlmock.AnyArg(), userID).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := dao.Delete(userID)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestCustomerDao_List 测试分页列表
func TestCustomerDao_List(t *testing.T) {
	gormDB, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	dao := &CustomerDao{db: gormDB}

	t.Run("Success", func(t *testing.T) {
		page := 1
		pageSize := 10
		testUUID1 := uuid.New()
		testUUID2 := uuid.New()
		now := time.Now()

		// Mock Count query
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "customers"`)).
			WillReturnRows(countRows)

		// Mock List query
		rows := sqlmock.NewRows([]string{
			"id", "uuid", "username", "email", "password_hash", "display_name",
			"avatar_url", "phone", "full_name", "company", "user_type",
			"account_type", "status", "email_verified", "phone_verified",
			"last_login_at", "created_at", "updated_at", "deleted_at",
		}).
			AddRow(1, testUUID1, "user1", "user1@example.com", "hash1", "User 1",
				"", "", "", "", "external", "individual", "active", false, false, nil, now, now, nil).
			AddRow(2, testUUID2, "user2", "user2@example.com", "hash2", "User 2",
				"", "", "", "", "external", "individual", "active", false, false, nil, now, now, nil)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers"`)).
			WithArgs(pageSize).
			WillReturnRows(rows)

		customers, total, err := dao.List(page, pageSize)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, customers, 2)
		assert.Equal(t, "user1", customers[0].Username)
		assert.Equal(t, "user2", customers[1].Username)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("CountError", func(t *testing.T) {
		page := 1
		pageSize := 10

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "customers"`)).
			WillReturnError(sql.ErrConnDone)

		customers, total, err := dao.List(page, pageSize)
		assert.Error(t, err)
		assert.Nil(t, customers)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ListError", func(t *testing.T) {
		page := 1
		pageSize := 10

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "customers"`)).
			WillReturnRows(countRows)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers"`)).
			WithArgs(pageSize).
			WillReturnError(sql.ErrConnDone)

		customers, total, err := dao.List(page, pageSize)
		assert.Error(t, err)
		assert.Nil(t, customers)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
