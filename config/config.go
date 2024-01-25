package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

const (
	difficultDefault  int    = 4
	expirationDefault int    = 60
	hostDefault       string = "127.0.0.1"
	portDefault       string = "8080"
	redisHost         string = "0.0.0.0"
	redisPort         string = "6379"
)

func LoadEnv(path string) {
	err := godotenv.Load(path)
	if err != nil {
		log.Printf("error loading %s file", path)
	}
}

func GetRedisHost() string {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		log.Printf("REDIS_HOST not found, will be used default: %d", redisHost)
		host = redisHost
	}

	port := os.Getenv("REDIS_PORT")
	if port == "" {
		log.Printf("REDIS_PORT not found, will be used default: %d", port)
		port = redisPort
	}

	return fmt.Sprintf("%s:%s", host, port)
}

func GetExpiration() int {
	exp := os.Getenv("HASHCASH_EXPIRATION")

	if exp == "" {
		log.Printf("HASHCASH_EXPIRATION not found, will be used default: %d", expirationDefault)
		return expirationDefault
	}
	expiration, err := strconv.Atoi(exp)
	if err != nil {
		log.Printf("failed to convert expiration %d to int, will be used default: %d", expiration, expirationDefault)
		return expirationDefault
	}

	return expiration
}

func GetAddress() string {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	if host == "" {
		host = hostDefault
		log.Printf("HOST not found in environment, will be used default %s", hostDefault)
	}
	if port == "" {
		port = portDefault
		log.Printf("PORT not found in environment, will be used default %s", portDefault)
	}

	return fmt.Sprintf("%s:%s", host, port)

}

func GetDifficulty() int {
	difficulty := os.Getenv("HASHCASH_DIFFICULTY")

	if difficulty == "" {
		log.Printf("HASHCASH_DIFFICULTY not found, will be used default, %d", difficultDefault)
		return difficultDefault
	}

	v, err := strconv.Atoi(difficulty)
	if err != nil {
		log.Printf("failed to convert %s to int, will be used default, %d", difficulty, difficultDefault)
		return difficultDefault
	}

	return v
}
