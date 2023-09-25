# gen-changelog

`gen-changelog` is a wrapper around the GitHub CLI [`gh`](https://cli.github.com/) which 
automatically generates boilerplate for [`go-changelog`](https://github.com/hashicorp/go-changelog) 
so you can easily write the changelog body without fiddling with files and figuring out PR numbers.

Currently the "Type" options for changelogs [here](https://github.com/kisunji/gen-changelog/blob/v1.0.0/main.go#L16-L24) are inferred from
[hashicorp/consul](https://github.com/hashicorp/consul)'s [template](https://github.com/hashicorp/consul/blob/main/.changelog/changelog.tmpl). Feel free to fork and adjust the options as needed.

[changelog-demo.webm](https://github.com/kisunji/gen-changelog/assets/30640057/ebc210ce-98a9-4478-8746-9243b984f0d7)

## Usage

### Prerequisites

You must have `gh` installed and authenticated. See https://cli.github.com/ for installation instructions.

### Installing

```sh
go install github.com/kisunji/gen-changelog@latest
```

### Usage

Currently, `gen-changelog` can only be run in a git repository root where it expects to find the `.changelog` directory.

You can optionally alias the `gen-changelog` to something like `gcl` for ease of typing.

## Contributing

I consider this complete enough to use in my daily workflow, but contributions are welcome.
