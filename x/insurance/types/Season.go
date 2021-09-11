package types

import (
	"time"
)

type Season string

const (
	Spring = "SPRING"
	Summer = "SUMMER"
	Fall   = "FALL"
	Winter = "WINTER"
)

func Now() Season {
	today := time.Now()

	for year := 2021; ; year++ {
		startOfNextYear := time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.UTC)

		if today.Before(startOfNextYear) {
			startOfSpring := time.Date(year, time.March, 21, 0, 0, 0, 0, time.UTC)
			if today.Before(startOfSpring) {
				return Winter
			}

			startOfSummer := time.Date(year, time.June, 21, 0, 0, 0, 0, time.UTC)
			if today.Before(startOfSummer) {
				return Spring
			}

			startOfFall := time.Date(year, time.September, 23, 0, 0, 0, 0, time.UTC)
			if today.Before(startOfFall) {
				return Summer
			}

			startOfWinter := time.Date(year, time.December, 21, 0, 0, 0, 0, time.UTC)
			if today.Before(startOfWinter) {
				return Fall
			}

			return Winter
		}
	}
}
