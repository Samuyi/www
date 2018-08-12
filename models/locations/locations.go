package locations

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	validator "github.com/asaskevich/govalidator"
	_ "github.com/lib/pq" // postgres driver
)

//Location where an item is based
type Location struct {
	LocationID  string    `json:"location_id"`
	City        string    `json:"city"`
	UserID      string    `json:"user_id,omitempty"`
	State       string    `json:"state"`
	Country     string    `json:"country"`
	CountryCode string    `json:"country_code"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

var countries = map[string]string{
	"AF": "Afghanistan",
	"AL": "Albania",
	"DZ": "Algeria",
	"AS": "American Samoa",
	"AD": "Andorra",
	"AO": "Angola",
	"AI": "Anguilla",
	"AQ": "Antarctica",
	"AG": "Antigua And Barbuda",
	"AR": "Argentina",
	"AM": "Armenia",
	"AW": "Aruba",
	"AU": "Australia",
	"AT": "Austria",
	"AZ": "Azerbaijan",
	"BS": "Bahamas",
	"BH": "Bahrain",
	"BD": "Bangladesh",
	"BB": "Barbados",
	"BY": "Belarus",
	"BE": "Belgium",
	"BZ": "Belize",
	"BJ": "Benin",
	"BM": "Bermuda",
	"BT": "Bhutan",
	"BO": "Bolivia",
	"BA": "Bosnia And Herzegovina",
	"BW": "Botswana",
	"BV": "Bouvet Island",
	"BR": "Brazil",
	"IO": "British Indian Ocean Territory",
	"BN": "Brunei Darussalam",
	"BG": "Bulgaria",
	"BF": "Burkina Faso",
	"BI": "Burundi",
	"KH": "Cambodia",
	"CM": "Cameroon",
	"CA": "Canada",
	"CV": "Cape Verde",
	"KY": "Cayman Islands",
	"CF": "Central African Republic",
	"TD": "Chad",
	"CL": "Chile",
	"CN": "China",
	"CX": "Christmas Island",
	"CC": "Cocos (keeling) Islands",
	"CO": "Colombia",
	"KM": "Comoros",
	"CG": "Congo",
	"CD": "Congo, The Democratic Republic Of The",
	"CK": "Cook Islands",
	"CR": "Costa Rica",
	"CI": "Cote D'ivoire",
	"HR": "Croatia",
	"CU": "Cuba",
	"CY": "Cyprus",
	"CZ": "Czech Republic",
	"DK": "Denmark",
	"DJ": "Djibouti",
	"DM": "Dominica",
	"DO": "Dominican Republic",
	"TP": "East Timor",
	"EC": "Ecuador",
	"EG": "Egypt",
	"SV": "El Salvador",
	"GQ": "Equatorial Guinea",
	"ER": "Eritrea",
	"EE": "Estonia",
	"ET": "Ethiopia",
	"FK": "Falkland Islands (malvinas)",
	"FO": "Faroe Islands",
	"FJ": "Fiji",
	"FI": "Finland",
	"FR": "France",
	"GF": "French Guiana",
	"PF": "French Polynesia",
	"TF": "French Southern Territories",
	"GA": "Gabon",
	"GM": "Gambia",
	"GE": "Georgia",
	"DE": "Germany",
	"GH": "Ghana",
	"GI": "Gibraltar",
	"GR": "Greece",
	"GL": "Greenland",
	"GD": "Grenada",
	"GP": "Guadeloupe",
	"GU": "Guam",
	"GT": "Guatemala",
	"GN": "Guinea",
	"GW": "Guinea-bissau",
	"GY": "Guyana",
	"HT": "Haiti",
	"HM": "Heard Island And Mcdonald Islands",
	"VA": "Holy See (vatican City State)",
	"HN": "Honduras",
	"HK": "Hong Kong",
	"HU": "Hungary",
	"IS": "Iceland",
	"IN": "India",
	"ID": "Indonesia",
	"IR": "Iran, Islamic Republic Of",
	"IQ": "Iraq",
	"IE": "Ireland",
	"IL": "Israel",
	"IT": "Italy",
	"JM": "Jamaica",
	"JP": "Japan",
	"JO": "Jordan",
	"KZ": "Kazakstan",
	"KE": "Kenya",
	"KI": "Kiribati",
	"KP": "Korea, Democratic People's Republic Of",
	"KR": "Korea, Republic Of",
	"KV": "Kosovo",
	"KW": "Kuwait",
	"KG": "Kyrgyzstan",
	"LA": "Lao People's Democratic Republic",
	"LV": "Latvia",
	"LB": "Lebanon",
	"LS": "Lesotho",
	"LR": "Liberia",
	"LY": "Libyan Arab Jamahiriya",
	"LI": "Liechtenstein",
	"LT": "Lithuania",
	"LU": "Luxembourg",
	"MO": "Macau",
	"MK": "Macedonia, The Former Yugoslav Republic Of",
	"MG": "Madagascar",
	"MW": "Malawi",
	"MY": "Malaysia",
	"MV": "Maldives",
	"ML": "Mali",
	"MT": "Malta",
	"MH": "Marshall Islands",
	"MQ": "Martinique",
	"MR": "Mauritania",
	"MU": "Mauritius",
	"YT": "Mayotte",
	"MX": "Mexico",
	"FM": "Micronesia, Federated States Of",
	"MD": "Moldova, Republic Of",
	"MC": "Monaco",
	"MN": "Mongolia",
	"MS": "Montserrat",
	"ME": "Montenegro",
	"MA": "Morocco",
	"MZ": "Mozambique",
	"MM": "Myanmar",
	"NA": "Namibia",
	"NR": "Nauru",
	"NP": "Nepal",
	"NL": "Netherlands",
	"AN": "Netherlands Antilles",
	"NC": "New Caledonia",
	"NZ": "New Zealand",
	"NI": "Nicaragua",
	"NE": "Niger",
	"NG": "Nigeria",
	"NU": "Niue",
	"NF": "Norfolk Island",
	"MP": "Northern Mariana Islands",
	"NO": "Norway",
	"OM": "Oman",
	"PK": "Pakistan",
	"PW": "Palau",
	"PS": "Palestinian Territory, Occupied",
	"PA": "Panama",
	"PG": "Papua New Guinea",
	"PY": "Paraguay",
	"PE": "Peru",
	"PH": "Philippines",
	"PN": "Pitcairn",
	"PL": "Poland",
	"PT": "Portugal",
	"PR": "Puerto Rico",
	"QA": "Qatar",
	"RE": "Reunion",
	"RO": "Romania",
	"RU": "Russian Federation",
	"RW": "Rwanda",
	"SH": "Saint Helena",
	"KN": "Saint Kitts And Nevis",
	"LC": "Saint Lucia",
	"PM": "Saint Pierre And Miquelon",
	"VC": "Saint Vincent And The Grenadines",
	"WS": "Samoa",
	"SM": "San Marino",
	"ST": "Sao Tome And Principe",
	"SA": "Saudi Arabia",
	"SN": "Senegal",
	"RS": "Serbia",
	"SC": "Seychelles",
	"SL": "Sierra Leone",
	"SG": "Singapore",
	"SK": "Slovakia",
	"SI": "Slovenia",
	"SB": "Solomon Islands",
	"SO": "Somalia",
	"ZA": "South Africa",
	"GS": "South Georgia And The South Sandwich Islands",
	"ES": "Spain",
	"LK": "Sri Lanka",
	"SD": "Sudan",
	"SR": "Suriname",
	"SJ": "Svalbard And Jan Mayen",
	"SZ": "Swaziland",
	"SE": "Sweden",
	"CH": "Switzerland",
	"SY": "Syrian Arab Republic",
	"TW": "Taiwan, Province Of China",
	"TJ": "Tajikistan",
	"TZ": "Tanzania, United Republic Of",
	"TH": "Thailand",
	"TG": "Togo",
	"TK": "Tokelau",
	"TO": "Tonga",
	"TT": "Trinidad And Tobago",
	"TN": "Tunisia",
	"TR": "Turkey",
	"TM": "Turkmenistan",
	"TC": "Turks And Caicos Islands",
	"TV": "Tuvalu",
	"UG": "Uganda",
	"UA": "Ukraine",
	"AE": "United Arab Emirates",
	"GB": "United Kingdom",
	"US": "United States",
	"UM": "United States Minor Outlying Islands",
	"UY": "Uruguay",
	"UZ": "Uzbekistan",
	"VU": "Vanuatu",
	"VE": "Venezuela",
	"VN": "VietNam",
	"VG": "Virgin Islands, British",
	"VI": "Virgin Islands, U.s.",
	"WF": "Wallis And Futuna",
	"EH": "Western Sahara",
	"YE": "Yemen",
	"ZM": "Zambia",
	"ZW": "Zimbabwe",
}

var db *sql.DB

const (
	host     = "localhost"
	port     = 5432
	user     = "help"
	password = "help"
	dbname   = "help.ng"
)

func init() {
	var err error

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Println(err)
	}
	err = db.Ping()

	if err != nil {
		log.Println(err)
	}

	log.Println("connected to database")
}

//Validate location struct
func (location *Location) Validate() map[string]string {
	var errors = make(map[string]string)

	if len(location.City) <= 2 {
		message := "Please supply a valid city"
		errors["Invalid City"] = message
	}

	if len(location.State) <= 2 {
		message := "Please supply a valid state"
		errors["Invalid State"] = message
	}

	if _, ok := countries[location.CountryCode]; !ok {
		message := "Please supply a valid country code"
		errors["Invalid Country code"] = message
	}

	if !validator.IsUUID(location.UserID) {
		message := "Please supply a valid user id"
		errors["Invalid user Id"] = message
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

//Create a location
func (location *Location) Create() error {
	query := "INSERT INTO locations (city, user_id, state, country) VALUES ($1, $2, $3, $4) returning location_id"

	stmt, err := db.Prepare(query)

	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return err
	}

	location.City = strings.ToUpper(location.City)
	location.State = strings.ToUpper(location.State)
	location.Country = countries[location.CountryCode]

	err = stmt.QueryRow(location.City, location.UserID, location.State, location.Country).Scan(&location.LocationID)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

//Get a location
func (location *Location) Get() error {
	query := "SELECT city, state, country, created_at FROM locations where location_id = $1"

	stmt, err := db.Prepare(query)

	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return err
	}

	err = stmt.QueryRow(location.LocationID).Scan(&location.City, &location.State, &location.Country, &location.CreatedAt)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}

//Update a location
func (location *Location) Update(changes map[string]string) error {
	var query bytes.Buffer
	query.Write([]byte("UPDATE locations SET"))
	for k, v := range changes {
		query.Write([]byte(" "))
		query.Write([]byte(k))
		query.Write([]byte(" ="))
		query.Write([]byte(" "))
		query.Write([]byte("'" + v + "'"))
		query.Write([]byte(", "))
	}

	query.Write([]byte(" updated_at = $1 WHERE location_id = $2"))

	stmt, err := db.Prepare(query.String())

	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return err
	}
	location.UpdatedAt = time.Now()

	_, err = stmt.Exec(location.UpdatedAt, location.LocationID)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

//Delete a location
func (location *Location) Delete() error {
	query := "DELETE FROM locations WHERE location_id = $1"

	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return err
	}

	_, err = stmt.Exec(location.LocationID)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

//GetAll locations
func (location *Location) GetAll() ([]Location, error) {
	query := "SELECT location_id, city, state, country FROM locations"

	stmt, err := db.Prepare(query)

	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	rows, err := stmt.Query()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var locations []Location

	defer rows.Close()
	for rows.Next() {
		var location Location

		if err := rows.Scan(&location.LocationID, &location.City, &location.State, &location.Country); err != nil {
			log.Println(err)
			return nil, err
		}
		locations = append(locations, location)
	}

	return locations, nil
}
