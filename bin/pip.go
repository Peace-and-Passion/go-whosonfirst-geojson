package main

import (
	"flag"
	"fmt"
	"github.com/kellydunn/golang-geo"
	geojson "github.com/whosonfirst/go-whosonfirst-geojson"
)

func main() {

	/*
		./bin/pip /usr/local/mapzen/whosonfirst-data/data/404/529/181/404529181.geojson /usr/local/mapzen/whosonfirst-data/data/857/848/31/85784831.geojson
		/usr/local/mapzen/whosonfirst-data/data/404/529/181/404529181.geojson has this many polygons: 1
		/usr/local/mapzen/whosonfirst-data/data/404/529/181/404529181.geojson polygon 1 has 3 interior rings
		contained: false
		/usr/local/mapzen/whosonfirst-data/data/857/848/31/85784831.geojson has this many polygons: 1
		/usr/local/mapzen/whosonfirst-data/data/857/848/31/85784831.geojson polygon 1 has 0 interior rings
		contained: true
	*/

	flag.Parse()
	args := flag.Args()

	pt := geo.NewPoint(45.523668, -73.600159)

	for _, path := range args {

		f, parse_err := geojson.UnmarshalFile(path)

		if parse_err != nil {
			panic(parse_err)
		}

		polygons := f.GeomToPolygons()

		fmt.Printf("%s has this many polygons: %d\n", path, len(polygons))

		for i, poly := range polygons {

			fmt.Printf("%s polygon %d has %d interior rings\n", path, (i + 1), len(poly.InteriorRings))

			c := poly.Contains(pt)
			fmt.Printf("contained: %t\n", c)
		}
	}

}
