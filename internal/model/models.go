package model

type Artist struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	Dates        string   `json:"dates"`
	Relation     string   `json:"relations"`
	Image        string   `json:"image"`
}

type ArtistInfo struct {
	Artist    Artist
	Locations Locations
	Dates     Dates
	Relations Relations
}

type Locations struct {
	ID        int      `json:"id"`
	Dates     string   `json:"dates"`
	Locations []string `json:"locations"`
}

type Dates struct {
	ID    int    `json:"id"`
	Dates string `json:"dates"`
}

type Relations struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type LocationsIndex struct {
	Index []Locations `json:"index"`
}

type DatesIndex struct {
	Index []Dates `json:"index"`
}

type RelationsIndex struct {
	Index []Relations `json:"index"`
}
