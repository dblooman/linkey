# Linkey

**Status checker for your websites**

The idea is to quickly check a page for broken links by doing a status check on all the relative URL's on the page.

## Install

To install, use `go get`:

```bash
$ go get -d github.com/DaveBlooman/linkey
```

## Usage

### Command Line

```sh
linkey check /path/to/config.yaml
```

**Examples**

```sh
linkey check config.yaml
```

**Output**

Once running, you'll see either a 200 with `Status is 200 for <URL>` or `Status is NOT GOOD for <URL>`.

### Config File

In some situations, you may be deploying applications that you don't want to be public facing, so ensuring they don't 200 is essential.  There is a status code option to allow a specific status code to be set against a group of URL's, ensuring builds fail if the right code conditions are met.

Example YAML Config:

```yaml
base: 'http://www.bbc.co.uk'

headers:
  -
   key: 'X-content-override'
   value: 'https://example.com'

statuscode: 200

paths:
  - /news
  - /news/uk

```

## Contribution

1. Fork ([https://github.com/DaveBlooman/linkey-go/fork](https://github.com/DaveBlooman/linkey-go/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[DaveBlooman](https://github.com/DaveBlooman)
