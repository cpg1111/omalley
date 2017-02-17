package network

import (
	"context"
	"errors"
	"strings"

	fnet "github.com/coreos/flannel/network"
	"github.com/coreos/flannel/remote"
	"github.com/coreos/flannel/subnet"

	_ "github.com/coreos/flannel/backend/vxlan"
)

type SDNOpts struct {
	EtcdEndpoints  string
	EtcdPrefix     string
	EtcdKeyFile    string
	EtcdCertFile   string
	EtcdCAFile     string
	EtcdUsername   string
	EtcdPassword   string
	Help           bool
	Version        bool
	Listen         string
	Remote         string
	RemoteKeyFile  string
	RemoteCertFile string
	RemoteCAFile   string
}

func getSubnetManager(opts *SDNOpts) (subnet.Manager, error) {
	if opts.Remote != "" {
		return remote.NewRemoteManager(opts.Remote, opts.RemoteCAFile, opts.RemoteCertFile, opts.RemoteKeyFile)
	}
	cfg := &subnet.EtcdConfig{
		Endpoints: strings.Split(opts.EtcdEndpoints, ","),
		Keyfile:   opts.EtcdKeyFile,
		Certfile:  opts.EtcdCertFile,
		CAFile:    opts.EtcdCAFile,
		Prefix:    opts.EtcdPrefix,
		Username:  opts.EtcdUsername,
		Password:  opts.EtcdPassword,
	}
	return subnet.NewLocalManager(cfg)
}

func RunSDN(opts *SDNOpts) error {
	subnetManager, err := getSubnetManager(opts)
	if err != nil {
		return err
	}
	ctx := context.Background()
	if opts.Listen != "" {
		if opts.Remote != "" {
			return errors.New("SDN can only listen or be remote, not both")
		}
		remote.RunServer(ctx, subnetManager, opts.Listen, opts.RemoteCAFile, opts.RemoteCertFile, opts.RemoteKeyFile)
	} else {
		netManager, err := fnet.NewNetworkManager(ctx, subnetManager)
		if err != nil {
			return err
		}
		netManager.Run(ctx)
	}
	return nil
}

func GetPublicIPAddr() (string, error) {
	return "", nil
}
