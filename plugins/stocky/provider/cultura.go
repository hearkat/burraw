package provider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

type CulturaConfigStore struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type CulturaConfigItem struct {
	Url     string `json:"url"`
	Name    string `json:"name"`
	Channel string `json:"channel"`
}

type CulturaConfigSearch struct {
	Url     string `json:"url"`
	Channel string `json:"channel"`
}

type CulturaConfig struct {
	Items  []CulturaConfigItem   `json:"items"`
	Stores []CulturaConfigStore  `json:"stores"`
	Search []CulturaConfigSearch `json:"search"`
}

type CulturaDocs struct {
	Description     string    `json:"Description"`
	Category        string    `json:"Category"`
	Size            string    `json:"Size"`
	Color           string    `json:"Color"`
	ImageURL        string    `json:"ImageUrl"`
	MinPrice        float64   `json:"MinPrice"`
	MainCategory    string    `json:"MainCategory"`
	ProductURL      string    `json:"ProductUrl"`
	Name            string    `json:"Name"`
	Brand           string    `json:"Brand"`
	DeliveryTime    string    `json:"DeliveryTime"`
	ShippingCost    string    `json:"ShippingCost"`
	Model           string    `json:"Model"`
	MinPriceSale    float64   `json:"MinPriceSale"`
	ID              string    `json:"Id"`
	SKU             string    `json:"SKU"`
	ShippingWeight  string    `json:"ShippingWeight"`
	StockStore      []string  `json:"StockStore"`
	StockCount      []int     `json:"StockCount"`
	StockPrice      []float64 `json:"StockPrice"`
	StockPriceBarre []float64 `json:"StockPriceBarre"`
}

type CulturaItem struct {
	Response struct {
		NumFound int           `json:"numFound"`
		Start    int           `json:"start"`
		Docs     []CulturaDocs `json:"docs"`
	} `json:"response"`
}

func getCultura(url string) (*CulturaItem, error) {
	rgx := regexp.MustCompile("-([0-9]+).html")
	strId := rgx.FindAllStringSubmatch(url, -1)
	id := strId[0][1]

	req, err := http.NewRequest("GET", "https://widget-cultura.proximis.com/products/select", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("q", fmt.Sprintf("SKU:%s", id))
	q.Add("wt", "json")
	q.Add("omitHeader", "true")
	req.URL.RawQuery = q.Encode()

	c := http.Client{}

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	s, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var item CulturaItem
	err = json.Unmarshal(s, &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func searchCultura(search *CulturaConfigSearch) (*string, error) {
	fmt.Println(search)

	res, err := http.Get(search.Url)
	if err != nil {
		return nil, err
	}

	s, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(s))

	return nil, err
}

func FindCultura(config *CulturaConfig, data *Data) error {
	/*for _, conf := range config.Items {
		item, err := getCultura(conf.Url)
		if err != nil {
			return err
		}

		for _, doc := range item.Response.Docs {
			for i := range doc.StockCount {
				store := doc.StockStore[i]
				for _, st := range config.Stores {
					if st.Id == store {
						data.SetStock(conf.Channel, st.Name, conf.Name, doc.StockCount[i],
							conf.Url)
					}
				}
			}
		}
	}*/

	for _, search := range config.Search {
		searchCultura(&search)
	}

	return nil
}
