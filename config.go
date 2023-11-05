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

	dbc.Host, err = getvar(defaultConfig.Host)
	if err != nil {
		return nil, err
	}

	dbc.Port, err = getvar(defaultConfig.Port)
	if err != nil {
		return nil, err
	}

	dbc.User, err = getvar(defaultConfig.User)
	if err != nil {
		return nil, err
	}

	dbc.Password, err = getvar(defaultConfig.Password)
	if err != nil {
		return nil, err
	}

	dbc.DBName, err = getvar(defaultConfig.Name)
	if err != nil {
		return nil, err
	}

	dbc.SslMode, err = getvar(defaultConfig.SslMode)
	if err != nil {
		return nil, err
	}

	dbc.Timezone, err = getvar(defaultConfig.Timezone)
	if err != nil {
		return nil, err
	}

	return dbc, nil
}

func (c Config) MakeDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&timezone=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SslMode, c.Timezone)
}

type KeysConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SslMode  string
	Timezone string
}

var defaultConfig *KeysConfig = &KeysConfig{
	Host:     "DB_HOST",
	Port:     "DB_PORT",
	User:     "DB_USER",
	Password: "DB_PASSWORD",
	Name:     "DB_NAME",
	SslMode:  "DB_SSL_MODE",
	Timezone: "DB_TIMEZONE",
}

func GetConfigFromEnvOption(keys *KeysConfig) (DBConfig, error) {
	if keys == nil {
		return nil, fmt.Errorf("keysConfig is nil")
	}
	dbc := &Config{}
	var err error

	dbc.Host, err = getvar(keys.Host)
	if err != nil {
		return nil, err
	}

	dbc.Port, err = getvar(keys.Port)
	if err != nil {
		return nil, err
	}

	dbc.User, err = getvar(keys.User)
	if err != nil {
		return nil, err
	}

	dbc.Password, err = getvar(keys.Password)
	if err != nil {
		return nil, err
	}

	dbc.DBName, err = getvar(keys.Name)
	if err != nil {
		return nil, err
	}

	dbc.SslMode, err = getvar(keys.SslMode)
	if err != nil {
		return nil, err
	}

	dbc.Timezone, err = getvar(keys.Timezone)
	if err != nil {
		return nil, err
	}

	return dbc, nil
}
