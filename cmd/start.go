package cmd

import (
	"errors"
	"flag"
	"fmt"
	"io"
)

var usageMsg = `
Usage: tea start [options] <path>

Starts the server in the <path>. The path will load any "_config.yaml" file it
finds. The options can be used to configure the server or to overwrite any value
from the configuration file. It is mandatory to define a domain name.

Options:
`

// StartCli CLI flag set for the start command. It defines most the options from
// the configurations but the ones that require authentication. They allow to
// overwrite the configuration from the _config.yaml file. It is loaded with the
// values of
type StartCli struct {
	Flagset  *flag.FlagSet
	Path     string
	Domain   string
	SSL      string
	Port     uint
	Markdown bool
	Spa      bool
	Dev      bool
	UI       bool
}

// NewStartCli Creates a new struct and initialises the fields for a StartCli.
func NewStartCli(out io.Writer) *StartCli {
	// Initialise the struct.
	start := &StartCli{
		Flagset:  flag.NewFlagSet("tea start", flag.ExitOnError),
		Path:     "",
		Domain:   "",
		SSL:      "",
		Port:     8080,
		Markdown: false,
		Spa:      false,
		Dev:      false,
		UI:       false,
	}

	// Definte the output of the CLI.
	start.Flagset.SetOutput(out)

	// Define the help message.
	start.Flagset.Usage = func() {
		fmt.Println(usageMsg)
		start.Flagset.PrintDefaults()
	}

	start.Flagset.StringVar(&start.Domain, "domain", "", "Site address that will be used by the server (mandatory)")
	start.Flagset.StringVar(&start.SSL, "cert", "", "SSL certificate to use")
	start.Flagset.UintVar(&start.Port, "port", 8080, "Server port (default 8080)")
	start.Flagset.BoolVar(&start.Markdown, "markdown", false, "Enable markdown (default false)")
	start.Flagset.BoolVar(&start.Spa, "spa", false, "Enable routes for single page applications (default false)")
	start.Flagset.BoolVar(&start.Dev, "dev", false, "Enable HTTP development server (default false)")
	start.Flagset.BoolVar(&start.UI, "ui", false, "Enable UI for folders (default false)")

	return start
}

// Parse Parse the provided arguments. It consumes the arguments and fills the
// struct with the stored values.
func (s *StartCli) Parse(arguments []string) error {
	err := s.Flagset.Parse(arguments)
	if err != nil {
		return err
	}

	args := s.Flagset.Args()

	if len(args) != 1 {
		return errors.New("more than one path provided")
	}

	s.Path = args[0]
	return nil
}
