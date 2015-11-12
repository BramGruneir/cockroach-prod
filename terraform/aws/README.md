# Deploy cockroach cluster on AWS using Terraform

This directory contains the [Terraform](https://terraform.io/) configuration
files needed to launch a cockroach cluster on AWS.

The following steps will create a three node cluster.

## One-time setup steps
1. Have an [AWS](http://aws.amazon.com/) account
2. [Download terraform](https://terraform.io/downloads.html), *version 0.6.6 or greater*, unzip, and add to your `PATH`.
3. [Valid AWS credentials file](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-set-up.html#cli-signup).
4. [Create an AWS keypair](https://console.aws.amazon.com/ec2/v2/home?region=us-east-1#KeyPairs:sort=keyName) named `cockroach` and save the file as `~/.ssh/cockroach.pem`.

## Variables

The following variables can be modified in `variables.tf` if necessary.
* `cockroach_port`: the port for the backends and load balancer
* `aws_region`: region to run in. Affects `aws_availability_zone` and `aws_ami_id`
* `aws_availability_zone`: availability zone for instances and load balancer
* `aws_ami_id`: image ID. depends on the region.
* `key_name`: base name of the AWS key
* `action`: default action. Defaults to `start`. Override is specified in initialization step

## Create the cluster

The following three commands will initialize all needed AWS infrastructure in the region `us-east-1`,
initialize the first cockroach node, then add two more nodes to the cluster.
All dynamic configuration is set through terraform command-line flags but can be set in `variables.tf`.

To see the actions expected to be performed by terraform, use `plan` instead of `apply`.

#### Initialize AWS infrastructure

```
$ terraform apply \
    --var=gossip=""                       \
    --var=num_instances=0

aws_security_group.default: Creating...
aws_security_group.default: Creation complete
aws_elb.elb: Creating...
aws_elb.elb: Creation complete
Outputs:
  elb             = elb-1289187553.us-east-1.elb.amazonaws.com
  gossip_variable = elb-1289187553.us-east-1.elb.amazonaws.com:26257
  instances       = 
  port            = 26257

```

This command creates a load balancer and displays the value of the gossip flag to be used in the
following steps.

Save the elb address and port as an environment variable:
```
$ export ELB="elb-1289187553.us-east-1.elb.amazonaws.com:26257"
```

#### Initialize the first node

```
$ terraform apply \
    --var=gossip="lb=${ELB}" \
    --var=num_instances=1                 \
    --var=action="init"

aws_security_group.default: Refreshing state... (ID: sg-828435e4)
aws_elb.elb: Refreshing state... (ID: elb)
aws_instance.cockroach: Creating...
aws_instance.cockroach: Provisioning with 'file'...
aws_instance.cockroach: Provisioning with 'file'...
aws_instance.cockroach: Provisioning with 'remote-exec'...
aws_instance.cockroach: Creation complete
aws_elb.elb: Modifying...
aws_elb.elb: Modifications complete
Outputs:
  elb             = elb-1289187553.us-east-1.elb.amazonaws.com
  gossip_variable = elb-1289187553.us-east-1.elb.amazonaws.com:26257
  instances       = ec2-54-85-12-159.compute-1.amazonaws.com
  port            = 26257

```

The `--var=action="init"` parameter causes the first node to be initialized for with a new cluster.
The cluster is now running with a single node and is reachable through the load balancer (see `Using the cluster`).

#### Add more nodes to the cluster

```
$ terraform apply \
    --var=gossip="lb=${ELB}" \
    --var=num_instances=3

aws_security_group.default: Refreshing state... (ID: sg-828435e4)
aws_instance.cockroach.0: Refreshing state... (ID: i-1d10fbca)
aws_elb.elb: Refreshing state... (ID: elb)
aws_instance.cockroach.1: Creating...
aws_instance.cockroach.2: Creating...
aws_instance.cockroach.2: Provisioning with 'file'...
aws_instance.cockroach.1: Provisioning with 'file'...
aws_instance.cockroach.2: Provisioning with 'remote-exec'...
aws_instance.cockroach.2: Creation complete
aws_instance.cockroach.1: Provisioning with 'remote-exec'...
aws_instance.cockroach.1: Creation complete
aws_elb.elb: Modifying...
aws_elb.elb: Modifications complete
Outputs:
  elb             = elb-1289187553.us-east-1.elb.amazonaws.com
  gossip_variable = elb-1289187553.us-east-1.elb.amazonaws.com:26257
  instances       = ec2-54-85-12-159.compute-1.amazonaws.com,ec2-54-175-192-198.compute-1.amazonaws.com,ec2-54-88-84-13.compute-1.amazonaws.com
  port            = 26257

```

## Use the cluster

#### Connect to the cluster

Use the load balancer address to connect to the cluster. You may need to wait a few minutes after
ELB creation for its DNS name to be resolvable.

```
$ ./cockroach sql --insecure --addr=${ELB}
elb-1289187553.us-east-1.elb.amazonaws.com:26257> show databases;
+----------+
| Database |
+----------+
| system   |
+----------+
```

#### Ssh into individual instances

The DNS names of AWS instances is shown as a comma-separated list in the terraform output.

```
$ ssh -i ~/.ssh/cockroach.pem ubuntu@ec2-54-85-12-159.compute-1.amazonaws.com

ubuntu@ip-172-31-15-87:~$ ps -Af|grep cockroach
ubuntu    1448     1  4 20:03 ?        00:00:39 ./cockroach start --log-dir=logs --logtostderr=false --stores=ssd=data --insecure --gossip=lb=${ELB}

ubuntu@ip-172-31-15-87:~$ ls logs
cockroach.ERROR
cockroach.INFO
cockroach.ip-172-31-15-87.ubuntu.log.ERROR.2015-11-02T20_03_14Z.1448
cockroach.ip-172-31-15-87.ubuntu.log.INFO.2015-11-02T20_03_09Z.1443
cockroach.ip-172-31-15-87.ubuntu.log.INFO.2015-11-02T20_03_09Z.1448
cockroach.ip-172-31-15-87.ubuntu.log.WARNING.2015-11-02T20_03_09Z.1448
cockroach.STDERR
cockroach.STDOUT
cockroach.WARNING

```

#### Profile servers

Using either the ELB address (will hit a random node), or a specific instance:
```
$ go tool pprof <address:port>/debug/pprof/profile
```

#### Running examples against the cockroach cluster

See `examples.tf` for sample examples and how to run them against the created cluster.

## Destroy the cluster

```
$ terraform destroy \
    --var=gossip="lb=${ELB}" \
    --var=num_instances=3
```

The destroy command requires confirmation.
