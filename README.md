## texplate - CLI wrapper around Golang text/template package

```sh
$ ./texplate execute -h

- Uses Golang's text/template syntax
- Includes Sprig template helpers
- The input files must contain a map in YAML/JSON format
- The template file format is flexible if '--output-format=preserve', otherwise the template must be YAML/JSON

Usage:
  texplate execute <template.yml> [flags]

Flags:
  -h, --help                     help for execute
  -f, --input-file stringSlice   (optional) an input file containing key-value pair to interpolate into the template
  -o, --output-format string     (optional) renders interpolated template in the given format. Accepts 'preserve', 'yaml', or 'json' (default "preserve")
```
