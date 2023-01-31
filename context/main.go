package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"context"
	"errors"
	"log"
	"math/rand"
)

const keyID = "id"

func operation1() error {
	time.Sleep(100 * time.Millisecond)
	return errors.New("failed")
}

func operation2(ctx context.Context) {
	select {
	case <-time.After(500 * time.Millisecond):
		fmt.Println("done")
	case <-ctx.Done():
		fmt.Println("halted operation2")
	}
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)


	go func() {
		err := operation1()
		if err != nil {
			cancel()
		}
	}()

	operation2(ctx)

	ctx = context.Background()
	ctx, _ = context.WithTimeout(ctx, 1000 * time.Millisecond)

	req, _ := http.NewRequest(http.MethodGet, "https://google.com", nil)
	req = req.WithContext(ctx)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request failed: ", err)
	}

	fmt.Println("Response received status code: ", resp.StatusCode)

	rand.Seed(time.Now().Unix())
	ctx = context.WithValue(context.Background(), keyID, rand.Int())
	operation3(ctx)

	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		fmt.Fprintf(os.Stdout, "processing request\n")

		select {
		case <-time.After(2 * time.Second):
			w.Write([]byte("request processed"))
		case <-ctx.Done():
			fmt.Fprintf(os.Stderr, "request cancelled\n")
		}
	}))
}

func operation3(ctx context.Context) {
	log.Println("operation3 for id: ", ctx.Value(keyID), " completed")
	operation4(ctx)
}

func operation4(ctx context.Context) {
	log.Println("operation4 for id: ", ctx.Value(keyID), " completed")
}
