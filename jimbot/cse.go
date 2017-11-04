package jimbot

import (
	"io/ioutil"
	"log"

	"golang.org/x/oauth2"
	google "golang.org/x/oauth2/google"
	customsearch "google.golang.org/api/customsearch/v1"
)

// Result : CSE search result type
type Result struct {
	Position int64
	Result   *customsearch.Result
}

const (
	noResult = HUH + " Google search returned no result"
)

// Search : CSE search, for external use
func Search(query string, image bool) string {
	data, err := ioutil.ReadFile("cse-search-key.json")
	if err != nil {
		log.Print("[---] google json not found")
		return noResult
	}

	//Get the config from the json key file with the correct scope
	conf, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/cse")
	if err != nil {
		log.Print("[---] Read key err ", err)
		return noResult
	}

	// Initiate an http.Client. The following GET request will be
	// authorized and authenticated on the behalf of
	// your service account.
	client := conf.Client(oauth2.NoContext)

	cseService, _ := customsearch.New(client)
	search := cseService.Cse.List(query)
	if image == true {
		search.SearchType("image")
	}

	// CSE id of your search engine
	cseID := ReadConfig().CSE
	search.Cx(cseID)

	result := doSearch(search)

	if result.Position == 0 {
		log.Println("No results found in the top 10 pages.")
	} else if result.Result.Link == "" {
		log.Println("[---] NO RESULTS")
		return noResult
	}

	log.Println("*********************Google Search****************************")
	log.Printf("Position: %d\n", result.Position)
	log.Printf("Url: %s\n", result.Result.Link)
	log.Printf("Title: %s\n", result.Result.Title)
	log.Printf("Snippet: %s\n", result.Result.Snippet)

	return ("Title: " + result.Result.Title + "\n\nSnippet: " + result.Result.Snippet + "\n\nURL: " + result.Result.Link)
}

func doSearch(search *customsearch.CseListCall) (result Result) {

	start := int64(1)

	// CSE Limits you to 10 pages of results with max 10 results per page
	for start < 100 {
		search.Start(start)
		call, err := search.Do()
		if err != nil {
			log.Print("[---] CSE search failed ", err)
		}

		position, csResult := getResults(call.Items, start)

		if csResult != nil {
			result = Result{
				Position: position,
				Result:   csResult,
			}
			return // need only the first result
		}

		// No more search results?
		if call.SearchInformation.TotalResults < start {
			return
		}
		start = start + 10
	}

	return
}

func getResults(results []*customsearch.Result, start int64) (position int64, result *customsearch.Result) {
	for index, r := range results {
		return int64(index) + start, r
	}
	return 0, nil
}
