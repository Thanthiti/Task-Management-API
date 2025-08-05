package repository_test

import (
	"log"
	// "mymodule/config"
	"mymodule/internal/task/model"
	"mymodule/internal/task/repository"
	"mymodule/pkg/logger"
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	logger.InitLogger()
	os.Exit(m.Run())
}


func WithRollback(db *gorm.DB, t *testing.T, testFunc func(tx *gorm.DB)) {
	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("failed to begin transaction: %v", tx.Error)
	}
	
	defer func() {
		err := tx.Rollback().Error
		if err != nil && err != gorm.ErrInvalidTransaction {
			t.Fatalf("failed to rollback transaction: %v", err)
		}
		}()

	testFunc(tx)
}

// PostgreSQL
// func setupTestDB() *gorm.DB {
// 	db := config.InitDB(".env.test")
// 	err := db.AutoMigrate(&model.Task{})
// 	if err != nil {
// 		log.Fatalf("failed to migrate: %v", err)
// 	}
// 	return db
// }

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	err = db.AutoMigrate(&model.Task{})
	if err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}
	return db
}

func TestSaveTask(t *testing.T) {
	db := setupTestDB()

	WithRollback(db, t, func(tx *gorm.DB) {
		repo := repository.NewGormTaskRepository(tx)

		task := model.Task{
			Title:       "Test Task",
			Description: "Test description",
			UserID:      1,
		}
		err := repo.Save(task)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		var check model.Task
		tx.First(&check, "title = ?", "Test Task")
		if check.Title != "Test Task" {
			t.Errorf("expected title 'Test Task', got: %v", check.Title)
		}
	})
}

func TestSaveTask_DBError(t *testing.T) {
	db := setupTestDB()
	sqlDB, _ := db.DB()
	sqlDB.Close()

	repo := repository.NewGormTaskRepository(db)
	task := model.Task{Title: "Should Fail", UserID: 1}

	err := repo.Save(task)
	if err == nil {
		t.Errorf("expected error due to closed DB, got nil")
	}
}

func TestFindTaskByIDAndUser(t *testing.T) {
	db := setupTestDB()
	repo := repository.NewGormTaskRepository(db)

	task := model.Task{
		Title:  "Find Me",
		UserID: 1,
	}
	db.Create(&task)

	found, err := repo.FindByIDAndUser(task.ID, task.UserID)
	if err != nil || found == nil {
		t.Fatalf("expected to find task, got error: %v", err)
	}
	if found.Title != "Find Me" {
		t.Errorf("expected 'Find Me', got: %v", found.Title)
	}
}

func TestFindTaskByIDAndUser_NotFound(t *testing.T) {
	db := setupTestDB()
	repo := repository.NewGormTaskRepository(db)

	found, err := repo.FindByIDAndUser(9999, 1)
	if err == nil || found != nil {
		t.Errorf("expected not found error and nil task, got: %v", err)
	}
}

func TestUpdateTask(t *testing.T) {
	db := setupTestDB()
	WithRollback(db, t, func(tx *gorm.DB) {
		repo := repository.NewGormTaskRepository(tx)

		task := model.Task{Title: "Old Title", UserID: 1}
		tx.Create(&task)

		task.Title = "Updated Title"
		err := repo.Update(&task)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		var check model.Task
		tx.First(&check, task.ID)
		if check.Title != "Updated Title" {
			t.Errorf("expected 'Updated Title', got: %v", check.Title)
		}
	})
}

func TestUpdateTask_DBError(t *testing.T) {
	db := setupTestDB()
	task := model.Task{Title: "Before", UserID: 1}
	db.Create(&task)

	sqlDB, _ := db.DB()
	sqlDB.Close()

	repo := repository.NewGormTaskRepository(db)
	task.Title = "After"
	err := repo.Update(&task)
	if err == nil {
		t.Errorf("expected error due to closed DB, got nil")
	}
}

func TestDeleteTask(t *testing.T) {
	db := setupTestDB()
	WithRollback(db, t, func(tx *gorm.DB) {
		repo := repository.NewGormTaskRepository(tx)

		task := model.Task{Title: "To Delete", UserID: 1}
		tx.Create(&task)

		err := repo.Delete(task.ID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		var count int64
		tx.Model(&model.Task{}).Where("id = ?", task.ID).Count(&count)
		if count != 0 {
			t.Errorf("expected task to be deleted")
		}
	})
}

func TestDeleteTask_NotExist(t *testing.T) {
	db := setupTestDB()
	repo := repository.NewGormTaskRepository(db)

	err := repo.Delete(9999)
	if err != nil && err != gorm.ErrRecordNotFound {
		t.Errorf("expected no error or ErrRecordNotFound, got: %v", err)
	}
}
