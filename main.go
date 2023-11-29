package main

import (
	"context"
	"log"

	"github.com/shinhagunn/mpa/config"
	"github.com/shinhagunn/mpa/models"
	"github.com/shinhagunn/mpa/pkg/mongo_fx"
	"github.com/zsmartex/pkg/v2/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type Command string

var (
	CommandCreate = Command("create")
	CommandUpdate = Command("update")
	CommandDelete = Command("delete")
)

type Action string

var (
	ActionBefore = Action("before")
	ActionAfter  = Action("after")
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	mpa, err := mongo_fx.New(cfg)
	if err != nil {
		panic(err)
	}

	userRepo := mongo_fx.NewRepository(mpa, models.User{})

	userRepo.AddCallback(mongo_fx.CallbackTypeBeforeCreate, func(db *mongo.Database, value *models.User) error {
		log.Println("Before create")
		value.Role = models.UserRoleSuperAdmin
		log.Println(value)

		return nil
	})

	userRepo.AddCallback(mongo_fx.CallbackTypeAfterCreate, func(db *mongo.Database, value *models.User) error {
		log.Println("After create")
		log.Println(value)

		return nil
	})

	err = userRepo.Create(
		context.Background(),
		&models.User{
			Email: "test1",
			Role:  models.UserRoleAdmin,
			UID:   utils.GenerateUID(),
		},
	)
	if err != nil {
		panic(err)
	}
}
