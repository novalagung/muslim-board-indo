package usecase

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"muslimboard-api.novalagung.com/repositories/aladhan"
	"muslimboard-api.novalagung.com/repositories/openstreetmap"
	"muslimboard-api.novalagung.com/repositories/unsplash"
)

func GetImage(ctx context.Context, imageUrl string) (string, io.ReadCloser, error) {
	return unsplash.GetImage(ctx, imageUrl)
}

// GetShalatScheduleByCoordinate is handler of get shalat schedule by coordinate
func GetShalatScheduleByCoordinate(ctx context.Context, method, latitude, longitude, month, year string) (map[string]interface{}, error) {
	namespace := "handler.GetShalatScheduleByCoordinate"

	// if lat long is invalid, then simply return true
	latInt, _ := strconv.ParseFloat(latitude, 64)
	lonInt, _ := strconv.ParseFloat(longitude, 64)

	schedules, err := aladhan.GetShalatScheduleByCoordinate(ctx, method, latInt, lonInt, month, year)
	if err != nil {
		log.Errorln(namespace, "getShalatScheduleByCoordinate", err.Error())
		return nil, err
	}

	res := map[string]interface{}{
		"schedules":   schedules,
		"address":     "",
		"countryCode": "id",
	}

	locationRes, err := openstreetmap.GetLocationByCoordinate(ctx, latitude, longitude)
	if err != nil {
		log.Errorln(namespace, "getLocationByCoordinate", err.Error())
		res["address"] = fmt.Sprintf("Location %v, %v", latitude, longitude)
		res["countryCode"] = ""
	} else {
		res["address"] = locationRes["address"]
		res["countryCode"] = locationRes["countryCode"]
	}

	return res, nil
}

// GetShalatScheduleByLocation is handler of get shalat schedule by location
// for now, immediately use aladhan.com api coz kemenag backend still under development
func GetShalatScheduleByLocation(ctx context.Context, method, province, city, month, year string) (map[string]interface{}, error) {
	namespace := "handler.GetShalatScheduleByLocation"

	location := fmt.Sprintf("%s,%s", city, province)
	location = strings.ToLower(location)
	location = strings.Replace(location, "d.i. ", "", -1)
	location = strings.Replace(location, "kab. ", "", -1)
	location = strings.Replace(location, "kota ", "", -1)
	location = strings.TrimSpace(location)

	// get coordinate by location
	coordinate, err := openstreetmap.GetCoordinateByLocation(ctx, location)
	if err != nil {
		log.Errorln(namespace, "getCoordinateByLocation", err.Error())
		return nil, err
	}

	// parse result
	latitude, _ := strconv.ParseFloat(coordinate["lat"].(string), 64)
	longitude, _ := strconv.ParseFloat(coordinate["lon"].(string), 64)

	schedules, err := aladhan.GetShalatScheduleByCoordinate(ctx, method, latitude, longitude, month, year)
	if err != nil {
		log.Errorln(namespace, "getShalatScheduleByCoordinate", err.Error())
		return nil, err
	}

	address := cases.Title(language.English).String(strings.ToLower(fmt.Sprintf("%s, %s", city, province)))
	res := map[string]interface{}{
		"schedules":   schedules,
		"address":     address,
		"countryCode": "id",
	}

	return res, nil
}
