package filter

import (
	"errors"
	"groupie-tracker/internal/asort"
	"groupie-tracker/internal/entity"
	"groupie-tracker/internal/webapi"
	"log/slog"
	"sort"
)

// PrepareFilterData
func PrepareFilterData(artists []entity.Artist) (*entity.Filters, error) {
	if len(artists) == 0 {
		return nil, errors.New("error, artists slice must be not lower then zero")
	}
	filter := &entity.Filters{}
	// AddCreationDate
	filter = AddCreationDate(artists, filter)

	// AddFirstAlbumPublishDate
	filter = AddFirstAlbumPublishDate(artists, filter)

	// AddNumMembers
	// filter.Members = []int{1, 2, 3, 4, 5, 6, 7}
	filter = AddNumMembers(artists, filter)

	// AddLocations - not need to sort
	filter = AddLocations(artists, filter)

	return filter, nil
}

// AddCreationDate
func AddCreationDate(artists []entity.Artist, filter *entity.Filters) *entity.Filters {
	sort.Sort(asort.SArtists(artists))
	lasti := len(artists) - 1

	filter.FirstCreationDate = artists[0].CreationDate
	filter.LastCreationDate = artists[lasti].CreationDate

	return filter
}

func AddFirstAlbumPublishDate(artists []entity.Artist, filter *entity.Filters) *entity.Filters {
	sort.Sort(asort.SArtists(artists))
	lastalb := len(artists) - 1

	filter.LowestFirstAlbum = artists[0].FirstAlbum
	filter.HighestFirstAlbum = artists[lastalb].FirstAlbum

	return filter
}

func AddNumMembers(artists []entity.Artist, filter *entity.Filters) *entity.Filters {
	alSize := len(artists)
	numbersMembers := make(map[int]int, alSize)
	numbers := make([]int, 0, alSize)

	for _, elem := range artists {
		number := len(elem.Members)
		numbers = append(numbers, number)
	}

	for _, elem := range numbers {
		if _, ok := numbersMembers[elem]; !ok {
			numbersMembers[elem] = elem
			filter.Members = append(filter.Members, elem)
		}

	}

	sort.Ints(filter.Members)

	return filter
}

func AddLocations(artists []entity.Artist, filter *entity.Filters) *entity.Filters {
	locations, err := webapi.New().GetAllUniqueLocations()
	if err != nil {
		slog.Error(err.Error())
		return nil
	}

	filter.Locations = locations

	return filter
}
