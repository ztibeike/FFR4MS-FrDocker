package frecovery

import (
	"context"

	"gitee.com/zengtao321/frdocker/config"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (app *FrecoveryApp) Persist() *cron.Cron {
	app.Logger.Info("register persistence scheduled task...")
	c := cron.New()
	c.AddFunc(config.FRECOVERY_PERSISTENCE_INTERVAL, app.persistenceTask)
	c.Start()
	return c
}

func (app *FrecoveryApp) persistenceTask() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	collection := app.DbCli.Collection(config.FRECOVERY_PERSISTENCE_COLLECTION)
	filter := bson.D{{Key: "networkInterface", Value: app.NetworkInterface}}
	opts := options.Replace().SetUpsert(true)
	_, err := collection.ReplaceOne(ctx, filter, app, opts)
	if err != nil {
		app.Logger.Error("failed to persist frecovery app")
	}
}

func (app *FrecoveryApp) getPersistedApp() *FrecoveryApp {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	collection := app.DbCli.Collection(config.FRECOVERY_PERSISTENCE_COLLECTION)
	filter := bson.D{{Key: "networkInterface", Value: app.NetworkInterface}}
	var persistedApp FrecoveryApp
	err := collection.FindOne(ctx, filter).Decode(&persistedApp)
	if err != nil {
		return nil
	}
	return &persistedApp
}
