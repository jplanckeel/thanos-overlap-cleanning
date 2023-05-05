package pkg

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

type client struct {
	minioClient *minio.Client
	endpoint    string
	bucketName  string
}

func newClient(
	provider string,
	bucketName string,
	accessKey,
	secretKey string,
	region string,
	maxTime string,
	minTime string,
	labels string,
) (*client, error) {

	if bucketName == "" {
		logrus.Fatal("no bucket bucket name")
	}

	endpoint := getS3Endpoint(provider, region)
	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Region: region,
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	return &client{
		endpoint:    endpoint,
		bucketName:  bucketName,
		minioClient: minioClient,
	}, nil
}

func (c *client) listMeta() []string {

	metadatas := make([]string, 0)

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	objectCh := c.minioClient.ListObjects(ctx, c.bucketName, minio.ListObjectsOptions{
		Recursive: true,
	})
	for object := range objectCh {
		if object.Err != nil {
			logrus.Error(object.Err)
			return nil
		}

		if strings.Contains(object.Key, "meta.json") {
			// call getobject directly to not create object ?
			metadatas = append(metadatas, object.Key)
		}

	}
	return metadatas
}

func (c *client) getObjectFileContent(objectName string) (string, error) {
	logrus.Debugf("Read object: %s/%s", c.bucketName, objectName)
	obj, err := c.minioClient.GetObject(
		context.Background(),
		c.bucketName,
		objectName,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return "", err
	}

	buffer := bytes.NewBuffer(make([]byte, 0))
	part := make([]byte, 1024)

	var count int
	for {
		count, err = obj.Read(part)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return "", err
			}
			buffer.Write(part[:count])
			break
		}
		buffer.Write(part[:count])
	}

	ret := buffer.String()
	logrus.Debug(ret)

	return ret, nil
}

func (c *client) removeObjects(objectName string) {
	logrus.Info("Remove object: ", c.bucketName, "/", objectName)

	objectsCh := make(chan minio.ObjectInfo)
	ctx := context.Background()

	// Send object names that are needed to be removed to objectsCh
	go func() {
		defer close(objectsCh)
		// List all objects from a bucket-name with a matching prefix.
		for object := range c.minioClient.ListObjects(ctx, c.bucketName, minio.ListObjectsOptions{
			Recursive: true,
			Prefix:    objectName,
		}) {
			if object.Err != nil {
				logrus.Fatalln(object.Err)
			}
			objectsCh <- object
		}
	}()

	opts := minio.RemoveObjectsOptions{
		GovernanceBypass: true,
	}

	for rErr := range c.minioClient.RemoveObjects(ctx, c.bucketName, objectsCh, opts) {
		logrus.Warn("Error detected during deletion: ", rErr)
	}

}

func getS3Endpoint(provider, region string) string {
	if provider == "scw" {
		return fmt.Sprintf("s3.%s.scw.cloud", region)
	}

	return fmt.Sprintf("s3.%s.amazonaws.com", region)
}
