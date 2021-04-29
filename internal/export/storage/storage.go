package storage

import (
	gcp "cloud.google.com/go/storage"
	"context"
	"encoding/csv"
	"fmt"
	"google.golang.org/api/iterator"
	"io"
	"log"
	"os"
	"time"
)

type Storage struct {
	bkt gcp.BucketHandle
	ctx context.Context
}

func Init(bucket string) (*Storage, error) {
	// TODO ctx.Timeout?
	ctx := context.Background()
	client, err := gcp.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	// TODO assert bucket exist if not create
	bkt := *client.Bucket(bucket)

	return &Storage{
		bkt: bkt,
		ctx: ctx,
	}, nil
}

func (s Storage) Write(file *os.File, path string) error {
	ctx, cancel := context.WithTimeout(s.ctx, time.Second*50)
	defer cancel()

	w := s.bkt.Object(path).NewWriter(ctx)
	if _, err := io.Copy(w, file); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	return nil
}

// TODO: pass fileformat?
// Read the object as csv
func (s Storage) Read(filepath string) ([][]string, error) {
	ctx, cancel := context.WithTimeout(s.ctx, time.Second*50)
	defer cancel()

	obj := s.bkt.Object(filepath)
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

// TODO:
// Should create folder symbol, and filename name: BTCUSDT/2019-01-01.csv
func Path(symbol, name string) string {
	return fmt.Sprintf("%s/%s.csv", symbol, name)
}
