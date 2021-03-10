# ISO 639-2

![Tests badge](https://github.com/barbashov/iso639-2/actions/workflows/test.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/barbashov/iso639-2)](https://goreportcard.com/report/github.com/barbashov/iso639-2)


A database of ISO 639-2 and ISO 639-1 languages.

Generated from official ISO 639-2 list, so no native names, unfortunately :(

# Motivation

There's an excellent [Go library for ISO 639-1](https://github.com/emvi/iso-639-1), but it lacks ISO 639-2 codes.

# Data source

Database is generated (see `cmd/generator.go`) from official ISO 639-2 data. See [The Library of Congress website](https://www.loc.gov/standards/iso639-2/) for details.

# Installation

```
go get github.com/barbashov/iso639-2
```

# Examples

```go
iso639_2.Languages // returns languages lookup table

iso639_2.FromCode("eng") // returns object representing English language
iso639_2.FromEnglishName("English") // returns object representing English language
iso639_2.FromFrenchName("anglais") // returns object representing English language
```

# Contribute

Feel free to open issues send pull requests.

# License

MIT