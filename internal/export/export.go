package export

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/jmrtzsn/coiner/internal"
	"google.golang.org/api/iterator"
	"io"
	"log"
	"os"
)

func Export(symbol, bucket, name string, ohlcv []internal.OHLCV) {
	ctx := context.Background()

	// The client will use the default application credentials.
	// Clients should be reused instead of created as needed.
	client, err := storage.NewClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	bkt := client.Bucket(bucket)

	obj := bkt.Object(name)

	// w implements io.Writer.
	w := obj.NewWriter(ctx)
	// w.Write()

	// Write some text to obj.
	//This will either create the object or overwrite whatever is there already.
	if _, err := fmt.Fprintf(w, "This object contains text.\n"); err != nil {
		// TODO: Handle error.
	}
	// Close, just like writing a file.
	if err := w.Close(); err != nil {
		// TODO: Handle error.
	}
}

func readObject(bkt *storage.BucketHandle, ctx context.Context, name string)  {
	obj := bkt.Object(name)
	// Read it back.
	r, err := obj.NewReader(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	defer r.Close()
	if _, err := io.Copy(os.Stdout, r); err != nil {
		// TODO: Handle error.
	}
}

func listBuckets(bkt *storage.BucketHandle, ctx context.Context) []string {
	// Optional Query filter, defaults a,b,c..
	query := &storage.Query{Prefix: ""}

	var names []string
	it := bkt.Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		names = append(names, attrs.Name)
	}
	return names
}
