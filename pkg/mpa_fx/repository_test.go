package mpa_fx_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/shinhagunn/mpa/config"
	"github.com/shinhagunn/mpa/filters"
	"github.com/shinhagunn/mpa/models"
	"github.com/shinhagunn/mpa/pkg/mpa_fx"
	"github.com/stretchr/testify/suite"
	"github.com/zsmartex/pkg/v2/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func InitUser() *models.User {
	return &models.User{
		UID:       utils.GenerateUID(),
		Email:     "ga@gmail.com",
		Role:      models.UserRoleAdmin,
		CreatedAt: time.Now().Local(),
		UpdatedAt: time.Now().Local(),
	}
}

type UserRepositorySuite struct {
	suite.Suite
	app *fxtest.App

	userRepo mpa_fx.Repository[models.User]
}

func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserRepositorySuite))
}

func (suite *UserRepositorySuite) SetupSuite() {
	suite.app = fxtest.New(
		suite.T(),
		config.Module,
		mpa_fx.Module,
		fx.Invoke(func(db *mongo.Database) {
			suite.userRepo = mpa_fx.NewRepository(db, models.User{})
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
	count2, err := suite.userRepo.Count(context.TODO(), filters.WithFieldGreaterThan("level", 4))
	suite.Require().NoError(err)

	log.Println(count2)

	// Count user have level <= 2
	count3, err := suite.userRepo.Count(context.TODO(), filters.WithFieldLesThanOrEqualTo("level", 2))
	suite.Require().NoError(err)

	log.Println(count3)
}

func (suite *UserRepositorySuite) TestFindUsers() {
	// Find all users
	users1, err := suite.userRepo.Find(context.TODO(), []filters.Filter{})
	suite.Require().NoError(err)
	log.Println(users1)

	// Find users with username = 'Ha'
	users2, err := suite.userRepo.Find(context.TODO(), []filters.Filter{
		filters.WithFieldEqual("username", "Ha"),
	})
	suite.Require().NoError(err)
	log.Println(users2)

	// Find users with role in ['anonymus', 'member']
	users3, err := suite.userRepo.Find(context.TODO(), []filters.Filter{
		filters.WithFieldIn("role", []string{"anonymus", "member"}),
	})
	suite.Require().NoError(err)
	log.Println(users3)

	// Find users with time.Now() greater than created_at
	users4, err := suite.userRepo.Find(context.TODO(), []filters.Filter{
		filters.WithFieldLessThan("created_at", time.Now()),
	})
	suite.Require().NoError(err)
	log.Println(users4)

	// Find users with not equal
	users5, err := suite.userRepo.Find(context.TODO(), []filters.Filter{
		filters.WithFieldNotEqual("username", "Ga"),
	})
	suite.Require().NoError(err)
	log.Println(users5)

	// Find users with role not in ['anonymus', 'member']
	users6, err := suite.userRepo.Find(context.TODO(), []filters.Filter{
		filters.WithFieldNotIn("role", []string{"anonymus", "member"}),
	})
	suite.Require().NoError(err)
	log.Println(users6)

	// Find users with field like
	users7, err := suite.userRepo.Find(context.TODO(), []filters.Filter{
		filters.WithFieldLike("username", "H"),
	})
	suite.Require().NoError(err)
	log.Println(users7)

	// Find users with is null
	users8, err := suite.userRepo.Find(context.TODO(), []filters.Filter{
		filters.WithFieldNotNull("username"),
	})
	suite.Require().NoError(err)
	log.Println(users8)
}

func (suite *UserRepositorySuite) TestFirstUsers() {
	// Find user with ID
	user, err := suite.userRepo.First(context.TODO(), filters.WithID("65603744fc427f07e4e4cf24"))
	suite.Require().NoError(err)
	log.Println(user)
}

func (suite *UserRepositorySuite) TestLastUsers() {
	// Find user with username
	user, err := suite.userRepo.Last(context.TODO(), filters.WithFieldEqual("role", "a"))
	suite.Require().NoError(err)
	log.Println(user)
}

func (suite *UserRepositorySuite) TestFirstOrCreateUser() {
	user, err := suite.userRepo.FirstOrCreate(
		context.TODO(),
		&models.User{
			Email: "ga123@gmail.com",
		},
		filters.WithFieldEqual("role", "asd"),
	)
	suite.Require().NoError(err)

	log.Println(user)
}

func (suite *UserRepositorySuite) TestUpdateUser() {
	user, err := suite.userRepo.First(context.TODO(), filters.WithFieldEqual("role", "member"))
	suite.Require().NoError(err)

	userUpdate := make(map[string]interface{})
	userUpdate["role"] = "ahihi"

	err = suite.userRepo.Updates(context.TODO(), user, userUpdate)
	suite.Require().NoError(err)
}

func (suite *UserRepositorySuite) TestDeleteUser() {
	user, err := suite.userRepo.First(context.TODO(), filters.WithFieldEqual("role", "ahihi"))
	suite.Require().NoError(err)

	err = suite.userRepo.DeleteByID(context.TODO(), user.ID)
	suite.Require().NoError(err)
}
