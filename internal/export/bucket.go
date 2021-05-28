package export

import (
	gcp "cloud.google.com/go/storage"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"google.golang.org/api/iterator"
	"io"
	"log"
	"os"
	"time"
)

const timeout = 50

// Storage type struct
type Bucket struct {
	bkt      gcp.BucketHandle
	ctx      context.Context
	exchange string
	symbol   string
}

// New is the constructor for Storage
func New(ctx context.Context, exchange, symbol string) (*Bucket, error) {
	client, err := gcp.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	// TODO assert bucket exist if not create
	bucket, ok := os.LookupEnv("CLOUD_BUCKET")
	if !ok {
		return nil, errors.New("invalid env var")
	}
	bkt := *client.Bucket(bucket)

	return &Bucket{
		bkt:      bkt,
		ctx:      ctx,
		exchange: exchange,
		symbol:   symbol,
	}, nil
}

func (b Bucket) Export(file *os.File, date string) error {
	ctx, cancel := context.WithTimeout(b.ctx, time.Second*timeout)
	defer cancel()

	path := b.Path(date)
	w := b.bkt.Object(path).NewWriter(ctx)
	file.Seek(0, 0)
	if _, err := io.Copy(w, file); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	return nil
}

// TODO: pass fileformat?
func (b Bucket) Read(date string) ([][]string, error) {
	ctx, cancel := context.WithTimeout(b.ctx, time.Second*timeout)
	defer cancel()

	path := b.Path(date)
	obj := b.bkt.Object(path)
	r, err := obj.NewReader(ctx)
	if err != nil {
		return nil, err
	}
	records, err := csv.NewReader(r).ReadAll()
	defer r.Close()
	if _, err := io.Copy(os.Stdout, r); err != nil {
		return nil, err
	}
	// TODO: Logging
	// fmt.Fprintf(w, "Read b %v uploaded.\n", path)
	return records, nil
}

// Delete the folder/file
func (b Bucket) Delete(filepath string) error {
	obj := b.bkt.Object(filepath)
	err := obj.Delete(b.ctx)
	if err != nil {
		return err
	}
	return nil
}

// List all files
func (b Bucket) List() []string {
	query := &gcp.Query{Prefix: ""}

	var names []string
	it := b.bkt.Objects(b.ctx, query)
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

// Path generates a exchange/symbol/date.csv path for gcp buckets
func (b Bucket) Path(date string) string {
	return fmt.Sprintf("%s/%s/%s.csv", b.exchange, b.symbol, date)
}
