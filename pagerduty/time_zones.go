package pagerduty

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

// TimeZone is a thin wrapper around time.Location to can marshal and unmarshal
// to the custom PagerDuty TimeZone type.
type TimeZone struct {
	*time.Location
}

func (tz TimeZone) EncodeValues(key string, v *url.Values) error {
	location := tz.Location
	timeZone, ok := ianaToPagerDuty[location.String()]
	if !ok {
		return fmt.Errorf("time zone %q does not exist", location.String())
	}

	v.Add(key, timeZone)
	return nil
}

func (tz *TimeZone) MarshalJSON() ([]byte, error) {
	location := tz.Location.String()
	timeZone, ok := ianaToPagerDuty[location]
	if !ok {
		return []byte{}, fmt.Errorf("time zone %q does not exist", location)
	}

	return []byte(`"` + timeZone + `"`), nil
}

func (tz *TimeZone) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	timeZone, ok := pagerdutyToIANA[s]
	if !ok {
		return fmt.Errorf("time zone %q does not exist", s)
	}

	location, err := time.LoadLocation(timeZone)
	if err != nil {
		return err
	}

	tz.Location = location

	return nil
}

// reverseMap is a helper function to reverse an existing map.
func reverseMap(m map[string]string) map[string]string {
	r := make(map[string]string, len(m))
	for k, v := range m {
		r[v] = k
	}
	return r
}

// Mapping from the custom PagerDuty TimeZone type to the IANA time zone names.
var pagerdutyToIANA = map[string]string{
	"Abu Dhabi":                    "Asia/Dubai",
	"Adelaide":                     "Australia/Adelaide",
	"Alaska":                       "America/Juneau",
	"Almaty":                       "Asia/Almaty",
	"Amsterdam":                    "Europe/Amsterdam",
	"Arizona":                      "America/Phoenix",
	"Astana":                       "Asia/Thimphu",
	"Athens":                       "Europe/Athens",
	"Atlantic Time (Canada)":       "America/Halifax",
	"Auckland":                     "Pacific/Auckland",
	"Azores":                       "Atlantic/Azores",
	"Baghdad":                      "Asia/Baghdad",
	"Baku":                         "Asia/Baku",
	"Bangkok":                      "Asia/Bangkok",
	"Beijing":                      "Asia/Shanghai",
	"Belgrade":                     "Europe/Belgrade",
	"Berlin":                       "Europe/Berlin",
	"Bern":                         "CET",
	"Bogota":                       "America/Bogota",
	"Brasilia":                     "America/Sao_Paulo",
	"Bratislava":                   "Europe/Bratislava",
	"Brisbane":                     "Australia/Brisbane",
	"Brussels":                     "Europe/Brussels",
	"Bucharest":                    "Europe/Bucharest",
	"Budapest":                     "Europe/Budapest",
	"Buenos Aires":                 "America/Argentina/Buenos_Aires",
	"Cairo":                        "Africa/Cairo",
	"Canberra":                     "Australia/Canberra",
	"Cape Verde Is.":               "Atlantic/Cape_Verde",
	"Caracas":                      "America/Caracas",
	"Casablanca":                   "Africa/Casablanca",
	"Central America":              "America/Guatemala",
	"Central Time (US & Canada)":   "America/Chicago",
	"Chennai":                      "Asia/Kolkata",
	"Chihuahua":                    "America/Chihuahua",
	"Chongqing":                    "Asia/Chongqing",
	"Copenhagen":                   "Europe/Copenhagen",
	"Darwin":                       "Australia/Darwin",
	"Dhaka":                        "Asia/Dhaka",
	"Dublin":                       "Europe/Dublin",
	"Eastern Time (US & Canada)":   "America/New_York",
	"Edinburgh":                    "Europe/Dublin",
	"Ekaterinburg":                 "Asia/Yekaterinburg",
	"Fiji":                         "Pacific/Fiji",
	"Georgetown":                   "America/Argentina/San_Juan",
	"Greenland":                    "America/Godthab",
	"Guadalajara":                  "America/Mexico_City",
	"Guam":                         "Pacific/Guam",
	"Hanoi":                        "Asia/Bangkok",
	"Harare":                       "Africa/Harare",
	"Hawaii":                       "Pacific/Honolulu",
	"Helsinki":                     "Europe/Helsinki",
	"Hobart":                       "Australia/Hobart",
	"Hong Kong":                    "Asia/Hong_Kong",
	"Indiana (East)":               "America/Indiana/Indianapolis",
	"International Date Line West": "Pacific/Midway",
	"Irkutsk":                      "Asia/Irkutsk",
	"Islamabad":                    "Asia/Karachi",
	"Istanbul":                     "Europe/Istanbul",
	"Jakarta":                      "Asia/Jakarta",
	"Jerusalem":                    "Asia/Jerusalem",
	"Kabul":                        "Asia/Kabul",
	"Kamchatka":                    "Asia/Kamchatka",
	"Karachi":                      "Asia/Karachi",
	"Kathmandu":                    "Asia/Katmandu",
	"Kolkata":                      "Asia/Kolkata",
	"Krasnoyarsk":                  "Asia/Krasnoyarsk",
	"Kuala Lumpur":                 "Asia/Kuala_Lumpur",
	"Kuwait":                       "Asia/Kuwait",
	"Kyev":                         "Europe/Kiev",
	"La Paz":                       "America/La_Paz",
	"Lima":                         "America/Lima",
	"Lisbon":                       "Europe/Lisbon",
	"Ljubljana":                    "Europe/Ljubljana",
	"London":                       "Europe/London",
	"Madrid":                       "Europe/Madrid",
	"Magadan":                      "Asia/Magadan",
	"Marshall Is.":                 "Pacific/Majuro",
	"Mazatlan":                     "America/Mazatlan",
	"Melbourne":                    "Australia/Melbourne",
	"Mexico City":                  "America/Mexico_City",
	"Mid-Atlantic":                 "Atlantic/South_Georgia",
	"Midway Island":                "Pacific/Midway",
	"Minsk":                        "Europe/Minsk",
	"Monrovia":                     "Africa/Monrovia",
	"Monterrey":                    "America/Monterrey",
	"Moscow":                       "Europe/Moscow",
	"Mountain Time (US & Canada)":  "America/Denver",
	"Mumbai":                       "Asia/Kolkata",
	"Muscat":                       "Asia/Muscat",
	"Nairobi":                      "Africa/Nairobi",
	"New Caledonia":                "Pacific/Noumea",
	"New Delhi":                    "Asia/Kolkata",
	"Newfoundland":                 "America/St_Johns",
	"Novosibirsk":                  "Asia/Novosibirsk",
	"Nuku'alofa":                   "Pacific/Tongatapu",
	"Osaka":                        "Asia/Tokyo",
	"Pacific Time (US & Canada)": "America/Los_Angeles",
	"Paris":               "Europe/Paris",
	"Perth":               "Australia/Perth",
	"Port Moresby":        "Pacific/Port_Moresby",
	"Prague":              "Europe/Prague",
	"Pretoria":            "Africa/Johannesburg",
	"Quito":               "America/Lima",
	"Rangoon":             "Asia/Rangoon",
	"Riga":                "Europe/Riga",
	"Riyadh":              "Asia/Riyadh",
	"Rome":                "Europe/Rome",
	"Samoa":               "Pacific/Samoa",
	"Santiago":            "America/Santiago",
	"Sapporo":             "Asia/Tokyo",
	"Sarajevo":            "Europe/Sarajevo",
	"Saskatchewan":        "Canada/Saskatchewan",
	"Seoul":               "Asia/Seoul",
	"Singapore":           "Asia/Singapore",
	"Skopje":              "Europe/Skopje",
	"Sofia":               "Europe/Sofia",
	"Solomon Is.":         "Pacific/Guadalcanal",
	"Sri Jayawardenepura": "Asia/Colombo",
	"St. Petersburg":      "Europe/Moscow",
	"Stockholm":           "Europe/Stockholm",
	"Sydney":              "Australia/Sydney",
	"Taipei":              "Asia/Taipei",
	"Tallinn":             "Europe/Tallinn",
	"Tashkent":            "Asia/Tashkent",
	"Tbilisi":             "Asia/Tbilisi",
	"Tehran":              "Asia/Tehran",
	"Tijuana":             "America/Tijuana",
	"Tokyo":               "Asia/Tokyo",
	"UTC":                 "UTC",
	"Ulaan Bataar":        "Asia/Ulaanbaatar",
	"Urumqi":              "Asia/Urumqi",
	"Vienna":              "Europe/Vienna",
	"Vilnius":             "Europe/Vilnius",
	"Vladivostok":         "Asia/Vladivostok",
	"Volgograd":           "Europe/Moscow",
	"Warsaw":              "Europe/Warsaw",
	"Wellington":          "Pacific/Auckland",
	"West Central Africa": "Africa/Algiers",
	"Yakutsk":             "Asia/Yakutsk",
	"Yerevan":             "Asia/Yerevan",
	"Zagreb":              "Europe/Zagreb",
}

// Mapping from the IANA time zone names to the custom PagerDuty TimeZone type.
var ianaToPagerDuty = reverseMap(pagerdutyToIANA)
