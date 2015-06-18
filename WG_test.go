package WG

import (
	"fmt"
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

	rd.Uri = []string{"https://api.worldoftanks.eu/wot/account/list/?application_id=demo&search=howthestt"}
	rd.Content, _ = ioutil.ReadFile("./testdata/account/list/howthestt_name.json")
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
	rd.Uri = []string{"https://api.worldoftanks.eu/wot/account/tanks/?account_id=507197901&application_id=demo&tank_id=1",
		"https://api.worldoftanks.eu/wot/account/tanks/?account_id=507197901&application_id=demo&tank_id=13%2C1"}
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
	// complete vehicle list multiple players
	rd.Uri = []string{"https://api.worldoftanks.eu/wot/account/tanks/?account_id=507197901%2C515080611&application_id=demo"}
	rd.Content, _ = ioutil.ReadFile("./testdata/account/tanks/complete_vehicle_list_multiple_players.json")
	res = append(res, rd)
	// multiple players, multiple vehicles
	rd.Uri = []string{"https://api.worldoftanks.eu/wot/account/tanks/?account_id=507197901%2C515080611&application_id=demo&tank_id=81%2C3329%2C321"}
	rd.Content, _ = ioutil.ReadFile("./testdata/account/tanks/multiple_players_multiple_vehicles.json")
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
func (s *WGSuite) TestconstructURL(c *C) {
	s.Wg.SetTransport("http")
	s.Wg.SetRegion("eu")
	// one key, one parameter
	var params = map[string]string{"one": "1",
		"two": "twee",
	}
	c.Check(s.Wg.constructURL("account/list", params), Equals, "http://api.worldoftanks.eu/wot/account/list/?one=1&two=twee")

	// special characters need to be escaped
	params = map[string]string{"1with_slashes": "co//ol", "2email": "harald.brinkhof@gmail.com"}
	c.Check(s.Wg.constructURL("account/help", params), Equals, "http://api.worldoftanks.eu/wot/account/help/?1with_slashes=co%2F%2Fol&2email=harald.brinkhof%40gmail.com")
}

func (s *WGSuite) TestretrieveData(c *C) {
	s.Wg.SetTransport("https")
	s.Wg.SetRegion("eu")
	ret, err := s.Wg.retrieveData("account/list", map[string]string{"application_id": "demo", "type": "exact", "search": "howthestoryends"})

	data, _ := ioutil.ReadFile("./testdata/account/list/howthestoryends_name.json")
	c.Check(err, Equals, nil)
	c.Check(ret, DeepEquals, data)

	// unretrievable data
	ret, err = s.Wg.retrieveData("@/@", map[string]string{})
	c.Check(err, NotNil)
	c.Check(ret, DeepEquals, []byte{})

}

func (s *WGSuite) TestSearchPlayersByName(c *C) {
	s.Wg.SetTransport("https")
	s.Wg.SetRegion("eu")
	// 'startswith' search yielding 1 result
	var player = []Player{{Nickname: "HowTheStoryEnds", AccountId: 507197901}}
	data, err := s.Wg.SearchPlayersByName("howthestoryends", false)
	c.Assert(err, Equals, nil)
	c.Check(data, DeepEquals, player)

	// exact search yielding 1 result
	data2, _ := s.Wg.SearchPlayersByName("howthestoryends", true)
	c.Check(player, DeepEquals, data2)

	// search that can yield multiple results
	var players = []Player{}
	data3, _ := s.Wg.SearchPlayersByName("howthe", false)
	for _, p := range data3 {
		switch p.Nickname {
		case "howtheblank":
			players = append(players, Player{Nickname: "howtheblank", AccountId: 502301211})
		case "HowTheGodsKill":
			players = append(players, Player{Nickname: "HowTheGodsKill", AccountId: 525427444})
		case "HowTheGuy":
			players = append(players, Player{Nickname: "HowTheGuy", AccountId: 506219127})
		case "HowTheStoryEnds":
			players = append(players, Player{Nickname: "HowTheStoryEnds", AccountId: 507197901})
		}
	}
	c.Check(players, DeepEquals, data3)

	//no result found
	data4, _ := s.Wg.SearchPlayersByName("howthestt", false)
	c.Check([]Player{}, DeepEquals, data4)
}

func (s *WGSuite) TestGetPlayerPersonalData(c *C) {
	s.Wg.SetTransport("https")
	s.Wg.SetRegion("eu")

	// player in a clan
	var HowTheStoryEnds = Player{Nickname: "HowTheStoryEnds", AccountId: 507197901}
	HowTheStoryEnds.ClientLanguage = "en"
	HowTheStoryEnds.CreatedAt = 1355749511
	HowTheStoryEnds.UpdatedAt = 1434133471
	HowTheStoryEnds.Private = PlayerPrivate{}
	HowTheStoryEnds.GlobalRating = 7728
	HowTheStoryEnds.ClanId = 500010805
	HowTheStoryEnds.LastBattleTime = 1434126614
	HowTheStoryEnds.LogoutAt = 1434133469

	ps := PlayerStatistics{}
	ps.TreesCut = 21573

	ps.Clan.Spotted = 36
	ps.Clan.AvgDamageAssistedTrack = 45.57
	ps.Clan.AvgDamageBlocked = 219.55
	ps.Clan.DirectHitsReceived = 121
	ps.Clan.ExplosionHits = 15
	ps.Clan.PiercingsReceived = 99
	ps.Clan.Piercings = 132
	ps.Clan.Xp = 37281
	ps.Clan.SurvivedBattles = 33
	ps.Clan.DroppedCapturePoints = 2
	ps.Clan.HitsPercents = 57
	ps.Clan.Draws = 2
	ps.Clan.Battles = 55
	ps.Clan.DamageReceived = 33255
	ps.Clan.AvgDamageAssisted = 184.95
	ps.Clan.Frags = 36
	ps.Clan.AvgDamageAssistedRadio = 139.39
	ps.Clan.CapturePoints = 26
	ps.Clan.Hits = 195
	ps.Clan.BattleAvgXp = 678
	ps.Clan.Wins = 46
	ps.Clan.Losses = 7
	ps.Clan.DamageDealt = 35846
	ps.Clan.NoDamageDirectHitsReceived = 22
	ps.Clan.Shots = 343
	ps.Clan.ExplosionHitsReceived = 1
	ps.Clan.TankingFactor = 0.35

	ps.Company.Spotted = 190
	ps.Company.AvgDamageAssistedTrack = 6.25
	ps.Company.AvgDamageBlocked = 0
	ps.Company.DirectHitsReceived = 52
	ps.Company.ExplosionHits = 0
	ps.Company.PiercingsReceived = 48
	ps.Company.Piercings = 58
	ps.Company.Xp = 82190
	ps.Company.SurvivedBattles = 60
	ps.Company.DroppedCapturePoints = 156
	ps.Company.HitsPercents = 61
	ps.Company.Draws = 3
	ps.Company.Battles = 152
	ps.Company.DamageReceived = 74793
	ps.Company.AvgDamageAssisted = 165.58
	ps.Company.Frags = 98
	ps.Company.AvgDamageAssistedRadio = 159.33
	ps.Company.CapturePoints = 229
	ps.Company.Hits = 502
	ps.Company.BattleAvgXp = 541
	ps.Company.Wins = 97
	ps.Company.Losses = 52
	ps.Company.DamageDealt = 81272
	ps.Company.NoDamageDirectHitsReceived = 4
	ps.Company.Shots = 819
	ps.Company.ExplosionHitsReceived = 0
	ps.Company.TankingFactor = 0.0

	ps.StrongholdDefense.Spotted = 41
	ps.StrongholdDefense.MaxFragsTankId = 7169
	ps.StrongholdDefense.MaxXp = 1062
	ps.StrongholdDefense.DirectHitsReceived = 177
	ps.StrongholdDefense.ExplosionHits = 0
	ps.StrongholdDefense.PiercingsReceived = 121
	ps.StrongholdDefense.Piercings = 121
	ps.StrongholdDefense.Xp = 21559
	ps.StrongholdDefense.SurvivedBattles = 17
	ps.StrongholdDefense.DroppedCapturePoints = 31
	ps.StrongholdDefense.HitsPercents = 74
	ps.StrongholdDefense.Draws = 0
	ps.StrongholdDefense.MaxXpTankId = 10785
	ps.StrongholdDefense.Battles = 33
	ps.StrongholdDefense.DamageReceived = 47316
	ps.StrongholdDefense.Frags = 17
	ps.StrongholdDefense.CapturePoints = 38
	ps.StrongholdDefense.MaxDamageTankId = 12369
	ps.StrongholdDefense.MaxDamage = 3900
	ps.StrongholdDefense.Hits = 179
	ps.StrongholdDefense.BattleAvgXp = 653
	ps.StrongholdDefense.Wins = 27
	ps.StrongholdDefense.Losses = 6
	ps.StrongholdDefense.DamageDealt = 45584
	ps.StrongholdDefense.NoDamageDirectHitsReceived = 56
	ps.StrongholdDefense.MaxFrags = 2
	ps.StrongholdDefense.Shots = 242
	ps.StrongholdDefense.ExplosionHitsReceived = 26
	ps.StrongholdDefense.TankingFactor = 0.33

	ps.StrongholdSkirmish.Spotted = 1350
	ps.StrongholdSkirmish.MaxFragsTankId = 4913
	ps.StrongholdSkirmish.MaxXp = 1980
	ps.StrongholdSkirmish.DirectHitsReceived = 7982
	ps.StrongholdSkirmish.ExplosionHits = 224
	ps.StrongholdSkirmish.PiercingsReceived = 5722
	ps.StrongholdSkirmish.Piercings = 8514
	ps.StrongholdSkirmish.Xp = 1899689
	ps.StrongholdSkirmish.SurvivedBattles = 1413
	ps.StrongholdSkirmish.DroppedCapturePoints = 1522
	ps.StrongholdSkirmish.HitsPercents = 69
	ps.StrongholdSkirmish.Draws = 9
	ps.StrongholdSkirmish.MaxXpTankId = 4385
	ps.StrongholdSkirmish.Battles = 2158
	ps.StrongholdSkirmish.DamageReceived = 1750773
	ps.StrongholdSkirmish.Frags = 1399
	ps.StrongholdSkirmish.CapturePoints = 2134
	ps.StrongholdSkirmish.MaxDamageTankId = 7169
	ps.StrongholdSkirmish.MaxDamage = 4144
	ps.StrongholdSkirmish.Hits = 11663
	ps.StrongholdSkirmish.BattleAvgXp = 880
	ps.StrongholdSkirmish.Wins = 1947
	ps.StrongholdSkirmish.Losses = 202
	ps.StrongholdSkirmish.DamageDealt = 2120855
	ps.StrongholdSkirmish.NoDamageDirectHitsReceived = 2260
	ps.StrongholdSkirmish.MaxFrags = 5
	ps.StrongholdSkirmish.Shots = 16898
	ps.StrongholdSkirmish.ExplosionHitsReceived = 146
	ps.StrongholdSkirmish.TankingFactor = 0.45

	ps.Team.Spotted = 615
	ps.Team.AvgDamageAssistedTrack = 66.12
	ps.Team.MaxXp = 1848
	ps.Team.AvgDamageBlocked = 224.65
	ps.Team.DirectHitsReceived = 2552
	ps.Team.ExplosionHits = 1
	ps.Team.PiercingsReceived = 1851
	ps.Team.Piercings = 2259
	ps.Team.MaxDamageTankId = 5377
	ps.Team.Xp = 474872
	ps.Team.SurvivedBattles = 278
	ps.Team.DroppedCapturePoints = 1301
	ps.Team.HitsPercents = 71
	ps.Team.Draws = 11
	ps.Team.MaxXpTankId = 5377
	ps.Team.Battles = 651
	ps.Team.DamageReceived = 392013
	ps.Team.AvgDamageAssisted = 172.07
	ps.Team.MaxFragsTankId = 4385
	ps.Team.Frags = 378
	ps.Team.AvgDamageAssistedRadio = 105.95
	ps.Team.CapturePoints = 1958
	ps.Team.MaxDamage = 3526
	ps.Team.Hits = 3268
	ps.Team.BattleAvgXp = 729
	ps.Team.Wins = 491
	ps.Team.Losses = 149
	ps.Team.DamageDealt = 517046
	ps.Team.NoDamageDirectHitsReceived = 701
	ps.Team.MaxFrags = 5
	ps.Team.Shots = 4603
	ps.Team.ExplosionHitsReceived = 2
	ps.Team.TankingFactor = 0.37

	ps.All.Spotted = 24553
	ps.All.AvgDamageAssistedTrack = 66.92
	ps.All.MaxXp = 2790
	ps.All.AvgDamageBlocked = 222.86
	ps.All.DirectHitsReceived = 50646
	ps.All.ExplosionHits = 4000
	ps.All.PiercingsReceived = 37905
	ps.All.Piercings = 79506
	ps.All.MaxDamageTankId = 7425
	ps.All.Xp = 12796688
	ps.All.SurvivedBattles = 6757
	ps.All.DroppedCapturePoints = 23554
	ps.All.HitsPercents = 59
	ps.All.Draws = 226
	ps.All.MaxXpTankId = 11265
	ps.All.Battles = 19468
	ps.All.DamageReceived = 9471075
	ps.All.AvgDamageAssisted = 307.41
	ps.All.MaxFragsTankId = 2369
	ps.All.Frags = 24133
	ps.All.AvgDamageAssistedRadio = 240.49
	ps.All.CapturePoints = 19541
	ps.All.MaxDamage = 6017
	ps.All.Hits = 159851
	ps.All.BattleAvgXp = 657
	ps.All.Wins = 11445
	ps.All.Losses = 7797
	ps.All.DamageDealt = 16903164
	ps.All.NoDamageDirectHitsReceived = 12741
	ps.All.MaxFrags = 11
	ps.All.Shots = 270136
	ps.All.ExplosionHitsReceived = 1564
	ps.All.TankingFactor = 0.35
	HowTheStoryEnds.Statistics = ps

	retrieved, err := s.Wg.GetPlayerPersonalData([]uint32{507197901})
	if err != nil {
		fmt.Println(err.Error())
		c.Fail()
	}
	c.Check(retrieved, DeepEquals, []Player{HowTheStoryEnds})
	retrieved, err = s.Wg.GetPlayerPersonalData([]uint32{507197901, 1})
	// id 1 does not exist so what should be returned is 1 player record for 507197901 and 1 empty player record
	compare := []Player{}
	for _, v := range retrieved {
		if v.Nickname == "HowTheStoryEnds" {
			compare = append(compare, HowTheStoryEnds)
		} else {
			compare = append(compare, Player{})
		}

	}
	c.Check(retrieved, DeepEquals, compare)
	if err != nil {
		fmt.Println(err.Error())
		c.Fail()
	}
	// clanless player
	var HowTheGodsKill = Player{Nickname: "HowTheGodsKill", AccountId: 525427444}

	HowTheGodsKill.CreatedAt = 1420581930
	HowTheGodsKill.UpdatedAt = 1433267697
	HowTheGodsKill.Private = PlayerPrivate{}
	HowTheGodsKill.GlobalRating = 1263
	HowTheGodsKill.LastBattleTime = 1423093994
	HowTheGodsKill.LogoutAt = 1433267693
	HowTheGodsKill.ClientLanguage = "pl"
	HowTheGodsKill.Statistics.TreesCut = 42

	HowTheGodsKill.Statistics.All.Spotted = 53
	HowTheGodsKill.Statistics.All.AvgDamageAssistedTrack = 4.34
	HowTheGodsKill.Statistics.All.MaxXp = 1407
	HowTheGodsKill.Statistics.All.AvgDamageBlocked = 51.58
	HowTheGodsKill.Statistics.All.DirectHitsReceived = 552
	HowTheGodsKill.Statistics.All.ExplosionHits = 0
	HowTheGodsKill.Statistics.All.PiercingsReceived = 356
	HowTheGodsKill.Statistics.All.Piercings = 462
	HowTheGodsKill.Statistics.All.MaxDamageTankId = 4945
	HowTheGodsKill.Statistics.All.Xp = 19371
	HowTheGodsKill.Statistics.All.SurvivedBattles = 14
	HowTheGodsKill.Statistics.All.DroppedCapturePoints = 43
	HowTheGodsKill.Statistics.All.HitsPercents = 57
	HowTheGodsKill.Statistics.All.Draws = 0
	HowTheGodsKill.Statistics.All.MaxXpTankId = 4945
	HowTheGodsKill.Statistics.All.Battles = 67
	HowTheGodsKill.Statistics.All.DamageReceived = 10986
	HowTheGodsKill.Statistics.All.AvgDamageAssisted = 25.39
	HowTheGodsKill.Statistics.All.MaxFragsTankId = 4945
	HowTheGodsKill.Statistics.All.Frags = 70
	HowTheGodsKill.Statistics.All.AvgDamageAssistedRadio = 21.04
	HowTheGodsKill.Statistics.All.CapturePoints = 83
	HowTheGodsKill.Statistics.All.MaxDamage = 1354
	HowTheGodsKill.Statistics.All.Hits = 691
	HowTheGodsKill.Statistics.All.BattleAvgXp = 289
	HowTheGodsKill.Statistics.All.Wins = 32
	HowTheGodsKill.Statistics.All.Losses = 35
	HowTheGodsKill.Statistics.All.DamageDealt = 14403
	HowTheGodsKill.Statistics.All.NoDamageDirectHitsReceived = 196
	HowTheGodsKill.Statistics.All.MaxFrags = 5
	HowTheGodsKill.Statistics.All.Shots = 1216
	HowTheGodsKill.Statistics.All.ExplosionHitsReceived = 3
	HowTheGodsKill.Statistics.All.TankingFactor = 0.33

	// check for good handling of clanless player
	retrieved, err = s.Wg.GetPlayerPersonalData([]uint32{525427444})
	if err != nil {
		fmt.Println(err.Error())
		c.Fail()
	}
	c.Check(retrieved, DeepEquals, []Player{HowTheGodsKill})

	// check for good handling of multiple players, all found in 1 request
	retrieved, err = s.Wg.GetPlayerPersonalData([]uint32{525427444, 507197901})
	if err != nil {
		fmt.Println(err.Error())
		c.Fail()
	}
	var HowTheGodsKill2, HowTheStoryEnds2 Player
	var playersFound []Player
	for _, v := range retrieved {
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
	c.Check(len(retrieved), Equals, 2)
	c.Check(retrieved, DeepEquals, playersFound)

}

func hasTank(TankId uint32, Owner uint32, TankCollection []Vehicle) bool {
	for _, t := range TankCollection {
		if t.TankId == TankId && t.Owner == Owner {
			return true
		}
	}
	return false
}
func (s *WGSuite) TestGetPlayerTanks(c *C) {
	s.Wg.SetRegion("eu")

	Vehicle_1 := Vehicle{Statistics: VehicleStatistics{Wins: 55, Battles: 88}, MarkOfMastery: 4, TankId: 1, Owner: 507197901}
	Vehicle_11601 := Vehicle{Statistics: VehicleStatistics{Wins: 213, Battles: 408}, MarkOfMastery: 4, TankId: 11601, Owner: 507197901}
	Vehicle_3089 := Vehicle{Statistics: VehicleStatistics{Wins: 190, Battles: 361}, MarkOfMastery: 4, TankId: 3089, Owner: 507197901}
	Vehicle_11777 := Vehicle{Statistics: VehicleStatistics{Wins: 177, Battles: 284}, MarkOfMastery: 4, TankId: 11777, Owner: 507197901}

	result, err := s.Wg.GetPlayerTanks([]uint32{507197901}, []uint32{1})
	if err != nil {
		fmt.Println(err.Error())
		c.Fail()
	}
	// 1 vehicle
	c.Check(result, DeepEquals, []Vehicle{Vehicle_1})
	//multiple but not all found
	result, err = s.Wg.GetPlayerTanks([]uint32{507197901}, []uint32{13, 1})
	if err != nil {
		fmt.Println(err.Error())
		c.Fail()
	}
	// 1 vehicle
	c.Check(result, DeepEquals, []Vehicle{Vehicle_1})

	// all found
	result, err = s.Wg.GetPlayerTanks([]uint32{507197901}, []uint32{11601, 3089, 11777})
	var compare []Vehicle
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
	result, err = s.Wg.GetPlayerTanks([]uint32{507197901}, []uint32{2})
	if err != nil {
		fmt.Println(err.Error())
		c.Fail()
	}
	c.Check(result, DeepEquals, []Vehicle{})

	// no vehicle ids given should return ALL vehicles
	result, err = s.Wg.GetPlayerTanks([]uint32{507197901}, []uint32{})
	if err != nil {
		fmt.Println(err.Error())
		c.Fail()
	}
	// so check it returns enough vehicles
	var pl = []uint32{11601, 1057, 3089, 10529, 4897, 1793, 11777, 11553, 10273, 12113, 3361, 2817, 2833, 3585, 5121, 5169, 7425, 3105, 1809, 10049, 57105, 8977, 3873, 1313}
	pl = append(pl, []uint32{11265, 7697, 4385, 9217, 1105, 7713, 2065, 4609, 5713, 3857, 11857, 54785, 11585, 4657, 6161, 16657, 1041, 6401, 6177, 545, 51713, 5969, 1889, 6433, 10497, 4097, 5409}...)
	pl = append(pl, []uint32{2625, 3633, 6465, 2897, 3601, 11025, 1537, 1121, 4913, 257, 513, 11009, 9553, 6721, 849, 2305, 801, 7761, 2577, 3329, 7201, 6673, 14097, 2369, 10241, 2881, 6417, 9041}...)
	pl = append(pl, []uint32{11041, 3073, 10817, 273, 6913, 8193, 785, 11521, 4673, 5185, 3153, 2129, 12561, 1, 11297, 6945, 15137, 16417, 1569, 16145, 14145, 1297, 6977, 3377, 11793, 1025, 2081}...)
	pl = append(pl, []uint32{15889, 1553, 55073, 10769, 3841, 5889, 529, 11089, 9505, 11345, 18433, 6481, 4689, 5457, 3137, 7249, 9249, 3921, 2321, 54353, 5649, 5905, 10785, 12289, 3121, 2561, 4129}...)
	pl = append(pl, []uint32{8785, 11281, 1073, 4369, 53585, 54609, 4929, 9793, 16641, 55569, 10753, 7185, 5665, 17953, 289, 13393, 1825, 51457, 4641, 53841, 5393, 5153, 4401, 10833, 15617, 52769}...)
	pl = append(pl, []uint32{18193, 57361, 5921, 12369, 5377, 7233, 1617, 8961, 54545, 4945, 55313, 16673, 6657, 8017, 4113, 81, 9761, 8257, 1089, 5953, 54801, 13345, 3409, 12545, 8273, 51745, 7169}...)
	pl = append(pl, []uint32{52481, 7969, 18177, 321, 6993, 1345, 10577, 64817, 55297, 2353, 577, 9473, 609, 3345, 593, 3617, 5201, 1601, 1329, 53537, 51553, 1361, 60689, 7745}...)
	c.Check(len(result), Equals, len(pl))
	//and check that it returned all the individual vehicles otherwise fail
	// lots of append because gvim otherwise chokes on it ;_; boo windows

	for _, v := range pl {
		if !hasTank(v, 507197901, result) {
			c.Fail()
		}
	}
	// all vehicles but multiple players will return a map containing tank array
	var pl2 = []uint32{81, 3089, 2065, 545, 3329}
	result, err = s.Wg.GetPlayerTanks([]uint32{507197901, 515080611}, []uint32{})
	if err != nil {
		fmt.Println(err.Error())
		c.Fail()
	}
	// check that all the tanks of their respective owners are returned
	for _, v := range pl {
		if !hasTank(v, 507197901, result) {
			fmt.Println("507197901: missing tanks, pl array")
			c.Fail()
		}
	}
	for _, v := range pl2 {
		if !hasTank(v, 515080611, result) {
			fmt.Println("515080611: missing tanks, pl2 array")
			c.Fail()
		}
	}
	// so check it returns enough vehicles
	c.Check(result, HasLen, len(pl)+len(pl2))

	// not every player has every tank, both have 81 and 3329 bigger account has 321
	result, err = s.Wg.GetPlayerTanks([]uint32{507197901, 515080611}, []uint32{81, 3329, 321})
	if err != nil {
		fmt.Println(err.Error())
		c.Fail()
	}
	c.Check(hasTank(81, 507197901, result) && hasTank(81, 515080611, result), Equals, true)
	c.Check(hasTank(3329, 507197901, result) && hasTank(3329, 515080611, result), Equals, true)
	c.Check(hasTank(321, 507197901, result) && !hasTank(321, 515080611, result), Equals, true)
	c.Check(result, HasLen, 5)

}
