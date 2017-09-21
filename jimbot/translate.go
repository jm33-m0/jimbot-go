package jimbot

import (
	"github.com/buger/jsonparser"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	apiURL = "https://translate.googleapis.com/translate_a/single?client=gtx&sl=auto&tl=en&dt=t&q="
)

// ToEnglish : translate any language to English
func ToEnglish(text string) string {
	var textTrans string

	queryURL := url.PathEscape(apiURL + text)
	req, err := http.Get(string(queryURL))
	if err != nil {
		log.Print("[*] Can't reach translate API")
	}
	defer req.Body.Close()
	readBody, _ := ioutil.ReadAll(req.Body)

	data, _, _, _ := jsonparser.Get(readBody, "[0]", "[0]", "[0]")

	log.Print("[+++] TRANSLATE response: ", string(readBody))
	log.Print("[+++] TRANSLATE response decoded: ", string(data))

	textTrans = string(data)
	return textTrans
}
