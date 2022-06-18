package common

type Status struct {
	Key     string `msgpack:"key"`
	Running bool   `msgpack:"running"`
}
