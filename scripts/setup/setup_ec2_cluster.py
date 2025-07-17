#!/usr/bin/env python3
import ec2_cluster
import sys

ec2_cluster.parse_args(sys.argv)
ec2_cluster.setup_key()
ec2_cluster.setup_sg()
ec2_cluster.setup_ec2()
