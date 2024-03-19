package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
	"graphql-project/config"
	"graphql-project/domain/repository"
)

var (
	Cfg        config.Config
	DataSource *repository.DataSource
)

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	err := Cfg.Load(config.DisableFlags(), config.Files("../.env", "../test.env"))
	if err != nil {
		fmt.Println(err)
		return 1
	}
	DataSource, err = repository.NewDataSource(&Cfg)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	defer DataSource.Close()

	db := DataSource.OpenDB()
	defer db.Close()

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("testdata/fixtures"),
	)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	fixtures.Load()
	_ = fixtures
	// https://github.com/helloticket/go-dbunit
	// https://github.com/cjwcjswo/dbunit
	// https://github.com/go-testfixtures/testfixtures
	// https://github.com/bluele/factory-go
	// https://github.com/Pallinder/go-randomdata
	return m.Run()
}
