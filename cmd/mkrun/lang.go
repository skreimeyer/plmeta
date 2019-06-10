package main

import (
	"strings"
)

type lang struct {
	Compiled bool
	BuildCmd string
	RunCmd   string
	Path string
}

// helper function to keep a long hashmap out of our main function. This may be
// poor practice, since it would by extension rebuild our map every time that
// fetchlang is called.
func fetchLang(s string) lang {
	langMeta := map[string]lang{
		"PYTHON":lang{
			Compiled:false,
			RunCmd:"python",
		},
		"PERL_6":lang{
			Compiled:false,
			RunCmd:"perl6",
		},
		"PERL":lang{
			Compiled:false,
			RunCmd:"perl",
		},
		"RUBY":lang{
			Compiled:false,
			RunCmd:"ruby",
		},
	}
	return langMeta[s]
}

// slugify
func slugUp(s string) string {
	return strings.ToUpper(strings.ReplaceAll(s, " ", "_"))
}