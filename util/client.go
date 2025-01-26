package util

import (
	"log"

	uuid "github.com/nu7hatch/gouuid"
)

type Client struct {
	Name string `json:"name"`
	// Apikey is used in netvigil tic of `config.toml`
	Apikey string `json:"apikey"`
}

func GetClients() []Client {
	rows, err := DB.Query("SELECT name, apikey FROM clients")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	clients := []Client{}
	for rows.Next() {
		client := Client{}
		err := rows.Scan(&client.Name, &client.Apikey)
		if err != nil {
			log.Fatal(err)
		}
		clients = append(clients, client)
	}
	return clients
}

func CreateClient(name string) error {
	u, err := uuid.NewV4()
	if err != nil {
		return err
	}
	_, err = DB.Exec("INSERT INTO clients (name, apikey) VALUES (?, ?)", name, u.String())
	return err
}

func DeleteClient(apikey string) error {
	_, err := DB.Exec("DELETE FROM clients WHERE apikey = ?", apikey)
	return err
}

func init() {
	DB.Exec("CREATE TABLE IF NOT EXISTS clients (name TEXT, apikey TEXT)")
}
