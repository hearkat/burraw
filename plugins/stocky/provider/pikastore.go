package provider

import (
	"bytes"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/valyala/fasthttp"
)

type PikastoreConfigItem struct {
	Url     string `json:"url"`
	Name    string `json:"name"`
	Channel string `json:"channel"`
}

type PikastoreConfig struct {
	Items []PikastoreConfigItem `json:"items"`
}

func getStockCount(url string) (int, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url)

	err := fasthttp.Do(req, resp)
	if err != nil {
		return -1, err
	}

	bodyBytes := resp.Body()

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(bodyBytes))
	if err != nil {
		return -1, err
	}

	stockText := doc.Find(".PBStockTbl tbody tr:nth-child(2) td span").Text()
	rgx := regexp.MustCompile("([0-9]+)")
	strCount := rgx.FindAllStringSubmatch(stockText, -1)
	if len(strCount) == 0 {
		return 0, nil
	}
	stockCount, err := strconv.Atoi(strCount[0][1])
	if err != nil {
		return -1, err
	}

	return stockCount, nil
}

func FindPikastore(config *PikastoreConfig, data *Data) error {
	for _, conf := range config.Items {
		stockCount, err := getStockCount(conf.Url)
		if err != nil {
			return err
		}

		data.SetStock(conf.Channel, "Pikastore", conf.Name, stockCount,
			conf.Url)
	}

	return nil
}
