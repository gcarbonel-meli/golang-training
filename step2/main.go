package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ApiResponse struct {
	ID      string              `json:"id"`
	Content *ApiResponseContent `json:"content,omitempty"`
	Partial bool                `json:"partial"`
}

type ApiResponseContent struct {
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
}

type CoinGeckoApi []struct {
	ID                           string    `json:"id"`
	Symbol                       string    `json:"symbol"`
	Name                         string    `json:"name"`
	Image                        string    `json:"image"`
	CurrentPrice                 float64   `json:"current_price"`
	MarketCap                    int64     `json:"market_cap"`
	MarketCapRank                int       `json:"market_cap_rank"`
	FullyDilutedValuation        int64     `json:"fully_diluted_valuation"`
	TotalVolume                  int64     `json:"total_volume"`
	High24H                      float64   `json:"high_24h"`
	Low24H                       float64   `json:"low_24h"`
	PriceChange24H               float64   `json:"price_change_24h"`
	PriceChangePercentage24H     float64   `json:"price_change_percentage_24h"`
	MarketCapChange24H           int64     `json:"market_cap_change_24h"`
	MarketCapChangePercentage24H float64   `json:"market_cap_change_percentage_24h"`
	CirculatingSupply            float64   `json:"circulating_supply"`
	TotalSupply                  float64   `json:"total_supply"`
	MaxSupply                    float64   `json:"max_supply"`
	Ath                          float64   `json:"ath"`
	AthChangePercentage          float64   `json:"ath_change_percentage"`
	AthDate                      time.Time `json:"ath_date"`
	Atl                          float64   `json:"atl"`
	AtlChangePercentage          float64   `json:"atl_change_percentage"`
	AtlDate                      time.Time `json:"atl_date"`
	Roi                          any       `json:"roi"`
	LastUpdated                  time.Time `json:"last_updated"`
}

func main() {
	ginEngine := gin.Default()
	ginEngine.GET("/myapi", myApiHandler)
	ginEngine.Run()
}

func myApiHandler(context *gin.Context) {
	fiatCurrency := "usd"
	cryptoCurrency := "bitcoin"
	fiatCurrencyQueried := context.Request.URL.Query().Get("fiat")
	if fiatCurrencyQueried != "" {
		fiatCurrency = fiatCurrencyQueried
	}
	cryptoCurrencyQueried := context.Request.URL.Query().Get("crypto")
	if cryptoCurrencyQueried != "" {
		cryptoCurrency = cryptoCurrencyQueried
	}
	partial := false
	// var httpStatus int
	httpStatus := http.StatusOK
	apiResponseContent, err := cryptoPrice(cryptoCurrency, fiatCurrency)
	if err != nil {
		httpStatus = http.StatusPartialContent
		partial = true
	}
	apiResponse := ApiResponse{}
	apiResponse.ID = cryptoCurrency
	if !partial {
		apiResponse.Content = apiResponseContent
	}
	apiResponse.Partial = partial
	context.JSON(httpStatus, apiResponse)
}

func cryptoPrice(crypto string, fiat string) (*ApiResponseContent, error) {
	response, err := http.Get("https://api.coingecko.com/api/v3/coins/markets?vs_currency=" + fiat + "&ids=" + crypto)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	coinGeckoApi := CoinGeckoApi{}
	err = json.Unmarshal(body, &coinGeckoApi)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if len(coinGeckoApi) == 0 {
		err = errors.New("There are no results for this parameters")
		log.Println(err)
		return nil, err
	}
	cryptoData := coinGeckoApi[0]
	apiResponseContent := ApiResponseContent{}
	apiResponseContent.Currency = fiat
	apiResponseContent.Price = cryptoData.CurrentPrice
	return &apiResponseContent, nil
}
