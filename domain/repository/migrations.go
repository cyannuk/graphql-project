package repository

import (
	"fmt"
	"net/url"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"

	"graphql-project/config"
	"graphql-project/core"
	"graphql-project/data"
)

func ApplyMigrations(config *config.Config) (err error) {
	u, err := url.Parse(fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", config.DbUser(), config.DbPassword(), config.DbHost(), config.DbPort(), config.DbName()))
	if err != nil {
		return
	}
	if config.DbMigrate() {
		db := dbmate.New(u)
		db.FS = data.AssetFS()
		db.Log = core.LoggerWriter()
		db.AutoDumpSchema = false
		db.Verbose = true
		err = db.CreateAndMigrate()
	}
	return
}
