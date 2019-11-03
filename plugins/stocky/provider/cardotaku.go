package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	ur "net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/valyala/fasthttp"
)

type CardotakuConfigItem struct {
	Url     string `json:"url"`
	Name    string `json:"name"`
	Channel string `json:"channel"`
}

type CardotakuConfigSearch struct {
	Url     string `json:"url"`
	Channel string `json:"channel"`
}

type CardotakuConfig struct {
	Items  []CardotakuConfigItem   `json:"items"`
	Search []CardotakuConfigSearch `json:"search"`
}

type CardotakuQueryItem struct {
	Status string `json:"status"`
	Items  []struct {
		Item struct {
			ID                   int64       `json:"id"`
			Title                string      `json:"title"`
			Handle               string      `json:"handle"`
			Description          string      `json:"description"`
			PublishedAt          time.Time   `json:"published_at"`
			CreatedAt            time.Time   `json:"created_at"`
			Vendor               string      `json:"vendor"`
			Type                 string      `json:"type"`
			Tags                 []string    `json:"tags"`
			Price                int         `json:"price"`
			PriceMin             int         `json:"price_min"`
			PriceMax             int         `json:"price_max"`
			Available            bool        `json:"available"`
			PriceVaries          bool        `json:"price_varies"`
			CompareAtPrice       interface{} `json:"compare_at_price"`
			CompareAtPriceMin    int         `json:"compare_at_price_min"`
			CompareAtPriceMax    int         `json:"compare_at_price_max"`
			CompareAtPriceVaries bool        `json:"compare_at_price_varies"`
			Variants             []struct {
				ID                  int64       `json:"id"`
				Title               string      `json:"title"`
				Option1             string      `json:"option1"`
				Option2             interface{} `json:"option2"`
				Option3             interface{} `json:"option3"`
				Sku                 string      `json:"sku"`
				RequiresShipping    bool        `json:"requires_shipping"`
				Taxable             bool        `json:"taxable"`
				FeaturedImage       interface{} `json:"featured_image"`
				Available           bool        `json:"available"`
				Name                string      `json:"name"`
				PublicTitle         interface{} `json:"public_title"`
				Options             []string    `json:"options"`
				Price               int         `json:"price"`
				Weight              int         `json:"weight"`
				CompareAtPrice      interface{} `json:"compare_at_price"`
				InventoryQuantity   int         `json:"inventory_quantity"`
				InventoryManagement string      `json:"inventory_management"`
				InventoryPolicy     string      `json:"inventory_policy"`
				Barcode             string      `json:"barcode"`
			} `json:"variants"`
			Images        []string `json:"images"`
			FeaturedImage string   `json:"featured_image"`
			Options       []string `json:"options"`
			Content       string   `json:"content"`
		} `json:"item"`
		Title       string `json:"title"`
		Description string `json:"description"`
	} `json:"items"`
}

func getStockCountCardOtaku(url string) (int, error) {
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

	stockText := doc.Find(".variant-inventory span").Text()
	if stockText == "" {
		return 0, nil
	}
	if stockText == "This product is available." {
		return 10, nil
	}
	rgx := regexp.MustCompile("([0-9]+)")
	strCount := rgx.FindAllStringSubmatch(stockText, -1)
	stockCount, err := strconv.Atoi(strCount[0][1])
	if len(strCount) == 0 {
		return 0, nil
	}
	if err != nil {
		return -1, err
	}
	return stockCount, nil
}

func getSearchResultsCardOtaku(url string) ([]DataItem, error) {
	u, err := ur.Parse(url)
	if err != nil {
		return nil, err
	}
	m, err := ur.ParseQuery(u.RawQuery)
	if err != nil {
		return nil, err
	}
	query := m["q"][0]

	list := make([]DataItem, 0)

	for i := 1; i < 9; i++ {
		apiURL := fmt.Sprintf("https://cardotaku.com/search?view=ls_products&q=%s&page=%d", query, i)
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		req.SetRequestURI(apiURL)

		err := fasthttp.Do(req, resp)
		if err != nil {
			return nil, err
		}

		bodyBytes := resp.Body()
		var items CardotakuQueryItem
		err = json.Unmarshal(bodyBytes, &items)
		if err != nil {
			return nil, err
		}
		for _, item := range items.Items {
			data := DataItem{item.Item.Title, fmt.Sprintf("https://cardotaku.com/products/%s", item.Item.Handle)}
			list = append(list, data)
		}
	}
	return list, nil
}

func FindCardotaku(config *CardotakuConfig, data *Data) error {
	for _, conf := range config.Items {
		stockCount, err := getStockCountCardOtaku(conf.Url)
		if err != nil {
			return err
		}
		data.SetStock(conf.Channel, "Cardotaku", conf.Name, stockCount,
			conf.Url)
	}

	for _, conf := range config.Search {
		dataItems, err := getSearchResultsCardOtaku(conf.Url)
		if err != nil {
			return err
		}
		data.SetSearch(conf.Channel, "Cardotaku", dataItems)
	}

	return nil
}
