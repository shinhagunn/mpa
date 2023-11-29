package mongo_fx_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/shinhagunn/mpa/config"
	"github.com/shinhagunn/mpa/models"
	filters "github.com/shinhagunn/mpa/mongodb/fitlers"
	"github.com/shinhagunn/mpa/pkg/mongo_fx"
	"github.com/stretchr/testify/suite"
	"github.com/zsmartex/pkg/v2/utils"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func InitUser() *models.User {
	return &models.User{
		UID:       utils.GenerateUID(),
		Email:     "ha@gmail.com",
		Role:      models.UserRoleMember,
		CreatedAt: time.Now().Local(),
		UpdatedAt: time.Now().Local(),
	}
}

type UserRepositorySuite struct {
	suite.Suite
	app *fxtest.App

	userRepo mongo_fx.Repository[models.User]
}

func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserRepositorySuite))
}

func (suite *UserRepositorySuite) SetupSuite() {
	suite.app = fxtest.New(
		suite.T(),
		config.Module,
		mongo_fx.Module,
		fx.Invoke(func(db *mongo_fx.Mongo) {
			suite.userRepo = mongo_fx.NewRepository(db, models.User{})
		}),
	).RequireStart()
}

func (suite *UserRepositorySuite) TearDownSuite() {
	suite.app.RequireStop()
}

func (suite *UserRepositorySuite) TestCreateUser() {
	initedUser := InitUser()

	err := suite.userRepo.Create(context.TODO(), initedUser)
	suite.Require().NoError(err)

	fmt.Println(initedUser)
}

func (suite *UserRepositorySuite) TestCountUsers() {
	// Count all user exist
	count1, err := suite.userRepo.Count(context.TODO())
	suite.Require().NoError(err)

	log.Println(count1)

	// Count user have level > 4
	count2, err := suite.userRepo.Count(context.TODO(), filters.WithFieldEqual("email", "ha@gmail.com"))
	suite.Require().NoError(err)

	log.Println(count2)

	// Count user have level <= 2
	count3, err := suite.userRepo.Count(context.TODO(), filters.WithFieldLesThanOrEqualTo("level", 2))
	suite.Require().NoError(err)

	log.Println(count3)
}

func (suite *UserRepositorySuite) TestFindUsers() {
	// Find all users
	users1, err := suite.userRepo.Find(context.TODO(), nil)
	suite.Require().NoError(err)
	log.Println(users1)

	users2, err := suite.userRepo.Find(context.TODO(), nil, filters.WithFieldEqual("email", "ha@gmail.com"))
	suite.Require().NoError(err)
	log.Println(users2)

	users3, err := suite.userRepo.Find(context.TODO(), nil, filters.WithFieldNotIn("role", []string{"f"}))
	suite.Require().NoError(err)
	log.Println(users3)

	users4, err := suite.userRepo.Find(context.TODO(), nil, filters.WithFieldLessThan("created_at", time.Now()))
	suite.Require().NoError(err)
	log.Println(users4)

	users5, err := suite.userRepo.Find(context.TODO(), nil, filters.WithFieldNotEqual("role", "admin"))
	suite.Require().NoError(err)
	log.Println(users5)

	users6, err := suite.userRepo.Find(context.TODO(), nil, filters.WithFieldNotIn("role", []string{"anonymus", "member"}))
	suite.Require().NoError(err)
	log.Println(users6)

	// Find users with field like
	users7, err := suite.userRepo.Find(context.TODO(), nil, filters.WithFieldLike("role", "m"))
	suite.Require().NoError(err)
	log.Println(users7)

	// Find users with is null
	users8, err := suite.userRepo.Find(context.TODO(), nil, filters.WithFieldIsNull("username"))
	suite.Require().NoError(err)
	log.Println(users8)
}

func (suite *UserRepositorySuite) TestFirstUsers() {
	// Find user with ID
	user, err := suite.userRepo.First(context.TODO(), filters.WithID("656442e608b5eb5270875c13"))
	suite.Require().NoError(err)
	log.Println(user)
}

func (suite *UserRepositorySuite) TestLastUsers() {
	// Find user with username
	user, err := suite.userRepo.Last(context.TODO(), filters.WithFieldEqual("role", "admin"))
	suite.Require().NoError(err)
	log.Println(user)
}

func (suite *UserRepositorySuite) TestFirstOrCreateUser() {
	var user models.User
	err := suite.userRepo.FirstOrCreate(
		context.TODO(),
		&user,
		&models.User{
			UID:       utils.GenerateUID(),
			Email:     "test@gmail.com",
			Role:      models.UserRoleMember,
			CreatedAt: time.Now().Local(),
			UpdatedAt: time.Now().Local(),
		},
		filters.WithFieldEqual("role", "admin"),
	)
	suite.Require().NoError(err)

	log.Println(user)
}

func (suite *UserRepositorySuite) TestUpdateUser() {
	user, err := suite.userRepo.First(context.TODO(), filters.WithID("656442e608b5eb5270875c13"))
	suite.Require().NoError(err)

	userUpdate := make(map[string]interface{})
	userUpdate["role"] = "ahihi"

	err = suite.userRepo.Updates(context.TODO(), user, userUpdate)
	suite.Require().NoError(err)

	log.Println(user)
}

func (suite *UserRepositorySuite) TestDeleteUser() {
	user, err := suite.userRepo.First(context.TODO(), filters.WithFieldEqual("role", "ahihi"))
	suite.Require().NoError(err)

	err = suite.userRepo.Delete(context.TODO(), filters.WithFieldEqual("_id", user.ID))
	suite.Require().NoError(err)
}
