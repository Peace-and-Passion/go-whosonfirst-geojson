package main

import (
	"flag"
	"fmt"
	geojson "github.com/whosonfirst/go-whosonfirst-geojson"
	"github.com/kellydunn/golang-geo"
)


func main() {

	flag.Parse()
	args := flag.Args()

	for _, path := range args {

		f, parse_err := geojson.UnmarshalFile(path)

		if parse_err != nil {
			panic(parse_err)
		}

		polygons := f.GeomToPolygons()

		fmt.Printf("%s has this many polygons: %d\n", path, len(polygons))

		pt := geo.NewPoint(40.681, -73.986)

		for _, poly := range polygons {

		    c := poly.Contains(pt)
		    fmt.Printf("contained: %v\n", c)
		}
	}

}
