package satcat

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

type SatelliteEntry struct {
	Designation       string
	CatalogNumber     int64
	PayloadFlag       bool
	OperationalStatus string
	Name              string
	Ownership         string

	LaunchDate string
	LaunchSite string

	DecayDate string

	Apogee      int64
	Perigee     int64
	Period      float64
	Inclination float64

	IsDebris bool
	IsFirstComponent bool
	HasDecayed bool
}

func (entry SatelliteEntry) Print() {
	fmt.Printf("%s : satellite %s (%d - %s) launched by %s from %s to %d/%d@%f %f\n",
		entry.LaunchDate,
		entry.Designation,
		entry.CatalogNumber,
		entry.Name,
		entry.Ownership,
		entry.LaunchSite,
		entry.Perigee,
		entry.Apogee,
		entry.Inclination,
		entry.Period)
}

func (entry *SatelliteEntry) FromString(line string) {
	var err error

	// 1957-001A    00001   D SL-1 R/B                  CIS    1957-10-04  TYMSC  1957-12-01     96.2   65.1     938     214   20.4200
	entry.Designation = strings.Trim(line[0:10], " ")

	entry.CatalogNumber, err = strconv.ParseInt(strings.Trim(line[13:17], " "), 10, 64)
	checkError(err)

	entry.PayloadFlag = line[20:21] == "*"

	entry.OperationalStatus = line[21:22]
	entry.Name = strings.Trim(line[23:46], " ")
	entry.Ownership = strings.Trim(line[49:53], " ")

	entry.LaunchDate = strings.Trim(line[56:66], " ")
	entry.LaunchSite = strings.Trim(line[68:72], " ")

	entry.Perigee, err = strconv.ParseInt(strings.Trim(line[111:116], " "), 10, 64)
	if err != nil {
		entry.Perigee = 0
	}
	entry.Apogee, err = strconv.ParseInt(strings.Trim(line[103:108], " "), 10, 64)
	if err != nil {
		entry.Apogee = 0
	}
	entry.Inclination, err = strconv.ParseFloat(strings.Trim(line[96:100], " "), 64)
	if err != nil {
		entry.Inclination = 0
	}

	debrisRegex := regexp.MustCompile("\\bDEB\\b")
	entry.IsDebris = debrisRegex.MatchString(entry.Name)

	entry.IsFirstComponent, _ = regexp.MatchString("\\dA$", entry.Designation)
	entry.HasDecayed = line[21:22] == "D"
}
