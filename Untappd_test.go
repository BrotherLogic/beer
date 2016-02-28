package main

import "errors"
import "net/http"
import "strings"
import "testing"

type stubFailUnmarshaller struct{}

func (unmarshaller stubFailUnmarshaller) Unmarshal(int []byte, resp *map[string]interface{}) error {
	return errors.New("Built to fail")
}

type stubPassFetcher struct{}

func (fetcher stubPassFetcher) Fetch(url string) (*http.Response, error) {
	var resp = &http.Response{}
	return resp, nil
}

type stubFailConverter struct{}

func (converter stubFailConverter) Convert(response *http.Response) ([]byte, error) {
	return make([]byte, 0), errors.New("Built to fail")
}

type stubFailFetcher struct{}

func (fetcher stubFailFetcher) Fetch(url string) (*http.Response, error) {
	var resp = &http.Response{}
	var err = errors.New("Built to fail")
	return resp, err
}

func TestGetBeerName(t *testing.T) {
	beerName := GetBeerName(7936)
	if beerName != "Firestone Walker Brewing Company - Parabola" {
		t.Errorf("Beer name %q is not firestone, parabola\n", beerName)
	}
}

func TestGetBeerPage(t *testing.T) {
	var fetcher = mainFetcher{}
	var converter = mainConverter{}
	beerPage := getBeerPage(fetcher, converter, 7936)
	if !strings.Contains(beerPage, "Firestone") {
		t.Errorf("Beer page is not being retrieved\n%q\n", beerPage)
	}
}

func TestGetBeerPageFailHttp(t *testing.T) {
	var fetcher httpResponseFetcher = stubFailFetcher{}
	var converter = mainConverter{}
	beerPage := getBeerPage(fetcher, converter, 7936)
	if !strings.Contains(beerPage, "Failed to retrieve") {
		t.Errorf("Beer page retrieve did not fail\n%q\n", beerPage)
	}
}

func TestGetBeerPageConvertHttp(t *testing.T) {
	var fetcher httpResponseFetcher = stubPassFetcher{}
	var converter = stubFailConverter{}
	beerPage := getBeerPage(fetcher, converter, 7936)
	if !strings.Contains(beerPage, "Failed to retrieve") {
		t.Errorf("Beer page retrieve did not fail\n%q\n", beerPage)
	}
}

func TestConvertPageToName(t *testing.T) {
	var unmarshaller = mainUnmarshaller{}
	beerName := convertPageToName(mockBeerPage, unmarshaller)
	if beerName != "Firestone Walker Brewing Company - Parabola" {
		t.Errorf("Beer name %q is not parabola\n", beerName)
	}
}

func TestFailingConvertPageToName(t *testing.T) {
	var unmarshaller unmarshaller = stubFailUnmarshaller{}
	beerName := convertPageToName(mockBeerPage, unmarshaller)
	if !strings.Contains(beerName, "Failed to unmarshal") {
		t.Errorf("Unmarshalling did not fail\n")
	}
}

const mockBeerPage string = "{\"meta\":{\"code\":200,\"response_time\":{\"time\":0.012,\"measure\":\"seconds\"},\"init_time\":{\"time\":0.005,\"measure\":\"seconds\"}},\"notifications\":[],\"response\":{\"beer\":{\"bid\":7936,\"beer_name\":\"Parabola\",\"beer_label\":\"https:\\/\\/untappd.akamaized.net\\/site\\/beer_logos\\/beer-firestoneWalkerParabola.jpg\",\"beer_label_hd\":\"\",\"beer_abv\":14,\"beer_ibu\":82,\"beer_description\":\"One of our most aggressive and sought-after offerings. Bold bourbon, tobacco and espresso aromas and a hint of American oak greet the nose. Rich, chewy roasted malts, charred oak and bourbon-like vanilla fill the palate and create a seamless finish. A remarkably complex brew that offers a transcendental drinking experience – enjoy with good company.\\n\",\"beer_style\":\"Stout - Russian Imperial\",\"is_in_production\":1,\"beer_slug\":\"firestone-walker-brewing-company-parabola\",\"is_homebrew\":0,\"created_at\":\"Sun, 07 Nov 2010 16:38:18 +0000\",\"rating_count\":10169,\"rating_score\":4.55809,\"stats\":{\"total_count\":14261,\"monthly_count\":574,\"total_user_count\":11946,\"user_count\":0},\"brewery\":{\"brewery_id\":524,\"brewery_name\":\"Firestone Walker Brewing Company\",\"brewery_label\":\"https:\\/\\/untappd.akamaized.net\\/site\\/brewery_logos\\/brewery-FirestoneWalkerBrewingCompany_524.jpeg\",\"country_name\":\"United States\",\"contact\":{\"twitter\":\"FirestoneWalker\",\"facebook\":\"http:\\/\\/www.facebook.com\\/#!\\/firestone.walker\",\"url\":\"http:\\/\\/www.firestonebeer.com\\/\"},\"location\":{\"brewery_city\":\"Paso Robles\",\"brewery_state\":\"CA\",\"lat\":35.5953,\"lng\":-120.694}},\"auth_rating\":0,\"wish_list\":false,\"weighted_rating_score\":4.53289,\"vintages\":{\"count\":5,\"items\":[{\"beer\":{\"bid\":65409,\"beer_label\":\"https:\\/\\/untappd.akamaized.net\\/site\\/beer_logos\\/beer-Parabola2010_65409.jpeg\",\"beer_slug\":\"firestone-walker-brewing-company-parabola-2010\",\"beer_name\":\"Parabola (2010)\",\"is_vintage\":1,\"is_variant\":0}},{\"beer\":{\"bid\":68254,\"beer_label\":\"https:\\/\\/untappd.akamaized.net\\/site\\/beer_logos\\/beer-Parabola2011_68254.jpeg\",\"beer_slug\":\"firestone-walker-brewing-company-parabola-2011\",\"beer_name\":\"Parabola (2011)\",\"is_vintage\":1,\"is_variant\":0}},{\"beer\":{\"bid\":149851,\"beer_label\":\"https:\\/\\/untappd.akamaized.net\\/site\\/beer_logos\\/beer-Parabola2012_149851.jpeg\",\"beer_slug\":\"firestone-walker-brewing-company-parabola-2012\",\"beer_name\":\"Parabola (2012)\",\"is_vintage\":1,\"is_variant\":0}},{\"beer\":{\"bid\":340423,\"beer_label\":\"https:\\/\\/untappd.akamaized.net\\/site\\/beer_logos\\/beer-_340423_sm_162fd646cdd75abe8426e494700948.jpeg\",\"beer_slug\":\"firestone-walker-brewing-company-parabola-2013\",\"beer_name\":\"Parabola (2013)\",\"is_vintage\":1,\"is_variant\":0}},{\"beer\":{\"bid\":345213,\"beer_label\":\"https:\\/\\/untappd.akamaized.net\\/site\\/assets\\/images\\/temp\\/badge-beer-default.png\",\"beer_slug\":\"firestone-walker-brewing-company-parabajava\",\"beer_name\":\"Parabajava\",\"is_vintage\":0,\"is_variant\":1}}]}}}}"
