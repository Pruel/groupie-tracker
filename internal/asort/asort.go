package asort // sort

import (
	"groupie-tracker/internal/entity"
	"time"
	// "sort"
)

type SArtists []entity.Artist

//  Type Alice
//  type MyInt int64

// sort - interface = len, swap, less

// Len ...
func (a SArtists) Len() int {
	// sort.Interface

	return len(a)
}

// Swap ...
func (a SArtists) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Less
func (a SArtists) Less(i, j int) bool {
	dateFormat := "01-01-1970"

	// first album date
	falbDate, _ := time.Parse(dateFormat, a[i].FirstAlbum)
	lalbDate, _ := time.Parse(dateFormat, a[j].FirstAlbum)

	if falbDate.Year() != lalbDate.Year() {
		return falbDate.Year() < lalbDate.Year()
	}

	return a[i].CreationDate < a[j].CreationDate
}
