package cnf

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"opg-reports/shared/dates"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const billingDay int = 13
const errUrlSubNotFound string = "Pattern [%s] is not supported"

// ApiSubFunc type used for signature mapping
// key - the section of the url to be replaced: {month} | {month:-2}
// fragment - would be the -2 in {month:-2}
// url - the original url
type ApiSubFunc func(key string, fragment string, url string, d *time.Time) string

func month(key string, fragment string, url string, d *time.Time) (m string) {
	m = url
	now := *d
	date := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	if fragment == "" {
		m = strings.ReplaceAll(url, key, date.Format(dates.FormatYM))
		return
	}

	if i, err := strconv.Atoi(fragment); err == nil {
		date = date.AddDate(0, i, 0)
		m = strings.ReplaceAll(url, key, date.Format(dates.FormatYM))
	}
	slog.Debug("month",
		slog.String("key", key),
		slog.String("fragment", fragment),
		slog.String("url", url),
		slog.String("m", m),
		slog.String("date", date.String()))
	return
}

func billingMonth(key string, fragment string, url string, d *time.Time) (m string) {
	var date time.Time
	m = url
	now := *d
	date = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	if now.Day() < billingDay {
		date = date.AddDate(0, -2, 0)
	} else {
		date = date.AddDate(0, -1, 0)
	}

	slog.Debug("billingMonth",
		slog.String("key", key),
		slog.String("fragment", fragment),
		slog.String("url", url),
		slog.String("m", m),
		slog.String("now", now.String()),
		slog.String("date", date.String()))
	m = strings.ReplaceAll(url, key, date.Format(dates.FormatYM))
	return
}

var ApiSubstitutionFuncs = map[string]ApiSubFunc{
	"month":        month,
	"billingMonth": billingMonth,
}

type RepoStandards struct {
	Baseline    []string `json:"baseline"`
	Extended    []string `json:"extended"`
	Information []string `json:"information"`
}
type StandardsCnf struct {
	Repository RepoStandards `json:"repository"`
}

type SiteSection struct {
	Name     string         `json:"name"`
	Href     string         `json:"href"`
	Header   bool           `json:"header"`
	Sections []*SiteSection `json:"sections"`

	Exclude bool `json:"exclude"`

	Api          map[string]string `json:"api"`
	TemplateName string            `json:"template"`

	Registered bool `json:"-"`
}

func (s *SiteSection) ClassName() string {
	str := "sect-"
	str = str + strings.ToLower(s.Name)
	str = strings.ReplaceAll(str, " ", "-")
	return str
}

func (s *SiteSection) ApiUrls() (res map[string]string, err error) {
	var re = regexp.MustCompile(`(?mi){(.*?)}`)
	res = map[string]string{}
	now := time.Now().UTC()

	for name, url := range s.Api {
		org := url
		for _, match := range re.FindAllString(url, -1) {
			index, fragment := getIndex(match)
			if subFunc, ok := ApiSubstitutionFuncs[index]; ok {
				url = subFunc(match, fragment, url, &now)
			} else {
				err = fmt.Errorf(errUrlSubNotFound, match)
			}
		}
		slog.Debug("[cnf] api url processing",
			slog.String("name", name),
			slog.String("original_url", org),
			slog.String("processed_url", url))
		res[name] = url
	}
	return
}

func getIndex(match string) (index string, fragment string) {
	fragment = ""
	index = strings.ReplaceAll(strings.ReplaceAll(match, "}", ""), "{", "")
	if strings.Contains(index, ":") {
		sp := strings.Split(index, ":")
		index = sp[0]
		fragment = sp[1]
	}
	return
}

type Config struct {
	Organisation string         `json:"organisation"`
	Sections     []*SiteSection `json:"sections"`
	// Standards is used ot allow customisable fields for baseline repo standards
	Standards *StandardsCnf `json:"standards"`
}

func Load(content []byte) (*Config, error) {
	cfg := &Config{}
	err := json.Unmarshal(content, &cfg)
	return cfg, err
}

func FlatSections(sects []*SiteSection, flat map[string]*SiteSection) {

	for _, sect := range sects {
		flat[sect.Href] = sect
		if len(sect.Sections) > 0 {
			FlatSections(sect.Sections, flat)
		}
	}
	return
}
