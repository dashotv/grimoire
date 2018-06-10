package parser

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"

	"github.com/dashotv/grimoire/config"
)

type Parser struct {
	cfg *config.Config
}

type Release struct {
	Raw         string
	Title       string
	Description string
	View        string
	Download    string
	Season      string
	Episode     string
	Size        string // bytes?
	GUID        string // infohash?
	Resolution  string
	Encoding    string
	Team        string
	Group       string
	Verified    bool
	Bluray      bool
	Uncensored  bool
	Checksum    string
	Source      string
	Type        string
	Published   *time.Time
}

func (r *Release) CalculateChecksum() {
	h := md5.New()
	h.Write([]byte(r.Download))
	r.Checksum = hex.EncodeToString(h.Sum(nil))
}

func NewParser(cfg *config.Config) *Parser {
	return &Parser{cfg: cfg}
}

func (p *Parser) Parse() []*Release {
	list := []*Release{}
	fp := gofeed.NewParser()
	for _, f := range p.cfg.Feeds {
		feed, err := fp.ParseURL(f.URL)
		if err != nil {
			fmt.Printf("failed to parse URL: %s\n", err)
			continue
		}

		for n, i := range feed.Items {
			r, err := p.ParseItem(f.Type, f.Source, i)
			fmt.Printf("%03d: %t: %s\n", n+1, err == nil, i.Title)
			if err != nil {
				//fmt.Printf("error: %s\n", err)
				continue
			}
			list = append(list, r)
		}
	}
	return list
}

func (p *Parser) ParseItem(t, source string, item *gofeed.Item) (*Release, error) {
	for _, r := range p.cfg.Patterns {
		if r.Type != t {
			continue
		}

		matches := r.Regex.FindAllStringSubmatch(item.Title, -1)
		if matches == nil {
			continue
		}

		fmt.Printf("matches: %#v\n", matches)
		r := &Release{}
		r.Raw = item.Title
		r.GUID = item.GUID
		r.Team = matches[0][2]
		r.Title = strings.Replace(matches[0][3], ".", " ", -1)
		r.Season = matches[0][5]
		r.Episode = matches[0][7]
		r.Resolution = matches[0][9]
		r.Group = matches[0][10]
		r.Published = item.PublishedParsed
		r.Download = item.Link
		r.Source = source
		r.Type = t

		//fmt.Printf("%#v\n", item.Extensions["newznab"]["attr"])
		for _, e := range item.Extensions["newznab"]["attr"] {
			switch e.Attrs["name"] {
			case "size":
				r.Size = e.Attrs["value"]
			//case "season":
			//	r.Season = e.Attrs["value"]
			//case "episode":
			//	r.Episode = e.Attrs["value"]
			case "showtitle":
				r.Title = e.Attrs["value"]
			case "group":
				r.Group = e.Attrs["value"]
			}
		}
		return r, nil
	}

	return nil, fmt.Errorf("no match: '%s'", item.Title)
}
