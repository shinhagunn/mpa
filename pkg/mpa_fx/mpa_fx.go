package mpa_fx

import (
	"go.uber.org/fx"
)

var (
	Module = fx.Module("mpa_fx", mpaProviders)

	mpaProviders = fx.Provide(New)

	// mpaInvokes = fx.Invoke(registerHooks)
)

// func registerHooks(lc fx.Lifecycle, client *mongo.Client) {
// 	lc.Append(fx.StartHook(func() {
// 		if err := client.Disconnect(context.TODO()); err != nil {
// 			panic(err)
// 		}
// 	}))
// }
