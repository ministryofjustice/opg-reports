// Package endpoiints is a string with magic values - `{NAME}` - that
// is converted into a uri
//
// This allows the string to contain place holders for things like
// the current date or version change. It also allows for modifiers
// on those values
//
// Example:
//
//	`/{version}/model/{month:-1}/{month:0}`
//	`/{version}/model/{month:-2,2024-01-01}/{month:2}`
//
// Currently supports the following magic values:
//
//	`{year}`
//	`{month}`
//	`{day}`
//	`{billing_date}`
//	`{version}`
//
// See the variable `parsers` for the mapping between keyword
// and function
package endpoints

import (
	"log/slog"
	"net/http"
	"regexp"
	"strings"
)

type parserGroup struct {
	Name      string
	Original  string
	Arguments []string
}

type parserFunc func(uri string, pg *parserGroup, request *http.Request) string

// pattern matches anything in between {}
//   - /test/{word:-1}/bar => {word:-1}
var pattern string = `(?mi){(.*?)}`

// parsers maps the keyword in string to funcs that process them
var parsers map[string]parserFunc = map[string]parserFunc{
	"year":         year,
	"month":        month,
	"day":          day,
	"billing_date": billingMonth,
	"version":      version,
}

type ApiEndpoint string

func (self ApiEndpoint) String() string {
	return string(self)
}

// parserGroups processes the string version of self and
// returns a list of elements within the string that need
// to be replaced and details about them
// Uses regex pattern to find them `{(.*?)}`
//   - /{will-be-replaced}/fixed-value/
//   - /{day:-10}/fixed-value/
//   - what-is-{your name}
func (self ApiEndpoint) parserGroups() (found []*parserGroup) {
	var (
		source  string = string(self)
		re             = regexp.MustCompile(pattern)
		matches        = re.FindAllString(source, -1)
	)

	// go over each match and get arguments
	for _, match := range matches {
		var pg = &parserGroup{Original: match}
		var trimmed = strings.Trim(strings.Trim(match, "}"), "{")
		var sp = strings.Split(trimmed, ":")
		// replace any spaces in the name with underscore
		pg.Name = strings.ReplaceAll(sp[0], " ", "_")
		// no look for any arguments to the function name
		if len(sp) > 1 {
			var str = strings.ReplaceAll(sp[1], " ", "")
			var args = strings.Split(str, ",")
			for _, arg := range args {
				if arg != "" {
					pg.Arguments = append(pg.Arguments, arg)
				}
			}
		}
		found = append(found, pg)
	}
	return
}

// Parse replaces all magic values - {NAME:x,y,z} - with their resolved values
// and therefore creates the uri to call the api with
func (self ApiEndpoint) Parse(request *http.Request) (u string) {
	u = string(self)

	var groups = self.parserGroups()

	for _, pg := range groups {
		if parser, ok := parsers[pg.Name]; ok {
			u = parser(u, pg, request)
		} else {
			slog.Error("unknown url group match:", slog.String("name", pg.Name), slog.String("url", string(self)))
		}
	}

	return
}
