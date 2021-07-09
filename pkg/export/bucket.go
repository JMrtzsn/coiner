package export

import (
	gcp "cloud.google.com/go/storage"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"
)

const timeout = 50

type Bucket struct {
	bkt      gcp.BucketHandle
	ctx      context.Context
	exchange string
}

// NewBucket constructor
func NewBucket(ctx context.Context, exchange, path string) (*Bucket, error) {
	client, err := gcp.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	bkt := *client.Bucket(path)

	return &Bucket{
		bkt:      bkt,
		ctx:      ctx,
		exchange: exchange,
	}, nil
}

func (b Bucket) String() string {
	return "Bucket"
}

// Export copies a file to a GCP bucket
func (b Bucket) Export(csv *os.File, date, symbol string) error {
	ctx, cancel := context.WithTimeout(b.ctx, time.Second*timeout)
	defer cancel()

	path := b.Path(date, symbol)
	w := b.bkt.Object(path).NewWriter(ctx)
	_, err := csv.Seek(0, 0)
	if err != nil {
		return err
	}
	if _, err := io.Copy(w, csv); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	return nil
}

// Read reads a bucket object into records
func (b Bucket) Read(date, symbol string) ([][]string, error) {
	ctx, cancel := context.WithTimeout(b.ctx, time.Second*timeout)
	defer cancel()

	path := b.Path(date, symbol)
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
	return records, nil
}

// Delete the folder/file
func (b Bucket) Delete(filepath string) error {
	ctx, cancel := context.WithTimeout(b.ctx, time.Second*timeout)
	defer cancel()

	obj := b.bkt.Object(filepath)
	err := obj.Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Path generates a exchange/symbol/date.csv path
func (b Bucket) Path(date, symbol string) string {
	return fmt.Sprintf("%s/%s/%s.csv", b.exchange, symbol, date)
}
