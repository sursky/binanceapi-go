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

package main

import (
	"fmt"
	"github.com/crankykernel/binanceapi-go"
	"os"
)

type Command struct {
	name    string
	help    string
	handler func(args []string)
}

func main() {
	commands := map[string]Command{
		"stream": Command{
			name:    "stream",
			help:    "Connect to a websocket stream",
			handler: streamHandler,
		},
	}

	if len(os.Args) < 2 {
		fmt.Printf("Nothing to do.\n")
		return
	}

	command, ok := commands[os.Args[1]]
	if !ok {
		fmt.Printf("error: unknown command\n")
		return
	}

	command.handler(os.Args[2:])
}

func streamHandler(args []string) {
	stream, err := binanceapi.OpenSingleStream(args[0])
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	for {
		body, err := stream.Next()
		if err != nil {
			fmt.Printf("error: %+v\n", err)
			return
		}
		fmt.Printf("%s\n", string(body))
	}
}
