package main

//
//// registerService registers the correct entry for this service in etcd
//func (app *Config) registerService() {
//	cli, _ := connectToEtcd()
//	kv := clientv3.NewKV(cli)
//
//	app.Etcd = cli
//
//	lease := clientv3.NewLease(cli)
//	grantResp, err := lease.Grant(context.TODO(), 10)
//	if err != nil {
//		log.Println("Error creating lease", err)
//	}
//
//	// insert something with the lease
//	_, err = kv.Put(context.TODO(), fmt.Sprintf("/auth/%s", app.randomString(32)), "authentication-service", clientv3.WithLease(grantResp.ID))
//	if err != nil {
//		log.Println("Error inserting using lease", err)
//	}
//
//	// keep lease alive
//	kalRes, err := lease.KeepAlive(context.TODO(), grantResp.ID)
//	if err != nil {
//		log.Println("Error with keepalive", err)
//	}
//	go app.listenToKeepAlive(kalRes)
//}
//
//// listenToKeepAlive just consumes channel responses from etcd's KeepAlive method
//func (app *Config) listenToKeepAlive(kalRes <-chan *clientv3.LeaseKeepAliveResponse) {
//	defer func() {
//		if r := recover(); r != nil {
//			log.Println("Error", fmt.Sprintf("%v", r))
//		}
//	}()
//
//	// the only reason this exists is to consume the response from etcd's KeepAlive, because
//	// if we don't, unexpected behaviour is the result.
//	for {
//		_ = <-kalRes
//	}
//}
//
//// connectToEtcd tries to connect to etcd, for up to 30 seconds
//func connectToEtcd() (*clientv3.Client, error) {
//	var cli *clientv3.Client
//	var counts = 0
//
//	for {
//		c, err := clientv3.New(clientv3.Config{Endpoints: []string{"etcd:2379"},
//			DialTimeout: 5 * time.Second,
//		})
//		if err != nil {
//			fmt.Println("etcd not ready...")
//			counts++
//		} else {
//			fmt.Println()
//			cli = c
//			break
//		}
//
//		if counts > 15 {
//			return nil, err
//		}
//		fmt.Println("Backing off for 2 seconds...")
//		time.Sleep(2 * time.Second)
//		continue
//	}
//	log.Println("Connected to etcd!")
//	return cli, nil
//}
