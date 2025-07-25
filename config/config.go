package config

import "fmt"

type ServerConfig struct {
	GRPCAddress  Address
	DbConnection Connection
	SecretKey    string `env:"SECRET_KEY"`
}
type Address struct {
	Host string `env:"GRPC_HOST"`
	Port string `env:"GRPC_PORT"`
}
type Connection struct {
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Database string `env:"DB_DATABASE"`
	Password string `env:"DB_PASSWORD"`
}

func (cfg ServerConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", cfg.DbConnection.Host, cfg.DbConnection.Port, cfg.DbConnection.User, cfg.DbConnection.Database, cfg.DbConnection.Password)
}
