package jimbot

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"unicode/utf8"

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

	client := &http.Client{}
	req, _ := http.NewRequest("GET", queryURL, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:64.0) Gecko/20100101 Firefox/64.0")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Add("Accept-Language", "zh-CN,en-US;q=0.7,en;q=0.3")

	res, err := client.Do(req)
	if err != nil {
		log.Print("[*] Can't reach translate API")
	}
	defer func() {
		err = res.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	readBody, _ := ioutil.ReadAll(res.Body)

	data, _, _, _ := jsonparser.Get(readBody, "[0]", "[0]", "[0]")

	log.Print("[+++] TRANSLATE response: ", readBody)

	for len(data) > 0 {
		r, size := utf8.DecodeRune(data)

		textTrans += string(r)
		data = data[size:]
	}

	log.Print("[+++] TRANSLATE response decoded: ", textTrans)

	return textTrans
}
