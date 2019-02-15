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
	"strings"
	"time"
)

type CombinedStreamBuilder struct {
	streams []string
}

func NewCombinedStreamBuilder() *CombinedStreamBuilder {
	return &CombinedStreamBuilder{}
}

func (b *CombinedStreamBuilder) SubscribeAggTrade(symbol string) *CombinedStreamBuilder {
	stream := fmt.Sprintf("%s@aggTrade", strings.ToLower(symbol))
	b.streams = append(b.streams, stream)
	return b
}

func (b *CombinedStreamBuilder) Connect() (*Stream, error) {
	endpoint := fmt.Sprintf("stream?streams=%s", strings.Join(b.streams, "/"))
	stream, err := OpenStream(endpoint)
	if err != nil {
		return nil, err
	}
	stream.Type = STREAM_TYPE_COMBINED
	return stream, nil
}

// Stream name: <symbol>@aggTrade.
type StreamAggTrade struct {
	EventType       string  `json:"e"`
	EventTimeMillis int64   `json:"E"`
	Symbol          string  `json:"s"`
	TradeID         int64   `json:"a"`
	Price           float64 `json:"p,string"`
	Quantity        float64 `json:"q,string"`
	FirstTradeID    int64   `json:"f"`
	LastTradeID     int64   `json:"l"`
	TradeTimeMillis int64   `json:"T"`
	BuyerMaker      bool    `json:"m"`
	Ignored         bool    `json:"M"`
}

func (t *StreamAggTrade) QuoteQuantity() float64 {
	return t.Quantity * t.Price
}

func (t *StreamAggTrade) Timestamp() time.Time {
	return time.Unix(0, t.TradeTimeMillis*int64(time.Millisecond))
}

type CombinedStreamAggTrade struct {
	Stream   string         `json:"stream"`
	AggTrade StreamAggTrade `json:"data"`
}

type CombinedStreamMessage struct {
	Type     StreamType
	Stream   string
	AggTrade *StreamAggTrade
}

func (r *CombinedStreamMessage) UnmarshalJSON(b []byte) error {
	prefix := string(b[0:40])

	if strings.Index(prefix, "@aggTrade") > -1 {
		var message CombinedStreamAggTrade
		if err := json.Unmarshal(b, &message); err != nil {
			return err
		}
		r.Type = STREAM_TYPE_AGGTRADE
		r.Stream = message.Stream
		r.AggTrade = &message.AggTrade
		return nil
	}

	return fmt.Errorf("unknown stream type")
}

func DecodeStreamMessage(b []byte) (CombinedStreamMessage, error) {
	var message CombinedStreamMessage
	err := json.Unmarshal(b, &message)
	return message, err
}
