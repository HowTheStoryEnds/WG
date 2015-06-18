package WG

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"strconv"
	//"github.com/davecgh/go-spew/spew"
)

const MasteryNone uint32 = 0
const MasteryThirdClass uint32 = 1
const MasterySecondClass uint32 = 2
const MasteryFirstClass uint32 = 3
const MasteryAceTanker uint32 = 4

type MetaCount struct {
	Count uint32 `json:"count"`
}
type ApiRequest struct {
	Status string    `json:"status"`
	Meta   MetaCount `json:"meta"`
}

type ApiResponse interface {
	PlayerList() []Player
	VehicleList() []Vehicle
}

type Vehicle struct {
	// https://eu.wargaming.net/developers/api_reference/wot/account/tanks/
	MarkOfMastery uint32            `json:"mark_of_mastery"`
	TankId        uint32            `json:"tank_id"`
	Statistics    VehicleStatistics `json:"statistics"`

	// custom field
	Owner uint32 `json:",omitempty"`
}
type VehicleStatistics struct {
	Battles uint32 `json:"battles"`
	Wins    uint32 `json:"wins"`
}
type Player struct {
	//https://eu.wargaming.net/developers/api_reference/wot/account/info/
	AccountId      uint32        `json:"account_id"`
	ClanId         uint32        `json:"clan_id"`
	ClientLanguage string        `json:"client_language"`
	CreatedAt      uint32        `json:"created_at"`
	GlobalRating   uint32        `json:"global_rating"`
	LastBattleTime uint32        `json:"last_battle_time"`
	LogoutAt       uint32        `json:"logout_at"`
	Nickname       string        `json:"nickname"`
	UpdatedAt      uint32        `json:"updated_at"`
	Private        PlayerPrivate `json:"private"`
	Statistics     PlayerStatistics

	// holders not part of the regular api
	Region   string
	Vehicles []Vehicle
}

type PlayerStatistics struct {
	//https://eu.wargaming.net/developers/api_reference/wot/account/info/
	TreesCut uint32 `json:"trees_cut"`

	All                SubStatistics `json:"all"`
	Company            SubStatistics `json:"company"`
	Clan               SubStatistics `json:"clan"`
	Historical         SubStatistics `json:"historical"`
	RegularTeam        SubStatistics `json:"regular_team"`
	StrongholdDefense  SubStatistics `json:"stronghold_defense"`
	StrongholdSkirmish SubStatistics `json:"stronghold_skirmish"`
	Team               SubStatistics `json:"team"`
}
type SubStatistics struct {
	//https://eu.wargaming.net/developers/api_reference/wot/account/info/
	AvgDamageAssisted          float64 `json:"avg_damage_assisted"`
	AvgDamageAssistedRadio     float64 `json:"avg_damage_assisted_radio"`
	AvgDamageAssistedTrack     float64 `json:"avg_damage_assisted_track"`
	AvgDamageBlocked           float64 `json:"avg_damage_blocked"`
	BattleAvgXp                uint32  `json:"battle_avg_xp"`
	Battles                    uint32  `json:"battles"`
	CapturePoints              uint32  `json:"capture_points"`
	DamageDealt                uint32  `json:"damage_dealt"`
	DamageReceived             uint32  `json:"damage_received"`
	DirectHitsReceived         uint32  `json:"direct_hits_received"`
	Draws                      uint32  `json:"draws"`
	DroppedCapturePoints       uint32  `json:"dropped_capture_points"`
	ExplosionHits              uint32  `json:"explosion_hits"`
	ExplosionHitsReceived      uint32  `json:"explosion_hits_received"`
	Frags                      uint32  `json:"frags"`
	Hits                       uint32  `json:"hits"`
	HitsPercents               uint32  `json:"hits_percents"`
	Losses                     uint32  `json:"losses"`
	MaxDamage                  uint32  `json:"max_damage"`
	MaxDamageTankId            uint32  `json:"max_damage_tank_id"`
	MaxFrags                   uint32  `json:"max_frags"`
	MaxFragsTankId             uint32  `json:"max_frags_tank_id"`
	MaxXp                      uint32  `json:"max_xp"`
	MaxXpTankId                uint32  `json:"max_xp_tank_id"`
	NoDamageDirectHitsReceived uint32  `json:"no_damage_direct_hits_received"`
	Piercings                  uint32  `json:"piercings"`
	PiercingsReceived          uint32  `json:"piercings_received"`
	Shots                      uint32  `json:"shots"`
	Spotted                    uint32  `json:"spotted"`
	SurvivedBattles            uint32  `json:"survived_battles"`
	TankingFactor              float64 `json:"tanking_factor"`
	Wins                       uint32  `json:"wins"`
	Xp                         uint32  `json:"xp"`
}

type PrivateGroupedContacts struct {
	//https://eu.wargaming.net/developers/api_reference/wot/account/info/
	Blocked   []uint32          `json:"blocked"`
	Groups    map[string]string `json:"groups"`
	Ignored   []uint32          `json:"ignored"`
	Muted     []uint32          `json:"muted"`
	Ungrouped []uint32          `json:"ungrouped"`
}
type PrivateRented struct {
	//https://eu.wargaming.net/developers/api_reference/wot/account/info/
	CompensationCredits uint32 `json:"compensation_credits"`
	CompensationGold    uint32 `json:"compensation_gold"`
	ExpirationTime      uint32 `json:"expiration_time"`
	TankId              uint32 `json:"tank_id"`
}
type PrivateRestrictions struct {
	//https://eu.wargaming.net/developers/api_reference/wot/account/info/
	ChatBanTime uint32 `json:"chat_ban_time"`
	ClanTime    uint32 `json:"clan_time"`
}
type PlayerPrivate struct {
	//https://eu.wargaming.net/developers/api_reference/wot/account/info/
	BanInfo          string                 `json:"ban_info"`
	BanTime          uint32                 `json:"ban_time"`
	BattleLifeTime   uint32                 `json:"battle_life_time"`
	Credits          uint32                 `json:"credits"`
	FreeXp           uint32                 `json:"free_xp"`
	Friends          []uint32               `json:"friends"`
	Gold             uint32                 `json:"gold"`
	IsBoundToPhone   bool                   `json:"is_bound_to_phone"`
	IsPremium        bool                   `json:"is_premium"`
	PersonalMissions map[string]string      `json:"personal_missions"`
	PremiumExpiresAt uint32                 `json:"premium_expires_at"`
	GroupedContacts  PrivateGroupedContacts `json:"grouped_contacts"`
	Rented           PrivateRented          `json:"rented"`
	Restrictions     PrivateRestrictions    `json:"restrictions"`
}

type RequestAccountList struct {
	ApiRequest
	Data []Player `json:"data"`
}

func (r RequestAccountList) PlayerList() []Player {
	return r.Data
}
func (r RequestAccountList) VehicleList() []Vehicle {
	panic("can return players only")
	return nil
}

type RequestAccountInfo struct {
	ApiRequest
	Data map[string]Player `json:"data"`
}

func (r RequestAccountInfo) PlayerList() []Player {
	players := []Player{}
	for _, v := range r.Data {
		players = append(players, v)
	}
	return players
}
func (r RequestAccountInfo) VehicleList() []Vehicle {
	panic("can return players only")
	return nil
}

type RequestAccountTanks struct {
	ApiRequest
	Data map[string][]Vehicle `json:"data"`
}

func (r RequestAccountTanks) PlayerList() []Player {
	panic("can return players only")
	return nil

}
func (r RequestAccountTanks) VehicleList() []Vehicle {
	vehicles := []Vehicle{}
	for k, v := range r.Data {
		for _, vv := range v {
			Owner, _ := strconv.ParseUint(k, 10, 32)
			vv.Owner = uint32(Owner)
			vehicles = append(vehicles, vv)
		}
	}
	return vehicles
}

// we request our data from WG server, we allow for net/http
type WG struct {
	region    string
	transport string
	method    string
	apiKey    string
}

func (w *WG) Init(region string, transport string, method string, ApiKey string) {
	w.SetRegion(region)
	w.SetTransport(transport)
	w.SetMethod(method)
	w.SetApiKey(ApiKey)
}

func (w *WG) constructURL(action string, parameters map[string]string) string {
	//sane defaults
	if w.transport != "http" {
		w.transport = "https"
	}

	var regional = make(map[string]string)
	regional["eu"] = "api.worldoftanks.eu/wot/"
	regional["na"] = "api.worldoftanks.com/wot/"
	regional["ru"] = "api.worldoftanks.ru/wot/"

	uri := w.transport + "://" + regional[w.region]

	baseUrl, err := url.Parse(uri)
	if err != nil {
		panic("constructURL: " + err.Error())

	}

	params := url.Values{}
	for key, value := range parameters {
		params.Add(key, value)

	}

	baseUrl.RawQuery = params.Encode()
	return uri + action + "/" + "?" + baseUrl.RawQuery

}

func (w *WG) retrieveData(command string, params map[string]string) ([]byte, error) {

	// make sure the apikey is set
	if _, present := params["application_id"]; !present {
		params["application_id"] = w.apiKey
	}

	uri := w.constructURL(command, params)

	// and retrieve the JSON data in a string
	var val, err = http.Get(uri)
	var res = []byte{}
	if err == nil {
		res, err = ioutil.ReadAll(val.Body)
		val.Body.Close()
	}

	return res, err
}

func (w *WG) SetRegion(val string) {
	w.region = val
}
func (w *WG) SetTransport(val string) {
	w.transport = val
}
func (w *WG) SetMethod(val string) {
	w.method = val
}
func (w *WG) SetApiKey(val string) {
	w.apiKey = val
}

func (w *WG) apiCall(command string, params map[string]string) (ApiResponse, error) {
	data, err := w.retrieveData(command, params)
	if err != nil {
		err = errors.New("apiCall <" + command + ">: could not retrieve data\n" + err.Error())
		return nil, err

	}

	var req ApiResponse
	switch command {
	case "account/list":
		reqt := RequestAccountList{}
		err = json.Unmarshal(data, &reqt)
		req = reqt

	case "account/info":
		reqt := RequestAccountInfo{}
		err = json.Unmarshal(data, &reqt)
		req = reqt

	case "account/tanks":
		reqt := RequestAccountTanks{}
		err = json.Unmarshal(data, &reqt)
		req = reqt
	default:
		panic("Unknown command, this shouldn't happen")
	}

	// no errors so we return the data retrieved
	if err == nil {
		return req, err
	}
	//otherwise return just the error
	return nil, err

}

// account/list
func (w *WG) SearchPlayersByName(name string, exact bool) ([]Player, error) {
	params := make(map[string]string)
	if exact {
		params["type"] = "exact"
	}
	params["search"] = name
	result, err := w.apiCall("account/list", params)
	if err != nil {
		fmt.Println("SearchPlayersByName: " + err.Error())
		return nil, err
	}
	return result.PlayerList(), err

}

// account/info
func (w *WG) GetPlayerPersonalData(accountid []uint32) ([]Player, error) {
	params := make(map[string]string)

	for _, v := range accountid {
		_, ok := params["account_id"]
		if len(params["account_id"]) == 0 || !ok {
			params["account_id"] = fmt.Sprint(v)
		} else {
			params["account_id"] += "," + fmt.Sprint(v)
		}
	}

	result, err := w.apiCall("account/info", params)
	if err != nil {
		fmt.Println("GetPlayerPersonalData: " + err.Error())
		return nil, err
	}
	return result.PlayerList(), err

}

// account/tanks
func (w *WG) GetPlayerTanks(accountid []uint32, vehicleid []uint32) ([]Vehicle, error) {
	params := make(map[string]string)
	for _, v := range vehicleid {
		_, ok := params["tank_id"]
		if len(params["tank_id"]) == 0 || !ok {
			params["tank_id"] = fmt.Sprint(v)
		} else {
			params["tank_id"] += "," + fmt.Sprint(v)
		}
	}
	for _, v := range accountid {
		_, ok := params["account_id"]
		if len(params["account_id"]) == 0 || !ok {
			params["account_id"] = fmt.Sprint(v)
		} else {
			params["account_id"] += "," + fmt.Sprint(v)
		}
	}
	result, err := w.apiCall("account/tanks", params)
	if err != nil {
		fmt.Println("GetPlayerTanks: " + err.Error())
		return nil, err
	}
	return result.VehicleList(), err

}
