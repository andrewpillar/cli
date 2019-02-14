# cli

Simple library for building command line programs in Go.

* [Quick Start](#quick-start)
* [Creating a Program](#creating-a-program)
* [Defining Commands](#defining-commands)
* [Adding Flags](#adding-flags)
* [Working with Flags](#working-with-flags)
* [Arguments](#arguments)

## Quick Start

To get started, simply add the repository to your project as an import path.

```go
package main

import (
    "fmt"
    "os"

    "github.com/andrewpillar/cli"
)

func main() {
    c := cli.New()

    c.MainCommand(func(cmd cli.Command) {
        fmt.Println("Say something!")
    })

    helloCmd := c.Command("hello", func(cmd cli.Command) {
        count, err := cmd.Flags.GetInt("count")

        if err != nil {
            fmt.Fprintf(os.Stderr, "%s\n", err)
            os.Exit(1)
        }

        for i := 0; i < count; i++ {
            fmt.Println("Hello", cmd.Args.Get(0))
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

And here is what the above program will produce once it has been built and run.

```
$ say hello world --count 5
Hello world
Hello world
Hello world
Hello world
Hello world
```

## Creating a Program

Before defining commands and flags for your program, we must first create a new one. This is done by calling the `New` function. This will return a pointer to the `cli.Cli` type, from which we will be able to define commands and flags for our program.

```go
package main

import (
    "fmt"
    "os"

    "github.com/andrewpillar/cli"
)

func main() {
    c := cli.New()
}
```

## Defining Commands

Within this library there are two types of commands, a single main command, and a regular command. The single main command has no name, and is run when no other command is specified for the program. Regular commands on the other hand, do have names.

Each command that is defined on the program will have a command handler associated with it. This handler is a simple callback function that takes a `cli.Command` type as its only argument. From this type, you will be able to access the arguments and flags that have been passed to the command.

The program's main command can be defined by calling the `MainCommand` method on the `cli.Cli` type we were returned from the `New` function. This method takes the command handler as its only argument.

```go
c := cli.New()

c.MainCommand(func(cmd cli.Command) {
    fmt.Println("Say something!")
})
```

Regular commands are defined bu calling the `Command` method on the `cli.Cli` type. Unlike the `MainCommand` method, this takes two parameters, the name of the command and the command handler.

```go
helloCmd := c.Command("hello", func(cmd cli.Command) {

})
```

## Adding Flags

Flags can either be added to the entire program or individual commands. When a flag is added to the entire program, it will be passed down to every command and sub-command in the program.

Each time a command is defined with either the `MainCommand` or `Command` method, a pointer to the recently created command will be returned. The `AddFlag` method can then be called on the returned command to add a flag to that specific command.

```go
helloCmd := c.Command("hello", func(cmd cli.Command) {

})

helloCmd.AddFlag(&cli.Flag{
    Name:     "count",
    Short:    "-c",
    Long:     "--count",
    Argument: true,
    Default:  1,
})
```

Program wide flags can be added by calling `AddFlag` directly on the `cli.Cli` type returned by `New`.

```go
c := cli.New()

c.AddFlag(&cli.Flag{
    Name: "help",
    Long: "--help",
})
```

## Working with Flags

Each flag added to a command will be given a name, along with a long and short version of that flag. The flag's name is what is used for accessing the underlying `cli.Flag` type. There are multiple methods on the `Flags` property on the `cli.Command` type. Each of the methods take the flag's name as their only argument.

If a flag has a value, then you would call the `Get<Type>` method on the `Flags` property to get the underlying value. A flag's value can either be `string`, `int`, `int8`, `int16`, `int32`, `int64`, `float32`, or `float64`.

```go
helloCmd := c.Command("hello", func(cmd cli.Command) {
    count, err := cmd.Flags.GetInt("count")
})
```

The methods that involved returning number types, return multiple values, the parsed value, and an `error`.

Boolean flags can be checked against by calling `IsSet`, this method takes the flag's name much like the prior `Get<Type>` methods. However this will return `true` or `false` depending on whether the flag was set on the command.

```go
helloCmd := c.Command("hello", func(cmd cli.Command) {
    if cmd.Flags.IsSet("help") {
    }

    count, err := cmd.Flags.GetInt("count")
})
```

If the same flag is passed multiple times to a single command, then you can retrieve all occurences of that flag with the `GetAll` method.

```go
helloCmd := c.Command("hello", func(cmd cli.Command) {
    for _, flag := range cmd.Flags.GetAll("count") {
        count, err := flag.GetInt()
    }
})
```

The `cli.Flag` type has a `Get<Type>` method for each type the flag could be, along with an `IsSet` method.

## Arguments

Arguments passed to the program, or an individual command, can be accessed via the `Args` property on the passed command. This type can be looped over, or individual arguments can be accessed by calling the `Get` method and passing the index of the argument you want.

```go
c.Command("hello", func(cmd cli.Command) {
    fmt.Println("Hello", cmd.Args.Get(0))

    for _, arg := range cmd.Args {
        fmt.Println("Hello", arg)
    }
})
```
