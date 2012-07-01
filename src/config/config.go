package config

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
)

var (
	def *Config
)

type Config struct {
	data     map[string]interface{}
	filename string
}

func LoadConfigFromFile(file string) (cfg *Config) {
	cfg = new(Config)
	cfg.data = make(map[string]interface{})
	cfg.filename = file

	err := cfg.parse()
	if err != nil {
		log.Fatalf("Error loading config file %s: %s", file, err)
	}

	if def == nil {
		def = cfg
	}

	return
}

func GetDefaultConfig() *Config {
	return def
}

func (x *Config) parse() error {
	f, err := os.Open(x.filename)
	if err != nil {
		return err
	}
	defer f.Close()
	b := new(bytes.Buffer)
	_, err = b.ReadFrom(f)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b.Bytes(), &x.data)
	if err != nil {
		return err
	}

	return nil
}

// Returns a string for the config variable key
func (x *Config) GetString(key string) string {
	result, present := x.data[key]
	if !present {
		return ""
	}
	return result.(string)
}

// Returns an int for the config variable key
func (x *Config) GetInt(key string) int {
	v, ok := x.data[key]
	if !ok {
		return -1
	}
	return int(v.(float64))
}

// Returns a boolean for the config variable key
func (x *Config) GetBool(key string) bool {
	if v, ok := x.data[key]; ok {
		if v == "yes" {
			return true
		}
	} 
	
	return false
}



// Returns an array for the config variable key
func (x *Config) GetArray(key string) []interface{} {
	result, present := x.data[key]
	if !present {
		return []interface{}(nil)
	}
	return result.([]interface{})
}
