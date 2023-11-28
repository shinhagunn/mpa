package main

import (
	"context"
	"log"
	"time"

	"github.com/shinhagunn/mpa/config"
	"github.com/shinhagunn/mpa/models"
	"github.com/shinhagunn/mpa/mongodb"
	"github.com/shinhagunn/mpa/pkg/mongo_fx"
	"github.com/zsmartex/pkg/v2/utils"
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

	db, err := mongo_fx.New(cfg)
	if err != nil {
		panic(err)
	}

	userRepo := mongo_fx.NewRepository(db, models.User{})

	userRepo.AddCallback(mongodb.BeforeCreate, func() {
		log.Println("Hello onichan!")
	})

	user := &models.User{
		UID:       utils.GenerateUID(),
		Email:     "ha@gmail.com",
		Role:      models.UserRoleMember,
		CreatedAt: time.Now().Local(),
		UpdatedAt: time.Now().Local(),
	}
	err = userRepo.Create(
		context.Background(),
		user,
	)
	if err != nil {
		panic(err)
	}

	// cur, err := db.Collection("users").Find(context.Background(), bson.M{
	// 	"$or": []bson.M{
	// 		{"role": "admin"},
	// 		{"role": "ahihi"},
	// 	},
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// var users []models.User
	// err = cur.All(context.Background(), &users)
	// if err != nil {
	// 	panic(err)
	// }

	log.Println(user)
}
