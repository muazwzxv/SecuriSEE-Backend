package service

import (
	"context"
	"fmt"
	"io"

	"github.com/oracle/oci-go-sdk/example/helpers"
	"github.com/oracle/oci-go-sdk/objectstorage"
	"github.com/oracle/oci-go-sdk/v47/common"
)

type ObjectStorageInstance struct {
	Client    objectstorage.ObjectStorageClient
	Context   context.Context
	Bucket    string
	Namespace string
}

func ConnectToObjectStorage() *ObjectStorageInstance {

	c, err := objectstorage.NewObjectStorageClientWithConfigurationProvider(common.DefaultConfigProvider())
	if err != nil {
		panic(err.Error())
	}

	ctx := context.Background()
	namespace := getNamespace(ctx, c)

	return &ObjectStorageInstance{
		Client:    c,
		Context:   ctx,
		Bucket:    "bucket-al-kolo-kontol",
		Namespace: namespace,
	}

}

// return namespace
func getNamespace(ctx context.Context, c objectstorage.ObjectStorageClient) string {
	request := objectstorage.GetNamespaceRequest{}
	r, err := c.GetNamespace(ctx, request)
	helpers.FatalIfError(err)
	fmt.Println("get namespace")
	return *r.Value
}

func (objectStorage *ObjectStorageInstance) UploadFile(fileName string, contentLen int64, content io.ReadCloser, metadata map[string]string) error {

	request := objectstorage.PutObjectRequest{
		NamespaceName: &objectStorage.Namespace,
		BucketName:    &objectStorage.Bucket,
		ObjectName:    &fileName,
		ContentLength: &contentLen,
		PutObjectBody: content,
		OpcMeta:       metadata,
	}
	_, err := objectStorage.Client.PutObject(objectStorage.Context, request)
	return err
}

func (objectStorage *ObjectStorageInstance) DownloadFile(fileName string) (objectstorage.GetObjectResponse, error) {

	req := objectstorage.GetObjectRequest{
		NamespaceName: &objectStorage.Namespace,
		BucketName:    &objectStorage.Bucket,
		ObjectName:    &fileName,
	}
	res, err := objectStorage.Client.GetObject(objectStorage.Context, req)
	if err != nil {
		return res, err
	}

	return res, nil
}
