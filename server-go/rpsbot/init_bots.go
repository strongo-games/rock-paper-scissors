package rpsbot

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"github.com/strongo/app"
	"github.com/prizarena/rock-paper-scissors/server-go/rpsbot/platforms/rpstgbots"
	"github.com/prizarena/rock-paper-scissors/server-go/rpsbot/rpsrouting"
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/bots-framework/platforms/telegram"
	"net/http"
	"github.com/prizarena/rock-paper-scissors/server-go/rpssecrets"
)

func InitBot(httpRouter *httprouter.Router, botHost bots.BotHost, appContext bots.BotAppContext) error {
	gaSettings := bots.AnalyticsSettings{
		GaTrackingID: rpssecrets.GaTrackingID,
		Enabled: func(r *http.Request) bool {
			return false
		},
	}

	driver := bots.NewBotDriver(gaSettings, appContext, botHost,
		"Please report any issues to @",
	)
	//routing.WebhooksRouter

	if driver == nil {
		panic("driver == nil")
	}

	newTranslator := func(c context.Context) strongo.Translator {
		return strongo.NewMapTranslator(c, nil)
	}

	driver.RegisterWebhookHandlers(httpRouter, "/bot",
		telegram.NewTelegramWebhookHandler(
			func(c context.Context) bots.SettingsBy {
				return rpstgbots.Bots(c, strongo.EnvProduction, rpsrouting.WebhooksRouter) //gaestandard.GetEnvironment(c)
			}, // Maps of bots by code, language, token, etc...
			newTranslator, // Creates newTranslator that gets a context.Context (for logging purpose)
		),
	)

	return nil
}
