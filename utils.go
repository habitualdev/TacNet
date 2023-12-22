package main

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

const r2m float64 = 1018.591636

func PolarToGrid(distance, azimuth, gunPos string) (string, error) {

	distanceFloat, err := strconv.ParseFloat(distance, 64)

	az, err := strconv.Atoi(azimuth)
	if err != nil {
		return "", err
	}
	sinDelta := math.Sin(float64(az)/r2m) * distanceFloat
	cosDelta := math.Cos(float64(az)/r2m) * distanceFloat

	switch {
	case len(gunPos) == 6:
		easting := gunPos[:3]
		northing := gunPos[3:]
		gunPos = easting + "00" + northing + "00"
	case len(gunPos) == 8:
		easting := gunPos[:4]
		northing := gunPos[4:]
		gunPos = easting + "0" + northing + "0"
	case len(gunPos) == 10:
		break
	default:
		return "", errors.New("Check your gun position grid.")
	}

	quadrant := 0

	switch {
	case az < 1600:
		quadrant = 1
	case az < 3200:
		quadrant = 2
	case az < 4800:
		quadrant = 3
	case az < 6400:
		quadrant = 4
	}

	eastingGun := gunPos[:5]
	northingGun := gunPos[5:]

	eastingGunInt, err := strconv.Atoi(eastingGun)
	northingGunInt, err := strconv.Atoi(northingGun)

	targetNorthing := 0
	targetEasting := 0

	switch {
	case quadrant == 1:
		targetNorthing = northingGunInt + int(cosDelta)
		targetEasting = eastingGunInt + int(sinDelta)
	case quadrant == 2:
		targetNorthing = northingGunInt + int(cosDelta)
		targetEasting = eastingGunInt + int(sinDelta)
	case quadrant == 3:
		targetNorthing = northingGunInt + int(cosDelta)
		targetEasting = eastingGunInt + int(sinDelta)
	case quadrant == 4:
		targetNorthing = northingGunInt + int(cosDelta)
		targetEasting = eastingGunInt + int(sinDelta)
	}

	targetEastingString := fmt.Sprintf("%d", targetEasting)
	if len(targetEastingString) < 5 {
		for len(targetEastingString) != 5 {
			targetEastingString = "0" + targetEastingString
		}
	}

	targetNorthingString := fmt.Sprintf("%d", targetNorthing)
	if len(targetNorthingString) < 5 {
		for len(targetNorthingString) != 5 {
			targetNorthingString = "0" + targetNorthingString
		}
	}

	targetGrid := targetEastingString + targetNorthingString

	return targetGrid, nil
}
