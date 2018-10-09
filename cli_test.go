package cli

import (
	"fmt"
	"testing"
)

func TestMainCommand(t *testing.T) {
	cli := New()

	cli.Main(func(c Command) {
		fmt.Println("TestMainCommand")
	})

	if err := cli.Run([]string{}); err != nil {
		t.Error(err)
	}
}

func TestCommand(t *testing.T) {
	cli := New()

	cli.Command("hello", func(c Command) {
		fmt.Println("hello", c.Args.Get(0))
	})

	if err := cli.Run([]string{"hello", "world"}); err != nil {
		t.Error(err)
	}

	if err := cli.Run([]string{"foo", "world"}); err == nil {
		t.Errorf("expected command 'foo world' to fail, it did not\n")
	} else {
		fmt.Println(err)
	}
}

func TestSubCommand(t *testing.T) {
	cli := New()

	cmd := cli.Command("remote", nil)
	cmd.Command("add", func(c Command) {
		fmt.Println("remote add", c.Args.Get(0))
	})

	if err := cli.Run([]string{"remote", "add", "origin"}); err != nil {
		t.Error(err)
	}

	if err := cli.Run([]string{"remote", "foo", "origin"}); err == nil {
		t.Errorf("expected command 'remote foo origin' to fail, it did not\n")
	} else {
		fmt.Println(err)
	}
}

func TestFlagArgument(t *testing.T) {
	cli := New()

	cmd := cli.Command("hello", func(c Command) {
		count, err := c.Flags.GetInt("count")

		if err != nil {
			t.Error(err)
		}

		if count != 5 {
			t.Errorf("expected '--count' flag to be '5', it was not\n")
		}

		for i := 0; i < count; i++ {
			fmt.Println("hello", c.Args.Get(0))
		}
	})

	cmd.AddFlag(&Flag{
		Name:     "count",
		Short:    "-c",
		Long:     "--count",
		Argument: true,
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

func TestFlagSliceArgument(t *testing.T) {
	cli := New()

	cmd := cli.Command("hello", func(c Command) {
		names := c.Flags.GetSlice("names", ",")

		for _, n := range names {
			fmt.Println("hello", n)
		}
	})

	cmd.AddFlag(&Flag{
		Name:     "names",
		Long:     "--names",
		Argument: true,
	})

	if err := cli.Run([]string{"hello", "--names=sam,bill,ted"}); err != nil {
		t.Error(err)
	}
}

func TestFlagDefaultArgument(t *testing.T) {
	cli := New()

	defaultValue := 1

	cmd := cli.Command("hello", func(c Command) {
		count, err := c.Flags.GetInt("count")

		if err != nil {
			t.Error(err)
		}

		if count != defaultValue {
			t.Errorf(
				"expected '--count' flag to be '%d' it was '%d'\n",
				defaultValue,
				count,
			)
		}
	})

	cmd.AddFlag(&Flag{
		Name:     "count",
		Short:    "-c",
		Long:     "--count",
		Argument: true,
		Default:  defaultValue,
	})

	if err := cli.Run([]string{"hello", "world"}); err != nil {
		t.Error(err)
	}
}

func TestGlobalFlag(t *testing.T) {
	cli := New()

	cli.AddFlag(&Flag{
		Name: "help",
		Long: "--help",
	})

	cli.Command("hello", func(c Command) {
		if c.Flags.IsSet("help") {
			fmt.Println("say hello")
			return
		} else {
			t.Errorf("expected flag '--help' to be set on 'hello'\n")
		}
	})

	cmd := cli.Command("remote", nil)
	cmd.Command("add", func(c Command) {
		if c.Flags.IsSet("help") {
			fmt.Println("add a remote")
		} else {
			t.Errorf("expected flag '--help' to be set on 'remote add'\n")
		}
	})

	if err := cli.Run([]string{"hello", "--help"}); err != nil {
		t.Error(err)
	}

	if err := cli.Run([]string{"remote", "add", "--help"}); err != nil {
		t.Error(err)
	}
}

func TestGlobalFlagHandler(t *testing.T) {
	cli := New()

	executed := false

	cli.AddFlag(&Flag{
		Name:    "help",
		Long:    "--help",
		Handler: func(f Flag, c Command) {
			fmt.Println("usage for command", c.Name)
		},
	})

	cli.Command("hello", func(c Command) {
		fmt.Println("hello", c.Args.Get(0))

		executed = true
	})

	if err := cli.Run([]string{"hello", "--help"}); err != nil {
		t.Error(err)
	}

	if !executed {
		t.Errorf("expected command 'hello' to run\n")
	}
}

func TestExclusiveGlobalFlagHandler(t *testing.T) {
	cli := New()

	executed := false

	cli.AddFlag(&Flag{
		Name:      "help",
		Long:      "--help",
		Exclusive: true,
		Handler:   func(f Flag, c Command) {
			fmt.Println("usage for command", c.Name)
		},
	})

	cli.Command("hello", func(c Command) {
		fmt.Println("hello", c.Args.Get(0))

		executed = true
	})

	if err := cli.Run([]string{"hello", "--help"}); err != nil {
		t.Error(err)
	}

	if executed {
		t.Errorf("expected command 'hello' to not run\n")
	}
}

func TestFlagNotFound(t *testing.T) {
	cli := New()

	cli.Command("hello", func(c Command) {})

	if err := cli.Run([]string{"hello", "--foo"}); err == nil {
		t.Errorf("expected 'hello --foo' to fail\n")
	}
}

func TestNilCommandHandler(t *testing.T) {
	cli := New()

	executed := false

	cli.NilHandler(func(c Command) {
		fmt.Println("nil handler", c.Name)

		executed = true
	})

	cli.Main(nil)

	if err := cli.Run([]string{}); err != nil {
		t.Error(err)
	}

	if !executed {
		t.Errorf("expected nil handler to execute\n")
	}
}

func TestCommandFullName(t *testing.T) {
	cli := New()

	fullName := "remote-add"

	cli.NilHandler(func(c Command) {
		actual := c.FullName()

		if actual != fullName {
			t.Errorf("expected full name to be %s go %s\n", fullName, actual)
		}
	})

	cmd := cli.Command("remote", nil)
	cmd.Command("add", nil)

	if err := cli.Run([]string{"remote", "add"}); err != nil {
		t.Error(err)
	}
}

func TestMultipleSameFlag(t *testing.T) {
	cli := New()

	cmd := cli.Main(func(c Command) {
		headers := c.Flags.GetAll("header")

		if len(headers) != 3 {
			t.Errorf("expected three header flags\n")
		}

		for _, h := range headers {
			fmt.Println(h.Value)
		}
	})

	cmd.AddFlag(&Flag{
		Name:     "header",
		Short:    "-H",
		Argument: true,
		Default:  "",
	})

	flags := []string{"-H", "foo", "-H", "bar", "-H", "zap"}

	if err := cli.Run(flags); err != nil {
		t.Error(err)
	}
}
