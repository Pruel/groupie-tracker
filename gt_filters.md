
	// 1. prepare for filtering data
	// Делаем структуру для наших данных.
	// mainData -> artists, filtersData
	// filter -> creationDate, firstAlbum, Members, Locations
	// Artists

	// 2.
	// FRONT all data in one form
	// <form>, begin

	// creationDate
	// inpute -> type range, value=1950, min=1900, max=2020, name=creationDate

	// firstAlbum
	// input type date, name=firstAlbum, value=1900

	// members
	// input type check box, value=1, name=members

	// locations
	// select name=locations
	// {{ range .locations }}
	// option value=some_location
	// {{ end }}
	// </form> end

	// 3. receive filters params from front and save this params in filters structure (Пишем структуру опять куда сохранять)

	// 4. check and validate data, and parse or converting filters data

	// 5. do sort by filters params

	// 6. create artists slice for saving after filtering

	// 7. after do filters -> for -> element equel by filter param -> true -> append filteredArtists type = model.Artists

	// 8. return filteredArtists
