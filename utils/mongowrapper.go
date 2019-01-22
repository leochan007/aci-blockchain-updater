package utils

import (
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"log"
	"time"
)

type MongoWrapper struct {
	client *mongo.Client
	ctx    context.Context
}

func (instance *MongoWrapper) InitClient(mongoUrl string) (err error) {
	instance.ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	instance.client, err = mongo.Connect(instance.ctx, mongoUrl)

	//fmt.Println("client", instance.client)
	return
}

func (instance *MongoWrapper) GetProcessedCreditInquiries(status string) (results []QueryResult, errInfo error) {
	if instance.client != nil {
		filter := bson.M{"status": status}
		collection := instance.client.Database("alphacar").Collection("credit_inquiry")
		cur, err := collection.Find(instance.ctx, filter)
		if err != nil {
			log.Fatal(err)
		}
		for cur.Next(instance.ctx) {
			var result bson.M
			err := cur.Decode(&result)
			if err != nil {
				log.Fatal(err)
				results = nil
				errInfo = err
				return
			}
			var item QueryResult
			item.Hash = result["hash"].(string)
			item.LocalTxId = result["localTxId"].(string)
			results = append(results, item)
		}
		if err := cur.Err(); err != nil {
			log.Fatal(err)
			results = nil
			errInfo = err
			return
		}
	}
	errInfo = nil
	return
}

func (instance *MongoWrapper) Close() (err error) {
	if instance.client != nil {
		err = instance.client.Disconnect(instance.ctx)
	}
	return
}
