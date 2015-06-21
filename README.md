# Linkey

**Link checker for BBC News & World Services sites.**

The idea is to quickly check a page for broken links by doing a status check on all the relative URL's on the page.

## Installation

    go get github.com/daveblooman/linkey

## Usage

### Command Line

```sh
linkey smoke /path/to/config.yaml
```

**Examples**

```sh
linkey smoke config.yaml
```

**Output**

Once running, you'll see either a 200 with `Status is 200 for <URL>` or `Status is NOT GOOD for <URL>`.

### Config File

If you have a lot of URLs that you want to check all the time using from a file is an alternative option.  This will utilise the smoke option, then point to a YAML file with the extension.  In some situations, we are deploying applications that we don't want public facing, so ensuring they 404 is essential.  There is a status code option to allow a specific status code to be set against a group of URL's, ensuring builds fail if the right code conditions are met.

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
