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
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
)

var STREAM_URL = "wss://stream.binance.com:9443"

type StreamType int

const (
	STREAM_TYPE_COMBINED          StreamType = 0
	STREAM_TYPE_AGGTRADE          StreamType = 1
	STREAM_TYPE_PARTIAL_BOOK      StreamType = 2
	STREAM_TYPE_ALL_MARKET_TICKER StreamType = 3
)

type Stream struct {
	Conn *websocket.Conn
	Type StreamType
}

func OpenStream(stream string) (*Stream, error) {
	url := fmt.Sprintf("%s/%s", STREAM_URL, stream)
	conn, response, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusSwitchingProtocols {
		return nil, fmt.Errorf("failed to upgrade to websocket: %s", response.Status)
	}
	return &Stream{
		Conn: conn,
	}, nil
}

func OpenSingleStream(stream string) (*Stream, error) {
	url := fmt.Sprintf("%s/ws/%s", STREAM_URL, stream)
	conn, response, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusSwitchingProtocols {
		return nil, fmt.Errorf("failed to upgrade to websocket: %s", response.Status)
	}
	return &Stream{
		Conn: conn,
	}, nil
}

func (s *Stream) Close() {
	s.Conn.Close()
}

func (s *Stream) Next() ([]byte, error) {
	_, payload, err := s.Conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func OpenPartialBookDepthStream(symbol string, depth int) (*Stream, error) {
	stream, err := OpenSingleStream(fmt.Sprintf("%s@depth%d", strings.ToLower(symbol), depth))
	if err != nil {
		return nil, err
	}
	stream.Type = STREAM_TYPE_PARTIAL_BOOK
	return stream, nil
}

func OpenAllMarketTickerStream() (*Stream, error) {
	stream, err := OpenSingleStream("!ticker@arr")
	if err != nil {
		return nil, err
	}
	stream.Type = STREAM_TYPE_PARTIAL_BOOK
	return stream, nil
}
