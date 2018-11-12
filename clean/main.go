package main

import (
	"github.com/ONSdigital/cmd-dev-utils/clean/config"
	"github.com/ONSdigital/go-ns/log"
	"github.com/globalsign/mgo"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"os"
)

func main() {
	log.HumanReadable = true

	cfg, err := config.Load()
	if err != nil {
		exit(err)
	}

	if err := dropMongo(cfg); err != nil {
		exit(err)
	}
	if err := dropNeo4j(cfg); err != nil {
		exit(err)
	}
}

func dropMongo(cfg *config.Model) error {
	log.Info("dropping mongo databases", nil)
	sess, err := mgo.Dial(cfg.MongoURL)
	if err != nil {
		return err
	}
	defer sess.Close()

	for _, db := range cfg.MongoCollections {
		log.Info("dropping "+db, nil)
		if err := sess.DB(db).DropDatabase(); err != nil {
			return err
		}
	}
	log.Info("drop mongo collections completed successfully", nil)
	return nil
}

func dropNeo4j(cfg *config.Model) error {
	log.Info("dropping neo4j database", nil)
	pool, err := bolt.NewDriverPool(cfg.Neo4jURL, 1)
	if err != nil {
		return err
	}

	conn, err := pool.OpenPool()
	if err != nil {
		return err
	}
	defer conn.Close()
	res, err := conn.ExecNeo("MATCH(n) DETACH DELETE n", nil)
	log.Debug("results", log.Data{
		"delete results": res.Metadata()["stats"],
	})
	if err != nil {
		return err
	}
	log.Info("drop neo4j completed successfully", nil)
	return nil
}

func exit(err error) {
	log.Error(err, nil)
	os.Exit(1)
}
