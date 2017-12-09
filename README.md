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

#### Included helper functions:

- All the [Sprig](http://masterminds.github.io/sprig/) helpers, e.g. `env`, `trim`, and `list`
- `cidrhost <cidr> <hostIndex>`:
  - Returns an IP at the given index within that CIDR, e.g. `cidrhost 10.0.1.0/24 2` returns `10.0.1.2`
  - Adapted from Terraform's [cidrhost](https://www.terraform.io/docs/configuration/interpolation.html#cidrhost-iprange-hostnum-)
