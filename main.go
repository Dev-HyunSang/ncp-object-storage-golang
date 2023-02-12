package main

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"log"
	"os"
)

var (
	ncpAccessKey  string = os.Getenv("NCP_ACCESS_KEY")
	ncpSecretKey  string = os.Getenv("NCP_SECURITY_KEY")
	ncpKrRegion   string = "kr-standard"
	ncpKrEndPoint string = "https://kr.object.ncloudstorage.com"
)

func Init() *s3.Client {
	creds := credentials.NewStaticCredentialsProvider(ncpAccessKey, ncpSecretKey, "")

	ncpResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           ncpKrEndPoint,
			SigningRegion: ncpKrRegion,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithEndpointResolverWithOptions(ncpResolver),
		config.WithCredentialsProvider(creds))
	if err != nil {
		log.Panicln(err)
	}

	client := s3.NewFromConfig(cfg)

	return client
}

func main() {
	bucketName := "hello-world"

	result := GetBucketList()
	log.Println(result)

	result2 := GetBucketInObject(bucketName)
	log.Println(result2)

	result3 := DeleteBucket(bucketName)
	log.Println(result3)

	uploadFile, err := os.ReadFile("./test.mp4")
	if err != nil {
		log.Panicln(err)
	}

	result4 := PutObjectInBucket(uploadFile, bucketName, "./test.mp4", "")
	log.Println(result4)
}

func GetBucketList() *s3.ListBucketsOutput {
	client := Init()

	result, err := client.ListBuckets(context.Background(), &s3.ListBucketsInput{}, func(options *s3.Options) {})
	if err != nil {
		log.Panicln(err)
	}

	return result
}

func GetBucketInObject(bucketName string) *s3.ListObjectsOutput {
	client := Init()
	result, err := client.ListObjects(context.Background(),
		&s3.ListObjectsInput{
			Bucket: &bucketName,
		},
		func(options *s3.Options) {})
	if err != nil {
		log.Panicln(err)
	}

	return result
}

func DeleteBucket(bucketName string) *s3.DeleteBucketOutput {
	client := Init()
	result, err := client.DeleteBucket(context.Background(), &s3.DeleteBucketInput{
		Bucket: &bucketName,
	}, func(options *s3.Options) {})
	if err != nil {
		log.Panicln(err)
	}

	return result
}

func PutObjectInBucket(file []byte, bucketName, fileName, acl string) *manager.UploadOutput {
	client := Init()

	uploader := manager.NewUploader(client)

	result, err := uploader.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &fileName,
		Body:   bytes.NewReader(file),
		ACL:    types.ObjectCannedACL(*aws.String(acl)),
	})
	if err != nil {
		log.Panicln(err)
	}

	return result
}
