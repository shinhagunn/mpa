package transaction

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

func TransactionMongo(ctx context.Context, tx *Transaction) {
	db := tx.Usecases.Mongo

	session, err := db.Client().StartSession()
	if err != nil {
		tx.Listen <- ListenMessage{
			Type:  "error",
			Name:  "mongo",
			Value: err,
		}
		return
	}
	defer session.EndSession(ctx)

	if err := mongo.WithSession(context.Background(), session, func(sessionContext mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}

		tx.Listen <- ListenMessage{
			Type:  "receive",
			Name:  "mongo",
			Value: sessionContext,
		}

		for sendErr := range tx.SendMongo {
			if sendErr != nil {
				if abortErr := session.AbortTransaction(context.Background()); abortErr != nil {
					return abortErr
				}

				return sendErr
			}

			break
		}

		if err := session.CommitTransaction(sessionContext); err != nil {
			return err
		}

		return nil
	}); err != nil {
		if abortErr := session.AbortTransaction(context.Background()); abortErr != nil {
			tx.Listen <- ListenMessage{
				Type:  "error",
				Name:  "mongo",
				Value: err,
			}
			return
		}

		tx.Listen <- ListenMessage{
			Type:  "error",
			Name:  "mongo",
			Value: err,
		}
	}
}
