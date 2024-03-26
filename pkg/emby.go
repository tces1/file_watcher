package pkg

import (
	"fmt"
	"net/http"
)

// RefreshEmbyLibrary sends a POST request to Emby to refresh the library.
func RefreshEmbyLibrary(url, api string) {
	r, e := http.Post(fmt.Sprintf("%s/Library/Refresh?api_key=%s", url, api), "application/json", nil)
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println(r.Status)
}
