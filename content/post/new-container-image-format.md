+++
date = "2015-07-01T10:00:00-07:00"
title = "Container Native Package System"

+++

A lot of exciting things happened at [Dockercon 2015](http://www.dockercon.com/) last week.  For me the most exciting was the announcement of the [Open Container Project](http://opencontainers.org/) foundation.  Not only is it great to see the community come together under one banner but also as a chance to entertain new ideas in this space.

These are some thoughts on how to improve what we consider a "container image." I'm looking at both the container format itself and what goes on inside of it. This obviously builds on ideas in other systems and I've tried to call those out.  These thoughts are still early so I'm hoping to find others of like mind and start a good discussion.

## History and Context

The thing I've been exploring most recently is the intersection between the container image format and package management.  While there has been plenty of attention on base OSs to host the container ([CoreOS](https://coreos.com/products/), [RancherOS](http://rancher.com/rancher-os/), [Project Atomic](http://www.projectatomic.io/), [Snappy Ubuntu](https://developer.ubuntu.com/en/snappy/)) and efforts to coordinate a cluster of hosts ([Kubernetes](http://kubernetes.io/), [Mesosphere](https://mesosphere.com/), [Docker Swarm](https://www.docker.com/docker-swarm)) we haven't paid as much attention as we could to what is going on **inside** the container.

**Docker Images** are great. Images are pretty efficient to push and pull and, with new focus on security, it is getting easier and easier to be sure that what you want in the image is actually what you are running.

**Dockerfiles** are also great.  They are a purpose built makefile analog that are super easy to understand and logically build on the layered nature of Docker images.  Like most of the Docker project, they are much more approachable than other efforts in this area and solve real customer needs. When constructed appropriately, they allow for an efficient dev flow where many of the time consuming steps can be reused.

One of the **best innovations of Docker is actually a bit of an awesome hack**.  It leverages the package managers for existing Linux distributions.  Reusing the package manager means that users can read any number of guides on how to get software installed and easily translate it into a Dockerfile.

Think of it this way: a typical Linux distribution is 2 things. First is a bunch of stuff to get the box booted.  Second is a package manager to install and manage software on that box.  Docker images typically only need the second one.  The first one is along for the ride even if the user never needs it.  There are package managers out there that are cleanly factored from the underlying OS ([Homebrew](http://brew.sh/), [Nix](https://nixos.org/nix/)) but they aren't typically used in Docker images.

## Problems

This all mostly works okay.  There is some cruft in the images that can easily be ignored and is "cheap" as the download and storage cost is amortized because of layer sharing for Docker images.

But we can do better.

Problems with the current state of the world:

* **No package introspection.**  When the next security issue comes along it is difficult to easily see which images are vulnerable.  Furthermore, it is hard to write automated policy to prevent those images from running.
* **No easy sharing of packages.**  If 2 images install the same package, the bits for that package are downloaded twice.  It isn't uncommon for users to construct complicated "inheritence" chains to help work around this issue[^dockerfile-chains].
* **No surgical package updating.**  Updating a package requires recreating an image and all re-running all downstream actions in the Dockerfile.  If users are good about tracking which sources go into which image[^tracking-inputs], it should be possible to just update the package but that is difficult and error prone.
* **Order dependent image builds.** Order matters in a Dockerfile --- even when it doesn't have to.  Often times two actions have zero interaction with each other.  But Docker has no way of knowing that so must assume that every action depends on everything coming previously.
* **Package manager cruft.** Most well built Dockerfiles have something like:

    ```
    RUN apt-get update \
      && apt-get install -y --no-install-recommends \
        build-essential \
      && rm -rf /var/lib/apt/lists/*
    ```

    This helps to minimize the size of the layer on disk.  This is confusing boilerplate that is likely just cargo-culted by many users.

[^dockerfile-chains]: The standard golang image is a great example of this.  `golang:1.4.2` &rarr; `buildpack-deps:jessie-scm` &rarr; `buildpack-deps:jessie-curl` &rarr; `debian:jessie`.  Most of this is done to enable efficient sharing of installed packages.
[^tracking-inputs]: Best practices should be to track every single input into your docker file.  That means that if you are pushing sources you should know which git commit, for example, those sources come from.  My guess is that this is rarely done.

## Ideas for Solutions

While I don't have a fully formed solution to all of these problems, I do have a bunch of ideas.

First, imagine that we take the idea of a container image and break it down a little.  The first thing we define is a `FileSet`.  A `FileSet` is a named, versioned and verified set of files.  Google has a system internally called the "Midas Package Manager" (MPM) that this[^decentralized-mpm].  Dinah McNutt gave a [great talk on MPM](https://www.youtube.com/watch?v=_uJlTllziPI) at a 2013 USENIX conference.  A further tweak would allow a `FileSet` to import other `FileSet`s into the file namespace of the host.  This allows for a `FileSet` to have multiple "parents" -- unlike the current Docker image format.

[^decentralized-mpm]: Actually, we need our system to be decentralized.  MPM, like may package management systems (including [Homebrew](https://github.com/Homebrew/homebrew/tree/master/Library/Formula) and [Nix](https://github.com/NixOS/nixpkgs)) has a single central repository/database of all packages.  Whatever is used here must be distributed --- probably in a namespace rooted with DNS.  Something like [Docker Notary](https://github.com/docker/notary) would play a role in signing and verifying packages.  Something like the [Nix archive format (NAR)](http://lethalman.blogspot.com/2014/08/nix-pill-9-automatic-runtime.html) will help make this more stable and predictable.

Second, we define a `Package` as a type of `FileSet`.  It would have a standard directory structin and include metadata on other packages required along with simple instructions to "install" the package[^package-install].Ideally, these packages would be built from verified sources and a verified tool chain.  This would enable true provenance for every bit.  This would be optional.

[^package-install]: Package install should consist of simply symlinking files into some common directories (`/bin`, `/lib`).  This would all be done via a declarative manifest. There are probably going to be cases where an "install" is a little bit more complex and a script is necessary.  I'd love to see how far we could get before that becomes absolutely necessary.  It is also assumed that the package directories themselves are only ever mounted read only. 

Finally, we would redefine a `ContainerImage` also as a type of `FileSet` that has metadata necessary to make it runnable.  The definition of this metadata is a big part of what the [Docker Image format](https://github.com/docker/docker/blob/master/image/spec/v1.md) and the [ACI](https://github.com/appc/spec/blob/master/SPEC.md) format are.

A `ContainerImage` that is using this container native package system would define a set of read-only imports of all required packages `FileSet`s.  Image construction tools would verify that all dependencies are satisfied.  Furthermore, the install steps would be run to symlink all of the packages into the appropriate places[^late-symlinking].

[^late-symlinking]: There is a choice on when the package install happens.  It could happen early as the container is created.  Or it could happen late as part of the container start process.  I'd prefer late binding as it makes surgical package updating simpler.  The directories that store the symlinks could be tmpfs directories to keep this all very speedy.

User code could either be packed up as a `Package` or just inserted directly into the `ContainerImage`.

## Analysis

By creating a container friendly packaging system and expanding the idea of what an container image is, we can solve most of the issues outlined again.  

* The list of `FileSet`s imported into, say, `/packages` would be the list of all packages versions that are included in that image.
* Individual `FileSet`s could be cached by hosts and easily and safely shared between disparate images.
* A package could be updated in a straightforward way.  The toolset would have to make sure that all dependencies are satisfied and that the install steps are run as necessary.
* Image build tools would list the packages necessary and order wouldn't matter.  Because there are multiple "parents" to an image, order cannot matter.
* The package install cruft (archived version of the package) would be handled on the host side similar to images themselves.  The only thing the container would see would be the actual files -- and they would be symlinked in.

There are some missing and underdefined parts to this story.

* _How are packages created?_  I'm thinking that we could do that by running a container with build time packages that produces output file into a specific directory.  Files in that directory are then used to create a package.  As part of this the inputs into the build container could be included in the package metadata and signed.
* _What does a package distribution look like?_  I imagine we'd have curated sets of packages that are known to work well together.  For instance, xyz.com could create `xyz.com/apache` that depends on `xyz.com/openssl`.
* _How do users override packages?_  Perhaps `abc.com/openssl` could specify that it can be used in place of `xyz.com/openssl`.  Any guarantees by `xyz.com` would be void but it would be a way to do custom versions and carry patches.
* _Opportunity: Kernel and capability requirements._  Packages could specify their requirements in a way that would be visible to the host.  This would provide a more direct requirement chain between the host and the code running in the container.

This solution obviously borrows from both Homebrew and Nix.  What I think is new here is the idea of expanding the definition of an container image with FileSets and making this be fundamentally decentralized.  We also need to ensure that the easy to approach spirit of Dockerfile isn't lost.  If we do this right we can make images much easier to efficiently create, verify, update and manage.

Ping back to me on twitter ([@jbeda](http://www.twitter.com/jbeda)) or we can talk over on [Hacker News](https://news.ycombinator.com/item?id=9814131)

(Thanks to [Kelsey Hightower](https://twitter.com/kelseyhightower) for reviewing an earlier version of this post.)
