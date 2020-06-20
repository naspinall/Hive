package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
)

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type Config struct {
	Port     int            `json:"port"`
	Env      string         `json:"env"`
	Pepper   string         `json:"pepper"`
	JWTKey   string         `json:"jwtKey"`
	Database PostgresConfig `json:"database"`
}

func (c PostgresConfig) Dialect() string {
	return "postgres"
}

func (c PostgresConfig) ConnectionInfo() string {
	if c.Password == "" {
		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Name)
	}

	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.Name)
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "hive",
		Name:     "hive",
	}
}

func (c Config) IsProd() bool {
	return c.Env == "production"
}

func DefaultConfig() Config {
	return Config{
		Port:     3001,
		Env:      "development",
		Pepper:   "salt-and-pepper-is-delicious",
		JWTKey:   "jwt-make-life-easy",
		Database: DefaultPostgresConfig(),
	}
}

// Load from environment variables
func LoadFromEnvironment() (pc PostgresConfig) {
	port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		pc.Port = 5432
	}
	pc.Host = os.Getenv("POSTGRES_HOST")
	pc.Port = port
	pc.User = os.Getenv("POSTGRES_USER")
	pc.Password = os.Getenv("POSTGRES_PASSWORD")
	pc.Name = os.Getenv("POSTGRES_DB")
	return
}

func LoadConfig() (c Config) {
	f, err := os.Open("./config.json")
	if err != nil {
		log.Println("No config file found, using default config")
		return DefaultConfig()
	}

	err = json.NewDecoder(f).Decode(&c)
	if err != nil {
		panic(err)
	}

	log.Println("Loaded config")
	return
}
