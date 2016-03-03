package main

import "encoding/json"
import "log"
import "net/http"
import "io/ioutil"
import "os"
import "strconv"
import "strings"

type unmarshaller interface {
	Unmarshal([]byte, *map[string]interface{}) error
}
type mainUnmarshaller struct{}

func (unmarshaller mainUnmarshaller) Unmarshal(inp []byte, resp *map[string]interface{}) error {
	return json.Unmarshal(inp, resp)
}

type responseConverter interface {
	Convert(*http.Response) ([]byte, error)
}
type mainConverter struct{}

func (converter mainConverter) Convert(response *http.Response) ([]byte, error) {
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}

type httpResponseFetcher interface {
	Fetch(url string) (*http.Response, error)
}
type mainFetcher struct{}

func (fetcher mainFetcher) Fetch(url string) (*http.Response, error) {
	log.Printf("Retrieving %q\n", url)
	return http.Get(url)
}

func getBeerPage(fetcher httpResponseFetcher, converter responseConverter, id int) string {
	url := "https://api.untappd.com/v4/beer/info/BID?client_id=CLIENTID&client_secret=CLIENTSECRET&compact=true"
	url = strings.Replace(url, "BID", strconv.Itoa(id), 1)
	url = strings.Replace(url, "CLIENTID", os.Getenv("CLIENTID"), 1)
	url = strings.Replace(url, "CLIENTSECRET", os.Getenv("CLIENTSECRET"), 1)

	response, err := fetcher.Fetch(url)

	if err != nil {
		log.Printf("Failed on getBeerPage: %q\n", err)
	} else {
		contents, err := converter.Convert(response)
		if err != nil {
			log.Printf("%q\n", err)
		} else {
			return string(contents)
		}
	}

	return "Failed to retrieve " + strconv.Itoa(id)
}

func getVenuePage(fetcher httpResponseFetcher, converter responseConverter, id int) string {
	url := "https://api.untappd.com/v4/venue/info/VID?client_id=CLIENTID&client_secret=CLIENTSECRET&compact=true"
	url = strings.Replace(url, "VID", strconv.Itoa(id), 1)
	url = strings.Replace(url, "CLIENTID", os.Getenv("CLIENTID"), 1)
	url = strings.Replace(url, "CLIENTSECRET", os.Getenv("CLIENTSECRET"), 1)

	response, err := fetcher.Fetch(url)

	log.Printf("Getting venue page: %v\n", url)

	if err != nil {
		log.Printf("Failed on getVenuePage: %q\n", err)
	} else {
		contents, err := converter.Convert(response)
		if err != nil {
			log.Printf("%q\n", err)
		} else {
			return string(contents)
		}
	}

	return "Failed to retrieve " + strconv.Itoa(id)
}

func convertPageToName(page string, unmarshaller unmarshaller) string {
	var mapper map[string]interface{}
	err := unmarshaller.Unmarshal([]byte(page), &mapper)
	if err != nil {
		log.Printf("%q\n", err)
		return "Failed to unmarshal"
	}

	response := mapper["response"].(map[string]interface{})
	beer := response["beer"].(map[string]interface{})
	brewery := beer["brewery"].(map[string]interface{})
	return brewery["brewery_name"].(string) + " - " + beer["beer_name"].(string)
}

func convertPageToDrinks(page string, unmarshaller unmarshaller) ([]int, error) {
	var mapper map[string]interface{}
	var values []int
	err := unmarshaller.Unmarshal([]byte(page), &mapper)
	if err != nil {
		log.Printf("%q\n", err)
		return values, err
	}

	response := mapper["response"].(map[string]interface{})
	checkins := response["checkins"].(map[string]interface{})
	items := checkins["items"].([]interface{})

	for _, v := range items {
		beer := v.(map[string]interface{})["beer"].(map[string]interface{})
		beerID := int(beer["bid"].(float64))
		values = append(values, beerID)
	}

	return values, nil
}

// GetBeerName Determines the name of the beer from the id
func GetBeerName(id int) string {
	var fetcher httpResponseFetcher = mainFetcher{}
	var converter responseConverter = mainConverter{}
	var unmarshaller unmarshaller = mainUnmarshaller{}
	text := getBeerPage(fetcher, converter, id)
	return convertPageToName(text, unmarshaller)
}
