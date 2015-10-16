+++
date = "2015-09-08T12:10:45-07:00"
title = "Anatomy of a Modern Production Stack"

+++

(I'm updating this post as folks comment.  You can look at the [history on github](https://github.com/jbeda/eightypercent/commits/master/content/post/layers-in-the-stack.md).)

I was chatting on an Xoogler message board the other day and Dennis Ordanov ([@daodennis](https://twitter.com/daodennis)) was asking about the basic moving parts of a production stack.  I just started enumerating them from memory and thought it might be a good blog post[^other-posts].  So, here is a mostly stream-of-consciousness dump of the parts a modern (container based) production environment[^caveats].

*A note on the term "modern":*  This is my view, based on experiences at Google, for a stack that delivers what I'd want for a major production system.  You can do this without containers, but I think it is hard to meet my criteria that way.  The full stack here probably isn't necessary for small applications and, as of today, is **way** too hard to get up and running.  The qualities that I'd look for in a "modern" stack:

* **Self healing and self managing.** If a machine fails, I don't want to have to think about it. The system should just work.
* **Supports microservices.** The idea of breaking your app into smaller components (regardless of the name) can help you to scale your engineering organization by keeping the dev team for each µs small enough that a 2 pizza team can own it.
* **Efficient.** I want a stack that doesn't require a lot of hand holding to make sure I'm not wasting a ton of resources.
* **Debuggable.** Complex applications can be hard to debug. Having good strategies for application specific monitoring and log collection/aggregation can really help to provide insights into the stack.

So, with that, here is a brain dump of the parts that make up a "modern" stack:

* **Production Host OS**.  This is a simplified and manageable Linux distribution.  Usually it is just enough to get a container engine up and running.  
  * Examples include [CoreOS](https://coreos.com/using-coreos/), [Red Hat Project Atomic](http://www.projectatomic.io/), [Ubuntu Snappy](https://developer.ubuntu.com/en/snappy/), and [Rancher OS](http://rancher.com/rancher-os/).
* **Bootstrapping system**.  Assuming you are starting with a generic VM image or bare metal hardware, something has to be able to bootstrap those machines and get them running as productive members of the cluster.  This becomes very important as you are dealing with lots machines that come and go as hardware fails.
  * [Cloud Foundry BOSH](https://bosh.io/docs) was created to do this for Cloud Foundry but is seeing new life as an independent product.
  * The standard config tools ([Puppet](https://puppetlabs.com/), [Chef](https://www.chef.io/), [Ansible](http://www.ansible.com/home), [Salt](http://saltstack.com/)) can serve this role.
  * [CoreOS Fleet](https://coreos.com/using-coreos/clustering/) is a lightweight clustering system that can also be used to bootstrap more comprehensive solutions.
* **Container Engine**. This is the system for setting up and managing containers.  It is the primary management agent on the node.  
  * Examples include [Docker Engine](https://www.docker.com/docker-engine), [CoreOS rkt](https://coreos.com/rkt/docs/latest/), and [LXC](https://linuxcontainers.org/) and [systemd-nspawn](http://www.freedesktop.org/software/systemd/man/systemd-nspawn.html).  
  * Some of these systems are more amenable to being directly controlled remotely than others. 
  * The [Open Container Initiative](https://www.opencontainers.org/) is working to standardize the input into these systems -- basically the root filesystem for the container along with some common parameters in a JSON file.
* **Container Image Packaging and Distribution.** A Container Image is a named and cloneable chroot that can be used to create container instances.  It is pretty much an efficient way to capture, name and distribute the set of files that make up a container at runtime.  
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
  * Hosted systems include [Google Container Engine (GKE)](https://cloud.google.com/container-engine/) (based on Kubernetes), [Mesosphere DCOS](https://mesosphere.com/product/) and [Amazon EC2 Container Service (ECS)](https://aws.amazon.com/ecs/).  Recently announced is the [Microsoft Azure Container Service](https://azure.microsoft.com/en-us/blog/azure-container-service-now-and-the-future/) based on Mesosphere DCOS and Docker Swarm.
* **Orchestration Config.** Many of the orchestration systems have small granular objects.  Creating and parameterizing these by hand can be difficult.  In this context, an orchestration config system can take higher level description and compile them down to the nuts and bolts that the orchestration systems works with.
  * The Google solutions to this problem have never been made public (to my knowledge).
  * [AWS CloudFormation](https://aws.amazon.com/cloudformation/) and [Google Cloud Deployment Manager](https://cloud.google.com/deployment-manager/overview) play this role for their respective cloud ecosystems (only).
  * [Hashicorp Terraform](https://github.com/hashicorp/terraform) and [Flabbergast](http://flabbergast.org/) look like they could be applied to container orchestration systems but haven't yet.
  * [Docker Compose](https://docs.docker.com/compose/) is a start to a more comprehensive config system.
  * The Kubernetes team (Brian Grant especially) have lots of [ideas and plans](https://github.com/kubernetes/kubernetes/labels/area%2Fapp-config-deployment) for this area.  There is a [Kubernetes SIG](https://github.com/kubernetes/kubernetes/wiki/Special-Interest-Groups-(SIGs%29) being formed.
* **Network Virtualization.** While not strictly necessary, clustered container systems are much easier to use if each container has full presence on the cluster network.  This has been referred to as "IP per Container".
  * Without a networking solution, orchestration systems must allocate and enforce port assignment as ports per host are a shared resource.
  * Examples here include [CoreOS Flannel](https://github.com/coreos/flannel), [Weave](http://weave.works/), [Project Calico](http://www.projectcalico.org/), and [Docker libnetwork](https://github.com/docker/libnetwork) (not ready for production yet).  I've also been pointed to [OpenContrail](http://www.opencontrail.org/) but haven't looked deeply.
* **Container Storage Systems.** As users move past special "pet" hosts storage becomes more difficult.
  * I have more to say on this that I'll probably put into a blog post at some point in the future.
  * [ClusterHQ Flocker](https://github.com/clusterhq/flocker) deals with migrating data between hosts (among other things).
  * I know there are other folks (someone pointed me at [Blockbridge](http://www.blockbridge.com/)) that are working on software defined storage systems that can work well in this world.
* **Discovery Service.** Discovery is a fancy term for naming.  Once you launch a bunch of containers, you need to figure out where they are so you can talk to them.
  * DNS is often used as a solution here but can cause issues in highly dynamic environments due to aggressive caching.  Java, in particular, is troublesome as it [doesn't honor DNS TTLs by default](https://www.google.com/search?btnG=1&pws=0&q=networkaddress.cache.ttl+default&gws_rd=ssl).
  * Many people build on top of highly consistent stores (lock servers) for this.  Examples include: [Apache Zookeeper](https://zookeeper.apache.org/), [CoreOS etcd](https://coreos.com/etcd/), [Hashicorp Consul](https://www.consul.io/).
  * Kubernetes supports [service definition and discovery](https://github.com/kubernetes/kubernetes/blob/master/docs/user-guide/services.md) (with a stable virtual IP with load balanced proxy).
  * Weave has a built in [DNS server](http://blog.weave.works/2015/09/08/weave-gossip-dns/) that stores data locally so that TTLs can be minimal.
  * Related is a system to configure wider facing load balancer to manage the interface between the cluster and the wider network.
* **Production Identity and Authentication.** As clustered deployments grow, an identity system becomes necessary.  When microservice A calls microservice B, microservice B needs some way to verify that it is actually microservice A calling.  Note that this is for computer to computer communication within the cluster 
  * This is not a well understood component of the stack.  I expect it to be an active area of development in the near future.  Ideally the orchestration system would automatically configure the identity for each running container in a secure way.
  * Related areas include secret storage and authorization.
  * I've used the term "Authentity" to describe this area. Please use it as I'm hoping it'll catch on.
  * [conjur.net](http://conjur.net) is a commercial offering that can help out in this situation.
* **Monitoring.** A modern production application has deep monitoring.  Not only should the operations folks make sure that the binaries continue to run, they should also be monitoring application specific metrics that are thrown off by each microservice that makes up that application.  
  * A modern monitoring solution is characterized by its ability to deal with a wide set of metrics along with flexible vector math necessary to do complex aggregations and tests.  Time series data should be sampled frequently (30-60s or less), stored for a long time and be easily explored and graphed.  A good monitoring system not only lets you know when things are down but also is a critical debugging tool to know where and how things are broken.
  * For open systems, [Prometheus](http://prometheus.io/) looks *very* interesting in this area.  
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
* **Deep Inspection and Tracing** There are a class of tools that help to do deep debugging.  
  * Inside of Google, [Dapper](http://research.google.com/pubs/pub36356.html) is a great example of tracing a user request across many microservices.  [Appdash](https://github.com/sourcegraph/appdash) and [Zipkin](http://zipkin.io/) are open source system inspired by Dapper. 
  * Startups like [Sysdig](http://www.sysdig.org/) in this category too by allowing deep inspection and capture of what is going on with a machine.

PaaS systems often help to bring this all together in an easy way.  Systems like [OpenShift 3](http://www.openshift.org/), [Deis](http://deis.io/), or [Flynn](https://flynn.io) build on top of some of the independent systems above.  Other PaaS such as [Heroku](https://www.heroku.com/), [Google App Engine](https://cloud.google.com/appengine/docs) or [Cloud Foundry](https://www.cloudfoundry.org/) are more vertically integrated without the component layers being broken out in a well supported way.

Next on the list would be to talk about continuous integration/continuous deployment (CI/CD) systems and systems for communicating between microservices (RPC and queues). But I think I'll stop here.  If this is useful (or if you think I'm missing anything huge) please let me know via [twitter](https://www.twitter.com/jbeda).  Or you can comment on the [Hacker News thread](https://news.ycombinator.com/item?id=10187598).

[^other-posts]: [Brandon Philips](https://twitter.com/brandonphilips) from CoreOS points me to a [similar post](https://coreos.com/blog/cluster-osi-model/) from [Barak Michener](https://twitter.com/barakmich).  I go into more minutia here and don't try and define a strict stack.

[^caveats]:
    Some caveats:

    * I'm sure I'm missing some parts of the stack.
    * The way I break this problem down is based on my experiences at Google. There are many ways to skin this cat.
    * I've listed example projects/products/companies/systems at different levels but this isn't meant to be exhaustive.
    * The fact that I've listed a system here doesn't mean that I've run it in production and it has my stamp of approval.

[^flume]: Don't confuse Apache Flume with [Google FlumeJava](http://research.google.com/pubs/pub35650.html).  I guess once you start processing logs some names are just obvious.  Also see [Google Sawzall](http://research.google.com/archive/sawzall.html) and [Google Dremel](http://research.google.com/pubs/pub36632.html).