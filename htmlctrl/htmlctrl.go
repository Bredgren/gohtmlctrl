// Package htmlctrl works with gopherjs to turn Go values into html. Changes made to the html will be automatically
// reflected in the value used to created it. After conversion you should refrain from modifying the value.
package htmlctrl

import (
	"fmt"
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
// The current value
//
// Struct tags recognized:
//  * desc - A description of the struct field. Becomes the title attribute of the html tag.
//  * min - Minimum value for a number
//  * max - Maximum value for a number
//  * choice - Comma separated list. This will created an html choice tag when used on a string type.
//  * valid - Name of a registered validator.
func Struct(structPtr interface{}, desc string) (jquery.JQuery, error) {
	return jq(), nil
}

// Slice takes a pointer to a slice and returns a JQuery object associated with it. A non-nil error is returned
// in the event the conversion fails. It includes buttons for adding and removing elements from the slice.
func Slice(slicePtr interface{}, desc string) (jquery.JQuery, error) {
	return jq(), nil
}

// Bool takes a pointer to a bool value and returns a JQuery object associated with it in the form of a checkbox.
// A non-nil error is returned in the event the conversion fails. The current value of the bool will be used as
// the initial value of the checkbox.
func Bool(b *bool, desc string, valid Validator) (jquery.JQuery, error) {
	j := jq("<input>").AddClass(ClassPrefix + "-bool")
	j.SetAttr("type", "checkbox")
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
