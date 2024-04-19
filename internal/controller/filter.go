package controller

import (
	"groupie-tracker/internal/entity"
	"log/slog"
	"net/http"
	"strconv"
)

// FilterController
func FilterController(w http.ResponseWriter, r *http.Request) {
	// read request body and save filter params in some structure
    if r.Method != http.MethodPost {
        slog.Debug("Invalid request")
        ErrorController(w, r)
    }
    
    flt := entity.Filters{}

    r.ParseForm()

    creationDate := r.FormValue("creationDate")
    firstAlbumDate := r.FormValue("firstAlbumDate")
    locations := r.FormValue("locations")
    numMembers := make([]int, 0, 8)

	for i := 1; i < 8; i++ {
        membersKey := fmt.Sprintf("members%d", i)
        if mem := r.FormValue(membersKey); mem != "" {
            mNum, err := strconv.Atoi(mem)
            if err != nil {
                slog.Debug(err.Error())
            }
            numMembers = append(numMembers, mNum)
        }
	}
    
	// GET = MainController

	
	// POST = FilterController -> post redirect to main page


	http.Redirect(w, r, "/")
}