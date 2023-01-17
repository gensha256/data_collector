package tests

import (
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/gensha256/data_collector/pkg/models"
	"github.com/gensha256/data_collector/pkg/store"
)

var numbers = make(chan int, 2)

func TestSemaphore(t *testing.T) {
	t.Skip()
	wg := sync.WaitGroup{}

	rs, err := store.NewRedisStore()
	if err != nil {
		t.Fail()
	}

	chSymbol := make(chan string)

	wg.Add(2)
	go Writer(chSymbol, rs, &wg)

	testCh := <-chSymbol
	if len(testCh) == 0 {
		t.Error("not added data in redis")
		t.Fail()
	}

	if len(chSymbol) == 0 {
		t.Fail()
	}

	go Listener(chSymbol, &wg)

	for {

	}
}

func Writer(c chan string, rs *store.RedisStore, wg *sync.WaitGroup) {

	ent := models.CmcEntity{
		Id:               -1,
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
		LastUpdated:      time.Date(2022, 12, 26, 13, 54, 00, 00, time.UTC),
	}

	for {

		randNum := rand.Intn(1000 - 100)
		time.Sleep(time.Duration(randNum) * time.Millisecond)

		err := rs.JustStore(ent)
		if err != nil {
			log.Println(err)
		}

		select {
		case _, ok := <-c:
			{
				if !ok {
					log.Println("cant write data in closed channel")
					return
				}
			}
		case c <- ent.Symbol:
		}
	}
}

func Listener(c chan string, wg *sync.WaitGroup) {
	ticker := time.NewTicker(2 * time.Second)

	counter := 0

	for {

		select {

		case symbol, ok := <-c:
			{
				if !ok {
					log.Printf("closed channel")
					return
				}
			}

			counter += len(symbol)
		case <-ticker.C:
			log.Printf("was added %d XLM symbols", counter)
			close(c)
		}
	}
}

func TestSem(t *testing.T) {
	t.Skip()
	wg := sync.WaitGroup{}
	defer wg.Wait()

	for num := range []int{0, 1, 2, 3, 4, 5} {

		wg.Add(1)
		go inParallel(&wg, num)
	}
}

func inParallel(wg *sync.WaitGroup, n int) {

	numbers <- n
	log.Println("Numbers", n)

	time.Sleep(1000 * time.Millisecond)
	<-numbers

	defer wg.Done()
}
