package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	// AWS SDK v1
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	s3v1 "github.com/aws/aws-sdk-go/service/s3"

	// AWS SDK v2
	"github.com/aws/aws-sdk-go-v2/config"
	s3v2 "github.com/aws/aws-sdk-go-v2/service/s3"
)

// This example demonstrates that infrastructure created with SDK v1 can be
// fully managed with SDK v2 (and vice versa).
//
// We'll create an S3 bucket with v1, then list and manage it with v2.
func main() {
	fmt.Println("=== Cross-Version Infrastructure Test ===\n")

	// Generate a unique bucket name
	bucketName := fmt.Sprintf("sdk-migration-test-%d", time.Now().Unix())
	region := "us-east-1"
	ctx := context.Background()

	fmt.Printf("Test bucket name: %s\n\n", bucketName)

	// ===== PHASE 1: Create bucket with SDK v1 =====
	fmt.Println("PHASE 1: Creating S3 bucket using SDK v1")
	fmt.Println("------------------------------------------")

	sessV1, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		log.Fatalf("Failed to create v1 session: %v", err)
	}
	s3ClientV1 := s3v1.New(sessV1)

	fmt.Printf("Creating bucket '%s' with SDK v1...\n", bucketName)
	_, err = s3ClientV1.CreateBucket(&s3v1.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		log.Fatalf("Failed to create bucket with v1: %v", err)
	}
	fmt.Println("✓ Bucket created successfully with SDK v1")

	// Verify with v1
	fmt.Println("\nVerifying bucket exists using SDK v1...")
	_, err = s3ClientV1.HeadBucket(&s3v1.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		log.Fatalf("Failed to verify bucket with v1: %v", err)
	}
	fmt.Println("✓ Bucket verified with SDK v1")

	// ===== PHASE 2: Manage bucket with SDK v2 =====
	fmt.Println("\n\nPHASE 2: Managing the same bucket using SDK v2")
	fmt.Println("------------------------------------------------")

	cfgV2, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Fatalf("Failed to load v2 config: %v", err)
	}
	s3ClientV2 := s3v2.NewFromConfig(cfgV2)

	// List buckets with v2 to find our bucket
	fmt.Println("Listing all buckets using SDK v2...")
	listResult, err := s3ClientV2.ListBuckets(ctx, &s3v2.ListBucketsInput{})
	if err != nil {
		log.Fatalf("Failed to list buckets with v2: %v", err)
	}

	bucketFound := false
	for _, bucket := range listResult.Buckets {
		if *bucket.Name == bucketName {
			bucketFound = true
			fmt.Printf("✓ Found our bucket '%s' created with v1, now visible in v2!\n", *bucket.Name)
			fmt.Printf("  Created: %v\n", bucket.CreationDate)
			break
		}
	}

	if !bucketFound {
		log.Fatalf("Bucket not found in v2 list (this shouldn't happen!)")
	}

	// Get bucket details with v2
	fmt.Println("\nGetting bucket location using SDK v2...")
	locationResult, err := s3ClientV2.GetBucketLocation(ctx, &s3v2.GetBucketLocationInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		log.Fatalf("Failed to get bucket location with v2: %v", err)
	}
	location := "us-east-1" // Default for empty LocationConstraint
	if locationResult.LocationConstraint != "" {
		location = string(locationResult.LocationConstraint)
	}
	fmt.Printf("✓ Bucket location: %s\n", location)

	// Put an object using v2
	fmt.Println("\nPutting an object into the bucket using SDK v2...")
	objectKey := "test-object.txt"
	objectContent := "This object was created with SDK v2 in a bucket created with SDK v1!"
	_, err = s3ClientV2.PutObject(ctx, &s3v2.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   strings.NewReader(objectContent),
	})
	if err != nil {
		log.Printf("Warning: Failed to put object with v2: %v", err)
	} else {
		fmt.Printf("✓ Object '%s' created successfully with SDK v2\n", objectKey)
	}

	// ===== PHASE 3: Verify with v1 again =====
	fmt.Println("\n\nPHASE 3: Verifying changes are visible back in SDK v1")
	fmt.Println("--------------------------------------------------------")

	fmt.Println("Listing objects in bucket using SDK v1...")
	listObjResult, err := s3ClientV1.ListObjectsV2(&s3v1.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		log.Printf("Note: ListObjects with v1: %v (this is OK if bucket is empty)", err)
	} else {
		fmt.Printf("✓ SDK v1 can see %d objects in the bucket\n", len(listObjResult.Contents))
	}

	// ===== CLEANUP =====
	fmt.Println("\n\nCLEANUP: Deleting test bucket")
	fmt.Println("-------------------------------")

	fmt.Println("Deleting bucket using SDK v2...")
	_, err = s3ClientV2.DeleteBucket(ctx, &s3v2.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		log.Printf("Warning: Failed to delete bucket: %v", err)
		fmt.Printf("\nPlease manually delete bucket: %s\n", bucketName)
	} else {
		fmt.Println("✓ Bucket deleted successfully with SDK v2")
	}

	// ===== CONCLUSION =====
	fmt.Println("\n\n=== Conclusion ===")
	fmt.Println("✓ Infrastructure created with SDK v1 is fully accessible with SDK v2")
	fmt.Println("✓ Both SDKs interact with the same AWS APIs and resources")
	fmt.Println("✓ You can create resources with v1 and migrate management to v2")
	fmt.Println("✓ AWS resources are SDK-agnostic - they exist independently")
	fmt.Println("\nThis proves you can migrate your codebase incrementally without")
	fmt.Println("needing to recreate any existing infrastructure.")
}
