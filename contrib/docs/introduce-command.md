## How to introduce new Command?

### Implements interfaces

* `/internal/cmd`.`CommandPreparator` 
    - `Prepare()` provides the `cobra.Command` and should attach any flags
* `/internal/cmd`.`ValidatedCommand`
    - `Validate()` should be implemented on any command that has to validate input parameters
* `/internal/cmd`.`Command`
    - `Run()` executes the main logic behind the command

### Optinal interfaces

* `/internal/cmd`.`FormattedCommand` should be implemented if the command supports different output formatting through a --format or -f flag
    - `SetOutputFormat` - sets the output format
* `/internal/cmd`.`ConfirmedCommand` should be implemented if the command should ask for user confirmation prior execution
    - `AskForConfirmation()` - asks the user for confirmation, see the helper function `CommonConfirmationPrompt`
* `/internal/cmd`.`HiddenUsageCommand`
    - `HideUsage()` should return true if the command should NOT return its usage doc without the *--help* flag

### How to handle flags?

There are several predefined flags that can be added to commands:

* `--mode` - can be added with `cmd.AddModeFlag()` it allows to change how requests should be executed async or sync
* `--output` or `-o` - can be added with `cmd.AddFormatFlag()` it allows to change the output of the command. For example json or yaml or table
* `--field-query` or `-f` and  `--label-query` or `-l` - can be added with `cmd.AddQueryingFlags()`, it allows to add query to the request to Service Manager

Arbitrary flags can be added in the `Prepare()` of each command as by Cobra framework https://github.com/spf13/cobra.

### Where to register the command?

#### Is this SM specific command?

Note: **If the command needs to authenticate and execute requests against Service Manager, then this is an SM specific command.**

In the `main.go` file in the `smCommandsGroup` variable, you can see all the commands that are meant to call Service Manager. Add your command in the same list and this should be enough.

### Is it common CLI command?

Note: **smctl login is an exception**

If the command is for example to get the version of the command line tool, then it is not SM specific.
Such a command can be added to the list of `normalCommandsGroup` variable.