package main

import (
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

// GetServiceURL will get a service's url from those listed as available in etcd
func (app *Config) GetServiceURL(serviceType string) (string, error) {
	var serviceURL string

	// TODO - get service URL from etcd
	switch serviceType {

	}

	return serviceURL, nil
}

// connectToRabbit tries to connect to etcd, for up to 30 seconds
func connectToEtcd() (*clientv3.Client, error) {
	var cli *clientv3.Client
	var counts = 0

	for {
		c, err := clientv3.New(clientv3.Config{Endpoints: []string{"etcd:2379"},
			DialTimeout: 5 * time.Second,
		})
		if err != nil {
			fmt.Println("etcd not ready...")
			counts++
		} else {
			fmt.Println()
			cli = c
			break
		}

		if counts > 15 {
			return nil, err
		}
		fmt.Println("Backing off for 2 seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
	log.Println("Connected to etcd!")
	return cli, nil
}
