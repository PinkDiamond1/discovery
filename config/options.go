// Copyright (c) 2021 BlockDev AG
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
)

type Options struct {
	DbDSN             string
	QualityOracleURL  url.URL
	BrokerURL         url.URL
	UniverseJWTSecret string
	RedisAddress      string
	RedisPass         string
	RedisDB           int
	BadgerAddress     url.URL
	LocationAddress   url.URL
	LocationUser      string
	LocationPass      string
}

func Read() (*Options, error) {
	dsn, err := RequiredEnv("DB_DSN")
	if err != nil {
		return nil, err
	}
	qualityOracleURL, err := RequiredEnvURL("QUALITY_ORACLE_URL")
	if err != nil {
		return nil, err
	}
	brokerURL, err := RequiredEnvURL("BROKER_URL")
	if err != nil {
		return nil, err
	}
	universeJWTSecret, err := RequiredEnv("UNIVERSE_JWT_SECRET")
	if err != nil {
		return nil, err
	}
	redisAddress, err := RequiredEnv("REDIS_ADDRESS")
	if err != nil {
		return nil, err
	}
	badgerAddress, err := RequiredEnvURL("BADGER_ADDRESS")
	if err != nil {
		return nil, err
	}

	locationUser := OptionalEnv("LOCATION_USER", "")
	locationPass := OptionalEnv("LOCATION_PASS", "")
	locationAddress, err := RequiredEnvURL("LOCATION_ADDRESS")
	if err != nil {
		return nil, err
	}

	redisPass := OptionalEnv("REDIS_PASS", "")

	redisDBint := 0
	redisDB := OptionalEnv("REDIS_DB", "0")
	if redisDB != "" {
		res, err := strconv.Atoi(redisDB)
		if err != nil {
			return nil, fmt.Errorf("could not parse redis db from %q: %w", redisDB, err)
		}
		redisDBint = res
	}
	return &Options{
		DbDSN:             dsn,
		QualityOracleURL:  *qualityOracleURL,
		BrokerURL:         *brokerURL,
		UniverseJWTSecret: universeJWTSecret,
		RedisAddress:      redisAddress,
		RedisPass:         redisPass,
		RedisDB:           redisDBint,
		BadgerAddress:     *badgerAddress,
		LocationAddress:   *locationAddress,
		LocationUser:      locationUser,
		LocationPass:      locationPass,
	}, nil
}

func RequiredEnv(key string) (string, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("required environment variable is misssing: %s", key)
	}
	return val, nil
}

func RequiredEnvURL(key string) (*url.URL, error) {
	strVal, err := RequiredEnv(key)
	if err != nil {
		return nil, err
	}
	parsedURL, err := url.Parse(strVal)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s from value '%s'", key, strVal)
	}
	return parsedURL, nil
}

func OptionalEnvURL(key string, defaults string) (*url.URL, error) {
	strVal := OptionalEnv(key, defaults)
	parsedURL, err := url.Parse(strVal)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s from value '%s'", key, strVal)
	}
	return parsedURL, nil
}

func OptionalEnv(key string, defaults string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaults
	}
	return val
}

func OptionalEnvBool(key string) bool {
	val, _ := strconv.ParseBool(os.Getenv(key))
	return val
}
