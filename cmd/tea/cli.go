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
	return fmt.Sprintf("{403:%s 404:%s 500:%s}", p.Page403, p.Page404, p.Page500)
}

func (p *pages) Set(value string) error {
	csv := strings.Split(value, ",")
	if len(csv) > 3 || len(csv) < 1 {
		return errors.New("pages: invalid amount of options")
	}

	re := regexp.MustCompile(`(?m)(403|404|500)=(\S+)`)
	for _, val := range csv {
		matches := re.FindStringSubmatch(val)
		if matches != nil {
			k := matches[1]
			v := matches[2]

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

	// HTML5
	html5            = "html5"
	html5Default     = false
	html5Description = "Enable routes for single page applications (default false)"

	// dev
	dev            = "dev"
	devDefault     = false
	devDescription = "Enable HTTP development server (default false)"

	// ui
	ui           = "ui"
	uiDefault    = false
	uiDescrition = "Enable UI for folders (default false)"
)

// CLI variables for the start command. It defines most the options from
// the configurations but the ones that require authentication. They allow to
// overwrite the configuration from the _config.yaml file.
type CLI struct {
	Domain string
	Port   uint
	SSL    string
	Pages  *pages
	MD     bool
	HTML5  bool
	Dev    bool
	UI     bool

	// unexported
	flag *flag.FlagSet
}

// Parse Wrapper around the flagset.
func (cli *CLI) Parse(args []string) error {
	return cli.flag.Parse(args)
}

// NewCLI Creates a new struct and initialises the fields for a cli.
func NewCLI(name string, out io.Writer) *CLI {
	// Initialise the struct.
	cli := &CLI{
		Domain: "",
		Port:   8080,
		SSL:    "",
		Pages:  &pages{},
		MD:     false,
		HTML5:  false,
		Dev:    false,
		UI:     false,
		flag:   flag.NewFlagSet(name, flag.ExitOnError),
	}

	// Definte the output of the CLI.
	cli.flag.SetOutput(out)

	// Define the help message.
	cli.flag.Usage = func() {
		fmt.Fprint(out, usageMsg)
		cli.flag.PrintDefaults()
	}

	cli.flag.StringVar(&cli.Domain, domain, domainDefault, domainDescription)
	cli.flag.UintVar(&cli.Port, port, portDefault, portDescription)
	cli.flag.StringVar(&cli.SSL, ssl, sslDefault, sslDescription)
	cli.flag.Var(cli.Pages, pg, pgDescription)
	cli.flag.BoolVar(&cli.MD, md, mdDefault, mdDescription)
	cli.flag.BoolVar(&cli.HTML5, html5, html5Default, html5Description)
	cli.flag.BoolVar(&cli.Dev, dev, devDefault, devDescription)
	cli.flag.BoolVar(&cli.UI, ui, uiDefault, uiDescrition)

	return cli
}
