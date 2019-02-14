package cli

import (
	"strconv"
	"strings"
)

type flags struct {
	expected map[string]*Flag
	received map[string][]Flag
}

type flagHandler func(f Flag, c Command)

type Flag struct {
	handler   flagHandler
	isSet     bool
	global    bool
	Name      string
	Short     string
	Long      string
	Argument  bool
	Value     string
	Default   interface{}
	Exclusive bool
}

func newFlags() flags {
	return flags{
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
	if f.Value == "" {
		if f.Default == nil {
			return 0.0, nil
		}

		return f.Default.(float64), nil
	}

	fl, err := strconv.ParseFloat(f.Value, bitSize)

	if err != nil {
		return 0.0, err
	}

	return fl, nil
}

func (f Flag) getInt(bitSize int) (int64, error) {
	if f.Value == "" {
		if f.Default == nil {
			return 0, nil
		}

		return f.Default.(int64), nil
	}

	i, err := strconv.ParseInt(f.Value, 10, bitSize)

	if err != nil {
		return 0, err
	}

	return i, nil
}

func (f Flag) GetFloat32() (float32, error) {
	fl, err := f.getFloat(32)

	if err != nil {
		return 0.0, err
	}

	return float32(fl), nil
}

func (f Flag) GetFloat64() (float64, error) {
	fl, err := f.getFloat(64)

	if err != nil {
		return 0.0, err
	}

	return fl, nil
}

func (f Flag) GetInt() (int, error) {
	i, err := f.getInt(0)

	if err != nil {
		return 0, err
	}

	return int(i), err
}

func (f Flag) GetInt8() (int8, error) {
	i, err := f.getInt(8)

	if err != nil {
		return 0, err
	}

	return int8(i), err
}

func (f Flag) GetInt16() (int16, error) {
	i, err := f.getInt(16)

	if err != nil {
		return 0, err
	}

	return int16(i), err
}

func (f Flag) GetInt32() (int32, error) {
	i, err := f.getInt(32)

	if err != nil {
		return 0, err
	}

	return int32(i), err
}

func (f Flag) GetInt64() (int64, error) {
	i, err := f.getInt(64)

	if err != nil {
		return 0, err
	}

	return i, err
}

func (f Flag) GetString() string {
	if f.Value == "" && f.Default != nil {
		return f.Default.(string)
	}

	return f.Value
}

func (f Flag) IsSet() bool {
	return f.isSet
}

func (f flags) first(name string) Flag {
	flags := f.GetAll(name)

	if len(flags) > 0 {
		return flags[0]
	}

	return *f.expected[name]
}

func (f flags) GetAll(name string) []Flag {
	return f.received[name]
}

func (f flags) GetInt(name string) (int, error) {
	return f.first(name).GetInt()
}

func (f flags) GetInt8(name string) (int8, error) {
	return f.first(name).GetInt8()
}

func (f flags) GetInt16(name string) (int16, error) {
	return f.first(name).GetInt16()
}

func (f flags) GetInt32(name string) (int32, error) {
	return f.first(name).GetInt32()
}

func (f flags) GetInt64(name string) (int64, error) {
	return f.first(name).GetInt64()
}

func (f flags) GetFloat32(name string) (float32, error) {
	return f.first(name).GetFloat32()
}

func (f flags) GetFloat64(name string) (float64, error) {
	return f.first(name).GetFloat64()
}

func (f flags) GetString(name string) string {
	return f.first(name).GetString()
}

func (f flags) IsSet(name string) bool {
	flags := f.GetAll(name)

	if len(flags) == 0 {
		return false
	}

	return flags[0].isSet
}
