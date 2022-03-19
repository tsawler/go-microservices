package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"strings"
	"time"
)

//
func (app *Config) watchEtcd() {
	for {
		// watch for service changes
		watchKey := app.Etcd.Watch(context.Background(), "/mail/", clientv3.WithPrefix())
		for resp := range watchKey {
			for _, item := range resp.Events {
				// get our values as strings so we can work with them
				eventType := item.Type.String()
				key := string(item.Kv.Key)
				value := string(item.Kv.Value)
				var deleteURL = false
				if strings.HasPrefix(eventType, "DELETE") {
					deleteURL = true
				}

				// add to or remove from service maps (using url as key, and empty string as value)
				switch {
				case strings.HasPrefix(key, "mail"):
					// mail
					if deleteURL {
						delete(app.MailServiceURLs, key)
					} else {
						app.MailServiceURLs[value] = ""
					}

				case strings.HasPrefix(key, "logger"):
					// logger
					if deleteURL {
						delete(app.LogServiceURLs, key)
					} else {
						app.LogServiceURLs[value] = ""
					}

				case strings.HasPrefix(key, "auth"):
					// authentication
					if deleteURL {
						delete(app.AuthServiceURLs, key)
					} else {
						app.AuthServiceURLs[value] = ""
					}
				}
			}
		}
	}
}

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
