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

        # server_ec2 = CustomEC2(self, "server_ec2", vpc, sec_group, "test-server", ["server_asset", "server_cmd.sh"], ["server_asset", "hw4.zip"])
        
        db_ec2 = CustomEC2(self, "databse_ec2", vpc, sec_group, "test-sql-server", ["db_asset", "sql_setup.sh"])
