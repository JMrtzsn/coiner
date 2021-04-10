package binance

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"github.com/jmrtzsn/coiner/internal"
	"strconv"
	"time"
)

type Binance struct {
	client binance.Client
	config Config
}

// -------------- Interface functions  -----------------

// Init loads env variables and creates a binance rest client
func (e *Binance) Init() {
	e.config.LoadEnv()
	e.client = *binance.NewClient(e.config.apiKey, e.config.apiSecret)
}

// OHLCV validats and parses inputs.
// Returns standard trade data, according to the Open, High, Low, Volume Format
// symbol: BUYSELL = BTCUSD
// interval: 1d, 1h, 15m, 1m
// start: datetime ISO RFC3339 - "2020-04-04 T12:07:00"
// end: UNIX datetime - 1499040000000
// limit: rows returned - 10
func (e *Binance) OHLCV(symbol string, interval, start, end string) ([]internal.OHLCV, error) {
	startTS, err := isoToUnix(start)
	if err != nil {
		return nil, err
	}
	endTS, err := isoToUnix(end)
	if err != nil {
		return nil, err
	}

	klines, err := e.klines(symbol, interval, startTS, endTS)
	if err != nil {
		return nil, err
	}

	ohlcv := toOHLCV(klines)
	return ohlcv, err
}

func isoToUnix(date string) (int64, error) {
	ts, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return 0, err
	}
	return ts.UnixNano() / int64(time.Millisecond), nil
}

func unixToISO(date int64) string {
	dt := time.Unix(date/1000, 0)
	return dt.UTC().Format(time.RFC3339)
}

func toOHLCV(blines []*binance.Kline) []internal.OHLCV {
	var ohlcvs []internal.OHLCV
	for _, b := range blines {

		dt := unixToISO(b.OpenTime)
		ts := strconv.Itoa(int(b.OpenTime) / 1000)
		o := internal.OHLCV{
			DATE:   dt,
			TS:     ts,
			OPEN:   b.Open,
			CLOSE:  b.Close,
			HIGH:   b.High,
			LOW:    b.Low,
			VOLUME: b.Volume,
		}
		ohlcvs = append(ohlcvs, o)
	}
	return ohlcvs
}

// -------------- Exchange functions  -----------------

func (e *Binance) prices(symbol string) ([]*binance.SymbolPrice, error) {
	order, err := e.client.NewListPricesService().Symbol(symbol).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	return order, nil
}

// symbol: BUYSELL = BTCUSD
// interval: 1day, 1hour, 15min, 1min
// start: UNIX datetime - 1499040000000
// end: UNIX datetime - 1499040000000
// limit how many rows returned - 10
func (e *Binance) klines(symbol, interval string, start, end int64) ([]*binance.Kline, error) {
	candles, err := e.client.
		NewKlinesService().
		Symbol(symbol).
		Interval(interval).
		StartTime(start).
		EndTime(end).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	return candles, nil
}

func (e *Binance) accountBalances() ([]binance.Balance, error) {
	res, err := e.client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return nil, err
	}
	return res.Balances, nil
}

func (e *Binance) depth(symbol string) (*binance.DepthResponse, error) {
	order, err := e.client.NewDepthService().Symbol(symbol).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	return order, nil
}

// -------------- Orders -----------------

func (e *Binance) marketOrder(symbol string, order string, quantity string) (*binance.CreateOrderResponse, error) {
	side, err := setSideType(order)
	if err != nil {
		return nil, err
	}
	result, err := e.client.NewCreateOrderService().
		Symbol(symbol).
		Side(side).
		Type(binance.OrderTypeMarket).
		Quantity(quantity).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (e *Binance) orderStatus(orderId int64, symbol string) (*binance.Order, error) {
	order, err := e.client.NewGetOrderService().Symbol(symbol).
		OrderID(orderId).Do(context.Background())
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (e *Binance) cancelOrder(orderId int64, symbol string) (*binance.CancelOrderResponse, error) {
	order, err := e.client.NewCancelOrderService().Symbol(symbol).
		OrderID(orderId).Do(context.Background())
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (e *Binance) openOrders(symbol string) ([]*binance.Order, error) {
	order, err := e.client.NewListOpenOrdersService().Symbol(symbol).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	return order, nil
}

func setSideType(side string) (binance.SideType, error) {
	var sideType binance.SideType
	if side == "BUY" {
		sideType = binance.SideTypeBuy
	} else if side == "SELL" {
		sideType = binance.SideTypeSell
	} else {
		return sideType, fmt.Errorf("received invalid order type%v", side)
	}
	return sideType, nil
}
