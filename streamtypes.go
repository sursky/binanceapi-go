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
	"time"
)

// Stream name: <symbol>@ticker.
type TickerStreamMessage struct {
	EventType            string  `json:"e"`
	EventTime            int64   `json:"E"`
	Symbol               string  `json:"s"`
	PriceChange          float64 `json:"p,string"`
	PriceChangePercent   float64 `json:"P,string"`
	WeightedAveragePrice float64 `json:"w,string"`
	PreviousDayClose     float64 `json:"x,string"`
	CurrentDayClose      float64 `json:"c,string"`
	CloseTradeQuantity   float64 `json:"Q,string"`
	Bid                  float64 `json:"b,string"`
	BidQuantity          float64 `json:"B,string"`
	Ask                  float64 `json:"a,string"`
	AskQuantity          float64 `json:"A,string"`
	OpenPrice            float64 `json:"o,string"`
	HighPrice            float64 `json:"h,string"`
	LowPrice             float64 `json:"l,string"`
	TotalBaseVolume      float64 `json:"v,string"`
	TotalQuoteVolume     float64 `json:"q,string"`
	StatsOpenTime        int64   `json:"O"`
	StatsCloseTime       int64   `json:"C"`
	FirstTradeID         int64   `json:"F"`
	LastTradeID          int64   `json:"L"`
	TotalNumberTrades    int64   `json:"n"`
}

func (t *TickerStreamMessage) Timestamp() time.Time {
	return time.Unix(0, t.EventTime*int64(time.Millisecond))
}

func DecodeAllMarketTickerStream(payload []byte) ([]TickerStreamMessage, error) {
	var message []TickerStreamMessage
	if err := json.Unmarshal(payload, &message); err != nil {
		return nil, err
	}
	return message, nil
}
