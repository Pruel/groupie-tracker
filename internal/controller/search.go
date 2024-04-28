package controller


// SearchController 
func SearchController(w http.ResponseWriter, r *http.Request) {
    // 1. template parse

    // getSearchValue - call second function

    // getMainData = recieve all groups from web api, and return []entity.Artists -> webapi.New().GetAllArtists()

    // Search -> searchByGroupName


    // last call template with found groups as context data of the template

}

// II
// getSearchValue
    // read request body = r.FormValue("search") // <input type="text" name="search"
    // check search value, if str is empty -> return error
    // if ok, return search value


// III
// Search(searchValue string, artists []entity.Artists), search by group name, and return found groups = foundGroups []entity.Artists
    //     

    // check artists on not nil (artists type slice, zero value of the slice is nil)
    // if ok
        // for _, group := range artists {
        //     if strings.Contains(group.Name, searchValue) -> ok = true {
                        // foundGroups = append(foundGroups, group)
        //        }
        //      for _, loc := group.Locations {
                    // if loc == searchValue {}
//         }
        // }
    //  after return foundGroups
    
 
