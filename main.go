package main

import (
	"context"
	"log"

	"github.com/shinhagunn/mpa/config"
	"github.com/shinhagunn/mpa/models"
	filters "github.com/shinhagunn/mpa/mongodb/fitlers"
	"github.com/shinhagunn/mpa/pkg/mongo_fx"
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

	users, err := userRepo.Find(
		context.Background(),
		nil,
		filters.WithOr(
			filters.WithFieldEqual("role", "admin"),
			filters.WithFieldEqual("role", "ahihi"),
		),
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

	log.Println(users)
}
