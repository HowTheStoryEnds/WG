package WG

import(
    "testing"
    . "gopkg.in/check.v1"
    "github.com/jarcoal/httpmock"
    "encoding/json"
    
)




// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }


type WGSuite struct{
     Wg WG
}
var _ = Suite(&WGSuite{})


func (s *WGSuite) SetUpSuite(c *C){
    //setup HTTP mocking service    
    httpmock.Activate()
    

    httpmock.RegisterResponder("GET", "https://api.mybiz.com/articles.json",
    httpmock.NewStringResponder(200, `[{"id": 1, "name": "My Great Article"}]`))

    s.Wg.Init("na","https","get","apikey")
}

func (s *WGSuite) TearDownSuite(c *C){
    //shutdown HTTP mocking
    httpmock.DeactivateAndReset()
}


// tests
func (s *WGSuite) TestConstructURL(c *C) {
    
    s.Wg.SetRegion("na")
    c.Check(s.Wg.constructURL(), Equals,"https://api.worldoftanks.com/wot/")    
    s.Wg.SetRegion("eu")
    c.Check(s.Wg.constructURL(), Equals,"https://api.worldoftanks.eu/wot/")

    s.Wg.SetTransport("http")
    c.Check(s.Wg.constructURL(), Equals,"http://api.worldoftanks.eu/wot/")

    
    
    
}

func (s *WGSuite) TestaddGetParams(c *C) {

    // one key, one parameter
    var params = map[string][]string{"one": {"1"},
                                "two": {"twee"},
                 }
    c.Check(s.Wg.addGetParams("http://www.example.com",params),Equals,"http://www.example.com?one=1&two=twee")

    // multiple parameters to the same key
    params = map[string][]string{"dienaren": {"Harald","Babs"},
                                 "dieren": {"Rickey","Tiger","Storm","Mazzeltje"},
             }
    c.Check(s.Wg.addGetParams("http://www.example.com",params),
                               Equals,
                              "http://www.example.com?dienaren=Harald&dienaren=Babs&dieren=Rickey&dieren=Tiger&dieren=Storm&dieren=Mazzeltje")

    // special characters need to be escaped
    params = map[string][]string{"1with_slashes": {"co//ol"},"2email": {"harald.brinkhof@gmail.com"}}
    c.Check(s.Wg.addGetParams("http://www.example.com",params),Equals,"http://www.example.com?1with_slashes=co%2F%2Fol&2email=harald.brinkhof%40gmail.com")

}

func (s *WGSuite) TestSearchPlayersByName(c *C){

   var returnValue map[string]interface{}
   testData := []byte(`{
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
   }`)
   err := json.Unmarshal(testData,&returnValue)
   if(err != nil){ 
          c.Error("testData Unmarshal: " + err.Error()) 
          
   }

   var result map[string]interface{}
   json.Unmarshal([]byte(s.Wg.SearchPlayersByName("howthestoryends",false)),&result)
  
   c.Check(result, DeepEquals, returnValue)


}
