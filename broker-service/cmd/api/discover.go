package main

import (
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

//func (app *Config) getServiceURLs() {
//	kv := clientv3.NewKV(app.Etcd)
//	app.MailServiceURLs = make(map[string]string)
//	app.LogServiceURLs = make(map[string]string)
//	app.AuthServiceURLs = make(map[string]string)
//
//	prefixes := []string{"/mail/", "/logger/", "/auth/"}
//
//	// range through all the services we want to discover
//	for _, curPrefix := range prefixes {
//		getResp, err := kv.Get(context.TODO(), curPrefix, clientv3.WithPrefix())
//		if err != nil {
//			log.Println(err)
//		}
//
//		for _, k := range getResp.Kvs {
//			//log.Println("Key", string(k.Key))
//			//log.Println("Adding", string(k.Value), "to", curPrefix, "service map; key was", string(k.Key))
//			switch curPrefix {
//			case "/mail/":
//				app.MailServiceURLs[string(k.Value)] = ""
//			case "/logger/":
//				app.LogServiceURLs[string(k.Value)] = ""
//			case "/auth/":
//				app.AuthServiceURLs[string(k.Value)] = ""
//			}
//		}
//	}
//}

// watchEtcd runs in the background, looking for changes in etcd. When it finds changes
// hosts, it updates the appropriate map in the *Config receiver.
//func (app *Config) watchEtcd() {
//	for {
//		// watch for service changes
//		watchKey := app.Etcd.Watch(context.Background(), "/mail/", clientv3.WithPrefix())
//		for resp := range watchKey {
//			for _, item := range resp.Events {
//				// get our values as strings so that we can work with them
//				eventType := item.Type.String()
//				key := string(item.Kv.Key)
//				value := string(item.Kv.Value)
//				var deleteURL = false
//				if strings.Contains(eventType, "DELETE") {
//					deleteURL = true
//				}
//
//				// add to or remove from service maps (using url as key, and empty string as value)
//				switch {
//				case strings.HasPrefix(key, "mail"):
//					// mail
//					if deleteURL {
//						log.Println("Removing", value, "from mail service map")
//						delete(app.MailServiceURLs, key)
//					} else {
//						log.Println("Adding", value, "to mail service map")
//						app.MailServiceURLs[value] = ""
//					}
//
//				case strings.HasPrefix(key, "logger"):
//					// logger
//					if deleteURL {
//						delete(app.LogServiceURLs, key)
//					} else {
//						app.LogServiceURLs[value] = ""
//					}
//
//				case strings.HasPrefix(key, "auth"):
//					// authentication
//					if deleteURL {
//						delete(app.AuthServiceURLs, key)
//					} else {
//						app.AuthServiceURLs[value] = ""
//					}
//				}
//			}
//		}
//	}
//}

// GetServiceURL will get a service's url from those listed as available in etcd
//func (app *Config) GetServiceURL(serviceType string) string {
//	var serviceURL string
//
//	// get service URL from etcd
//	switch serviceType {
//	case "mail":
//		serviceURL = getUrlFromMap(app.MailServiceURLs)
//	case "logger":
//		serviceURL = getUrlFromMap(app.LogServiceURLs)
//	case "auth":
//		serviceURL = getUrlFromMap(app.AuthServiceURLs)
//	}
//
//	return serviceURL
//}

// getUrlFromMap returns a random value from available urls in
// service maps. Since maps are never guaranteed to be in the same order,
// grabbing the first value is sufficient for our purposes.
func getUrlFromMap(m map[string]string) string {
	var u string
	for k := range m {
		u = k
		break
	}
	return u
}

// connectToEtcd tries to connect to etcd, for up to 30 seconds
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
	fmt.Println("Connected to etcd!")
	return cli, nil
}
