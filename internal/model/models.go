package model

type Artist struct {
	ID           int
	Name         string
	Members      []string
	CreationDate int
	FirstAlbum   string
	Locations    string
	Dates        string
	Relation     string
	Image        string
}

type Locations struct {
	ID        int
	Dates     string
	Locations []string
}

type Dates struct {
	ID    int
	Dates []string
}

type Relations struct {
	ID             int
	DatesLocations map[string][]string
}

type LocationsIndex struct {
	Index []Locations
}

type DatesIndex struct {
	Index []Dates
}

type RelationsIndex struct {
	Index []Relations
}
