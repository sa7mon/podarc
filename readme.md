# Podarc

A simple podcast archiver

## Usage

```text
podarc --feedUrl https://provider.here/show.xml
```

```text
 -feedUrl string
        URL of podcast feed to archive. (Required)
  -outputDir string
        Directory to save the files into. (Required)
  -overwrite
        Overwrite episodes already downloaded. Default: false
```

## Development

Run tests:

```shell
go test tests/providers/
go test tests/utils/
```