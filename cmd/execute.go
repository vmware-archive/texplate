package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/ljfranklin/texplate-cli/interpolater"
	"github.com/spf13/cobra"
)

var (
	inputFiles   []string
	outputFormat string
)

var executeCmd = &cobra.Command{
	Use:   "execute <template.yml>",
	Short: "Interpolate input files into given template",
	Long: `
- Uses Golang's text/template syntax
- Includes Sprig template helpers
- The input files must contain a map in YAML/JSON format
- The template file format is flexible if '--output-format=preserve', otherwise the template must be YAML/JSON
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("must specify template path as first positional arg")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		basePath := args[0]

		i := interpolater.Interpolater{
			Writer:       os.Stdout,
			OutputFormat: outputFormat,
		}
		err := i.Execute(basePath, inputFiles)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(executeCmd)
	executeCmd.Flags().StringSliceVarP(&inputFiles, "input-file", "f", []string{}, "(optional) an input file containing key-value pair to interpolate into the template")
	executeCmd.Flags().StringVarP(&outputFormat, "output-format", "o", interpolater.FormatPreserve, fmt.Sprintf("(optional) renders interpolated template in the given format. Accepts '%s', '%s', or '%s'", interpolater.FormatPreserve, interpolater.FormatYAML, interpolater.FormatJSON))
}
