package main

import (
	"flag"
	"fmt"
	geojson "github.com/whosonfirst/go-whosonfirst-geojson"
)

func main() {

	flag.Parse()
	args := flag.Args()

	for _, path := range args {

		f, parse_err := geojson.UnmarshalFile(path)

		if parse_err != nil {
			panic(parse_err)
		}

		sp, _ := f.EnSpatialize()
		fmt.Printf("%v\n", sp)
		fmt.Printf("%d\n", sp.Id)
	}

}
