package main

import (
	"fmt"

	"github.com/Bredgren/gohtmlctrl/htmlctrl"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
)

var jq = jquery.NewJQuery
var console = js.Global.Get("console")

func logError(i ...interface{}) {
	console.Call("error", i...)
}

func log(i ...interface{}) {
	console.Call("log", i...)
}

func main() {
	js.Global.Set("onBodyLoad", onBodyLoad)
}

func onBodyLoad() {
	body := jq("body")
	funcs := []func(jquery.JQuery){
		testBool,
	}
	for _, fn := range funcs {
		fn(body)
	}
}

func testBool(body jquery.JQuery) {
	log("begin testBool")
	cases := []struct {
		name  string
		b     bool
		valid func(interface{}) bool
	}{
		{"b1", false, nil},
		{"b2", true, nil},
		{"b3", true, func(b interface{}) bool {
			log("b3 is locked at true")
			return b.(bool)
		}},
		{"b4", false, func(b interface{}) bool {
			log("b4 is locked at false")
			return !b.(bool)
		}},
	}
	bools := jq("<div>")
	for _, c := range cases {
		log(fmt.Sprintf("test case: %#v", c))
		j, e := htmlctrl.Bool(&c.b, c.name, c.valid)
		if e != nil {
			logError(fmt.Sprintf("%s: ", c.name))
		}
		if b := j.Prop("checked").(bool); b != c.b {
			logError(fmt.Sprintf("%s: checked was %t, expected %t", c.name, b, c.b))
		}
		bools.Append(j)
		c := &c
		bools.Append(jq("<button>").SetText("verify "+c.name).Call(jquery.CLICK, func() {
			log(c.name, c.b)
		}))
	}
	body.Append(bools)
	log("end testBool")
}

// func testStruct(body jquery.JQuery) {
// 	log("begin testStruct")
// 	log("end testStruct")
// }
