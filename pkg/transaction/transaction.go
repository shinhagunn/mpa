package transaction

import (
	"context"
	"sync"

	"github.com/shinhagunn/mpa/models"
	"github.com/zsmartex/pkg/v2/usecase"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type ListenMessage struct {
	// Type message
	Type string

	// Database name
	Name string

	// Value
	Value interface{}
}

type Transaction struct {
	// Channel error for all database
	Errors map[string]chan error

	// Usecases per database
	Usecases TransactionUsecase

	// Data per database
	Datas TransactionData

	// Listen event from databases
	Listen chan ListenMessage

	// Send status to databases
	SendMongo chan error

	// Send status to databases
	SendGorm chan error

	// Cancel channel
	Cancel chan interface{}
}

type TransactionUsecase struct {
	gorm  usecase.IUsecase[models.User]
	mongo *mongo.Database
}

type TransactionData struct {
	gorm  *gorm.DB
	mongo mongo.SessionContext
}

type TransactionConfig struct {
	GormUsecase usecase.IUsecase[models.User]
	MongoDB     *mongo.Database
}

func New(config TransactionUsecase) *Transaction {
	errors := make(map[string]chan error)

	errors["gorm"] = make(chan error)
	errors["mongo"] = make(chan error)

	return &Transaction{
		Errors:    errors,
		Usecases:  config,
		Datas:     TransactionData{},
		Listen:    make(chan ListenMessage),
		SendMongo: make(chan error),
		SendGorm:  make(chan error),
	}
}

func (t *Transaction) Run(ctx context.Context, fn func(data TransactionData) error) error {
	var wg sync.WaitGroup

	wg.Add(2)
	go t.ListenEvent(&wg)
	go TransactionGorm(ctx, t)
	go TransactionMongo(ctx, t)

	wg.Wait()

	if err := fn(t.Datas); err != nil {
		t.SendMongo <- err
		t.SendGorm <- err
	} else {
		t.SendMongo <- nil
		t.SendGorm <- nil
	}

	t.Listen <- ListenMessage{
		Type: "done",
		Name: "transaction",
	}

	return nil
}

func (t *Transaction) ListenEvent(wg *sync.WaitGroup) {
	for mess := range t.Listen {
		switch mess.Type {
		case "receive":
			if mess.Name == "gorm" {
				t.Datas.gorm = mess.Value.(*gorm.DB)
			} else if mess.Name == "mongo" {
				t.Datas.mongo = mess.Value.(mongo.SessionContext)
			}
			wg.Done()
		case "error":
			t.Cancel <- mess.Value
		case "done":
			break
		}
	}
}
