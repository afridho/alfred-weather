package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jason0x43/go-alfred"
)

// HourlyCommand gets a weather forecast
type HourlyCommand struct{}

// About returns information about a command
func (c HourlyCommand) About() alfred.CommandDef {
	return alfred.CommandDef{
		Keyword:     "hourly",
		Description: "Get a forecast for the next few hours",
		IsEnabled:   true,
	}
}

// Items returns the items for the command
func (c HourlyCommand) Items(arg, data string) (items []alfred.Item, err error) {
	dlog.Printf("Running HourlyCommand")

	var cfg hourlyConfig
	if data != "" {
		if err := json.Unmarshal([]byte(data), &cfg); err != nil {
			dlog.Printf("Invalid hourly config")
		}
	}

	var weather Weather
	var loc Location
	if loc, weather, err = getWeather(arg); err != nil {
		return
	}

	var startTime time.Time
	if cfg.Start != nil {
		startTime = *cfg.Start
	} else if len(weather.Hourly) > 0 {
		startTime = weather.Hourly[0].Time
	}

	heading := alfred.Item{
		Title: "Weather for " + loc.Name,
		// Subtitle: alfred.Line,
		Subtitle: loc.Name,
		Arg: &alfred.ItemArg{
			Keyword: "daily",
		},
	}

	if weather.URL != "" {
		heading.AddMod(alfred.ModCmd, alfred.ItemMod{
			Subtitle: "Open this forecast in a browser",
			Arg: &alfred.ItemArg{
				Keyword: "daily",
				Mode:    alfred.ModeDo,
				Data:    alfred.Stringify(&dailyCfg{ToOpen: weather.URL}),
			},
		})
	}

	items = append(items, heading)

	deg := "F"
	if config.Units == unitsMetric {
		deg = "C"
	}

	addAlertItems(&weather, &items)

	for _, entry := range weather.Hourly {
		if entry.Time.Before(startTime) {
			continue
		}

		conditions := entry.Summary
		icon := entry.Icon

		item := alfred.Item{
			Title:    entry.Time.Format("Mon "+config.TimeFormat) + ": " + conditions,
			Subtitle: fmt.Sprintf("%d°%s (%d°%s)   ☂ %d%%", entry.Temp.Int64(), deg, entry.ApparentTemp.Int64(), deg, entry.Precip),
			Icon:     getIconFile(icon),
		}

		items = append(items, item)
	}

	return
}

type hourlyConfig struct {
	Start *time.Time
}
