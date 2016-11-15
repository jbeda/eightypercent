+++
date = "2015-09-18T15:00:11-07:00"
draft = true
title = "Modern Cluster Storage Taxonomy"

+++

Following up on my previous article on the [Anatomy of a Modern Production Stack](http://www.eightypercent.net/post/layers-in-the-stack.html) I want to explore the different types of storage systems that can be used with and by that stack.

First, a disclaimer: I'm not a storage guy.  I've used these systems and I know enough to be dangerous but I'm probably missing some subtlety.  This map is useful for understanding how the storage systems can relate to a computing cluster and how to approach picking your production stack.

The **point of view** being used here is that of a large cluster where we value the following:

  * **Self healing and self managing**.  A storage system that can deal with individual components (disks, servers) failing and recover gracefully and automatically is to be prized.  At scale, these situations happen all the time.
  * **Efficient**. Using every part of a resource is better than not.  At the end of the day, this comes down to $ cost.
  * **Predictable performance**.  Fast storage is great but consistently fast storage is even better.

There are vastly simpler ways to solve these problems at a small scale.  For instance, if you only have one machine, you can buy super reliable hardware and wake up in the middle of the night when it fails (every couple of years).  You over provision for good and reliable performance and don't worry about the cost as it pales in comparison to your other costs.  As systems scale up though other costs start to dominate.

With that in mind here is how I categorize storage solutions in my mind:

* **Local Storage** Local storage is colocated with the compute workload.
  * While it is possible to buy 
  * **Ephemeral Local Storage** It is easy to have ephemeral storage that sits next to a compute workload.  If that workload dies or needs to move the storage is lost.  Workload placement doesn't need to think about storage.  Predictable performance can be difficult on shared devices as it is hard to manage bandwidth across workloads and things like disk syncing can slow everyone down.
  * **Persistent Local Storage** In this case a set of storage is allocated to a workload and has a lifetime independent of that workload.  If the workload needs to get started   