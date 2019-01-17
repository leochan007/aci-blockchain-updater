package main

import (
	"fmt"
	. "github.com/leochan007/aci-blockchain-updater1/utils"
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

	mongoWrapper := &MongoWrapper{}

	err := mongoWrapper.InitClient(mongoUrl)

	if err != nil {
		fmt.Println(err.Error())
	}

	eosWrapper := &EosWrapper{BaseUrl: baseUrl}

	resp, err := eosWrapper.GetTransaction("83eb311b599276ff3988ba043cae390a549c39d9cfee60a12cda28e292040548")
	if err == nil {
		fmt.Println(resp)
	} else {
		fmt.Println(err.Error())
	}

	resp, err = eosWrapper.GetTransaction("83eb311b599276ff3988ba043cae390a549c39d9cfee60a12cda28e292040541")
	if err == nil {
		fmt.Println(resp)
	} else {
		fmt.Println(err.Error())
	}

	err = mongoWrapper.Close()

}
