---
title: 万字带你深度了解组件 etcd
author: NuyoahCh
date: 2025-04-01 15:27:56
description: 深入探索 etcd 奥妙
tags: [etcd] ##标签
categories:
 - [万字详解系列文章分享〄]
 - [etcd] ##目录
top: false
# cover:
---

近几年，云原生越来越火，你在各种大会或博客的标题里都可以见到 “云原生” 的字样，我们这次要学习的 etcd 也是云原生架构中重要的基础组件，因为 etcd 项目是 Kubernetes 内部的一大关键组件，目前有很多项目都依赖 etcd 进行可靠的分布式数据存储。

etcd 是 CoreOS 团队于 2013 年 6 月发起的开源项目，2018 年底正式加入云原生计算基金会（CNCF）。etcd 组件基于 Go 语言实现，目前最新版本为 V3.4.9。

### 为什么需要 etcd

在具体讲解 etcd 前，我们还是先谈谈分布式系统存在的问题。

从本质上来讲，云原生中的微服务应用属于分布式系统的一种落地实践。在分布式环境中，由于网络的复杂性、不确定性以及节点故障等情况，会产生一系列的问题。最常见的、最大的难点就是**数据存储不一致**的问题，即多个服务实例自身的数据或者获取到的数据各不相同。因此我们需要基于一致性的存储组件构建可靠的分布式系统。

#### 分布式中的 CAP 理论

CAP 原理是描述分布式系统下节点数据同步的基本定理，分别指 **Consistency（一致性）、Availability（可用性）和 Partition tolerance（分区容错性）**，这三个要素最多只能同时实现两点，不能三者兼顾。

基于分布式系统的基本特质，P（分区容错性）是必须要满足的，所以接下来需要考虑满足 C（一致性）还是 A（可用性）。

在类似银行之类对金额数据要求强一致性的系统中，要优先考虑满足数据一致性；而在大众网页之类的系统中，用户对网页版本的新旧不会有特别的要求，在这种场景下服务可用性会高于数据一致性。

了解了分布式系统中的问题，接下来让我们结合官网中的定义，看看为什么在分布式系统中需要 etcd？

#### etcd 是什么

根据 [etcd 官网](https://etcd.io/)的介绍，我找到了如下定义：

\> A highly-available key value store for shared configuration and service discovery.
\> 即一个用于配置共享和服务发现的键值存储系统。

从定义上你也可以发现，etcd 归根结底是一个存储组件，且可以实现配置共享和服务发现。

在分布式系统中，各种服务配置信息的管理共享和服务发现是一个很基本也是很重要的问题，无论你调用服务还是调度容器，都需要知道对应的服务实例和容器节点地址信息。etcd 就是这样一款**实现了元数据信息可靠存储的组件**。

**etcd 可集中管理配置信息**。服务端将配置信息存储于 etcd，客户端通过 etcd 得到服务配置信息，etcd 监听配置信息的改变，发现改变通知客户端。

而 etcd 满足 CAP 理论中的 CP（一致性和分区容错性） 指标，由此我们知道，etcd 解决了分布式系统中一致性存储的问题。

### etcd 中常用的术语

为了我们接下来更好地学习 etcd，我在这里给你列举了常用的 etcd 术语，尽快熟悉它们也会对接下来的学习有所助益。

下面我们具体了解一下 etcd 的相关特性、架构和使用场景。

### etcd 的特性

etcd 可以用来构建高可用的分布式键值数据库，总结来说有如下特点。

-   **简单**：etcd 的安装简单，且为用户提供了 HTTP API，使用起来也很简单。
-   **存储**：etcd 的基本功能，数据分层存储在文件目录中，类似于我们日常使用的文件系统。
-   **Watch 机制**：Watch 指定的键、前缀目录的更改，并对更改时间进行通知。
-   **安全通信**：支持 SSL 证书验证。
-   **高性能**：etcd 单实例可以支持 2K/s 读操作，官方也有提供基准测试脚本。
-   **一致可靠**：基于 Raft 共识算法，实现分布式系统内部数据存储、服务调用的一致性和高可用性。

etcd 是一个实现了分布式一致性键值对存储的中间件，支持跨平台，拥有活跃用户的技术社区。etcd 集群中的节点基于 Raft 算法进行通信，Raft 算法保证了微服务实例或机器集群所访问的数据的可靠一致性。

在分布式系统或者 Kubernetes 集群中，etcd 可以**作为服务注册与发现和键值对存储组件**。不管是简单应用程序，还是复杂的容器集群，都可以很方便地从 etcd 中读取数据，满足了各种场景的需求。

### etcd 的应用场景

etcd 在稳定性、可靠性和可伸缩性上表现极佳，同时也为云原生应用系统提供了协调机制。etcd 经常用于服务注册与发现的场景，此外还有键值对存储、消息发布与订阅、分布式锁等场景。

-   **键值对存储**

etcd 是一个用于**键值存储**的组件，存储是 etcd 最基本的功能，其他应用场景都建立在 etcd 的可靠存储上。比如 Kubernetes 将一些元数据存储在 etcd 中，将存储状态数据的复杂工作交给 etcd，Kubernetes 自身的功能和架构就能更加稳定。

etcd 基于 Raft 算法，能够有力地保证分布式场景中的一致性。各个服务启动时注册到 etcd 上，同时为这些服务配置键的 TTL 时间。注册到 etcd 上面的各个服务实例通过心跳的方式定期续租，实现服务实例的状态监控。

-   **消息发布与订阅**

在分布式系统中，服务之间还可以通过消息通信，即消息的发布与订阅，如下图所示：

消息发布与订阅流程图

通过构建 etcd 消息中间件，服务提供者发布对应主题的消息，消费者则订阅他们关心的主题，一旦对应的主题有消息发布，就会产生订阅事件，消息中间件就会通知该主题所有的订阅者。

-   **分布式锁**

分布式系统中涉及多个服务实例，存在跨进程之间资源调用，对于资源的协调分配，单体架构中的锁已经无法满足需要，需要引入分布式锁的概念。etcd 基于 Raft 算法，实现分布式集群的一致性，存储到 etcd 集群中的值必然是全局一致的，因此基于 etcd 很容易实现分布式锁。

### etcd 的核心架构

etcd 作为一个如此重要的部件，我们只有深入理解其架构设计才能更好地学习。下面还是先来看看 etcd 总体的架构图。

etcd 总体架构图

从上图可知，etcd 有 etcd Server、gRPC Server、存储相关的 MVCC 、Snapshot、WAL，以及 Raft 模块。

其中：

-   etcd Server 用于对外接收和处理客户端的请求；
-   gRPC Server 则是 etcd 与其他 etcd 节点之间的通信和信息同步；
-   MVCC，即多版本控制，etcd 的存储模块，键值对的每一次操作行为都会被记录存储，这些数据底层存储在 BoltDB 数据库中；
-   WAL，预写式日志，etcd 中的数据提交前都会记录到日志；
-   Snapshot 快照，以防 WAL 日志过多，用于存储某一时刻 etcd 的所有数据；
-   Snapshot 和 WAL 相结合，etcd 可以有效地进行数据存储和节点故障恢复等操作。

虽然 etcd 内部实现机制复杂，但对外提供了简单的 API 接口，方便客户端调用。我们可以通过 **etcdctl 客户端命令行**操作和访问 etcd 中的数据，或者通过 **HTTP API** 接口直接访问 etcd。

etcd 中的数据结构很简单，它的数据存储其实就是键值对的有序映射。etcd 还提供了一种键值对监测机制，即 Watch 机制，客户端通过订阅相关的键值对，获取其更改的事件信息。Watch 机制实时获取 etcd 中的增量数据更新，使数据与 etcd 同步。

etcd 目前有 V2.x 和 V3.x 两个大版本。etcd V2 和 V3 是在底层使用同一套 Raft 算法的两个独立应用，但相互之间实现原理和使用方法上差别很大，接口不一样、存储不一样，两个版本的数据互相隔离。

至于由 etcd V2 升级到 etcd V3 的情况，原有数据只能通过 etcd V2 接口访问，V3 接口创建的数据只能通过新的 V3 的接口访问。我们的专栏重点讲解**当前常用且主流的 V3 版本**。

### 小结

这一讲我主要介绍了 etcd 相关的概念。关于 etcd 你需要记住以下三点：

-   etcd 是云原生架构中的存储基石，可以有效保证存储数据的一致性和可靠性；
-   etcd 内部实现机制复杂，但是对外提供了简单直接的 API 接口；
-   使用 etcd 的常见分布式场景包括键值对存储、服务注册与发现、消息订阅与发布、分布式锁等。

### etcd 单机安装部署

etcd 的安装有多种方式，这里我以 CentOS 7 为例，可以通过`yum install etcd`进行安装。然而通过系统工具安装的 etcd 版本比较滞后，如果需要安装最新版本的 etcd ，我们可以通过二进制包、源码编译以及 Docker 容器安装。

#### 二进制安装

目前最新的 etcd API 版本为 v3.4，我们基于 3.4.4 版本进行实践，API 版本与最新版保持一致，在 CentOS 7 上面使用如下脚本进行安装：

```
ETCD_VER=v3.4.4


GITHUB_URL=https:


DOWNLOAD_URL=${GITHUB_URL}


rm -f /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz


rm -rf /tmp/etcd-download-test && mkdir -p /tmp/etcd-download-test


curl -L ${DOWNLOAD_URL}/${ETCD_VER}/etcd-${ETCD_VER}-linux-amd64.tar.gz -o /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz


tar xzvf /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz -C /tmp/etcd-download-test --strip-components=1


rm -f /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz


/tmp/etcd-download-test/etcd --version


/tmp/etcd-download-test/etcdctl version
```

下载可能比较慢，执行后，查看 etcd 版本的结果如下：

```
etcd Version: 3.4.4


Git SHA: e784ba73c


Go Version: go1.12.12


Go OS/Arch: linux/amd64
```

根据上面的执行结果可以看到，我们在 Linux 上安装成功，macOS 的二进制安装也类似，这里不重复演示了。关于 Windows 系统的安装比较简单，下载好安装包后，直接执行。其中 etcd.exe 是服务端，etcdctl.exe 是客户端，如下图所示：

#### 源码安装

使用源码安装，首先你需要确保本地的 Go 语言环境。如未安装，请参考 https://golang.org/doc/install 安装 Go 语言环境。我们需要 Go 版本为 1.13+，来构建最新版本的 etcd。如果你想尝试最新版本，也可以从 master 分支构建 etcd。

首先查看一下我们本地的 Go 版本：

```
$ go version


Go version go1.13.6 linux/amd64
```

经检测，本地的 Go 版本满足要求。将制定版本的 etcd 项目 clone 到本地之后，在项目文件夹构建 etcd。构建好之后，执行测试命令，确保 etcd 编译安装成功：

```
$ ./etcdctl version


etcdctl version: 3.4.4


API version: 3.4
```

经过上述步骤，我们已经通过源码编译成功安装 etcd。

除此之外，还可以使用 Docker 容器安装，这种方式更为简单，但这种方式有一个弊端：启动时会将 etcd 的端口暴露出来。

### etcd 集群安装部署

刚刚我们介绍了 etcd 的单机安装，但是**在实际生产环境中，为了整个集群的高可用，etcd 通常都会以集群方式部署，以避免单点故障**。下面我们看看如何进行 etcd 集群部署。

引导 etcd 集群的启动有以下三种方式：

-   静态启动
-   etcd 动态发现
-   DNS 发现

其中静态启动 etcd 集群的方式要求每个成员都知道集群中的其他成员。然而很多场景中，群集成员的 IP 可能未知。因此需要借助发现服务引导 etcd 群集启动。下面我们将会分别介绍这几种方式。

#### 静态方式启动 etcd 集群

如果我们想在一台机器上实践 etcd 集群的搭建，可以通过 goreman 工具。

goreman 是一个 Go 语言编写的多进程管理工具，是对 Ruby 下广泛使用的 Foreman 的重写。下面我们使用 goreman 来演示在一台物理机上面以静态方式启动 etcd 集群。

前面我们已经确认过 Go 语言的安装环境，现在可以直接执行：

```
go get github.com/mattn/goreman
```

编译后的文件放在`$GOPATH/bin`中，`$GOPATH/bin`目录已经添加到了系统`$PATH`中，所以我们可以方便地执行命令`goreman`命令。

下面就是编写 Procfile 脚本，我们启动三个 etcd，具体对应如下：

Procfile 脚本中 infra1 启动命令如下：

```
etcd1: etcd --name infra1 --listen-client-urls http:
```

infra2 和 infra3 的启动命令类似。下面我们看一下其中各配置项的说明。

-   --name：etcd 集群中的节点名，这里可以随意，方便区分且不重复即可。
-   --listen-client-urls：监听用于客户端通信的 url，同样可以监听多个。
-   --advertise-client-urls：建议使用的客户端通信 url，该值用于 etcd 代理或 etcd 成员与 etcd 节点通信。
-   --listen-peer-urls：监听用于节点之间通信的 url，可监听多个，集群内部将通过这些 url 进行数据交互 (如选举、数据同步等)。
-   --initial-advertise-peer-urls：建议用于节点之间通信的 url，节点间将以该值进行通信。
-   --initial-cluster-token： etcd-cluster-1，节点的 token 值，设置该值后集群将生成唯一 ID，并为每个节点也生成唯一 ID。当使用相同配置文件再启动一个集群时，只要该 token 值不一样，etcd 集群就不会相互影响。
-   --initial-cluster：集群中所有的 initial-advertise-peer-urls 的合集。
-   --initial-cluster-state：new，新建集群的标志。

注意上面的脚本，**etcd 命令执行时需要根据本地实际的安装地址进行配置**。下面我们使用`goreman`命令启动 etcd 集群：

```
goreman -f /opt/procfile start
```

启动完成后查看集群内的成员：

```
$ etcdctl --endpoints=http:


8211f1d0f64f3269, started, infra1, http:


91bc3c398fb3c146, started, infra2, http:


fd422379fda50e48, started, infra3, http:
```

现在我们的 etcd 集群已经搭建成功。需要注意的是：**在集群启动时，我们通过静态的方式指定集群的成员。但在实际环境中，集群成员的 IP 可能不会提前知道，此时需要采用动态发现的机制**。

#### 动态发现启动 etcd 集群

Discovery Service，即发现服务，帮助新的 etcd 成员使用共享 URL 在集群引导阶段发现所有其他成员。

Discovery Service 使用已有的 etcd 集群来协调新集群的启动，主要操作如下：

-   首先，所有新成员都与发现服务交互，并帮助生成预期的成员列表；
-   之后，每个新成员使用此列表引导其服务器，该列表执行与`--initial-cluster`标志相同的功能，即设置所有集群的成员信息。

我们的实验中启动三个 etcd 实例，具体对应如下：

下面就开始以动态发现的方式来启动 etcd 集群，具体步骤如下。

**获取 discovery 的 token**

**首先需要生成标识新集群的唯一令牌**。该令牌将用于键空间中的唯一前缀，比较简单的方法是使用 uuidgen 生成 UUID：

**指定集群的大小**

获取令牌时，必须指定群集大小。**发现服务使用该数值来确定组成集群的所有成员**：

我们需要把该 url 作为`--discovery`参数来启动 etcd，节点会自动使用该路径对应的目录进行 etcd 的服务注册与发现。

**公共发现服务**

当我们本地没有可用的 etcd 集群，etcd 官网提供了一个可以公网访问的 etcd 存储地址。你可以通过如下命令得到 etcd 服务的目录，并把它作为`--discovery`参数使用。

公共发现服务`discovery.etcd.io`以相同的方式工作，但是有一层修饰，在此之上仍然使用 etcd 群集作为数据存储。

**以动态发现方式启动集群**

etcd 发现模式下，启动 etcd 的命令如下：

```
# etcd1 启动


$ /opt/etcd/bin/etcd  --name etcd1 --initial-advertise-peer-urls http://192.168.202.128:2380 \


  --listen-peer-urls http://192.168.202.128:2380 \


  --data-dir /opt/etcd/data \


  --listen-client-urls http://192.168.202.128:2379,http://127.0.0.1:2379 \


  --advertise-client-urls http://192.168.202.128:2379 \


  --discovery https://discovery.etcd.io/3e86b59982e49066c5d813af1c2e2579cbf573de
```

etcd2 和 etcd3 启动类似，替换 listen-peer-urls 和 advertise-client-urls 即可。需要注意的是：**在我们完成了集群的初始化后，**`--discovery`就失去了作用。当需要增加节点时，需要使用 etcdctl 进行操作。

为了安全，每次启动新 etcd 集群时，都会使用新的 discovery token 进行注册。当出现初始化启动的节点超过了指定的数量时，多余的节点会自动转化为 Proxy 模式的 etcd（在后面课时会详细介绍）。

**结果验证**

集群启动后，进行验证，我们看一下集群节点的健康状态：

```
$ /opt/etcd/bin/etcdctl  --endpoints=&quot;http://192.168.202.128:2379,http://192.168.202.129:2379,http://192.168.202.130:2379&quot;  endpoint  health


# 结果如下


    http://192.168.202.128:2379 is healthy: successfully committed proposal: took = 3.157068ms


    http://192.168.202.130:2379 is healthy: successfully committed proposal: took = 3.300984ms


    http://192.168.202.129:2379 is healthy: successfully committed proposal: took = 3.263923ms
```

可以看到，集群中的三个节点都是健康的正常状态，可以说明以动态发现方式启动集群成功。

除此之外，**etcd 还支持使用 DNS SRV 记录启动 etcd 集群**。使用 DNSmasq 创建 DNS 服务，实际上是利用 DNS 的 SRV 记录不断轮训查询实现，DNS SRV 是 DNS 数据库中支持的一种资源记录的类型，它记录了计算机与所提供服务信息的对应关系。

至此，我们就介绍完了 etcd 安装部署的两种方式，对于可靠的系统来说，还需要考虑 etcd 集群通信的安全性，为我们的数据安全增加防护。

### etcd 通信安全

etcd 支持通过 TLS 协议进行的加密通信，TLS 通道可用于对等体（指 etcd 集群中的服务实例）之间的加密内部群集通信以及加密的客户端流量。这一节我们看一下客户端 TLS 设置群集的实现。

想要实现数据 HTTPS 加密协议访问、保障数据的安全，就需要 SSL 证书，TLS 是 SSL 与 HTTPS 安全传输层协议名称。

#### 进行 TLS 加密实践

为了进行实践，我们将安装一些实用的命令行工具，这包括 CFSSL 、cfssljson。

CFSSL 是 CloudFlare 的 PKI/TLS 加密利器。它既是命令行工具，也可以用于签名、验证和捆绑 TLS 证书的 HTTP API 服务器，需要 Go 1.12+ 版本才能构建。

**环境配置**

**安装 CFSSL**

```
$ ls ~/Downloads/cfssl


cfssl-certinfo_1.4.1_linux_amd64 cfssl_1.4.1_linux_amd64          cfssljson_1.4.1_linux_amd64


chmod +x cfssl_1.4.1_linux_amd64 cfssljson_1.4.1_linux_amd64 cfssl-certinfo_1.4.1_linux_amd64


mv cfssl_1.4.1_linux_amd64 /usr/local/bin/cfssl


mv cfssljson_1.4.1_linux_amd64 /usr/local/bin/cfssljson


mv cfssl-certinfo_1.4.1_linux_amd64 /usr/bin/cfssl-certinfo
```

安装完成后，查看版本信息如下所示：

```
$ cfssl version


Version: 1.4.1


Runtime: go1.12.12
```

**配置 CA 并创建 TLS 证书**

我们将使用 CloudFlare's PKI 工具 CFSSL 来配置 PKI 安全策略，然后用它创建 Certificate Authority（CA，即证书机构），并为 etcd 创建 TLS 证书。

首先创建 SSL 配置目录：

```
mkdir /opt/etcd/{bin,cfg,ssl} -p


cd /opt/etcd/ssl/
```

接下来完善 etcd CA 配置，写入 ca-config.json 如下的配置：

```
cat &lt;&lt; EOF | tee ca-config.json


{


&quot;signing&quot;: {


&quot;default&quot;: {


&quot;expiry&quot;: &quot;87600h&quot;


    },


&quot;profiles&quot;: {


&quot;etcd&quot;: {


&quot;expiry&quot;: &quot;87600h&quot;,


&quot;usages&quot;: [


&quot;signing&quot;,


&quot;key encipherment&quot;,


&quot;server auth&quot;,


&quot;client auth&quot;


        ]


      }


    }


  }


}


EOF
```

生成获取 etcd ca 证书，需要证书签名的请求文件，因此在 ca-csr.json 写入如下的配置：

```
cat &lt;&lt; EOF | tee ca-csr.json


{


&quot;CN&quot;: &quot;etcd CA&quot;,


&quot;key&quot;: {


&quot;algo&quot;: &quot;rsa&quot;,


&quot;size&quot;: 2048


    },


&quot;names&quot;: [


        {


&quot;C&quot;: &quot;CN&quot;,


&quot;L&quot;: &quot;Shanghai&quot;,


&quot;ST&quot;: &quot;Shanghai&quot;


        }


    ]


}


EOF
```

根据上面的配置，生成 CA 凭证和私钥：

```
$ cfssl gencert -initca ca-csr.json | cfssljson -bare ca
```

生成 etcd server 证书，写入 server-csr.json 如下的配置：

```
cat &lt;&lt; EOF | tee server-csr.json


{


&quot;CN&quot;: &quot;etcd&quot;,


&quot;hosts&quot;: [


&quot;192.168.202.128&quot;,


&quot;192.168.202.129&quot;,


&quot;192.168.202.130&quot;


    ],


&quot;key&quot;: {


&quot;algo&quot;: &quot;rsa&quot;,


&quot;size&quot;: 2048


    },


&quot;names&quot;: [


        {


&quot;C&quot;: &quot;CN&quot;,


&quot;L&quot;: &quot;Beijing&quot;,


&quot;ST&quot;: &quot;Beijing&quot;


        }


    ]


}


EOF
```

最后就可以生成 server 证书了：

```
cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=etcd server-csr.json | cfssljson -bare server
```

启动 etcd 集群，配置如下：

```
#etcd1 启动


$ /opt/etcd/bin/etcd --name etcd1 --initial-advertise-peer-urls https://192.168.202.128:2380 \


     --listen-peer-urls https://192.168.202.128:2380 \


     --listen-client-urls https://192.168.202.128:2379,https://127.0.0.1:2379 \


     --advertise-client-urls https://192.168.202.128:2379 \


     --initial-cluster-token etcd-cluster-1 \


     --initial-cluster etcd1=https://192.168.202.128:2380, etcd2=https://192.168.202.129:2380, etcd3=https://192.168.202.130:2380 \


     --initial-cluster-state new \


     --client-cert-auth --trusted-ca-file=/opt/etcd/ssl/ca.pem \


     --cert-file=/opt/etcd/ssl/server.pem --key-file=/opt/etcd/ssl/server-key.pem \


     --peer-client-cert-auth --peer-trusted-ca-file=/opt/etcd/ssl/ca.pem \


     --peer-cert-file=/opt/etcd/ssl/server.pem --peer-key-file=/opt/etcd/ssl/server-key.pem
```

etcd2 和 etcd3 启动类似，**注意替换 listen-peer-urls 和 advertise-client-urls**。通过三台服务器的控制台可以知道，集群已经成功建立。

下面我们进行验证：

```
$ /opt/etcd/bin/etcdctl --cacert=/opt/etcd/ssl/ca.pem --cert=/opt/etcd/ssl/server.pem --key=/opt/etcd/ssl/server-key.pem --endpoints=&quot;https://192.168.202.128:2379,https://192.168.202.129:2379,https://192.168.202.130:2379&quot;  endpoint health


# 输出如下：


https://192.168.202.129:2379 is healthy: successfully committed proposal: took = 9.492956ms


https://192.168.202.130:2379 is healthy: successfully committed proposal: took = 12.805109ms


https://192.168.202.128:2379 is healthy: successfully committed proposal: took = 13.036091ms
```

查看三个节点的健康状况，`endpoint health`，输出的结果符合我们的预期。

经过 TLS 加密的 etcd 集群，在进行操作时，需要加上认证相关的信息。

**自动证书**

如果集群需要加密的通信但不需要经过身份验证的连接，可以将 etcd 配置为自动生成其密钥。在初始化时，每个成员都基于其通告的 IP 地址和主机创建自己的密钥集。

在每台机器上，etcd 将使用以下标志启动：

```
$ etcd --name etcd1 --initial-advertise-peer-urls https:


  --listen-peer-urls https:


  --listen-client-urls https:


  --advertise-client-urls https:


  --initial-cluster-token etcd-cluster-1 \


  --initial-cluster infra0=https:


  --initial-cluster-state new \


  --auto-tls \


  --peer-auto-tls
```

由于自动签发证书并不能认证身份，直接 curl 会返回错误，**需要使用 curl 的`-k`命令屏蔽对证书链的校验**。

### 小结

这一讲我们主要介绍了 etcd 的安装部署方式，包括单机和集群的部署。为了高可用，在生产环境中我们一般选择使用集群的方式部署 etcd。在集群部署的时候，有静态和动态两种方式发现集群中的成员。静态的方式是直接在启动时指定各个成员的地址，但在实际环境中，集群成员的 IP 可能不会提前知道，这时候就需要采用动态发现的机制。

Discovery Service 用于生成集群的发现令牌，需要注意的是，**该令牌仅用于集群引导阶段，不能用于运行时重新配置或集群监视**。一个发现令牌只能代表一个 etcd 集群，只要此令牌上的发现协议启动，即使它中途失败，也不能用于引导另一个 etcd 集群。我在这一讲最后也介绍了 etcd 集群通过 TLS 协议进行的加密通信，来保证 etcd 通信的安全。



### etcdctl 客户端

etcdctl 是一个命令行客户端，便于我们进行**服务测试或手动修改数据库内容**，我们刚开始熟悉 etcd 功能时可以通过 etdctl 客户端熟悉相关操作。etcdctl 在两个不同的 etcd 版本（v2 和 v3）下的功能和使用方式也完全不同。一般通过如下方式来指定使用 etcd 的版本：

```
export ETCDCTL_API=2


export ETCDCTL_API=3
```

我们的专栏课程主要讲解 API 3。etcd 项目二进制发行包中已经包含了 etcdctl 工具，通过 etcd 安装包中的 etcdctl 可执行文件可以进行调用。下面我们先来看看 etcd 的常用命令有哪些，并进行实践应用。

### 常用命令介绍

我们首先来看下 etcdctl 支持哪些命令，通过`etcdctl -h`命令查看：

```
$ etcdctl -h


NAME:


	etcdctl - A simple command line client for etcd3.


USAGE:


	etcdctl [flags]


VERSION:


3.4.4


API VERSION:


3.4
```

COMMANDS:

OPTIONS:

etcdctl 支持的命令大体上分为**数据库操作和非数据库操作**两类。其中数据库的操作命令是最常用的命令，我们将在下面具体介绍。其他的命令如用户、角色、授权、认证相关，你可以根据语法自己尝试一下。

### 数据库操作

数据库操作基本围绕着对键值和目录的 CRUD 操作（即增删改查），及其对应的生命周期管理。我们上手这些操作其实很方便，因为这些操作是符合 REST 风格的一套 API 操作。

etcd 在键的组织上采用了类似文件系统中目录的概念，即**层次化的空间结构**，我们指定的键可以作为键名，如：testkey，实际上，此时键值对放于根目录 / 下面。我们也可以为键的存储**指定目录结构**，如 /cluster/node/key；如果不存在 /cluster/node 目录，则 etcd Server 将会创建相应的目录结构。

下面我们基于键操作、watch、lease 三类分别介绍 etcdctl 的使用与实践。

#### 键操作

键操作包括最常用的增删改查操作，包括 PUT、GET、DELETE 等命令。

**PUT 设置或者更新某个键的值**。例如：

```
$ etcdctl put /test/foo1 &quot;Hello world&quot;


$ etcdctl put /test/foo2 &quot;Hello world2&quot;


$ etcdctl put /test/foo3 &quot;Hello world3&quot;
```

成功写入三对键值，/test/foo1、/test/foo2 和 /test/foo3。

**GET 获取指定键的值**。例如获取 /testdir/testkey 对应的值：

```
$ etcdctl get /testdir/testkey


Hello world
```

除此之外， etcdctl 的 GET 命令还提供了根据指定的键（key），获取其对应的十六进制格式值，即以十六进制格式返回：

```
$ etcdctl get /test/foo1 --hex


\x2f\x74\x65\x73\x74\x64\x69\x72\x2f\x74\x65\x73\x74\x6b\x65\x79 #键


\x48\x65\x6c\x6c\x6f\x20\x77\x6f\x72\x6c\x64 #值
```

加上`--print-value-only`可以读取对应的值。十六进制在 etcd 中有多处使用，如**租约 ID** 也是十六进制。

GET 范围内的值：

```
$ etcdctl get /test/foo1 /test/foo3


/test/foo1


Hello world


/test/foo2


Hello world2
```

可以看到，上述操作获取了大于等于 /test/foo1，且小于 /test/foo3 的键值对。foo3 不在范围之内，因为范围是半开区间 [foo1, foo3)，不包含 foo3。

获取某个前缀的所有键值对，通过 --prefix 可以指定前缀：

```
$ etcdctl get --prefix /test/foo


/test/foo1


Hello world


/test/foo2


Hello world2


/test/foo3


Hello world3
```

这样就能获取所有以 /test/foo 开头的键值对，当前缀获取的结果过多时，还可以通过 --limit=2 限制获取的数量：

```
etcdctl get --prefix --limit=2 /test/foo
```

读取键过往版本的值，应用可能想读取键的被替代的值。

例如，应用可能想**通过访问键的过往版本回滚到旧的配置**。或者，应用可能想**通过多个请求得到一个覆盖多个键的统一视图**，而这些请求可以通过访问键历史记录而来。因为 etcd 集群上键值存储的每个修改都会增加 etcd 集群的全局修订版本，应用可以通过提供旧有的 etcd 修改版本来读取被替代的键。现有如下键值对：

```
foo = bar         # revision = 2


foo1 = bar2       # revision = 3


foo = bar_new     # revision = 4


foo1 = bar1_new   # revision = 5
```

以下是访问以前版本 key 的示例：

```
$ etcdctl get --prefix foo # 访问最新版本的 key


foo


bar_new


foo1


bar1_new


$ etcdctl get --prefix --rev=4 foo # 访问第 4 个版本的 key


foo


bar_new


foo1


bar1


$ etcdctl get --prefix --rev=3 foo #  访问第 3 个版本的 key


foo


bar


foo1


bar1


$ etcdctl get --prefix --rev=2 foo #  访问第 3 个版本的 key


foo


bar


$ etcdctl get --prefix --rev=1 foo #  访问第 1 个版本的 key
```

应用可能想读取大于等于指定键的 byte 值的键。假设 etcd 集群已经有如下列键：

读取大于等于键 b 的 byte 值的键的命令：

```
$ etcdctl get --from-key b


b


456


z


789
```

**DELETE 键，应用可以从 etcd 集群中删除一个键或者特定范围的键**。

假设 etcd 集群已经有下列键：

```
foo = bar


foo1 = bar1


foo3 = bar3


zoo = val


zoo1 = val1


zoo2 = val2


a = 123


b = 456


z = 789
```

删除键 foo 的命令：

```
$ etcdctl del foo


1 # 删除了一个键
```

删除从 foo 到 foo9 范围的键的命令：

```
$ etcdctl del foo foo9


2 # 删除了两个键
```

删除键 zoo 并返回被删除的键值对的命令：

```
$ etcdctl del --prev-kv zoo


1   # 一个键被删除


zoo # 被删除的键


val # 被删除的键的值
```

删除前缀为 zoo 的键的命令：

```
$ etcdctl del --prefix zoo


2 # 删除了两个键
```

删除大于等于键 b 的 byte 值的键的命令：

```
$ etcdctl del --from-key b


2 # 删除了两个键
```

#### watch 键值对的改动

etcd 的 watch 功能是一个常用的功能，我们来看看通过 etcdctl 如何实现 watch 指定的键值对。

watch 监测一个键值的变化，一旦**键值发生更新，就会输出最新的值并退出**。例如：用户更新 testkey 键值为 Hello watch。

```
$ etcdctl watch testkey


# 在另外一个终端: etcdctl put testkey Hello watch


testkey


Hello watch
```

从 foo to foo9 范围内键的命令：

```
$ etcdctl watch foo foo9


# 在另外一个终端: etcdctl put foo bar


PUT


foo


bar


# 在另外一个终端: etcdctl put foo1 bar1


PUT


foo1


bar1
```

以 16 进制格式在键 foo 上进行观察的命令：

```
$ etcdctl watch foo --hex


# 在另外一个终端: etcdctl put foo bar


PUT


\x66\x6f\x6f          # 键


\x62\x61\x72          # 值
```

观察多个键 foo 和 zoo 的命令：

```
$ etcdctl watch -i


$ watch foo


$ watch zoo


# 在另外一个终端: etcdctl put foo bar


PUT


foo


bar


# 在另外一个终端: etcdctl put zoo val


PUT


zoo


val
```

查看 key 的历史改动，应用可能想观察 etcd 中键的历史改动。

例如，应用服务想要获取某个键的所有修改。如果应用客户端一直与 etcd 服务端保持连接，使用 watch 命令就能够实现了。但是当应用或者 etcd 实例出现异常，该键的改动可能发生在出错期间，这样导致了应用客户端没能实时接收这个更新。因此，**应用客户端必须观察键的历史变动**，为了做到这点，应用客户端可以在观察时指定一个历史修订版本。

首先我们需要完成下述序列的操作：

```
$ etcdctl put foo bar         # revision = 2


OK


$ etcdctl put foo1 bar1       # revision = 3


OK


$ etcdctl put foo bar_new     # revision = 4


OK


$ etcdctl put foo1 bar1_new   # revision = 5


OK
```

观察历史改动：

```
# 从修订版本 2 开始观察键 `foo` 的改动


$ etcdctl watch --rev=2 foo


PUT


foo


bar


PUT


foo


bar_new
```

从上一次历史修改开始观察：

```
# 在键 `foo` 上观察变更并返回被修改的值和上个修订版本的值


$ etcdctl watch --prev-kv foo


# 在另外一个终端: etcdctl put foo bar_latest


PUT


foo         # 键


bar_new     # 在修改前键 foo 的上一个值


foo         # 键


bar_latest  # 修改后键 foo 的值
```

压缩修订版本。

参照上述内容，etcd 保存修订版本以便应用客户端可以读取键的历史版本。但是，为了避免积累无限数量的历史数据，需要对历史的修订版本进行压缩。**经过压缩，etcd 删除历史修订版本，释放存储空间，且在压缩修订版本之前的数据将不可访问**。

下述命令实现了压缩修订版本：

```
$ etcdctl compact 5


compacted revision 5 #在压缩修订版本之前的任何修订版本都不可访问


$ etcdctl get --rev=4 foo


{&quot;level&quot;:&quot;warn&quot;,&quot;ts&quot;:&quot;2020-05-04T16:37:38.020+0800&quot;,&quot;caller&quot;:&quot;clientv3/retry_interceptor.go:62&quot;,&quot;msg&quot;:&quot;retrying of unary invoker failed&quot;,&quot;target&quot;:&quot;endpoint://client-c0d35565-0584-4c07-bfeb-034773278656/127.0.0.1:2379&quot;,&quot;attempt&quot;:0,&quot;error&quot;:&quot;rpc error: code = OutOfRange desc = etcdserver: mvcc: required revision has been compacted&quot;}


Error: etcdserver: mvcc: required revision has been compacted
```

#### lease（租约）

lease 意为租约，类似于 Redis 中的 TTL(Time To Live)。etcd 中的键值对可以绑定到租约上，实现**存活周期控制**。在实际应用中，常用来实现服务的心跳，即服务在启动时获取租约，将租约与服务地址绑定，并写入 etcd 服务器，为了保持心跳状态，服务会定时刷新租约。

**授予租约**

应用客户端可以为 etcd 集群里面的键授予租约。当键被附加到租约时，它的存活时间被绑定到租约的存活时间，而租约的存活时间相应的被 TTL 管理。在授予租约时，每个租约的最小 TTL 值由应用客户端指定。**一旦租约的 TTL 到期，租约就会过期并且所有附带的键都将被删除**。

```
# 授予租约，TTL 为 100 秒


$ etcdctl lease grant 100


lease 694d71ddacfda227 granted with TTL(10s)


# 附加键 foo 到租约 694d71ddacfda227


$ etcdctl put --lease=694d71ddacfda227 foo10 bar


OK
```

在实际的操作中，**建议 TTL 时间设置久一点**，避免来不及操作而出现如下错误：

```
{&quot;level&quot;:&quot;warn&quot;,&quot;ts&quot;:&quot;2020-12-04T17:12:27.957+0800&quot;,&quot;caller&quot;:&quot;clientv3/retry_interceptor.go:62&quot;,&quot;msg&quot;:&quot;retrying of unary invoker failed&quot;,&quot;target&quot;:&quot;endpoint://client-f87e9b9e-a583-453b-8781-325f2984cef0/127.0.0.1:2379&quot;,&quot;attempt&quot;:0,&quot;error&quot;:&quot;rpc error: code = NotFound desc = etcdserver: requested lease not found&quot;}
```

**撤销租约**

应用通过租约 ID 可以撤销租约。撤销租约将删除所有附带的 key。

我们进行下列操作：

```
$ etcdctl lease revoke 694d71ddacfda227


lease 694d71ddacfda227 revoked


$ etcdctl get foo10
```

**刷新租期**

应用程序可以通过刷新其 TTL 保持租约存活，因此不会过期。

```
$ etcdctl lease keep-alive 694d71ddacfda227


lease 694d71ddacfda227 keepalived with TTL(100)


lease 694d71ddacfda227 keepalived with TTL(100)


...
```

**查询租期**

应用客户端可以查询租赁信息，检查续订或租赁的状态，是否存在或者是否已过期。应用客户端还可以查询特定租约绑定的 key。

我们进行下述操作：

```
$ etcdctl lease grant 300


lease 694d71ddacfda22c granted with TTL(300s)


$ etcdctl put --lease=694d71ddacfda22c foo10 bar


OK
```

获取有关租赁信息以及哪些 key 绑定了租赁信息：

```
$ etcdctl lease timetolive 694d71ddacfda22c


lease 694d71ddacfda22c granted with TTL(300s), remaining(282s)


$ etcdctl lease timetolive --keys 694d71ddacfda22c


lease 694d71ddacfda22c granted with TTL(300s), remaining(220s), attached keys([foo10])
```

### 小结

这一讲我们主要介绍了 etcdctl 相关命令的说明以及数据库命令的使用实践。etcdctl 为用户提供一些简洁的命令，用户通过 etcdctl 可以直接与 etcd 服务端交互。etcdctl 客户端提供的操作与 HTTP API 基本上是对应的，甚至可以替代 HTTP API 的方式。通过 etcdctl 客户端工具的学习，对于我们快速熟悉 etcd 组件的功能和入门使用非常有帮助。
