# ncp-object-storage-golang
본 레파지스토리는 NCP(Naver Cloud Platform)에서 [`aws-sdk-go-v2`](https://github.com/aws/aws-sdk-go-v2)를 이용해서 Object Storage를 이용하는 방법에 대해서 서술하였습니다.  
[20230208 - 네이버 클라우드 플랫폼과 AWS SDK for Go](https://hyunsang.dev/TIL/Golang/20230208-%EB%84%A4%EC%9D%B4%EB%B2%84-%ED%81%B4%EB%9D%BC%EC%9A%B0%EB%93%9C-%ED%94%8C%EB%9E%AB%ED%8F%BC%EA%B3%BC-AWS-SDK-for-Go.html)에서의 서술한 내용을 정리하여 본 레파지스토리에 서술하였습니다.

## Getting Started
```shell
$ export NCP_ACCESS_KEY = ""
$ export NCP_SECURITY_KEY = ""
```
본격적으로 `aws-sdk-go-v2`에 접근하기 위해서는 네이버 클라우드 플랫폼에서 발급 받은 인증키가 필요합니다.     
네이버 클라우드 플랫폼 포털의 마이페이지 > 계정 관리 > 인증키 관리에서 발급 받으실 수 있습니다.  
혹은 자신의 용도에 맞게 변경하여서 사용하시면 되겠습니다.  

## 코드 분석
### 네이버 클라우드 플랫폼 연결
```go
var (
    ncpAccessKey  string = os.Getenv("NCP_ACCESS_KEY")
    ncpSecretKey  string = os.Getenv("NCP_SECURITY_KEY")
    ncpKrRegion   string = "kr-standard"
    ncpKrEndPoint string = "https://kr.object.ncloudstorage.com"
)

func Init() *s3.Client {
	// Access Key and Secret Key
	creds := credentials.NewStaticCredentialsProvider(ncpAccessKey, ncpSecretKey, "")

	ncpResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           ncpKrEndPoint, // End Point
			SigningRegion: ncpKrRegion, // Region
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
```
- 인증키가 필요합니다. 인증키는 환경 변수로 불러오고 있습니다.   
인증키는 `Access Key`, `Secret Key`가 필요하며 위의 [Getting Started](#getting-started)를 참고하여 발급 받아주세요.  
- 리전별로 엔드포인트도 다릅니다. 다른 리전 및 엔드포인드를 알고 싶으시면 [Object Storage](https://api.ncloud-docs.com/docs/storage-objectstorage)를 참고해 주세요.  

### 새로운 버킷 생성하기
> `AccessDenied: Access Denied` 오류로 인해서 아직 개발하지 못 하였습니다.  
> 해법을 알고 계시면 제 메일로 알려주시거나 PR로 추가해 주시면 감사하겠습니다 :)

### 생성되어 있는 버킷의 리스트를 불러오기
```go
func GetBucketList() *s3.ListBucketsOutput {
	client := Init()

	result, err := client.ListBuckets(context.Background(), &s3.ListBucketsInput{}, func(options *s3.Options) {})
	if err != nil {
		log.Panicln(err)
	}

	return result
}
```
- 현재 생성되어 있는 버킷의 항목을 불러오는 기능입니다.

### 생성되어 있는 버킷 안의 존재하는 오브젝트 리스트 불러오기 
```go
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
```
- 버킷 내부에 업로드 되어 있는 오브젝트의 리스트를 가지고 오는 기능입니다.  
버킷의 이름을 알고 있어야지만 버킷 내부의 오브젝트의 리스트를 가지고 올 수 있습니다.

### 버킷 내부에 오브젝트 업로드하기
```go
func main() {
    uploadFile, err := os.ReadFile("./test.mp4")
    if err != nil {
        log.Panicln(err)
    }
}	
```

```go
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
```
- `os.ReadFile`를 통해서 파일을 읽고 `PutObjectInBucket()`를 이용해서 업로드 할 수 있습니다.

## LICENSE
```text
MIT License

Copyright (c) 2023 HyunSang Park

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```