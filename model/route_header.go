package model

type RouteHeader struct {
	ID      uint   `json:"id"`
	RouteId uint   `json:"route_id"`
	Name    string `json:"name"`
	Value   string `json:"value"`
}
