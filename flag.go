package cli

import (
	"strconv"
	"strings"
)

type Flags struct {
	expected map[string]*Flag
	received map[string][]Flag
}

type flagHandler func(f Flag, c Command)

// A flag for your program, or for the program's individual command.
type Flag struct {
	isSet     bool
	global    bool
	value     string

	// The name of the flag. This is used to retrieve the flag from the flags
	// type that stores either the program's flags, or a command's flags.
	Name string

	// The short version of the flag. Typically a single letter value preprended
	// with a single '-'.
	Short string

	// The long version of the flag. Typically a hyphenated string value
	// prepended with a '--'.
	Long string

	// Whether the flag takes an argument. If set to true, then the flag's
	// argument will be parsed from the input strings. Listed below are the
	// three valid ways of specifying a flag's value:
	//
	//  -f arg
	//  --flag arg
	//  --flag=arg
	Argument  bool

	// The default value of the flag if no value is given to the flag itself
	// during program program invocation.
	Default   interface{}

	// The handler to invoke whenever the flag is set.
	Handler   flagHandler

	// If a flag is exclusive then the flag's handler will be invoked, and the
	// flag's command will not be invoked.
	Exclusive bool
}

func newFlags() Flags {
	return Flags{
		expected: make(map[string]*Flag),
		received: make(map[string][]Flag),
	}
}

func (f Flag) matches(arg string) bool {
	if strings.Contains(arg, "=") {
		arg = strings.Split(arg, "=")[0]
	}

	return f.Short == arg || f.Long == arg
}

func (f Flag) getFloat(bitSize int) (float64, error) {
	if f.value == "" {
		if f.Default == nil {
			return 0.0, nil
		}

		return f.Default.(float64), nil
	}

	fl, err := strconv.ParseFloat(f.value, bitSize)

	if err != nil {
		return 0.0, err
	}

	return fl, nil
}

func (f Flag) getInt(bitSize int) (int64, error) {
	if f.value == "" {
		if f.Default == nil {
			return 0, nil
		}

		return f.Default.(int64), nil
	}

	i, err := strconv.ParseInt(f.value, 10, bitSize)

	if err != nil {
		return 0, err
	}

	return i, nil
}

// Attempt to parse the underlying flag string value to a float32 type.
func (f Flag) GetFloat32() (float32, error) {
	fl, err := f.getFloat(32)

	if err != nil {
		return 0.0, err
	}

	return float32(fl), nil
}

// Attempt to parse the underlying flag string value to a float64 type.
func (f Flag) GetFloat64() (float64, error) {
	fl, err := f.getFloat(64)

	if err != nil {
		return 0.0, err
	}

	return fl, nil
}

// Attempt to parse the underlying flag string value to an int type.
func (f Flag) GetInt() (int, error) {
	i, err := f.getInt(0)

	if err != nil {
		return 0, err
	}

	return int(i), err
}

// Attempt to parse the underlying flag string value to an int8 type.
func (f Flag) GetInt8() (int8, error) {
	i, err := f.getInt(8)

	if err != nil {
		return 0, err
	}

	return int8(i), err
}

// Attempt to parse the underlying flag string value to an int16 type.
func (f Flag) GetInt16() (int16, error) {
	i, err := f.getInt(16)

	if err != nil {
		return 0, err
	}

	return int16(i), err
}

// Attempt to parse the underlying flag string value to an int32 type.
func (f Flag) GetInt32() (int32, error) {
	i, err := f.getInt(32)

	if err != nil {
		return 0, err
	}

	return int32(i), err
}

// Attempt to parse the underlying flag string value to an int64 type.
func (f Flag) GetInt64() (int64, error) {
	i, err := f.getInt(64)

	if err != nil {
		return 0, err
	}

	return i, err
}

// Get the underlying flag value as a string.
func (f Flag) GetString() string {
	if f.value == "" && f.Default != nil {
		return f.Default.(string)
	}

	return f.value
}

// Return whether the flag has been set.
func (f Flag) IsSet() bool {
	return f.isSet
}

func (f Flags) first(name string) Flag {
	flags := f.GetAll(name)

	if len(flags) > 0 {
		return flags[0]
	}

	return *f.expected[name]
}

// Get all of the flags for the given flag name.
func (f Flags) GetAll(name string) []Flag {
	return f.received[name]
}

// Attempt to parse the first flag's value as an int for the given flag name.
func (f Flags) GetInt(name string) (int, error) {
	return f.first(name).GetInt()
}

// Attempt to parse the first flag's value as an int8 for the given flag name.
func (f Flags) GetInt8(name string) (int8, error) {
	return f.first(name).GetInt8()
}

// Attempt to parse the first flag's value as an int16 for the given flag name.
func (f Flags) GetInt16(name string) (int16, error) {
	return f.first(name).GetInt16()
}

// Attempt to parse the first flag's value as an int32 for the given flag name.
func (f Flags) GetInt32(name string) (int32, error) {
	return f.first(name).GetInt32()
}

// Attempt to parse the first flag's value as an int64 for the given flag name.
func (f Flags) GetInt64(name string) (int64, error) {
	return f.first(name).GetInt64()
}

// Attempt to parse the first flag's value as a float32 for the given flag name.
func (f Flags) GetFloat32(name string) (float32, error) {
	return f.first(name).GetFloat32()
}

// Attempt to parse the first flag's value as a float64 for the given flag name.
func (f Flags) GetFloat64(name string) (float64, error) {
	return f.first(name).GetFloat64()
}

// Get the underlying flag value as a string for the given flag name.
func (f Flags) GetString(name string) string {
	return f.first(name).GetString()
}

// Return whether the flag has been set for the given flag name.
func (f Flags) IsSet(name string) bool {
	flags := f.GetAll(name)

	if len(flags) == 0 {
		return false
	}

	return flags[0].isSet
}
