package dao

import (
	"testing"
)

func tearUp() (DAO, error) {
	daoTest := DAO{}

	err := daoTest.Init("test.sqlite")

	if err != nil {
		return daoTest, err
	}

	err = daoTest.DropTables()

	if err != nil {
		return daoTest, err
	}

	err = daoTest.CreateTables()

	return daoTest, err
}

func tearDown(dao DAO) error {
	err := dao.DropTables()

	if err != nil {
		return err
	}

	return dao.Close()
}

func TestDAO_Init(t *testing.T) {
	dao, err := tearUp()

	if err != nil {
		t.Fatal(err)
	}

	defer tearDown(dao)
}
