package s3bucket

import (
	"context"
	"fmt"
	"os"
	"time"
)

func RunBucket() {
	AWS_ACCESS_KEY_ID := os.Getenv("AWS_ACCESS_KEY_ID")
	AWS_SECRET_ACCESS_KEY := os.Getenv("AWS_SECRET_ACCESS_KEY")
	AWS_REGION := os.Getenv("AWS_REGION")
	AWS_S3_BUCKET := os.Getenv("AWS_S3_BUCKET")

	s3Client, err := NewWithCredentials(AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_REGION)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	filekey := "myimage"
	// fileName := "/home/iamnitesh/Desktop/work_space/golang-practice/s3bucket/myimage.jpg"
	// err = s3Client.UploadFile(ctx, AWS_S3_BUCKET, "myimage", fileName)
	// if err != nil {
	// 	panic(err)
	// }

	objects, err := s3Client.ListObjects(ctx, AWS_S3_BUCKET)
	if err != nil {
		panic(err)
	}

	for _, object := range objects {
		fmt.Println("Object Key:", *object.Key)
		fmt.Println("Object Size:", *object.Size)
		fmt.Println("Object Last Modified:", *object.LastModified)
		fmt.Println("Object ETag:", *object.ETag)

	}

	// err = s3Client.DownloadFile(ctx, AWS_S3_BUCKET, filekey, "downloaded_image.jpg")
	// if err != nil {
	// 	panic(err)
	// // }

	publicURL := s3Client.GetPublicURL(AWS_S3_BUCKET, filekey, AWS_REGION)
	fmt.Println("$$$Public URL:", publicURL)
	fmt.Println("\n\n")

	presignedURL, err := s3Client.GetPresignedURL(ctx, AWS_S3_BUCKET, filekey, time.Minute)
	if err != nil {
		fmt.Println("Error generating presigned URL:", err)
	} else {
		fmt.Println("$$Presigned URL:", presignedURL)
	}

	fmt.Println("File downloaded successfully:", filekey)

}
