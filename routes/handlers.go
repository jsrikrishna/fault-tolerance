package routes

import (
	"net/http"
	"fmt"
)

func Resources(w http.ResponseWriter, r *http.Request){
	fmt.Println("Got the request")
}
