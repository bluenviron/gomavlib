// dialect-import command.
package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"

	"github.com/bluenviron/gomavlib/v2/pkg/conversion"
)

var cli struct {
	Link bool   `help:"Link included definitions instead of including them into the main definition"`
	XML  string `arg:"" help:"Path or url pointing to a XML Mavlink dialect"`
}

func run(args []string) error {
	parser, err := kong.New(&cli,
		kong.Description("Convert Mavlink dialects from XML format to Go format."),
		kong.UsageOnError())
	if err != nil {
		return err
	}

	_, err = parser.Parse(args)
	if err != nil {
		return err
	}

	return conversion.Convert(cli.XML, cli.Link)
}

func main() {
	err := run(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERR: %s\n", err)
		os.Exit(1)
	}
}
