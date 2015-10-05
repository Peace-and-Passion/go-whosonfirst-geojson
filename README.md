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

## The long version

This is really a "GeoJSON" specific library, yet. Right now it's just a thin wrapper around the [Gabs](https://github.com/jeffail/gabs) utility for wrangling unknown JSON structures in to a Go `WOFFeature` struct.

Eventually it would be nice to make Gabs hold hands with Paul Mach's [go.geojson](https://github.com/paulmach/go.geojson) and use the former to handle the GeoJSON properties dictionary. But that day is not today.

## See also

* https://www.github.com/jeffail/gabs
