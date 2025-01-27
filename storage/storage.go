package storage

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"

	_ "go.beyondstorage.io/services/azblob/v3"
	_ "go.beyondstorage.io/services/fs/v4"

	_ "go.beyondstorage.io/services/ftp" // Ensure the correct version is imported
	_ "go.beyondstorage.io/services/gcs/v3"
	"go.beyondstorage.io/v5/pkg/randbytes"
	"go.beyondstorage.io/v5/services"
	"go.beyondstorage.io/v5/types"

	_ "go.beyondstorage.io/services/gcs/v3"

	"go.beyondstorage.io/v5/pairs"
)

func Storage() {

	server, err := NewFTPFromString()
	if err != nil {
		panic(err)
	}

	// List will create an iterator of object under path.
	it, err := server.List(".", pairs.WithListMode(types.ListModeDir))
	if err != nil {
		panic(err)
	}

	for {
		o, err := it.Next()
		if err != nil {
			if err == types.IterateDone {
				break
			}
			fmt.Println("Error: ", err.Error())
		}
		fmt.Println("Object:", o.ID)

		fmt.Println("Object:", o.Path)

	}
	fmt.Println("Working dir : ", server.Metadata().WorkDir)
	server.Write("test.txt", randbytes.NewRand(), 1024)

	// Write will write data to object.
	// fmt.Println(object.Path)

	fs_server, err := NewFsFromString()
	if err != nil {
		panic(err)
	}

	// List will create an iterator of object under path.
	fs_it, err := fs_server.List(".", pairs.WithListMode(types.ListModeDir))
	if err != nil {
		panic(err)
	}

	fmt.Println("FS STORAGE STARTED")

	for {
		o, err := fs_it.Next()
		if err != nil {
			if err == types.IterateDone {
				break
			}
			fmt.Println("Error: ", err.Error())
		}
		fmt.Println("Object:", o.ID)

		fmt.Println("Object:", o.Path)
	}

	fs_server.Write("test.txt", randbytes.NewRand(), 1024)

}

func NewFTPFromString() (types.Storager, error) {
	str := fmt.Sprintf(
		"ftp://%s?credential=%s&endpoint=%s",
		"/app",
		"basic:newftpuser:newftpuser",
		"tcp:localhost:21",
	)
	return services.NewStoragerFromString(str)
}

func NewFsFromString() (types.Storager, error) {
	connStr := fmt.Sprintf(
		"fs://%s",
		"/home/iamnitesh/Desktop/work_space/golang-practice/fsdir",
	)
	return services.NewStoragerFromString(connStr)
}

func NewMinioFromString() (types.Storager, error) {
	str := fmt.Sprintf(
		"minio://%s%s?credential=%s&endpoint=%s",
		os.Getenv("STORAGE_MINIO_NAME"),
		os.Getenv("STORAGE_MINIO_WORKDIR"),
		os.Getenv("STORAGE_MINIO_CREDENTIAL"),
		os.Getenv("STORAGE_MINIO_ENDPOINT"),
	)
	return services.NewStoragerFromString(str)
}

func NewS3FromString() (types.Storager, error) {
	connStr := fmt.Sprintf(
		"s3://%s%s?credential=%s&endpoint=%s&location=%s&enbale_virtual_dir",
		os.Getenv("STORAGE_S3_NAME"),
		os.Getenv("STORAGE_S3_WORKDIR"),
		os.Getenv("STORAGE_S3_CREDENTIAL"),
		os.Getenv("STORAGE_S3_ENDPOINT"),
		os.Getenv("STORAGE_S3_LOCATION"),
	)
	return services.NewStoragerFromString(connStr)
}

// func NewGCSFromString() (types.Storager, error) {
// 	// str := "gcs://bucket_name/path/to/workdir?credential=file:<absolute_path_to_token_file>&project_id=<google_cloud_project_id>"
// 	str := fmt.Sprintf(
// 		"gcs://%s%s?credential=file:%s&project_id=%s",
// 		os.Getenv("STORAGE_GCS_BUCKET_NAME"),
// 		os.Getenv("STORAGE_GCS_WORKDIR"),
// 		os.Getenv("STORAGE_GCS_CREDENTIAL"),
// 		os.Getenv("STORAGE_GCS_PROJECT_ID"),
// 	)
// }

func WriteData(store types.Storager, path string) {
	// content to write
	size := rand.Int63n(4 * 1024 * 1024)
	r := io.LimitReader(randbytes.NewRand(), size)

	// Write needs at least three arguments.
	// `path` is the path of object.
	// `r` is io.Reader instance for reading the data for uploading.
	// `size` is the length, in bytes, of the data for uploading.
	//
	// Write will return two values.
	// `n` is the size of write operation.
	// `err` is the error during this operation.
	n, err := store.Write(path, r, size)
	if err != nil {
		log.Fatalf("write %v: %v", path, err)
	}

	log.Printf("write size: %d", n)
}
