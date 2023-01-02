package common

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"cgt.name/pkg/go-mwclient"
	"cgt.name/pkg/go-mwclient/params"
	"github.com/antonholmquist/jason"
)

var gameNames = map[string]string{
	"yume":                  "Yume Nikki",
	"2kki":                  "Yume 2kki",
	"flow":                  "Dotflow",
	"someday":               "Someday",
	"deepdreams":            "Deep Dreams",
	"prayers":               "Answered Prayers",
	"amillusion":            "Amillusion",
	"unevendream":           "Uneven Dream",
	"braingirl":             "Braingirl",
	"collectiveunconscious": "Collective Unconscious",
	"muma":                  "Muma Rope",
	"genie":                 "Dream Genie",
}

func createClient() (client *mwclient.Client, err error) {
	client, err = mwclient.New("https://yume.wiki/api.php", "yumeWikiAPIBot")
	return client, err
}

func fetchAllResultsFromSmwQuery(smwQuery *SmwQuery) (results []*jason.Object, err error) {
	for smwQuery.Next() {
		result := smwQuery.Resp()
		fetchedResults, err := result.GetObjectArray("query", "results")

		if err != nil {
			return results, err
		}

		results = append(results, fetchedResults...)
	}

	if smwQuery.Err() != nil {
		return results, smwQuery.Err()
	}

	return results, err
}

func GetLocations(gameCode string) (locations []*Location, err error) {
	gameName, ok := gameNames[gameCode]
	if !ok {
		return locations, errors.New("game not supported")
	}

	client, err := createClient()
	if err != nil {
		return locations, err
	}

	condition := fmt.Sprintf("Category:%s Locations", gameName)
	printouts := []string{"Has location image", "Header background color", "Header font color", "Has primary author", "Has contributing author", "Japanese name", "Has BGM", "Has location map", "Version added", "Versions updated", "Version removed", "Version gaps"}

	parameters := params.Values{
		"conditions":  condition,
		"printouts":   strings.Join(printouts, "|"),
		"parameters":  "limit=500",
		"format":      "json",
		"api_version": "3",
	}
	results := NewSmwQuery(client, parameters)

	locationsToProcess, err := fetchAllResultsFromSmwQuery(results)

	if err != nil {
		return locations, err
	}

	for _, locationToProcess := range locationsToProcess {
		if err != nil {
			return locations, err
		}

		for _, value := range locationToProcess.Map() {
			value, err := value.Object()
			if err != nil {
				return locations, err
			}

			location, err := processLocation(gameName, value)
			if err != nil {
				return locations, err
			}

			locations = append(locations, location)
		}
	}

	return locations, err
}

func GetConnections(gameCode string) (connections []*Connection, err error) {
	gameName, ok := gameNames[gameCode]
	if !ok {
		return connections, errors.New("game not supported")
	}

	client, err := createClient()
	if err != nil {
		return connections, err
	}

	conditions := []string{fmt.Sprintf("%s:+", gameName), "Is subobject type::connection"}
	printouts := []string{"Connection/Origin", "Connection/Location", "Connection/Attribute", "Connection/Unlock Conditions"}

	parameters := params.Values{
		"conditions":  strings.Join(conditions, "|"),
		"printouts":   strings.Join(printouts, "|"),
		"parameters":  "limit=500",
		"format":      "json",
		"api_version": "3",
	}

	results := NewSmwQuery(client, parameters)

	connectionsToProcess, err := fetchAllResultsFromSmwQuery(results)
	if err != nil {
		return connections, err
	}

	for _, connectionToProcess := range connectionsToProcess {
		if err != nil {
			return connections, err
		}

		for _, value := range connectionToProcess.Map() {
			value, err := value.Object()
			if err != nil {
				return connections, err
			}

			connection, err := processConnection(gameCode, value)
			if err != nil {
				return connections, err
			}

			connections = append(connections, connection)
		}
	}

	return connections, err
}

func GetAuthors(gameCode string) (authors []*Author, err error) {
	gameName, ok := gameNames[gameCode]
	if !ok {
		return authors, errors.New("game not supported")
	}

	client, err := createClient()
	if err != nil {
		return authors, err
	}

	conditions := fmt.Sprintf("-Has subobject::%s:Authors", gameName)
	printouts := []string{"Author/Name", "Author/Original Name"}
	queryParams := []string{"sort=Author/Date Joined", "order=asc", "limit=500"}

	parameters := params.Values{
		"conditions":  conditions,
		"printouts":   strings.Join(printouts, "|"),
		"parameters":  strings.Join(queryParams, "|"),
		"format":      "json",
		"api_version": "3",
	}

	results := NewSmwQuery(client, parameters)
	authorsToProcess, err := fetchAllResultsFromSmwQuery(results)
	if err != nil {
		return authors, err
	}

	for _, authorToProcess := range authorsToProcess {
		if err != nil {
			return authors, err
		}

		for _, value := range authorToProcess.Map() {
			value, err := value.Object()
			if err != nil {
				return authors, err
			}

			author, err := processAuthor(gameCode, value)
			if err != nil {
				return authors, err
			}

			authors = append(authors, author)
		}
	}
	return authors, err
}

func GetMaps(gameCode string, locationTitle string) (locationMaps []*LocationMap, err error) {
	locationMaps = []*LocationMap{}
	gameName, ok := gameNames[gameCode]
	if !ok {
		return locationMaps, errors.New("game not supported")
	}

	client, err := createClient()
	if err != nil {
		return locationMaps, err
	}

	conditions := fmt.Sprintf("%s:%s", gameName, locationTitle)
	printouts := "Has location map"

	parameters := params.Values{
		"conditions": conditions,
		"printouts": printouts,
		"format": "json",
		"api_version": "3",
	}

	results := NewSmwQuery(client, parameters)
	locationsToProcess, err := fetchAllResultsFromSmwQuery(results)
	if err != nil {
		return locationMaps, err
	}

	for _, locationToProcess := range locationsToProcess {
		if err != nil {
			return locationMaps, err
		}

		for _, value := range locationToProcess.Map() {
			value, err := value.Object()
			if err != nil {
				return locationMaps, err
			}

			printouts, err := value.GetObject("printouts")
			if err != nil {
				return nil, err
			}

			locationMapObjects, err := printouts.GetObjectArray("Has location map")
			if err != nil {
				log.Print("SERVER", "locationMaps", err.Error())
				return nil, err
			}

			maps, err := processLocationMaps(locationMapObjects)
			if err != nil {
				log.Print("SERVER", "locationMaps", err.Error())
				return nil, err
			}

			locationMaps = maps
		}
	}
	return locationMaps, err
}

func processLocation(gameCode string, value *jason.Object) (location *Location, err error) {
	printouts, err := value.GetObject("printouts")
	if err != nil {
		return nil, err
	}

	fulltext, err := value.GetString("fulltext")
	if err != nil {
		log.Print("SERVER", "fulltext", err.Error())
		return nil, err
	}

	title := strings.Split(fulltext, ":")[1]

	location = &Location{
		Title:           title,
		Game:            gameCode,
		BGMs:            []*BGM{},
		LocationMaps:    []*LocationMap{},
		VersionsUpdated: []string{},
	}

	locationImage, err := printouts.GetStringArray("Has location image")
	if err != nil {
		log.Print("SERVER", "locationImage", err.Error())
		return nil, err
	}

	if len(locationImage) > 0 {
		location.LocationImage = locationImage[0]
	}

	headerBackgroundColor, err := printouts.GetStringArray("Header background color")
	if err != nil {
		log.Print("SERVER", "headerBackgroundColor", err.Error())
		return nil, err
	}

	if len(headerBackgroundColor) > 0 {
		location.BackgroundColor = headerBackgroundColor[0]
	}

	headerFontColor, err := printouts.GetStringArray("Header font color")
	if err != nil {
		log.Print("SERVER", "headerFontColor", err.Error())
		return nil, err
	}

	if len(headerFontColor) > 0 {
		location.FontColor = headerFontColor[0]
	}

	primaryAuthor, err := printouts.GetStringArray("Has primary author")
	if err != nil {
		log.Print("SERVER", "primaryAuthor", err.Error())
		return nil, err
	}

	if len(primaryAuthor) > 0 {
		location.PrimaryAuthor = primaryAuthor[0]
	}

	contributingAuthors, err := printouts.GetStringArray("Has contributing author")
	if err != nil {
		log.Print("SERVER", "contributingAuthors", err.Error())
		return nil, err
	}
	location.ContributingAuthors = contributingAuthors

	japaneseName, err := printouts.GetStringArray("Japanese name")
	if err != nil {
		log.Print("SERVER", "japaneseName", err.Error())
		return nil, err
	}

	if len(japaneseName) > 0 {
		location.OriginalName = japaneseName[0]
	}

	versionAdded, err := printouts.GetStringArray("Version added")
	if err != nil {
		log.Print("SERVER", "versionAdded", err.Error())
		return nil, err
	}

	if len(versionAdded) > 0 {
		location.VersionAdded = versionAdded[0]
	}

	versionsUpdated, err := printouts.GetStringArray("Versions updated")
	if err != nil {
		log.Print("SERVER", "versionUpdated", err.Error())
		return nil, err
	}

	location.VersionsUpdated = versionsUpdated

	versionRemoved, err := printouts.GetStringArray("Version removed")
	if err != nil {
		log.Print("SERVER", "versionRemoved", err.Error())
		return nil, err
	}

	if len(versionRemoved) > 0 {
		location.VersionRemoved = versionRemoved[0]
	}

	versionGaps, err := printouts.GetStringArray("Version gaps")
	if err != nil {
		log.Print("SERVER", "versionGaps", err.Error())
		return nil, err
	}

	location.VersionGaps = versionGaps

	bgmObjects, err := printouts.GetObjectArray("Has BGM")
	if err != nil {
		log.Print("SERVER", "bgms", err.Error())
		return nil, err
	}

	bgms, err := processBGMs(bgmObjects)
	if err != nil {
		log.Print("SERVER", "bgms", err.Error())
		return nil, err
	}

	if len(bgms) > 0 {
		location.BGMs = bgms
	}

	locationMapObjects, err := printouts.GetObjectArray("Has location map")
	if err != nil {
		log.Print("SERVER", "locationMaps", err.Error())
		return nil, err
	}

	locationMaps, err := processLocationMaps(locationMapObjects)
	if err != nil {
		log.Print("SERVER", "locationMaps", err.Error())
		return nil, err
	}

	if len(locationMaps) > 0 {
		location.LocationMaps = locationMaps
	}

	return location, err
}

func processBGMs(bgmObjects []*jason.Object) (bgms []*BGM, err error) {
	bgms = []*BGM{}
	for _, bgm := range bgmObjects {
		var bgmPath string
		var bgmTitle string
		var bgmLabel string
		path, err := bgm.GetStringArray("Has media path", "item")
		if err != nil {
			return bgms, err
		}

		if len(path) > 0 {
			bgmPath = path[0]
		}

		title, err := bgm.GetStringArray("BGM/Title", "item")
		if err != nil {
			return bgms, err
		}

		if len(title) > 0 {
			bgmTitle = title[0]
		}
		label, err := bgm.GetStringArray("BGM/Label", "item")
		if err != nil {
			return bgms, err
		}

		if len(label) > 0 {
			bgmLabel = label[0]
		}
		bgms = append(bgms, &BGM{
			Path:  bgmPath,
			Title: bgmTitle,
			Label: bgmLabel,
		})
	}
	return bgms, nil
}

func processLocationMaps(locationMapObjects []*jason.Object) (locationMaps []*LocationMap, err error) {
	locationMaps = []*LocationMap{}
	for _, locationMapObject := range locationMapObjects {
		var locationMapPath string
		var locationMapCaption string
		path, err := locationMapObject.GetStringArray("Has image path", "item")
		if err != nil {
			return locationMaps, err
		}
		if len(path) > 0 {
			locationMapPath = path[0]
		}

		caption, err := locationMapObject.GetStringArray("Location Map/Caption", "item")
		if err != nil {
			return locationMaps, err
		}
		if len(caption) > 0 {
			locationMapCaption = caption[0]
		}

		locationMaps = append(locationMaps, &LocationMap{
			Path:    locationMapPath,
			Caption: locationMapCaption,
		})
	}
	return locationMaps, nil
}

func processConnection(gameCode string, value *jason.Object) (connection *Connection, err error) {
	printouts, err := value.GetObject("printouts")

	if err != nil {
		return connection, err
	}

	var origin string
	var destination string

	connectionOrigin, err := printouts.GetObjectArray("Connection/Origin")
	if err != nil {
		log.Print("SERVER", "origin", err.Error())
		return connection, err
	}

	if len(connectionOrigin) > 0 {
		originText, err := connectionOrigin[0].GetString("fulltext")
		if err != nil {
			log.Print("SERVER", "origin", err.Error())
			return connection, err
		}
		origin = strings.Split(originText, ":")[1]
	}

	connectionDestination, err := printouts.GetObjectArray("Connection/Location")
	if err != nil {
		return connection, err
	}

	if len(connectionDestination) > 0 {
		destinationText, err := connectionDestination[0].GetString("fulltext")
		if err != nil {
			log.Print("SERVER", "destination", err.Error())
			return connection, err
		}
		destination = strings.Split(destinationText, ":")[1]
	}

	attributes, err := printouts.GetStringArray("Connection/Attribute")
	if err != nil {
		log.Print("SERVER", "attributes", err.Error())
		return connection, err
	}

	unlockConditions, err := printouts.GetStringArray("Connection/Unlock Conditions")
	if err != nil {
		log.Print("SERVER", "unlockConditions", err.Error())
		return connection, err
	}

	connection = &Connection{
		Game:        gameCode,
		Origin:      origin,
		Destination: destination,
		Attributes:  attributes,
	}

	if len(unlockConditions) > 0 {
		connection.UnlockConditions = unlockConditions[0]
	}

	return connection, err
}

func processAuthor(gameCode string, value *jason.Object) (author *Author, err error) {
	author = &Author{}
	printouts, err := value.GetObject("printouts")
	if err != nil {
		return nil, err
	}

	authorName, err := printouts.GetStringArray("Author/Name")
	if err != nil {
		log.Print("SERVER", "authorName", err.Error())
		return nil, err
	}

	if len(authorName) > 0 {
		author.Name = authorName[0]
	}

	originalNameObject, err := printouts.GetObjectArray("Author/Original Name")
	if err != nil {
		log.Print("SERVER", "originalNameObject", err.Error())
		return nil, err
	}

	if len(originalNameObject) > 0 {
		originalName, err := originalNameObject[0].GetStringArray("Text", "item")

		if err != nil {
			log.Print("SERVER", "originalName", err.Error())
			return nil, err
		}

		if len(originalName) > 0 {
			author.OriginalName = originalName[0]
		}
	}

	return author, err
}
