package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mmcdole/gofeed"

	"github.com/dashotv/grimoire/config"
	"github.com/dashotv/server/models"
)

type Parser struct {
	cfg *config.Config
}

func NewParser(cfg *config.Config) *Parser {
	return &Parser{cfg: cfg}
}

func (p *Parser) Parse() []*models.Release {
	list := []*models.Release{}
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

func (p *Parser) ParseItem(t, source string, item *gofeed.Item) (*models.Release, error) {
	for _, r := range p.cfg.Patterns {
		if r.Type != t {
			continue
		}

		matches := r.Regex.FindAllStringSubmatch(item.Title, -1)
		if matches == nil {
			continue
		}

		fmt.Printf("matches: %#v\n", matches)
		r := &models.Release{}
		r.Raw = item.Title
		r.Guid = item.GUID
		r.Team = matches[0][2]
		r.Title = strings.Replace(matches[0][3], ".", " ", -1)
		r.Season, _ = strconv.Atoi(matches[0][5])
		r.Episode, _ = strconv.Atoi(matches[0][7])
		r.Resolution, _ = strconv.Atoi(matches[0][9])
		r.Team = matches[0][10]
		r.Published = *item.PublishedParsed
		r.Download = item.Link
		r.Source = source
		r.Type = t

		//fmt.Printf("%#v\n", item.Extensions["newznab"]["attr"])
		for _, e := range item.Extensions["newznab"]["attr"] {
			switch e.Attrs["name"] {
			case "size":
				r.Size, _ = strconv.Atoi(e.Attrs["value"])
			//case "season":
			//	r.Season = e.Attrs["value"]
			//case "episode":
			//	r.Episode = e.Attrs["value"]
			case "showtitle":
				r.Title = e.Attrs["value"]
			case "group":
				r.Team = e.Attrs["value"]
			}
		}
		return r, nil
	}

	return nil, fmt.Errorf("no match: '%s'", item.Title)
}
