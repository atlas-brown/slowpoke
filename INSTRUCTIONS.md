# Overview

The paper makes the following contributions:

1. **Throughput predictor(§2)**: A system, Slowpoke, for predicting throughput optimizations in complex microservice architectures.
2. **Performance model (§3)**: a mathematical model underpinning the Slowpoke optimizations; this model is described in the paper.
3. **Distributed slowdown mechanism (§4)**: a service-local controller for slowing down microservices, Poker, providing runtime information to the model.

Slowpoke is characterized via (1) four real-world microservice applications and (2) synthetic microbenchmarks that cover many different micorservice configurations. The artifact focuses on claims no. 1 and 3 (the mathematical model is described the paper), and primarily evaluations with real-world applications (but validating microbenchmarks is tagged as optional, for reviewers).

This artifact targets the following badges:

* [ ] [Artifact available](#artifact-available): Reviewers are expected to confirm that Slowpoke system, benchmarks, and testing scripts are all publicly available (~10 minutes).
* [ ] [Artifact functional](#artifact-functional): Reviewers are expected to confirm sufficient documentation, key components described in the paper, and execution with one experiment (about 20 minutes).
* [ ] [Results reproducible](#results-reproducible): Reviewers are expected to reproduce key results of section 5 of the paper (about 2.5 hours for real-world benchmarks, optionally 2–3 days for the synthetic benchmarks).

> [!IMPORTANT]
> To reproduce results, this artifact uses a real distributed system on AWS. Remember:
> * Each reviewer's username and password has been shared with HotCRP, *please do not share them outside the AEC*.
> * Remember to freeze or turn off evaluation when you're done, as this evaluation is expensive! If you do not know how, ask us via HotCRP—thank you!

# Artifact Available (10 minutes)

Reviewers are expected to confirm that Slowpoke system, benchmarks, and testing scripts are all publicly available:

* Slowpoke is available at [https://github.com/atlas-brown/slowpoke](https://github.com/atlas-brown/slowpoke) (`nsdi26-ae` will be frozen).

* The [Slowpoke analysis component](pkg/slowpoke), [the Poker slowdown component](slowpoke/poker/poker.c), the [`benchmarks`](cmd/) (command-line entry points) and [`internal/`](inernal) (request handlers).

* Top-level scripts in [`slowpoke/`](slowpoke) and ones related to the artifact reproducilbility: [`run_function.sh`](run_functional.sh) and [`run_reproducible.sh`](run_reproducible.sh).

# Artifact Functional (20 minutes)

Confirm sufficient documentation, key components as described in the paper, and execution with minimal inputs:

* Documentation: The top-level [README](README.md) file provides instructions for setting up Kubernetes clusters, installing dependencies, building application images with Slowpoke, generating synthetic benchmarks, and running experiments.
 
* Completeness: (1)Slowpoke user-level library ([initialization](pkg/slowpoke/utils.go), [request handler](pkg/wrapper/wrappers.go)), (2) Poker slowdown mechanism ([Pause](pkg/slowpoke/pause.go), [poker](slowpoke/poker/poker.c)), (3) four real-world benchmarks [[1](XXX)_, [2](XXX), [3](XXX), [4](XXX)], and [108 synthetic configuration files](slowpoke/synthetic/) used with an ([emulator](cmd/synthetic/service)) that dynamically changes behavior based on configuration files.
 
* Exercisability: Instructions below access an AWS cluster via a gateaway (to allow multiple reviewers to log in at the same time without interferring with each other).

**Exercisability**: To run Slowpoke, we prepared distributed clusters on AWS. To `ssh` into AWS, replace `ae1` and  `<!PerReviewerPassword!>` with the user ID and password shared over HotCRP. 
* `ssh ae1@3.133.138.10` (use `<!PerReviewerPassword!>` when asked for a password).
* To start the cluster and ssh into the cluster control node, run `python3 ./scripts/start_ec2_cluster.py -d cluster_info` and `ssh -i ~/cluster_info/slowpoke-expr.pem`—this starts a cluster and `ssh` into it.
* To confirm it's functional, clone the repo, `cd` into it, and run `./run_functional.sh`. This will predict the throughput of the `boutique` benchmark after optimizing the execution time of the `cart` service by 1 ms and compare the result against the ground truth.
* To stop the cluster, run `exit` (to exit the cluster) and then (back into the original machine) `python3 ./scripts/stop_ec2_cluster.py -d cluster_info`

<details>
 <summary>Explaination</summary>

The cluster is already set up using scripts in this repo under [`scripts/setup/`]() (see. The cluster contains 2 AWS `m5.xlarge` and 12 `m5.large` EC2 instances. The public IPs of the EC2 machines will be stored in `~/cluster_info/ec2_ips`, first one is the kubernetes control node, the second one is worker node that runs the workload generator, the rest are worker nodes that run the services in each benchmark.

</details>

> While testing the artifact, we discovered kubernetes issue that shows up non-deterministically. Specifically, the `wrk`'s worker node occationally stop responding. If you notice the time spent in one of the steps take much longer than our estimation below, we advise going back to the gateway machine, stopping and restarting the cluster, and then trying again.

You should expect a file created at `results/boutique_tiny.log` ending with:
```console
$ tail -n 20 results/boutique_tiny.log
...
    Baseline throughput: 4008.7682284172083
    Groundtruth: [5829.637692207481]
    Slowdown:    [2545.670613587085]
    Predicted:   [5989.933238198921]
    Error Perc:  [2.7496656645696422]
```

<details>
 <summary>Explaination</summary>

`./run_functional.sh` runs [`./slowpoke/boutique/run-boutique-tiny.sh`](), which runs the main testing script with appropriate arguments

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
</details>

# Results Reproducible (about 2.5 hours)

For this step we strongly recommend the reviewers using `screen` or `tmux` to avoid accidental disconnect.
The key results of Slowpoke's accuracy is represented by the following experiments

**(§5.1, Fig.8) Across four real-world benchmarks**

The results in the paper is done with each benchmark having 10 hypothetical optimization points and 3 repetitions. To make the artifact evaluation process faster, we only run them once. The entire process should take about 2.5 hours. To run them, make sure you are in `~/slowpoke`, then

```console
$ # Recommend doing the following command in screen or tmux
$ ./run_reproducible.sh
...
The results are stored in /home/ubuntu/slowpoke/results
To visualize the results, run: 

python3 /home/ubuntu/slowpoke/draw.py /home/ubuntu/slowpoke/results
```

The log file `results/boutique_medium.log`, `results/hotel_medium.log`, `results/social_medium.log`, and `results/movie_medium.log` will grow overtime and should make progress at least once per minute.

To see the results, run (as the previous script suggested)

```console
$ python3 /home/ubuntu/slowpoke/draw.py /home/ubuntu/slowpoke/results
Result for /home/ubuntu/slowpoke/results/boutique_medium.log is available at
http://xx.xx.xx.xx/boutique_medium.png
...
```
The reviewer can expect to see plots comparing predicted throughput and groundtruth by opening the links. The log files can also be inspected directly. The reviewer should see the relative prediction error (the `Error Perc:` row in the log file) to be within 10%, and most of them centers around 0-4%. 

<details>
 <summary>
  Sample results for individual plots
 </summary>
 
We did a run on the same environment and the results are stored in [`sample_output/`](sample_output)

Boutique

![boutique](sample_output/boutique_medium.png)

Movie

![movie](sample_output/movie_medium.png)

Hotel

![hotel](sample_output/hotel_medium.png)

Social

![social](sample_output/social_medium.png)
</details>

To recreate Fig. 8 similar to the paper, run

```console
$ python3 plot_macro.py
Result for plot_macro.pdf is available at
http://xx.xx.xx.xx/plot_macro.pdf
```

Sample output
![Sample macro](sample_output/plot_macro.pdf)

# Optional: Applying Slowpoke to All Benchmarks (2–3 days)

We also provide scripts to reproduce all experiments from the paper using the full set of optimization samples, including the scenarios where the bottleneck is due to mutex contention (as described in §5.3).  
To generate the complete set of results, invoke the main script for each case with the `--full` flag:

```bash
./main.sh --real-world --full
./main.sh --synthetic --full
./main.sh --mutex
