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

type PriceTickerResponse struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price,string"`
}

func (c *RestClient) GetPriceTicker(symbol string) (PriceTickerResponse, error) {
	endpoint := "/api/v3/ticker/price"
	var response PriceTickerResponse
	params := map[string]interface{}{
		"symbol": symbol,
	}
	err := c.GetAndDecode(endpoint, params, &response)
	return response, err
}

func (c *RestClient) GetPriceTickerAll() ([]PriceTickerResponse, error) {
	endpoint := "/api/v3/ticker/price"
	var response []PriceTickerResponse
	err := c.GetAndDecode(endpoint, nil, &response)
	return response, err
}

type BookTickerResponse struct {
	Symbol    string  `json:"symbol"`
	BidPrice  float64 `json:"bidPrice,string"`
	BidVolume float64 `json:"bidQty,string"`
	AskPrice  float64 `json:"askPrice,string"`
	AskVolume float64 `json:"askQty,string"`
}

func (c *RestClient) GetBookTicker(symbol string) (BookTickerResponse, error) {
	endpoint := "/api/v3/ticker/bookTicker"
	var response BookTickerResponse
	params := map[string]interface{}{
		"symbol": symbol,
	}
	err := c.GetAndDecode(endpoint, params, &response)
	return response, err
}

type UserDataStreamResponse struct {
	ListenKey string `json:"listenKey"`
}

func (c *RestClient) GetUserDataStream() (string, error) {
	httpResponse, err := c.PostWithApiKey("/api/v1/userDataStream", nil)
	if err != nil {
		return "", err
	}

	if httpResponse.StatusCode >= 400 {
		return "", NewRestApiErrorFromResponse(httpResponse)
	}

	var response UserDataStreamResponse
	if err := c.decodeBody(httpResponse, &response); err != nil {
		return "", err
	}

	return response.ListenKey, nil
}

func (c *RestClient) PutUserStreamKeepAlive(listenKey string) error {
	queryString := c.BuildQueryString(map[string]interface{}{
		"listenKey": listenKey,
	})
	path := fmt.Sprintf("/api/v1/userDataStream?%s", queryString)
	httpResponse, err := c.PutWithApiKey(path)
	if err != nil {
		return err
	}
	if httpResponse.StatusCode != http.StatusOK {
		return NewRestApiErrorFromResponse(httpResponse)
	}
	return nil
}

type QueryOrderResponse struct {
	Symbol        string      `json:"symbol"`
	OrderId       int64       `json:"orderId"`
	ClientOrderId string      `json:"clientOrderId"`
	Price         float64     `json:"price,string"`
	OrigQty       float64     `json:"origQty,string"`
	ExecutedQty   float64     `json:"executeQty,string"`
	Status        OrderStatus `json:"status"`
	TimeInForce   TimeInForce `json:"timeInForce"`
	Type          OrderType   `json:"type"`
	Side          OrderSide   `json:"side"`
	StopPrice     float64     `json:"stopPrice,string"`
	IcebergQty    float64     `json:"icebergQty,string"`
	TimeMillis    int64       `json:"time"`
	IsWorking     bool        `json:"isWorking"`
}

func (c *RestClient) GetOrderByOrderId(symbol string, orderId int64) (QueryOrderResponse, error) {
	var response QueryOrderResponse
	params := map[string]interface{}{
		"symbol":  symbol,
		"orderId": orderId,
	}
	httpResponse, err := c.GetWithAuth("/api/v3/order", params)
	if err != nil {
		return response, err
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode != http.StatusOK {
		return response, NewRestApiErrorFromResponse(httpResponse)
	}
	decoder := json.NewDecoder(httpResponse.Body)
	if err := decoder.Decode(&response); err != nil {
		return response, err
	}
	return response, nil
}

func (c *RestClient) GetOrderByClientId(symbol string, clientId string) (QueryOrderResponse, error) {
	var response QueryOrderResponse
	params := map[string]interface{}{
		"symbol":            symbol,
		"origClientOrderId": clientId,
	}
	httpResponse, err := c.GetWithAuth("/api/v3/order", params)
	if err != nil {
		return response, err
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode != http.StatusOK {
		return response, NewRestApiErrorFromResponse(httpResponse)
	}
	decoder := json.NewDecoder(httpResponse.Body)
	if err := decoder.Decode(&response); err != nil {
		return response, err
	}
	return response, nil
}

type MyTradesResponseEntry struct {
	ID              int64   `json:"id"`
	OrderID         int64   `json:"orderId"`
	Price           float64 `json:"price,string"`
	Quantity        float64 `json:"qty,string"`
	Commission      float64 `json:"commission,string"`
	CommissionAsset string  `json:"commissionAsset"`
	TimeMillis      int64   `json:"time"`
	IsBuyer         bool    `json:"isBuyer"`
	IsMaker         bool    `json:"isMaker"`
	IsBestMatch     bool    `json:"isBestMatch"`
}

func (c *RestClient) GetMytrades(symbol string, limit int64, fromId int64) ([]MyTradesResponseEntry, error) {
	endpoint := "/api/v3/myTrades"
	params := map[string]interface{}{
		"symbol": symbol,
	}
	if limit > 0 {
		params["limit"] = limit
	}
	if fromId > -1 {
		params["fromId"] = fromId
	}
	var response []MyTradesResponseEntry
	err := c.AuthGetAndDecode(endpoint, params, &response)
	return response, err
}

type AccountInfoBalance struct {
	Asset  string  `json:"asset"`
	Free   float64 `json:"free,string"`
	Locked float64 `json:"locked,string"`
}

type AccountInfoResponse struct {
	MakerCommission  int64                `json:"makerCommission"`
	TakerCommission  int64                `json:"takerCommission"`
	BuyerCommission  int64                `json:"buyerCommission"`
	SellCommission   int64                `json:"sellCommission"`
	CanTrade         bool                 `json:"canTrade"`
	CanWithdraw      bool                 `json:"canWithdraw"`
	CanDeposit       bool                 `json:"canDeposit"`
	UpdateTimeMillis int64                `json:"updateTime"`
	Balances         []AccountInfoBalance `json:"balances"`
}

func (c *RestClient) GetAccount() (*AccountInfoResponse, error) {
	httpResponse, err := c.GetWithAuth("/api/v3/account", nil)
	if err != nil {
		return nil, err
	}
	if httpResponse.StatusCode >= 400 {
		return nil, NewRestApiErrorFromResponse(httpResponse)
	}

	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}

	var response AccountInfoResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
