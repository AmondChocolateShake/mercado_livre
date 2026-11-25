package config

import (
	"net"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		Port:       getEnv("PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "mercado_1999"),
	}
}

func (c Config) DSN() string {
	return c.mysqlConfig(c.DBName).FormatDSN()
}

func (c Config) ServerDSN() string {
	return c.mysqlConfig("").FormatDSN()
}

func (c Config) mysqlConfig(dbName string) *mysql.Config {
	m := mysql.NewConfig()
	m.User = c.DBUser
	m.Passwd = c.DBPassword
	m.Net = "tcp"
	m.Addr = net.JoinHostPort(c.DBHost, c.DBPort)
	m.DBName = dbName
	m.ParseTime = true
	m.ClientFoundRows = true
	m.Params = map[string]string{"time_zone": "'+00:00'"}
	return m
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
