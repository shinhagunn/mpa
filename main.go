package main

import (
	"context"
	"log"

	"github.com/shinhagunn/mpa/config"
	"github.com/shinhagunn/mpa/filters"
	"github.com/shinhagunn/mpa/models"
	"github.com/shinhagunn/mpa/pkg/mpa_fx"
	"github.com/shinhagunn/mpa/repo"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	db, err := mpa_fx.New(cfg)
	if err != nil {
		panic(err)
	}

	userRepo := repo.New(db, &models.User{})

	var users []models.User
	if err := userRepo.Find(context.TODO(), &users, []filters.Filter{}); err != nil {
		panic(err)
	}

	// cur, err := db.Collection("users").Find(context.TODO(), filters.ApplyFilters(flts...))
	// if err != nil {
	// 	panic(err)
	// }

	// var users []models.User
	// if err := cur.All(context.TODO(), &users); err != nil {
	// 	panic(err)
	// }

	log.Println(users)
}
