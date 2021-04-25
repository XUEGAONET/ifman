package wireguard

import (
	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"net"
)

func Update(wg *WireGuard) error {
	err := wg.check()
	if err != nil {
		return err
	}

	link, err := netlink.LinkByName(wg.name)
	if err != nil {
		return err
	}

	err = netlink.LinkSetMTU(link, wg.mtu)
	if err != nil {
		return err
	}

	err = netlink.LinkSetTxQLen(link, wg.txQueueLen)
	if err != nil {
		return err
	}

	err = netlink.LinkSetMasterByIndex(link, wg.masterId)
	if err != nil {
		return err
	}

	cli, err := wgctrl.New()
	if err != nil {
		return err
	}
	defer cli.Close()

	_, allow, err := net.ParseCIDR("0.0.0.0/0")
	if err != nil {
		return err
	}

	wgPeer := wgtypes.PeerConfig{
		PublicKey:                   *wg.peerPublic,
		Remove:                      false,
		UpdateOnly:                  false,
		PresharedKey:                nil,
		Endpoint:                    wg.endpoint,
		PersistentKeepaliveInterval: wg.hsInterval,
		ReplaceAllowedIPs:           true,
		AllowedIPs:                  []net.IPNet{*allow},
	}

	wgConf := wgtypes.Config{
		PrivateKey:   wg.private,
		ListenPort:   nil,
		FirewallMark: nil,
		ReplacePeers: true,
		Peers:        []wgtypes.PeerConfig{wgPeer},
	}

	if wg.listenPort != 0 {
		wgConf.ListenPort = &wg.listenPort
	}

	err = cli.ConfigureDevice(wg.name, wgConf)
	return err
}
