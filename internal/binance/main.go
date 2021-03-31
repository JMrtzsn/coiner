package binance

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2"
)

type Binance struct {
	client binance.Client
	config Config
}

func (e *Binance) Init() {
	client := *binance.NewClient(e.config.apiKey, e.config.apiSecret)
	e.client = client
}


// -------------- Historical Data  -----------------

func (e *Binance) Prices(symbol string) ([]*binance.SymbolPrice, error) {
	order, err := e.client.NewListPricesService().Symbol(symbol).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	return order, nil
}

// Klines creates a historical
func (e *Binance) Klines(symbol string, interval string, start, end int64, limit int) ([]*binance.Kline, error) {
	candles, err := e.client.
		NewKlinesService().
		Symbol(symbol).
		Interval(interval).
		StartTime(start).
		EndTime(end).
		Limit(limit).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	return candles, nil
}

func (e *Binance) AccountBalance() ([]binance.Balance, error) {
	res, err := e.client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return nil, err
	}
	return res.Balances, nil
}


func (e *Binance) Depth(symbol string) (*binance.DepthResponse, error) {
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
