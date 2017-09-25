package jimbot

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/buger/jsonparser"
)

const (
	apiURL = "https://translate.googleapis.com/translate_a/single?client=gtx&sl=auto&tl=en&dt=t&q="
)

// ToEnglish : translate any language to English
func ToEnglish(text string) string {
	var textTrans string

	query := url.PathEscape(text)
	queryURL := apiURL + query
	log.Print("[*] TRANSLATE API URL: ", queryURL)
	req, err := http.Get(string(queryURL))
	if err != nil {
		log.Print("[*] Can't reach translate API")
	}
	defer req.Body.Close()
	readBody, _ := ioutil.ReadAll(req.Body)

	data, _, _, _ := jsonparser.Get(readBody, "[0]", "[0]", "[0]")

	// TODO : Parse raw []byte as UTF-8 string

	log.Print("[+++] TRANSLATE response: ", string(readBody))
	log.Print("[+++] TRANSLATE response decoded: ", string(data))

	textTrans = string(data)
	return textTrans
}
