# Slowpoke: End-to-end Throughput Optimization Modeling for Microservice Applications

Slowpoke is a causal profiling tool for microservice applications. It answers what-if throughput questions by predicting impact of service-level throughput optimization on end-to-end application-level throughput.

## Citing Slowpoke
```
@inproceedings{slowpoke:nsdi:2026,
  author = {Yizheng Xie and Di Jin and Oğuzhan Çölkesen and Vasiliki Kalavri and John Liagouris and Nikos Vasilakis},
  title = {Slowpoke: End-to-end Throughput Optimization Modeling for Microservice Applications},
  booktitle = {23rd USENIX Symposium on Networked Systems Design and Implementation (NSDI 26)},
  year = {2026},
  address = {Renton, WA},
  publisher = {USENIX Association},
  month = may
}
```
## Requirement
Using Slowpoke requires a Kubernetes cluster.

<details><summary>(Optional) AWS Setup</summary>
AWS Cluster setup for different account and environment can be quite different. The following setup procedure is just for reference. We do not claim our infrastructure for AWS EC2 setup to be production-ready or reusable as is.

1. Create AWS access token from AWS account.
2. Install `aws-cli` command line tool and set the access token appropriately with aws configure. Also set region explicitly.
3. Modify the `IMAGE_ID` global variable in `scripts/setup/ec2_cluster.py` to be the Ubuntu image in the region you selected (For example, it is `ami-0d1b5a8c13042c939` for us-east-2)
4. Create an empty directory for storing cluster information (For example `~/mycluster`)
5. Run cluster setup with `python3 ./scripts/setup/setup_ec2_cluster.py -d ~/mycluster/ -n 12`. Replace 12 with desired worker count.
6. After the script finishes, run `./scripts/setup/initialize-aws.sh ~/mycluster` to setup kubernetes on each node
7. Connect to the control node by running `ssh -i ~/test_cluster/slowpoke-expr.pem ubuntu@$(head -n 1 ~/mycluster/ec2_ips)`
8. (Optional) If you wish to stop the cluster, but not delete them, run `python3 ./scripts/setup/stop_ec2_cluster.py -d ~/mycluster`
9. (Optional) Similarly, to restart `python3 ./scripts/setup/start_ec2_cluster.py -d ~/mycluster`
10. Finally, if you wish to delete the cluster, run `python3 ./scripts/setup/teardown_ec2_cluster.py -d ~/mycluster`
</details>

## Using Slowpoke

To use Slowpoke for a microservice application, changes in three places are needed.

1. Application: the current implementation of Slowpoke requires the application to be written in Go. To add in Slowpoke's runtime, modify http request handler registration to use slowpoke handler wrapper, similar to [`cmd/boutique/cart/main.go/`](cmd/boutique/cart/main.go/)
2. Container configuration change: the application process now needs to start as a subprocess to the poker controller [`slowpoke/poker/poker.c`](slowpoke/poker/poker.c). To do this, include Poker's compilation in the container configuration file as well as the entrypoint command change, similar to [`build/PrebuiltDockerfile`](build/PrebuiltDockerfile). 
3. Deployment configuration change: Slowpoke's runtime changes inserts artificial pauses based on environment variables `SLOWPOKE_DELAY_MICROS`, `SLOWPOKE_PRERUN`, `SLOWPOKE_POKER_BATCH_THRESHOLD`, and `SLOWPOKE_IS_TARGET_SERVICE`. Set these values accordingly when launching a deployment, or programatically pass in the values similar to [`slowpoke/boutique/yamls/cart.yaml`](slowpoke/boutique/yamls/cart.yaml).
4. Run slowpoke's end-to-end profiling by running `test.py`
```console
$ python3 slowpoke/test.py --help
usage: test.py [-h] [-b BENCHMARK] [-t NUM_THREADS] [-c NUM_CONNS] [-x TARGET_SERVICE] [-r REQUEST_TYPE]
               [--repetitions REPETITIONS] [--num_exp NUM_EXP] [--pre_run] [--range RANGE RANGE] [--num_req NUM_REQ]
               [--clien_cpu_quota CLIEN_CPU_QUOTA] [--random_seed RANDOM_SEED] [--poker_batch POKER_BATCH]
               [--poker_batch_req POKER_BATCH_REQ] [--poker_relative_batch POKER_RELATIVE_BATCH]

options:
  -h, --help            show this help message and exit
  -b BENCHMARK, --benchmark BENCHMARK
                        the benchmark to run the experiment on
  -t NUM_THREADS, --num_threads NUM_THREADS
                        the number of threads to run the experiment on
  -c NUM_CONNS, --num_conns NUM_CONNS
                        the number of connections to run the experiment on
  -x TARGET_SERVICE, --target_service TARGET_SERVICE
                        the target service to run the experiment on
  -r REQUEST_TYPE, --request_type REQUEST_TYPE
                        the request type to run the experiment on
  --repetitions REPETITIONS
                        the number of repetitions to run the experiment on
  --num_exp NUM_EXP     the number of experiments to run
  --pre_run             whether to run the experiment with prerun
  --range RANGE RANGE
  --num_req NUM_REQ
  --clien_cpu_quota CLIEN_CPU_QUOTA
  --random_seed RANDOM_SEED
  --poker_batch POKER_BATCH
  --poker_batch_req POKER_BATCH_REQ
  --poker_relative_batch POKER_RELATIVE_BATCH
```

For example
```
python3 test.py -b hotel -r mix -x profile --num_exp 10 -t 8 -c 1024 --poker_batch_req 100 --repetition 3 --num_req 50000
```

## Repo Structure

The top level repo is the Go package source code for benchmark applications, combined with Slowpoke's Go runtime [`pkg/slowpoke`](pkg/slowpoke).
Slowpoke utility and model prediction scripts are in [`slowpoke`]
