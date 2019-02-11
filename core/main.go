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

func getParams() (mongoUrl string, mongoOpt string, baseUrl string, sleepTime int64) {
	mongoUrl = getEnv("MONGODB_URL", "mongodb://127.0.0.1:27017")
	mongoOpt = getEnv("MONGODB_OPT", "?authSource=admin")
	baseUrl = getEnv("EOSIO_HTTP_URL", "https://eos.greymass.com")
	val, err := strconv.ParseInt(getEnv("SLEEP_TIME", "5000"), 10, 64)
	if err == nil {
		sleepTime = val
	} else {
		sleepTime = 5000
	}
	return
}

func fetchAndUpdate(mongoUrl string, mongoOpt string) {
	
	err := mongoWrapper.InitClient(mongoUrl, "alphacar", mongoOpt)

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
			if _, ok := resp["trx"]; ok {
				txMap1 := resp["trx"].(map[string]interface{})
				if _, ok := txMap1["receipt"]; ok {
					txMap2 := txMap1["receipt"].(map[string]interface{})
					if _, ok := txMap2["status"]; ok {
						status := txMap2["status"].(string)

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

					}
				}
			}
			//status := resp["trx"].(map[string]interface{})["receipt"].(map[string]interface{})["status"]
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

func run(sleepTime int64, mongoUrl string, mongoOpt string) {

	tickTimer := time.NewTicker(time.Duration(sleepTime) * time.Millisecond)

LOOP:
	for {
		select {
		case s := <-c:
			fmt.Println()
			fmt.Println("run get:", s)
			break LOOP
		case <-tickTimer.C:
			fmt.Println("with TICKER begin fetch... ", time.Now())
			fetchAndUpdate(mongoUrl, mongoOpt)
		default:
		}

	}
	wg.Done()

}

func main() {
	fmt.Println("start")
	mongoUrl, mongoOpt, baseUrl, sleepTime := getParams()

	eosWrapper = &EosWrapper{BaseUrl: baseUrl}

	fmt.Println("baseUrl:", baseUrl)

	mongoWrapper = &MongoWrapper{}

	c = make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	wg.Add(1)

	go run(sleepTime, mongoUrl, mongoOpt)

	wg.Wait()

}
