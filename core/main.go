package main

import (
	"encoding/json"
	"fmt"
	. "github.com/leochan007/aci-blockchain-updater/utils"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"
)

var eosWrapper *EosWrapper = nil

var mongoWrapper *MongoWrapper = nil

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getParams() (mongoUrl string, baseUrl string, sleepTime int64) {
	mongoUrl = getEnv("MONGODB_URL", "mongodb://127.0.0.1:27017/alphacar")
	baseUrl = getEnv("BASE_URL", "https://eos.greymass.com")
	val, err := strconv.ParseInt(getEnv("SLEEP_TIME", "5000"), 10, 64)
	if err == nil {
		sleepTime = val
	} else {
		sleepTime = 5000
	}
	return
}

func fetchAndUpdate(mongoUrl string) {

	err := mongoWrapper.InitClient(mongoUrl)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	result, err := mongoWrapper.GetProcessedCreditInquiries("processed")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for k, v := range result {
		fmt.Println(k, " ", v)

		resp, err := eosWrapper.GetTransaction(v.LocalTxId)
		if err == nil {
			b, err := json.Marshal(resp)
			if err != nil {
				fmt.Println("json.Marshal failed:", err)
			}
			fmt.Println(string(b))
			status := resp["trx"].(map[string]interface{})["receipt"].(map[string]interface{})["status"]

			if status == "executed" {
				blockNum := resp["block_num"].(float64)
				lastIrreversibleBlock := resp["last_irreversible_block"].(float64)
				txId := resp["id"].(string)

				if blockNum <= lastIrreversibleBlock {
					fmt.Println("EXECUTED IRREVERSIBLE")
					err := mongoWrapper.UpdateTxId(v.Hash, txId)
					if err != nil {
						fmt.Println(err.Error())
					}
				}
			}
		} else {
			fmt.Println(err.Error())
		}

	}

	err = mongoWrapper.Close()

	if err != nil {
		fmt.Println(err.Error())
	}

}

var c chan os.Signal
var wg sync.WaitGroup

func run(sleepTime int64, mongoUrl string) {

LOOP:
	for {
		select {
		case s := <-c:
			fmt.Println()
			fmt.Println("run get:", s)
			break LOOP
		default:
		}

		fmt.Println("begin fetch... ", time.Now())
		fetchAndUpdate(mongoUrl)

		time.Sleep(time.Duration(sleepTime) * time.Millisecond)

	}
	wg.Done()

}

func main() {
	fmt.Println("start")
	mongoUrl, baseUrl, sleepTime := getParams()

	eosWrapper = &EosWrapper{BaseUrl: baseUrl}

	fmt.Println("mongoUrl:", mongoUrl)

	mongoWrapper = &MongoWrapper{}

	c = make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	wg.Add(1)
	go run(sleepTime, mongoUrl)

	wg.Wait()

}
