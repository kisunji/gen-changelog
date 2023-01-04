# gen-changelog

`gen-changelog` is a wrapper around the GitHub CLI [`gh`](https://cli.github.com/) which 
automatically generates boilerplate for [`go-changelog`](https://github.com/hashicorp/go-changelog) 
so you can easily write the changelog body without fiddling with files and figuring out PR numbers.


[gen-changelog-demo.webm](https://user-images.githubusercontent.com/30640057/210468619-dc9374e8-541b-43c1-a245-587b0a69fa99.webm)

## Installation

### Prerequisites

You must have `gh` installed and authenticated. See https://cli.github.com/ for installation instructions.

```sh
go install github.com/kisunji/gen-changelog@latest
```

## Contributing

I consider this complete enough to use in my daily workflow, but contributions are welcome.