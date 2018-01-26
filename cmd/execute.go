package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/pivotal-cf/texplate/interpolater"
	"github.com/spf13/cobra"
)

var (
	inputFiles   []string
	outputFormat string
	outputFile   string
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

		var writer io.Writer
		if outputFile == "-" {
			writer = os.Stdout
		} else {
			f, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to open %s: %s", outputFile, err)
				os.Exit(1)
			}
			defer func() {
				err := f.Close()
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
			}()
			writer = f
		}

		i := interpolater.Interpolater{
			Writer:       writer,
			OutputFormat: outputFormat,
		}
		err := i.Execute(basePath, inputFiles)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if outputFile != "-" {
			fmt.Fprintf(os.Stderr, "Wrote output to %s", outputFile)
		}
	},
}

func init() {
	rootCmd.AddCommand(executeCmd)
	executeCmd.Flags().StringSliceVarP(&inputFiles, "input-file", "f", []string{}, "(optional) an input file containing key-value pair to interpolate into the template")
	executeCmd.Flags().StringVarP(&outputFormat, "output-format", "o", interpolater.FormatPreserve, fmt.Sprintf("(optional) renders interpolated template in the given format. Accepts '%s', '%s', or '%s'", interpolater.FormatPreserve, interpolater.FormatYAML, interpolater.FormatJSON))
	executeCmd.Flags().StringVarP(&outputFile, "output-file", "", "-", "(optional) writes output to given filepath. Defaults to writing to STDOUT")
}
