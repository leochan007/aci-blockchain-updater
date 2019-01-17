package utils

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"time"
)

type MongoWrapper struct {
	client *mongo.Client
	ctx context.Context
}

func (instance *MongoWrapper) InitClient(mongoUrl string) (err error) {
	instance.ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	instance.client, err = mongo.Connect(instance.ctx, mongoUrl)

	fmt.Println("client", instance.client)
	return
}

func(instance *MongoWrapper) Close() (err error) {
	if instance.client != nil {
		err = instance.client.Disconnect(instance.ctx)
	}
	return
}
