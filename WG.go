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

const OrderById = "id"
const OrderByIdDesc = "-id"
const OrderByName = "name"
const OrderByNameDesc = "-name"
const OrderBySize = "members_count"
const OrderBySizeDesc = "-members_count"
const OrderByTag = "tag"
const OrderByTagDesc = "-tag"
const OrderByCreationDate = "created_at"
const OrderByCreationDateDesc = "-created_at"

const NoAccessToken = ""

type ApiRequestError struct {
	ApiRequest
	WGError ApiError `json:"error"`
	Sibling error
}
type ApiError struct {
	Code    uint32 `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field"`
	Value   string `json:"value"`
}

func (e ApiRequestError) Error() string {
	return "Code: " + fmt.Sprint(e.WGError.Code) + "\n" + e.WGError.Message + "------------\n" + e.Sibling.Error()
}

type ApiRequest struct {
	Status string    `json:"status"`
	Meta   MetaCount `json:"meta"`
	Total  uint32    `json:"total"`
}
type MetaCount struct {
	Count uint32 `json:"count"`
}
type ApiResponse interface {
	PlayerList() []Player
	VehicleList() []Vehicle
	ClanList() []Clan
}
type NoClanList struct{}

func (n NoClanList) ClanList() []Clan {
	panic("can not return clans, incorrect call")
	return nil
}

type NoVehicleList struct{}

func (n NoVehicleList) VehicleList() []Vehicle {
	panic("can not return vehicles, incorrect call")
	return nil
}

type NoPlayerList struct{}

func (n NoPlayerList) PlayerList() []Player {
	panic("can not return players, incorrect call")
	return nil
}

type Clan struct {
	AcceptsJoinRequests bool         `json:"accepts_join_requests"`
	ClanId              uint32       `json:"clan_id"`
	Color               string       `json:"color"`
	CreatedAt           uint32       `json:"created_at"`
	CreatorId           uint32       `json:"creator_id"`
	CreatorName         string       `json:"creator_name"`
	Description         string       `json:"description"`
	DescriptionHtml     string       `json:"description_html"`
	IsClanDisbanded     bool         `json:"is_clan_disbanded"`
	LeaderId            uint32       `json:"leader_id"`
	LeaderName          string       `json:"leader_name"`
	MembersCount        uint32       `json:"members_count"`
	Motto               string       `json:"motto"`
	Name                string       `json:"name"`
	OldTag              string       `json:"old_tag"`
	OldName             string       `json:"old_name"`
	RenamedAt           uint32       `json:"renamed_at"`
	Tag                 string       `json:"tag"`
	UpdatedAt           uint32       `json:"updated_at"`
	Emblems             ClanEmblems  `json:"emblems"`
	Members             []ClanMember `json:"members"`
	Private             ClanPrivate  `json:"private"`
}
type ClanEmblems struct {
	X195 map[string]string `json:"x195"`
	X24  map[string]string `json:"x24"`
	X256 map[string]string `json:"x256"`
	X32  map[string]string `json:"x32"`
	X64  map[string]string `json:"x64"`
}
type ClanPrivate struct {
	Treasury uint32 `json:"treasury"`
}
type ClanMember struct {
	AccountId   uint32 `json:"account_id"`
	AccountName string `json:"account_name"`
	JoinedAt    uint32 `json:"joined_at"`
	Role        string `json:"role"`
	RoleI18n    string `json:"role_i18n"`
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
	NoVehicleList
	NoClanList
	ApiRequest
	Data []Player `json:"data"`
}

func (r RequestAccountList) PlayerList() []Player {
	return r.Data
}

type RequestAccountInfo struct {
	NoVehicleList
	NoClanList
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

type RequestAccountTanks struct {
	NoPlayerList
	NoClanList
	ApiRequest
	Data map[string][]Vehicle `json:"data"`
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

type RequestClanList struct {
	NoPlayerList
	NoVehicleList
	ApiRequest
	Data []Clan `json:"data"`
}

func (r RequestClanList) ClanList() []Clan {
	return r.Data
}

type RequestClanInfo struct {
	NoPlayerList
	NoVehicleList
	ApiRequest
	Data map[string]Clan `json:"data"`
}

func (r RequestClanInfo) ClanList() []Clan {
	Clans := []Clan{}
	for _, v := range r.Data {
		// don't append 'empty' clans (those that aren't found), since we don't return the key
		if v.ClanId != 0 {
			Clans = append(Clans, v)
		}
	}
	return Clans
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
	regional["eu"] = "api.worldoftanks.eu/"
	regional["na"] = "api.worldoftanks.com/"
	regional["ru"] = "api.worldoftanks.ru/"

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
	case "wot/account/list":
		reqt := RequestAccountList{}
		err = json.Unmarshal(data, &reqt)
		req = reqt

	case "wot/account/info":
		reqt := RequestAccountInfo{}
		err = json.Unmarshal(data, &reqt)
		req = reqt

	case "wot/account/tanks":
		reqt := RequestAccountTanks{}
		err = json.Unmarshal(data, &reqt)
		req = reqt
	case "wgn/clans/list":
		reqt := RequestClanList{}
		err = json.Unmarshal(data, &reqt)
		req = reqt
	case "wgn/clans/info":
		reqt := RequestClanInfo{}
		err = json.Unmarshal(data, &reqt)
		req = reqt
	default:
		panic("Unknown command, this shouldn't happen")
	}

	// no errors so we return the data retrieved
	if err == nil {
		return req, err
	}
	// if an error is found check whether it's a WG API error by unmarshaling the data in an ApiError struct,
	// so we can pass it along
	apiError := ApiRequestError{}
	ae := json.Unmarshal(data, &apiError)
	if ae == nil && apiError.Status == "error" {
		// wrap any WG API error around possible previous error
		apiError.Sibling = err
		// and we return the API error
		err = apiError
	}

	//otherwise return just the error
	return nil, err

}

// wot/account/list
func (w *WG) SearchPlayersByName(name string, exact bool) (RequestAccountList, error) {
	params := make(map[string]string)
	if exact {
		params["type"] = "exact"
	}
	params["search"] = name
	result, err := w.apiCall("wot/account/list", params)
	if err != nil {
		return RequestAccountList{}, err
	}
	return result.(RequestAccountList), err

}

// wot/account/info
func (w *WG) GetPlayerInfo(accountid []uint32) (RequestAccountInfo, error) {
	params := make(map[string]string)

	for _, v := range accountid {
		_, ok := params["account_id"]
		if len(params["account_id"]) == 0 || !ok {
			params["account_id"] = fmt.Sprint(v)
		} else {
			params["account_id"] += "," + fmt.Sprint(v)
		}
	}

	result, err := w.apiCall("wot/account/info", params)
	if err != nil {
		return RequestAccountInfo{}, err
	}
	return result.(RequestAccountInfo), err

}

// wot/account/tanks
func (w *WG) GetPlayerTanks(accountid []uint32, vehicleid []uint32) (RequestAccountTanks, error) {
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
	result, err := w.apiCall("wot/account/tanks", params)
	if err != nil {
		return RequestAccountTanks{}, err
	}
	return result.(RequestAccountTanks), err

}

// wgn/clans/list
func (w *WG) SearchClansByName(name string, orderby string, pageno uint32, limit uint32) (RequestClanList, error) {
	params := make(map[string]string)
	params["search"] = name
	params["order_by"] = orderby
	params["page_no"] = fmt.Sprint(pageno)
	params["limit"] = fmt.Sprint(limit)
	result, err := w.apiCall("wgn/clans/list", params)
	if err != nil {
		return RequestClanList{}, err
	}
	return result.(RequestClanList), err
}

// wgn/clans/info
func (w *WG) GetClanInfo(clanid []uint32, accessToken string) (RequestClanInfo, error) {
	params := make(map[string]string)
	if len(accessToken) > 0 {
		params["access_token"] = accessToken
	}
	for _, v := range clanid {
		_, ok := params["clan_id"]
		if len(params["clan_id"]) == 0 || !ok {
			params["clan_id"] = fmt.Sprint(v)
		} else {
			params["clan_id"] += "," + fmt.Sprint(v)
		}
	}
	result, err := w.apiCall("wgn/clans/info", params)
	if err != nil {
		return RequestClanInfo{}, err
	}
	return result.(RequestClanInfo), err
}
