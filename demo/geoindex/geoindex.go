package hello

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	index "github.com/mariotoffia/go-geoindex"
)

func init() {
	http.HandleFunc("/points", points)
	http.HandleFunc("/knearest", knearest)
}

var geoindex *index.ClusteringIndex

func sign() float64 {
	if rand.Float64() > 0.5 {
		return 1
	}
	return -1
}

func getIndex(context context.Context) *index.ClusteringIndex {
	if geoindex == nil {
		geoindex = index.NewClusteringIndex()

		capitals := worldCapitals(context)
		id := 1

		for _, capital := range capitals {
			for i := 0; i < 300; i++ {
				id++

				geoindex.Add(index.NewGeoPoint(
					fmt.Sprintf("%d", id),
					capital.Lat()+rand.Float64()/6.0*sign(),
					capital.Lon()+rand.Float64()/6.0*sign(),
				))
			}
		}
	}

	return geoindex
}

func worldCapitals(context context.Context) []index.Point {
	file, err := os.OpenFile("static/capitals.csv", os.O_RDONLY, 0)

	if err != nil {
		log.Default().Printf("Unable to open capitals.csv: %v", err)
		return make([]index.Point, 0)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'

	records, _ := reader.ReadAll()
	capitals := make([]index.Point, 0)

	for _, record := range records {
		id := record[0]
		lat, _ := strconv.ParseFloat(record[3], 64)
		lon, _ := strconv.ParseFloat(record[4], 64)

		capital := index.NewGeoPoint(id, lat, lon)
		capitals = append(capitals, capital)
	}

	return capitals
}

func points(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	topLeftLat, _ := strconv.ParseFloat(r.Form["topLeftLat"][0], 64)
	topLeftLon, _ := strconv.ParseFloat(r.Form["topLeftLon"][0], 64)
	bottomRightLat, _ := strconv.ParseFloat(r.Form["bottomRightLat"][0], 64)
	bottomRightLon, _ := strconv.ParseFloat(r.Form["bottomRightLon"][0], 64)

	c := context.Background()
	visiblePoints := getIndex(c).Range(index.NewGeoPoint("topLeft", topLeftLat, topLeftLon), index.NewGeoPoint("bottomRight", bottomRightLat, bottomRightLon))

	data, _ := json.Marshal(visiblePoints)
	fmt.Fprintln(w, string(data))
}

func knearest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	lat, _ := strconv.ParseFloat(r.Form["lat"][0], 64)
	lon, _ := strconv.ParseFloat(r.Form["lon"][0], 64)
	k, _ := strconv.ParseInt(r.Form["k"][0], 10, 32)

	c := context.Background()
	nearest := getIndex(c).KNearest(index.NewGeoPoint("query", lat, lon), int(k), index.Km(5), func(_ index.Point) bool { return true })
	data, _ := json.Marshal(nearest)
	fmt.Fprintln(w, string(data))
}
