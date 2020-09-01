package types

type LoggingPush struct {
	Identity string
	HasError bool
	Message  interface{}
}
