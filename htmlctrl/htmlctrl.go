// Package htmlctrl works with gopherjs to turn Go values into html. Changes made to the html will be automatically
// reflected in the value used to created it. After conversion you should refrain from modifying the value.
package htmlctrl

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/gopherjs/jquery"
)

// ClassPrefix is used to prefix the CSS classes. They will be of the form ClassPrefix-GoType
var ClassPrefix = "go"

var jq = jquery.NewJQuery

var validators = make(map[string]Validator)

// Validator is a function used to validate changes made via html objects. It is given the requested new value
// and should return true only when it is an acceptable value. If it returns false then the change is reverted
type Validator func(interface{}) bool

// RegisterValidator associates a name with the validator function so that it may be referenced in a struct tag.
func RegisterValidator(name string, fn Validator) {
	validators[name] = fn
}

// Struct takes a pointer to a struct and returns a JQuery object associated with it. A non-nil error is returned
// in the event the conversion fails.
//
// All exported fields of the struct will recursively converted. Fields that whose types don't support conversion
// are ignored. A type is supported if it has it's own conversion function in this package.
//
// Struct tags recognized
//  desc - A description of the struct field. Becomes the title attribute of the html tag.
//  min - Minimum value for a number
//  max - Maximum value for a number
//  step - How much the up and down buttons change a number by
//  choice - Comma separated list. This will created an html choice tag when used on a string type.
//  valid - Name of a registered validator.
func Struct(structPtr interface{}, desc string) (jquery.JQuery, error) {
	return jq(), nil
}

// Slice takes a pointer to a slice and returns a JQuery object associated with it as a list tag. A non-nil error
// is returned in the event the conversion fails. It includes buttons for adding and removing elements from the
// slice. The slice's type must be among those supported by this package (or a pointer to one). An error will be
// returned if the slice's type is not supported.
//
// min, max, step, and valid will be applied if the slices element type supports it.
func Slice(slicePtr interface{}, desc string, min, max, step float64, valid Validator) (jquery.JQuery, error) {
	t, v := reflect.TypeOf(slicePtr), reflect.ValueOf(slicePtr)
	if t.Kind() != reflect.Ptr {
		return jq(), fmt.Errorf("slicePtr should be a pointer, got %s instead", t.Kind())
	}
	if t.Elem().Kind() != reflect.Slice {
		return jq(), fmt.Errorf("slicePtr should be a pointer to slice, got pointer to %s instead", t.Elem().Kind())
	}
	sliceType, sliceValue := t.Elem(), v.Elem()
	_, _ = sliceType, sliceValue

	newLi := func(ji jquery.JQuery) jquery.JQuery {
		li := jq("<li>").Append(ji)
		delBtn := jq("<button>").SetText("-")
		delBtn.Call(jquery.CLICK, func() {
			li.Remove()
		})
		li.Append(delBtn)
		return li
	}
	j := jq("<list>").AddClass(ClassPrefix + "-slice")
	j.SetAttr("title", desc)
	for i := 0; i < sliceValue.Len(); i++ {
		elem := sliceValue.Index(i)
		ji, e := convert(elem, "", min, max, step, valid)
		if e != nil {
			return jq(), nil
		}
		j.Append(newLi(ji))
	}
	addBtn := jq("<button>").SetText("+")
	addBtn.Call(jquery.CLICK, func() {
		fmt.Println("add")
		// j.Append(newLi(ji))
	})
	j.Append(addBtn)

	return j, nil
}

// Bool takes a pointer to a bool value and returns a JQuery object associated with it in the form of a checkbox.
// A non-nil error is returned in the event the conversion fails. The current value of the bool will be used as
// the initial value of the checkbox.
func Bool(b *bool, desc string, valid Validator) (jquery.JQuery, error) {
	j := jq("<input>").AddClass(ClassPrefix + "-bool")
	j.SetAttr("type", "checkbox")
	j.SetAttr("title", desc)
	j.SetProp("checked", *b)
	j.SetData("prev", *b)
	j.Call(jquery.CHANGE, func(event jquery.Event) {
		val := event.Target.Get("checked").String()
		bNew, e := strconv.ParseBool(val)
		if e != nil {
			// Theorectially impossible
			panic(fmt.Sprintf("value '%s' has invalid type, expected bool", val))
		}
		if valid != nil && !valid(bNew) {
			bNew = j.Data("prev").(bool)
			j.SetProp("checked", bNew)
		}
		*b = bNew
		j.SetData("prev", bNew)
	})
	return j, nil
}

// Int takes a pointer to a int value and returns a JQuery object associated with it in the form of an input of
// number type. Attempt to fill in a non-int value will result in it being truncated to an integer. A non-nil
// error is returned in the event the conversion fails. The current value of the int will be used as the initial
// value of the input.
func Int(i *int, desc string, min, max int, valid Validator) (jquery.JQuery, error) {
	return jq(), nil
}

// Float64 takes a pointer to a float64 value and returns a JQuery object associated with it in the form of an
// input of number type. A non-nil error is returned in the event the conversion fails. The current value of the
// float64 will be used as the initial value of the input.
func Float64(f *float64, desc string, min, max float64, valid Validator) (jquery.JQuery, error) {
	return jq(), nil
}

// String takes a pointer to a string value and returns a JQuery object associated with it in the form of an
// input of text type. A non-nil error is returned in the event the conversion fails. The
// current value of the string will be used as the initial value of the input.
func String(s *string, desc string, valid Validator) (jquery.JQuery, error) {
	return jq(), nil
}

// Choice is a special string that can only be one of the values in options. It returns a JQuery object
// associated with it in the form of a choice tag. A non-nil error is returned in the event the conversion
// fails. If s is the empty string then the initial value is options[0]. If is is not empty but not in options
// then A non-nil error is returned. If s is in options then it is used as the intial value.
func Choice(s *string, options []string, valid Validator) (jquery.JQuery, error) {
	return jq(), nil
}

func convert(val reflect.Value, desc string, min, max, step float64, valid Validator) (jquery.JQuery, error) {
	kind := val.Type().Kind()
	intf := val.Addr().Interface()
	if val.Type().Kind() == reflect.Ptr {
		kind = val.Type().Elem().Kind()
		intf = val.Interface()
	}
	switch kind {
	case reflect.Struct:
		return jq(), fmt.Errorf("unimplemented type %s", val.Type().Kind())
	case reflect.Slice:
		return jq(), fmt.Errorf("unimplemented type %s", val.Type().Kind())
	case reflect.Bool:
		return Bool(intf.(*bool), desc, valid)
	case reflect.Int:
		return jq(), fmt.Errorf("unimplemented type %s", val.Type().Kind())
	case reflect.Float64:
		return jq(), fmt.Errorf("unimplemented type %s", val.Type().Kind())
	case reflect.String:
		return jq(), fmt.Errorf("unimplemented type %s", val.Type().Kind())
	}
	return jq(), fmt.Errorf("unsupported type %s", val.Type().Kind())
}
