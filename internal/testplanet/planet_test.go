package testplanet

import (
	"context"
	"strconv"
	"testing"
)

func TestBasic(t *testing.T) {
	t.Log("New")
	planet, err := New(1, 4)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		t.Log("Shutdown")
		err = planet.Shutdown()
		if err != nil {
			t.Fatal(err)
		}
	}()

	t.Log("Start")
	planet.Start(context.Background())
}

func BenchmarkCreate(b *testing.B) {
	storageNodes := []int{4, 10, 100}
	for _, count := range storageNodes {
		b.Run(strconv.Itoa(count), func(b *testing.B) {
			ctx := context.Background()
			for i := 0; i < b.N; i++ {
				planet, err := New(1, 100)
				if err != nil {
					b.Fatal(err)
				}

				planet.Start(ctx)

				err = planet.Shutdown()
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}