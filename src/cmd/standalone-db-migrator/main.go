package main

import (
	"os"
	"strconv"

	"github.com/goharbor/harbor/src/common"
	"github.com/goharbor/harbor/src/common/dao"
	"github.com/goharbor/harbor/src/common/models"
	"github.com/goharbor/harbor/src/lib/log"
	"github.com/goharbor/harbor/src/migration"
)

// key: env var, value: default value
var defaultAttrs = map[string]string{

	common.DatabaseType: common.DatabaseType_PostGreSQL,

	"POSTGRESQL_HOST":     "localhost",
	"POSTGRESQL_PORT":     "5432",
	"POSTGRESQL_USERNAME": "postgres",
	"POSTGRESQL_PASSWORD": "password",
	"POSTGRESQL_DATABASE": "registry",
	"POSTGRESQL_SSLMODE":  "disable",

	"MYSQL_HOST": "localhost",
	"MYSQL_PORT": "3306",
	"MYSQL_USERNAME": "root",
	"MYSQL_PASSWORD": "password",
	"MYSQL_DATABASE": "registry",
}

func main() {
	p, _ := strconv.Atoi(getAttr("POSTGRESQL_PORT"))
	mysql_port, _ := strconv.Atoi(getAttr("MYSQL_PORT"))
	db := &models.Database{
		Type: getAttr(common.DatabaseType),
		PostGreSQL: &models.PostGreSQL{
			Host:         getAttr("POSTGRESQL_HOST"),
			Port:         p,
			Username:     getAttr("POSTGRESQL_USERNAME"),
			Password:     getAttr("POSTGRESQL_PASSWORD"),
			Database:     getAttr("POSTGRESQL_DATABASE"),
			SSLMode:      getAttr("POSTGRESQL_SSLMODE"),
			MaxIdleConns: 5,
			MaxOpenConns: 5,
		},
		MySQL: &models.MySQL{
			Host:     getAttr("MYSQL_HOST"),
			Port:     mysql_port,
			Username: getAttr("MYSQL_USERNAME"),
			Password: getAttr("MYSQL_PASSWORD"),
			Database: getAttr("MYSQL_DATABASE"),
		},
	}

	log.Info("Migrating the data to latest schema...")
	switch db.Type {
	case "", common.DatabaseType_PostGreSQL:
		log.Infof("DB info: postgres://%s@%s:%d/%s?sslmode=%s", db.PostGreSQL.Username, db.PostGreSQL.Host,
			db.PostGreSQL.Port, db.PostGreSQL.Database, db.PostGreSQL.SSLMode)
	case common.DatabaseType_MySQL:
		log.Infof("DB info: mysql://%s@%s:%d/%s", db.MySQL.Username, db.MySQL.Host,
			db.MySQL.Port, db.MySQL.Database)
	default:
		log.Fatalf("DB info: unknown type - %s", db.Type)
	}

	if err := dao.InitDatabase(db); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	if err := migration.MigrateDB(db); err != nil {
		log.Fatalf("failed to migrate DB: %v", err)
	}
	log.Info("Migration done.  The data schema in DB is now update to date.")
}

func getAttr(k string) string {
	v := os.Getenv(k)
	if len(v) > 0 {
		return v
	}
	return defaultAttrs[k]
}
