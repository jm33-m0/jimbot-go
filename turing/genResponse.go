package turing

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/buger/jsonparser"
)

const apiURL = "http://openapi.tuling123.com/openapi/api/v2"

// GetResponse : get response from turing api
func GetResponse(input string) string {
	// json data to post
	data := []byte(fmt.Sprintf(`{
		"reqType":0,
		"perception": {
			"inputText": {
				"text": "%s"
			},
			"inputImage": {
				"url": ""
			},
			"selfInfo": {
				"location": {
					"city": "",
					"province": "",
					"street": ""
				}
			}
		},
		"userInfo": {
			"apiKey": "063487b9cffb41adbebe25c86b56b807",
			"userId": "jimdeb3"
		}}`, input))

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "Could not reach Turing API"
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	responseRaw, _ := ioutil.ReadAll(resp.Body)

	response, _ := jsonparser.GetString(responseRaw, "results", "[0]", "values", "text")
	return response
}
