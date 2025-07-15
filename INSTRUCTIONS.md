# Overview

The paper makes the following claims on pg. 2 (Comments to AEC reviewers after `:`):

**(C1) Slowpoke accurately quantifies throughput optimizations on four real-world microservice applications.**

**(C2) Slowpoke accurately quantifies throughput after scaling optimizations or when the bottleneck is caused by mutex contention.**

**(C3) Slowpoke accurately quantifies throguhput optimizations across synthetic microservice applications covering a wide range of microservice characteristics.**

This artifact targets the following badges:

* [ ] [Artifact available](#artifact-available): Reviewers are expected to confirm that Slowpoke system, benchmarks, and testing scripts are all publicly available (about 10 minutes).
* [ ] [Artifact functional](#artifact-functional): Reviewers are expected to confirm sufficient documentation, key components as described in the paper, and execution with one experiment (about 20 minutes).
* [ ] [Results reproducible](#results-reproducible): Reviewers are expected to reproduce _key_ results of section 5 of the paper (3 hours).

# Artifact Available (10 minutes)

Confirm that the benchmark programs, their inputs, and automation scripts are all publicly available:

1. The code is hosted at: [https://github.com/atlas-brown/slowpoke](https://github.com/atlas-brown/slowpoke)

2. Slowpoke is available in the [`pkg/slowpoke/`]() directory of the repository, including a third-party lib and Poker runtime.

3. Benchmarks are available in the [`benchmarks/`]() directory of the repository.

4. Additional scripts are available in the [`scripts/`]() directory of the repository.

> AEC Reviewers: From this point on, scripts use the provided AWS EC2 instances. All preprofiling results, Docker images are provided for efficiency.
> We provide a kubernete cluster with X machines for each reviwer. To request access to the control node, please send a request to slowpoke@brown.edu. Once the access is granted, reviwers can start/stop the clusters whenever.

# Artifact Functional (20 minutes)

Confirm sufficient documentation, key components as described in the paper, and execution with min inputs (about 20 minutes):

<!-- * Exercisability: Instructions below set up an Debian-based container and run _all_ benchmarks on `min` inputs (`all.sh`) or run specific benchmarks. -->

* Documentation:
* Key components:
* Benchmarks:

> At this point, **run `git clone https://github.com/atlas-brown/slowpoke`**

**Quickstart: Running a demo experiment (e.g., `online-boutique`):** To quickly execute a specific benchmark such as `online-boutique`... (2 minutes):

```sh
./script/demo.sh
```

# Results Reproducible (about 3 hours)