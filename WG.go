package WG
import(
   "net/http"
   "net/url"
   "io/ioutil"
)

// we request our data from WG server, we allow for net/http
type WG struct {
   region string
   transport string
   method string
   apiKey string
}
func (w *WG) Init(region string,transport string,method string, ApiKey string){
     w.SetRegion(region)
     w.SetTransport(transport)
     w.SetMethod(method)
     w.SetApiKey(ApiKey)
}


func (w *WG) constructURL() string{
     //sane defaults
     if(w.transport != "http"){ w.transport = "https" }

     var regional = make(map[string]string) 
     regional["eu"] = "api.worldoftanks.eu/wot/"
     regional["na"] = "api.worldoftanks.com/wot/"
     regional["ru"] = "api.worldoftanks.ru/wot/"
                
      
     return w.transport + "://" + regional[w.region]
    
}

func (w *WG) retrieveData(action string, parameters map[string]string) string{

        baseUrl, err := url.Parse("http://google.com/search")
	if err != nil {
                //panic!
		
	}

	params := url.Values{}
	params.Add("pass%word", "key%20word")

	baseUrl.RawQuery = params.Encode()
	

    var val,_ = http.Get("https://api.mybiz.com/articles.json")
    var res,_ = ioutil.ReadAll(val.Body)
    val.Body.Close()

    return string(res)
}

func (w *WG) SetRegion(val string){
     w.region = val
}
func (w *WG) SetTransport(val string){
     w.transport = val
}
func (w *WG) SetMethod(val string){
     w.method = val
}
func (w *WG) SetApiKey(val string){
     w.apiKey = val
}

type Tank struct {
    Name string
    InGarage bool
    AvgDmg uint32
    AvgSpotting uint32
    AvgExp uint32
}


type Player struct {
    id uint32
    nickname string
    region string
    Tanks []Tank
}


func (w* WG) GetPlayer(ApiKey string, player uint32) Player {
   return Player{}
}


