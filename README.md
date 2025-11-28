# AWS SDK Migration Tests

A collection of test programs demonstrating AWS SDK v1 to v2 migration patterns and compatibility.

## Overview

This project provides practical examples showing how to migrate from AWS SDK for Go v1 to v2. It demonstrates that:

- Infrastructure created with SDK v1 can be fully managed with SDK v2 (and vice versa)
- Both SDKs can coexist in the same application
- AWS resources are SDK-agnostic and exist independently of the SDK version used to create them
- You can migrate your codebase incrementally without recreating existing infrastructure

## Test Programs

### 1. cross_version_infrastructure

Demonstrates cross-version infrastructure compatibility using S3.

**What it does:**
- Creates an S3 bucket using SDK v1
- Lists and manages the bucket using SDK v2
- Puts objects with v2 into the v1-created bucket
- Verifies changes are visible back in v1
- Cleans up resources

**Key takeaway:** Resources created with one SDK version are fully accessible and manageable by the other version.

### 2. mixed_sdk

Demonstrates running both SDKs side-by-side in the same application.

**What it does:**
- Initializes both v1 and v2 clients for EC2
- Lists EC2 instances, VPCs, and Subnets using v1
- Lists the same resources using v2
- Compares the results and highlights API differences

**Key takeaway:** Both SDKs can work independently in the same application, allowing for gradual migration.

## Prerequisites

- Go 1.24 or later
- AWS credentials configured (via environment variables, shared credentials file, or IAM role)
- Appropriate AWS permissions for the operations being tested

## Building

Build all binaries:
```bash
make all
```

Build individual binaries:
```bash
make cross_version    # Build cross_version_infrastructure
make mixed_sdk        # Build mixed_sdk
```

## Running

Run the cross-version infrastructure test:
```bash
./cross_version_infrastructure
```

Run the mixed SDK test:
```bash
./mixed_sdk
```

## AWS Credentials

These programs require valid AWS credentials. Configure them using one of these methods:

1. Environment variables:
   ```bash
   export AWS_ACCESS_KEY_ID=your_access_key
   export AWS_SECRET_ACCESS_KEY=your_secret_key
   export AWS_REGION=us-east-1
   ```

2. Shared credentials file (`~/.aws/credentials`)

3. IAM role (if running on EC2)

## Required Permissions

### For cross_version_infrastructure:
- `s3:CreateBucket`
- `s3:DeleteBucket`
- `s3:ListBuckets`
- `s3:PutObject`
- `s3:ListObjects`
- `s3:GetBucketLocation`

### For mixed_sdk:
- `ec2:DescribeInstances`
- `ec2:DescribeVpcs`
- `ec2:DescribeSubnets`

## Key Differences Between SDK v1 and v2

| Aspect | SDK v1 | SDK v2 |
|--------|--------|--------|
| Pointer handling | Uses `aws.String()`, `aws.StringValue()` extensively | Uses native types, requires explicit nil checks |
| Context | Optional | Required for all operations |
| Enums | String pointers | Strongly-typed enums |
| Configuration | Session-based | Config-based |
| Error handling | Standard Go errors | More detailed error types |

## Development

Run tests:
```bash
make test
```

Clean build artifacts:
```bash
make clean
```

Get help:
```bash
make help
```

## Project Structure

```
.
├── cross_version_infrastructure.go  # Cross-version compatibility test
├── mixed_sdk.go                     # Side-by-side SDK comparison
├── Makefile                         # Build automation
├── go.mod                           # Go module dependencies
└── README.md                        # This file
```

## Dependencies

- [AWS SDK for Go v1](https://github.com/aws/aws-sdk-go) (v1.55.8)
- [AWS SDK for Go v2](https://github.com/aws/aws-sdk-go-v2) (latest)

## License

This is a test project for demonstrating AWS SDK migration patterns.
