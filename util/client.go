package util

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	// Apikey is used in netvigil tic of `config.toml`
	Apikey string `json:"apikey"`
	// Admin is a flag to indicate whether the user is an administrator
	Admin bool `json:"admin"`
}

func CreateUser(username, password, apikey string, admin bool) error {
	_, err := DB.Exec("INSERT INTO users (username, password, apikey, admin) VALUES (?, ?, ?, ?)", username, password, apikey, admin)
	return err
}

func GetUser(username string) (*User, error) {
	var u User
	err := DB.QueryRow("SELECT username, password, apikey, admin FROM users WHERE username = ?", username).Scan(&u.Username, &u.Password, &u.Apikey, &u.Admin)
	return &u, err
}

func DeleteUser(username string) error {
	_, err := DB.Exec("DELETE FROM users WHERE username = ?", username)
	return err
}

func init() {
	DB.Exec("CREATE TABLE IF NOT EXISTS users (username TEXT PRIMARY KEY, password TEXT, apikey TEXT, admin INTEGER)")
}
