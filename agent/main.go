package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"consts"
)

type OperationData struct {
	Id   uint          `json:"id"`
	Op   string        `json:"op"`
	Time time.Duration `json:"time"`
	Args []float64     `json:"args"`
}

type Result struct {
	Id     uint    `json:"id"`
	Result float64 `json:"result"`
}

func main() {

	master, ok := os.LookupEnv(consts.EnvMaster)
	if !ok {
		master = consts.AgentDefaultMaster
	}

	workers := consts.AgentDefaultWorkers
	if workersStr, ok := os.LookupEnv(consts.EnvNWorkers); ok {
		val, err := strconv.ParseUint(workersStr, 10, 0)
		if err == nil {
			workers = int(val)
		}
	}

	batch := consts.AgentDefaultBatch
	if batchStr, ok := os.LookupEnv(consts.EnvBatch); ok {
		val, err := strconv.ParseUint(batchStr, 10, 0)
		if err == nil {
			batch = int(val)
		}
	}

	delay := consts.AgentDefaultDelay * time.Second
	if delayStr, ok := os.LookupEnv(consts.EnvDelay); ok {
		val, err := strconv.ParseUint(delayStr, 10, 0)
		if err == nil {
			delay = time.Duration(val) * time.Second
		}
	}

	var mx sync.Mutex
	ch := make(chan OperationData, workers)

	for i := 0; i < workers; i++ {
		go func() {
			for {
				r := <-ch
				mx.Lock()
				workers -= 1
				mx.Unlock()
				time.Sleep(r.Time)

				a, b := r.Args[0], r.Args[1]
				c := float64(0)
				switch r.Op {
				case "+":
					c = a + b
				case "-":
					c = a - b
				case "*":
					c = a * b
				case "/":
					c = a / b
				}

				go func() {
					url := fmt.Sprintf(`%s/result`, master)
					data := Result{
						Id:     r.Id,
						Result: c,
					}

					binary, err := json.Marshal(data)
					if err != nil {
						log.Print(err)
						return
					}

					req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(binary))
					if err != nil {
						log.Print(err)
						return
					}

					req.Header.Set("Content-Type", "application/json")

					client := &http.Client{}
					resp, err := client.Do(req)
					if err != nil {
						log.Print(err)
						return
					}
					err = resp.Body.Close()
					if err != nil {
						log.Print(err)
						return
					}
				}()

				mx.Lock()
				workers += 1
				mx.Unlock()
			}
		}()
	}

	t := time.Tick(delay)
	for _ = range t {
		mx.Lock()

		reqCount := workers
		if batch < reqCount {
			reqCount = batch
		}

		mx.Unlock()

		binary, err := json.Marshal(data)
		if err != nil {
			log.Print(err)
			return
		}

		req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(binary))
		if err != nil {
			log.Print(err)
			return
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Print(err)
			return
		}
		err = resp.Body.Close()
		if err != nil {
			log.Print(err)
			return
		}

	}
}
