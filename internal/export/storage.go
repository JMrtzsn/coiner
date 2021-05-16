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
type Storage struct {
	bkt gcp.BucketHandle
	ctx context.Context
	exchange string
	symbol   string
}

// newStorage is the constructor for Storage
func newStorage(ctx context.Context, exchange, symbol string) (*Storage, error) {
	client, err := gcp.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	// TODO assert bucket exist if not create
	bucket, ok := os.LookupEnv("CLOUD_BUCKET")
	if !ok{
		return nil, errors.New("invalid env var")
	}
	bkt := *client.Bucket(bucket)

	return &Storage{
		bkt: bkt,
		ctx: ctx,
		exchange: exchange,
		symbol: symbol,
	}, nil
}

func (s Storage) Export(file *os.File, date string) error {
	ctx, cancel := context.WithTimeout(s.ctx, time.Second*timeout)
	defer cancel()

	path := storagePath(s.exchange, s.symbol, date)
	w := s.bkt.Object(path).NewWriter(ctx)
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
func (s Storage) Read(date string) ([][]string, error) {
	ctx, cancel := context.WithTimeout(s.ctx, time.Second*timeout)
	defer cancel()

	path := storagePath(s.exchange, s.symbol, date)
	obj := s.bkt.Object(path)
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
func (s Storage) Delete(filepath string) error {
	obj := s.bkt.Object(filepath)
	err := obj.Delete(s.ctx)
	if err != nil {
		return err
	}
	return nil
}

// List all files
func (s Storage) List() []string {
	// Optional Query filter, defaults a,b,c..
	query := &gcp.Query{Prefix: ""}

	var names []string
	it := s.bkt.Objects(s.ctx, query)
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


// storagePath generates a folder/file.csv
func storagePath(exchange, symbol, name string) string {
	return fmt.Sprintf("%s/%s/%s.csv", exchange, symbol, name)
}
