# go-whosonfirst-geojson

Go tools for working with Who's On First documents

## Usage

```
package main

import (
	"flag"
	"fmt"
	geojson "github.com/whosonfirst/go-whosonfirst-geojson/whosonfirst"
)

func main() {

	flag.Parse()
	args := flag.Args()

	for _, path := range args {

		// This is mostly just helper code to read the file
		// and call geojson.UnmarshalFeature (from a bag of bytes)

		f, parse_err := geojson.UnmarshalFile(path)

		if parse_err != nil {
			panic(parse_err)
		}

		fmt.Printf("# %s\n", path)
		fmt.Println(f.Dumps())
	}

}
```

## The longer version

This isn't really a "GeoJSON" specific library, yet. Right now it's just a thin wrapper around the [Gabs](https://github.com/jeffail/gabs) utility for wrangling unknown JSON structures in to a Go `WOFFeature` struct.

Eventually it would be nice to make Gabs hold hands with Paul Mach's [go.geojson](https://github.com/paulmach/go.geojson) and use the former to handle the GeoJSON properties dictionary. But that day is not today.

## The longer longer version

Right now this library has evolved and grown functionality on as-needed basis, targeting on Who's On First specific use-cases. As such it consists of a handful of WOF struct types - `WOFFeature` and `WOFPolygon` and `WOFSpatial` - that are wrappers around other people's heavy-lifting. There are not any WOF related interfaces but that's really the direction we want to head in... but we're not there yet. So things will probably change in the short-term. Not too much , hopefully.

## See also

* https://www.github.com/jeffail/gabs
