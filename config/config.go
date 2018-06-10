package config

import (
	"regexp"
)

type Config struct {
	Feeds    []*Feed
	Patterns []*Pattern
}

type Pattern struct {
	Type    string
	Content string
	Regex   *regexp.Regexp
}

type Feed struct {
	Source string
	Type   string
	URL    string
}

func (c *Config) Compile() {
	for _, p := range c.Patterns {
		p.Regex = regexp.MustCompile(p.Content)
	}
}
