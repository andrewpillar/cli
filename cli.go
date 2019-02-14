package cli

import (
	"errors"
	"strings"
)

type args []string

type Cli struct {
	cmds  commands
	flags flags
	main  *Command
}

func addCommand(name string, parent *Command, handler commandHandler, cmds commands) *Command {
	cmd := &Command{
		parent:  parent,
		cmds:    commands(make(map[string]*Command)),
		handler: handler,
		Name:    name,
		Args:    args([]string{}),
		Flags:   newFlags(),
	}

	cmds[name] = cmd

	return cmd
}

func findCommand(argv []string, cmds commands, main *Command) (*Command, error) {
	if len(argv) == 0 {
		if main == nil || main.handler == nil {
			return nil, errors.New("no command to run")
		}

		return main, nil
	}

	name := argv[0]

	if strings.HasPrefix(name, "--") || strings.HasPrefix(name, "-") {
		main.Args = argv

		return main, nil
	}

	cmd, ok := cmds[name]

	if !ok {
		if main == nil {
			return nil, errors.New("command '" + name + "' not found")
		}

		main.Args = argv

		return main, nil
	}

	cmd.Args = argv[1:]

	return findCommand(cmd.Args, cmd.cmds, cmd)
}

func (a args) Get(i int) string {
	if i >= len(a) {
		return ""
	}

	return a[i]
}

func (a *args) set(i int, s string) {
	if i >= len(*a) {
		return
	}

	(*a)[i] = s
}

func New() *Cli {
	return &Cli{cmds: commands(make(map[string]*Command)), flags: newFlags()}
}

func (c *Cli) parseLong(i int, arg string, cmd *Command, flag *Flag) error {
	cmd.Args.set(i, "")

	if flag.Argument {
		val := ""

		if strings.Contains(arg, "=") {
			val = strings.Split(arg, "=")[1]
		} else {
			val = cmd.Args.Get(i + 1)

			if !strings.HasPrefix(val, "--") && !strings.HasPrefix(val, "-") {
				val = cmd.Args.Get(i + 1)

				cmd.Args.set(i + 1, "")
			} else {
				val = ""
			}
		}

		if val == "" && flag.Default == nil {
			return errors.New("option '" + arg + "' requires an argument")
		}

		flag.Value = val
		return nil
	}

	flag.isSet = true
	return nil
}

func (c *Cli) parseShort(i int, arg string, cmd *Command, flag *Flag) error {
	cmd.Args.set(i, "")

	if flag.Argument {
		val := cmd.Args.Get(i + 1)

		cmd.Args.set(i + 1, "")

		if val == "" && flag.Default == nil {
			return errors.New("option '" + arg + "' requires an argument")
		}

		flag.Value = val
		return nil
	}

	flag.isSet = true
	return nil
}

func (c *Cli) AddFlag(f *Flag) {
	f.global = true

	c.flags.expected[f.Name] = f
}

func (c *Cli) Command(name string, handler commandHandler) *Command {
	return addCommand(name, nil, handler, c.cmds)
}

func (c *Cli) MainCommand(handler commandHandler) *Command {
	c.main = &Command{
		cmds:    commands(make(map[string]*Command)),
		handler: handler,
		Args:    args([]string{}),
		Flags:   newFlags(),
	}

	return c.main
}

func (c *Cli) Run(argv []string) error {
	cmd, err := findCommand(args(argv), c.cmds, c.main)

	if err != nil {
		return err
	}

	for _, flag := range c.flags.expected {
		cmd.AddFlag(flag)
	}

	for i, arg := range cmd.Args {
		if arg == "--" {
			break
		}

		if strings.HasPrefix(arg, "--") || strings.HasPrefix(arg, "-") {
			var flag *Flag

			for _, f := range cmd.Flags.expected {
				if f.matches(arg) {
					flag = f
					break
				}
			}

			if flag == nil {
				return errors.New("unknown option '" + arg + "'")
			}

			if strings.HasPrefix(arg, "--") {
				if err := c.parseLong(i, arg, cmd, flag); err != nil {
					return err
				}
			} else if err := c.parseShort(i, arg, cmd, flag); err != nil {
				return err
			}

			cmd.Flags.received[flag.Name] = append(cmd.Flags.received[flag.Name], *flag)
		}
	}

	trimmed := make([]string, 0, len(argv))

	for _, a := range cmd.Args {
		if a != "" {
			trimmed = append(trimmed, a)
		}
	}

	cmd.Args = trimmed
	cmd.Run()

	return nil
}
