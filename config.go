package pgxhelper

import (
	"fmt"
	"os"
)

type DBConfig interface {
	MakeDSN() string
}

var _ DBConfig = &Config{}

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SslMode  string
	Timezone string
}

const (
	Host     = "DB_HOST"
	Port     = "DB_PORT"
	User     = "DB_USER"
	Password = "DB_PASSWORD"
	Name     = "DB_NAME"
	Sslmode  = "DB_SSL_MODE"
	Timezone = "DB_TIMEZONE"
)

func message(variable string) error {
	return fmt.Errorf("env variable %q is not set and is required", variable)
}
func getvar(env string) (string, error) {
	val, ok := os.LookupEnv(env)
	if !ok {
		return "", message(env)
	}
	return val, nil
}
func GetConfigFromEnv() (DBConfig, error) {
	dbc := &Config{}
	var err error

	dbc.Host, err = getvar(Host)
	if err != nil {
		return nil, err
	}

	dbc.Port, err = getvar(Port)
	if err != nil {
		return nil, err
	}

	dbc.User, err = getvar(User)
	if err != nil {
		return nil, err
	}

	dbc.Password, err = getvar(Password)
	if err != nil {
		return nil, err
	}

	dbc.DBName, err = getvar(Name)
	if err != nil {
		return nil, err
	}

	dbc.SslMode, err = getvar(Sslmode)
	if err != nil {
		return nil, err
	}

	dbc.Timezone, err = getvar(Timezone)
	if err != nil {
		return nil, err
	}

	return dbc, nil
}

func (c Config) MakeDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&timezone=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SslMode, c.Timezone)
}
