# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.6.0] - 2022-09-26
### Added
- Initial support for Return to Monkey Island ggpack files

### Changed
- ggpack: implement the io/fs.FS interface introduced in Go 1.16

## [0.5.0] - 2020-10-19
### Added
- Savegame support (new `savegame` package, new `ggsavegame` command)

## [0.4.0] - 2020-09-06
### Fixed
- Fix encoding of .bnut files in ggpack

## [0.3.0] - 2020-03-03
### Added
- New tool to run yack dialogs

## [0.2.0] - 2020-02-24
### Fixed
- Fix bnut file de-/encryption

## [0.1.0] - 2020-02-15
### Added
- First release of the command line tools as binaries for Linux, macOS and Windows
