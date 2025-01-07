package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"cgt.name/pkg/go-mwclient"
	"cgt.name/pkg/go-mwclient/params"
	"github.com/antonholmquist/jason"
)

type GameParams struct {
	GameCode, Protag, ContinueKey string
}

var gameNames = map[string]string{
	"yume":        "Yume Nikki",
	"2kki":        "Yume 2kki",
	"flow":        "Dotflow",
	"someday":     "Someday",
	"deepdreams":  "Deep Dreams",
	"prayers":     "Answered Prayers",
	"amillusion":  "Amillusion",
	"unevendream": "Uneven Dream",
	"braingirl":   "Braingirl",
	"unconscious": "Collective Unconscious",
	"cerasus":     "Cerasus",
	"muma":        "Muma Rope",
	"genie":       "Dream Genie",
	"mikan":       "Mikan Muzou",
	"ultraviolet": "Ultra Violet",
	"sheawaits":   "She Awaits",
	"oversomnia":  "Oversomnia",
	"tagai":       "Yume Tagai",
	"tsushin":     "Yume Tsushin",
	"nostalgic":   "NostAlgic",
	"if":          "If",
}

var namespaceNumbers = map[string]string{
	"yume":        "3000",
	"2kki":        "3002",
	"flow":        "3004",
	"someday":     "3006",
	"deepdreams":  "3008",
	"prayers":     "3010",
	"amillusion":  "3012",
	"unevendream": "3014",
	"braingirl":   "3016",
	"unconscious": "3018",
	"cerasus":     "3020",
	"muma":        "3022",
	"genie":       "3026",
	"mikan":       "3028",
	"ultraviolet": "3030",
	"sheawaits":   "3032",
	"oversomnia":  "3034",
	"tagai":       "3036",
	"tsushin":     "3038",
	"nostalgic":   "3040",
	"if":          "3042",
}

var protagCategoriesPerGame = map[string]map[string]string{
	"unevendream": {
		"kubotsuki":  "Category:Kubotsuki's Worlds",
		"totsutsuki": "Category:Totsutsuki's Worlds",
	},
	"tagai": {
		"makitsuki": "Category:Makitsuki's Worlds",
		"sakiyuki":  "Category:Sakiyuki's Worlds",
	},
}

func createClient() (client *mwclient.Client, err error) {
	client, err = mwclient.New("https://yume.wiki/api.php", "yumeWikiAPIBot")
	client.SetHTTPTimeout(60000000000)
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

func GetLocations(gameParams GameParams) (locations *Locations, err error) {
	gameName, ok := gameNames[gameParams.GameCode]
	if !ok {
		return locations, errors.New("game not supported")
	}

	protagCategories, hasMultipleProtags := protagCategoriesPerGame[gameParams.GameCode]

	if !hasMultipleProtags && gameParams.Protag != "" {
		return locations, errors.New("game has only one protagonist")
	}

	if hasMultipleProtags && len(gameParams.Protag) == 0 {
		return locations, errors.New("game has multiple protagonists, please specify one")
	}

	locations = &Locations{
		Game: gameParams.GameCode,
	}

	protagCategory := ""
	if hasMultipleProtags && gameParams.Protag != "" {
		protagCategory, ok = protagCategories[gameParams.Protag]

		if !ok {
			return locations, errors.New("protagonist does not exist or is misspelled")
		}
	}

	client, err := createClient()
	if err != nil {
		return locations, err
	}

	condition := fmt.Sprintf("Category:%s Locations", gameName)
	if protagCategory != "" {
		condition += "|" + protagCategory
	}
	printouts := []string{"Has location image", "Header background color", "Header font color", "Has primary author", "Has contributing author", "Japanese name", "Has BGM", "Map IDs", "Has location map", "Version added", "Versions updated", "Version removed", "Version gaps"}

	parameters := params.Values{
		"action":      "askargs",
		"conditions":  condition,
		"printouts":   strings.Join(printouts, "|"),
		"parameters":  "limit=250",
		"format":      "json",
		"api_version": "3",
	}

	if gameParams.ContinueKey != "" {
		offset := fmt.Sprintf("|offset=%s", gameParams.ContinueKey)
		currentParams := parameters.Get("parameters")
		parameters.Set("parameters", currentParams+offset)
	}

	query, err := client.Get(parameters)
	if err != nil {
		return locations, err
	}

	continueKey, err := query.GetNumber("query-continue-offset")
	if err == nil {
		locations.ContinueKey = string(continueKey)
	}

	locationsToProcess, err := query.GetObjectArray("query", "results")
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

			locations.Locations = append(locations.Locations, location)
		}
	}

	return locations, err
}

func GetConnections(gameParams GameParams) (connections *Connections, err error) {
	gameName, ok := gameNames[gameParams.GameCode]
	if !ok {
		return connections, errors.New("game not supported")
	}

	connections = &Connections{
		Game: gameParams.GameCode,
	}
	protagCategories, hasMultipleProtags := protagCategoriesPerGame[gameParams.GameCode]

	if !hasMultipleProtags && gameParams.Protag != "" {
		return connections, errors.New("game has only one protagonist")
	}

	if hasMultipleProtags && len(gameParams.Protag) == 0 {
		return connections, errors.New("game has multiple protagonists, please specify one")
	}

	protagCategory := ""
	if hasMultipleProtags && gameParams.Protag != "" {
		protagCategory, ok = protagCategories[gameParams.Protag]

		if !ok {
			return connections, errors.New("protagonist does not exist or is misspelled")
		}
	}

	client, err := createClient()
	if err != nil {
		return connections, err
	}

	conditions := []string{fmt.Sprintf("%s:+", gameName), "Is subobject type::connection"}
	if protagCategory != "" {
		conditions = append(conditions, fmt.Sprintf("-Has subobject::<q>[[%s]]</q>", protagCategory))
	}
	printouts := []string{"Connection/Origin", "Connection/Location", "Connection/Attribute", "Connection/Unlock conditions", "Connection/Effects needed", "Connection/Season available", "Connection/Chance percentage", "Connection/Chance description", "Connection/Is removed"}

	parameters := params.Values{
		"action":      "askargs",
		"conditions":  strings.Join(conditions, "|"),
		"printouts":   strings.Join(printouts, "|"),
		"parameters":  "limit=500",
		"format":      "json",
		"api_version": "3",
	}

	if gameParams.ContinueKey != "" {
		offset := fmt.Sprintf("|offset=%s", gameParams.ContinueKey)
		currentParams := parameters.Get("parameters")
		parameters.Set("parameters", currentParams+offset)
	}

	query, err := client.Get(parameters)
	if err != nil {
		return connections, err
	}

	continueKey, err := query.GetNumber("query-continue-offset")
	if err == nil {
		connections.ContinueKey = string(continueKey)
	}

	connectionsToProcess, err := query.GetObjectArray("query", "results")
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

			connection, err := processConnection(gameParams.GameCode, value)
			if err != nil {
				return connections, err
			}

			connections.Connections = append(connections.Connections, connection)
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
	queryParams := []string{"sort=Author/Name", "order=asc", "limit=500"}

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
		"conditions":  conditions,
		"printouts":   printouts,
		"format":      "json",
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

func GetVendingMachines(gameCode string) (vendingMachines []*VendingMachine, err error) {
	vendingMachines = []*VendingMachine{}
	gameName, ok := gameNames[gameCode]
	if !ok {
		return vendingMachines, errors.New("game not supported")
	}

	client, err := createClient()
	if err != nil {
		return vendingMachines, err
	}

	conditions := []string{fmt.Sprintf("-Has subobject::%s:Vending Machine", gameName), "Vending Machine/Is implemented::true", "Vending Machine/Is accessible::true", "Vending Machine/Is secret::false"}
	printouts := []string{"Has image path", "Vending Machine/Map ID", "Vending Machine/Event ID"}
	queryParams := []string{"sort=Vending Machine/Location", "order=asc", "limit=500"}

	parameters := params.Values{
		"conditions":  strings.Join(conditions, "|"),
		"printouts":   strings.Join(printouts, "|"),
		"parameters":  strings.Join(queryParams, "|"),
		"format":      "json",
		"api_version": "3",
	}

	results := NewSmwQuery(client, parameters)
	vmsToProcess, err := fetchAllResultsFromSmwQuery(results)
	if err != nil {
		return vendingMachines, err
	}

	for _, vmToProcess := range vmsToProcess {
		if err != nil {
			return vendingMachines, err
		}

		for _, value := range vmToProcess.Map() {
			value, err := value.Object()
			if err != nil {
				return vendingMachines, err
			}

			vm, err := processVendingMachine(gameCode, value)
			if err != nil {
				return vendingMachines, err
			}

			vendingMachines = append(vendingMachines, vm)
		}
	}

	return vendingMachines, err
}

func GetImages(gameParams GameParams) (images *LocationImages, err error) {
	gameName, ok := gameNames[gameParams.GameCode]
	if !ok {
		return images, errors.New("game not supported")
	}
	images = &LocationImages{
		Game: gameParams.GameCode,
	}

	parameters := params.Values{
		"action":      "query",
		"format":      "json",
		"list":        "categorymembers",
		"cmtitle":     fmt.Sprintf("Category:%s Locations", gameName),
		"cmprop":      "title",
		"cmnamespace": namespaceNumbers[gameParams.GameCode],
		"cmlimit":     "50",
	}

	if gameParams.ContinueKey != "" {
		parameters.Set("continue", "-||")
		parameters.Set("cmcontinue", gameParams.ContinueKey)
	}

	client, err := createClient()
	if err != nil {
		return images, err
	}

	query, err := client.Get(parameters)
	if err != nil {
		return images, err
	}

	continueKey, err := query.GetString("continue", "cmcontinue")
	if err == nil {
		images.ContinueKey = continueKey
	}

	pagesToProcess, err := query.GetObjectArray("query", "categorymembers")

	for _, pageToProcess := range pagesToProcess {
		pageTitle, err := pageToProcess.GetString("title")
		if err != nil {
			return images, err
		}

		title := strings.Split(pageTitle, ":")[1]
		pageImage := &LocationImage{
			Title: title,
			Game:  gameParams.GameCode,
		}

		parameters := params.Values{
			"action":      "query",
			"format":      "json",
			"prop":        "imageinfo",
			"list":        "",
			"titles":      pageTitle,
			"generator":   "images",
			"iiprop":      "size|url",
			"iiurlwidth":  "320",
			"iiurlheight": "240",
			"gimlimit":    "max",
		}

		results, err := client.Get(parameters)
		if err != nil {
			return images, err
		}

		pageImagesToProcess, err := results.GetObjectArray("query", "pages")
		if err != nil {
			return images, err
		}

		for _, pageImageToProcess := range pageImagesToProcess {
			imageInfoToProcess, err := pageImageToProcess.GetObjectArray("imageinfo")
			if err != nil {
				continue
			}

			for _, imageInfo := range imageInfoToProcess {
				image := &Image{}
				url, err := imageInfo.GetString("url")
				if err != nil {
					return images, err
				}

				width, err := imageInfo.GetNumber("width")
				if err != nil {
					return images, err
				}

				height, err := imageInfo.GetNumber("height")
				if err != nil {
					return images, err
				}

				thumburl, err := imageInfo.GetString("thumburl")
				if err == nil {
					image.Url = thumburl
				} else {
					image.Url = url
				}

				thumbwidth, err := imageInfo.GetNumber("thumbwidth")
				if err == nil {
					image.Width = thumbwidth
				} else {
					image.Width = width
				}

				thumbheight, err := imageInfo.GetNumber("thumbheight")
				if err == nil {
					image.Height = thumbheight
				} else {
					image.Height = height
				}

				pageImage.Images = append(pageImage.Images, image)
			}
		}
		images.LocationImages = append(images.LocationImages, pageImage)
	}
	return images, err
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
		location.PrimaryAuthor = strings.Join(primaryAuthor, ", ")
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

	mapIdObjects, err := printouts.GetObjectArray("Map IDs")
	if err != nil {
		log.Print("SERVER", "mapIds", err.Error())
		return nil, err
	}

	mapIds, err := processMapIdInfo(mapIdObjects)
	if err != nil {
		log.Print(title)
		log.Print("SERVER", "mapIds", err.Error())
		return nil, err
	}

	location.MapIds = mapIds

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

func processMapIdInfo(mapIdObjects []*jason.Object) (mapIds []int, err error) {
	mapIds = []int{}
	for _, info := range mapIdObjects {
		var mapId json.Number
		outputMapId, err := info.GetNumberArray("Has map ID", "item")

		if err != nil {
			log.Print("SERVER", "mapId", err.Error())
			return mapIds, err
		}

		if len(outputMapId) > 0 {
			mapId = outputMapId[0]
		}

		data := mapId.String()

		if id, err := strconv.Atoi(data); err == nil {
			mapIds = append(mapIds, id)
		}
	}

	return mapIds, nil
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

	unlockConditions, err := printouts.GetStringArray("Connection/Unlock conditions")
	if err != nil {
		log.Print("SERVER", "unlockConditions", err.Error())
		return connection, err
	}

	effectsNeeded, err := printouts.GetStringArray("Connection/Effects needed")
	if err != nil {
		log.Print("SERVER", "effectsNeeded", err.Error())
		return connection, err
	}

	seasonAvailable, err := printouts.GetStringArray("Connection/Season available")
	if err != nil {
		log.Print("SERVER", "seasonAvailable", err.Error())
		return connection, err
	}

	chancePercentage, err := printouts.GetStringArray("Connection/Chance percentage")
	if err != nil {
		log.Print("SERVER", "chancePercentage", err.Error())
		return connection, err
	}

	chanceDescription, err := printouts.GetStringArray("Connection/Chance description")
	if err != nil {
		log.Print("SERVER", "chanceDescription", err.Error())
		return connection, err
	}

	isRemoved, err := printouts.GetStringArray("Connection/Is removed")
	if err != nil {
		log.Print("SERVER", "isRemoved", err.Error())
		return connection, err
	}

	connection = &Connection{
		Game:        gameCode,
		Origin:      origin,
		Destination: destination,
		Attributes:  attributes,
	}

	if len(isRemoved) > 0 && isRemoved[0] == "t" {
		connection.IsRemoved = true
	}

	if len(unlockConditions) > 0 {
		connection.UnlockConditions = unlockConditions[0]
	}

	if len(effectsNeeded) > 0 {
		connection.EffectsNeeded = effectsNeeded
	}

	if len(seasonAvailable) > 0 {
		connection.SeasonAvailable = seasonAvailable[0]
	}

	if len(chancePercentage) > 0 {
		connection.ChancePercentage = chancePercentage[0]
	}

	if len(chanceDescription) > 0 {
		connection.ChanceDescription = chanceDescription[0]
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

func processVendingMachine(gameCode string, value *jason.Object) (vendingMachine *VendingMachine, err error) {
	vendingMachine = &VendingMachine{
		Game: gameCode,
	}
	printouts, err := value.GetObject("printouts")
	if err != nil {
		return nil, err
	}

	path, err := printouts.GetStringArray("Has image path")
	if err != nil {
		log.Print("SERVER", "path", err.Error())
		return nil, err
	}
	if len(path) > 0 {
		vendingMachine.Path = path[0]
	}

	mapId, err := printouts.GetStringArray("Vending Machine/Map ID")
	if err != nil {
		log.Print("SERVER", "mapId", err.Error())
		return nil, err
	}

	if len(mapId) > 0 {
		vendingMachine.MapId = mapId[0]
	}

	eventIds, err := printouts.GetStringArray("Vending Machine/Event ID")
	if err != nil {
		log.Print("SERVER", "eventIds", err.Error())
		return nil, err
	}

	if len(eventIds) > 0 {
		vendingMachine.EventIds = eventIds
	}

	return vendingMachine, err
}
