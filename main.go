package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-numb/go-ftx/auth"
	"github.com/go-numb/go-ftx/rest"
	"github.com/go-numb/go-ftx/rest/public/markets"
	"github.com/hannut91/upbit-go"
	"io/ioutil"
	"log"
	"net/http"
)

func ftxbtc() int {
	client := rest.New(auth.New("", ""))
	market, err := client.Markets(&markets.RequestForMarkets{
		ProductCode: "BTC/USD",
	})
	if err != nil {
		log.Fatal(err)
	}
	return int((*market)[0].Last)
}

func upbbtc() int {
	client := upbit.NewClient("", "")
	mark2ets, err := client.Ticker("KRW-BTC")
	if err != nil {
		log.Fatal(err)
	}
	return int(mark2ets[0].TradePrice)
}

func kawaserate() int {
	url := "https://quotation-api-cdn.dunamu.com/v1/forex/recent?codes=FRX.KRWUSD"
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	byteArray, _ := ioutil.ReadAll(resp.Body)
	var objmap []map[string]interface{}
	if err := json.Unmarshal(byteArray, &objmap); err != nil {
		log.Fatal(err)
	}
	return int(objmap[0]["basePrice"].(float64))
}

func calc() float32 {
	ftxbtc1 := ftxbtc()
	usdkrw := kawaserate()
	uptbtc1 := upbbtc()

	globalprice := ftxbtc1 * usdkrw
	kimp := (1 - (float32(globalprice) / float32(uptbtc1))) * 100
	return kimp
}

func main() {

	fmt.Println(calc())
}
