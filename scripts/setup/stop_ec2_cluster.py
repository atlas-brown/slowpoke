#!/usr/bin/env python3
import ec2_cluster
import sys

ec2_cluster.parse_args(sys.argv)
ec2_cluster.stop_ec2()

