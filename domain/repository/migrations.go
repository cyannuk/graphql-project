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

func ApplyMigrations(config *config.Config) error {
	if config.DbMigrate() {
		if u, err := url.Parse(fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", config.DbUser(), config.DbPassword(), config.DbHost(), config.DbPort(), config.DbName())); err != nil {
			return err
		} else {
			db := dbmate.New(u)
			db.FS = data.AssetFS()
			db.Log = core.LoggerWriter()
			db.AutoDumpSchema = false
			db.Verbose = true
			if err = db.CreateAndMigrate(); err != nil {
				return err
			}
		}
	}
	return nil
}
