package WG

import (
	//"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/jarcoal/httpmock"
	. "gopkg.in/check.v1"
	"io/ioutil"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type WGSuite struct {
	Wg WG
}

var _ = Suite(&WGSuite{})

type ResponderData struct {
	Uri     []string
	Content []byte
}

func (s *WGSuite) SetUpSuite(c *C) {

	var res = make([]ResponderData, 0, 0)
	var rd = ResponderData{}

	/*
	 *     account/list
	 */
	// search player single result found && search single player with exact name
	rd.Uri = []string{"https://api.worldoftanks.eu/wot/account/list/?application_id=demo&search=howthestoryends",
		"https://api.worldoftanks.eu/wot/account/list/?application_id=demo&search=howthestoryends&type=exact"}
	rd.Content, _ = ioutil.ReadFile("./testdata/account/list/howthestoryends_name.json")
	res = append(res, rd)

	rd.Uri = []string{"https://api.worldoftanks.eu/wot/account/list/?application_id=demo&search=howthe"}
	rd.Content, _ = ioutil.ReadFile("./testdata/account/list/howthe_name.json")
	res = append(res, rd)

	/*
	 *           account/info
	 */
	// single player
	rd.Uri = []string{"https://api.worldoftanks.eu/wot/account/info/?account_id=507197901&application_id=demo"}
	rd.Content, _ = ioutil.ReadFile("./testdata/account/info/single_player.json")
	res = append(res, rd)
	//single clanless player
	rd.Uri = []string{"https://api.worldoftanks.eu/wot/account/info/?account_id=525427444&application_id=demo"}
	rd.Content, _ = ioutil.ReadFile("./testdata/account/info/single_player_clanless.json")
	res = append(res, rd)
	// 2 players, 1 can not be found
	rd.Uri = []string{"https://api.worldoftanks.eu/wot/account/info/?account_id=507197901%2C1&application_id=demo"}
	rd.Content, _ = ioutil.ReadFile("./testdata/account/info/two_players_only_one_found.json")
	res = append(res, rd)
	// 2 players, both found
	rd.Uri = []string{"https://api.worldoftanks.eu/wot/account/info/?account_id=525427444%2C507197901&application_id=demo"}
	rd.Content, _ = ioutil.ReadFile("./testdata/account/info/two_players.json")
	res = append(res, rd)

	/*
	 *   account/tanks
	 */
	// single vehicle
	rd.Uri = []string{"https://api.worldoftanks.eu/wot/account/tanks/?account_id=507197901&application_id=demo&tank_id=1"}
	rd.Content, _ = ioutil.ReadFile("./testdata/account/tanks/single_vehicle.json")
	res = append(res, rd)
	// three vehicles
	rd.Uri = []string{"https://api.worldoftanks.eu/wot/account/tanks/?account_id=507197901&application_id=demo&tank_id=11601%2C3089%2C11777"}
	rd.Content, _ = ioutil.ReadFile("./testdata/account/tanks/three_vehicles.json")
	res = append(res, rd)
	// vehicle not found
	rd.Uri = []string{"https://api.worldoftanks.eu/wot/account/tanks/?account_id=507197901&application_id=demo&tank_id=2"}
	rd.Content, _ = ioutil.ReadFile("./testdata/account/tanks/vehicle_not_found.json")
	res = append(res, rd)
	// complete vehicle list
	rd.Uri = []string{"https://api.worldoftanks.eu/wot/account/tanks/?account_id=507197901&application_id=demo"}
	rd.Content, _ = ioutil.ReadFile("./testdata/account/tanks/complete_vehicle_list.json")
	res = append(res, rd)

	//setup HTTP mocking service
	httpmock.Activate()
	//setup the urls with their content
	for _, v := range res {
		for _, url := range v.Uri {
			httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(200, string(v.Content)))
		}
	}

	// initialize a WG api
	s.Wg.Init("na", "https", "get", "demo")
}

func (s *WGSuite) TearDownSuite(c *C) {
	//shutdown HTTP mocking
	httpmock.DeactivateAndReset()
}

// tests
func (s *WGSuite) TestConstructURL(c *C) {

	s.Wg.SetRegion("na")
	c.Check(s.Wg.constructURL(), Equals, "https://api.worldoftanks.com/wot/")
	s.Wg.SetRegion("eu")
	c.Check(s.Wg.constructURL(), Equals, "https://api.worldoftanks.eu/wot/")

	s.Wg.SetTransport("http")
	c.Check(s.Wg.constructURL(), Equals, "http://api.worldoftanks.eu/wot/")

}

func (s *WGSuite) TestaddGetParams(c *C) {
	s.Wg.SetTransport("http")
	// one key, one parameter
	var params = map[string][]string{"one": {"1"},
		"two": {"twee"},
	}
	c.Check(s.Wg.addGetParams("http://www.example.com", params), Equals, "http://www.example.com?one=1&two=twee")

	// multiple parameters to the same key
	params = map[string][]string{"dienaren": {"Harald", "Babs"},
		"dieren": {"Rickey", "Tiger", "Storm", "Mazzeltje"},
	}
	c.Check(s.Wg.addGetParams("http://www.example.com", params),
		Equals,
		"http://www.example.com?dienaren=Harald&dienaren=Babs&dieren=Rickey&dieren=Tiger&dieren=Storm&dieren=Mazzeltje")

	// special characters need to be escaped
	params = map[string][]string{"1with_slashes": {"co//ol"}, "2email": {"harald.brinkhof@gmail.com"}}
	c.Check(s.Wg.addGetParams("http://www.example.com", params), Equals, "http://www.example.com?1with_slashes=co%2F%2Fol&2email=harald.brinkhof%40gmail.com")

}

func (s *WGSuite) TestretrieveData(c *C) {
	s.Wg.SetTransport("https")
	ret, err := s.Wg.retrieveData("account/list", map[string][]string{"type": {"exact"}, "search": {"howthestoryends"}})
	var data = make(map[string]interface{})
	data["status"] = "ok"
	data["meta"] = make(map[string]interface{})
	data["meta"].(map[string]interface{})["count"] = 1.0
	data["data"] = make([]interface{}, 1, 1)
	data["data"].([]interface{})[0] = make(map[string]interface{})
	data["data"].([]interface{})[0].(map[string]interface{})["nickname"] = "HowTheStoryEnds"
	data["data"].([]interface{})[0].(map[string]interface{})["account_id"] = 507197901.0

	c.Check(ret, DeepEquals, data)

	// unretrievable data
	ret, err = s.Wg.retrieveData("@/@", map[string][]string{})
	c.Check(err, NotNil)
	var notFound map[string]interface{}
	c.Check(ret, DeepEquals, notFound)

}

func (s *WGSuite) TestSearchPlayersByName(c *C) {
	s.Wg.SetTransport("https")

	// 'startswith' search yielding 1 result
	var player = []Player{{Nickname: "HowTheStoryEnds", AccountId: 507197901, Region: "eu", Tanks: []Tank{}}}
	var data = s.Wg.SearchPlayersByName("howthestoryends", false)
	c.Check(player, DeepEquals, data)

	// exact search yielding 1 result
	data2 := s.Wg.SearchPlayersByName("howthestoryends", true)
	c.Check(player, DeepEquals, data2)

	// search that can yield multiple results
	var players = []Player{}
	var data3 []Player = s.Wg.SearchPlayersByName("howthe", false)
	for _, p := range data3 {
		switch p.Nickname {
		case "howtheblank":
			players = append(players, Player{Nickname: "howtheblank", AccountId: 502301211, Region: s.Wg.region, Tanks: []Tank{}})
		case "HowTheGodsKill":
			players = append(players, Player{Nickname: "HowTheGodsKill", AccountId: 525427444, Region: s.Wg.region, Tanks: []Tank{}})
		case "HowTheGuy":
			players = append(players, Player{Nickname: "HowTheGuy", AccountId: 506219127, Region: s.Wg.region, Tanks: []Tank{}})
		case "HowTheStoryEnds":
			players = append(players, Player{Nickname: "HowTheStoryEnds", AccountId: 507197901, Region: s.Wg.region, Tanks: []Tank{}})
		}
	}
	c.Check(players, DeepEquals, data3)

	//no result found
	data4 := s.Wg.SearchPlayersByName("howthestt", false)
	c.Check([]Player{}, DeepEquals, data4)
}

func (s *WGSuite) TestGetPlayerPersonalData(c *C) {
	s.Wg.SetTransport("https")

	// player in a clan
	var HowTheStoryEnds = Player{Nickname: "HowTheStoryEnds", AccountId: 507197901, Region: "eu"}
	HowTheStoryEnds.Frags = 0
	HowTheStoryEnds.MaxXp = 2790
	HowTheStoryEnds.Region = "eu"
	HowTheStoryEnds.TreesCut = 21573
	HowTheStoryEnds.MaxFrags = 11
	HowTheStoryEnds.MaxDamage = 6017
	HowTheStoryEnds.MaxXpTankId = 11265
	HowTheStoryEnds.MaxFragsTankId = 2369
	HowTheStoryEnds.MaxDamageTankId = 7425
	HowTheStoryEnds.MaxDamageVehicle = 7425
	HowTheStoryEnds.CreatedAt = 1355749511
	HowTheStoryEnds.UpdatedAt = 1434133471
	HowTheStoryEnds.Private = PlayerPrivate{}
	HowTheStoryEnds.GlobalRating = 7728
	HowTheStoryEnds.ClanId = 500010805
	HowTheStoryEnds.Tanks = []Tank{}
	HowTheStoryEnds.Statistics = map[string]PlayerStatistics{}
	HowTheStoryEnds.LastBattleTime = 1434126614
	HowTheStoryEnds.LogoutAt = 1434133469

	clan := PlayerStatistics{}
	clan.Spotted = 36
	clan.AvgDamageAssistedTrack = 45.57
	clan.AvgDamageBlocked = 219.55
	clan.DirectHitsReceived = 121
	clan.ExplosionHits = 15
	clan.PiercingsReceived = 99
	clan.Piercings = 132
	clan.Xp = 37281
	clan.SurvivedBattles = 33
	clan.DroppedCapturePoints = 2
	clan.HitsPercents = 57
	clan.Draws = 2
	clan.Battles = 55
	clan.DamageReceived = 33255
	clan.AvgDamageAssisted = 184.95
	clan.Frags = 36
	clan.AvgDamageAssistedRadio = 139.39
	clan.CapturePoints = 26
	clan.Hits = 195
	clan.BattleAvgXp = 678
	clan.Wins = 46
	clan.Losses = 7
	clan.DamageDealt = 35846
	clan.NoDamageDirectHitsReceived = 22
	clan.Shots = 343
	clan.ExplosionHitsReceived = 1
	clan.TankingFactor = 0.35
	HowTheStoryEnds.Statistics["clan"] = clan

	HowTheStoryEnds.Statistics["regular_team"] = PlayerStatistics{}

	company := PlayerStatistics{}
	company.Spotted = 190
	company.AvgDamageAssistedTrack = 6.25
	company.AvgDamageBlocked = 0
	company.DirectHitsReceived = 52
	company.ExplosionHits = 0
	company.PiercingsReceived = 48
	company.Piercings = 58
	company.Xp = 82190
	company.SurvivedBattles = 60
	company.DroppedCapturePoints = 156
	company.HitsPercents = 61
	company.Draws = 3
	company.Battles = 152
	company.DamageReceived = 74793
	company.AvgDamageAssisted = 165.58
	company.Frags = 98
	company.AvgDamageAssistedRadio = 159.33
	company.CapturePoints = 229
	company.Hits = 502
	company.BattleAvgXp = 541
	company.Wins = 97
	company.Losses = 52
	company.DamageDealt = 81272
	company.NoDamageDirectHitsReceived = 4
	company.Shots = 819
	company.ExplosionHitsReceived = 0
	company.TankingFactor = 0.0
	HowTheStoryEnds.Statistics["company"] = company

	strongholdDefense := PlayerStatistics{}
	strongholdDefense.Spotted = 41
	strongholdDefense.MaxFragsTankId = 7169
	strongholdDefense.MaxXp = 1062
	strongholdDefense.DirectHitsReceived = 177
	strongholdDefense.ExplosionHits = 0
	strongholdDefense.PiercingsReceived = 121
	strongholdDefense.Piercings = 121
	strongholdDefense.Xp = 21559
	strongholdDefense.SurvivedBattles = 17
	strongholdDefense.DroppedCapturePoints = 31
	strongholdDefense.HitsPercents = 74
	strongholdDefense.Draws = 0
	strongholdDefense.MaxXpTankId = 10785
	strongholdDefense.Battles = 33
	strongholdDefense.DamageReceived = 47316
	strongholdDefense.Frags = 17
	strongholdDefense.CapturePoints = 38
	strongholdDefense.MaxDamageTankId = 12369
	strongholdDefense.MaxDamage = 3900
	strongholdDefense.Hits = 179
	strongholdDefense.BattleAvgXp = 653
	strongholdDefense.Wins = 27
	strongholdDefense.Losses = 6
	strongholdDefense.DamageDealt = 45584
	strongholdDefense.NoDamageDirectHitsReceived = 56
	strongholdDefense.MaxFrags = 2
	strongholdDefense.Shots = 242
	strongholdDefense.ExplosionHitsReceived = 26
	strongholdDefense.TankingFactor = 0.33
	HowTheStoryEnds.Statistics["stronghold_defense"] = strongholdDefense

	strongholdSkirmish := PlayerStatistics{}
	strongholdSkirmish.Spotted = 1350
	strongholdSkirmish.MaxFragsTankId = 4913
	strongholdSkirmish.MaxXp = 1980
	strongholdSkirmish.DirectHitsReceived = 7982
	strongholdSkirmish.ExplosionHits = 224
	strongholdSkirmish.PiercingsReceived = 5722
	strongholdSkirmish.Piercings = 8514
	strongholdSkirmish.Xp = 1899689
	strongholdSkirmish.SurvivedBattles = 1413
	strongholdSkirmish.DroppedCapturePoints = 1522
	strongholdSkirmish.HitsPercents = 69
	strongholdSkirmish.Draws = 9
	strongholdSkirmish.MaxXpTankId = 4385
	strongholdSkirmish.Battles = 2158
	strongholdSkirmish.DamageReceived = 1750773
	strongholdSkirmish.Frags = 1399
	strongholdSkirmish.CapturePoints = 2134
	strongholdSkirmish.MaxDamageTankId = 7169
	strongholdSkirmish.MaxDamage = 4144
	strongholdSkirmish.Hits = 11663
	strongholdSkirmish.BattleAvgXp = 880
	strongholdSkirmish.Wins = 1947
	strongholdSkirmish.Losses = 202
	strongholdSkirmish.DamageDealt = 2120855
	strongholdSkirmish.NoDamageDirectHitsReceived = 2260
	strongholdSkirmish.MaxFrags = 5
	strongholdSkirmish.Shots = 16898
	strongholdSkirmish.ExplosionHitsReceived = 146
	strongholdSkirmish.TankingFactor = 0.45
	HowTheStoryEnds.Statistics["stronghold_skirmish"] = strongholdSkirmish

	HowTheStoryEnds.Statistics["historical"] = PlayerStatistics{}

	team := PlayerStatistics{}
	team.Spotted = 615
	team.AvgDamageAssistedTrack = 66.12
	team.MaxXp = 1848
	team.AvgDamageBlocked = 224.65
	team.DirectHitsReceived = 2552
	team.ExplosionHits = 1
	team.PiercingsReceived = 1851
	team.Piercings = 2259
	team.MaxDamageTankId = 5377
	team.Xp = 474872
	team.SurvivedBattles = 278
	team.DroppedCapturePoints = 1301
	team.HitsPercents = 71
	team.Draws = 11
	team.MaxXpTankId = 5377
	team.Battles = 651
	team.DamageReceived = 392013
	team.AvgDamageAssisted = 172.07
	team.MaxFragsTankId = 4385
	team.Frags = 378
	team.AvgDamageAssistedRadio = 105.95
	team.CapturePoints = 1958
	team.MaxDamage = 3526
	team.Hits = 3268
	team.BattleAvgXp = 729
	team.Wins = 491
	team.Losses = 149
	team.DamageDealt = 517046
	team.NoDamageDirectHitsReceived = 701
	team.MaxFrags = 5
	team.Shots = 4603
	team.ExplosionHitsReceived = 2
	team.TankingFactor = 0.37
	HowTheStoryEnds.Statistics["team"] = team

	all := PlayerStatistics{}
	all.Spotted = 24553
	all.AvgDamageAssistedTrack = 66.92
	all.MaxXp = 2790
	all.AvgDamageBlocked = 222.86
	all.DirectHitsReceived = 50646
	all.ExplosionHits = 4000
	all.PiercingsReceived = 37905
	all.Piercings = 79506
	all.MaxDamageTankId = 7425
	all.Xp = 12796688
	all.SurvivedBattles = 6757
	all.DroppedCapturePoints = 23554
	all.HitsPercents = 59
	all.Draws = 226
	all.MaxXpTankId = 11265
	all.Battles = 19468
	all.DamageReceived = 9471075
	all.AvgDamageAssisted = 307.41
	all.MaxFragsTankId = 2369
	all.Frags = 24133
	all.AvgDamageAssistedRadio = 240.49
	all.CapturePoints = 19541
	all.MaxDamage = 6017
	all.Hits = 159851
	all.BattleAvgXp = 657
	all.Wins = 11445
	all.Losses = 7797
	all.DamageDealt = 16903164
	all.NoDamageDirectHitsReceived = 12741
	all.MaxFrags = 11
	all.Shots = 270136
	all.ExplosionHitsReceived = 1564
	all.TankingFactor = 0.35
	HowTheStoryEnds.Statistics["all"] = all

	c.Check(s.Wg.GetPlayerPersonalData([]uint32{507197901}), DeepEquals, []Player{HowTheStoryEnds})

	// id 1 does not exist so what should be returned is only 1 player record for 507197901
	c.Check(s.Wg.GetPlayerPersonalData([]uint32{507197901, 1}), DeepEquals, []Player{HowTheStoryEnds})

	// clanless player
	var HowTheGodsKill = Player{Nickname: "HowTheGodsKill", AccountId: 525427444, Region: "eu"}
	HowTheGodsKill.Frags = 0
	HowTheGodsKill.MaxXp = 1407
	HowTheGodsKill.Region = "eu"
	HowTheGodsKill.TreesCut = 42
	HowTheGodsKill.MaxFrags = 5
	HowTheGodsKill.MaxDamage = 1354
	HowTheGodsKill.MaxXpTankId = 4945
	HowTheGodsKill.MaxFragsTankId = 4945
	HowTheGodsKill.MaxDamageTankId = 4945
	HowTheGodsKill.MaxDamageVehicle = 4945
	HowTheGodsKill.CreatedAt = 1420581930
	HowTheGodsKill.UpdatedAt = 1433267697
	HowTheGodsKill.Private = PlayerPrivate{}
	HowTheGodsKill.GlobalRating = 1263
	HowTheGodsKill.Tanks = []Tank{}
	HowTheGodsKill.Statistics = map[string]PlayerStatistics{}
	HowTheGodsKill.LastBattleTime = 1423093994
	HowTheGodsKill.LogoutAt = 1433267693

	HowTheGodsKill.Statistics["clan"] = PlayerStatistics{}
	HowTheGodsKill.Statistics["regular_team"] = PlayerStatistics{}
	HowTheGodsKill.Statistics["company"] = PlayerStatistics{}
	HowTheGodsKill.Statistics["stronghold_defense"] = PlayerStatistics{}
	HowTheGodsKill.Statistics["stronghold_skirmish"] = PlayerStatistics{}
	HowTheGodsKill.Statistics["historical"] = PlayerStatistics{}
	HowTheGodsKill.Statistics["team"] = PlayerStatistics{}

	all = PlayerStatistics{}
	all.Spotted = 53
	all.AvgDamageAssistedTrack = 4.34
	all.MaxXp = 1407
	all.AvgDamageBlocked = 51.58
	all.DirectHitsReceived = 552
	all.ExplosionHits = 0
	all.PiercingsReceived = 356
	all.Piercings = 462
	all.MaxDamageTankId = 4945
	all.Xp = 19371
	all.SurvivedBattles = 14
	all.DroppedCapturePoints = 43
	all.HitsPercents = 57
	all.Draws = 0
	all.MaxXpTankId = 4945
	all.Battles = 67
	all.DamageReceived = 10986
	all.AvgDamageAssisted = 25.39
	all.MaxFragsTankId = 4945
	all.Frags = 70
	all.AvgDamageAssistedRadio = 21.04
	all.CapturePoints = 83
	all.MaxDamage = 1354
	all.Hits = 691
	all.BattleAvgXp = 289
	all.Wins = 32
	all.Losses = 35
	all.DamageDealt = 14403
	all.NoDamageDirectHitsReceived = 196
	all.MaxFrags = 5
	all.Shots = 1216
	all.ExplosionHitsReceived = 3
	all.TankingFactor = 0.33
	HowTheGodsKill.Statistics["all"] = all

	// check for good handling of clanless player
	c.Check(s.Wg.GetPlayerPersonalData([]uint32{525427444}), DeepEquals, []Player{HowTheGodsKill})

	// check for good handling of multiple players, all found in 1 request
	result := s.Wg.GetPlayerPersonalData([]uint32{525427444, 507197901})

	var HowTheGodsKill2, HowTheStoryEnds2 Player
	var playersFound []Player
	for _, v := range result {
		if v.Nickname == "HowTheGodsKill" {
			HowTheGodsKill2 = v
			playersFound = append(playersFound, HowTheGodsKill)
		}
		if v.Nickname == "HowTheStoryEnds" {
			HowTheStoryEnds2 = v
			playersFound = append(playersFound, HowTheStoryEnds)
		}

	}
	c.Check(HowTheGodsKill2, DeepEquals, HowTheGodsKill)
	c.Check(HowTheStoryEnds2, DeepEquals, HowTheStoryEnds)
	c.Check([]Player{HowTheStoryEnds2, HowTheGodsKill2}, DeepEquals, []Player{HowTheStoryEnds, HowTheGodsKill})
	c.Check(len(result), Equals, 2)
	c.Check(result, DeepEquals, playersFound)

}

func hasTank(TankId uint32, TankCollection []Tank) bool {
	for _, t := range TankCollection {
		if t.TankId == TankId {
			return true
		}
	}
	return false
}
func (s *WGSuite) TestGetPlayerTanks(c *C) {
	s.Wg.SetRegion("eu")
	Vehicle_1 := Tank{Wins: 55, Battles: 88, MarkOfMastery: 4, TankId: 1}
	Vehicle_11601 := Tank{Wins: 213, Battles: 408, MarkOfMastery: 4, TankId: 11601}
	Vehicle_3089 := Tank{Wins: 190, Battles: 361, MarkOfMastery: 4, TankId: 3089}
	Vehicle_11777 := Tank{Wins: 177, Battles: 284, MarkOfMastery: 4, TankId: 11777}

	// 1 vehicle
	c.Check(s.Wg.GetPlayerTanks(507197901, []uint32{1}), DeepEquals, []Tank{Vehicle_1})
	//multiple vehicles
	result := s.Wg.GetPlayerTanks(507197901, []uint32{11601, 3089, 11777})
	var compare []Tank
	for _, v := range result {
		switch v.TankId {
		case 11601:
			compare = append(compare, Vehicle_11601)
		case 3089:
			compare = append(compare, Vehicle_3089)
		case 11777:
			compare = append(compare, Vehicle_11777)
		}
	}
	c.Check(len(result), Equals, 3)
	c.Check(result, DeepEquals, compare)
	// unknown vehicle, should return an empty array
	c.Check(s.Wg.GetPlayerTanks(507197901, []uint32{2}), DeepEquals, []Tank{})
	// no vehicle ids given should return ALL vehicles
	all := s.Wg.GetPlayerTanks(507197901, []uint32{})
	// so check it returns enough vehicles
	c.Check(len(all), Equals, 210)
	//and check that it returned all the individual vehicles otherwise fail
	// lots of append because gvim otherwise chokes on it ;_; boo windows
	var pl = []uint32{11601, 1057, 3089, 10529, 4897, 1793, 11777, 11553, 10273, 12113, 3361, 2817, 2833, 3585, 5121, 5169, 7425, 3105, 1809, 10049, 57105, 8977, 3873, 1313}
	pl = append(pl, []uint32{11265, 7697, 4385, 9217, 1105, 7713, 2065, 4609, 5713, 3857, 11857, 54785, 11585, 4657, 6161, 16657, 1041, 6401, 6177, 545, 51713, 5969, 1889, 6433, 10497, 4097, 5409}...)
	pl = append(pl, []uint32{2625, 3633, 6465, 2897, 3601, 11025, 1537, 1121, 4913, 257, 513, 11009, 9553, 6721, 849, 2305, 801, 7761, 2577, 3329, 7201, 6673, 14097, 2369, 10241, 2881, 6417, 9041}...)
	pl = append(pl, []uint32{11041, 3073, 10817, 273, 6913, 8193, 785, 11521, 4673, 5185, 3153, 2129, 12561, 1, 11297, 6945, 15137, 16417, 1569, 16145, 14145, 1297, 6977, 3377, 11793, 1025, 2081}...)
	pl = append(pl, []uint32{15889, 1553, 55073, 10769, 3841, 5889, 529, 11089, 9505, 11345, 18433, 6481, 4689, 5457, 3137, 7249, 9249, 3921, 2321, 54353, 5649, 5905, 10785, 12289, 3121, 2561, 4129}...)
	pl = append(pl, []uint32{8785, 11281, 1073, 4369, 53585, 54609, 4929, 9793, 16641, 55569, 10753, 7185, 5665, 17953, 289, 13393, 1825, 51457, 4641, 53841, 5393, 5153, 4401, 10833, 15617, 52769}...)
	pl = append(pl, []uint32{18193, 57361, 5921, 12369, 5377, 7233, 1617, 8961, 54545, 4945, 55313, 16673, 6657, 8017, 4113, 81, 9761, 8257, 1089, 5953, 54801, 13345, 3409, 12545, 8273, 51745, 7169}...)
	pl = append(pl, []uint32{52481, 7969, 18177, 321, 6993, 1345, 10577, 64817, 55297, 2353, 577, 9473, 609, 3345, 593, 3617, 5201, 1601, 1329, 53537, 51553, 1361, 60689, 7745}...)
	for _, v := range pl {
		if !hasTank(v, all) {
			c.Fail()
		}
	}

}
