package main

import (
	"fmt"
	"github.com/pedrocb/random-album-picker/pkg"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("No username provided")
		return
	}
	user := args[1]
	album, err := random_album.GetRandomAlbumFromRYM(user)
	if err != nil {
		fmt.Printf("Got error %s", err.Error())
	} else {
		fmt.Println(album)
	}

}
