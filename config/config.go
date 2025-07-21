package config

import "fmt"

type ServerConfig struct {
	GRPCAddress struct {
		Host string
		Port string
	}
	DbConnection struct {
		Host     string
		Port     string
		User     string
		Database string
		Password string
	}
	SecretKey string
}

func (cfg ServerConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", cfg.DbConnection.Host, cfg.DbConnection.Port, cfg.DbConnection.User, cfg.DbConnection.Database, cfg.DbConnection.Password)
}
