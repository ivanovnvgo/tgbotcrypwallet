package getprice

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type bnResp struct {
	Price float64 `json:"price,string"`
	Code  int64   `json:"code"`
}

func GetPrice(symbol string) (price float64, err error) {

	resp, err := http.Get(fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%sUSDT", symbol))
	fmt.Println("resp = ", resp)
	if err != nil {
		log.Printf("\nhttp.Get read error: %s", err)
		return
	}
	defer resp.Body.Close()
	var jsonResp bnResp
	err = json.NewDecoder(resp.Body).Decode(&jsonResp)
	if err != nil {
		log.Printf("\njson.NewDecoder error write: %s", err)
		return
	}
	if jsonResp.Code != 0 {
		err = errors.New("invalid cryptocurrency symbol")
	}

	price = jsonResp.Price
	return
}
