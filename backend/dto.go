package backend

type CommonAck[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

type ServerInfo struct {
	AreaId        int32    `json:"serverId"`
	RunState      int32    `json:"state"`
	Name          string   `json:"serverName"`
	OpenTimestamp int64    `json:"openTimestamp"`
	MergeAreaIds  []int32  `json:"mergedServers"`
	WhiteIps      []string `json:"whiteIps"`
}

type ActivationCodeItem struct {
	GoodsType int32 `json:"type"`
	GoodsId   int32 `json:"id"`
	Count     int32 `json:"size"`
}

type ActivationCodeAck struct {
	Code           int                  `json:"code"`
	ActivationCode string               `json:"activate_code"`
	Items          []ActivationCodeItem `json:"item_array"`
	GroupId        int32                `json:"group_id"`
}
