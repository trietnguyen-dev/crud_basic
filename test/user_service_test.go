package test

import (
	"book-project/config"
	"book-project/daos"
	pb "book-project/protobuf/gen/go"
	"book-project/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

type UserSuite struct {
	suite.Suite
	conf    *config.Config
	userSvc *service.UserScv
	logger  *zap.Logger
}

// setup config
func (suite *UserSuite) SetupSuite() {
	suite.conf = config.GetConfig()

	dao, err := daos.NewDAO(suite.conf)
	if err != nil {
		suite.logger.Fatal("failed to init daos", zap.Error(err))
	}

	suite.userSvc = service.NewUserSvc(dao, suite.conf)
}

func (suite *UserSuite) TestCreateUser() {
	res_1 := rand.Int()
	s1 := strconv.FormatInt(int64(res_1), 10)
	createUser := struct {
		FullName    string
		PhoneNumber string
		Email       string
		Password    string
		RePassword  string
		CreatedAt   time.Time
	}{
		FullName:    "111111",
		PhoneNumber: "0919651145",
		Email:       s1 + "@gmail.com",
		Password:    "123",
		RePassword:  "123",
		CreatedAt:   time.Now(),
	}
	// Create a protobuf message
	req := &pb.RegisterUserReq{
		FullName:    createUser.FullName,
		PhoneNumber: createUser.PhoneNumber,
		Email:       createUser.Email,
		Password:    createUser.Password,
		RePassword:  createUser.RePassword,
	}
	// Perform the user registration
	res, err := suite.userSvc.RegisterUser(req)

	// Check for errors or conditions
	assert.NoError(suite.T(), err, "Failed to create user ")
	assert.NotNil(suite.T(), res, false)
}
func (suite *UserSuite) TestCreateUserWithExistEmail() {
	createUser := struct {
		FullName    string
		PhoneNumber string
		Email       string
		Password    string
		RePassword  string
		CreatedAt   time.Time
	}{
		FullName:    "111111",
		PhoneNumber: "0919651145",
		Email:       "123123@gmail.com",
		Password:    "123",
		RePassword:  "123",
		CreatedAt:   time.Now(),
	}

	// Create a protobuf message
	req := &pb.RegisterUserReq{
		FullName:    createUser.FullName,
		PhoneNumber: createUser.PhoneNumber,
		Email:       createUser.Email,
		Password:    createUser.Password,
		RePassword:  createUser.RePassword,
	}

	res, err := suite.userSvc.RegisterUser(req)

	assert.NoError(suite.T(), err, "Exist Email")
	assert.NotNil(suite.T(), res, false)
}
func (suite *UserSuite) TestCreateUserWithErrorPass() {
	res_1 := rand.Int()
	s1 := strconv.FormatInt(int64(res_1), 10)
	createUser := struct {
		FullName    string
		PhoneNumber string
		Email       string
		Password    string
		RePassword  string
		CreatedAt   time.Time
	}{
		FullName:    "111111",
		PhoneNumber: "0919651145",
		Email:       s1 + "@gmail.com",
		Password:    "123",
		RePassword:  "1234",
		CreatedAt:   time.Now(),
	}
	// Create a protobuf message
	req := &pb.RegisterUserReq{
		FullName:    createUser.FullName,
		PhoneNumber: createUser.PhoneNumber,
		Email:       createUser.Email,
		Password:    createUser.Password,
		RePassword:  createUser.RePassword,
	}

	res, err := suite.userSvc.RegisterUser(req)

	assert.NoError(suite.T(), err, "Failed to create user because of mistake re-password")
	assert.NotNil(suite.T(), res, false)
}
func (suite *UserSuite) TestUserLoginSuccess() {

	loginUser := struct {
		Email    string
		Password string
	}{
		Email:    "6554199082744098033@gmail.com",
		Password: "123",
	}
	// Hash the plaintext password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(loginUser.Password), 14)
	assert.NoError(suite.T(), err, "Error hashing password")

	// Create a protobuf message with the hashed password
	req := &pb.LoginUserReq{
		Email:    loginUser.Email,
		Password: string(hashedPassword), // Store the hashed password in the database
	}

	res, err := suite.userSvc.LoginUser(req)

	assert.NoError(suite.T(), err, "Error password")
	assert.NotNil(suite.T(), res, false)
}
func (suite *UserSuite) TestUserLoginError() {
	loginUser := struct {
		Email    string
		Password string
	}{
		Email:    "2744098033@gmail.com",
		Password: "123",
	}
	// Hash the plaintext password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(loginUser.Password), 14)
	assert.NoError(suite.T(), err, "Error hashing password")

	// Create a protobuf message with the hashed password
	req := &pb.LoginUserReq{
		Email:    loginUser.Email,
		Password: string(hashedPassword), // Store the hashed password in the database
	}

	res, err := suite.userSvc.LoginUser(req)

	assert.NoError(suite.T(), err, "Error password")
	assert.NotNil(suite.T(), res, false)
}
func (suite *UserSuite) TestGetUserSuccess() {
	loginFilter := "5225894317891256713@gmail.com"

	filter := &pb.UserInfo{
		Email: loginFilter,
	}

	// Create a protobuf message with the filter
	req := &pb.GetInfoUserReq{
		Filter: filter,
	}
	res, err := suite.userSvc.GetInfoUser(req)

	assert.NoError(suite.T(), err, "No information")
	assert.NotNil(suite.T(), res, "Response should not be nil") // Change the third argument to true
}
func (suite *UserSuite) TestGetUserFailed() {
	loginFilter := "52258943178912567d13@gmail.com"

	filter := &pb.UserInfo{
		Email: loginFilter,
	}

	// Create a protobuf message with the filter
	req := &pb.GetInfoUserReq{
		Filter: filter,
	}
	res, err := suite.userSvc.GetInfoUser(req)

	assert.NoError(suite.T(), err, "No information")
	assert.NotNil(suite.T(), res, "Response should not be nil") // Change the third argument to true
}
func (suite *UserSuite) TestUpdateUserSuccess() {

	userIdStr := "650d374453373cc78cff7afd"

	updateUser := struct {
		FullName    string
		PhoneNumber string
		Email       string
		UpdatedAt   time.Time
	}{
		FullName:    "triet nguyen",
		PhoneNumber: "0919651145",
		Email:       "minhtriet@gmail.com",
		UpdatedAt:   time.Now(),
	}

	// Create a protobuf message with the filter
	req := &pb.UpdateUserReq{
		Id:          userIdStr,
		FullName:    updateUser.FullName,
		Email:       updateUser.Email,
		PhoneNumber: updateUser.PhoneNumber,
	}
	res, err := suite.userSvc.UpdateUser(req)

	assert.NoError(suite.T(), err, "No find id")
	assert.NotNil(suite.T(), res, "Response should not be nil") // Change the third argument to true
}
func (suite *UserSuite) TestUpdateUserFailedWithPhoneNumber() {

	userIdStr := "650d374453373cc78cff7afd"

	updateUser := struct {
		FullName    string
		PhoneNumber string
		Email       string
		UpdatedAt   time.Time
	}{
		FullName:    "triet nguyen",
		PhoneNumber: "091961145",
		Email:       "minhtriet@gmail.com",
		UpdatedAt:   time.Now(),
	}

	// Create a protobuf message with the filter
	req := &pb.UpdateUserReq{
		Id:          userIdStr,
		FullName:    updateUser.FullName,
		Email:       updateUser.Email,
		PhoneNumber: updateUser.PhoneNumber,
	}
	res, err := suite.userSvc.UpdateUser(req)

	assert.NoError(suite.T(), err, "Number invalid")
	assert.NotNil(suite.T(), res, "Response should not be nil") // Change the third argument to true
}
func (suite *UserSuite) TestUpdateUserFailedWithEmail() {

	userIdStr := "650d374453373cc78cff7afd"

	updateUser := struct {
		FullName    string
		PhoneNumber string
		Email       string
		UpdatedAt   time.Time
	}{
		FullName:    "triet nguyen",
		PhoneNumber: "0919611145",
		Email:       "minhtrietmail.com",
		UpdatedAt:   time.Now(),
	}

	// Create a protobuf message with the filter
	req := &pb.UpdateUserReq{
		Id:          userIdStr,
		FullName:    updateUser.FullName,
		Email:       updateUser.Email,
		PhoneNumber: updateUser.PhoneNumber,
	}
	res, err := suite.userSvc.UpdateUser(req)

	assert.NoError(suite.T(), err, "Email invalid")
	assert.NotNil(suite.T(), res, "Response should not be nil") // Change the third argument to true
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}
