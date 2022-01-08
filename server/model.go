package server

type SetuInfo struct {
	Id    int    `json:"pid"`
	Title string `json:"title"`
	Uid   int    `json:"uid"`
	Url   string `json:"url"`
	IsR18 bool   `json:"r18"`
}
