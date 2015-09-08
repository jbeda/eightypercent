+++
date = "2015-09-07T16:23:25-07:00"
draft = true
title = "Anatomy of a Modern Production Stack"

+++

I was chatting on an Xoogler message board the other day and Dennis Ordanov ([@daodennis](https://twitter.com/daodennis)) was asking about the basic moving parts of a production stack.  I just started enumerating them from memory and thought it might be a good blog post.  So, here is a mostly stream-of-consciousness dump of the parts a modern (container based) production environment[^caveats].

* **Production Host OS**.  This is a simplified and managable Linux distribution.  Usually it is just enough to get a container engine up and running.  
  * Examples include [CoreOS](https://coreos.com/using-coreos/), [Red Hat Project Atomic](http://www.projectatomic.io/), [Ubuntu Snappy](https://developer.ubuntu.com/en/snappy/), and [Rancher OS](http://rancher.com/rancher-os/).
* **Container Engine**. This is the system for setting up and managing contianers.  It is the primary management agent on the node.  
  * Examples include [Docker Engine](https://www.docker.com/docker-engine), [CoreOS rkt](https://coreos.com/rkt/docs/latest/), and [LXC](https://linuxcontainers.org/) and [systemd-nspawn](http://www.freedesktop.org/software/systemd/man/systemd-nspawn.html).  
  * Some of these systems are more amenable to being directly controlled remotely than others. 
  * The [Open Container Initiative](https://www.opencontainers.org/) is working to standarize the input into these systems -- basically the root filesystem for the container along with some common parameters in a JSON file.
* **Container Image Packaging and Distribution.** A Container Image is a named and cloneable chroot that can be used to create container instances.  It is pretty much an efficient way to caputure, name and distribute the set of files that make up a container at runtime.  
  * Both Docker and CoreOS rkt solve this problem.  It is built into the Docker Engine but is broken out for rkt as a separate tool set call [acbuild](https://github.com/appc/acbuild).  
  * Inside of Google this was done slightly differently with a file package distribution system called [MPM](https://www.youtube.com/watch?v=_uJlTllziPI).
  * Personally, I'm hoping that we can define a widely adopted spec for this, hopefully as part of the OCI.
* **Container Image Registry/Repository.** This is a central place to store and load Container Images.
  * Hosted versions of this include the [Docker Hub](https://hub.docker.com/), [Quay.io (owned by CoreOS)](https://quay.io), and [Google Container Registry](https://cloud.google.com/container-registry/).
  * Docker also has an open source registry called [Docker Distribution](https://github.com/docker/distribution).
  * Personally, I'm hoping that the state of the art will evolve past centralized solutions with specialized APIs to solutions that are simpler by working regular HTTP and more transport agnostic so that protocols like BitTorrent can be used to distribute images.
* **Container Distribution.** This is the system for structuring what is running inside of a container.  Many people don't talk about this as a separate thing but instead reuse OS distributions such as Ubuntu, Debian, or CentOS.
  * Many folks are working to build minimal container distributions by either using distributions based in the embedded world (BusyBox or Alpine) or by building [static binaries](https://medium.com/@kelseyhightower/optimizing-docker-images-for-static-binaries-b5696e26eb07) and not needing anything else.
  * Personally, I'd love to see the idea of a Container Distribution be further developed and take advantage of features only available in the container world.  I wrote a [blog post](http://www.eightypercent.net/post/new-container-image-format.html) on this.
* **Container Orchestration System.** Once you have containers running on a single host, you need to get them running across multiple hosts.  
  * This is a super hot area of interest with lots of innovation.
  * Open source deployable examples include [Kubernetes](http://kubernetes.io/), [Docker Swarm](https://docs.docker.com/swarm/), and [Apache Mesos](http://mesos.apache.org/).
  * Hosted systems include [Google Container Engine (GKE)](https://cloud.google.com/container-engine/) (based on Kubernetes), [Mesosphere DCOS](https://mesosphere.com/product/) and [Amazon EC2 Container Service (ECS)](https://aws.amazon.com/ecs/).
* **Orchestration Config.** Many of the orchestration systems have small granular objects.  Creating and parameterizing these by hand can be difficult.  In this context, an orchestration config system can take higher level description and compile them down to the nuts and bolts that the orchestration systems works with.
  * The Google solutions to this problem have never been made public (to my knowledge).
  * [AWS CloudFormation](https://aws.amazon.com/cloudformation/) and [Google Cloud Deployment Manager](https://cloud.google.com/deployment-manager/overview) play this role for their respective cloud ecosystems (only).
  * [Hashicorp Terraform](https://github.com/hashicorp/terraform) and [Flabbergast](http://flabbergast.org/) look like they could be applied to container orchestration systems but haven't yet.
  * [Docker Compose](https://docs.docker.com/compose/) is a start to a more comprehensive config system.
  * The Kubernetes team (Brian Grant especially) have lots of [ideas and plans](https://github.com/kubernetes/kubernetes/labels/area%2Fapp-config-deployment) for this area.  There is a [Kubernetes SIG](https://github.com/kubernetes/kubernetes/wiki/Special-Interest-Groups-(SIGs%29) being formed.
* **Network Virtualization.** While not strictly necessary, clustered container systems are much easier to use if each container has full presence on the cluster network.  This has been referred to as "IP per Container".
  * Without a networking solution, orchestration systems must allocate and enforce port assignment as ports per host are a shared resource.
  * Examples here include [CoreOS Flannel](https://github.com/coreos/flannel), [Weave](http://weave.works/), [Project Calico](http://www.projectcalico.org/), and [Docker libnetwork](https://github.com/docker/libnetwork) (not ready for production yet).
* **Container Storage Systems.** As users move past special "pet" hosts storage becomes more difficult.
  * I have more to say on this that I'll probably put into a blog post at some point in the future.
  * [ClusterHQ Flocker](https://github.com/clusterhq/flocker) deals with migrating data between hosts (among other things).
* **Discovery Service.** Discovery is a fancy term for naming.  Once you launch a bunch of containers, you need to figure out where they are so you can talk to them.
  * DNS is often used as a solution here but can cause issues in highly dynamic environments due to aggressive caching.  Java, in particular, is troublesome as it [doesn't honor DNS TTLs by default](https://www.google.com/search?btnG=1&pws=0&q=networkaddress.cache.ttl+default&gws_rd=ssl).
  * Many people build on top of highly consistent stores (lock servers) for this.  Examples include: [Apache Zookeeper](https://zookeeper.apache.org/), [CoreOS etcd](https://coreos.com/etcd/), [Hashicorp Consul](https://www.consul.io/).
  * Kubernetes supports [service definition and discovery](https://github.com/kubernetes/kubernetes/blob/master/docs/user-guide/services.md) (with a stable virtual IP with load balanced proxy).
  * Related is a system to configure wider facing load balancer to manage the interface between the cluster and the wider network.
* **Production Identity and Authentication.** As clustered deployments grow, an identity system becomes necessary.  When microservice A calls microservice B, microservice B needs some way to verify that it is actually microservice A calling.  Note that this is for computer to computer communication within the cluster 
  * This is not a well understood component of the stack.  I expect it to be an active area of development in the near future.  Ideally the orchestration system would automatically configure the identity for each running container in a secure way.
  * Related areas include secret storage and authorization.
  * I've used the term "Authentity" to describe this area. Please use it as I'm hoping it'll catch on.
  * [conjur.net](http://conjur.net) is a commercial offering that can help out in this situation.
* **Monitoring.** A modern production application has deep monitoring.  Not only should the operations folks make sure that the binaries continue to run, they should also be monitoring application specific metrics that are thrown off by each microservice that makes up that application.  
  * A modern monitoring solution is characterized by its ability to deal with a wide set of metrics along with flexible vector math necessary to do complex aggregations and tests.  Time series data should be sampled frequently (30-60s or less), stored for a long time and be easily explored and graphed.  A good monitoring system not only lets you know when things are down but also is a critical debugging tool to know where and how things are broken.
  * For open systems, [Prometheus](http://prometheus.io/) looks *very* intersting in this area.  
  * There are also ways to break this apart into smaller parts such as [Grafana](http://grafana.org/) as a frontend backed by a dedicated time series database like [InfluxDB](https://influxdb.com/) or [OpenTSDB](http://opentsdb.net/).
  * [Heapster](https://github.com/kubernetes/heapster) is a container specific monitoring system that surfaces data collected by [cAdvisor](https://github.com/google/cadvisor).
  * There are hosted systems such as [Google Cloud Monitoring](https://cloud.google.com/monitoring/) and [SignalFx](https://signalfx.com/).
  * I'm not an expert here so I'm sure I'm missing some of the awesome stuff going on in this area.
* **Logging.** Logging can generally be broken into two types: unstructured debug logs and structured application logs.  The debug logs are used to figure out what is going on with the system while structured logs are usually used to capture important events for later processing and aggregation.  One might use structured logs for ad impressions that are critical for reconciling real money.
  * Systems such as [fluentd](http://www.fluentd.org/) and [logstash](https://www.elastic.co/products/logstash) are agents that collect and upload logs.
  * There are a ton of systems for storing and indexing logs.  This includes [elasticsearch](https://www.elastic.co/) along with more traditional databases (MySQL, Mongo, etc.).
  * Hosted systems include [Google Cloud Logging](https://cloud.google.com/logging/docs/).
  * Logging can also throw off monitoring signals.  For instance, while processing saved logs the local agent can count 500s and feed those into a monitoring system.
  * Systems like [Apache Flume](http://flume.apache.org/)[^flume] can be used to collect and reliably save structured logs for processing in the Hadoop ecosystem.  [Google BigQuery](https://cloud.google.com/bigquery/) and [Google Cloud Dataflow](https://cloud.google.com/dataflow/) are also well suited to ingesting and analyzing structured log data.

Next on the list would be to talk about continuous integration/continuous deployment (CI/CD) systems and systems for communicating between microservices (RPC and queues). But I think I'll stop here.  If this is useful (or if you think I'm missing anything huge) please let me know via [twitter](https://www.twitter.com/jbeda).

[^caveats]:
    Some caveats:

    * I'm sure I'm missing some parts of the stack.
    * The way I break this problem down is based on my experiences at Google. There are many ways to skin this cat.
    * I've listed example projects/products/companies/systems at different levels but this isn't meant to be exhaustive.
    * The fact that I've listed a system here doesn't mean that I've run it in production and it has my stamp of approval.

[^flume]: Don't confuse Apache Flume with [Google FlumeJava](http://research.google.com/pubs/pub35650.html).  I guess once you start processing logs some names are just obvious.  Also see [Google Sawzall](http://research.google.com/archive/sawzall.html) and [Google Dremel](http://research.google.com/pubs/pub36632.html).