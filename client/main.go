package main

import (
	"fmt"
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func main()  {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "up")
	})

	http.HandleFunc("/outlier", func(w http.ResponseWriter, r *http.Request) {
		client := &http.Client{
			Timeout: time.Second * 10,
		}

		for i := 0; i < 5; i++ {
			req, _ := http.NewRequest("GET", "http://envoy:10001/failed", new(bytes.Buffer))

			res, err := client.Do(req)

			if err != nil {
				http.Error(w, err.Error(), 500)
				continue
			}
			if res.StatusCode == 500 {
				http.Error(w, "error", 500)
			}
			fmt.Printf("%d\n", res.StatusCode)
			defer res.Body.Close()
		}
		fmt.Println("outlier")
	})

	http.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		params := request.URL.Query()
		times, err := strconv.Atoi(params.Get("times"))
		if err != nil {
			times = 1
		}
		instances, err := strconv.Atoi(params.Get("instances"))
		if err != nil {
			instances = 1
		}
		if err = testCircuitBreaker(writer, instances, times); err != nil {
			http.Error(writer, "failed testing circuitbreaker", 500)
		}
	})

	fmt.Println("stared")
	http.ListenAndServe(":6061", nil)
	fmt.Println("finished")
}

func testCircuitBreaker(w http.ResponseWriter, instance, number int) error {
	successCounter := 0
	failureCounter := 0
	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(instance)
	for i := 0; i < instance; i++ {
		go func(instanceID int) {
			defer wg.Done()

			client := &http.Client{
				Timeout: time.Second * 10,
			}

			for j := 1; j <= number; j++ {
				start := time.Now()
				elaspsed := time.Since(start).String()

				req, err := http.NewRequest("GET", "http://envoy:10001/ping", new(bytes.Buffer))

				res, err := client.Do(req)
				if err != nil || res.StatusCode >= 400 {
					fmt.Printf("failed %d, %d, %v, %v\n", instanceID, j, elaspsed, err)
					failureCounter++
					continue
				}
				defer res.Body.Close()

				respBody, err := ioutil.ReadAll(res.Body)

				successCounter++
				fmt.Printf("success: %s %d, %d, %v\n", respBody, instanceID, j, elaspsed)
			}
		}(i)
	}
	wg.Wait()
	elaspsed := time.Since(start)
	fmt.Printf("circuitbreaker completed in %v\n", elaspsed.String())
	fmt.Printf("success %d, failed: %d\n", successCounter, failureCounter)
	fmt.Fprint(w, "done")
	return nil
}
