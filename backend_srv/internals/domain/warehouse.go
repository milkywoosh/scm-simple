package domain

// implement location repo

type Warehouse struct {
	Description    string `json:"description"`
	LocationCode   string `json:"location_code"`
	ParentLocation string `json:"parent_location"`
	Point          any    `json:"point"`
}


type WarehouseInfo interface {
	
}