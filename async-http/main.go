package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
}

const URL string = "https://dummyjson.com/users"

func GetUserById(ctx context.Context, id int, delaySeconds int) (*User, error) {
	url := fmt.Sprintf(URL+"/%d?delay=%d", id, delaySeconds)
	fmt.Printf("fetching from [%s]\n", url)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	fmt.Printf("received response for %s\n", url)

	var user User
	err = json.NewDecoder(response.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

type Response struct {
	User  *User `json:"user"`
	Error error `json:"error"`
}

func Do() (resCh chan Response) {
	ctx := context.Background()
	const NUM_OF_REQUESTS int = 3
	var wg sync.WaitGroup

	resCh = make(chan Response, NUM_OF_REQUESTS)

	for i := 1; i <= NUM_OF_REQUESTS; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			ctx, cancelFunc := context.WithTimeout(ctx, time.Second*3)
			defer cancelFunc()
			user, err := GetUserById(ctx, id, 3)
			resCh <- Response{
				User:  user,
				Error: err,
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(resCh)
	}()

	return resCh
}

func main() {
	res := Do()
	for r := range res {
		json, err := json.Marshal(r)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%s\n", string(json))
	}
}
