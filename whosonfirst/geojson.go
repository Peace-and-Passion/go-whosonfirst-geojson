package geojson

import (
	"github.com/jeffail/gabs"
	"io/ioutil"
)

/*
Something like this using "github.com/paulmach/go.geojson" seems
like it would be a good thing but I don't think I have the stamina
to figure out how to parse the properties separately right now...
(20151005/thisisaaronland)
/*

type WOFProperties struct {
     Raw    []byte
     Parsed *gabs.Container
}

type WOFFeature struct {
     ID          json.Number            `json:"id,omitempty"`
     Type        string                 `json:"type"`
     BoundingBox []float64              `json:"bbox,omitempty"`
     Geometry    *gj.Geometry      `json:"geometry"`
     Properties  WOFProperties			`json:"properties"`
     // Properties  map[string]interface{} `json:"properties"`
     CRS         map[string]interface{} `json:"crs,omitempty"` // Coordinate Reference System Objects are not currently supported
}
*/

type WOFFeature struct {
	Raw    []byte
	Parsed *gabs.Container
}

func (wof WOFFeature) Body() *gabs.Container {
	return wof.Parsed
}

func (wof WOFFeature) Dumps() string {
	return wof.Parsed.String()
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
		Raw:    raw,
		Parsed: parsed,
	}

	return &rsp, nil
}
