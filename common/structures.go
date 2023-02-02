package common

import "encoding/json"

type Location struct {
	Title               string         `json:"title"`
	LocationImage       string         `json:"locationImage"`
	Game                string         `json:"game"`
	BackgroundColor     string         `json:"backgroundColor"`
	FontColor           string         `json:"fontColor"`
	OriginalName        string         `json:"originalName,omitempty"`
	BGMs                []*BGM         `json:"bgms"`
	LocationMaps        []*LocationMap `json:"locationMaps"`
	PrimaryAuthor       string         `json:"primaryAuthor,omitempty"`
	ContributingAuthors []string       `json:"contributingAuthors,omitempty"`
	VersionAdded        string         `json:"versionAdded"`
	VersionsUpdated     []string       `json:"versionsUpdated"`
	VersionRemoved      string         `json:"versionRemoved,omitempty"`
	VersionGaps         []string       `json:"versionGaps"`
}

type BGM struct {
	Path  string `json:"path"`
	Title string `json:"title"`
	Label string `json:"label,omitempty"`
}

type LocationMap struct {
	Path    string `json:"path"`
	Caption string `json:"caption"`
}

type Connection struct {
	Game              string   `json:"game"`
	Origin            string   `json:"origin"`
	Destination       string   `json:"destination"`
	Attributes        []string `json:"attributes"`
	UnlockConditions  string   `json:"unlockCondition,omitempty"`
	EffectsNeeded     []string `json:"effectsNeeded,omitempty"`
	SeasonAvailable   string   `json:"seasonAvailable,omitempty"`
	ChancePercentage  string   `json:"chancePercentage,omitempty"`
	ChanceDescription string   `json:"chanceDescription,omitempty"`
	IsRemoved         bool     `json:"isRemoved,omitempty"`
}

type Author struct {
	Name         string `json:"name"`
	OriginalName string `json:"originalName,omitempty"`
}

type VendingMachine struct {
	Game     string   `json:"game"`
	Path     string   `json:"path"`
	MapId    string   `json:"mapId"`
	EventIds []string `json:"eventIds"`
}

type Effect struct {
	Name           string   `json:"name"`
	OriginalName   string   `json:"originalName,omitempty"`
	AlternateNames []string `json:"alias,omitempty"`
	Location       string   `json:"location,omitempty"`
}

type MenuType struct {
	Name       string `json:"name"`
	Location   string `json:"location,omitempty"`
	Conditions string `json:"omitempty"`
}

type VersionHistory struct {
	VersionNumber string `json:"versionNumber"`
	CreatedBy     string `json:"createdBy"`
	CreatedAt     string `json:"createdAt"`
}

type LocationImage struct {
	Title  string   `json:"title"`
	Game   string   `json:"game"`
	Images []*Image `json:"images"`
}

type Image struct {
	Url    string      `json:"url"`
	Width  json.Number `json:"width"`
	Height json.Number `json:"height"`
}
