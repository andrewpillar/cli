package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainCommand(t *testing.T) {
	argv := []string{"argv1", "argv2", "argv3", "argv4"}

	cli := New()

	cli.MainCommand(func(cmd Command) {
		assert.Equal(t, len(cmd.Args), len(argv), "Expected cmd.Args to be %d, it was %d\n", len(argv), len(cmd.Args))
	})

	err := cli.Run(argv)

	assert.Nilf(t, err, "Expected cli.Run to not error, it did: %s\n", err)
}

func TestCommandNilHandler(t *testing.T) {
	cmdName := "cmd-no-handler"

	cli := New()

	cli.Command(cmdName, nil)

	err := cli.Run([]string{cmdName})

	assert.NotNilf(t, err, "Expected cli.Run to error, it did not\n")
}

func TestCommandNilHandlerSubCommand(t *testing.T) {
	cmdName := "cmd-no-handler"
	subCmdName := "sub-cmd-handler"

	argv := []string{cmdName, subCmdName, "argv1"}

	cli := New()

	cli.Command(cmdName, nil).Command(subCmdName, func(cmd Command) {
		assert.Equalf(t, len(cmd.Args), 1, "Expected cmd.Args to be 1, it was not: %s\n", cmd.Args)
	})

	err := cli.Run(argv)

	assert.Nilf(t, err, "Expected cli.Run to not error, it did: %s\n")
}

func TestCommandSubCommand(t *testing.T) {
	cmdName := "cmd-handler"
	subCmdName := "sub-cmd-handler"

	argv := []string{cmdName, subCmdName, "argv1"}

	cli := New()

	cli.Command(cmdName, func(cmd Command) {
		assert.FailNowf(t, "command '%s' should not have executed\n", cmd.Name)
	}).Command(subCmdName, func(cmd Command) {
		assert.Equalf(t, len(cmd.Args), 1, "Expected cmd.Args to be 1, it was not: %s\n", cmd.Args)
	})

	err := cli.Run(argv)

	assert.Nilf(t, err, "Expected cli.Run to not error, it did: %s\n")
}

func TestCommandFullName(t *testing.T) {
	cmdName := "cmd1"
	subCmdName := "cmd2"

	cli := New()

	cmd := cli.Command(cmdName, nil).Command(subCmdName, nil)

	assert.Equalf(t, cmd.FullName(), "cmd1-cmd2", "Expected cmd.Name to be 'cmd1-cmd2', it was not: %s\n", cmd.FullName())
}

func TestCommandNotFound(t *testing.T) {
	cli := New()

	err := cli.Run([]string{"cmd"})

	assert.NotNilf(t, err, "Expected cli.Run to error, it did not\n")
}

func TestGlobalFlag(t *testing.T) {
	argv := []string{"cmd", "sub-cmd", "--help"}

	cli := New()

	cli.AddFlag(&Flag{
		Name: "help",
		Long: "--help",
	})

	cli.Command("cmd", nil).Command("sub-cmd", func(cmd Command) {
		assert.True(t, cmd.Flags.IsSet("help"), "Expected help flag to be set on command '%s', it was not\n", cmd.FullName())
	})

	err := cli.Run(argv)

	assert.Nilf(t, err, "Expected cli.Run to not error, it did: %s\n", err)
}

func TestCommandFlagValue(t *testing.T) {
	argvShort := []string{"cmd", "-f", "value"}
	argvLong := []string{"cmd", "--flag", "value"}
	argvLongEq := []string{"cmd", "--flag=value"}

	cli := New()

	cli.MainCommand(func(cmd Command) {
		assert.Equalf(t, cmd.Flags.GetString("flag"), "value", "Expected flag to have value, it did not\n")
	}).AddFlag(&Flag{
		Name:     "flag",
		Short:    "-f",
		Long:     "--flag",
		Argument: true,
	})

	err := cli.Run(argvShort)

	assert.Nilf(t, err, "Exepcted cli.Run to not error, it did: %s\n", err)

	err = cli.Run(argvLong)

	assert.Nilf(t, err, "Exepcted cli.Run to not error, it did: %s\n", err)

	err = cli.Run(argvLongEq)

	assert.Nilf(t, err, "Exepcted cli.Run to not error, it did: %s\n", err)
}

func TestCommandFlagDefaultValue(t *testing.T) {
	argv := []string{"cmd"}

	cli := New()

	cli.MainCommand(func(cmd Command) {
		assert.Equalf(t, cmd.Flags.GetString("flag"), "value", "Expected flag to have value, it did not\n")
	}).AddFlag(&Flag{
		Name:     "flag",
		Long:     "--flag",
		Argument: true,
		Default:  "value",
	})

	err := cli.Run(argv)

	assert.Nilf(t, err, "Exepcted cli.Run to not error, it did: %s\n", err)
}

func TestCommandMixedFlagValues(t *testing.T) {
	argv := []string{"cmd", "--flag-one", "--flag-two", "value-two"}

	cli := New()

	cli.MainCommand(func(cmd Command) {
		assert.Equal(t, cmd.Flags.GetString("flag-one"), "value-one", "Expected flag-one to have value, it did not\n")
		assert.Equal(t, cmd.Flags.GetString("flag-two"), "value-two", "Expected flag-two to have value, it did not\n")
	}).AddFlag(&Flag{
		Name:     "flag-one",
		Long:     "--flag-one",
		Argument: true,
		Default: "value-one",
	}).AddFlag(&Flag{
		Name:     "flag-two",
		Long:     "--flag-two",
		Argument: true,
	})

	err := cli.Run(argv)

	assert.Nilf(t, err, "Expected cli.Run to not error, it did: %s\n", err)
}

func TestCommandRepeatedFlags(t *testing.T) {
	argv := []string{"cmd", "-f", "val1", "-f", "val2", "-f", "val3"}

	cli := New()

	cli.MainCommand(func(cmd Command) {
		flags := cmd.Flags.GetAll("flag")

		assert.Equalf(t, len(flags), 3, "Expected number of flags to be 3, it was: %d\n", len(flags))

		first := cmd.Flags.GetString("flag")

		assert.Equalf(t, first, "val1", "Expected first flag to be val1, it was: %s\n", first)
	}).AddFlag(&Flag{
		Name:     "flag",
		Short:    "-f",
		Argument: true,
	})

	err := cli.Run(argv)

	assert.Nilf(t, err, "Expected cli.Run to not error, it did: %s\n", err)
}

func TestCommandFlagNumberValues(t *testing.T) {
	argv := []string{"--int", "10", "--float", "25.08"}

	cli := New()

	cli.MainCommand(func(cmd Command) {
		i, err := cmd.Flags.GetInt("int")

		assert.Nilf(t, err, "Expected cmd.Flags.GetInt to not error, it did: %s\n", err)

		fl, err := cmd.Flags.GetFloat64("float")

		assert.Nilf(t, err, "Expected cmd.Flags.GetFloat to not error, it did: %s\n", err)

		assert.Equalf(t, i, 10, "Expected int flag to be 10, it was: %d\n", i)
		assert.Equalf(t, fl, 25.08, "Expected float flag to be 25.08, it was: %f\n", fl)
	}).AddFlag(&Flag{
		Name:     "int",
		Long:     "--int",
		Argument: true,
	}).AddFlag(&Flag{
		Name:     "float",
		Long:     "--float",
		Argument: true,
	})

	err := cli.Run(argv)

	assert.Nilf(t, err, "Expected cli.Run to not error, it did: %s\n", err)
}

func TestCommandFlagHandler(t *testing.T) {
	argv := []string{"--help"}

	cli := New()

	cli.AddFlag(&Flag{
		Name:      "help",
		Long:      "--help",
		Exclusive: true,
		Handler:   func(f Flag, cmd Command) {},
	})

	cli.MainCommand(func(cmd Command) {
		assert.FailNowf(t, "command should not have executed\n", cmd.Name)
	})

	err := cli.Run(argv)

	assert.Nilf(t, err, "Expected cli.Run to not error, it did: %s\n", err)
}
