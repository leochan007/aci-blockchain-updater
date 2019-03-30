package utils

import (
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"log"
	"strings"
	"time"
)

type MongoWrapper struct {
	client *mongo.Client
	ctx    context.Context
	collection *mongo.Collection
}

func (instance *MongoWrapper) InitClient(mongoUrl string, dbName string, collectionName string, mongoOpt string) (err error) {
	instance.ctx, _ = context.WithTimeout(context.Background(), 10 * time.Second)
	tmpMongoUrl := mongoUrl
	if !strings.HasSuffix(mongoUrl, "/") {
		tmpMongoUrl += "/"
	}
	url := tmpMongoUrl + dbName + mongoOpt
	instance.client, err = mongo.Connect(instance.ctx, url)
	if err == nil {
		instance.collection = instance.client.Database(dbName).Collection(collectionName)
	}
	return
}

func (instance *MongoWrapper) GetProcessedCreditInquiries(status string) (results []QueryResult, errInfo error) {
	if instance.client != nil {
		filter := bson.M{"status": status}
		cur, err := instance.collection.Find(instance.ctx, filter)
		if err != nil {
			log.Fatal(err)
			results = nil
			errInfo = err
			return
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

func (instance *MongoWrapper) UpdateTxId(hash string, txId string) (errInfo error) {
	if instance.client != nil {
		filter := bson.M{"hash": hash}
		update := bson.D{
			{"$set", bson.D{
				{"txId", txId},
				{"status", "confirmed"},
			}},
		}
		_, err := instance.collection.UpdateOne(instance.ctx, filter, update)
		if err != nil {
			log.Fatal(err)
		}
		return
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
