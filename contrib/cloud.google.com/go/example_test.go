package cloud_test

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	cloudtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/cloud.google.com/go"
)

func Example() {
	ctx := context.Background()

	// create an http client which will trace requests
	client, _ := cloudtrace.NewClient(ctx)

	// create a storage client using the given http client
	svc, _ := storage.NewClient(ctx, option.WithHTTPClient(client))

	// call storage methods as usual
	it := svc.Buckets(ctx, "some-project-id")
	for {
		bucket, err := it.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			panic(err)
		}
		fmt.Println(bucket.Name)
	}
}
