package repository_test

import (
	"log"
	// "mymodule/config"
	"mymodule/internal/user/model"
	"mymodule/internal/user/repository"
	"mymodule/pkg/logger"
	"os"
	"testing"

	// "gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Testlog(m *testing.M) {
	logger.InitLogger()
	os.Exit(m.Run())
}

func WithRollback(db *gorm.DB, t *testing.T, testFunc func(tx *gorm.DB)) {
	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("failed to begin transaction: %v", tx.Error)
	}

	// defer rollback เพื่อให้ rollback อัตโนมัติเมื่อฟังก์ชันนี้จบ
	defer func() {
		err := tx.Rollback().Error
		if err != nil && err != gorm.ErrInvalidTransaction {
			t.Fatalf("failed to rollback transaction: %v", err)
		}
	}()

	// เรียกฟังก์ชันทดสอบพร้อม transaction นี้
	testFunc(tx)
}

//// PostgreSQL
// func setupTestDB() *gorm.DB {
// 	db := config.InitDB(".env.test")
// 	err := db.AutoMigrate(&model.User{})
// 	if err != nil {
// 		log.Fatalf("failed to migrate: %v", err)
// 	}
// 	return db
// }

// Sqlite
func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}
	return db
}
func TestSaveUser(t *testing.T) {
	logger.InitLogger()
	db := setupTestDB()

	WithRollback(db, t, func(tx *gorm.DB) {
		repo := repository.NewGormUserRepository(tx)

		user := model.User{Name: "John", Email: "john@example.com", Password: "1234"}
		err := repo.Save(user)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		var check model.User
		tx.First(&check, "email = ?", "john@example.com")
		if check.Name != "John" {
			t.Errorf("expected name 'John', got: %v", check.Name)
		}
	})
}

func TestSaveUser_DBError(t *testing.T) {
	logger.InitLogger()
	db := setupTestDB()
	sqlDB, _ := db.DB()
	sqlDB.Close()

	repo := repository.NewGormUserRepository(db)
	user := model.User{Name: "Error", Email: "err@example.com", Password: "pw"}
	err := repo.Save(user)

	if err == nil {
		t.Errorf("expected error due to closed DB, got nil")
	}
}

func TestFindByEmail(t *testing.T) {
	logger.InitLogger()
	db := setupTestDB()
	repo := repository.NewGormUserRepository(db)

	// Set data
	db.Create(&model.User{Name: "Jane", Email: "jane@example.com", Password: "pw"})

	user, err := repo.FindByEmail("jane@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user == nil || user.Email != "jane@example.com" {
		t.Errorf("expected to find user, got: %v", user)
	}
}

func TestFindByEmail_NotFound(t *testing.T) {
	logger.InitLogger()
	db := setupTestDB()
	repo := repository.NewGormUserRepository(db)

	user, err := repo.FindByEmail("unknown@example.com")
	if err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
	if user != nil {
		t.Errorf("expected nil user, got: %v", user)
	}
}

func TestFindByID(t *testing.T) {
	logger.InitLogger()
	db := setupTestDB()
	repo := repository.NewGormUserRepository(db)

	u := model.User{Name: "A", Email: "a@example.com", Password: "pw"}
	db.Create(&u)

	found, err := repo.FindByID(u.ID)
	if err != nil || found == nil {
		t.Errorf("expected to find user, got error: %v", err)
	}
	if found.Email != "a@example.com" {
		t.Errorf("email mismatch: %v", found.Email)
	}
}
func TestFindByID_NotFound(t *testing.T) {
	logger.InitLogger()
	db := setupTestDB()
	repo := repository.NewGormUserRepository(db)

	user, err := repo.FindByID(99999)
	if err == nil || user != nil {
		t.Errorf("expected not found error and nil user, got: %v", err)
	}
}

func TestUpdateUser(t *testing.T) {
	logger.InitLogger()
	db := setupTestDB()
	WithRollback(db, t, func(tx *gorm.DB) {
		repo := repository.NewGormUserRepository(db)

		u := model.User{Name: "Old", Email: "old@example.com", Password: "pw"}
		db.Create(&u)

		u.Name = "Updated"
		err := repo.Update(u)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		var check model.User
		db.First(&check, u.ID)
		if check.Name != "Updated" {
			t.Errorf("expected Updated, got: %v", check.Name)
		}
	})
}
func TestUpdateUser_DBError(t *testing.T) {
	logger.InitLogger()
	db := setupTestDB()

	user := model.User{Name: "Before", Email: "update@example.com", Password: "pw"}
	db.Create(&user)

	sqlDB, _ := db.DB()
	sqlDB.Close()

	repo := repository.NewGormUserRepository(db)
	user.Name = "After"
	err := repo.Update(user)

	if err == nil {
		t.Errorf("expected error due to closed DB, got nil")
	}
}

func TestDeleteUser(t *testing.T) {
	logger.InitLogger()
	db := setupTestDB()
	WithRollback(db, t, func(tx *gorm.DB) {
		repo := repository.NewGormUserRepository(db)

		u := model.User{Name: "ToDelete", Email: "del@example.com", Password: "pw"}
		db.Create(&u)

		err := repo.Delete(u.ID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		var count int64
		db.Model(&model.User{}).Where("id = ?", u.ID).Count(&count)
		if count != 0 {
			t.Errorf("expected user to be deleted")
		}
	})
}

func TestDeleteUser_NotExist(t *testing.T) {
	logger.InitLogger()
	db := setupTestDB()
	repo := repository.NewGormUserRepository(db)

	err := repo.Delete(99999)

	if err != nil && err != gorm.ErrRecordNotFound {
		t.Errorf("expected no error or gorm.ErrRecordNotFound, got: %v", err)
	}
}
