// Simple library for building CLI applications in Go
package cli

import (
	"strconv"
	"strings"
)

type flagHandler func(f Flag, c Command)

type flags struct {
	expected map[string]*Flag

	received map[string][]Flag
}

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
	return flags{
		expected: make(map[string]*Flag),
		received: make(map[string][]Flag),
	}
}

func (f *flags) putReceived(received Flag) {
	f.received[received.Name] = append(f.received[received.Name], received)
}

func (f flags) GetAll(name string) []Flag {
	return f.received[name]
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
	flag := *f.expected[name]
	flags := f.GetAll(name)

	if len(flags) > 0 {
		flag = flags[0]
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
	flag := *f.expected[name]
	flags := f.GetAll(name)

	if len(flags) > 0 {
		flag = flags[0]
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
	flags := f.GetAll(name)

	if len(flags) == 0 {
		return false
	}

	return flags[0].isSet
}

func (f Flag) Matches(arg string) bool {
	if strings.Contains(arg, "=") {
		arg = strings.Split(arg, "=")[0]
	}

	return f.Short == arg || f.Long == arg
}
