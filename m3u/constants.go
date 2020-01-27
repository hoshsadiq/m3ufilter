package m3u

import "strings"

var definitionReplaces = `SD|HD|FHD`

// todo plenty of options are missing here
// this gets additionally populated with the keys and values of countryOverrides below it
var countryReplaces = `LAT|RO|AF|ARB|IT|DE|PH|IRE`

var countryOverrides = map[string]string{
	"GB":  "UK",
	"GBR": "UK",
	"CAN": "CA",
	"USA": "US",
	"NLD": "NL",
	"BRA": "BR",
	"PRT": "PT",
	"ESP": "ES",
	"POL": "PL",
	"PAK": "PK",
	"IND": "IN",
	"FRA": "FR",
	"AUS": "AU",
}

func init() {
	countriesAdded := map[string]bool{}
	extraCountries := make([]string, 0, len(countryOverrides)*2)

	for alpha3, alpha2 := range countryOverrides {
		if _, ok := countriesAdded[alpha2]; !ok {
			extraCountries = append(extraCountries, alpha2)
			countriesAdded[alpha2] = true
		}
		if _, ok := countriesAdded[alpha3]; !ok {
			extraCountries = append(extraCountries, alpha3)
			countriesAdded[alpha3] = true
		}
	}
	countryReplaces = countryReplaces + "|" + strings.Join(extraCountries, "|")
}
