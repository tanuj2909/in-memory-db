package cmd

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/tanuj2909/in-memory-db/app/types"
)

const earthRadius = 6372797.560856 // meters

func GEOADD(server *types.ServerState, args ...string) []byte {
	if len(args) < 4 || (len(args)-1)%3 != 0 {
		return respHandler.Error.Encode("ERR wrong number of arguments for 'geoadd' command")
	}

	key := args[0]
	added := 0

	if _, ok := server.GeoSets[key]; !ok {
		server.GeoSets[key] = map[string]types.GeoPoint{}
	}

	for i := 1; i < len(args); i += 3 {
		lon, err := strconv.ParseFloat(args[i], 64)
		if err != nil || lon < -180 || lon > 180 {
			return respHandler.Error.Encode("ERR invalid longitude")
		}
		lat, err := strconv.ParseFloat(args[i+1], 64)
		if err != nil || lat < -85.05112878 || lat > 85.05112878 {
			return respHandler.Error.Encode("ERR invalid latitude")
		}
		member := args[i+2]

		if _, exists := server.GeoSets[key][member]; !exists {
			added++
		}
		server.GeoSets[key][member] = types.GeoPoint{Lon: lon, Lat: lat}
	}

	return respHandler.Integer.Encode(added)
}

func GEOPOS(server *types.ServerState, args ...string) []byte {
	if len(args) < 2 {
		return respHandler.Error.Encode("ERR wrong number of arguments for 'geopos' command")
	}

	key := args[0]
	members := args[1:]
	geoSet := server.GeoSets[key]

	res := []byte(fmt.Sprintf("*%d\r\n", len(members)))
	for _, member := range members {
		if geoSet == nil {
			res = append(res, []byte("*-1\r\n")...)
			continue
		}
		p, ok := geoSet[member]
		if !ok {
			res = append(res, []byte("*-1\r\n")...)
			continue
		}
		lonStr := strconv.FormatFloat(p.Lon, 'f', 17, 64)
		latStr := strconv.FormatFloat(p.Lat, 'f', 17, 64)
		pair := fmt.Sprintf("*2\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n",
			len(lonStr), lonStr, len(latStr), latStr)
		res = append(res, []byte(pair)...)
	}
	return res
}

func GEODIST(server *types.ServerState, args ...string) []byte {
	if len(args) < 3 || len(args) > 4 {
		return respHandler.Error.Encode("ERR wrong number of arguments for 'geodist' command")
	}

	key, member1, member2 := args[0], args[1], args[2]
	unit := "m"
	if len(args) == 4 {
		unit = strings.ToLower(args[3])
	}

	geoSet, ok := server.GeoSets[key]
	if !ok {
		return respHandler.Null.Encode()
	}
	p1, ok1 := geoSet[member1]
	p2, ok2 := geoSet[member2]
	if !ok1 || !ok2 {
		return respHandler.Null.Encode()
	}

	distMeters := haversine(p1.Lat, p1.Lon, p2.Lat, p2.Lon)
	dist := fromMeters(distMeters, unit)
	return respHandler.BulkString.Encode(strconv.FormatFloat(dist, 'f', 4, 64))
}

// GEOSEARCH key FROMMEMBER member | FROMLONLAT lon lat
//
//	BYRADIUS radius m|km|mi|ft
//	[ASC|DESC] [COUNT count] [WITHCOORD] [WITHDIST]
func GEOSEARCH(server *types.ServerState, args ...string) []byte {
	if len(args) < 5 {
		return respHandler.Error.Encode("ERR wrong number of arguments for 'geosearch' command")
	}

	key := args[0]
	args = args[1:]

	var centerLat, centerLon float64
	switch strings.ToUpper(args[0]) {
	case "FROMMEMBER":
		if len(args) < 2 {
			return respHandler.Error.Encode("ERR syntax error")
		}
		geoSet, ok := server.GeoSets[key]
		if !ok {
			return respHandler.Array.Encode([]string{})
		}
		p, ok := geoSet[args[1]]
		if !ok {
			return respHandler.Error.Encode("ERR could not perform this operation on a key that doesn't exist")
		}
		centerLat, centerLon = p.Lat, p.Lon
		args = args[2:]
	case "FROMLONLAT":
		if len(args) < 3 {
			return respHandler.Error.Encode("ERR syntax error")
		}
		var err error
		centerLon, err = strconv.ParseFloat(args[1], 64)
		if err != nil {
			return respHandler.Error.Encode("ERR invalid longitude")
		}
		centerLat, err = strconv.ParseFloat(args[2], 64)
		if err != nil {
			return respHandler.Error.Encode("ERR invalid latitude")
		}
		args = args[3:]
	default:
		return respHandler.Error.Encode("ERR syntax error")
	}

	if len(args) < 3 || strings.ToUpper(args[0]) != "BYRADIUS" {
		return respHandler.Error.Encode("ERR syntax error: expected BYRADIUS")
	}
	radius, err := strconv.ParseFloat(args[1], 64)
	if err != nil || radius < 0 {
		return respHandler.Error.Encode("ERR invalid radius")
	}
	unit := strings.ToLower(args[2])
	radiusMeters := toMeters(radius, unit)
	args = args[3:]

	order := "ASC"
	count := -1
	withDist := false
	withCoord := false

	for i := 0; i < len(args); {
		switch strings.ToUpper(args[i]) {
		case "ASC":
			order = "ASC"
			i++
		case "DESC":
			order = "DESC"
			i++
		case "COUNT":
			if i+1 >= len(args) {
				return respHandler.Error.Encode("ERR syntax error")
			}
			count, err = strconv.Atoi(args[i+1])
			if err != nil {
				return respHandler.Error.Encode("ERR value is not an integer or out of range")
			}
			i += 2
		case "WITHDIST":
			withDist = true
			i++
		case "WITHCOORD":
			withCoord = true
			i++
		default:
			return respHandler.Error.Encode(fmt.Sprintf("ERR syntax error: unknown option '%s'", args[i]))
		}
	}

	geoSet, ok := server.GeoSets[key]
	if !ok {
		return respHandler.Array.Encode([]string{})
	}

	type geoResult struct {
		member string
		dist   float64
		lon    float64
		lat    float64
	}

	var results []geoResult
	for member, p := range geoSet {
		dist := haversine(centerLat, centerLon, p.Lat, p.Lon)
		if dist <= radiusMeters {
			results = append(results, geoResult{
				member: member,
				dist:   fromMeters(dist, unit),
				lon:    p.Lon,
				lat:    p.Lat,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if order == "DESC" {
			return results[i].dist > results[j].dist
		}
		return results[i].dist < results[j].dist
	})

	if count >= 0 && count < len(results) {
		results = results[:count]
	}

	if !withDist && !withCoord {
		members := make([]string, len(results))
		for i, r := range results {
			members[i] = r.member
		}
		return respHandler.Array.Encode(members)
	}

	outer := []byte(fmt.Sprintf("*%d\r\n", len(results)))
	for _, r := range results {
		fieldCount := 1
		if withDist {
			fieldCount++
		}
		if withCoord {
			fieldCount++
		}
		inner := []byte(fmt.Sprintf("*%d\r\n", fieldCount))
		inner = append(inner, respHandler.BulkString.Encode(r.member)...)
		if withDist {
			inner = append(inner, respHandler.BulkString.Encode(strconv.FormatFloat(r.dist, 'f', 4, 64))...)
		}
		if withCoord {
			lonStr := strconv.FormatFloat(r.lon, 'f', 17, 64)
			latStr := strconv.FormatFloat(r.lat, 'f', 17, 64)
			coord := fmt.Sprintf("*2\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n", len(lonStr), lonStr, len(latStr), latStr)
			inner = append(inner, []byte(coord)...)
		}
		outer = append(outer, inner...)
	}
	return outer
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const deg2rad = math.Pi / 180
	dLat := (lat2 - lat1) * deg2rad
	dLon := (lon2 - lon1) * deg2rad
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*deg2rad)*math.Cos(lat2*deg2rad)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	return earthRadius * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

func toMeters(dist float64, unit string) float64 {
	switch unit {
	case "km":
		return dist * 1000
	case "mi":
		return dist * 1609.344
	case "ft":
		return dist * 0.3048
	default:
		return dist
	}
}

func fromMeters(dist float64, unit string) float64 {
	switch unit {
	case "km":
		return dist / 1000
	case "mi":
		return dist / 1609.344
	case "ft":
		return dist / 0.3048
	default:
		return dist
	}
}
