package common

import (
	"syscall/js"
)

var IsLoading = true

func StartLoading() {
	if IsLoading {
		return
	}
	IsLoading = true

	jquery := js.Global().Get("$")
	if !jquery.Equal(js.Undefined()) {
		jquery.Invoke("#loadingScreen").Call("fadeIn", "fast")
		jquery.Invoke("#actions").Call("fadeOut", "fast")
		jquery.Invoke("#loadingIndicator").Call("addClass", "active")
	}

}

func StopLoading() {
	if !IsLoading {
		return
	}
	IsLoading = false

	jquery := js.Global().Get("$")
	if !jquery.Equal(js.Undefined()) {
		jquery.Invoke("#loadingScreen").Call("fadeOut", "fast")
		jquery.Invoke("#actions").Call("fadeIn", "fast")
		jquery.Invoke("#loadingIndicator").Call("fadeOut", "fast")
	}
}
