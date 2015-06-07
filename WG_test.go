package github.com/howthestoryends/WG

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
    s.Wg.Init("na","https","get","apikey")
}

func (s *WGSuite) TestConstructURL(c *C) {
    s.Wg.SetRegion("na")
    c.Check(s.Wg.constructURL(), Equals,"https://api.worldoftanks.com/wot/")    
    s.Wg.SetRegion("eu")
    c.Check(s.Wg.constructURL(), Equals,"https://api.worldoftanks.eu/wot/")
    s.Wg.SetTransport("http")
    c.Check(s.Wg.constructURL(), Equals,"http://api.worldoftanks.eu/wot/")
    c.Check(42, Equals, 42)
    c.Check(42, Equals, 42)
    
}


func (s *WGSuite) TestHelloWorld2(c *C) {
    
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()

    httpmock.RegisterResponder("GET", "https://api.mybiz.com/articles.json",
        httpmock.NewStringResponder(200, `[{"id": 1, "name": "My Great Article"}]`))
    //var hc = http.Client{}
    var val,_ = http.Get("https://api.mybiz.com/articles.json")
    var _,_ = ioutil.ReadAll(val.Body)
    val.Body.Close()
    c.Check(42, Equals, 42)
   
    
}

