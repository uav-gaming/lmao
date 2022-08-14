package commands

import (
	"errors"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/sirupsen/logrus"
)

type CommandHandler func(discord.UserID) string

// Registory for keeping track of commands and its handlers.
// Absolutly **NOT** threadsafe
type CommandRegistrar struct {
	Commands []discord.Command
	handlers map[string]*Command
}

func NewCommandRegistrar() CommandRegistrar {
	return CommandRegistrar{
		make([]discord.Command, 0),
		make(map[string]*Command),
	}
}

func (reg *CommandRegistrar) AddCommand(name, description string) *Command {
	if _, exists := reg.handlers[name]; exists {
		logrus.Fatal("Adding command ", name, " but it already exists: %+v", reg.Commands)
	}

	reg.Commands = append(reg.Commands, discord.NewCommand(name, description))
	// Pointer to the command we just inserted
	cmdPtr := &reg.Commands[len(reg.Commands)-1]
	cmd := NewCommand(cmdPtr)
	reg.handlers[name] = &cmd
	return &cmd
}

// Handle a command request
func (reg *CommandRegistrar) HandleCommand(command *discord.CommandInteraction) (*string, error) {
	handler, found := reg.handlers[command.Name]
	if !found {
		logrus.Errorf("Received unknown command.\nRegistered commands:\n%+v\nReceived command:\n%+v", reg.Commands, command)
		return nil, errors.New("unknown command " + command.Name)
	}
	return handler.HandleSubcommand(command.Options)
}

type Command struct {
	command  *discord.Command
	handlers map[string]CommandHandler
}

// Create a new command.
func NewCommand(command *discord.Command) Command {
	return Command{
		command,
		make(map[string]CommandHandler),
	}
}

// Handle a request.
func (cmd *Command) HandleSubcommand(options discord.CommandInteractionOptions) (*string, error) {
	// Validate request.
	if len(options) != 1 {
		logrus.Errorf("Expect only one option but got %+v", options)
		return nil, errors.New("too much options")
	}
	option := &options[0] // TODO: is it a copy?
	if option.Type != discord.SubcommandOptionType {
		logrus.Errorf("Expect a subcommand but got %+v", options)
		return nil, errors.New("unexpected option type")
	}

	// Handle subcommand.
	subCmd := option.String()
	_, found := cmd.handlers[subCmd]
	if !found {
		logrus.Errorf("Received unknown command.\nRegistered commands:\n%+v\nReceived command:\n%+v", cmd.command, subCmd)
		return nil, errors.New("unknown subcommand " + subCmd)
	}
	// TODO: actually handle it
	logrus.Error("bruh you forgot to implment the handling logic")
	return nil, errors.New("go fix your `HandleSubcommand`")
}

// Add a subcommand within the group.
func (cmd *Command) AddSubcommand(name, description string, handler CommandHandler) {
	if _, exists := cmd.handlers[name]; exists {
		logrus.Fatal("Adding subcommand ", name, " but it already exists: %+v", *cmd)
	}
	cmd.handlers[name] = handler
}
