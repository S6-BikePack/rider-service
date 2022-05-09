package domain

import (
	"database/sql/driver"
	"encoding/hex"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"github.com/twpayne/go-geom/encoding/wkt"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Location struct {
	Latitude  float64
	Longitude float64
}

func (l Location) Value() (driver.Value, error) {

	g := geom.NewPointFlat(geom.XY, geom.Coord{l.Longitude, l.Latitude})

	e, err := wkt.Marshal(g)

	if err != nil {
		return nil, err
	}

	return "SRID=4326;" + e, nil
}

func (l *Location) Scan(value interface{}) error {
	t, err := hex.DecodeString(value.(string))
	if err != nil {
		return err
	}

	gt, err := ewkb.Unmarshal(t)
	if err != nil {
		return err
	}
	p := Location{
		Latitude:  gt.FlatCoords()[1],
		Longitude: gt.FlatCoords()[0],
	}
	*l = p

	return nil
}

func (Location) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return "geometry(Point, 4326)"
}
