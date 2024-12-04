package util

import (
	"database/sql"
	"log"

	"github.com/IncSW/geoip2"
	_ "github.com/mattn/go-sqlite3"
)

var (
	DB          *sql.DB
	GeoLiteCity *geoip2.CityReader
	err         error
)

func init() {
	DB, err = sql.Open("sqlite3", "file:netvigil.db")
	if err != nil {
		log.Panicln("Failed to open database:", err)
	}
	GeoLiteCity, err = geoip2.NewCityReaderFromFile("GeoLite2-City.mmdb")
	if err != nil {
		log.Panicln("Failed to open GeoLite2-City.mmdb:", err)
	}
}
