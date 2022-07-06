package util

import (
	"fmt"
	"sync"
)

type Change struct {
	// new wants requested
	Add string
}

type SessionWantSender struct {
	Changes chan string
}

var instance *SessionWantSender
var once sync.Once

const (
	// Maximum number of changes to accept before blocking
	changesBufferSize = 128
)

func GetInstance() *SessionWantSender {
	once.Do(func() {
		instance = &SessionWantSender{
			Changes:                  make(chan string, changesBufferSize),
		}
	})
	return instance
}

func (sws *SessionWantSender) AddChange(c string) {
	fmt.Println("hi aaa")
	select {
	case sws.Changes <- c:
	}
	fmt.Println("hi bbb")
}

func (sws *SessionWantSender) Run() {
	//ch := <-sws.Changes
	//fmt.Println("hi here %s", ch)
	for {
		select {
		case ch := <-sws.Changes:
			fmt.Println("hi here %s", ch)
			//case <-sws.ctx.Done():
			//	fmt.Println("hi good")
			//	return
		}
	}
}