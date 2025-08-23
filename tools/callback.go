package tools

import (
	"sync"
)

var (
	foundCallback     func(name, url string)
	foundCallbackLock sync.Mutex
)

// SetFoundCallback sets the callback function for found profiles
func SetFoundCallback(callback func(name, url string)) {
	foundCallbackLock.Lock()
	defer foundCallbackLock.Unlock()
	foundCallback = callback
}

// NotifyFound calls the callback function if set
func NotifyFound(name, url string) {
	foundCallbackLock.Lock()
	defer foundCallbackLock.Unlock()
	if foundCallback != nil {
		foundCallback(name, url)
	}
}
