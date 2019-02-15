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
	"bytes"
	"encoding/json"
	"fmt"
)

type PartialBookDepthStreamMessage struct {
	LastUpdateID int64
	Bids         []BidEntry
	Asks         []AskEntry
}

type BidEntry struct {
	Price  float64
	Volume float64
}

type AskEntry BidEntry

func DecodePartialBookDepthStream(payload []byte) (message PartialBookDepthStreamMessage, err error) {
	var object map[string]interface{}

	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.UseNumber()

	if err := decoder.Decode(&object); err != nil {
		return message, err
	}

	lastUpdateId, ok := object["lastUpdateId"].(json.Number)
	if !ok {
		return message, fmt.Errorf("invalid datatype for lastUpdateId")
	}
	message.LastUpdateID, _ = lastUpdateId.Int64()

	bids, ok := object["bids"].([]interface{})
	if !ok {
		return message, fmt.Errorf("invalid type for bids")
	}
	for _, bid := range bids {
		price, err := decodeStringFloat64(bid.([]interface{})[0])
		if err != nil {
			return message, fmt.Errorf("failed to decode bid price: %v", err)
		}

		volume, err := decodeStringFloat64(bid.([]interface{})[1])
		if err != nil {
			return message, fmt.Errorf("failed to decode bid volume: %v", err)
		}

		message.Bids = append(message.Bids, BidEntry{
			Price:  price,
			Volume: volume,
		})
	}

	asks, ok := object["asks"].([]interface{})
	if !ok {
		return message, fmt.Errorf("invalid type for asks")
	}
	for _, ask := range asks {
		price, err := decodeStringFloat64(ask.([]interface{})[0])
		if err != nil {
			return message, fmt.Errorf("failed to decode ask price: %v", err)
		}

		volume, err := decodeStringFloat64(ask.([]interface{})[1])
		if err != nil {
			return message, fmt.Errorf("failed to decode ask volume: %v", err)
		}

		message.Asks = append(message.Asks, AskEntry{
			Price:  price,
			Volume: volume,
		})
	}

	return message, nil
}
