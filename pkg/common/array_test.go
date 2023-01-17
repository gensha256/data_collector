package common

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var testCmcEntity = CmcEntity{
	Id:               1,
	Name:             "Stellar Lumen",
	Symbol:           "XLM",
	MaxSupply:        1,
	CmcRank:          -1,
	Price:            1,
	Volume24h:        1,
	VolumeChange24h:  -1,
	PercentChange1h:  1,
	PercentChange24h: 1,
	MarketCap:        -1,
}

func GenerateEntityArray(size int) []CmcEntity {

	result := make([]CmcEntity, 0, size)

	for i := 0; i < size; i++ {

		randDay := rand.Intn(30)
		randHours := rand.Intn(24)

		testCmcEntity.LastUpdated = time.Date(2022, 12, randDay, randHours, 0, 0, 0, time.UTC)

		result = append(result, testCmcEntity)
	}

	return result
}

func TestSortingArray(t *testing.T) {
	notSortArr := GenerateEntityArray(10)
	sortArr := SortArray(notSortArr)

	if len(notSortArr) != len(sortArr) {
		t.Error("not equal array size")
		t.Fail()
	}

	for i := 0; i < len(sortArr)-1; i++ {

		first := sortArr[i]
		second := sortArr[i+1]

		if first.LastUpdated.Unix() > second.LastUpdated.Unix() {
			t.Error("not sorted array")
			t.Fail()
		}
	}
}

func TestDoesArrayContainString(t *testing.T) {
	arr := []string{"hello", "name", "have"}
	search := "name"
	fmt.Printf("%p \n", &arr)
	fmt.Printf("%p \n", &search)

	result := DoesArrayContainString(arr, search)

	fmt.Printf("%p \n", &arr)
	fmt.Printf("%p \n", &search)

	if result != true {
		t.Fail()
	}
}
