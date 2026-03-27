package tools

import (
	"sync"
)

var (
	foundCallback    func(name, url string)
	progressCallback func(sitesChecked int)
	callbackLock     sync.Mutex
)

func SetFoundCallback(callback func(name, url string)) {
	callbackLock.Lock()
	defer callbackLock.Unlock()
	foundCallback = callback
}

func SetProgressCallback(callback func(sitesChecked int)) {
	callbackLock.Lock()
	defer callbackLock.Unlock()
	progressCallback = callback
}

func NotifyFound(name, url string) bool {
	callbackLock.Lock()
	defer callbackLock.Unlock()

	if foundCallback != nil {
		foundCallback(name, url)
		return true
	}

	return false
}

func NotifyProgress(sitesChecked int) {
	callbackLock.Lock()
	defer callbackLock.Unlock()

	if progressCallback != nil {
		progressCallback(sitesChecked)
	}
}
