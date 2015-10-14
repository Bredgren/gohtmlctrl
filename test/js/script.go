package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
)

var jq = jquery.NewJQuery
var console = js.Global.Get("console")

func main() {
	js.Global.Set("onBodyLoad", onBodyLoad)
}

func onBodyLoad() {
	// body := jq("body")
	console.Call("log", "hi")
}
