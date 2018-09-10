package cli

import (
	"fmt"
	"strings"
	"testing"
)

func TestMainCommand(t *testing.T) {
	cli := New()

	cli.Main(func(c Command) {
		fmt.Println("main command")
	})

	if err := cli.Run([]string{}); err != nil {
		t.Error(err)
	}
}

func TestCommand(t *testing.T) {
	cli := New()

	cli.Command("hello", func(c Command) {
		fmt.Println("Hello " + c.Args.Get(0))
	})

	if err := cli.Run([]string{"hello", "world"}); err != nil {
		t.Error(err)
	}
}

func TestCommandArgs(t *testing.T) {
	cli := New()

	cli.Command("hello", func(c Command) {
		for _, a := range c.Args {
			fmt.Println("Hello " + a)
		}
	})

	if err := cli.Run([]string{"hello", "world", "foo", "bar"}); err != nil {
		t.Error(err)
	}
}

func TestSubCommand(t *testing.T) {
	cli := New()

	cmd := cli.Command("remote", func(c Command) {})

	cmd.Command("add", func(c Command) {
		fmt.Println("add " + c.Args.Get(0))
	})

	if err := cli.Run([]string{"remote", "add", "origin"}); err != nil {
		t.Error(err)
	}
}

func TestCommandNotFound(t *testing.T) {
	cli := New()

	cli.Command("hello", nil)

	err := cli.Run([]string{})

	if err == nil {
		t.Error("expected command to fail")
	}

	fmt.Println(err)
}

func TestSubCommandNotFound(t *testing.T) {
	cli := New()

	cmd := cli.Command("remote", nil)

	cmd.Command("add", func(c Command) {
		fmt.Println("add " + c.Args.Get(0))
	})

	err := cli.Run([]string{"remote", "foo", "origin"})

	if err == nil {
		t.Error("expected command to fail")
	}

	fmt.Println(err)
}

func TestFlagArg(t *testing.T) {
	cli := New()

	cmd := cli.Command("hello", func(c Command) {
		cnt, err := c.Flags.GetInt("count")

		if err != nil {
			t.Error(err)
		}

		if cnt != 5 {
			t.Errorf("expected count to be 5, it was %d\n", cnt)
		}

		for i := 0; i < cnt; i++ {
			fmt.Println("hello " + c.Args.Get(0))
		}
	})

	cmd.AddFlag(&Flag{
		Name:     "count",
		Short:    "-c",
		Long:     "--count",
		Argument: true,
		Default:  1,
	})

	if err := cli.Run([]string{"hello", "world", "-c", "5"}); err != nil {
		t.Error(err)
	}

	if err := cli.Run([]string{"hello", "world", "--count", "5"}); err != nil {
		t.Error(err)
	}

	if err := cli.Run([]string{"hello", "world", "--count=5"}); err != nil {
		t.Error(err)
	}
}

func TestFlagBool(t *testing.T) {
	cli := New()

	cmd := cli.Command("hello", func(c Command) {
		arg := c.Args.Get(0)

		if c.Flags.IsSet("uppercase") {
			arg = strings.ToUpper(arg)
		} else {
			t.Errorf("expected uppercase flag to be true, it was false\n")
		}

		fmt.Println("hello " + arg)
	})

	cmd.AddFlag(&Flag{
		Name:  "uppercase",
		Short: "-u",
		Long:  "--uppercase",
	})

	if err := cli.Run([]string{"hello", "world", "-u"}); err != nil {
		t.Error(err)
	}

	if err := cli.Run([]string{"hello", "world", "--uppercase"}); err != nil {
		t.Error(err)
	}
}

func TestGlobalFlag(t *testing.T) {
	cli := New()

	cli.AddFlag(&Flag{
		Name:  "help",
		Short: "-h",
		Long:  "--help",
	})

	cmd := cli.Command("hello", func(c Command) {
		if c.Flags.IsSet("help") {
			fmt.Println("say hello")
		}
	})

	cmd.Command("sub", func(c Command) {
		if c.Flags.IsSet("help") {
			fmt.Println("this is a sub-command")
		}
	})

	if err := cli.Run([]string{"hello", "-h"}); err != nil {
		t.Error(err)
	}

	if err := cli.Run([]string{"hello", "--help"}); err != nil {
		t.Error(err)
	}

	if err := cli.Run([]string{"hello", "sub", "--help"}); err != nil {
		t.Error(err)
	}
}

func TestFlagHandler(t *testing.T) {
	cli := New()

	cli.AddFlag(&Flag{
		Name:    "help",
		Short:   "-h",
		Long:    "--help",
		Handler: func(f Flag, c Command) {
			if c.Name == "" {
				fmt.Println("usage for main command")
			} else {
				fmt.Println("usage for " + c.Name)
			}
		},
	})

	cli.Main(func(c Command) {
		fmt.Println("main command")
	})

	cli.Command("foo", func(c Command) {
		fmt.Println("foo command")
	})

	if err := cli.Run([]string{"--help"}); err != nil {
		t.Error(err)
	}

	if err := cli.Run([]string{"foo", "--help"}); err != nil {
		t.Error(err)
	}
}

func TestExclusiveTestHandler(t *testing.T) {
	cli := New()

	cli.AddFlag(&Flag{
		Name:      "help",
		Short:     "-h",
		Long:      "--help",
		Exclusive: true,
		Handler:   func(f Flag, c Command) {
			if c.Name == "" {
				fmt.Println("usage for main command")
			} else {
				fmt.Println("usage for " + c.Name)
			}
		},
	})

	cli.Main(func(c Command) {})
	cli.Command("foo", func(c Command) {})

	if err := cli.Run([]string{"--help"}); err != nil {
		t.Error(err)
	}

	if err := cli.Run([]string{"foo", "--help"}); err != nil {
		t.Error(err)
	}
}
