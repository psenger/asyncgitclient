package main

import (
	"net/url"
	"net/http"
	"encoding/json"
	"fmt"
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

func f( page string, user chan *UsersPaginationWrapper ) {
	rel := &url.URL{
		Path: "/api/users",
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

func main ()  {

	go func () {
		messages := make(chan *UsersPaginationWrapper, 1)
		go f("1", messages)
		msg := <- messages
		close(messages)
		fmt.Printf("done %+v\n", msg)
	}()

}