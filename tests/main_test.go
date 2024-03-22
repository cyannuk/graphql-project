package tests

import (
	"fmt"
	"os"
	"testing"

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

	return m.Run()
}
