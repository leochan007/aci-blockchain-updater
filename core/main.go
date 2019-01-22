package main

import (
	"encoding/json"
	"fmt"
	. "github.com/leochan007/aci-blockchain-updater/utils"
	"os"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getParams() (mongoUrl string, baseUrl string) {
	mongoUrl = getEnv("MONGODB_URL", "mongodb://127.0.0.1:27017/alphacar")
	baseUrl = getEnv("BASE_URL", "https://eos.greymass.com")
	return
}

func main() {
	fmt.Println("start")
	mongoUrl, baseUrl := getParams()

	eosWrapper := &EosWrapper{BaseUrl: baseUrl}

	fmt.Println("mongoUrl:", mongoUrl)

	mongoWrapper := &MongoWrapper{}

	err := mongoWrapper.InitClient(mongoUrl)

	if err != nil {
		fmt.Println(err.Error())
	}

	result, err := mongoWrapper.GetProcessedCreditInquiries("processed")

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

				if blockNum <= lastIrreversibleBlock {
					fmt.Println("EXECUTED IRREVERSIBLE")
				}
			}
		} else {
			fmt.Println(err.Error())
		}

	}

	if err != nil {
		fmt.Println(err.Error())
	}

	err = mongoWrapper.Close()

}
