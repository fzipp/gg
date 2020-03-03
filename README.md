# gg

A set of command line tools and [Go](https://golang.org) packages to work with
[Thimbleweed Park](https://thimbleweedpark.com/) data files.

The project name "gg" was chosen, because the names of some of these formats
start with those two letters, e.g. "ggpack" or "GGDictionary". They were
conceived by [Grumpy Gamer](https://grumpygamer.com/) (Ron Gilbert).
This project is not related to him or Terrible Toybox, Inc.

## Command line tools

* [ggpack](https://pkg.go.dev/github.com/fzipp/gg/cmd/ggpack) A tool to inspect, unpack or create "ggpack" files.
* [ggdict](https://pkg.go.dev/github.com/fzipp/gg/cmd/ggdict) A tool to convert back and forth between the GGDictionary format and JSON.
* [retext](https://pkg.go.dev/github.com/fzipp/gg/cmd/retext) A tool to replace ID placeholders like @12345 in files with texts from a text table file in TSV format.
* [nutfmt](https://pkg.go.dev/github.com/fzipp/gg/cmd/nutfmt) A tool to indent [Squirrel](http://squirrel-lang.org/) script files.
* [yack](https://pkg.go.dev/github.com/fzipp/gg/cmd/yack) A tool to run yack dialogs.

### Installation

Either download binaries for your operating system from the [latest release](https://github.com/fzipp/gg/releases/latest) or build from source with Go:

```
go get github.com/fzipp/gg/cmd/ggpack
go get github.com/fzipp/gg/cmd/ggdict
go get github.com/fzipp/gg/cmd/retext
go get github.com/fzipp/gg/cmd/nutfmt
go get github.com/fzipp/gg/cmd/yack
```

## Go packages

* [ggpack](https://pkg.go.dev/github.com/fzipp/gg/ggpack) Read and write ggpack files.
* [ggdict](https://pkg.go.dev/github.com/fzipp/gg/ggdict) Read and write the GGDictionary format.
* [texts](https://pkg.go.dev/github.com/fzipp/gg/texts) Replace text ID placeholders with texts.
* [yack](https://pkg.go.dev/github.com/fzipp/gg/yack) Parse and run yack dialog files.

Related, but independent packages:

* [texturepacker](https://pkg.go.dev/github.com/fzipp/texturepacker) Read sprite sheet information from [TexturePacker](https://www.codeandweb.com/texturepacker)'s JSON (Hash) export format.
* [bmfont](https://pkg.go.dev/github.com/fzipp/bmfont) Load and render bitmap fonts in the format of [AngelCode's bitmap font generator](https://www.angelcode.com/products/bmfont/).

## Related Work

Projects by other people with similar objectives:

* [NGGPack](https://github.com/scemino/NGGPack)
  .NET based tool for reading and writing ggpack archives.
* [ggdump](https://github.com/mstr-/twp-ggdump)
  Python based tool for listing and extracting files from ggpack archives.
* [r2-ggpack](https://github.com/mrmacete/r2-ggpack)
  Radare2 plugins to manipulate ggpack archives.
* [engge](https://github.com/scemino/engge)
  Experimental game engine for Thimbleweed Park by the same author as NGGPack.
* [Thimbleweed Park Explorer](https://github.com/bgbennyboy/Thimbleweed-Park-Explorer)
  An explorer/viewer/dumper tool for Thimbleweed Park
* [ggpack](https://github.com/s-l-teichmann/ggpack)
  Command line tool for inspecting ggpack files, written in Go.

## License

This project is free and open source software licensed under the
[BSD 3-Clause License](LICENSE).
