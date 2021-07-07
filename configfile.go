package main

import (
	_ "embed"
	"os"

	"github.com/ppreeper/dbq/database"
	"gopkg.in/yaml.v2"
)

// Conf array of Dbase
type Conf struct {
	Dbases []database.Database `json:"dbases,omitempty"`
}

func (c *Conf) getConf(configFile string) (*Conf, error) {
	yamlFile, err := os.ReadFile(configFile)
	checkErr(err)
	err = yaml.Unmarshal(yamlFile, c)
	checkErr(err)
	return c, err
}

func (c *Conf) getDB(name string) (d database.Database) {
	for _, v := range c.Dbases {
		if v.Name == name {
			d = v
		}
	}
	return
}
