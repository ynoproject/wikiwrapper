package common

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
}

type Author struct {
	Name         string `json:"name"`
	OriginalName string `json:"originalName,omitempty"`
}

type VendingMachine struct {
	Game     string   `json:"game"`
	Path     string   `json:"path"`
	MapId    string   `json:"mapId"`
	EventIds []string `json:"eventId"`
}
