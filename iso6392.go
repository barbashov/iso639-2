package iso639_2

type Language struct {
	Alpha3  string
	Alpha2  string
	English []string
	French  []string
}

//go:generate go run cmd/generator.go -o lang-db.go

// FromCode looks up language for given ISO639-1 or ISO639-2 code
// Returns nil if not found
func FromCode(code string) *Language {
	if l, ok := Languages[code]; ok {
		return &l
	}
	return nil
}

// FromEnglishName looks up language for given english name
// Returns nil if not found
func FromEnglishName(name string) *Language {
	for _, l := range Languages {
		if sliceContainsString(l.English, name) {
			return &l
		}
	}
	return nil
}

// FromFrenchName looks up language for given french name
// Returns nil if not found
func FromFrenchName(name string) *Language {
	for _, l := range Languages {
		if sliceContainsString(l.French, name) {
			return &l
		}
	}
	return nil
}
