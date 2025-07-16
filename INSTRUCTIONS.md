# Overview

The paper makes the following claims on pg. 2 (Comments to AEC reviewers after `:`):

* **Slowpoke**: the system for accurate what-if analysis of throughput optimizations in complex microservice architectures.

* **A performance model**: a formal model that explains why Slowpoke can accurately predict end-to-end throughput improvements.

* **A lightweight distributed slowdown mechanism**: a distributed protocol improving prediction accuracy by coordinating pauses across microservices.


<!-- **(C1) Slowpoke accurately quantifies throughput optimizations on four real-world microservice applications.** -->

<!-- **(C2) Slowpoke accurately quantifies throughput after scaling optimizations or when the bottleneck is caused by mutex contention.** -->

<!-- **(C2) Slowpoke accurately quantifies throguhput optimizations across synthetic microservice applications covering a wide range of microservice characteristics.** -->

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
> We provide a kubernete cluster with X machines for each reviwer, with all dependencies satisfied.
> To request access to the control node, please comment your public keys on hotcrp. 
> Once the access is granted, reviwers can start/stop the clusters as needed.

# Artifact Functional (20 minutes)

Confirm sufficient documentation, key components as described in the paper, and execution with minimal inputs (approximately 20 minutes):

* Documentation: The top-level [README]() file provides instructions for setting up Kubernetes clusters, installing dependencies, building application images with Slowpoke, generating synthetic benchmarks, and running experiments. 
* Completeness: 
  * Slowpoke: [User-level library]() and [Poker runtime]().
  * Four real-world benchmarks (i.e., [hotel-res](), [online-boutique](), [social-net](), and [media-review]()) and a [synthetic benchmark emulator]() that dynamically changes behavior based on configuration files, e.g., [108 example configuration files]().
* **Exercisability**:
  * **EC2 cluster setup**: We provide automation scripts and instructions in `scripts/setup/` to create, initialize, start, stop, and terminate EC2 clusters.
  * **Building and deploying applications with Slowpoke**: Instructions are available in `scripts/build/` for instrumenting applications and deploying them with Slowpoke, including modifying YAML configuration files.
  * **Automated testing framework**: Scripts in `scripts/test/` support end-to-end experiment orchestration with Slowpoke, enabling reproducible and automated testing.


> At this point, **run `git clone https://github.com/atlas-brown/slowpoke`**

**Setup (optional):** We provide fully initialized and configured Kubernetes clusters for convenience.  
Optionally, reviewers may set up their own EC2 machines using the following script:
```
# Create a EC2 cluster with 8 worker nodes and one control node
python3 script/setup/ec2_cluster -n 8
# Initialize the Kubernete clusters and return the control node IP
IP=$(./script/setup/initialize-aws.sh)
# ssh into the control node
ssh -i slowpoke.pem ubuntu@$IP
```

**Quickstart:** To quickly execute a demo experiment, invoke the top-level `main.sh` script with the `--demo` flag after `ssh` into the control node. This will predict the throughput of the `online-boutique` benchmark after optimizing the execution time of the `productCatalog` service by 1 ms and compare the result against the ground truth.

```sh
./main.sh --demo
```

# Results Reproducible (about 2 hours)
The key results of Slowpoke's accuracy are the following: 

**(C1, §5.1, Fig.8) Across four real-world benchmarks**

<!-- **(C2, §5.3, Fig.11) After scaling optimizations or when the bottleneck is caused by mutex contention** -->

**(C2, §5.2, Fig.9) Across synthetic microservice applications**

The results presented in the paper are based on 10 optimizations, each reducing the target service's processing time by increments ranging from 10\% to 100\%.  
Executing the full set of experiments takes several days to complete.  
To enable more efficient reproduction without loss of insight, we sample 5 optimizations from the same range.  

**C1 (X minutes):** Invoke the top-level `main.sh` script with the `--real-world` flag. This will run experiments for both predictions and ground truth measurements, and generate Figure 7 using 5 optimizations.
```bash
./main.sh --real-world
```

The results in the paper are generated from 9 synthetic topologies (Fig. 7), each evaluated under three different configuration parameters, resulting in a total of 108 applications.  
This exhaustive exploration is time-consuming. 
To reduce evaluation time, we sample one application from each topology.  
**C2 (X minutes):** Invoke the top-level `main.sh` script with the `--synthetic` flag. 
This runs experiments for both predictions and ground truth, and outputs Figure 9 with 45 data points.  

```bash
./main.sh --synthetic
```

# Optional: Applying Slowpoke to All Benchmarks (2–3 days)

We also provide scripts to reproduce all experiments from the paper using the full set of optimization samples, including the scenarios where the bottleneck is due to mutex contention (as described in §5.3).  
To generate the complete set of results, invoke the main script for each case with the `--full` flag:

```bash
./main.sh --real-world --full
./main.sh --synthetic --full
./main.sh --mutex
