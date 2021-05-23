package config

// DBConfig holds the database configuration information
type DBConfig struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	DBname   string `json:"dbname"`
}

type APIConfig struct {
	BindIP   string `json:"bindÏP"`
	BindPort int    `json:"bindPort"`
}

type LogConfig struct {
	Level int `json:"level"`
}

// Config holds the main configuration of the application
type Config struct {
	DB  DBConfig  `json:"db"`
	API APIConfig `json:"API"`
	Log LogConfig `json:"log"`
}
