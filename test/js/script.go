package main

import (
	"fmt"
	"math"

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
		testInt,
		testFloat64,
		testString,
		testChoice,
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
		valid htmlctrl.Validator
	}{
		{"b1", false, nil},
		{"b2", true, nil},
		{"b3", true, htmlctrl.ValidateBool(func(b bool) bool {
			log("b3 is locked at true")
			return b
		})},
		{"b4", false, htmlctrl.ValidateBool(func(b bool) bool {
			log("b4 is locked at false")
			return !b
		})},
	}
	bools := jq("<div>").AddClass("bools")
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

func testInt(body jquery.JQuery) {
	logInfo("begin testInt")
	cases := []struct {
		name           string
		i              int
		min, max, step float64
		valid          htmlctrl.Validator
	}{
		{"i1", 0, -10, 10, 3, nil},
		{"i2", 2, -100, 100, 1, htmlctrl.ValidateInt(func(i int) bool {
			if i == 5 {
				log("i can't be 5")
			}
			return i != 5
		})},
		{"i3", 0, math.NaN(), math.NaN(), math.NaN(), nil},
	}
	ints := jq("<div>").AddClass("ints")
	for _, c := range cases {
		logInfo(fmt.Sprintf("test case: %#v", c))
		j, e := htmlctrl.Int(&c.i, c.name, c.min, c.max, c.step, c.valid)
		if e != nil {
			logError(fmt.Sprintf("%s: unexpected error: %s", c.name, e))
		}
		if title := j.Attr("title"); title != c.name {
			logError(fmt.Sprintf("%s: title is %s, expected %s", c.name, title, c.name))
		}
		ints.Append(j)
		c := &c
		ints.Append(jq("<button>").SetText("verify "+c.name).Call(jquery.CLICK, func() {
			log(c.name, c.i)
		}))
	}
	body.Append(ints)
	logInfo("end testInt")
}

func testFloat64(body jquery.JQuery) {
	logInfo("begin testFloat64")
	cases := []struct {
		name           string
		f              float64
		min, max, step float64
		valid          htmlctrl.Validator
	}{
		{"f1", 0.5, -10, 10, 1.5, nil},
		{"f2", 2.1, -100, 100, 1, htmlctrl.ValidateFloat64(func(f float64) bool {
			if f == 5.5 {
				log("f can't be 5.5")
			}
			return f != 5.5
		})},
		{"f3", 0, math.NaN(), math.NaN(), math.NaN(), nil},
	}
	float64s := jq("<div>").AddClass("float64s")
	for _, c := range cases {
		logInfo(fmt.Sprintf("test case: %#v", c))
		j, e := htmlctrl.Float64(&c.f, c.name, c.min, c.max, c.step, c.valid)
		if e != nil {
			logError(fmt.Sprintf("%s: unexpected error: %s", c.name, e))
		}
		if title := j.Attr("title"); title != c.name {
			logError(fmt.Sprintf("%s: title is %s, expected %s", c.name, title, c.name))
		}
		float64s.Append(j)
		c := &c
		float64s.Append(jq("<button>").SetText("verify "+c.name).Call(jquery.CLICK, func() {
			log(c.name, c.f)
		}))
	}
	body.Append(float64s)
	logInfo("end testFloat64")
}

func testString(body jquery.JQuery) {
	logInfo("begin testString")
	cases := []struct {
		name  string
		s     string
		valid htmlctrl.Validator
	}{
		{"s1", "abc", nil},
		{"s2", "", htmlctrl.ValidateString(func(s string) bool {
			if s == "hello" {
				log("s2 can't be 'hello'")
			}
			return s != "hello"
		})},
	}
	strings := jq("<div>").AddClass("strings")
	for _, c := range cases {
		logInfo(fmt.Sprintf("test case: %#v", c))
		j, e := htmlctrl.String(&c.s, c.name, c.valid)
		if e != nil {
			logError(fmt.Sprintf("%s: unexpected error: %s", c.name, e))
		}
		if title := j.Attr("title"); title != c.name {
			logError(fmt.Sprintf("%s: title is %s, expected %s", c.name, title, c.name))
		}
		strings.Append(j)
		c := &c
		strings.Append(jq("<button>").SetText("verify "+c.name).Call(jquery.CLICK, func() {
			log(c.name, c.s)
		}))
	}
	body.Append(strings)
	logInfo("end testString")
}

func testChoice(body jquery.JQuery) {
	logInfo("begin testChoice")
	opts := []string{
		"def",
		"abc",
		"invalid",
		"hi",
	}
	cases := []struct {
		name  string
		s     string
		valid htmlctrl.Validator
	}{
		{"c1", "abc", nil},
		{"c2", "", htmlctrl.ValidateString(func(c string) bool {
			if c == "invalid" {
				log("c2 can't be 'invalid'")
			}
			return c != "invalid"
		})},
	}
	choices := jq("<div>").AddClass("choices")
	for _, c := range cases {
		logInfo(fmt.Sprintf("test case: %#v", c))
		j, e := htmlctrl.Choice(&c.s, opts, c.name, c.valid)
		if e != nil {
			logError(fmt.Sprintf("%s: unexpected error: %s", c.name, e))
		}
		if title := j.Attr("title"); title != c.name {
			logError(fmt.Sprintf("%s: title is %s, expected %s", c.name, title, c.name))
		}
		choices.Append(j)
		c := &c
		choices.Append(jq("<button>").SetText("verify "+c.name).Call(jquery.CLICK, func() {
			log(c.name, c.s)
		}))
	}
	body.Append(choices)
	logInfo("end testChoice")
}

type sliceCase interface {
	name() string
	slice() interface{}
	mms() (min, max, step float64)
	valid() htmlctrl.Validator
}

type sliceBoolCase struct {
	n string
	s []bool
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

type sliceBoolPtrCase struct {
	n string
	s []*bool
	v htmlctrl.Validator
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

type sliceIntCase struct {
	n              string
	s              []int
	min, max, step int
	v              htmlctrl.Validator
}

func (s *sliceIntCase) name() string {
	return s.n
}

func (s *sliceIntCase) slice() interface{} {
	return interface{}(&s.s)
}

func (s *sliceIntCase) mms() (min, max, step float64) {
	return float64(s.min), float64(s.max), float64(s.step)
}

func (s *sliceIntCase) valid() htmlctrl.Validator {
	return s.v
}

type sliceIntPtrCase struct {
	n              string
	s              []*int
	min, max, step int
	v              htmlctrl.Validator
}

func (s *sliceIntPtrCase) name() string {
	return s.n
}

func (s *sliceIntPtrCase) slice() interface{} {
	return interface{}(&s.s)
}

func (s *sliceIntPtrCase) mms() (min, max, step float64) {
	return float64(s.min), float64(s.max), float64(s.step)
}

func (s *sliceIntPtrCase) valid() htmlctrl.Validator {
	return s.v
}

type sliceFloat64Case struct {
	n              string
	s              []float64
	min, max, step float64
	v              htmlctrl.Validator
}

func (s *sliceFloat64Case) name() string {
	return s.n
}

func (s *sliceFloat64Case) slice() interface{} {
	return interface{}(&s.s)
}

func (s *sliceFloat64Case) mms() (min, max, step float64) {
	return float64(s.min), float64(s.max), float64(s.step)
}

func (s *sliceFloat64Case) valid() htmlctrl.Validator {
	return s.v
}

type sliceFloat64PtrCase struct {
	n              string
	s              []*float64
	min, max, step int
	v              htmlctrl.Validator
}

func (s *sliceFloat64PtrCase) name() string {
	return s.n
}

func (s *sliceFloat64PtrCase) slice() interface{} {
	return interface{}(&s.s)
}

func (s *sliceFloat64PtrCase) mms() (min, max, step float64) {
	return float64(s.min), float64(s.max), float64(s.step)
}

func (s *sliceFloat64PtrCase) valid() htmlctrl.Validator {
	return s.v
}

type sliceStringCase struct {
	n string
	s []string
	v htmlctrl.Validator
}

func (s *sliceStringCase) name() string {
	return s.n
}

func (s *sliceStringCase) slice() interface{} {
	return interface{}(&s.s)
}

func (s *sliceStringCase) mms() (min, max, step float64) {
	return 0, 0, 0
}

func (s *sliceStringCase) valid() htmlctrl.Validator {
	return s.v
}

type sliceStringPtrCase struct {
	n string
	s []*string
	v htmlctrl.Validator
}

func (s *sliceStringPtrCase) name() string {
	return s.n
}

func (s *sliceStringPtrCase) slice() interface{} {
	return interface{}(&s.s)
}

func (s *sliceStringPtrCase) mms() (min, max, step float64) {
	return 0, 0, 0
}

func (s *sliceStringPtrCase) valid() htmlctrl.Validator {
	return s.v
}

type sliceIntSliceCase struct {
	n              string
	s              [][]int
	min, max, step int
	v              htmlctrl.Validator
}

func (s *sliceIntSliceCase) name() string {
	return s.n
}

func (s *sliceIntSliceCase) slice() interface{} {
	return interface{}(&s.s)
}

func (s *sliceIntSliceCase) mms() (min, max, step float64) {
	return float64(s.min), float64(s.max), float64(s.step)
}

func (s *sliceIntSliceCase) valid() htmlctrl.Validator {
	return s.v
}

type sliceIntPtrSliceCase struct {
	n              string
	s              []*[]*int
	min, max, step int
	v              htmlctrl.Validator
}

func (s *sliceIntPtrSliceCase) name() string {
	return s.n
}

func (s *sliceIntPtrSliceCase) slice() interface{} {
	return interface{}(&s.s)
}

func (s *sliceIntPtrSliceCase) mms() (min, max, step float64) {
	return float64(s.min), float64(s.max), float64(s.step)
}

func (s *sliceIntPtrSliceCase) valid() htmlctrl.Validator {
	return s.v
}

func testSlices(body jquery.JQuery) {
	logInfo("begin testSlices")
	logInfo("begin testSlice bool")
	cases := []sliceCase{
		&sliceBoolCase{"bool1", []bool{}},
		&sliceBoolCase{"bool2", []bool{true, false}},
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
		&sliceBoolPtrCase{"[]*bool1", []*bool{&b1, &b2}, htmlctrl.ValidateBool(func(b bool) bool {
			log("bool is locked at true")
			return b
		})},
		&sliceBoolPtrCase{"[]*bool2", []*bool{}, nil},
	}
	testSlice(body, cases)

	logInfo("begin testSlice int")
	cases = []sliceCase{
		&sliceIntCase{"[]int1", []int{2, 4}, 0, 50, 2, htmlctrl.ValidateInt(func(i int) bool {
			allowed := i != 3 && i != 5 && i != 7
			if !allowed {
				log("int may not be 3, 5, or 7")
			}
			return allowed
		})},
		&sliceIntCase{"[]int2", []int{}, 0, 0, 1, nil},
	}
	testSlice(body, cases)

	logInfo("begin testSlice *int")
	i1, i2 := 1, 22
	cases = []sliceCase{
		&sliceIntPtrCase{"[]*int1", []*int{&i1, &i2}, 0, 50, 2, htmlctrl.ValidateInt(func(i int) bool {
			allowed := i != 3 && i != 5 && i != 7
			if !allowed {
				log("int may not be 3, 5, or 7")
			}
			return allowed
		})},
		&sliceIntPtrCase{"[]*int2", []*int{}, 0, 0, 1, nil},
	}
	testSlice(body, cases)

	logInfo("begin testSlice float64")
	cases = []sliceCase{
		&sliceFloat64Case{"[]float64 1", []float64{2.1, 4.2}, 0, 50, 2.1,
			htmlctrl.ValidateFloat64(func(f float64) bool {
				allowed := f != 3 && f != 5 && f != 7
				if !allowed {
					log("float64 may not be 3, 5, or 7")
				}
				return allowed
			})},
		&sliceFloat64Case{"[]float64 2", []float64{}, 0, 0, 1, nil},
	}
	testSlice(body, cases)

	logInfo("begin testSlice *float64")
	f1, f2 := 1.1, 22.2
	cases = []sliceCase{
		&sliceFloat64PtrCase{"[]*float64 1", []*float64{&f1, &f2}, 0, 50, 2,
			htmlctrl.ValidateFloat64(func(f float64) bool {
				allowed := f != 3 && f != 5 && f != 7
				if !allowed {
					log("float64 may not be 3, 5, or 7")
				}
				return allowed
			})},
		&sliceFloat64PtrCase{"[]*float64 2", []*float64{}, 0, 0, 1, nil},
	}
	testSlice(body, cases)

	logInfo("begin testSlice string")
	cases = []sliceCase{
		&sliceStringCase{"[]string1", []string{"a", "b"},
			htmlctrl.ValidateString(func(s string) bool {
				allowed := s != "c" && s != "d"
				if !allowed {
					log("string may not be c, d")
				}
				return allowed
			})},
		&sliceStringCase{"[]string2", []string{}, nil},
	}
	testSlice(body, cases)

	logInfo("begin testSlice *string")
	s1, s2 := "ab", "cd"
	cases = []sliceCase{
		&sliceStringPtrCase{"[]*string1", []*string{&s1, &s2},
			htmlctrl.ValidateString(func(s string) bool {
				allowed := s != "c" && s != "d"
				if !allowed {
					log("string may not be c, d")
				}
				return allowed
			})},
		&sliceStringPtrCase{"[]*string2", []*string{}, nil},
	}
	testSlice(body, cases)

	logInfo("begin testSlice []int")
	cases = []sliceCase{
		&sliceIntSliceCase{"[][]int1", [][]int{{2, 4}, {8, 16}}, 0, 50, 2, htmlctrl.ValidateInt(func(i int) bool {
			allowed := i != 3 && i != 5 && i != 7
			if !allowed {
				log("int may not be 3, 5, or 7")
			}
			return allowed
		})},
		&sliceIntSliceCase{"[][]int2", [][]int{}, 0, 0, 1, nil},
	}
	testSlice(body, cases)

	logInfo("begin testSlice *[]*int")
	is1, is2 := []*int{&i1, &i2}, []*int{}
	cases = []sliceCase{
		&sliceIntPtrSliceCase{"[]*[]*int1", []*[]*int{&is1, &is2}, 0, 50, 2, htmlctrl.ValidateInt(func(i int) bool {
			allowed := i != 3 && i != 5 && i != 7
			if !allowed {
				log("int may not be 3, 5, or 7")
			}
			return allowed
		})},
		&sliceIntPtrSliceCase{"[]*[]*int2", []*[]*int{}, 0, 0, 1, nil},
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
			logError(fmt.Sprintf("%s: unexpected error: %s", c.name(), e))
		}
		if title := j.Attr("title"); title != c.name() {
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
	Bptr := true
	Iptr := 11
	Fptr := 1.1
	Sptr := "abc"
	struct1 := struct {
		b    bool
		B    bool     `desc:"a bool"`
		Bptr *bool    `desc:"bool ptr"`
		Bt   bool     `desc:"Always true" valid:"BoolTrue"`
		I    int      `desc:"an int"`
		Iptr *int     `desc:"int ptr"`
		Ilim int      `desc:"limited int" min:"1" max:"10" step:"2" valid:"IntNot5"`
		F    float64  `desc:"an float64"`
		Fptr *float64 `desc:"float64 ptr"`
		Flim float64  `desc:"limited float64" min:"1.2" max:"10.5" step:"1.2" valid:"Float64Not5"`
		S    string   `desc:"a string"`
		Sptr *string  `desc:"string ptr"`
		Slim string   `desc:"limited string" valid:"StringNotHello"`
		C    string   `desc:"a choice" choice:"def,abc,invalid,hi"`
		Cptr *string  `desc:"choice ptr" choice:"def,abc,invalid,hi"`
		Clim string   `desc:"limited choice" choice:"def,abc,invalid,hi" valid:"ChoiceNotInvalid"`
	}{
		false, false, &Bptr, true,
		2, &Iptr, 1,
		2.5, &Fptr, 1.2,
		"a", &Sptr, "def",
		"", &Sptr, "hi",
	}
	htmlctrl.RegisterValidator("BoolTrue", htmlctrl.ValidateBool(func(b bool) bool {
		log("bool is locked at true")
		return b
	}))
	htmlctrl.RegisterValidator("IntNot5", htmlctrl.ValidateInt(func(i int) bool {
		not5 := i != 5
		if !not5 {
			log("int can't be 5")
		}
		return not5
	}))
	htmlctrl.RegisterValidator("Float64Not5", htmlctrl.ValidateFloat64(func(f float64) bool {
		not5 := f != 5
		if !not5 {
			log("float can't be 5")
		}
		return not5
	}))
	htmlctrl.RegisterValidator("StringNotHello", htmlctrl.ValidateString(func(s string) bool {
		notHello := s != "hello"
		if !notHello {
			log("string can't be 'hello'")
		}
		return notHello
	}))
	htmlctrl.RegisterValidator("ChoiceNotInvalid", htmlctrl.ValidateString(func(c string) bool {
		if c == "invalid" {
			log("choice can't be 'invalid'")
		}
		return c != "invalid"
	}))
	_, e := htmlctrl.Struct(struct1, "error")
	if e == nil {
		logError("expected error when passing non-ptr")
	}
	_, e = htmlctrl.Struct(&e, "error")
	if e == nil {
		logError("expected error when passing ptr to non-slice")
	}

	j, e := htmlctrl.Struct(&struct1, "struct1")
	if e != nil {
		logError(fmt.Sprintf("%s: unexpected error: %s", "struct1", e))
	}
	if title := j.Attr("title"); title != "struct1" {
		logError(fmt.Sprintf("%s: title is %s, expected %s", "struct1", title, "struct1"))
	}
	body.Append(j)
	body.Append(jq("<button>").SetText("verify struct1").Call(jquery.CLICK, func() {
		log("struct1", struct1)
	}))

	logInfo("end testStruct")
}
