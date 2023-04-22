package cobrax

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/zeromicro/go-zero/tools/goctl/internal/flags"
)

type Option func(*cobra.Command)

func WithRunE(runE func(*cobra.Command, []string) error) Option {
	return func(cmd *cobra.Command) {
		cmd.RunE = runE
	}
}

func WithRun(run func(*cobra.Command, []string)) Option {
	return func(cmd *cobra.Command) {
		cmd.Run = run
	}
}

func WithArgs(arg cobra.PositionalArgs) Option {
	return func(command *cobra.Command) {
		command.Args = arg
	}
}

func WithHidden() Option {
	return func(command *cobra.Command) {
		command.Hidden = true
	}
}

type Command struct {
	*cobra.Command
}

type FlagSet struct {
	*pflag.FlagSet
}

func (f *FlagSet) StringVar(p *string, name string) {
	f.StringVarWithDefaultValue(p, name, "")
}

func (f *FlagSet) StringVarWithDefaultValue(p *string, name string, value string) {
	f.FlagSet.StringVar(p, name, value, "")
}

func (f *FlagSet) StringVarP(p *string, name, shorthand string) {
	f.StringVarPWithDefaultValue(p, name, shorthand, "")
}

func (f *FlagSet) StringVarPWithDefaultValue(p *string, name, shorthand string, value string) {
	f.FlagSet.StringVarP(p, name, shorthand, value, "")
}

func (f *FlagSet) BoolVar(p *bool, name string) {
	f.BoolVarWithDefaultValue(p, name, false)
}

func (f *FlagSet) BoolVarWithDefaultValue(p *bool, name string, value bool) {
	f.FlagSet.BoolVar(p, name, value, "")
}

func (f *FlagSet) BoolVarP(p *bool, name, shorthand string) {
	f.BoolVarPWithDefaultValue(p, name, shorthand, false)
}

func (f *FlagSet) BoolVarPWithDefaultValue(p *bool, name, shorthand string, value bool) {
	f.FlagSet.BoolVarP(p, name, shorthand, value, "")
}

func (f *FlagSet) IntVar(p *int, name string) {
	f.IntVarWithDefaultValue(p, name, 0)
}

func (f *FlagSet) IntVarWithDefaultValue(p *int, name string, value int) {
	f.FlagSet.IntVar(p, name, value, "")
}

func (f *FlagSet) StringSliceVarP(p *[]string, name, shorthand string) {
	f.FlagSet.StringSliceVarP(p, name, shorthand, []string{}, "")
}

func (f *FlagSet) StringSliceVarPWithDefaultValue(p *[]string, name, shorthand string, value []string) {
	f.FlagSet.StringSliceVarP(p, name, shorthand, value, "")
}

func (f *FlagSet) StringSliceVar(p *[]string, name string) {
	f.StringSliceVarWithDefaultValue(p, name, []string{})
}

func (f *FlagSet) StringSliceVarWithDefaultValue(p *[]string, name string, value []string) {
	f.FlagSet.StringSliceVar(p, name, value, "")
}

func NewCommand(use string, opts ...Option) *Command {
	c := &Command{
		Command: &cobra.Command{
			Use: use,
		},
	}

	for _, opt := range opts {
		opt(c.Command)
	}

	return c
}

func (c *Command) AddCommand(cmds ...*Command) {
	for _, cmd := range cmds {
		c.Command.AddCommand(cmd.Command)
	}
}

func (c *Command) Flags() *FlagSet {
	set := c.Command.Flags()
	return &FlagSet{
		FlagSet: set,
	}
}

func (c *Command) PersistentFlags() *FlagSet {
	set := c.Command.PersistentFlags()
	return &FlagSet{
		FlagSet: set,
	}
}

func (c *Command) MustInit() {
	commands := append([]*cobra.Command{c.Command}, getCommandsRecursively(c.Command)...)
	for _, command := range commands {
		commandKey := getCommandName(command)
		if len(command.Short) == 0 {
			command.Short = flags.Get(commandKey + ".short")
		}
		if len(command.Long) == 0 {
			command.Long = flags.Get(commandKey + ".long")
		}
		if len(command.Example) == 0 {
			command.Example = flags.Get(commandKey + ".example")
		}
		command.Flags().VisitAll(func(flag *pflag.Flag) {
			flag.Usage = flags.Get(fmt.Sprintf("%s.%s", commandKey, flag.Name))
		})
		command.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
			flag.Usage = flags.Get(fmt.Sprintf("%s.%s", commandKey, flag.Name))
		})
	}
}

func getCommandName(cmd *cobra.Command) string {
	if cmd.HasParent() {
		return getCommandName(cmd.Parent()) + "." + cmd.Name()
	}
	return cmd.Name()
}

func getCommandsRecursively(parent *cobra.Command) []*cobra.Command {
	var commands []*cobra.Command
	for _, cmd := range parent.Commands() {
		commands = append(commands, cmd)
		commands = append(commands, getCommandsRecursively(cmd)...)
	}
	return commands
}
