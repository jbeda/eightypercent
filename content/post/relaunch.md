+++
date = "2015-06-10T17:08:02-07:00"
title = "A New Beginning..."

+++

Welcome to the new [80%](http://www.eightypercent.net)!

I've rewritten my ancient blog on top of [Hugo](http://gohugo.io/) using [Bootstrap](http://getbootstrap.com/) and [Google Fonts](https://www.google.com/fonts).  The whole thing is hosted on [GitHub](https://github.com/jbeda/eightypercent) if you want to check out the source.

For headlines, I'm using [Economica](https://www.google.com/fonts/specimen/Economica).  I picked it because it is narrow and I tend to write long headlines.  The body text is [Lora](https://www.google.com/fonts/specimen/Lora).  I like a nice Serif font for the body text and it has a bit of style while still being very readable.  Code is in [Roboto Mono](https://www.google.com/fonts/specimen/Roboto+Mono) and the Logo and Navigation is in [Coda](https://www.google.com/fonts/specimen/Coda).

I'm not in love with the color theme I have now (dark gray background with red and blue highlights) but it'll do for now.  I played with some CSS gradients to show a preview on the main index page.  I'm still not sure -- it may be a little too cute.

I've done the work of importing my old blog posts and making Hugo create compatible URLs.  My original system used XML files for each day so I wrote a [short Go program](https://github.com/jbeda/eightypercent/blob/master/old-blog/convert.go) to convert from the XML to something Hugo can consume (markdown files with a JSON header).  I had to introduce a [verbatim HTML shortcode](https://github.com/jbeda/eightypercent/blob/master/layouts/shortcodes/verbatim.html) to make Hugo pass through the HTML from the old blog directly.  

The other tricky thing was to generate a page per *day* instead of a page per post for these old posts.  Rolling things up on a daily basis was the way things were done back when I started the blog in [January 2003](http://new.eightypercent.net/post/old/00001.html).  Doing this daily roll up used the [taxonomy](http://gohugo.io/taxonomies/overview/) feature of Hugo where the taxonomy was "archive" and the term was the day.  This was easy to generate from the Go conversion program.

I'm hosting this on GCE under Docker and NGINX.  I'm going to write up a post on how I'm doing that and automatically syncing it to the git repro on each submit.

Now that I have it all set up, we'll see if I actually have anything to say...