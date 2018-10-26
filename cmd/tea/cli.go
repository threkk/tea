package tea

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// Custom flag to parse the paths to the error codes html files. It first
// splits the string by commas and after it tries to match the values of the
// key with one error code.
type pages struct {
	Page403 string
	Page404 string
	Page500 string
}

func (p *pages) String() string {
	return p.String()
}

func (p *pages) Set(value string) error {
	csv := strings.Split(value, ",")
	if len(csv) > 3 || len(csv) < 1 {
		return errors.New("pages: invalid amount of options")
	}

	re := regexp.MustCompile(`(?m)(403|404|500)=(\S+)`)
	for _, val := range csv {
		matches := re.FindAllString(val, -1)
		if matches != nil {
			k := matches[0]
			v := matches[1]

			switch k {
			case "403":
				(*p).Page403 = v
			case "404":
				(*p).Page404 = v
			case "500":
				(*p).Page500 = v
			default:
				return errors.New("pages: invalid key")
			}
		}
	}
	return nil
}

// Messages
const usageMsg = `
Usage: tea [options] <path>

Starts the server in the <path>. The path will load any "_config.yaml" file it
finds. The options can be used to configure the server or to overwrite any value
from the configuration file. It is mandatory to define a domain name.

Options:
`

const (
	// domain
	domain            = "domain"
	domainDefault     = ""
	domainDescription = "Site address that will be used by the server (mandatory)"

	// port
	port            = "port"
	portDefault     = 8080
	portDescription = "Server port (default 8080)"

	// ssl
	ssl            = "cert"
	sslDefault     = ""
	sslDescription = "SSL certificate to use"

	// pages
	pg            = "pages"
	pgDescription = "Comma-separated list of error codes and pages to render."

	// md
	md            = "markdown"
	mdDefault     = false
	mdDescription = "Enable markdown (default false)"

	// spa
	spa            = "spa"
	spaDefault     = false
	spaDescription = "Enable routes for single page applications (default false)"

	// dev
	dev            = "dev"
	devDefault     = false
	devDescription = "Enable HTTP development server (default false)"

	// ui
	ui           = "ui"
	uiDefault    = false
	uiDescrition = "Enable UI for folders (default false)"
)

// CLIValues variables for the start command. It defines most the options from
// the configurations but the ones that require authentication. They allow to
// overwrite the configuration from the _config.yaml file.
type CLIValues struct {
	Domain string
	Port   uint
	SSL    string
	Pages  *pages
	MD     bool
	Spa    bool
	Dev    bool
	UI     bool
}

// NewCLI Creates a new struct and initialises the fields for a cli.
func NewCLI(name string, out io.Writer) (*flag.FlagSet, *CLIValues) {
	// Initialise the flagset.
	flagset := flag.NewFlagSet(name, flag.ExitOnError)

	// Initialise the struct.
	values := &CLIValues{
		Domain: "",
		Port:   8080,
		SSL:    "",
		Pages:  &pages{},
		MD:     false,
		Spa:    false,
		Dev:    false,
		UI:     false,
	}

	// Definte the output of the CLI.
	flagset.SetOutput(out)

	// Define the help message.
	flagset.Usage = func() {
		fmt.Fprint(out, usageMsg)
		flagset.PrintDefaults()
	}

	flagset.StringVar(&values.Domain, domain, domainDefault, domainDescription)
	flagset.UintVar(&values.Port, port, portDefault, portDescription)
	flagset.StringVar(&values.SSL, ssl, sslDefault, sslDescription)
	flagset.Var(values.Pages, pg, pgDescription)
	flagset.BoolVar(&values.MD, md, mdDefault, mdDescription)
	flagset.BoolVar(&values.Spa, spa, spaDefault, spaDescription)
	flagset.BoolVar(&values.Dev, dev, devDefault, devDescription)
	flagset.BoolVar(&values.UI, ui, uiDefault, uiDescrition)

	return flagset, values
}
