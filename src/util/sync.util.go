package util

import (
	"sync"
)

func Wait(handlers ...func()) {
	var wg sync.WaitGroup
	wg.Add(len(handlers))
	for _, handler := range handlers {
		go func(handler func()) {
			handler()
			wg.Done()
		}(handler)
	}
	wg.Wait()
}