#!/usr/bin/env python3
import os

import aws_cdk as cdk

from cdk.cdk_stack import CdkStack

env=cdk.Environment(account='', region='us-west-2')

app = cdk.App()
CdkStack(app, "CdkStack", env=env)

app.synth()
