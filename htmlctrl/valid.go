package htmlctrl

var validators = make(map[string]Validator)

// RegisterValidator associates a name with the validator function so that it may be referenced in a struct tag.
func RegisterValidator(name string, fn Validator) {
	validators[name] = fn
}

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

// ValidateFloat64 is a function that validates float64 types.
type ValidateFloat64 func(float64) bool

// Validate implements the Validator interface but type asserts that the argument is an float.
func (v ValidateFloat64) Validate(i interface{}) bool {
	return v(i.(float64))
}

// ValidateString is a function that validates string types.
type ValidateString func(string) bool

// Validate implements the Validator interface but type asserts that the argument is an string.
func (v ValidateString) Validate(i interface{}) bool {
	return v(i.(string))
}
