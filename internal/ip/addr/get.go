package addr

//func Get(name string) (*Address, error) {
//	link, err := netlink.LinkByName(name)
//	if err != nil {
//		return nil, err
//	}
//
//	res := GetAttr()
//	res.name = link.Attrs().Name
//
//	addr, err := netlink.AddrList(link, netlink.FAMILY_ALL)
//	if err != nil {
//		return nil, err
//	}
//	for i, _ := range addr {
//		res.addr = append(res.addr, *addr[i].IPNet)
//	}
//
//	return res, nil
//}
