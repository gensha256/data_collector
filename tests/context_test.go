package tests

import (
	"context"
	"log"
	"math/rand"
	"testing"
	time "time"
)

func TestPrintInt(t *testing.T) {
	t.Skip()
	testCh := make(chan int)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	go PrintInt(ctxTimeout, testCh)

	for {
		time.Sleep(time.Second * 2)
		select {
		case <-ctxTimeout.Done():
			{
				if ctxTimeout.Err() != context.DeadlineExceeded {
					t.Fail()
				}
				return
			}
		case res := <-testCh:
			time.Sleep(time.Second * 1)
			log.Println("Numbers from channel:", res)
		}
	}
}

func TestPrint(t *testing.T) {
	t.Skip()
	chanTest := make(chan int)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	go PrintInt(ctx, chanTest)

	select {

	case <-time.After(time.Millisecond * 2200):
		if ctx.Err() != context.DeadlineExceeded || ctx.Done() == nil {
			t.Error("context not canceled")
			t.Fail()
		}
	}
}

func TestContextWithValue(t *testing.T) {
	t.Skip()
	ctx := context.Background()

	ctxValue := SendValue(ctx)

	testData := "Mark"

	if ctxValue != testData {
		t.Error("not correct value")
		t.Fail()
	}
}

func PrintInt(ctx context.Context, ch chan int) {

	log.Println("Add data in channel")
	for i := 0; i < 5; i++ {
		randNum := rand.Intn(100)

		select {
		case <-ctx.Done():
			{
				log.Println(ctx.Err())
				return
			}
		case ch <- randNum:
		}
	}
}

func SendValue(ctx context.Context) any {
	data := "Mark"
	ctxVal := context.WithValue(ctx, "user", data)
	return GetValue(ctxVal)
}

func GetValue(ctx context.Context) any {
	value := ctx.Value("user")
	return value
}
