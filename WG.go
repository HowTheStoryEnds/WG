package WG

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	//"strconv"
	//"github.com/davecgh/go-spew/spew"
)

const MasteryNone uint32 = 0
const MasteryThirdClass uint32 = 1
const MasterySecondClass uint32 = 2
const MasteryFirstClass uint32 = 3
const MasteryAceTanker uint32 = 4

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

func (w *WG) constructURL() string {
	//sane defaults
	if w.transport != "http" {
		w.transport = "https"
	}

	var regional = make(map[string]string)
	regional["eu"] = "api.worldoftanks.eu/wot/"
	regional["na"] = "api.worldoftanks.com/wot/"
	regional["ru"] = "api.worldoftanks.ru/wot/"

	return w.transport + "://" + regional[w.region]

}
func (w *WG) addGetParams(uri string, parameters map[string][]string) string {

	baseUrl, err := url.Parse(uri)
	if err != nil {
		//panic!

	}

	params := url.Values{}
	for key, v := range parameters {
		for _, value := range v {
			params.Add(key, value)
		}
	}

	baseUrl.RawQuery = params.Encode()
	return uri + "?" + baseUrl.RawQuery
}

func (w *WG) retrieveData(action string, parameters map[string][]string) (map[string]interface{}, error) {

	// make sure the apikey is set
	if _, present := parameters["application_id"]; !present {
		parameters["application_id"] = []string{w.apiKey}
	}

	//combine the identical parameter in a , (comma) seperated string
	for k, _ := range parameters {
		params := ""
		for _, iv := range parameters[k] {
			if params == "" {
				params = iv
			} else {
				params += "," + iv
			}

		}
		parameters[k] = []string{params}
	}

	//construct the full URL
	uri := w.addGetParams(w.constructURL()+action+"/", parameters)

	//fmt.Println("requesting uri: " + uri)

	// and retrieve the JSON data in a string
	var val, err = http.Get(uri)
	var res = []byte{}
	if err == nil {
		res, err = ioutil.ReadAll(val.Body)
		val.Body.Close()
	}

	var result map[string]interface{}
	if err == nil {
		err = json.Unmarshal(res, &result)
	}
	return result, err
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

type Tank struct {
	Name          string
	InGarage      bool
	AvgDmg        uint32
	AvgSpotting   uint32
	AvgExp        uint32
	Battles       uint32
	BattlesLost   uint32
	BattlesWon    uint32
	Winrate       uint32
	Wins          uint32
	MarkOfMastery uint32
	TankId        uint32
}

type Player struct {
	AccountId        uint32
	Nickname         string
	Region           string
	Tanks            []Tank
	Statistics       map[string]PlayerStatistics
	MaxFragsTankId   uint32
	MaxXpTankId      uint32
	MaxXp            uint32
	TreesCut         uint32
	MaxFrags         uint32
	MaxDamageTankId  uint32
	Frags            uint32
	MaxDamage        uint32
	MaxDamageVehicle uint32
	CreatedAt        uint32
	UpdatedAt        uint32
	Private          PlayerPrivate
	GlobalRating     uint32
	ClanId           uint32
	LastBattleTime   uint32
	LogoutAt         uint32
}

type PlayerStatistics struct {
	Spotted                    uint32
	AvgDamageAssistedTrack     float64
	AvgDamageBlocked           float64
	DirectHitsReceived         uint32
	ExplosionHits              uint32
	PiercingsReceived          uint32
	Piercings                  uint32
	Xp                         uint32
	SurvivedBattles            uint32
	DroppedCapturePoints       uint32
	HitsPercents               uint32
	Draws                      uint32
	Battles                    uint32
	DamageReceived             uint32
	AvgDamageAssisted          float64
	Frags                      uint32
	AvgDamageAssistedRadio     float64
	CapturePoints              uint32
	BaseXp                     uint32
	Hits                       uint32
	BattleAvgXp                uint32
	Wins                       uint32
	Losses                     uint32
	DamageDealt                uint32
	NoDamageDirectHitsReceived uint32
	Shots                      uint32
	ExplosionHitsReceived      uint32
	TankingFactor              float64
	MaxXp                      uint32
	MaxDamage                  uint32
	MaxFragsTankId             uint32
	MaxXpTankId                uint32
	MaxDamageTankId            uint32
	MaxFrags                   uint32
}

type PlayerPrivate struct {
	BanInfo                   string
	BanTime                   uint32
	BattleLifeTime            uint32
	Credits                   uint32
	FreeXp                    uint32
	Friends                   []uint32
	Gold                      uint32
	IsBoundToPhone            bool
	IsPremium                 bool
	PersonalMissions          map[string]string
	PremiumExpiresAt          uint32
	GroupedContactsBlocked    []uint32
	GroupedContactsGroups     map[string][]uint32
	GroupedContactsIgnored    []uint32
	GroupedContactsMuted      []uint32
	GroupedContactsUngrouped  []uint32
	RentedCompensationCredits uint32
	RentedCompensationGold    uint32
	RentedExpirationTime      uint32
	RentedTankId              uint32
	RestrictionsChatBanTime   uint32
	RestrictionsClanTime      uint32
}

func (w *WG) processTank(t map[string]interface{}) (interface{}, bool) {
	var tank = Tank{}
	ok := true

	for k, v := range t {
		if v == nil {
			continue
		}
		switch k {
		case "statistics":
			for sk, sv := range v.(map[string]interface{}) {
				if sv == nil {
					continue
				}
				switch sk {
				case "wins":
					tank.Wins = uint32(sv.(float64))
				case "battles":
					tank.Battles = uint32(sv.(float64))
				}
			}
		case "mark_of_mastery":
			tank.MarkOfMastery = uint32(v.(float64))
		case "tank_id":
			tank.TankId = uint32(v.(float64))

		}
	}

	// do not return empty Tank structures
	if tank.TankId == 0 {
		ok = false
	}

	return tank, ok

}

func (w *WG) processPlayer(p map[string]interface{}) (interface{}, bool) {
	//nickname := result["data"].([]interface{})[x].(map[string]interface{})["nickname"].(string)
	var player = Player{}
	ok := true
	//var p = make(map[string]interface{})
	//var p map[string]interface{} = result["data"].([]interface{})[x].(map[string]interface{})

	// insert empty holder, tanks are to be retrieved by different API call
	player.Tanks = []Tank{}
	for k, v := range p {
		if v == nil {
			continue
		}
		switch k {
		case "nickname":
			player.Nickname = v.(string)
		case "account_id":
			player.AccountId = uint32(v.(float64))
		case "last_battle_time":
			player.LastBattleTime = uint32(v.(float64))
		case "created_at":
			player.CreatedAt = uint32(v.(float64))
		case "updated_at":
			player.UpdatedAt = uint32(v.(float64))
		case "private":
			//TODO: Private data insertion
			player.Private = PlayerPrivate{}
		case "global_rating":
			player.GlobalRating = uint32(v.(float64))
		case "clan_id":
			player.ClanId = uint32(v.(float64))
		case "logout_at":
			player.LogoutAt = uint32(v.(float64))
		case "statistics":
			player.Statistics = make(map[string]PlayerStatistics)
			for sk, sv := range v.(map[string]interface{}) {
				if sv == nil {
					continue
				}
				switch sk {
				case "clan", "regular_team", "company", "all", "stronghold_defense", "stronghold_skirmish", "historical", "team":
					ps := PlayerStatistics{}
					for ssk, ssv := range sv.(map[string]interface{}) {
						if ssv == nil {
							continue
						}
						switch ssk {
						case "spotted":
							ps.Spotted = uint32(ssv.(float64))
						case "avg_damage_assisted_track":
							ps.AvgDamageAssistedTrack = ssv.(float64)
						case "max_xp":
							ps.MaxXp = uint32(ssv.(float64))
						case "avg_damage_blocked":
							ps.AvgDamageBlocked = ssv.(float64)
						case "direct_hits_received":
							ps.DirectHitsReceived = uint32(ssv.(float64))
						case "explosion_hits":
							ps.ExplosionHits = uint32(ssv.(float64))
						case "piercings_received":
							ps.PiercingsReceived = uint32(ssv.(float64))
						case "piercings":
							ps.Piercings = uint32(ssv.(float64))
						case "max_damage_tank_id":
							ps.MaxDamageTankId = uint32(ssv.(float64))
						case "xp":
							ps.Xp = uint32(ssv.(float64))
						case "survived_battles":
							ps.SurvivedBattles = uint32(ssv.(float64))
						case "dropped_capture_points":
							ps.DroppedCapturePoints = uint32(ssv.(float64))
						case "hits_percents":
							ps.HitsPercents = uint32(ssv.(float64))
						case "draws":
							ps.Draws = uint32(ssv.(float64))
						case "max_xp_tank_id":
							ps.MaxXpTankId = uint32(ssv.(float64))
						case "battles":
							ps.Battles = uint32(ssv.(float64))
						case "damage_received":
							ps.DamageReceived = uint32(ssv.(float64))
						case "avg_damage_assisted":
							ps.AvgDamageAssisted = ssv.(float64)
						case "max_frags_tank_id":
							ps.MaxFragsTankId = uint32(ssv.(float64))
						case "frags":
							ps.Frags = uint32(ssv.(float64))
						case "avg_damage_assisted_radio":
							ps.AvgDamageAssistedRadio = ssv.(float64)
						case "capture_points":
							ps.CapturePoints = uint32(ssv.(float64))
						case "max_damage":
							ps.MaxDamage = uint32(ssv.(float64))
						case "hits":
							ps.Hits = uint32(ssv.(float64))
						case "battle_avg_xp":
							ps.BattleAvgXp = uint32(ssv.(float64))
						case "wins":
							ps.Wins = uint32(ssv.(float64))
						case "losses":
							ps.Losses = uint32(ssv.(float64))
						case "damage_dealt":
							ps.DamageDealt = uint32(ssv.(float64))
						case "no_damage_direct_hits_received":
							ps.NoDamageDirectHitsReceived = uint32(ssv.(float64))
						case "max_frags":
							ps.MaxFrags = uint32(ssv.(float64))
						case "shots":
							ps.Shots = uint32(ssv.(float64))
						case "explosion_hits_received":
							ps.ExplosionHitsReceived = uint32(ssv.(float64))
						case "tanking_factor":
							ps.TankingFactor = ssv.(float64)

						}
					}

					player.Statistics[sk] = ps

				case "max_frags_tank_id":
					player.MaxFragsTankId = uint32(sv.(float64))

				case "max_xp_tank_id":
					player.MaxXpTankId = uint32(sv.(float64))
				case "max_xp":
					player.MaxXp = uint32(sv.(float64))
				case "trees_cut":
					player.TreesCut = uint32(sv.(float64))
				case "max_frags":
					player.MaxFrags = uint32(sv.(float64))
				case "max_damage_tank_id":
					player.MaxDamageTankId = uint32(sv.(float64))
				case "frags":
					player.Frags = uint32(sv.(float64))
				case "max_damage":
					player.MaxDamage = uint32(sv.(float64))
				case "max_damage_vehicle":
					player.MaxDamageVehicle = uint32(sv.(float64))

				}
			}

		}

	}

	player.Region = w.region
	// do not return empty Player structures
	if player.AccountId == 0 {
		ok = false
	}
	return player, ok
}

func (w *WG) resultsToData(result map[string]interface{}, toCallMap func(map[string]interface{}) (interface{}, bool), toCallArray func([]interface{}) (interface{}, bool)) []interface{} {
	//v := result["meta"].(map[string]interface{})
	//var fc float64 = v["count"].(float64)
	//var count uint = uint(fc)

	//fmt.Println("we found " + string(count) + "/" + strconv.FormatFloat(fc, 'f', -1, 64) + " players")

	var content []interface{}

	_, found := result["data"].([]interface{})
	if found {
		for _, v := range result["data"].([]interface{}) {
			data, ok := toCallMap(v.(map[string]interface{}))
			if ok {
				content = append(content, data)
			}
		}
	}
	_, found = result["data"].(map[string]interface{})
	if found {

		for _, v := range result["data"].(map[string]interface{}) {
			_, contentFound := v.(map[string]interface{})
			if contentFound {
				data, ok := toCallMap(v.(map[string]interface{}))
				if ok {
					content = append(content, data)
				}
			}
			_, contentFound = v.([]interface{})
			if contentFound {
				for _, vv := range v.([]interface{}) {
					data, ok := toCallMap(vv.(map[string]interface{}))
					if ok {
						content = append(content, data)
					}
				}

			}
		}
	}
	return content
}

func (w *WG) resultsToPlayer(results []interface{}) (solution []Player) {
	// make sure we always return an empty array and not a nil array
	solution = []Player{}
	for _, r := range results {
		solution = append(solution, r.(Player))

	}
	return
}
func (w *WG) resultsToTank(results []interface{}) (solution []Tank) {
	// make sure we always return an empty array and not a nil array
	solution = []Tank{}
	for _, r := range results {
		solution = append(solution, r.(Tank))

	}
	return
}

// account/list
func (w *WG) SearchPlayersByName(name string, exact bool) []Player {
	params := make(map[string][]string)
	if exact {
		params["type"] = []string{"exact"}
	}
	params["search"] = []string{name}

	var result, err = w.retrieveData("account/list", params)
	if err != nil {
		return []Player{}
	}

	return w.resultsToPlayer(w.resultsToData(result, w.processPlayer, nil))

}

// account/info
func (w *WG) GetPlayerPersonalData(accountid []uint32) (PlayersFound []Player) {
	params := make(map[string][]string)
	var idHolder []string
	for _, v := range accountid {
		idHolder = append(idHolder, fmt.Sprint(v))
	}

	params["account_id"] = idHolder

	var result, err = w.retrieveData("account/info", params)
	if err != nil {
		//TODO panic instead?
		fmt.Println("GetPlayerPersonalData: " + err.Error())
		return []Player{}
	}

	return w.resultsToPlayer(w.resultsToData(result, w.processPlayer, nil))

}

func (w *WG) GetPlayerTanks(accountid uint32, vehicleid []uint32) []Tank {
	params := make(map[string][]string)
	var idHolder []string
	for _, v := range vehicleid {
		idHolder = append(idHolder, fmt.Sprint(v))
	}

	if len(idHolder) > 0 {
		params["tank_id"] = idHolder
	}
	params["account_id"] = []string{fmt.Sprint(accountid)}
	var result, err = w.retrieveData("account/tanks", params)
	if err != nil {
		//TODO panic instead?
		fmt.Println("GetPlayerTanks: " + err.Error())
		return []Tank{}
	}

	return w.resultsToTank(w.resultsToData(result, w.processTank, nil))

}
