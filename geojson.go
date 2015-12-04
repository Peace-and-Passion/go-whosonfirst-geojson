package geojson

import (
	"fmt"
	rtreego "github.com/dhconnelly/rtreego"
	gabs "github.com/jeffail/gabs"
	geo "github.com/kellydunn/golang-geo"
	ioutil "io/ioutil"
	"strconv"
	"sync"
)

/*

- gabs is what handles marshaling a random bag of GeoJSON
- rtreego is imported to convert a WOFFeature in to a handy rtreego.Spatial object for indexing by go-whosonfirst-pip
- geo is imported to convert a WOFFeature geometry into a list of geo.Polygon objects for doing containment checks in go-whosonfirst-pip
  (only Polygons and MultiPolygons are supported at the moment)

*/

type WOFError struct {
	s string
}

func (e *WOFError) Error() string {
	return e.s
}

// See also
// https://github.com/dhconnelly/rtreego#storing-updating-and-deleting-objects

type WOFSpatial struct {
	bounds    *rtreego.Rect
	Id        int
	Name      string
	Placetype string
	Offset    int // used when calling EnSpatializeGeom in order to know which polygon we care about
}

type WOFPolygon struct {
	OuterRing     geo.Polygon
	InteriorRings []geo.Polygon
}

func (p *WOFPolygon) CountPoints() int {

	count := len(p.OuterRing.Points())

	for _, r := range p.InteriorRings {
		count += len(r.Points())
	}

	return count
}

func (p *WOFPolygon) Contains(latitude float64, longitude float64) bool {

	pt := geo.NewPoint(latitude, longitude)
	contains := false

	if p.OuterRing.Contains(pt) {
		contains = true
	}

	if contains && len(p.InteriorRings) > 0 {

		wg := new(sync.WaitGroup)

		for _, r := range p.InteriorRings {

			wg.Add(1)

			go func(poly geo.Polygon, point *geo.Point) {

				defer wg.Done()

				/*

					File under yak-shaving: Some way to send an intercept to poly.Contains
					to stop the raycasting if any one of these goroutines gets the answer
					it needs independent the results of the others. Like I said... yaks.
					(20151028/thisisaaronland)
				*/

				if poly.Contains(point) {
					contains = false
				}

			}(r, pt)
		}

		wg.Wait()
	}

	return contains
}

func (sp WOFSpatial) Bounds() *rtreego.Rect {
	return sp.bounds
}

type WOFFeature struct {
	// Raw    []byte
	Parsed *gabs.Container
}

func (wof WOFFeature) Body() *gabs.Container {
	return wof.Parsed
}

func (wof WOFFeature) Dumps() string {
	return wof.Parsed.String()
}

/*
func (wof WOFFeature) Id(path string) int {

	return wof.id(path)
}
*/

func (wof WOFFeature) Id() int {

	path := "id"
	return wof.id(path)
}

func (wof WOFFeature) WOFId() int {

	path := "properties.wof:id"
	return wof.id(path)
}

func (wof WOFFeature) id(path string) int {

	body := wof.Body()

	var id_float float64
	var id_str string
	var id int

	var ok bool

	// what follows shouldn't be necessary but appears to be
	// for... uh, reasons (20151013/thisisaaronland)

	id_float, ok = body.Path(path).Data().(float64)

	if ok {
		id = int(id_float)
	} else {
		id, ok = body.Path(path).Data().(int)
	}

	// But wait... there's more (20151028/thisisaaronland)

	if !ok {

		id_str, ok = body.Path(path).Data().(string)

		if ok {

			id_int, err := strconv.Atoi(id_str)

			if err != nil {
				ok = false
			} else {
				id = id_int
			}
		}
	}

	if !ok {
		id = -1
	}

	return id
}

func (wof WOFFeature) Name(path string) string {

	return wof.name(path)
}

func (wof WOFFeature) WOFName() string {

	path := "properties.wof:name"
	return wof.name(path)
}

func (wof WOFFeature) name(path string) string {

	name, ok := wof.StringValue(path)

	if !ok {
		name = "A Place With No Name"
	}

	return name
}

func (wof WOFFeature) Placetype(path string) string {

	return wof.placetype(path)
}

func (wof WOFFeature) WOFPlacetype() string {

	path := "properties.wof:placetype"
	return wof.placetype(path)
}

func (wof WOFFeature) placetype(path string) string {

	placetype, ok := wof.StringValue(path)

	if !ok {
		placetype = "unknown"
	}

	return placetype
}

func (wof WOFFeature) StringProperty(prop string) (string, bool) {

	path := fmt.Sprintf("properties.%s", prop)
	return wof.StringValue(path)
}

func (wof WOFFeature) StringValue(path string) (string, bool) {

	body := wof.Body()

	var value string
	var ok bool

	value, ok = body.Path(path).Data().(string)
	return value, ok
}

func (wof WOFFeature) EnSpatialize() (*WOFSpatial, error) {

	id := wof.WOFId()
	name := wof.WOFName()
	placetype := wof.WOFPlacetype()

	body := wof.Body()

	var swlon float64
	var swlat float64
	var nelon float64
	var nelat float64

	children, _ := body.S("bbox").Children()

	if len(children) != 4 {
		return nil, &WOFError{"weird and freaky bounding box"}
	}

	swlon = children[0].Data().(float64)
	swlat = children[1].Data().(float64)
	nelon = children[2].Data().(float64)
	nelat = children[3].Data().(float64)

	llat := nelat - swlat
	llon := nelon - swlon

	/*
		fmt.Printf("%f - %f = %f\n", nelat, swlat, llat)
		fmt.Printf("%f - %f = %f\n", nelon, swlon, llon)
	*/

	pt := rtreego.Point{swlon, swlat}
	rect, err := rtreego.NewRect(pt, []float64{llon, llat})

	if err != nil {
		return nil, err
	}

	return &WOFSpatial{rect, id, name, placetype, -1}, nil
}

func (wof WOFFeature) EnSpatializeGeom() ([]*WOFSpatial, error) {

	id := wof.WOFId()
	name := wof.WOFName()
	placetype := wof.WOFPlacetype()

	spatial := make([]*WOFSpatial, 0)
	polygons := wof.GeomToPolygons()

	for offset, poly := range polygons {

		swlat := 0.0
		swlon := 0.0
		nelat := 0.0
		nelon := 0.0

		for _, pt := range poly.OuterRing.Points() {

			lat := pt.Lat()
			lon := pt.Lng()

			if swlat == 0.0 || swlat > lat {
				swlat = lat
			}

			if swlon == 0.0 || swlon > lon {
				swlon = lon
			}

			if nelat == 0.0 || nelat < lat {
				nelat = lat
			}

			if nelon == 0.0 || nelon < lon {
				nelon = lon
			}

			// fmt.Println(lat, lon, swlat, swlon, nelat, nelon)
		}

		llat := nelat - swlat
		llon := nelon - swlon

		pt := rtreego.Point{swlon, swlat}
		rect, err := rtreego.NewRect(pt, []float64{llon, llat})

		if err != nil {
			return nil, err
		}

		sp := WOFSpatial{rect, id, name, placetype, offset}
		spatial = append(spatial, &sp)
	}

	return spatial, nil
}

func (wof WOFFeature) Contains(latitude float64, longitude float64) bool {

	polygons := wof.GeomToPolygons()
	contains := false

	wg := new(sync.WaitGroup)

	for _, p := range polygons {

		wg.Add(1)

		go func(poly *WOFPolygon, lat float64, lon float64) {

			defer wg.Done()

			if poly.Contains(lat, lon) {
				contains = true
			}

		}(p, latitude, longitude)
	}

	wg.Wait()

	return contains
}

func (wof WOFFeature) GeomToPolygons() []*WOFPolygon {

	body := wof.Body()

	var geom_type string

	geom_type, _ = body.Path("geometry.type").Data().(string)
	children, _ := body.S("geometry").ChildrenMap()

	polygons := make([]*WOFPolygon, 0)

	for key, child := range children {

		if key != "coordinates" {
			continue
		}

		var coordinates []interface{}
		coordinates, _ = child.Data().([]interface{})

		if geom_type == "Polygon" {
			polygons = append(polygons, wof.DumpPolygon(coordinates))
		} else if geom_type == "MultiPolygon" {
			polygons = wof.DumpMultiPolygon(coordinates)
		} else {
			// pass
		}
	}

	return polygons
}

func (wof WOFFeature) DumpMultiPolygon(coordinates []interface{}) []*WOFPolygon {

	polygons := make([]*WOFPolygon, 0)

	for _, ipolys := range coordinates {

		polys := ipolys.([]interface{})

		polygon := wof.DumpPolygon(polys)
		polygons = append(polygons, polygon)
	}

	return polygons
}

func (wof WOFFeature) DumpPolygon(coordinates []interface{}) *WOFPolygon {

	polygons := make([]geo.Polygon, 0)

	for _, ipoly := range coordinates {

		poly := ipoly.([]interface{})
		polygon := wof.DumpCoords(poly)
		polygons = append(polygons, polygon)
	}

	return &WOFPolygon{
		OuterRing:     polygons[0],
		InteriorRings: polygons[1:],
	}
}

func (wof WOFFeature) DumpCoords(poly []interface{}) geo.Polygon {

	polygon := geo.Polygon{}

	for _, icoords := range poly {

		coords := icoords.([]interface{})

		lon := coords[0].(float64)
		lat := coords[1].(float64)

		pt := geo.NewPoint(lat, lon)
		polygon.Add(pt)
	}

	return polygon
}

func UnmarshalFile(path string) (*WOFFeature, error) {

	body, read_err := ioutil.ReadFile(path)

	if read_err != nil {
		return nil, read_err
	}

	return UnmarshalFeature(body)
}

func UnmarshalFeature(raw []byte) (*WOFFeature, error) {

	parsed, parse_err := gabs.ParseJSON(raw)

	if parse_err != nil {
		return nil, parse_err
	}

	rsp := WOFFeature{
		// Raw:    raw,
		Parsed: parsed,
	}

	return &rsp, nil
}
