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
