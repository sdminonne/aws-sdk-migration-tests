package main

import (
	"context"
	"fmt"
	"log"

	// AWS SDK v1
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"

	// AWS SDK v2
	"github.com/aws/aws-sdk-go-v2/config"
	ec2v2 "github.com/aws/aws-sdk-go-v2/service/ec2"
)

// This example demonstrates using both SDK v1 and v2 in the same application.
// We'll use v1 for EC2 operations and v2 for the same EC2 operations to compare.
func main() {
	fmt.Println("=== Mixed SDK Test: EC2 with v1 and v2 ===\n")

	// Initialize SDK v1 for EC2
	fmt.Println("1. Initializing AWS SDK v1 for EC2...")
	sessV1, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		log.Fatalf("Failed to create v1 session: %v", err)
	}
	ec2ClientV1 := ec2.New(sessV1)
	fmt.Println("   ✓ SDK v1 session and EC2 client created")

	// Initialize SDK v2 for EC2
	fmt.Println("\n2. Initializing AWS SDK v2 for EC2...")
	ctx := context.Background()
	cfgV2, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("Failed to load v2 config: %v", err)
	}
	ec2ClientV2 := ec2v2.NewFromConfig(cfgV2)
	fmt.Println("   ✓ SDK v2 config and EC2 client created")

	// Use v1 to list EC2 instances
	fmt.Println("\n3. Using SDK v1 to list EC2 instances...")
	instancesV1, err := ec2ClientV1.DescribeInstances(&ec2.DescribeInstancesInput{})
	if err != nil {
		log.Printf("   ✗ Failed to list instances with v1: %v", err)
	} else {
		count := 0
		for _, reservation := range instancesV1.Reservations {
			count += len(reservation.Instances)
		}
		fmt.Printf("   ✓ Found %d EC2 instances using SDK v1\n", count)
		shown := 0
		for _, reservation := range instancesV1.Reservations {
			for _, instance := range reservation.Instances {
				if shown < 3 {
					name := "N/A"
					for _, tag := range instance.Tags {
						if aws.StringValue(tag.Key) == "Name" {
							name = aws.StringValue(tag.Value)
							break
						}
					}
					fmt.Printf("     - %s (State: %s, Type: %s)\n",
						aws.StringValue(instance.InstanceId),
						aws.StringValue(instance.State.Name),
						aws.StringValue(instance.InstanceType))
					if name != "N/A" {
						fmt.Printf("       Name: %s\n", name)
					}
					shown++
				}
			}
		}
		if count > 3 {
			fmt.Printf("     ... and %d more\n", count-3)
		}
	}

	// Use v1 to list VPCs
	fmt.Println("\n4. Using SDK v1 to list VPCs...")
	vpcsV1, err := ec2ClientV1.DescribeVpcs(&ec2.DescribeVpcsInput{})
	if err != nil {
		log.Printf("   ✗ Failed to list VPCs with v1: %v", err)
	} else {
		fmt.Printf("   ✓ Found %d VPCs using SDK v1\n", len(vpcsV1.Vpcs))
		for i, vpc := range vpcsV1.Vpcs {
			if i < 3 {
				name := "N/A"
				for _, tag := range vpc.Tags {
					if aws.StringValue(tag.Key) == "Name" {
						name = aws.StringValue(tag.Value)
						break
					}
				}
				fmt.Printf("     - %s (CIDR: %s, Default: %v)\n",
					aws.StringValue(vpc.VpcId),
					aws.StringValue(vpc.CidrBlock),
					aws.BoolValue(vpc.IsDefault))
				if name != "N/A" {
					fmt.Printf("       Name: %s\n", name)
				}
			}
		}
		if len(vpcsV1.Vpcs) > 3 {
			fmt.Printf("     ... and %d more\n", len(vpcsV1.Vpcs)-3)
		}
	}

	// Use v1 to list Subnets
	fmt.Println("\n5. Using SDK v1 to list Subnets...")
	subnetsV1, err := ec2ClientV1.DescribeSubnets(&ec2.DescribeSubnetsInput{})
	if err != nil {
		log.Printf("   ✗ Failed to list subnets with v1: %v", err)
	} else {
		fmt.Printf("   ✓ Found %d Subnets using SDK v1\n", len(subnetsV1.Subnets))
		for i, subnet := range subnetsV1.Subnets {
			if i < 3 {
				name := "N/A"
				for _, tag := range subnet.Tags {
					if aws.StringValue(tag.Key) == "Name" {
						name = aws.StringValue(tag.Value)
						break
					}
				}
				fmt.Printf("     - %s (VPC: %s, CIDR: %s, AZ: %s)\n",
					aws.StringValue(subnet.SubnetId),
					aws.StringValue(subnet.VpcId),
					aws.StringValue(subnet.CidrBlock),
					aws.StringValue(subnet.AvailabilityZone))
				if name != "N/A" {
					fmt.Printf("       Name: %s\n", name)
				}
			}
		}
		if len(subnetsV1.Subnets) > 3 {
			fmt.Printf("     ... and %d more\n", len(subnetsV1.Subnets)-3)
		}
	}

	// Use v2 to list EC2 instances
	fmt.Println("\n6. Using SDK v2 to list EC2 instances...")
	instancesV2, err := ec2ClientV2.DescribeInstances(ctx, &ec2v2.DescribeInstancesInput{})
	if err != nil {
		log.Printf("   ✗ Failed to list instances with v2: %v", err)
	} else {
		count := 0
		for _, reservation := range instancesV2.Reservations {
			count += len(reservation.Instances)
		}
		fmt.Printf("   ✓ Found %d EC2 instances using SDK v2\n", count)
		shown := 0
		for _, reservation := range instancesV2.Reservations {
			for _, instance := range reservation.Instances {
				if shown < 3 {
					name := "N/A"
					for _, tag := range instance.Tags {
						if tag.Key != nil && *tag.Key == "Name" && tag.Value != nil {
							name = *tag.Value
							break
						}
					}
					stateStr := "unknown"
					if instance.State != nil && instance.State.Name != "" {
						stateStr = string(instance.State.Name)
					}
					typeStr := "unknown"
					if instance.InstanceType != "" {
						typeStr = string(instance.InstanceType)
					}
					instanceID := "N/A"
					if instance.InstanceId != nil {
						instanceID = *instance.InstanceId
					}
					fmt.Printf("     - %s (State: %s, Type: %s)\n", instanceID, stateStr, typeStr)
					if name != "N/A" {
						fmt.Printf("       Name: %s\n", name)
					}
					shown++
				}
			}
		}
		if count > 3 {
			fmt.Printf("     ... and %d more\n", count-3)
		}
	}

	// Use v2 to list VPCs
	fmt.Println("\n7. Using SDK v2 to list VPCs...")
	vpcsV2, err := ec2ClientV2.DescribeVpcs(ctx, &ec2v2.DescribeVpcsInput{})
	if err != nil {
		log.Printf("   ✗ Failed to list VPCs with v2: %v", err)
	} else {
		fmt.Printf("   ✓ Found %d VPCs using SDK v2\n", len(vpcsV2.Vpcs))
		for i, vpc := range vpcsV2.Vpcs {
			if i < 3 {
				name := "N/A"
				for _, tag := range vpc.Tags {
					if tag.Key != nil && *tag.Key == "Name" && tag.Value != nil {
						name = *tag.Value
						break
					}
				}
				vpcID := "N/A"
				if vpc.VpcId != nil {
					vpcID = *vpc.VpcId
				}
				cidr := "N/A"
				if vpc.CidrBlock != nil {
					cidr = *vpc.CidrBlock
				}
				isDefault := false
				if vpc.IsDefault != nil {
					isDefault = *vpc.IsDefault
				}
				fmt.Printf("     - %s (CIDR: %s, Default: %v)\n", vpcID, cidr, isDefault)
				if name != "N/A" {
					fmt.Printf("       Name: %s\n", name)
				}
			}
		}
		if len(vpcsV2.Vpcs) > 3 {
			fmt.Printf("     ... and %d more\n", len(vpcsV2.Vpcs)-3)
		}
	}

	// Use v2 to list Subnets
	fmt.Println("\n8. Using SDK v2 to list Subnets...")
	subnetsV2, err := ec2ClientV2.DescribeSubnets(ctx, &ec2v2.DescribeSubnetsInput{})
	if err != nil {
		log.Printf("   ✗ Failed to list subnets with v2: %v", err)
	} else {
		fmt.Printf("   ✓ Found %d Subnets using SDK v2\n", len(subnetsV2.Subnets))
		for i, subnet := range subnetsV2.Subnets {
			if i < 3 {
				name := "N/A"
				for _, tag := range subnet.Tags {
					if tag.Key != nil && *tag.Key == "Name" && tag.Value != nil {
						name = *tag.Value
						break
					}
				}
				subnetID := "N/A"
				if subnet.SubnetId != nil {
					subnetID = *subnet.SubnetId
				}
				vpcID := "N/A"
				if subnet.VpcId != nil {
					vpcID = *subnet.VpcId
				}
				cidr := "N/A"
				if subnet.CidrBlock != nil {
					cidr = *subnet.CidrBlock
				}
				az := "N/A"
				if subnet.AvailabilityZone != nil {
					az = *subnet.AvailabilityZone
				}
				fmt.Printf("     - %s (VPC: %s, CIDR: %s, AZ: %s)\n", subnetID, vpcID, cidr, az)
				if name != "N/A" {
					fmt.Printf("       Name: %s\n", name)
				}
			}
		}
		if len(subnetsV2.Subnets) > 3 {
			fmt.Printf("     ... and %d more\n", len(subnetsV2.Subnets)-3)
		}
	}

	fmt.Println("\n=== Conclusion ===")
	fmt.Println("✓ Both SDKs work independently in the same application")
	fmt.Println("✓ Each SDK maintains its own session/config")
	fmt.Println("✓ Both SDKs can authenticate using the same AWS credentials")
	fmt.Println("\nKey differences between v1 and v2:")
	fmt.Println("  - v1 uses pointers extensively (aws.String, aws.StringValue)")
	fmt.Println("  - v2 uses native types and requires explicit nil checks")
	fmt.Println("  - v2 requires context.Context for all operations")
	fmt.Println("  - v2 uses strongly-typed enums instead of string pointers")
	fmt.Println("\nThis demonstrates that you can gradually migrate services")
	fmt.Println("from v1 to v2 without having to migrate everything at once.")
}
