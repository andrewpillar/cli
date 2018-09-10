package cli

import (
	"strconv"
	"strings"
)

type flags map[string]*Flag

type flagHandler func(f Flag, c Command)

type Flag struct {
	isSet bool

	global bool

	Name string

	Short string

	Long string

	Argument bool

	Value string

	Default interface{}

	Exclusive bool

	Handler flagHandler
}

func newFlags() flags {
	return flags(make(map[string]*Flag))
}

func (f flags) GetInt64(name string) (int64, error) {
	i, err := f.GetInt(name)

	if err != nil {
		return 0, err
	}

	return int64(i), nil
}

func (f flags) GetInt32(name string) (int32, error) {
	i, err := f.GetInt64(name)

	if err != nil {
		return 0, err
	}

	return int32(i), nil
}

func (f flags) GetInt(name string) (int, error) {
	flag, ok := f[name]

	if !ok {
		return 0, nil
	}

	if flag.Value == "" {
		if flag.Default == nil {
			return 0, nil
		}

		return flag.Default.(int), nil
	}

	i, err := strconv.ParseInt(flag.Value, 10, 64)

	if err != nil {
		return 0, err
	}

	return int(i), nil
}

func (f flags) GetSlice(name, sep string) []string {
	return strings.Split(f.GetString(name), sep)
}

func (f flags) GetString(name string) string {
	flag, ok := f[name]

	if !ok {
		return ""
	}

	if flag.Value == "" {
		if flag.Default == nil {
			return ""
		}

		return flag.Default.(string)
	}

	return flag.Value
}

func (f flags) IsSet(name string) bool {
	flag, ok := f[name]

	if !ok {
		return false
	}

	return flag.isSet
}

func (f Flag) Matches(arg string) bool {
	if strings.Contains(arg, "=") {
		arg = strings.Split(arg, "=")[0]
	}

	return f.Short == arg || f.Long == arg
}
