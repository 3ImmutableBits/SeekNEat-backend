package models

type Coords struct {
	Latitude  float64
	Longitude float64
}

type Meal struct {
	ID             uint   `gorm:"primaryKey"`
	Location       Coords `gorm:"serializer:json"`
	HostId         uint
	Timestamp      int64
	Price          string
	Clients        []*User `gorm:"many2many:user_meals;"`
	AvailableSpots uint
	Name           string
	Description    string
}
