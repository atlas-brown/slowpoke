# Slowpoke: End-to-end Throughput Optimization Modeling for Microservice Applications

## Using Slowpoke

To use Slowpoke for a microservice application, changes in three places are needed

### (Optional) AWS Setup
AWS Cluster setup for different account and environment can be quite different. The following setup procedure is just for reference. We do **not** claim our infrastructure for AWS EC2 setup to be production-ready or reusable as is.

1. Create AWS access token from AWS account.
2. Install `aws-cli` command line tool and set the access token appropriately with `aws configure`. Also set region explicitly.
3. Modify the `IMAGE_ID` global variable in `scripts/setup/ec2_cluster.py` to be the Ubuntu image in the region you selected (For example, it is  `ami-0d1b5a8c13042c939` for `us-east-2`)
4. Create an empty directory for storing cluster information (For example `~/mycluster`)
5. Run cluster setup with
`python3 ./scripts/setup/setup_ec2_cluster.py -d ~/mycluster/ -n 12`
Replace 12 with desired worker count.
6. After the script finishes, run `./scripts/setup/initialize-aws.sh ~/mycluster` to setup kubernetes on each node
7. Connect to the control node by running `ssh -i ~/test_cluster/slowpoke-expr.pem ubuntu@$(head -n 1 ~/mycluster/ec2_ips)`
8. (Optional) If you wish to stop the cluster, but not delete them, run `python3 ./scripts/setup/stop_ec2_cluster.py -d ~/mycluster`
9. (Optional) Similarly, to restart `python3 ./scripts/setup/start_ec2_cluster.py -d ~/mycluster`
10. Finally, if you wish to delete the cluster, run `python3 ./scripts/setup/teardown_ec2_cluster.py -d ~/mycluster`

#### Application Code
The application needs to include Slowpoke's runtime as dependency, and register request handler to the Slowpoke wrapper. Currently we only support Go (See example in XX)

#### Change to Kubernetes Deployment Configuration
Slowpoke takes the Kubernetes YAML files 
