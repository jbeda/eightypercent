+++
date = "2016-05-29T09:54:29-07:00"
title = "Recipe: Docker Logs \u2192 Google Stackdriver Logging"

+++

Docker logs are a pain to deal with.  If you don't do anything they'll fill up
your disk[^docker-log-rotation].  In addition logs are easily lost if you need
to delete and recreate a container.

[^docker-log-rotation]: It isn't that hard to set up log rotation but you have to do it.  And after you do, `docker logs` can't see old logs.  Example on how to do that [here](https://github.com/jbeda/docker-postfix-forwarder#gotcha-logging).

Google Cloud Platform has a [debug log
service](https://cloud.google.com/logging/docs/) that has evolved from the
(excellent) logging system built in to App Engine.  With the acquisition and
integration of Stackdriver it looks like this has been rebranded/merged.

Recently a new [`gcplogs` Docker log
driver](https://docs.docker.com/engine/admin/logging/gcplogs/) was merged (by
[Mike Danese](https://github.com/docker/docker/pull/18766)) and it looks like it
is available in Docker 1.11.  Getting it up and working is pretty easy but there
are couple of gotchas.

## Service Account Scopes

First, I did this on a GCE instance that already had the right service account scopes set up.  It [looks like](https://developers.google.com/identity/protocols/googlescopes#loggingv2beta1) you need the `https://www.googleapis.com/auth/logging.write` or `https://www.googleapis.com/auth/cloud-platform`[^cloud-platform-scope] scopes.

[^cloud-platform-scope]: This scope is new since from when I worked at Google. I *think* it is a "catch all" scope across many GCP services.  Kind of like root.  I can't find any documentation on it specifically.

My instance already had the `logging.write` scope so I was all set.  I was able to tell by looking at the instance in the cloud console.  Or you can use gcloud (as below).  If you don't see the right scope here you'll need to restart your instance or pretend you aren't on GCE.

```
$ gcloud compute instances describe --format='yaml(serviceAccounts[].scopes[])' my-instance
serviceAccounts:
- scopes:
  - https://www.googleapis.com/auth/compute
  - https://www.googleapis.com/auth/devstorage.read_write
  - https://www.googleapis.com/auth/logging.write
  - https://www.googleapis.com/auth/userinfo.email
```

If you aren't on GCE, you'll want to get a service account key and use that.  I haven't done this but it is documented as part of the Google Cloud [Application Default Credentials](https://developers.google.com/identity/protocols/application-default-credentials#howtheywork) system.

## Activate the Logging API

You need to "enable" the logging API before you can use it.  The fact that you
have to manually toggle APIs is probably a historical artifact at GCP.  Just go
[here](https://console.developers.google.com/apis/api/logging/overview) and flip
it on.  Make sure you are pointed to the right project.

If you don't do this first you may tickle [a bug in
Docker/aufs](https://github.com/docker/docker/issues/21704#issuecomment-222365187)
that will make the container un-deletable.

## Launch a container

Now just launch a container and have it log to GCP.

```
$ docker run -d --name my-container \
  --log-driver=gcplogs \
  --log-opt gcp-log-cmd=true \
  [...]
```

You can also set gcplogs as the default for all containers in the daemon but I
haven't tried that yet.

## View Logs on Cloud Console

You can now view these logs on the [cloud
console](https://console.cloud.google.com/logs?project=jbeda-personal&service=custom.googleapis.com&key1=&key2=&logName=).
Again, make sure you are pointed at the right project.

{{<figure src="/images/gcp-log-ui.png" title="Google Cloud Console Logging" link="/images/gcp-log-ui.png">}}

To be honest, there is a lot of structure here and the web UI isn't as smooth as
it could be.  You can *filter* on these fields but you can't control what is
shown.

You can
[filter](https://cloud.google.com/logging/docs/view/advanced_filters#introduction)
which events are shown.  This bit is a bit confusing.  It looks like there are 2
filter UIs.  I was using the "advanced" filter UI that you can get to from the
little down arrow icon to the right of the filter box.  The easiest thing to
filter on is the container name like so:

```
structPayload.container.name="/my-container"
```

In addition, there is a v1 syntax and a v2 syntax.  I think this box uses the v1
syntax.  It isn't clear how deep the differences are.


## View logs from the command line

Because of the lack of control on which fields are show in the log viewer, it is
probably best to view logs from the command line.  Note that this is beta level
in `gcloud` so the examples here may break in the future.

First, install the beta `gcloud` tools if you don't already have them.

```
gcloud components install beta
```

Now you can list the logs easily using [`gcloud beta logging
read`](https://cloud.google.com/sdk/gcloud/reference/beta/logging/read).  I
bundled this up into a shell function that is in my `.bashrc`.

```bash
function gcloud-docker-logs {
  local dayago
  if [ $(uname -s) = "Darwin" ]; then
    dayago=$(date -v-1d -u '+%Y-%m-%dT%H:%M:%SZ')
  else
    dayago=$(date  --date="-1 day" -Isec -u)
  fi

  gcloud beta logging read \
    --order ASC \
    --format='value(timestamp, jsonPayload.data)' \
    "jsonPayload.container.name=\"/$1\" AND timestamp > \"${dayago}\""
}
```

Decoder ring:

* `--order ASC` makes it look like you are `cat`ing a log file.
* `--format='value(timestamp, jsonPayload.data)'` picks the fields we want to show.  This is just the timestamp and the log line.  For some reason in the web UI this is called `structPayload` but here it is `jsonPayload`.  Perhaps this is v1 vs v2 of the API and filter syntax?
* `"jsonPayload.container.name=\"/$1\""` is the filter to apply.  Here we just slap in the name.
* We limit this to the last day or logs.  This happens by default when `order` is set to `DESC` but we want ascending order.  Because of this we have to explicitly limit the date range.  `date` is different on BSD systems (MacOS) vs. GNU systems (Linux). Ug.

## Feedback for product groups

* Docker logging should have rotation by default.  The fact that out of the box it fills up your disk is horrible.
* I'm sad that you still can't change the scopes on a GCE instance while it is running.  This has been in the bug database since scopes were first added.
* There is a bunch of structure in that isn't making it into the logging system.  For instance, severity isn't reflected (at least put stderr as WARN?).  It also might be good to hoist the instance name and container name into the primary and secondary keys for logging.  `custom.googleapis.com/primary_key: "primary_key"` looks like something was left stubbed out.
* The logging UI should let you select what to show in addition to filtering the events.  As more structure gets added the simple JSON serialization that is there now won't work.  In addition, that JSON serialization seems to be non-deterministic with respect to field order.
* `gcloud beta logging read` should have a "follow" mode so that you can tail the logs.  This mode exists in the UI.
* The GCP logging filter syntax is confusing and inconsistent.  Simple vs advanced. v1 vs v2.
* Why am I the one writing this guide? I don't work for Google anymore.
