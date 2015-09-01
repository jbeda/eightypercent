+++
date = "2015-09-01T14:51:41-07:00"
title = "Operations Lock-in vs. Development Lock-in"

+++

I'm not a large enterprise developer.  My experiences are at the unique environments of Google and Microsoft[^exp].  Over my career I've helped build _platforms_ and I've learned that to build a great developer platform (or any good product at all) requires you to talk to many customers and put yourself in their mindset.  This extends from the day to day usability of the product to the decision of a customer to bet a part of his/her business and career on the product in the first place.

[^exp]: I started my career on Internet Explorer and moved on to create graphics APIs (including XAML) for Windows Presentation Foundation (WPF).  At Google I started on Google Talk building out [XMPP specifications for call signaling](http://xmpp.org/about-xmpp/technology-overview/jingle/) and eventually started Google Compute Engine and helped to start Kubernetes. 

Lock-in is obviously one of the critical things that many customers take into account. But lock-in isn't absolute.  Different types of lock-in are defined by the friction and cost incurred in switching.

One interesting way to view lock-in is by looking at who will bear the brunt that switching cost.  You can look at this from the point of view of Ops vs. Dev[^data-lock-in].

[^data-lock-in]: There are, of course, other types of lock-in.  One that is really interesting is data locality lock-in.  As data gets bigger and bigger, simply moving that data around can be complicated and costly.  This gets even more complicated as that data continues to grow while it is being moved.

* With **Operations Lock-in** the operations team is primarily involved in a migration from one place to another.  The operations playbooks, policies, and procedures may have to change but the code will be substantially unchanged.  Ideally, nothing needs to be recompiled or rebuilt.
* With **Developer Lock-in** someone has to actually go in and crack the code.  This could be something minor (swapping in a new library with a standard language interface) to more architectural changes (changing your database).

Different organizations have different skill sets and one type of lock-in might be more serious than another.  While someplace like Google or Microsoft may have high quality software engineers to spare, most large enterprises are severly constrained here.  Often **projects were written by developers that are long gone** -- they were independent contractors, left the company, or have moved on to other projects.  In this case developer lock-in can be much more costly than operations lock-in.

Think about it this way -- a user has a large system built on open software packages running on on-prem hardware.  **Switching to a VM/IaaS offering is a relatively cheap project** as it can be an operations focused project with little to no developer involvement.  As IaaS offerings evolved to be more and more similar to real hardware (persistent network block stores, familiar networking) this trend really accelerated.

On the other hand, if a user builds a software product around a very proprietary set of APIs (either language APIs or service APIs) switching can be much more expensive.  Many PaaS offerings would be in this category.  For example, **moving a large application away from AWS Lambda would require a rewrite**.  Google App Engine has similar issues if the native datastore is being used.

This is all critical to keep in mind as we create new ways to build and run systems.  Many times betting on open APIs[^open-apis] and open source systems helps to fight lock-in by making switching an operations project vs. a development project.  This is why, despite the eye-rolls, I'm excited about the emergance of the [Open Container Initiative](https://www.opencontainers.org/) and the [Cloud Native Computing Foundation](https://cncf.io/).

What other ways do you think about lock-in?  Tweet at me ([@jbeda](https://twitter.com/jbeda)) and let's start a discussion.

[^open-apis]: With the advent of cloud APIs, we've actually taken a bit of a step back in terms of open APIs.  While there are some companies ([such as Slack](https://twitter.com/stewart/status/634533296555339777)) that encourate alternate implementations of their APIs, many companies are silent or hostile to others implementing those APIs.  I'd love to see a pattern of companies releasing their APIs with an open license or donating the IP to a foundation. This is made only more critical in light of the [Oracle API copyright decision](https://en.wikipedia.org/wiki/Oracle_America,_Inc._v._Google,_Inc.).
