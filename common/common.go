package common

import (
	"syscall/js"
)

var LoadingActive = true

func StartLoading() {
	jquery := js.Global().Get("$")
	if !jquery.Equal(js.Undefined()) {
		jquery.Invoke("#loadingScreen").Call("fadeIn", "slow")
		jquery.Invoke("#actions").Call("fadeOut", "slow")
		jquery.Invoke("#loadingIndicator").Call("addClass", "active")
	}
	LoadingActive = true
}

func StopLoading() {
	LoadingActive = false
	jquery := js.Global().Get("$")
	if !jquery.Equal(js.Undefined()) {
		jquery.Invoke("#loadingScreen").Call("fadeOut", "slow")
		jquery.Invoke("#actions").Call("fadeIn", "slow")
		jquery.Invoke("#loadingIndicator").Call("fadeOut", "slow")
	}
}
