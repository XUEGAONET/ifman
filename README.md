# ifman

ifman全称Interface Manager，用于管理Linux系统的接口及其相关功能。

## 模块

目前ifman支持了如下模块，这些模块均为独立的，在运行时通过控制台参数选择要使用的模块，具体的参数请通过参数`-h`查看。

* **test** - 配置测试模块，用于检查新配置是否能够正常解析。该模块仅仅检查语法错误，无法检查值是否合法，无法检查是否漏掉了某个值
* **core** - 核心服务模块，用于支持接口配置、核查和自动修复。该模块建议由systemd托管运行
* **key** - WireGuard的Key Chain生成器模块。该模块会生成两个KeyChain，在客户端和服务端的WireGuard PtP接口上任意配置其一即可，无需再单独生成私钥公钥去交叉配置
* **reload** - 配置重载通知模块。该模块会向已经运行的ifman进程发送`SIGUSR1`信号，通知其刷新配置文件

## 注意

请阅读注意内容，再使用该工具。

### Issue相关

**开Issue提bug时请务必提交日志，无日志无法处理~**

### 实例相关

**请勿同时运行多个ifman。**

**由于ifman启动时会生成固定位置的pid文件（/var/run/ifman.pid），多个实例启动时可能会相互覆盖，导致问题。**

### 配置核查

配置核查仅支持常规配置修复，主要为如下：

* 状态（UP or DOWN）
* 发送队列长度
* MTU
* 主接口名
* MAC地址

其中，Layer3接口自动忽略MAC地址属性（Unmanaged归属于Layer3接口），Layer2接口全部支持。WireGuard接口核查

### 双栈支持相关

* 现有代码仅完成了IPv4地址族特性的验证，IPv6目前还未验证。如有问题可以先开Issue，但是会等IPv6排期时再统一处理

### 配置相关

* 请留意，不要在配置文件后出现多余的空格，暂时不确定会不会有潜在的问题
* 地址以接口为单位，核查时按照每个接口依次检查和修复
* 配置reload仅支持interface、addr、rp_filter、learning部分重载，其他部分暂不支持
* 配置reload生效需要等待下一次核查时生效
* ifman不存在删除接口的调用。因此，请注意，假设配置中原接口名为vxlan0，其VNI配置为10，当修改该名称为vxlan1之后（其他配置不变），那么ifman核查时会认为vxlan1接口不存在，需要新建，而不会再去管原本的vxlan0接口。这时，由于系统中已经存在同样VNI的接口，因此vxlan1会创建失败

* 当接口不存在时，新建接口操作后还会追加一次强制配置核查，以检查额外配置项并进行修复。好比WireGuard接口，New的时候是仅创建了接口，但是Peer并没有创建，直到强制执行配置核查的时，Update检查到没有Peer配置，才会进行Peer安装
* 如果使用了KeyChain（即KeyChain不为空字符串），即便是指定了Public Key和Private Key，也仍然会去解析KeyChain并使用解析结果覆盖Public Key和Private Key

### 开发相关

* 如添加新的Link支持，在完成NewLink()和UpdateLink()更新后，请记得同时改动如下部分：
    * module_core.go - getLinkFromYaml()
    * link_type.go - getLinkType()

## 使用指引

### 已支持的接口类型

接口配置需要填写在配置文件的interface块中

#### 公共配置

公共配置即所有类型的接口都拥有的配置字段

```yaml
    name: dummy0
    link_up: true
    tx_queue_len: 1024
    mtu: 1500
    master_name: ""
    mac: 94:94:26:a7:c8:11
```

| 字段名称     | 内容释义     | 数据类型 | 其他                                                         |
| ------------ | ------------ | -------- | ------------------------------------------------------------ |
| name         | 接口名       | string   | 接口名，请勿重复。不能为空                                   |
| link_up      | 接口状态     | bool     | 接口UP（true）或者DOWN（false）                              |
| tx_queue_len | 发送队列长度 | uint16   | 建议根据接口业务类型调整，吞吐量优先可加大到4096，转发速度优先可减小到128，或者再适当调小（暂时没测试最小值）。为空（即为0）则由系统自动设置 |
| mtu          | MTU          | uint16   | 为空（即为0）则由系统自动设置                                |
| master_name  | 父接口名     | string   | 即接口的Parent接口名。为空需要使用""，即为无父接口           |
| mac          | MAC地址      | string   | 该项目在Unmanaged和Layer3类型的接口中不会生效，会始终忽略。为空则由系统自动生成 |

####  bridge

```yaml
  - type: bridge
    name: bridge0
    link_up: true
    tx_queue_len: 1024
    mtu: 1500
    master_name: ""
    mac: 94:94:26:a7:c8:ee
    multicast_snooping_on: false
    vlan_filtering_on: true
```

| 字段名称              | 内容释义               | 数据类型 | 其他                                                      |
| --------------------- | ---------------------- | -------- | --------------------------------------------------------- |
| multicast_snooping_on | Multicast Snooping开关 | bool     | 会影响性能，仅在有需要的时候打开                          |
| vlan_filtering_on     | VLAN Filter开关        | bool     | 目前Bridge VLAN支持并不是非常完善，仅在有需要的时候再打开 |

####  dummy

一般用来做loopback接口

```yaml
  - type: dummy
    name: dummy0
    link_up: true
    tx_queue_len: 1024
    mtu: 1500
    master_name: ""
    mac: 94:94:26:a7:c8:11
```

#### iptun

即IPIP Tunnel

```yaml
  - type: iptun
    name: ipip0
    link_up: true
    tx_queue_len: 1024
    mtu: 1500
    master_name: ""
    ttl: 16
    tos: 7
    local_ip: 1.1.1.1
    remote_ip: 1.1.1.2

```

| 字段名称  | 内容释义   | 数据类型 | 其他                                                     |
| --------- | ---------- | -------- | -------------------------------------------------------- |
| ttl       | TTL        | uint8    | 可能会和NoPmtuDisc冲突，暂时还未测试                     |
| tos       | TOS        | uint8    |                                                          |
| local_ip  | 本地IP地址 | string   | 仅IP地址，无需填写端口号。请注意iptables放行ip-encap的包 |
| remote_ip | 远端IP地址 | string   | 仅IP地址，无需填写端口号。请注意iptables放行ip-encap的包 |

####  unmanaged

非由ifman创建，但是需要被ifman进行配置核查的接口。名字可能有点奇怪

不论该接口为什么类型，程序始终认为其为Layer3类型

```yaml
  - type: unmanaged
    name: eth0
    link_up: true
    tx_queue_len: 1024
    mtu: 1500
    master_name: ""
```

####  tun

```yaml
  - type: tun
    name: tun0
    link_up: true
    tx_queue_len: 1024
    mtu: 1500
    master_name: ""
    multi_queue_on: true
    persist_on: true
    queues: 8
```

| 字段名称       | 内容释义   | 数据类型 | 其他         |
| -------------- | ---------- | -------- | ------------ |
| multi_queue_on | 多队列支持 | bool     |              |
| persist_on     | 持久模式   | bool     |              |
| queues         | 队列数量   | uint8    | 暂未测试范围 |

#### vlan

```yaml
  - type: vlan
    name: vlan0
    link_up: true
    tx_queue_len: 1024
    mtu: 1500
    master_name: ""
    mac: 94:94:26:a7:c8:11
    vlan_id: 10
    stacking_on: false
```

| 字段名称    | 内容释义    | 数据类型 | 其他 |
| ----------- | ----------- | -------- | ---- |
| vlan_id     | VLAN ID     | uint16   |      |
| stacking_on | 802.1ad模式 | bool     |      |

#### vrf

```yaml
  - type: vrf
    name: vrf0
    link_up: true
    tx_queue_len: 1024
    mtu: 1500
    master_name: ""
    mac: 94:94:26:a7:c8:11
    table_id: 200
```

| 字段名称 | 内容释义     | 数据类型 | 其他 |
| -------- | ------------ | -------- | ---- |
| table_id | VRF Table ID | uint32   |      |

#### vxlan

```yaml
  - type: vxlan
    name: vxlan0
    link_up: true
    tx_queue_len: 1024
    mtu: 1500
    master_name: ""
    mac: 94:94:26:a7:c8:11
    vni: 1000
    src_ip: 1.1.1.1
    dst_ip: 2.2.2.2
    ttl: 16
    tos: 7
    learning_on: false
    src_port_low: 10240
    src_port_high: 40960
    port: 4789
    vtep_name: eth0
```

| 字段名称      | 内容释义         | 数据类型 | 其他                                                         |
| ------------- | ---------------- | -------- | ------------------------------------------------------------ |
| vni           | VXLAN ID         | uint32   |                                                              |
| src_ip        | 本地IP地址       | string   | 可与vtep_name二选一填写                                      |
| dst_ip        | 远端IP地址       | string   | 做VTEP用时此项可填写`0.0.0.0/0`                              |
| ttl           | TTL              | uint8    |                                                              |
| tos           | TOS              | uint8    |                                                              |
| learning_on   | MAC Learning开关 | bool     | 只是决定VXLAN接口内置的MAC Learning是否开启，与Bridge Port的学习无关 |
| src_port_low  | 源端口范围       | uint16   | 请注意，当源端口范围过大时，VTEP如果启用了连接跟踪则可能会导致conntrack表爆表的问题。但是端口范围大时可以提高ECMP的随机度 |
| src_port_high | 源端口范围       | uint16   | 同上                                                         |
| port          | 端口号           | uint16   | 端口，建议4789                                               |
| vtep_name     | 绑定的接口名称   | string   | 物理接口名称，即VXLAN接口绑定的上游接口。可与`src_ip`二选一填写。 |

#### wireguard_ptp_server

```yaml
  - type: wireguard_ptp_server
    name: wgs0
    link_up: true
    tx_queue_len: 1024
    mtu: 1500
    master_name: ""
    listen_port: 10000
    private: sss
    peer_public: ssss
    key_chain: sssssssss
```

| 字段名称    | 内容释义        | 数据类型 | 其他                                           |
| ----------- | --------------- | -------- | ---------------------------------------------- |
| listen_port | 监听端口（UDP） | uint16   |                                                |
| private     | 私钥（本地）    | string   |                                                |
| peer_public | 公钥（远端）    | string   |                                                |
| key_chain   | 密码链          | string   | 当指定该字段时，private与peer_public则自动失效 |

#### wireguard_ptp_client

```yaml
  - type: wireguard_ptp_client
    name: wgc0
    link_up: true
    tx_queue_len: 1024
    mtu: 1500
    master_name: ""
    endpoint: 1.1.1.1:6666
    heartbeat_interval: 10
    private: sssss
    peer_public: ssssss
    key_chain: ssssssssssssss
```

| 字段名称           | 内容释义       | 数据类型 | 其他                                           |
| ------------------ | -------------- | -------- | ---------------------------------------------- |
| endpoint           | 目的地址       | string   | 带端口号，如1.1.1.1:3333                       |
| heartbeat_interval | 心跳周期（秒） | uint32   | 用于保持长连接，存在NAT时建议设置为5s或者更短  |
| private            | 私钥（本地）   | string   |                                                |
| peer_public        | 公钥（远端）   | string   |                                                |
| key_chain          | 密码链         | string   | 当指定该字段时，private与peer_public则自动失效 |





