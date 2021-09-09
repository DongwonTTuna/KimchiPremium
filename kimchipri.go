package main

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"github.com/hannut91/upbit-go"
	ex "github.com/me-io/go-swap/pkg/exchanger"
	"github.com/me-io/go-swap/pkg/swap"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	menu = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	btnadd    = menu.Text("⊕ 추가")
	btnview   = menu.Text("Θ 보기")
)
var (
	off = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	offbutton = off.Text("● 끄기")
)

var bp *tb.Bot
var client *upbit.Client
var user = &tb.User{ID: USERID}
var poller = &tb.LongPoller{Timeout: 15 * time.Second}
var spamProtected = tb.NewMiddlewarePoller(poller, func(upd *tb.Update) bool {
	if upd.Message == nil {
		return true
	}

	if strings.Contains(upd.Message.Text, "추가") {
		return true
	} else if strings.Contains(upd.Message.Text, "끄기") {
		return true
	} else if strings.Contains(upd.Message.Text, "보기") {
		return true
	}

	dd := strings.Split(upd.Message.Text, " ")
	fmt.Println(dd[1],dd[2])
	_, err := strconv.Atoi(dd[1])
	if err != nil{
		_, _ = bp.Send(user, "퍼센트를 잘못 입력하셨습니다.")
		return true
	}
	_, err = strconv.Atoi(dd[2])
	if err == nil {
		_, _ = bp.Send(user, "방향을 잘못 입력하셨습니다.")
		return true
	}
	if dd[2] != "이상"{
		if dd[2] != "이하"{
			_, _ = bp.Send(user, "방향을 잘못 입력하셨습니다.")
			return true
		}
	}
	err = ioutil.WriteFile("kimchipre.txt", []byte(dd[1]+ " " + dd[2]), os.FileMode(644))
	if err != nil {
		_, _ = bp.Send(user, "퍼센트를 추가하는데 오류가 발생하였습니다.")
	}
	err = ioutil.WriteFile("off_k.txt", []byte("0"), os.FileMode(644))
	if err != nil {
		_, _ = bp.Send(user, "알람을 활성화시키는데 오류가 발생하였습니다.")
	}
	_, _ = bp.Send(user, " 추가 완료")
	return true
})
func main() {
	_, err := os.Open("off_k.txt")
	if err != nil {
		_, _ = os.Create("off_k.txt")
	}
	_, err = os.Open("kimchipre.txt")
	if err != nil {
		_, _ = os.Create("kimchipre.txt")
	}
	bot, err := tb.NewBot(tb.Settings{
		Token:  "TELEGRAM_TOKEN",
		Poller: spamProtected,
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	bp = bot
	menu.Reply(
		menu.Row(btnadd),
		menu.Row(btnview),
	)
	off.Reply(
		off.Row(offbutton),
	)
	_, _ = bp.Send(user, "", menu)
	bp.Handle(&btnadd, func(m *tb.Message) {
		_, _ = bp.Send(user, "김프퍼센트 이상/이하 양식으로 입력해주세요.")
	})
	bp.Handle(&btnview, func(m *tb.Message) {
		data, _ := ioutil.ReadFile("kimchipre.txt")
		_, _ = bp.Send(user, "김프 감시 설정 : " + string(data))
	})
	bp.Handle(&offbutton, func(m *tb.Message) {
		err = ioutil.WriteFile("off_k.txt", []byte("1"), os.FileMode(644))
		if err != nil {
			_, _ = bp.Send(user, "알람을 끄는데 실패하였습니다.", menu)
		}
		err = ioutil.WriteFile("kimchipre.txt", []byte(""), os.FileMode(644))
		if err != nil {
			_, _ = bp.Send(user, "김프 조건을 초기화시키는데 실패하였습니다.")
		}
		_, _ = bp.Send(user, "꺼짐!", menu)
	})
	go NeverExit(kimchi)
	client = upbit.NewClient("", "")
	bp.Start()
}
func kimchipri() float64{
	markets, err := client.MinuteCandles(1,"KRW-BTC")
	if err != nil {
		_, _ = bp.Send(user, "업비트 가격을 불러오는데 실패하였습니다.")
		return -101
	}
	upprice := markets[0].TradePrice
	biprice := func() float64{
		client := binance.NewClient("", "")
		prices, err := client.NewListPricesService().Symbol("BTCUSDT").Do(context.Background())
		if err != nil {
			_, _ = bp.Send(user, "바이낸스 비트코인 가격을 요청할 수 없습니다.")
			return 0
		}
		SwapTest := swap.NewSwap()
		SwapTest.
			AddExchanger(ex.NewYahooApi(nil)).
			Build()
		euroToUsdRate := SwapTest.Latest("USD/KRW")
		a,_ := strconv.ParseFloat(prices[0].Price,64)
		return a * euroToUsdRate.GetRateValue()
	}()
	if biprice == 0{
		return -101
	}
	kimchi_va := (upprice - biprice) / biprice * 100
	return kimchi_va
}
func kimchi() {
	time.Sleep(5*time.Second)
	menu.Reply(
		menu.Row(btnadd),
		menu.Row(btnview),
	)

	for {
		data, _ := ioutil.ReadFile("kimchipre.txt")
		if string(data) == ""{
			time.Sleep(5*time.Second)
			continue
		}
		datea := strings.Split(string(data), " ")
		targetpre,_ := strconv.ParseFloat(datea[0],64)
		way := datea[1]
		pre := kimchipri()
		if pre == -101{
			time.Sleep(3600*time.Second)
			continue
		}
		if way == "이상" && pre >= targetpre {
			for {
				offm, err := ioutil.ReadFile("off_k.txt")
					if err != nil {
						_, _ = bp.Send(user, "알림을 껐는지 확인하는데 실패하였습니다.", menu)
						break
					}
					if string(offm) == "1" {
						break
					}
					_, _ = bp.Send(user, "일어나세요!", off)
				}
				time.Sleep(5 * time.Second)
		} else if way == "이하" && pre <= targetpre{
			for {
				offm, err := ioutil.ReadFile("off_k.txt")
				if err != nil {
					_, _ = bp.Send(user, "알림을 껐는지 확인하는데 실패하였습니다.", menu)
					break
				}
				if string(offm) == "1" {
					break
				}
				_, _ = bp.Send(user, "일어나세요!", off)
			}
			time.Sleep(5 * time.Second)
		}
	}
}
func NeverExit(f func()) {
	defer func() { if v := recover(); v != nil {
		log.Println(v)
		go NeverExit(f)
	}
	}()
	f()
}
