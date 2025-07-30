package model

type Location struct {
	Latitude  float64
	Longitude float64
	Room      string
}

type Lesson struct {
	ID        string
	TeacherID string
	Subject   string
	StartTime int64
	EndTime   int64
	Location  Location
}
