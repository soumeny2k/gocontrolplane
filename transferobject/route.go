package transferobject

type Route struct {
	Name              string   `json:"name"`
	Path              string   `json:"path"`
	PathPattern       string   `json:"pathPattern"`
	Method            string   `json:"method"`
	Retries           uint     `json:"retries"`
	RateLimit         uint     `json:"rateLimit"`
	ConnectionTimeout uint     `json:"connectionTimeout"`
	CacheEnabled      bool     `json:"cacheEnabled"`
	Headers           []Header `json:"headers"`
}

type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
