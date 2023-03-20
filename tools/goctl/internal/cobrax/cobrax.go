package cobrax

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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

type Command struct {
	*cobra.Command
}

type FlagSet struct {
	*pflag.FlagSet
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
