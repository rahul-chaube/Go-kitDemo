package user

type Profile struct {
	Id      string  `json:"id"`
	Name    string  `json:"name"`
	Address Address `json:"address"`
}
type Address struct {
	Id       string `json:"id"`
	Location string `json:"location"`
}
