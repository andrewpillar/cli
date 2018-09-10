# cli

Simple library for parsing commands, and flags from the command line in Go.

- [Quickstart](#quickstart)
- [Writing Commands](#writing-commands)
  * [Main Command](#main-command)
  * [Commands](#commands)
  * [Sub-commands](#sub-commands)
  * [Flags](#flags)
  * [Arguments](#arguments)

## Quickstart

This is a simple library that makes parsing commands, flags, arguments from the command line easier to do in Go. Here's what it looks like:

```go
package main

import (
    "fmt"
    "os"

    "github.com/andrewpillar/cli"
)

func main() {
    c := cli.New()

    c.Main(func(c cli.Command) {
        fmt.Println("Say something"!)
    })

    helloCmd := c.Command("hello", func(c cli.Command) {
        cnt, err := c.Flags.GetInt("count")

        if err != nil {
            fmt.Fprintf(os.Stderr, "%s\n", err)
            os.Exit(1)
        }

        for i := 0; i < cnt; i++ {
            fmt.Println("Hello " + c.Args.Get(0))
        }
    })

    helloCmd.AddFlag(&cli.Flag{
        Name:     "count",
        Short:    "-c",
        Long:     "--count",
        Argument: true,
        Default:  1,
    })

    if err := c.Run(os.Args[1:]); err != nil {
        fmt.Fprintf(os.Stderr, "%s\n", err)
        os.Exit(1)
    }
}
```

And here is what the above program will produce once it has been built, and run:

```
$ say hello world --count=5
Hello world
Hello world
Hello world
Hello world
Hello world
```
You can start using the library by installing it with `go get`:

```
go get github.com/andrewpillar/cli
```

## Writing Commands

This library can allow you to define a single main command, or multiple commands throughout. Each command which is created with this library can also have sub-commands.

Every command defined will take a call-back for executing that command. This call-back will receive the command which is being executed as its only parameter.

### Main Command

The main command for the program can be defined with the `Main` method. This is the default command that will be executed if no other command is specified.

```go
c.Main(func(c cli.Command) {
    fmt.Println("Say something!")
})
```

```
$ say
Say something!
```

### Commands

Commands can be defined with a name, and a call-back with the `Command` method. This can allow for command line program to have multiple commands throughout.

```go
c.Command("hello", func(c cli.Command) {
    fmt.Println("Hello!")
})
```

### Sub-commands

Each command created with the library can have further sub-commands. The `Command` method will return a pointer to the newly created `cli.Command` type. This type has a method called `Command` which can allow for commands to be defined on the command itself.

```go
cmd := c.Command("remote", func(c cli.Command) {
    fmt.Println("Doing something with a remote")
})

cmd.Command("add", func(c cli.Command) {
    fmt.Println("Adding a remote")
})
```

### Flags

Command line flags can either be added to individual commands, or to the entire command line program itself. This is done by calling the `AddFlag` method and passing a pointer to the `cli.Flag` type.

```go
helloCmd := c.Command("hello", func(c cli.Command) {
    cnt, err := c.Flags.GetInt("count")

    if err != nil {
        fmt.Fprintf(os.Stderr, "%s\n", err)
        os.Exit(1)
    }

    for i := 0; i < cnt; i++ {
        fmt.Println("Hello " + c.Args.Get(0))
    }
})

helloCmd.AddFlag(&cli.Flag{
    Name:     "count",
    Short:    "-c",
    Long:     "--count",
    Argument: true,
    Default:  1,
})
```

In the above example we added the `count` flag to the `hello` command in our program. We specified the short, and long flags for this flag, `-c`, and `--count` respectively, and stated that this flag takes an argument with a default value of `1`.

If this flag is passed to the `hello` command and given no argument, or not passed to the command at all, then the default value of `1` will be returned when we retrieve the flag's value from the command.

There are multiple ways of retrieving a flag's value, since its value can be of varying types. Listed below are the different getter methods for accessing a flag's value:

* `GetInt`
* `GetIn32`
* `GetIn64`
* `GetString`
* `GetSlice`

The `GetSlice` getter is a utility getter that wraps `GetString`, and takes an additional parameter for the delimiter on which to split the string.

```go
c.Command("hello", func(c cli.Command) {
    names := c.Flags.GetSlice("names", ",")
})
```

For flags which do not take values, and merely serve as boolean flags then the `IsSet` method should be used for determining whether the flag was passed to the command.

```go
c.Command("hello", func(c cli.Command) {
    if c.Flags.IsSet("help") {
        fmt.Println("Say hello")
        return
    }
})
```

Should a non-existent flag try and be accessed then the default zero-value for the type trying to be retrieved will be returned.

Global flags can be declared on the program itself. These will be added to every command which is created.

```go
c := cli.New()

c.AddFlag(&cli.Flag{
    Name: "help",
    Long: "--help",
})
```

Handlers can also be set on flags. A flag handler will be passed the flag, and the command which had that flag passed to it. Setting the `Exclusive` property on the given `cli.Flag` type to `true` will prevent the command which passed the flag from running. This is useful if you want to display usage information for that command via a `--help` flag without having the command itself be run.

```go
c := cli.New()

c.AddFlag(&cli.Flag{
    Name:      "help",
    Long:      "--help",
    Exclusive: true,
    Handler:   func(f cli.Flag, c cli.Command) {
        fmt.Println("usage for " + c.Name)
    },
})
```

### Arguments

Command line arguments can be retrieved by accessing the `Args` property on the given command, and calling the `Get` method passing through the index of the argument you want to get.

```go
c.Command("hello", func(c cli.Command) {
    fmt.Println("Hello " + c.Args.Get(0))
})
```

The `Args` property can also be looped through.

```go
c.Command("hello", func(c cli.Command) {
    for _, a := range c.Args {
        fmt.Println("Hello " + a)
    }
})
