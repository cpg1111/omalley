package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/pullrequestrfb/omalley/addrbook"
	"github.com/pullrequestrfb/omalley/network"
)

var (
	name       = flag.String("name", "nobody", "The name everyone identifies you with")
	masterAddr = flag.String("masterAddr", "127.0.0.1:4088", "address of the master instance, defaults to 127.0.0.1:4088")
	master     = flag.Bool("master", false, "whether or not this instance is the master instance")
	dbPath     = flag.String("dbpath", "/var/lib/omalley", "path to the persistent store for master, defaults to /var/lib/omalley")
	host       = flag.String("bind-host", "0.0.0.0", "host interface to bind on, defaults to 0.0.0.0")
	port       = flag.Int("port", 8044, "port to listen on, defaults to 8044")
)
var sdnOpts network.SdnOpts

func init() {
	flag.StringVar(&sdnOpts.EtcdEndpoints, "etcd-endpoints", "http://127.0.0.1:4001,http://127.0.0.1:2379", "a comma-delimited list of etcd endpoints")
	flag.StringVar(&sdnOpts.EtcdPrefix, "etcd-prefix", "/coreos.com/network", "etcd prefix")
	flag.StringVar(&sdnOpts.EtcdKeyFile, "etcd-keyfile", "", "SSL key file used to secure etcd communication")
	flag.StringVar(&sdnOpts.EtcdCertFile, "etcd-certfile", "", "SSL certification file used to secure etcd communication")
	flag.StringVar(&sdnOpts.EtcdCAFile, "etcd-cafile", "", "SSL Certificate Authority file used to secure etcd communication")
	flag.StringVar(&sdnOpts.EtcdUsername, "etcd-username", "", "Username for BasicAuth to etcd")
	flag.StringVar(&sdnOpts.EtcdPassword, "etcd-password", "", "Password for BasicAuth to etcd")
	flag.StringVar(&sdnOpts.Listen, "listen", "", "run as server and listen on specified address (e.g. ':8080')")
	flag.StringVar(&sdnOpts.Remote, "remote", "", "run as client and connect to server on specified address (e.g. '10.1.2.3:8080')")
	flag.StringVar(&sdnOpts.RemoteKeyFile, "remote-keyfile", "", "SSL key file used to secure client/server communication")
	flag.StringVar(&sdnOpts.RemoteCertFile, "remote-certfile", "", "SSL certification file used to secure client/server communication")
	flag.StringVar(&sdnOpts.RemoteCAFile, "remote-cafile", "", "SSL Certificate Authority file used to secure client/server communication")
}

func handleMainErr(err error) {
	if strings.Compare(os.Getenv("OMALLEY_ENV"), "production") == 0 {
		log.Fatal(err)
		return
	}
	panic(err)
}

func main() {
	abook, err := addrbook.New(*master, *dbPath)
	if err != nil {
		handleMainErr(err)
	}
	go network.RunSDN(nil)
	srv := network.New(*master, *masterAddr, *name, *host, *port, abook, nil)
}
