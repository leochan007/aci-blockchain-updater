package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type EosWrapper struct {
	BaseUrl string
}

func (instance *EosWrapper) post(url string, reader *bytes.Reader) (result map[string]interface{}, err error){
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		result = nil
		return
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := http.Client{}
	resp := &http.Response{}
	resp, err = client.Do(request)
	if err != nil {
		result = nil
		return
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		result = nil
		return
	}

	err = json.Unmarshal(respBytes, &result)

	if err != nil {
		result = nil
		return
	}

	return
}

func(instance *EosWrapper) GetTransaction(id string) (result map[string]interface{}, err error) {

	txInfo := make(map[string]interface{})
	txInfo["id"] = id
	bytesData, err := json.Marshal(txInfo)

	if err != nil {
		result = nil
		return
	}
	reader := bytes.NewReader(bytesData)

	result, err = instance.post(instance.BaseUrl + "/v1/history/get_transaction", reader)
	return
}
