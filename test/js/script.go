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
		testSlices,
		testStruct,
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
			logError(fmt.Sprintf("%s: unexpected error: %s", c.name, e))
		}
		if b := j.Prop("checked").(bool); b != c.b {
			logError(fmt.Sprintf("%s: checked was %t, expected %t", c.name, b, c.b))
		}
		if title := j.Attr("title"); title != c.name {
			logError(fmt.Sprintf("%s: title is %s, expected %s", c.name, title, c.name))
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

type sliceCase interface {
	name() string
	slice() interface{}
}

type sliceBoolCase struct {
	n string
	s []bool
}

func (s sliceBoolCase) name() string {
	return s.n
}

func (s sliceBoolCase) slice() interface{} {
	return interface{}(&s.s)
}

type sliceBoolPtrCase struct {
	n string
	s []*bool
}

func (s sliceBoolPtrCase) name() string {
	return s.n
}

func (s sliceBoolPtrCase) slice() interface{} {
	return interface{}(&s.s)
}

func testSlices(body jquery.JQuery) {
	log("begin testSlices")
	log("begin testSlice bool")
	cases := []sliceCase{
		sliceBoolCase{"s1", []bool{}},
		sliceBoolCase{"s2", []bool{true, false}},
	}
	_, e := htmlctrl.Slice(cases[0], "error")
	if e == nil {
		logError("expected error when passing non-ptr to slice")
	}
	_, e = htmlctrl.Slice(&e, "error")
	if e == nil {
		logError("expected error when passing ptr to non-slice")
	}
	testSlice(body, cases)

	log("begin testSlice *bool")
	b1, b2 := true, false
	cases = []sliceCase{
		sliceBoolPtrCase{"s1", []*bool{}},
		sliceBoolPtrCase{"s2", []*bool{&b1, &b2}},
	}
	testSlice(body, cases)
	log("end testSlices")
}

func testSlice(body jquery.JQuery, cases []sliceCase) {
	slices := jq("<div>")
	for _, c := range cases {
		log(fmt.Sprintf("test case: %#v", c))
		j, e := htmlctrl.Slice(c.slice(), c.name())
		if e != nil {
			logError(fmt.Sprintf("%s: unexpected error: %s", c.name(), e))
		}
		if title := j.Attr("title"); title != c.name() {
			logError(fmt.Sprintf("%s: title is %s, expected %s", c.name(), title, c.name()))
		}
		slices.Append(j)
	}
	body.Append(slices)
}

func testStruct(body jquery.JQuery) {
	log("begin testStruct")
	log("end testStruct")
}
