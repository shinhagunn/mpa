package transaction

import (
	"context"

	"gorm.io/gorm"
)

func TransactionGorm(ctx context.Context, trans *Transaction) {
	usecase := trans.Usecases.Gorm

	if err := usecase.Transaction(func(tx *gorm.DB) error {
		trans.Listen <- ListenMessage{
			Type:  "receive",
			Name:  "gorm",
			Value: tx,
		}

		for sendErr := range trans.SendGorm {
			if sendErr != nil {
				return sendErr
			}

			break
		}

		return nil
	}); err != nil {
		trans.Listen <- ListenMessage{
			Type:  "error",
			Name:  "gorm",
			Value: err,
		}
		return
	}
}
