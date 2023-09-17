package users

import (
	"golang.org/x/net/webdav"
	"regexp"
	"strings"
)

// Rule is a disallow/allow rule.
type Rule struct {
	Regex  bool           `yaml:"regex"`
	Allow  bool           `yaml:"allow"`
	Modify bool           `yaml:"modify"`
	Path   string         `yaml:"path"`
	Regexp *regexp.Regexp `yaml:"regexp"`
}

// User contains the settings of each user.
type User struct {
	Username string  `yaml:"username"`
	Password string  `yaml:"password"`
	Scope    string  `yaml:"scope"`
	Modify   bool    `yaml:"modify"`
	Rules    []*Rule `yaml:"rules"`

	Handler *webdav.Handler
}

// ACL checks if the user has permission to access a directory/file
func (u User) ACL(url string, noModification bool) bool {
	for _, rule := range u.Rules {
		isAllowed := rule.Allow && (noModification || rule.Modify)
		if rule.Regex {
			if rule.Regexp.MatchString(url) {
				return isAllowed
			}
		}

		if strings.HasPrefix(url, rule.Path) {
			return isAllowed
		}
	}

	return noModification || u.Modify
}
