package mpa_fx

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
)

var (
	Module = fx.Module("mpa_fx", mpaProviders, mpaInvokes)

	mpaProviders = fx.Provide(New)

	mpaInvokes = fx.Invoke(registerHooks)
)

func registerHooks(lc fx.Lifecycle, db *mongo.Database) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			if err := db.Client().Disconnect(ctx); err != nil {
				return err
			}

			return nil
		},
	})
}
