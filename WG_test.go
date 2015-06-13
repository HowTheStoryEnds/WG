package WG

import (
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

func (s *WGSuite) SetUpSuite(c *C) {

	//load json data
	var filenames = make(map[string]string)
	filenames["507197901"] = "./testdata/single_player.json"
	filenames["507197901,1"] = "./testdata/two_players_only_one_found.json"

	var content = make(map[string]string)

	for k, v := range filenames {
		var result, err = ioutil.ReadFile(v)
		if err == nil {
			content[k] = string(result)
		}

	}

	//setup HTTP mocking service
	httpmock.Activate()

	httpmock.RegisterResponder("GET", "https://api.mybiz.com/articles.json",
		httpmock.NewStringResponder(200, `[{"id": 1, "name": "My Great Article"}]`))

	httpmock.RegisterResponder("GET", "https://api.worldoftanks.eu/wot/account/list/?application_id=demo&search=howthestoryends",
		httpmock.NewStringResponder(200, `{
    "status": "ok",
    "meta": {
        "count": 1
    },
    "data": [
        {
            "nickname": "HowTheStoryEnds",
            "account_id": 507197901
        }
    ]
    }`))
	httpmock.RegisterResponder("GET", "https://api.worldoftanks.eu/wot/account/list/?application_id=demo&search=howthestoryends&type=exact",
		httpmock.NewStringResponder(200, `{
    "status": "ok",
    "meta": {
        "count": 1
    },
    "data": [
        {
            "nickname": "HowTheStoryEnds",
            "account_id": 507197901
        }
    ]
    }`))

	httpmock.RegisterResponder("GET", "https://api.worldoftanks.eu/wot/account/list/?application_id=demo&search=howthe",
		httpmock.NewStringResponder(200, `{
    "status": "ok",
    "meta": {
        "count": 4
    },
    "data": [
        {
            "nickname": "howtheblank",
            "account_id": 502301211
        },
        {
            "nickname": "HowTheGodsKill",
            "account_id": 525427444
        },
        {
            "nickname": "HowTheGuy",
            "account_id": 506219127
        },
        {
            "nickname": "HowTheStoryEnds",
            "account_id": 507197901
        }
    ]
    }`))

	httpmock.RegisterResponder("GET", "https://api.worldoftanks.eu/wot/account/info/?account_id=507197901&application_id=demo",
		httpmock.NewStringResponder(200, content["507197901"]))

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
	var players = []Player{{Nickname: "howtheblank", AccountId: 502301211, Region: s.Wg.region, Tanks: []Tank{}},
		{Nickname: "HowTheGodsKill", AccountId: 525427444, Region: s.Wg.region, Tanks: []Tank{}},
		{Nickname: "HowTheGuy", AccountId: 506219127, Region: s.Wg.region, Tanks: []Tank{}},
		{Nickname: "HowTheStoryEnds", AccountId: 507197901, Region: s.Wg.region, Tanks: []Tank{}},
	}
	data3 := s.Wg.SearchPlayersByName("howthe", false)
	c.Check(players, DeepEquals, data3)

	//no result found
	data4 := s.Wg.SearchPlayersByName("howthestt", false)
	c.Check([]Player{}, DeepEquals, data4)
}

func (s *WGSuite) TestGetPlayerPersonalData(c *C) {
	s.Wg.SetTransport("https")

	var players = []Player{{Nickname: "HowTheStoryEnds", AccountId: 507197901, Region: "eu"}}
	players[0].Frags = 0
	players[0].MaxXp = 2790
	players[0].Region = "eu"
	players[0].TreesCut = 21573
	players[0].MaxFrags = 11
	players[0].MaxDamage = 6017
	players[0].MaxXpTankId = 11265
	players[0].MaxFragsTankId = 2369
	players[0].MaxDamageTankId = 7425
	players[0].MaxDamageVehicle = 7425
	players[0].CreatedAt = 1355749511
	players[0].UpdatedAt = 1434133471
	players[0].Private = PlayerPrivate{}
	players[0].GlobalRating = 7728
	players[0].ClanId = 500010805
	players[0].Tanks = []Tank{}
	players[0].Statistics = map[string]PlayerStatistics{}
	players[0].LastBattleTime = 1434126614
	players[0].LogoutAt = 1434133469

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
	players[0].Statistics["clan"] = clan

	players[0].Statistics["regular_team"] = PlayerStatistics{}

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
	players[0].Statistics["company"] = company

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
	players[0].Statistics["stronghold_defense"] = strongholdDefense

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
	players[0].Statistics["stronghold_skirmish"] = strongholdSkirmish

	players[0].Statistics["historical"] = PlayerStatistics{}

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
	players[0].Statistics["team"] = team

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
	players[0].Statistics["all"] = all

	c.Check(s.Wg.GetPlayerPersonalData(507197901), DeepEquals, players)
}
