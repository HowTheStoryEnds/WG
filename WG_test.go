package WG

import(
    "testing"
    . "gopkg.in/check.v1"
    "github.com/jarcoal/httpmock"
    "net/http"
    "io/ioutil"
   
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




