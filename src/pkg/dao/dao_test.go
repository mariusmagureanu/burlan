package dao

import "testing"

func TestDAO_Init(t *testing.T) {
	daoTest := DAO{}

	err := daoTest.Init("demo.sqlite")

	if err != nil {
		t.Error(err)
	}

}
