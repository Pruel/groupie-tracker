

1. Создаём структуру для наших будущих данных. +++

type MainData struct {  Эта структура объединяющий контейнер который мы потом передадим в index.html.  
	Artists     []Artist
	FiltersData Filters
}

type Filters struct {  
	FirstCreationDate int
	LastCreationDate  int
	LowestFirstAlbum  string
	HighestFirstAlbum string
	Members           []int
	Locations         []string
	Мы создаём эту структуру в зависимости от того что нам нужно, в нашем случае нам требовалось всё что в ней описано по заданию groupier-tracker-filter 
}

2. Создаём фронтенд для проекта. +++ 
3. --- 
4. ---
5. Сортировка фильтрационных параметров. 
6.  
Сразу ссылаемся на нашу структуру что-бы сохранитьт результат фильтрации и сортировки, но нам не нужно ничего сохранить так как slice под капотом это ссылка на массив и поэтому все экземпляры artists ссылаются на один и тот-же участок памяти  


	// 3. receive filters params from front and save this params in filters structure (Пишем структуру опять куда сохранять) ---

	// 4. check and validate data, and parse or converting filters data ---

	// 5. do sort by filters params +++

	// 6. create artists slice for saving after filtering +++

	// Artists = type slice = *Pointer

	// 7. after do filters -> for -> element equel by filter param -> true -> append filteredArtists type = model.Artists ---

	// 8. return filteredArtists ---
