#!/usr/bin/env python3

'''
Require a default VPC that has internet setup and public IP by
default.  Also `aws` properly setup, e.g. `aws sts
get-caller-identity` should work properly

'''

import os
from time import sleep
from pathlib import Path
from subprocess import run, PIPE
import random
import string
import json

import argparse

OBJDIR = Path('.')
ARGS = None
WORKER_NUM = 4
IMAGE_ID = 'ami-0d1b5a8c13042c939'
SCRIPT_BASE = Path(__file__).parent
KEY_FILE = 'slowpoke-expr.pem'
KEYNAME_FILE = 'key_name'
VPCID_FILE = 'vpc_id'
SUBNET_FILE = 'subnet_id'
SG_FILE = 'sg_id'
EC2_FILE = 'instances'
IP_FILE = 'ec2_ips'
SGNAME = 'slowpokesg' + ''.join(
    random.choices('0123456789', k=10))

def parse_args(args):
    # Create ArgumentParser object
    parser = argparse.ArgumentParser(description="A simple argument parser example")

    # Add arguments
    parser.add_argument("-d", "--temp-dir", type=str, default='.', help="resource id storage basedir")
    parser.add_argument("-n", "--num", type=int, default=4, help="number of worker instances to create")
    parser.add_argument("-t", "--type", type=str, default='m5.large', help="instance type")

    print(f"[parse_args] args: {args}")

    # Parse arguments
    parsed_args = parser.parse_args(args[1:])

    global OBJDIR
    OBJDIR = Path(parsed_args.temp_dir)

    global WORKER_NUM
    WORKER_NUM = parsed_args.num

def get_vpc():
    out = run(f'aws ec2 describe-vpcs --filters Name=isDefault,Values=true '
              f'--query Vpcs[*].VpcId --output text'.split(), stdout=PIPE)
    vpcid = out.stdout.decode('ascii').strip()
    print(f'using default VPC {vpcid}')
    return vpcid

VPCID = get_vpc()

def get_s(fname):
    with open(OBJDIR / fname) as f:
        s = f.read().strip()
    return s

def save_s(fname, s, mode='w'):
    with open(OBJDIR / fname, mode) as f:
        f.write(s)
        if not s.endswith('\n'):
            f.write('\n')
        f.flush()

def des_s(fname):
    s = get_s(fname)
    os.remove(OBJDIR / fname)
    return s

def setup_key():
    keyname = 'slowpoke-expr' + ''.join(
        random.choices('0123456789', k=10))
    out = run(f'aws ec2 create-key-pair --key-name {keyname} --query KeyMaterial --output text'.split(),
              stdout=PIPE)
    save_s(KEY_FILE, out.stdout.decode('ascii'))
    save_s(KEYNAME_FILE, keyname)
    os.chmod(OBJDIR / KEY_FILE, 0o600)
    print(f'key written into {OBJDIR / KEY_FILE}')

def remove_key():
    keyname = des_s(KEYNAME_FILE)
    run(f'aws ec2 delete-key-pair --key-name {keyname}'.split())
    des_s(KEY_FILE)

def setup_sg():
    vpcid = VPCID
    out = run(f'aws ec2 create-security-group --group-name {SGNAME} '
              f'--description slowpoke_sg --vpc-id {vpcid}'.split(),
              stdout=PIPE)
    outs = out.stdout.decode('ascii')
    groupid = json.loads(outs)['GroupId']
    print(f'security group {groupid}')
    save_s(SG_FILE, groupid)
    run(f'aws ec2 authorize-security-group-ingress --group-id {groupid} '
        f'--protocol all --port -1 --cidr 0.0.0.0/0'.split(), stdout=PIPE)

def remove_sg():
    groupid = des_s(SG_FILE)
    run(f'aws ec2 delete-security-group --group-id {groupid}'.split())

def query_ec2(iid):
    out = run(f'aws ec2 describe-instances --instance-ids {iid} '
              f'--query Reservations[*].Instances[*].State.Name --output text'.split(),
              stdout=PIPE)
    stdout = out.stdout.decode('ascii').strip()
    return stdout

def query_ec2_status(iid):
    out = run(f'aws ec2 describe-instance-status --instance-ids {iid} '
              f'--query InstanceStatuses[*].InstanceStatus.Status --output text'.split(),
              stdout=PIPE)
    stdout = out.stdout.decode('ascii').strip()
    return stdout

def create_ec2_instance(num, itype, keyname, groupid, instance_name):
    out = run(f'aws ec2 run-instances --image-id {IMAGE_ID} '
              f'--count {num} --instance-type {itype} '
              f'--key-name {keyname} --security-group-ids {groupid} '
              f'--associate-public-ip-address '
              '--block-device-mappings {"DeviceName":"/dev/sda1","Ebs":{"Encrypted":false,"DeleteOnTermination":true,"Iops":3000,"VolumeSize":32,"VolumeType":"gp3","Throughput":125}} '
              f'--tag-specifications ResourceType=instance,Tags=[{{Key=Name,Value={instance_name}}}] '
              f'--query Instances[*].[InstanceId] --output text'.split(),
              stdout=PIPE)
    iids = out.stdout.decode('ascii')
    save_s(EC2_FILE, iids, "a")
    print(f'instances saved to {OBJDIR / EC2_FILE}')

def setup_ec2():
    groupid = get_s(SG_FILE)
    keyname = get_s(KEYNAME_FILE)
    foldername = OBJDIR
    print(f'[setup_ec2] creating control node and one worker node of type mx.2xlarge')
    create_ec2_instance(2, 'm5.2xlarge', keyname, groupid, f"{foldername}-control")
    num, itype = WORKER_NUM, 'm5.large'
    print(f'[setup_ec2] creating {num} workers of type {itype}')
    create_ec2_instance(num, itype, keyname, groupid, f"{foldername}-worker")
    print(f'instances saved to {OBJDIR / EC2_FILE}')
    iids = get_s(EC2_FILE).strip().split()
    for iid in iids:
        print(f'checking {iid}')
        while query_ec2(iid) != 'running':
            sleep(3)
        print('done')
    outs = []
    for iid in iids:
        out = run((f'aws ec2 describe-instances '
                   f'--query Reservations[*].Instances[*].[PublicIpAddress] '
                   f'--output text --instance-id ' + iid).split(),
                  stdout=PIPE)
        outs.append(out.stdout.decode('ascii').strip())
    save_s(IP_FILE, '\n'.join(outs))
    # print(f'ec2 ips saved to {OBJDIR / IP_FILE}')
    print('waiting for fully initalized')
    for iid in iids:
        print(f'checking {iid}')
        while query_ec2_status(iid) != 'ok':
            sleep(3)
        print('done')

def stop_ec2():
    iids = get_s(EC2_FILE).strip().split()
    run((f'aws ec2 stop-instances --instance-ids ' + ' '.join(iids)).split(), stdout=PIPE)
    for iid in iids:
        print(f'checking {iid}')
        while query_ec2(iid) != 'stopped':
            sleep(3)
        print('done')

def start_ec2():
    iids = get_s(EC2_FILE).strip().split()
    run((f'aws ec2 start-instances --instance-ids ' + ' '.join(iids)).split(), stdout=PIPE)
    for iid in iids:
        print(f'checking {iid}')
        while query_ec2(iid) != 'running':
            sleep(3)
        print('done')
    outs = []
    for iid in iids:
        out = run((f'aws ec2 describe-instances '
                   f'--query Reservations[*].Instances[*].[PublicIpAddress] '
                   f'--output text --instance-id ' + iid).split(),
                  stdout=PIPE)
        outs.append(out.stdout.decode('ascii').strip())
    save_s(IP_FILE, '\n'.join(outs))
    for iid in iids:
        print(f'checking {iid}')
        while query_ec2_status(iid) != 'ok':
            sleep(3)
        print('done')
        
def remove_ec2():
    iids = des_s(EC2_FILE).strip().split()
    run((f'aws ec2 terminate-instances --instance-ids ' + ' '.join(iids)).split(), stdout=PIPE)
    for iid in iids:
        print(f'checking {iid}')
        while query_ec2(iid) != 'terminated':
            sleep(3)
        print('done')
    try:
        des_s(IP_FILE)
    except Exception as e:
        return



if __name__ == '__main__':
    setup_key()
    setup_sg()
    setup_ec2()
    # remove_ec2()
    # remove_sg()
    # remove_key()
