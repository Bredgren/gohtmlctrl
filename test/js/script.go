package main

import (
	"fmt"

	"github.com/Bredgren/gohtmlctrl/htmlctrl"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
)

var jq = jquery.NewJQuery
var console = js.Global.Get("console")

func log(i ...interface{}) {
	console.Call("log", i...)
}

func logError(i ...interface{}) {
	console.Call("error", i...)
}

func logInfo(i ...interface{}) {
	console.Call("info", i...)
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
	logInfo("begin testBool")
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
		logInfo(fmt.Sprintf("test case: %#v", c))
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
	logInfo("end testBool")
}

type sliceCase interface {
	name() string
	slice() interface{}
	mms() (min, max, step float64)
	valid() htmlctrl.Validator
	error() bool
}

type sliceBoolCase struct {
	n string
	s []bool
	e bool
}

func (s *sliceBoolCase) name() string {
	return s.n
}

func (s *sliceBoolCase) slice() interface{} {
	return interface{}(&s.s)
}

func (s *sliceBoolCase) mms() (min, max, step float64) {
	return 0, 0, 0
}

func (s *sliceBoolCase) valid() htmlctrl.Validator {
	return nil
}

func (s *sliceBoolCase) error() bool {
	return s.e
}

type sliceBoolPtrCase struct {
	n string
	s []*bool
	v htmlctrl.Validator
	e bool
}

func (s *sliceBoolPtrCase) name() string {
	return s.n
}

func (s *sliceBoolPtrCase) slice() interface{} {
	return interface{}(&s.s)
}

func (s *sliceBoolPtrCase) mms() (min, max, step float64) {
	return 0, 0, 0
}

func (s *sliceBoolPtrCase) valid() htmlctrl.Validator {
	return s.v
}

func (s *sliceBoolPtrCase) error() bool {
	return s.e
}

func testSlices(body jquery.JQuery) {
	logInfo("begin testSlices")
	logInfo("begin testSlice bool")
	cases := []sliceCase{
		&sliceBoolCase{"bool1", []bool{}, true},
	}
	_, e := htmlctrl.Slice(cases[0], "error", 0, 0, 0, nil)
	if e == nil {
		logError("expected error when passing non-ptr to slice")
	}
	_, e = htmlctrl.Slice(&e, "error", 0, 0, 0, nil)
	if e == nil {
		logError("expected error when passing ptr to non-slice")
	}
	testSlice(body, cases)

	logInfo("begin testSlice *bool")
	b1, b2 := true, false
	cases = []sliceCase{
		&sliceBoolPtrCase{"*bool1", []*bool{&b1, &b2}, func(b interface{}) bool {
			log("bool is locked at true")
			return b.(bool)
		}, false},
		&sliceBoolPtrCase{"*bool2", []*bool{}, nil, false},
	}
	testSlice(body, cases)
	logInfo("end testSlices")
}

func testSlice(body jquery.JQuery, cases []sliceCase) {
	slices := jq("<div>")
	for _, c := range cases {
		logInfo(fmt.Sprintf("test case: %#v", c))
		min, max, step := c.mms()
		j, e := htmlctrl.Slice(c.slice(), c.name(), min, max, step, c.valid())
		if e != nil {
			if c.error() {
				log(fmt.Sprintf("%s: expected error: %s", c.name(), e))
			} else {
				logError(fmt.Sprintf("%s: unexpected error: %s", c.name(), e))
			}
		}
		if title := j.Attr("title"); !c.error() && title != c.name() {
			logError(fmt.Sprintf("%s: title is %s, expected %s", c.name(), title, c.name()))
		}
		slices.Append(j)
		c := c
		slices.Append(jq("<button>").SetText("verify "+c.name()).Call(jquery.CLICK, func() {
			log(c.name(), c.slice())
		}))
	}
	body.Append(slices)
}

func testStruct(body jquery.JQuery) {
	logInfo("begin testStruct")
	logInfo("end testStruct")
}
