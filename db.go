package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func GetPlacesFromDB(DB *sql.DB, location Location) ([]Place, error) {
	var places []Place

	sqlStatement := fmt.Sprintf("SELECT googleid, name FROM places WHERE location <@> POINT(%v, %v) < 0.5/1.6;", location.Longitude, location.Latitude)
	rows, err := DB.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	count := 0

	defer rows.Close()
	for rows.Next() {
		if count < 3 {
			var place Place
			if err := rows.Scan(&place.ID, &place.Name); err != nil {
				return nil, err
			}
			places = append(places, place)
		}
		count++

	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return places, nil
}
