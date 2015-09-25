package main

import (
	"dsp/bidding"
	"dsp/campaign"
	"dsp/user"
	"fmt"
	"net/http"
	// _ "net/http/pprof"
	"encoding/json"
	"runtime"
	"strconv"
)

const (
	MAX_NUMBER_OF_CAMPAIGNS  = 10000
	MAX_NUMBER_OF_TARGETS    = 26
	MAX_NUMBER_OF_ATTRIBUTES = 100
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	http.HandleFunc("/campaign", GetCampaigns)
	http.HandleFunc("/user", GetUser)
	http.HandleFunc("/import_camp", ImportCampaigns)
	http.HandleFunc("/search", Search)
	http.HandleFunc("/search_auto", SearchAuto)
	http.HandleFunc("/dump", DumpData)

	http.ListenAndServe(":3000", nil)
}

// Handler generates JSON array of campaigns
// Has 3 not required arguments:
//   x - max number of attributes per target in campaign target list (max: 100)
//   y - max number of targets in campaign target list (max: 26)
//   z - number of campaigns to generate (max: 10000)
// Responses:
//   200 OK - return JSON array of campaigns
//   400 Bad Request
func GetCampaigns(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		sendString(w, http.StatusNotFound, "Not Found")
		return
	}

	maxAttributesNumber, err := getReqIntArg(req, "x", MAX_NUMBER_OF_ATTRIBUTES)
	if err != nil {
		sendString(w, http.StatusBadRequest, "Argument 'x' must be a number")
		return
	}

	maxTargetsNumber, err := getReqIntArg(req, "y", MAX_NUMBER_OF_TARGETS)
	if err != nil {
		sendString(w, http.StatusBadRequest, "Argument 'y' must be a number")
		return
	}

	campaignsNumber, err := getReqIntArg(req, "z", MAX_NUMBER_OF_CAMPAIGNS)
	if err != nil {
		sendString(w, http.StatusBadRequest, "Argument 'z' must be a number")
		return
	}

	if err := validateCampaignArguments(campaignsNumber, maxTargetsNumber, maxAttributesNumber); err != nil {
		sendString(w, http.StatusBadRequest, err.Error())
		return
	}

	sendJson(w, http.StatusOK, campaign.GenerateMany(campaignsNumber, maxTargetsNumber, maxAttributesNumber))
}

// Generates randomly filled User JSON object
// Responses:
//   200 OK - return JSON object of user
func GetUser(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		sendString(w, http.StatusNotFound, "Not Found")
		return
	}
	sendJson(w, http.StatusOK, user.Generate())
}

// Imports campaigns and updates bidding.data struct
// Json array of campaigns should be send in request body (POST)
// Responses:
//    200 OK - return nothing
//    400 Bad Request - return error string
func ImportCampaigns(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		sendString(w, http.StatusNotFound, "Not Found")
		return
	}

	campaigns := make([]campaign.Campaign, 0)
	decoder := json.NewDecoder(req.Body)
	if decoder.Decode(&campaigns) != nil {
		sendString(w, http.StatusBadRequest, "Bad Request")
		return
	}

	if err := bidding.ReloadData(campaigns); err != nil {
		sendString(w, http.StatusBadRequest, err.Error())
		return
	}

	sendString(w, http.StatusOK, "")
}

// Returns won campaign by user profile and campaign max price
// Json array of campaigns should be send in request body (POST)
// Responses:
//    200 OK - return winner struct as Json object
//    400 Bad Request - return error string
func Search(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		sendString(w, http.StatusNotFound, "Not Found")
		return
	}

	u := &user.User{}
	decoder := json.NewDecoder(req.Body)
	if decoder.Decode(u) != nil {
		sendString(w, http.StatusBadRequest, "Bad Request")
		return
	}

	sendJson(w, http.StatusOK, bidding.ProcessBid(u))
}

// Generates user and tries to find campaign by user profile and max price
// Responses:
//    200 OK - return winner struct as Json object
//    400 Bad Request - return error string
func SearchAuto(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		sendString(w, http.StatusNotFound, "Not Found")
		return
	}
	sendJson(w, http.StatusOK, bidding.ProcessBid(user.Generate()))
}

// Returns dump of bidding.data object as JSON
// Responses:
//    200 OK
func DumpData(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		sendString(w, http.StatusNotFound, "Not Found")
		return
	}
	sendJson(w, http.StatusOK, bidding.GetData())
}

// Validates arguments of GetCampaigns handler
func validateCampaignArguments(campaignsNumber, maxTargetsNumber, maxAttributesNumber int) error {
	if maxAttributesNumber > MAX_NUMBER_OF_ATTRIBUTES {
		return fmt.Errorf("Argument 'x' cannot be more than %d", MAX_NUMBER_OF_ATTRIBUTES)
	}

	if maxTargetsNumber > MAX_NUMBER_OF_TARGETS {
		return fmt.Errorf("Argument 'y' cannot be more than %d", MAX_NUMBER_OF_TARGETS)
	}

	if campaignsNumber > MAX_NUMBER_OF_CAMPAIGNS {
		return fmt.Errorf("Argument 'z' cannot be more than %d", MAX_NUMBER_OF_CAMPAIGNS)
	}

	return nil
}

// sends json to client
func sendJson(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(data); err != nil {
		panic(err)
	}
}

// sends string to client
func sendString(w http.ResponseWriter, code int, str string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(code)
	fmt.Fprint(w, str)
}

// returns int argument from request query
func getReqIntArg(req *http.Request, name string, defaultVal int) (int, error) {
	val := req.URL.Query().Get(name)
	if val == "" {
		return defaultVal, nil
	} else {
		return strconv.Atoi(val)
	}
}
