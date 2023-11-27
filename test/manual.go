package test

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// This is a example about transaction manual
// But it's not working, when transaction 1 commit success, but transaction 2 commit failed, the transaction 1 will be rollback

type UsecaseGORM struct {
	db *gorm.DB
}

func (u UsecaseGORM) TransactionGORM(handler func(tx *gorm.DB) error) (*gorm.DB, error) {
	// Create transaction manual
	trans := u.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			trans.Rollback()
		}
	}()

	// Check error when create
	if err := trans.Error; err != nil {
		return nil, err
	}

	// Handle all command
	if err := handler(trans); err != nil {
		trans.Rollback()
		return nil, err
	}

	// Return trans for commit
	return trans, nil
}

type UsecaseMongoDB struct {
	db *mongo.Database
}

func (u UsecaseMongoDB) TransactionMongo(ctx context.Context, handler func(tx mongo.SessionContext) error) (mongo.Session, error) {
	session, err := u.db.Client().StartSession()
	if err != nil {
		return nil, err
	}

	sessionContext := mongo.NewSessionContext(ctx, session)
	if err := session.StartTransaction(); err != nil {
		return nil, err
	}

	if err := handler(sessionContext); err != nil {
		if abortErr := session.AbortTransaction(ctx); abortErr != nil {
			return nil, abortErr
		}
		return nil, err
	}

	return session, nil
}

func RandomAPI(c *fiber.Ctx) error {
	ctx := c.Context()

	// Above handle validate params
	usecaseGORM := UsecaseGORM{}
	usecaseMongoDB := UsecaseMongoDB{}

	txGORM, err := usecaseGORM.TransactionGORM(func(tx *gorm.DB) error {
		// Handle usecase here

		return nil
	})
	if err != nil {
		return err
	}

	txMongoDB, err := usecaseMongoDB.TransactionMongo(context.Background(), func(tx mongo.SessionContext) error {
		// Handle usecase here

		return nil
	})
	if err != nil {
		return err
	}

	if err := txGORM.Commit().Error; err != nil {
		return err
	}

	if err := txMongoDB.CommitTransaction(ctx); err != nil {
		return err
	}

	// Below handle data after transaction success

	return nil
}
