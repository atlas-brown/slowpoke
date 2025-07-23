# Slowpoke: End-to-end Throughput Optimization Modeling for Microservice Applications

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
4. Run slowpoke's end-to-end profiling with XXX

### (Optional) AWS Setup

#### Application Code
The application needs to include Slowpoke's runtime as dependency, and register request handler to the Slowpoke wrapper. Currently we only support Go (See example in XX)

#### Change to Kubernetes Deployment Configuration
Slowpoke takes the Kubernetes YAML files 
