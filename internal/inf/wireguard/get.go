package wireguard

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/wgctrl"
)

func Get(name string) (*WireGuard, error) {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return nil, err
	}

	d, ok := link.(netlink.Link)
	if !ok {
		return nil, fmt.Errorf("LinkByName type asserting failed")
	}

	res := GetAttr()
	res.name = d.Attrs().Name
	res.mtu = d.Attrs().MTU
	res.txQueueLen = d.Attrs().TxQLen
	res.masterId = d.Attrs().MasterIndex

	cli, err := wgctrl.New()
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	dev, err := cli.Device(res.name)
	if err != nil {
		return nil, err
	}

	res.listenPort = dev.ListenPort
	res.private = &dev.PrivateKey

	if len(dev.Peers) == 0 {
		return res, nil
	}
	if len(dev.Peers) > 1 {
		logrus.WithField("module", "wireguard").
			Warningf("more than one peer under wireguard interface %s", res.name)
	}
	res.endpoint = dev.Peers[0].Endpoint
	res.hsInterval = &dev.Peers[0].PersistentKeepaliveInterval
	res.peerPublic = &dev.Peers[0].PublicKey

	return res, nil
}
