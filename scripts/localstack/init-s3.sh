#!/bin/bash

# Wait for LocalStack to be fully ready
echo "Waiting for LocalStack to be ready..."
sleep 2

# Create S3 bucket for acacia
echo "Creating S3 bucket: acacia-issues"
awslocal s3 mb s3://acacia-issues --region us-east-1

# Verify bucket was created
echo "Listing S3 buckets:"
awslocal s3 ls

echo "LocalStack S3 initialization complete!"
