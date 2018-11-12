package config

import (
	"encoding/json"
	"github.com/ONSdigital/go-ns/log"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Model struct {
	MongoURL         string   `yaml:"mongo-url"`
	MongoCollections []string `yaml:"mongo-collections"`
	Neo4jURL         string   `yaml:"neo4j-url"`
}

func Load() (*Model, error) {
	source, err := ioutil.ReadFile("config.yml")
	if err != nil {
		return nil, errors.Wrap(err, "error reading config.yml")
	}

	var config Model
	if err := yaml.Unmarshal(source, &config); err != nil {
		return nil, errors.Wrap(err, "error marshalling config.yml")
	}

	log.Debug("cmd-cli configuration", log.Data{"": config.String()})
	return &config, nil
}

func (c Model) String() string {
	s, _ := json.MarshalIndent(&c, "", "    ")
	return string(s)
}
