// Package htmlctrl works with gopherjs to turn Go values into html. Changes made to the html will be automatically
// reflected in the value used to created it. After conversion you should refrain from modifying the value.
package htmlctrl

import (
	"fmt"
	"math"
	"reflect"
	"strconv"

	"github.com/gopherjs/jquery"
)

// ClassPrefix is used to prefix the CSS classes. They will be of the form ClassPrefix-GoType
var ClassPrefix = "go"

var (
	// SliceAddText is used to fill the add button for a slice
	SliceAddText = "+"
	// SliceDelText is used to fill the delete button for a slice
	SliceDelText = "-"
)

var jq = jquery.NewJQuery

var validators = make(map[string]Validator)

// Validator is used to validate changes made via html objects. The Valid function is given the requested new value
// and should return true only when it is an acceptable value. If it returns false then the change is reverted
type Validator interface {
	Validate(interface{}) bool
}

// ValidatorFunc describes an abitrary function that implements the Validator interface.
type ValidatorFunc func(interface{}) bool

// Validate implements the Validator interface
func (v ValidatorFunc) Validate(i interface{}) bool {
	return v(i)
}

// ValidateBool is a function that validates bool types.
type ValidateBool func(bool) bool

// Validate implements the Validator interface but type asserts that the argument is a bool.
func (v ValidateBool) Validate(i interface{}) bool {
	return v(i.(bool))
}

// ValidateInt is a function that validates int types.
type ValidateInt func(int) bool

// Validate implements the Validator interface but type asserts that the argument is an int.
func (v ValidateInt) Validate(i interface{}) bool {
	return v(i.(int))
}

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
	t, v := reflect.TypeOf(structPtr), reflect.ValueOf(structPtr)
	if t.Kind() != reflect.Ptr {
		return jq(), fmt.Errorf("structPtr should be a pointer, got %s instead", t.Kind())
	}
	if t.Elem().Kind() != reflect.Struct {
		return jq(), fmt.Errorf("structPtr should be a pointer to struct, got pointer to %s instead", t.Elem().Kind())
	}
	structType, structValue := t.Elem(), v.Elem()

	j := jq("<div>").AddClass(ClassPrefix + "-struct")
	j.SetAttr("title", desc)
	for i := 0; i < structType.NumField(); i++ {
		fieldType := structType.Field(i)
		// Ignore unexported fields
		if fieldType.PkgPath != "" {
			continue
		}
		fieldValue := structValue.Field(i)
		tag := fieldType.Tag
		validName := tag.Get("valid")
		valid, ok := validators[validName]
		if validName != "" && !ok {
			return jq(), fmt.Errorf("unregistered validator '%s'", validName)
		}
		min, e := strconv.ParseFloat(tag.Get("min"), 64)
		if e != nil {
			if tag.Get("min") != "" {
				return jq(), fmt.Errorf("min as value '%s' expected a number", tag.Get("min"))
			}
			min = math.NaN()
		}
		max, e := strconv.ParseFloat(tag.Get("max"), 64)
		if e != nil {
			if tag.Get("max") != "" {
				return jq(), fmt.Errorf("max as value '%s' expected a number", tag.Get("max"))
			}
			max = math.NaN()
		}
		step, e := strconv.ParseFloat(tag.Get("step"), 64)
		if e != nil {
			if tag.Get("step") != "" {
				return jq(), fmt.Errorf("step as value '%s' expected a number", tag.Get("step"))
			}
			step = math.NaN()
		}

		field, e := convert(fieldValue, tag.Get("desc"), min, max, step, valid)
		if e != nil {
			return jq(), fmt.Errorf("converting struct field %s (%s): %s", fieldType.Name, fieldType.Type.Kind(), e)
		}
		jf := jq("<div>").AddClass(ClassPrefix + "-struct-field")
		jf.Append(jq("<label>").SetText(fieldType.Name))
		jf.Append(field)
		j.Append(jf)
		fmt.Println(fieldType, fieldValue, tag)
	}
	return j, nil
}

// Slice takes a pointer to a slice and returns a JQuery object associated with it as a list tag. A non-nil error
// is returned in the event the conversion fails. It includes buttons for adding and removing elements from the
// slice. The slice's type must be among those supported by this package or a pointer to one. An error will be
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
	sliceElemType := sliceType.Elem()

	j := jq("<list>").AddClass(ClassPrefix + "-slice")
	j.SetAttr("title", desc)

	var populate func() error
	populate = func() error {
		newLi := func(j, ji jquery.JQuery) jquery.JQuery {
			li := jq("<li>").Append(ji)
			delBtn := jq("<button>").SetText(SliceDelText)
			delBtn.Call(jquery.CLICK, func() {
				i := li.Call("index").Get().Int()
				li.Remove()
				begin := sliceValue.Slice(0, i)
				end := sliceValue.Slice(i+1, sliceValue.Len())
				sliceValue.Set(reflect.AppendSlice(begin, end))
				// Just delete and redo everything to work with non-pointers when the slice resizes
				j.Empty()
				e := populate()
				if e != nil {
					panic(e)
				}
			})
			li.Append(delBtn)
			return li
		}

		for i := 0; i < sliceValue.Len(); i++ {
			elem := sliceValue.Index(i)
			ji, e := convert(elem, "", min, max, step, valid)
			if e != nil {
				return fmt.Errorf("converting slice element %d (%s): %s", i, elem.Type().Kind(), e)
			}
			j.Append(newLi(j, ji))
		}
		addBtn := jq("<button>").SetText(SliceAddText)
		addBtn.Call(jquery.CLICK, func() {
			if sliceElemType.Kind() == reflect.Ptr {
				newElem := reflect.New(sliceElemType.Elem())
				sliceValue.Set(reflect.Append(sliceValue, newElem))
			} else {
				newElem := reflect.New(sliceElemType)
				sliceValue.Set(reflect.Append(sliceValue, newElem.Elem()))
			}
			// Just delete and redo everything to work with non-pointers when the slice resizes
			j.Empty()
			e := populate()
			if e != nil {
				panic(e)
			}
		})
		j.Append(addBtn)
		return nil
	}

	e := populate()
	if e != nil {
		return jq(), e
	}

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
		if valid != nil && !valid.Validate(bNew) {
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
//
// min, max, and step are float64 to allow the use of math.NaN() to indicate not to set the corresponding html
// attribute. They will be truncated to ints otherwise.
func Int(i *int, desc string, min, max, step float64, valid Validator) (jquery.JQuery, error) {
	j := jq("<input>").AddClass(ClassPrefix + "-int")
	j.SetAttr("title", desc)
	j.SetAttr("type", "number")
	if !math.IsNaN(min) {
		j.SetAttr("min", int(min))
	}
	if !math.IsNaN(max) {
		j.SetAttr("max", int(max))
	}
	if !math.IsNaN(step) {
		j.SetAttr("step", int(step))
	}
	j.SetAttr("value", *i)
	j.SetData("prev", *i)
	j.Call(jquery.CHANGE, func(event jquery.Event) {
		val := event.Target.Get("value").String()
		newI, e := strconv.Atoi(val)
		if e != nil {
			f, e := strconv.ParseFloat(val, 64)
			if e != nil {
				panic(fmt.Errorf("value '%s' has invalid type, expected a number", val))
			}
			// Truncate to int
			newI = int(f)
			j.SetVal(newI)
		}
		// Need to check for min and max ourselves because html min and max are easy to get around
		isValid := valid == nil || valid.Validate(newI)
		isToLow := !math.IsNaN(min) && newI < int(min)
		isToHigh := !math.IsNaN(max) && newI > int(max)
		if !isValid || isToLow || isToHigh {
			fmt.Println(isValid, isToLow, isToHigh)
			newI = int(j.Data("prev").(float64))
			j.SetVal(newI)
		}
		*i = newI
		j.SetData("prev", newI)
	})
	return j, nil
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
		return Int(intf.(*int), desc, min, max, step, valid)
	case reflect.Float64:
		return jq(), fmt.Errorf("unimplemented type %s", val.Type().Kind())
	case reflect.String:
		return jq(), fmt.Errorf("unimplemented type %s", val.Type().Kind())
	}
	return jq(), fmt.Errorf("unsupported type %s", val.Type().Kind())
}
