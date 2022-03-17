from aws_cdk import Duration, Stack
from aws_cdk import aws_iam as iam
from aws_cdk import aws_ec2 as ec2
from aws_cdk.aws_s3_assets import Asset

from constructs import Construct

import os.path

dirname = os.path.dirname(__file__)


class UbuntuEC2:
    def __init__(self, scope: Construct, construct_id: str, vpc: ec2.Vpc, sec_group: ec2.SecurityGroup, instance_name: str,
                 cmd_path: list, file_path: list = None, **kwargs) -> None:

        # define a new ec2 instance
        ec2_instance = ec2.Instance(
            scope,
            f"{construct_id}-v0.2.1",
            instance_name=f"{instance_name}-v0.2.1",
            instance_type=ec2.InstanceType("t2.micro"),
            machine_image=ec2.MachineImage.generic_linux({
                'us-west-2': 'ami-0892d3c7ee96c0bf7'
            }),
            block_devices=[ec2.BlockDevice(
                device_name='/dev/sda1',
                volume=ec2.BlockDeviceVolume.ebs(10)
            )],
            vpc=vpc,
            security_group=sec_group,
        )

        # Script in S3 as Asset
        # file_asset = Asset(scope, f"{instance_name}-zip",
        #                    path=os.path.join(dirname, *file_path))
        # ec2_instance.user_data.add_s3_download_command(
        #     bucket=file_asset.bucket, bucket_key=file_asset.s3_object_key, local_file='home/ec2-user/app.zip'
        # )
        # file_asset.grant_read(ec2_instance.role)

        with open(os.path.join(dirname, *cmd_path)) as file:
            cmds = [line.rstrip() for line in file]
            ec2_instance.user_data.add_commands(*cmds)
