// MIT License
//
// Copyright (c) 2019 Cranky Kernel
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package binanceapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type OrderSide string

const (
	OrderSideBuy  OrderSide = "BUY"
	OrderSideSell OrderSide = "SELL"
)

type OrderType string

const (
	OrderTypeLimit  OrderType = "LIMIT"
	OrderTypeMarket OrderType = "MARKET"
)

type TimeInForce string

const (
	TimeInForceGTC TimeInForce = "GTC"
	TimeInForceIOC TimeInForce = "IOC"
	TimeInForceFOK TimeInForce = "FOK"
)

type OrderStatus string

const (
	OrderStatusNew             OrderStatus = "NEW"
	OrderStatusCanceled        OrderStatus = "CANCELED"
	OrderStatusFilled          OrderStatus = "FILLED"
	OrderStatusPartiallyFilled OrderStatus = "PARTIALLY_FILLED"
)

type OrderParameters struct {
	Symbol           string
	Side             OrderSide
	Type             OrderType
	TimeInForce      TimeInForce
	Quantity         float64
	Price            float64
	NewClientOrderId string
}

// TODO: Implement RESULT and FULL response types. Currently only ACK implemented.
type PostOrderResponse struct {
	Symbol                string `json:"symbol"`
	OrderId               int64  `json:"orderId"`
	ClientOrderId         string `json:"clientOrderId"`
	TransactionTimeMillis int64  `json:"transactTime"`
}

func (c *RestClient) PostOrder(order OrderParameters) (*http.Response, error) {
	params := map[string]interface{}{}
	params["symbol"] = order.Symbol
	params["side"] = order.Side
	params["type"] = order.Type
	params["quantity"] = fmt.Sprintf("%.8f", order.Quantity)

	switch order.Type {
	case OrderTypeMarket:
	default:
		params["price"] = fmt.Sprintf("%.8f", order.Price)
	}
	params["newClientOrderId"] = order.NewClientOrderId
	if order.TimeInForce != "" {
		params["timeInForce"] = order.TimeInForce
	}

	response, err := c.Post("/api/v3/order", params)
	if err != nil {
		return nil, err
	}
	if response.StatusCode >= 400 {
		return response, NewRestApiErrorFromResponse(response)
	}
	return response, nil
}

type CancelOrderResponse struct {
	Symbol            string `json:"symbol"`
	OrigClientOrderID string `json:"origClientOrderId"`
	OrderID           int64  `json:"orderId"`
	ClientOrderID     string `json:"clientOrderId"`
}

func (c *RestClient) CancelOrderById(symbol string, orderId int64) (CancelOrderResponse, error) {
	var cancelOrderResponse CancelOrderResponse
	params := map[string]interface{}{}
	params["symbol"] = symbol
	params["orderId"] = orderId

	httpResponse, err := c.Delete("/api/v3/order", params)
	if err != nil {
		return cancelOrderResponse, err
	}

	if httpResponse.StatusCode != http.StatusOK {
		return cancelOrderResponse, NewRestApiErrorFromResponse(httpResponse)
	}

	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return cancelOrderResponse, err
	}

	if err := json.Unmarshal(body, &cancelOrderResponse); err != nil {
		return cancelOrderResponse, err
	}
	return cancelOrderResponse, nil
}
