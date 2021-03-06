from aws_cdk import Duration, Stack
from aws_cdk import aws_iam as iam
from aws_cdk import aws_ec2 as ec2
from aws_cdk.aws_s3_assets import Asset

from constructs import Construct

from .custom_ec2 import CustomEC2

import os.path

dirname = os.path.dirname(__file__)


class CdkStack(Stack):
    def __init__(self, scope: Construct, construct_id: str, **kwargs) -> None:
        super().__init__(scope, construct_id, **kwargs)

        vpc = ec2.Vpc(
            self,
            "midterm-Vpc",
            cidr="10.0.0.0/16",
            max_azs=1,
            subnet_configuration=[
                ec2.SubnetConfiguration(
                    subnet_type=ec2.SubnetType.PUBLIC, name="Public", cidr_mask=24
                )
            ],
        )
        
        # ? create a new security group
        sec_group = ec2.SecurityGroup(
            self,
            "midterm-sg",
            vpc=vpc,
            allow_all_outbound=True,
        )

        sec_group.add_ingress_rule(
            peer=ec2.Peer.any_ipv4(),
            description="Allow SSH connection",
            connection=ec2.Port.tcp(22),
        )

        sec_group.add_ingress_rule(
            peer=ec2.Peer.any_ipv4(),
            description="Allow HTTP connection",
            connection=ec2.Port.tcp(80),
        )

        sec_group.add_ingress_rule(
            peer=ec2.Peer.any_ipv4(),
            description="Allow HTTPS connection",
            connection=ec2.Port.tcp(443),
        )
        
        sec_group.add_ingress_rule(
            peer=ec2.Peer.ipv4('10.0.0.0/16'),
            description="Allow internal 5432",
            connection=ec2.Port.tcp(5432),
        )
        
        sec_group.add_ingress_rule(
            peer=ec2.Peer.ipv4('10.0.0.0/16'),
            description="Allow internal connection",
            connection=ec2.Port.all_icmp()
        )
        version = "v0.6.0"
        server_ec2 = CustomEC2(self, "server_ec2", vpc, sec_group, "test-server", version, ["server_asset", "server_cmd.sh"], ["..","..","server"])
        
        db_ec2 = CustomEC2(self, "databse_ec2", vpc, sec_group, "test-sql-server", version, ["db_asset", "sql_setup.sh"])
