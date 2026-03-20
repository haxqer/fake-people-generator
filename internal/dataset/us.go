package dataset

type WeightedValue struct {
	Value  string
	Weight int
}

type City struct {
	Name       string
	PostalCode string
	AreaCodes  []string
	Weight     int
}

type State struct {
	Code   string
	Name   string
	Weight int
	Cities []City
}

type Data struct {
	States          []State
	StatesByCode    map[string]State
	MaleFirstNames  []WeightedValue
	FemaleFirstName []WeightedValue
	LastNames       []WeightedValue
	StreetNames     []WeightedValue
	StreetSuffixes  []WeightedValue
}

var usData = Data{
	States: []State{
		{Code: "AL", Name: "Alabama", Weight: 5, Cities: []City{{Name: "Birmingham", PostalCode: "35203", AreaCodes: []string{"205"}, Weight: 60}, {Name: "Montgomery", PostalCode: "36104", AreaCodes: []string{"334"}, Weight: 40}}},
		{Code: "AK", Name: "Alaska", Weight: 1, Cities: []City{{Name: "Anchorage", PostalCode: "99501", AreaCodes: []string{"907"}, Weight: 100}}},
		{Code: "AZ", Name: "Arizona", Weight: 7, Cities: []City{{Name: "Phoenix", PostalCode: "85004", AreaCodes: []string{"602"}, Weight: 70}, {Name: "Tucson", PostalCode: "85701", AreaCodes: []string{"520"}, Weight: 30}}},
		{Code: "AR", Name: "Arkansas", Weight: 3, Cities: []City{{Name: "Little Rock", PostalCode: "72201", AreaCodes: []string{"501"}, Weight: 100}}},
		{Code: "CA", Name: "California", Weight: 39, Cities: []City{{Name: "Los Angeles", PostalCode: "90012", AreaCodes: []string{"213"}, Weight: 40}, {Name: "San Diego", PostalCode: "92101", AreaCodes: []string{"619"}, Weight: 20}, {Name: "San Francisco", PostalCode: "94103", AreaCodes: []string{"415"}, Weight: 20}, {Name: "San Jose", PostalCode: "95113", AreaCodes: []string{"408"}, Weight: 20}}},
		{Code: "CO", Name: "Colorado", Weight: 6, Cities: []City{{Name: "Denver", PostalCode: "80202", AreaCodes: []string{"303"}, Weight: 70}, {Name: "Colorado Springs", PostalCode: "80903", AreaCodes: []string{"719"}, Weight: 30}}},
		{Code: "CT", Name: "Connecticut", Weight: 4, Cities: []City{{Name: "Bridgeport", PostalCode: "06604", AreaCodes: []string{"203"}, Weight: 55}, {Name: "Hartford", PostalCode: "06103", AreaCodes: []string{"860"}, Weight: 45}}},
		{Code: "DE", Name: "Delaware", Weight: 1, Cities: []City{{Name: "Wilmington", PostalCode: "19801", AreaCodes: []string{"302"}, Weight: 100}}},
		{Code: "DC", Name: "District of Columbia", Weight: 1, Cities: []City{{Name: "Washington", PostalCode: "20001", AreaCodes: []string{"202"}, Weight: 100}}},
		{Code: "FL", Name: "Florida", Weight: 22, Cities: []City{{Name: "Miami", PostalCode: "33130", AreaCodes: []string{"305"}, Weight: 35}, {Name: "Orlando", PostalCode: "32801", AreaCodes: []string{"407"}, Weight: 20}, {Name: "Tampa", PostalCode: "33602", AreaCodes: []string{"813"}, Weight: 20}, {Name: "Jacksonville", PostalCode: "32202", AreaCodes: []string{"904"}, Weight: 25}}},
		{Code: "GA", Name: "Georgia", Weight: 11, Cities: []City{{Name: "Atlanta", PostalCode: "30303", AreaCodes: []string{"404"}, Weight: 75}, {Name: "Savannah", PostalCode: "31401", AreaCodes: []string{"912"}, Weight: 25}}},
		{Code: "HI", Name: "Hawaii", Weight: 1, Cities: []City{{Name: "Honolulu", PostalCode: "96813", AreaCodes: []string{"808"}, Weight: 100}}},
		{Code: "ID", Name: "Idaho", Weight: 2, Cities: []City{{Name: "Boise", PostalCode: "83702", AreaCodes: []string{"208"}, Weight: 100}}},
		{Code: "IL", Name: "Illinois", Weight: 12, Cities: []City{{Name: "Chicago", PostalCode: "60601", AreaCodes: []string{"312"}, Weight: 80}, {Name: "Springfield", PostalCode: "62701", AreaCodes: []string{"217"}, Weight: 20}}},
		{Code: "IN", Name: "Indiana", Weight: 7, Cities: []City{{Name: "Indianapolis", PostalCode: "46204", AreaCodes: []string{"317"}, Weight: 100}}},
		{Code: "IA", Name: "Iowa", Weight: 3, Cities: []City{{Name: "Des Moines", PostalCode: "50309", AreaCodes: []string{"515"}, Weight: 100}}},
		{Code: "KS", Name: "Kansas", Weight: 3, Cities: []City{{Name: "Wichita", PostalCode: "67202", AreaCodes: []string{"316"}, Weight: 100}}},
		{Code: "KY", Name: "Kentucky", Weight: 4, Cities: []City{{Name: "Louisville", PostalCode: "40202", AreaCodes: []string{"502"}, Weight: 65}, {Name: "Lexington", PostalCode: "40507", AreaCodes: []string{"859"}, Weight: 35}}},
		{Code: "LA", Name: "Louisiana", Weight: 5, Cities: []City{{Name: "New Orleans", PostalCode: "70112", AreaCodes: []string{"504"}, Weight: 60}, {Name: "Baton Rouge", PostalCode: "70801", AreaCodes: []string{"225"}, Weight: 40}}},
		{Code: "ME", Name: "Maine", Weight: 1, Cities: []City{{Name: "Portland", PostalCode: "04101", AreaCodes: []string{"207"}, Weight: 100}}},
		{Code: "MD", Name: "Maryland", Weight: 6, Cities: []City{{Name: "Baltimore", PostalCode: "21201", AreaCodes: []string{"410"}, Weight: 70}, {Name: "Annapolis", PostalCode: "21401", AreaCodes: []string{"410"}, Weight: 30}}},
		{Code: "MA", Name: "Massachusetts", Weight: 7, Cities: []City{{Name: "Boston", PostalCode: "02108", AreaCodes: []string{"617"}, Weight: 75}, {Name: "Worcester", PostalCode: "01608", AreaCodes: []string{"508"}, Weight: 25}}},
		{Code: "MI", Name: "Michigan", Weight: 10, Cities: []City{{Name: "Detroit", PostalCode: "48226", AreaCodes: []string{"313"}, Weight: 65}, {Name: "Grand Rapids", PostalCode: "49503", AreaCodes: []string{"616"}, Weight: 35}}},
		{Code: "MN", Name: "Minnesota", Weight: 6, Cities: []City{{Name: "Minneapolis", PostalCode: "55401", AreaCodes: []string{"612"}, Weight: 70}, {Name: "Saint Paul", PostalCode: "55102", AreaCodes: []string{"651"}, Weight: 30}}},
		{Code: "MS", Name: "Mississippi", Weight: 3, Cities: []City{{Name: "Jackson", PostalCode: "39201", AreaCodes: []string{"601"}, Weight: 100}}},
		{Code: "MO", Name: "Missouri", Weight: 6, Cities: []City{{Name: "Kansas City", PostalCode: "64106", AreaCodes: []string{"816"}, Weight: 55}, {Name: "St. Louis", PostalCode: "63101", AreaCodes: []string{"314"}, Weight: 45}}},
		{Code: "MT", Name: "Montana", Weight: 1, Cities: []City{{Name: "Billings", PostalCode: "59101", AreaCodes: []string{"406"}, Weight: 100}}},
		{Code: "NE", Name: "Nebraska", Weight: 2, Cities: []City{{Name: "Omaha", PostalCode: "68102", AreaCodes: []string{"402"}, Weight: 100}}},
		{Code: "NV", Name: "Nevada", Weight: 3, Cities: []City{{Name: "Las Vegas", PostalCode: "89101", AreaCodes: []string{"702"}, Weight: 75}, {Name: "Reno", PostalCode: "89501", AreaCodes: []string{"775"}, Weight: 25}}},
		{Code: "NH", Name: "New Hampshire", Weight: 1, Cities: []City{{Name: "Manchester", PostalCode: "03101", AreaCodes: []string{"603"}, Weight: 100}}},
		{Code: "NJ", Name: "New Jersey", Weight: 9, Cities: []City{{Name: "Newark", PostalCode: "07102", AreaCodes: []string{"973"}, Weight: 55}, {Name: "Jersey City", PostalCode: "07302", AreaCodes: []string{"201"}, Weight: 45}}},
		{Code: "NM", Name: "New Mexico", Weight: 2, Cities: []City{{Name: "Albuquerque", PostalCode: "87102", AreaCodes: []string{"505"}, Weight: 100}}},
		{Code: "NY", Name: "New York", Weight: 20, Cities: []City{{Name: "New York", PostalCode: "10001", AreaCodes: []string{"212"}, Weight: 65}, {Name: "Buffalo", PostalCode: "14202", AreaCodes: []string{"716"}, Weight: 20}, {Name: "Rochester", PostalCode: "14604", AreaCodes: []string{"585"}, Weight: 15}}},
		{Code: "NC", Name: "North Carolina", Weight: 11, Cities: []City{{Name: "Charlotte", PostalCode: "28202", AreaCodes: []string{"704"}, Weight: 60}, {Name: "Raleigh", PostalCode: "27601", AreaCodes: []string{"919"}, Weight: 40}}},
		{Code: "ND", Name: "North Dakota", Weight: 1, Cities: []City{{Name: "Fargo", PostalCode: "58102", AreaCodes: []string{"701"}, Weight: 100}}},
		{Code: "OH", Name: "Ohio", Weight: 12, Cities: []City{{Name: "Columbus", PostalCode: "43215", AreaCodes: []string{"614"}, Weight: 40}, {Name: "Cleveland", PostalCode: "44114", AreaCodes: []string{"216"}, Weight: 35}, {Name: "Cincinnati", PostalCode: "45202", AreaCodes: []string{"513"}, Weight: 25}}},
		{Code: "OK", Name: "Oklahoma", Weight: 4, Cities: []City{{Name: "Oklahoma City", PostalCode: "73102", AreaCodes: []string{"405"}, Weight: 100}}},
		{Code: "OR", Name: "Oregon", Weight: 4, Cities: []City{{Name: "Portland", PostalCode: "97204", AreaCodes: []string{"503"}, Weight: 100}}},
		{Code: "PA", Name: "Pennsylvania", Weight: 13, Cities: []City{{Name: "Philadelphia", PostalCode: "19102", AreaCodes: []string{"215"}, Weight: 60}, {Name: "Pittsburgh", PostalCode: "15222", AreaCodes: []string{"412"}, Weight: 40}}},
		{Code: "RI", Name: "Rhode Island", Weight: 1, Cities: []City{{Name: "Providence", PostalCode: "02903", AreaCodes: []string{"401"}, Weight: 100}}},
		{Code: "SC", Name: "South Carolina", Weight: 5, Cities: []City{{Name: "Charleston", PostalCode: "29401", AreaCodes: []string{"843"}, Weight: 45}, {Name: "Columbia", PostalCode: "29201", AreaCodes: []string{"803"}, Weight: 55}}},
		{Code: "SD", Name: "South Dakota", Weight: 1, Cities: []City{{Name: "Sioux Falls", PostalCode: "57104", AreaCodes: []string{"605"}, Weight: 100}}},
		{Code: "TN", Name: "Tennessee", Weight: 7, Cities: []City{{Name: "Nashville", PostalCode: "37203", AreaCodes: []string{"615"}, Weight: 55}, {Name: "Memphis", PostalCode: "38103", AreaCodes: []string{"901"}, Weight: 45}}},
		{Code: "TX", Name: "Texas", Weight: 30, Cities: []City{{Name: "Houston", PostalCode: "77002", AreaCodes: []string{"713"}, Weight: 35}, {Name: "Dallas", PostalCode: "75201", AreaCodes: []string{"214"}, Weight: 25}, {Name: "Austin", PostalCode: "78701", AreaCodes: []string{"512"}, Weight: 20}, {Name: "San Antonio", PostalCode: "78205", AreaCodes: []string{"210"}, Weight: 20}}},
		{Code: "UT", Name: "Utah", Weight: 3, Cities: []City{{Name: "Salt Lake City", PostalCode: "84101", AreaCodes: []string{"801"}, Weight: 100}}},
		{Code: "VT", Name: "Vermont", Weight: 1, Cities: []City{{Name: "Burlington", PostalCode: "05401", AreaCodes: []string{"802"}, Weight: 100}}},
		{Code: "VA", Name: "Virginia", Weight: 9, Cities: []City{{Name: "Virginia Beach", PostalCode: "23451", AreaCodes: []string{"757"}, Weight: 55}, {Name: "Richmond", PostalCode: "23219", AreaCodes: []string{"804"}, Weight: 45}}},
		{Code: "WA", Name: "Washington", Weight: 8, Cities: []City{{Name: "Seattle", PostalCode: "98101", AreaCodes: []string{"206"}, Weight: 70}, {Name: "Spokane", PostalCode: "99201", AreaCodes: []string{"509"}, Weight: 30}}},
		{Code: "WV", Name: "West Virginia", Weight: 2, Cities: []City{{Name: "Charleston", PostalCode: "25301", AreaCodes: []string{"304"}, Weight: 100}}},
		{Code: "WI", Name: "Wisconsin", Weight: 6, Cities: []City{{Name: "Milwaukee", PostalCode: "53202", AreaCodes: []string{"414"}, Weight: 60}, {Name: "Madison", PostalCode: "53703", AreaCodes: []string{"608"}, Weight: 40}}},
		{Code: "WY", Name: "Wyoming", Weight: 1, Cities: []City{{Name: "Cheyenne", PostalCode: "82001", AreaCodes: []string{"307"}, Weight: 100}}},
	},
	MaleFirstNames: []WeightedValue{
		{Value: "James", Weight: 10}, {Value: "Michael", Weight: 9}, {Value: "Robert", Weight: 8}, {Value: "John", Weight: 8}, {Value: "David", Weight: 7},
		{Value: "William", Weight: 7}, {Value: "Joseph", Weight: 6}, {Value: "Daniel", Weight: 6}, {Value: "Matthew", Weight: 5}, {Value: "Christopher", Weight: 5},
		{Value: "Andrew", Weight: 4}, {Value: "Joshua", Weight: 4}, {Value: "Anthony", Weight: 4}, {Value: "Ryan", Weight: 4}, {Value: "Nicholas", Weight: 3},
		{Value: "Ethan", Weight: 3}, {Value: "Noah", Weight: 3}, {Value: "Benjamin", Weight: 3}, {Value: "Alexander", Weight: 3}, {Value: "Samuel", Weight: 3},
		{Value: "Elijah", Weight: 2}, {Value: "Liam", Weight: 2}, {Value: "Mason", Weight: 2}, {Value: "Jacob", Weight: 2}, {Value: "Logan", Weight: 2},
	},
	FemaleFirstName: []WeightedValue{
		{Value: "Mary", Weight: 10}, {Value: "Patricia", Weight: 8}, {Value: "Jennifer", Weight: 7}, {Value: "Linda", Weight: 7}, {Value: "Elizabeth", Weight: 6},
		{Value: "Barbara", Weight: 6}, {Value: "Susan", Weight: 5}, {Value: "Jessica", Weight: 5}, {Value: "Sarah", Weight: 5}, {Value: "Karen", Weight: 4},
		{Value: "Nancy", Weight: 4}, {Value: "Lisa", Weight: 4}, {Value: "Emily", Weight: 4}, {Value: "Ashley", Weight: 3}, {Value: "Emma", Weight: 3},
		{Value: "Olivia", Weight: 3}, {Value: "Sophia", Weight: 3}, {Value: "Ava", Weight: 3}, {Value: "Isabella", Weight: 3}, {Value: "Mia", Weight: 2},
		{Value: "Charlotte", Weight: 2}, {Value: "Amelia", Weight: 2}, {Value: "Evelyn", Weight: 2}, {Value: "Abigail", Weight: 2}, {Value: "Harper", Weight: 2},
	},
	LastNames: []WeightedValue{
		{Value: "Smith", Weight: 12}, {Value: "Johnson", Weight: 10}, {Value: "Williams", Weight: 9}, {Value: "Brown", Weight: 8}, {Value: "Jones", Weight: 8},
		{Value: "Garcia", Weight: 7}, {Value: "Miller", Weight: 7}, {Value: "Davis", Weight: 6}, {Value: "Rodriguez", Weight: 6}, {Value: "Martinez", Weight: 6},
		{Value: "Hernandez", Weight: 5}, {Value: "Lopez", Weight: 5}, {Value: "Gonzalez", Weight: 5}, {Value: "Wilson", Weight: 5}, {Value: "Anderson", Weight: 4},
		{Value: "Thomas", Weight: 4}, {Value: "Taylor", Weight: 4}, {Value: "Moore", Weight: 4}, {Value: "Jackson", Weight: 4}, {Value: "Martin", Weight: 4},
		{Value: "Lee", Weight: 3}, {Value: "Perez", Weight: 3}, {Value: "Thompson", Weight: 3}, {Value: "White", Weight: 3}, {Value: "Harris", Weight: 3},
		{Value: "Sanchez", Weight: 3}, {Value: "Clark", Weight: 3}, {Value: "Ramirez", Weight: 3}, {Value: "Lewis", Weight: 2}, {Value: "Robinson", Weight: 2},
		{Value: "Walker", Weight: 2}, {Value: "Young", Weight: 2}, {Value: "Allen", Weight: 2}, {Value: "King", Weight: 2}, {Value: "Wright", Weight: 2},
		{Value: "Scott", Weight: 2}, {Value: "Torres", Weight: 2}, {Value: "Nguyen", Weight: 2}, {Value: "Hill", Weight: 2}, {Value: "Flores", Weight: 2},
	},
	StreetNames: []WeightedValue{
		{Value: "Oak", Weight: 8}, {Value: "Maple", Weight: 7}, {Value: "Pine", Weight: 7}, {Value: "Cedar", Weight: 6}, {Value: "Elm", Weight: 6},
		{Value: "Washington", Weight: 6}, {Value: "Lake", Weight: 5}, {Value: "Hill", Weight: 5}, {Value: "Park", Weight: 5}, {Value: "Walnut", Weight: 4},
		{Value: "Sunset", Weight: 4}, {Value: "Ridge", Weight: 4}, {Value: "Main", Weight: 7}, {Value: "Highland", Weight: 4}, {Value: "Meadow", Weight: 4},
		{Value: "Lincoln", Weight: 4}, {Value: "Jefferson", Weight: 3}, {Value: "Cherry", Weight: 3}, {Value: "Forest", Weight: 3}, {Value: "Willow", Weight: 3},
	},
	StreetSuffixes: []WeightedValue{
		{Value: "Street", Weight: 10}, {Value: "Avenue", Weight: 8}, {Value: "Drive", Weight: 6}, {Value: "Lane", Weight: 5},
		{Value: "Road", Weight: 7}, {Value: "Court", Weight: 4}, {Value: "Boulevard", Weight: 4}, {Value: "Place", Weight: 3},
	},
}

func init() {
	usData.StatesByCode = make(map[string]State, len(usData.States))
	for _, state := range usData.States {
		usData.StatesByCode[state.Code] = state
	}
}

func US() *Data {
	return &usData
}
