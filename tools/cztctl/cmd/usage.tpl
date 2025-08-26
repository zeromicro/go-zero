{{blue "Usage:"}}{{if .Runnable}}
  {{green .UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{green .CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

{{blue "Aliases:"}}
  {{green .NameAndAliases}}{{end}}{{if .HasExample}}

{{blue "Examples:"}}
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

{{blue "Available Commands:"}}{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpadx .Name .NamePadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

{{blue "Flags:"}}
{{green .LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

{{blue "Global Flags:"}}
{{green .InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

{{blue "Additional help topics:"}}{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{green .CommandPath}} [command] --help" for more information about a command.{{end}}
