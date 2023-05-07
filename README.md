# gomdtoc
CLI program to generate TOC(table of content) for markdown notes directory, build with [Cobra](https://github.com/spf13/cobra)
![](https://political-capable-roll.glitch.me/get/@dreamjz-gomdtoc?theme=rule34)
## 1. Example

- [Example](./example)

  ```sh
  $ gomdtoc ./example -r -s sub3-skip 
  ```

- [Existed Repo](https://github.com/dreamjz/my-notes)
## 2. Installation

### Install with `go`

```sh
$ go install "github.com/dreamjz/gomdtoc@latest"
```

## Binary file

Download the binary file from [Release](https://github.com/dreamjz/gomdtoc/releases) 

## 3. Usage

```sh
$ gomdtoc directory_path [flags]
```

```sh
$ gomdtoc -h
CLI program to generate toc for markdown notes directory

Usage:
  gomdtoc [flags]

Flags:
  -h, --help           help for gomdtoc
  -r, --recursive      --recursive; generate TOC file for every sub-directory
  -s, --skip strings   --skip dir_name1,dir_name2, ...; skip specified directories
  -t, --title string   --title title_field, specify the title field in frontmatter  (default "title")
```





