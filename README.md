# :ocean: `spate`

[![GoDoc](https://godoc.org/github.com/mtneug/spate?status.png)](https://godoc.org/github.com/mtneug/spate)
[![Build Status](https://travis-ci.org/mtneug/spate.svg?branch=master)](https://travis-ci.org/mtneug/spate)
[![codecov](https://codecov.io/gh/mtneug/spate/branch/master/graph/badge.svg)](https://codecov.io/gh/mtneug/spate)
[![Docker Image Version](https://images.microbadger.com/badges/version/mtneug/spate.svg)](https://hub.docker.com/r/mtneug/spate/)
[![Docker Image Layers](https://images.microbadger.com/badges/image/mtneug/spate.svg)](https://microbadger.com/images/mtneug/spate)

> **spate** `/speɪt/`<br>
> noun
>
> 1.  freshet, flood
> 2.  a) a large number or amount<br>
>     b) a sudden or strong outburst rush
>
> — [Merriam-Webster Dictionary](https://www.merriam-webster.com/dictionary/spate)

`spate` is a horizontal service autoscaler for [Docker Swarm mode](https://docs.docker.com/engine/swarm/) inspired by Kubernetes' [Horizontal Pod Autoscaler](http://kubernetes.io/docs/user-guide/horizontal-pod-autoscaling/).

Currently `spate` can scale services based on exposed [Prometheus](https://prometheus.io/) metrics. However, the foundations already have been layed for different types. Future versions will, for instance, be able to scale based on the CPU or memory usage.

## Installation

For every release static Linux binaries can be downloaded from [the release page in GitHub](https://github.com/mtneug/spate/releases/latest). But the easiest way to run `spate` is, of course, as a Docker Swarm service with the [`mtneug/spate` image](https://hub.docker.com/r/mtneug/spate/):

```sh
$ docker service create \
    --constraint 'node.role == manager' \
    --mount 'type=bind,src=/var/run/docker.sock,dst=/var/run/docker.sock' \
    mtneug/spate
```

`spate` sets the replica count via the Docker API and needs therefore access to the Unix socket of a manager node. Also, make sure to put `spate` in the necessary networks so that it can poll data from the running containers is should scale. Prometheus metrics are exposed on port 8080.

## Quick Start

After `spate` is running a reconciliation loop constantly looks for changes in the Docker Swarm cluster. Autoscaling is entirely controlled through service labels. In this way users and scripts only have to directly interact with the Docker CLI client.

```sh
$ docker service create \
    --name my-worker \
    --network spate-demo \
    --label "de.mtneug.spate=enable" \
    --label "de.mtneug.spate.metric.load.type=prometheus" \
    --label "de.mtneug.spate.metric.load.kind=replica" \
    --label "de.mtneug.spate.metric.load.prometheus.endpoint=http://localhost:8080/metrics" \
    --label "de.mtneug.spate.metric.load.prometheus.name=jobs_total" \
    --label "de.mtneug.spate.metric.load.target=3" \
    my/worker
```

In this example a `my-worker` service is created. Each replica exposes the `jobs_total` Prometheus metric on port 8080, which should count the number of jobs currently processed. It is the aim, that each replica, in average, processes three jobs. A metric is therefore added with the `de.mtneug.spate.metric.load` labels. This way, if too many jobs are in the system, `spate` automatically increase the replica count and vice versa.

`spate` is highly configurable. Have a look at the [documentation](doc/README.md) for a complete list of options.

## Example

An example application is provided with [spate-demo](https://github.com/mtneug/spate-demo). Consult this repository's README for more information.

## License

Apache 2.0 (c) Matthias Neugebauer
