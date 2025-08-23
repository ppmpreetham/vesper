package tools

import (
	"sync"
)

var (
	foundCallback    func(name, url string)
	progressCallback func(sitesChecked int)
	callbackLock     sync.Mutex
)

// SetFoundCallback sets the callback function for found profiles
func SetFoundCallback(callback func(name, url string)) {
	callbackLock.Lock()
	defer callbackLock.Unlock()
	foundCallback = callback
}

// SetProgressCallback sets the callback function for progress updates
func SetProgressCallback(callback func(sitesChecked int)) {
	callbackLock.Lock()
	defer callbackLock.Unlock()
	progressCallback = callback
}

// NotifyFound calls the callback function if set and returns whether it was handled
func NotifyFound(name, url string) bool {
	callbackLock.Lock()
	defer callbackLock.Unlock()

	if foundCallback != nil {
		foundCallback(name, url)
		return true
	}

	return false
}

// NotifyProgress calls the progress callback function if set
func NotifyProgress(sitesChecked int) {
	callbackLock.Lock()
	defer callbackLock.Unlock()

	if progressCallback != nil {
		progressCallback(sitesChecked)
	}
}
