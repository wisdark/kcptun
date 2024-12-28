// The MIT License (MIT)
//
// # Copyright (c) 2016 xtaci
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

//go:build linux || darwin || freebsd

package std

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	kcp "github.com/xtaci/kcp-go/v5"
)

const (
	EXIT_WAIT = 5 // max seconds to wait before exit
)

func init() {
	go sigHandler()
}

func sigHandler() {
	var exitOnce sync.Once
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGUSR1, syscall.SIGTERM, syscall.SIGINT)
	signal.Ignore(syscall.SIGPIPE)

	for {
		sig := <-ch
		switch sig {
		case syscall.SIGUSR1:
			log.Printf("KCP SNMP:%+v", kcp.DefaultSnmp.Copy())
		case syscall.SIGTERM, syscall.SIGINT:
			postProcess()
			signal.Stop(ch)
			syscall.Kill(syscall.Getpid(), syscall.SIGTERM)

			// wait for max EXIT_WAIT seconds before exit
			exitOnce.Do(func() {
				go func() {
					<-time.After(EXIT_WAIT * time.Second)
					os.Exit(0)
				}()
			})
		}
	}
}
