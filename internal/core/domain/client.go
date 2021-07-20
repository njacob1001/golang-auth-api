package domain

type Client struct {
	ID        string `json: "id"`
	Name      string `json:"name"`
	LastName  string `json:"last_name"`
	Birthday  string `json:"birthday"`
	Email     string `json:"email"`
	City      string `json:"city"`
	Address   string `json:"address"`
	Cellphone string `json:"cellphone"`
}