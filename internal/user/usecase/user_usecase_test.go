package usecase_test

import (
	"errors"
	"os"
	"testing"

	"mymodule/internal/user/model"
	"mymodule/internal/user/usecase"
	"mymodule/pkg/logger"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

// Mock UserRepository
type MockUserRepo struct {
	usersByEmail map[string]*model.User
	usersByID    map[uint]*model.User

	SaveErr        error
	FindByEmailErr error
	FindByIDErr    error
	UpdateErr      error
	DeleteErr      error
}

func (m *MockUserRepo) Save(user model.User) error {
	if m.SaveErr != nil {
		return m.SaveErr
	}
	m.usersByEmail[user.Email] = &user
	m.usersByID[user.ID] = &user
    return nil
}

func (m *MockUserRepo) FindByEmail(email string) (*model.User, error) {
	if m.FindByEmailErr != nil {
		return nil, m.FindByEmailErr
	}
	user, ok := m.usersByEmail[email]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return user, nil
}

func (m *MockUserRepo) FindByID(userID uint) (*model.User, error) {
	if m.FindByIDErr != nil {
		return nil, m.FindByIDErr
	}
	user, ok := m.usersByID[userID] 
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return user, nil
}

func (m *MockUserRepo) Update(user model.User) error {
	if m.UpdateErr != nil {
		return m.UpdateErr
	}
	m.usersByEmail[user.Email] = &user
	m.usersByID[user.ID] = &user
	return nil  
}

func (m *MockUserRepo) Delete(userID uint) error {
	if m.DeleteErr != nil {
		return m.DeleteErr
	}
	user, ok := m.usersByID[userID] 
	if !ok {
		return gorm.ErrRecordNotFound
	}
	delete(m.usersByID, userID)
	delete(m.usersByEmail, user.Email) 
	return nil
}

// Mock CryptoService
type MockCryptoService struct {
    HashErr error
}

func (m *MockCryptoService) HashedPassword(pw string) (string, error) {
    if m.HashErr != nil {
        return "", m.HashErr  
    }
    return "hashedpassword", nil  
}

func (m *MockCryptoService) ComparePassword(hash string, pw string) bool {
    return pw == "correctpassword"
}

// Mock TokenService
type MockTokenService struct {
    TokenToReturn string
    ErrToReturn   error
}

func (m *MockTokenService) GenerateToken(userID uint) (string, error) {
    if m.ErrToReturn != nil {
        return "", m.ErrToReturn
    }
    return m.TokenToReturn, nil
}

func (m *MockTokenService) VerifyToken(token string) (*jwt.Token, error) {
    if token == "validtoken" {
        return nil, nil
    }
    return nil, errors.New("invalid token")
}

func Testlog(m *testing.M) {
    logger.InitLogger() 
    os.Exit(m.Run())
}

func TestUserUsecase_Register(t *testing.T) {
    logger.InitLogger()
	mockRepo := &MockUserRepo{
	usersByEmail: make(map[string]*model.User),
	usersByID:    make(map[uint]*model.User),
    }
    mockCrypto := &MockCryptoService{}
    mockToken := &MockTokenService{}

    uc := usecase.NewUserUsecase(mockRepo, mockCrypto, mockToken)

    // existing email
    mockRepo.usersByEmail["existing@example.com"] = &model.User{ID: 1, Email: "existing@example.com"}
    err := uc.Register(model.User{Email: "existing@example.com", Password: "123456"})
    if err == nil || err.Error() != "email already registered" {
        t.Errorf("expected email already registered error, got %v", err)
    }

    // Case hash password error
    mockCrypto.HashErr = errors.New("hash error")
    err = uc.Register(model.User{Email: "new@example.com", Password: "123456"})
    if err == nil {
        t.Errorf("expected hash error, got nil")
    }
    // Cleanup  
    mockCrypto.HashErr = nil 
    
    // Case save error
    mockRepo.SaveErr = errors.New("db save error")
    err = uc.Register(model.User{Email: "new@example.com", Password: "123456"})
    if err == nil {
        t.Errorf("expected save error, got nil")
    }
    // Cleanup 
    mockRepo.SaveErr = nil 

    // Case sucess
    err = uc.Register(model.User{Email: "new@example.com", Password: "123456"})
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
}

func TestUserUsecase_Login(t *testing.T) {
	logger.InitLogger()

	mockRepo := &MockUserRepo{
	usersByEmail: make(map[string]*model.User),
	usersByID:    make(map[uint]*model.User),
    }
	mockCrypto := &MockCryptoService{}
	mockToken := &MockTokenService{TokenToReturn: "mocktoken"}
	uc := usecase.NewUserUsecase(mockRepo, mockCrypto, mockToken)

	// Case can't find my email
	_, err := uc.Login("notfound@example.com", "123456")
	if err == nil || err.Error() != "invalid credentials" {
		t.Errorf("expected invalid credentials, got: %v", err)
	}

	// Case user nil 
	mockRepo.FindByEmailErr = nil
	mockRepo.usersByEmail["niluser@example.com"] = nil
	_, err = uc.Login("niluser@example.com", "123456")
	if err == nil || err.Error() != "invalid credentials" {
		t.Errorf("expected error for nil user, got: %v", err)
	}
	delete(mockRepo.usersByEmail, "niluser@example.com") // cleanup

	// Case wrong password 
	mockRepo.usersByEmail["user@example.com"] = &model.User{ID: 1, Email: "user@example.com", Password: "hashedpassword"}
	_, err = uc.Login("user@example.com", "wrongpassword")
	if err == nil || err.Error() != "invalid credentials" {
		t.Errorf("expected error for wrong password, got: %v", err)
	}

	// Case gen token  fail
	mockCrypto = &MockCryptoService{} // reset
	mockToken.ErrToReturn = errors.New("token error")
	uc = usecase.NewUserUsecase(mockRepo, mockCrypto, mockToken)
	_, err = uc.Login("user@example.com", "correctpassword")
	if err == nil || err.Error() != "internal server error" {
		t.Errorf("expected token error, got: %v", err)
	}

	// Case login success
	mockToken.ErrToReturn = nil
	uc = usecase.NewUserUsecase(mockRepo, mockCrypto, mockToken)
	token, err := uc.Login("user@example.com", "correctpassword")
	if err != nil || token != "mocktoken" {
		t.Errorf("expected token, got token: %v, err: %v", token, err)
	}
}

func TestUserUsecase_UpdateUser(t *testing.T) {
	logger.InitLogger()

	mockRepo := &MockUserRepo{
		usersByEmail: make(map[string]*model.User),
		usersByID:    make(map[uint]*model.User),
	}
	mockCrypto := &MockCryptoService{}
	mockToken := &MockTokenService{}
	uc := usecase.NewUserUsecase(mockRepo, mockCrypto, mockToken)

	// Set data: 
	userOld := &model.User{ID: 1, Email: "old@example.com", Name: "Old"}
	mockRepo.usersByEmail[userOld.Email] = userOld
	mockRepo.usersByID[userOld.ID] = userOld

	// 1. Can't find user by ID (mock error forced)
	mockRepo.FindByIDErr = gorm.ErrRecordNotFound
	err := uc.UpdateUser(model.User{ID: 99, Email: "x@example.com"})
	if err == nil || err.Error() != "user not found" {
		t.Errorf("expected not found error, got: %v", err)
	}
	mockRepo.FindByIDErr = nil

	// 2. nil user (simulate user missing by deleting from both maps)
	delete(mockRepo.usersByEmail, "old@example.com")
	delete(mockRepo.usersByID, 1)
	err = uc.UpdateUser(model.User{ID: 1, Email: "old@example.com"})
	if err == nil || err.Error() != "user not found" {
		t.Errorf("expected nil user error, got: %v", err)
	}
	// Add back for next tests
	mockRepo.usersByEmail[userOld.Email] = userOld
	mockRepo.usersByID[userOld.ID] = userOld

	// 3. Email is duplicated by another user
	userTaken := &model.User{ID: 2, Email: "taken@example.com"}
	mockRepo.usersByEmail[userTaken.Email] = userTaken
	mockRepo.usersByID[userTaken.ID] = userTaken

	err = uc.UpdateUser(model.User{ID: 1, Email: "taken@example.com", Name: "New"})
	if err == nil || err.Error() != "email already registered" {
		t.Errorf("expected email already registered, got: %v", err)
	}

	// 4. Update error
	mockRepo.UpdateErr = errors.New("update fail")
	err = uc.UpdateUser(model.User{ID: 1, Email: "old@example.com", Name: "New"})
	if err == nil || err.Error() != "update fail" {
		t.Errorf("expected update fail error, got: %v", err)
	}
	mockRepo.UpdateErr = nil

	// 5. Success
	err = uc.UpdateUser(model.User{ID: 1, Email: "updated@example.com", Name: "Updated"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}


func TestUserUsecase_DeleteUser(t *testing.T) {
	logger.InitLogger()

	mockRepo := &MockUserRepo{
		usersByEmail: make(map[string]*model.User),
		usersByID:    make(map[uint]*model.User),
	}
	mockCrypto := &MockCryptoService{}
	mockToken := &MockTokenService{}
	uc := usecase.NewUserUsecase(mockRepo, mockCrypto, mockToken)

	// Set data
	user := &model.User{ID: 1, Email: "user@example.com"}
	mockRepo.usersByEmail[user.Email] = user
	mockRepo.usersByID[user.ID] = user

	// 1. Can't find user by ID
	mockRepo.FindByIDErr = gorm.ErrRecordNotFound
	err := uc.DeleteUser(2)
	if err == nil || err.Error() != "user not found" {
		t.Errorf("expected not found error, got: %v", err)
	}
	mockRepo.FindByIDErr = nil

	// 2. user is nil
	delete(mockRepo.usersByID, 1)
	delete(mockRepo.usersByEmail, "user@example.com")
	err = uc.DeleteUser(1)
	if err == nil || err.Error() != "user not found" {
		t.Errorf("expected nil user error, got: %v", err)
	}

	// Return
	mockRepo.usersByEmail[user.Email] = user
	mockRepo.usersByID[user.ID] = user

	// 3. Fail delete
	mockRepo.DeleteErr = errors.New("cannot delete")
	err = uc.DeleteUser(1)
	if err == nil || err.Error() != "cannot delete" {
		t.Errorf("expected delete error, got: %v", err)
	}
	mockRepo.DeleteErr = nil

	// 4. Success delete
	err = uc.DeleteUser(1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

