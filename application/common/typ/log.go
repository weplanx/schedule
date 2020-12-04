package typ

type Log struct {
	Identity string                 `json:"Identity"`
	Task     string                 `json:"Task"`
	Url      string                 `json:"Url"`
	Header   map[string]string      `json:"Header"`
	Body     map[string]interface{} `json:"Body"`
	Status   bool                   `json:"Status"`
	Response map[string]interface{} `json:"Response"`
	Time     int64                  `json:"Time"`
}
