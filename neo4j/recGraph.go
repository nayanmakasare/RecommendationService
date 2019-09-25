package neo4j

import (
	"github.com/micro/go-micro/util/log"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

var (
	driver neo4j.Driver
	session neo4j.Session
	result neo4j.Result
	err error
)


func InitNeo4j() error{
	log.Info("Neo Trigggerd ")
	if driver, err = neo4j.NewDriver("bolt://localhost:7687", neo4j.BasicAuth("neo4j", "tlwn722n", "")); err != nil {
		return err // handle error
	}
	// handle driver lifetime based on your application lifetime requirements
	// driver's lifetime is usually bound by the application lifetime, which usually implies one driver instance per application
	defer driver.Close()

	if session , err = driver.Session(neo4j.AccessModeWrite); err != nil {
		return err
	}

	defer session.Close()

	result, err := session.Run("CREATE (m:Movies {title: $title, released: $year}) RETURN m.title, m.released", map[string]interface{}{
		"title":"Cloudwalker",
		"year":2019,
	})
	if err != nil {
		return err
	}
	log.Info("Bingo neo4j ")
	if err = result.Err(); err != nil {
		return err
	}
	return nil
}
























