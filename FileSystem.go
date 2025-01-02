package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.Handle("/",http.FileServer(http.Dir("/HDD_WD/dir_hdd/Sessions/")))

	e := http.ListenAndServe(":8082", nil)
	fmt.Println(e)
}

