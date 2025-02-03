package jimbot

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/buger/jsonparser"
)

// NowWeather : Current weather info from his/her city
func NowWeather(city string) string {
	apiURL := fmt.Sprintf("https://api.seniverse.com/v3/weather/now.json?key=jkfkmayhcqnvsrjr&location=%s&language=en&unit=c", city)

	resp, err := http.Get(apiURL)
	if err != nil {
		log.Print("[===] WEATHER API cant be reached")
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	readBody, _ := io.ReadAll(resp.Body)

	status, _, _, _ := jsonparser.Get(readBody, "results", "[0]", "now", "text")
	temp, _, _, _ := jsonparser.Get(readBody, "results", "[0]", "now", "temperature")
	retVal := "`Now in " + strings.ToUpper(city) + ": " + string(status) + "\nTemp: " + string(temp) + " °C`"
	return retVal
}
