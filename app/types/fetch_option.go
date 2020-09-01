package types

type FetchOption struct {
	Url     string
	Headers map[string]string
	Body    interface{}
}
