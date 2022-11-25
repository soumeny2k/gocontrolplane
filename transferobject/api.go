package transferobject

type Api struct {
	Name              string    `json:"name"`
	TeamId            uint      `json:"teamId"`
	RateLimit         uint      `json:"rateLimit"`
	ConnectionTimeout uint      `json:"connectionTimeout"`
	Protocol          string    `json:"protocol" default:"http"`
	Type              string    `json:"type" default:"REST"`
	Backends          []Backend `json:"backends"`
}

type Backend struct {
	Url    string `json:"url"`
	Weight uint   `json:"weight"`
}
