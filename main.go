package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// https://reqres.in/
type UsersPaginationWrapper struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
	Data       []struct {
		ID        int    `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Avatar    string `json:"avatar"`
	} `json:"data"`
}

func fetch(page string, user chan *UsersPaginationWrapper) {
	rel := &url.URL{
		Path:   "/api/users",
		Scheme: "https",
		Host:   "reqres.in",
	}
	q := rel.Query()
	q.Add("page", page)
	rel.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", rel.String(), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil || resp.StatusCode >= 300 {
		panic(err)
	}
	var usersPaginationWrapper UsersPaginationWrapper
	err = json.NewDecoder(resp.Body).Decode(&usersPaginationWrapper)

	user <- &usersPaginationWrapper
}

func main() {

	var outputs [4]chan *UsersPaginationWrapper
	for i := range outputs {
		outputs[i] = make(chan *UsersPaginationWrapper)
		go fetch(strconv.Itoa(i+1), outputs[i])
		msg := <-outputs[i]
		close(outputs[i])
		fmt.Printf("done %d %+v\n", i+1, msg)
	}
}
