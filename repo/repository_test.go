package repo_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/shinhagunn/mpa/config"
	"github.com/shinhagunn/mpa/filters"
	"github.com/shinhagunn/mpa/models"
	"github.com/shinhagunn/mpa/pkg/mpa_fx"
	"github.com/shinhagunn/mpa/repo"
	"github.com/stretchr/testify/suite"
	"github.com/zsmartex/pkg/v2/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func InitUser() *models.User {
	return &models.User{
		UID:            utils.GenerateUID(),
		Username:       "Ga",
		Email:          "ga@gmail.com",
		PasswordDigest: "$2a$10$VVTqd7vt.YtgEwQ.nb3OPea.a9L9FYp8siRLYkP9.lCOsigL8je7u",
		Level:          int64(gofakeit.Number(1, 10)),
		OTP:            false,
		Role:           models.UserRoleAdmin,
		State:          models.UserStateActive,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

type UserRepositorySuite struct {
	suite.Suite
	app *fxtest.App

	userRepo repo.Repository
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
			suite.userRepo = repo.New(db, &models.User{})
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
	var users1 []models.User
	err := suite.userRepo.Find(context.TODO(), &users1, []filters.Filter{})
	suite.Require().NoError(err)
	log.Println(users1)

	// Find users with username = 'Ha'
	var users2 []models.User
	err = suite.userRepo.Find(context.TODO(), &users2, []filters.Filter{
		filters.WithFieldEqual("username", "Ha"),
	})
	suite.Require().NoError(err)
	log.Println(users2)

	// Find users with role in ['anonymus', 'member']
	var users3 []models.User
	err = suite.userRepo.Find(context.TODO(), &users3, []filters.Filter{
		filters.WithFieldIn("role", []string{"anonymus", "member"}),
	})
	suite.Require().NoError(err)
	log.Println(users3)

	// Find users with time.Now() greater than created_at
	var users4 []models.User
	err = suite.userRepo.Find(context.TODO(), &users4, []filters.Filter{
		filters.WithFieldLessThan("created_at", time.Now()),
	})
	suite.Require().NoError(err)
	log.Println(users4)

	// Find users with not equal
	var users5 []models.User
	err = suite.userRepo.Find(context.TODO(), &users5, []filters.Filter{
		filters.WithFieldNotEqual("username", "Ga"),
	})
	suite.Require().NoError(err)
	log.Println(users5)

	// Find users with role not in ['anonymus', 'member']
	var users6 []models.User
	err = suite.userRepo.Find(context.TODO(), &users6, []filters.Filter{
		filters.WithFieldNotIn("role", []string{"anonymus", "member"}),
	})
	suite.Require().NoError(err)
	log.Println(users6)

	// Find users with field like
	var users7 []models.User
	err = suite.userRepo.Find(context.TODO(), &users7, []filters.Filter{
		filters.WithFieldLike("username", "H"),
	})
	suite.Require().NoError(err)
	log.Println(users7)

	// Find users with is null
	var users8 []models.User
	err = suite.userRepo.Find(context.TODO(), &users8, []filters.Filter{
		filters.WithFieldNotNull("username"),
	})
	suite.Require().NoError(err)
	log.Println(users8)
}

func (suite *UserRepositorySuite) TestFirstUsers() {
	// Find user with ID
	var user *models.User
	err := suite.userRepo.First(context.TODO(), &user, filters.WithID("655f3a9a3385f8f5a987fc9b"))
	suite.Require().NoError(err)
	log.Println(user)
}

func (suite *UserRepositorySuite) TestLastUsers() {
	// Find user with username
	var user *models.User
	err := suite.userRepo.Last(context.TODO(), &user, filters.WithFieldEqual("username", "Ga"))
	suite.Require().NoError(err)
	log.Println(user)
}

// TODO:
func (suite *UserRepositorySuite) TestFirstOrCreateUser() {
	user := &models.User{
		Username: "Ha",
	}

	err := suite.userRepo.First(context.TODO(), &user)
	suite.Require().NoError(err)
}

func (suite *UserRepositorySuite) TestUpdateUser() {
	var user *models.User
	err := suite.userRepo.First(context.TODO(), &user, filters.WithFieldEqual("level", 5))
	suite.Require().NoError(err)

	user.Username = "Ha ga"

	err = suite.userRepo.UpdateByID(context.TODO(), user.ID, user)
	suite.Require().NoError(err)
}

func (suite *UserRepositorySuite) TestDeleteUser() {
	var user *models.User
	err := suite.userRepo.First(context.TODO(), &user, filters.WithFieldEqual("level", 5))
	suite.Require().NoError(err)

	err = suite.userRepo.DeleteByID(context.TODO(), user.ID)
	suite.Require().NoError(err)
}
