# Podarc

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/c543edcdb7e94c67a1859d0147600629)](https://app.codacy.com/gh/sa7mon/podarc?utm_source=github.com&utm_medium=referral&utm_content=sa7mon/podarc&utm_campaign=Badge_Grade)

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