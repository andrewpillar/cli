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

	// Name specifies the name of the flag. This should be a string, and will
	// be what is used to access the flag when passed to the command's handler.
	Name string

	// Short specifies the short version of a flag, for example '-h'.
	Short string

	// Long specifies the long version of a flag, for example '--help'.
	Long string

	// Argument specifies whether or not the flag takes an argument.
	Argument bool

	// If the flag takes an argument then the Value property will be set during
	// parsing of the input arguments.
	Value string

	// Default specifies the default value the flag should be if no value is
	// given to the flag. This is an interface, and before accessing the flag's
	// value you should know what it's expected type should be.
	Default interface{}

	// Exclusive specifies whether or not a flag with a handler should be
	// exclusive in its execution. Setting this to true means that no command
	// will be executed if an exclusive flag, with a handler has been set on
	// that command, and passed to that command.
	//
	// For example, the '--help' flag could be considered an exclusive flag.
	// When passed to a command you do not want the command itself to be
	// executed along with the '--help' flag.
	Exclusive bool

	// Handler specifies the handler for the flag should a flag be given to
	// a command. This handler will be passed the flag itself, and the command
	// on which the flag was passed.
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

// Matches determins if the given argument matches against the current flag,
// based on the short, and long values of the current flag.
func (f Flag) Matches(arg string) bool {
	if strings.Contains(arg, "=") {
		arg = strings.Split(arg, "=")[0]
	}

	return f.Short == arg || f.Long == arg
}
