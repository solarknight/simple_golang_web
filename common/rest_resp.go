package common

type RestResponse struct {
	State int         `json:"state"`
	Msg   string      `json:"msg,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}
