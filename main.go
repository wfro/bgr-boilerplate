//go:generate protoc --gogo_out=. internal/internal.proto

package main

import (
	// "encoding/json"
	"encoding/binary"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	"github.com/wfro/burger-boilerplate/internal"
	"net/http"
)

type Burger struct {
	ID       int `json:"id"`
	Price    int `json:"price"`
	Calories int `json:"calories"`
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// burgers := BurgersResponse{
	// 	Burgers: Burgers{
	// 		Burger{
	// 			Toppings: []string{"cheddar cheese", "lettuce", "mushrooms"},
	// 			Price:    100,
	// 			Calories: 1000,
	// 		},
	// 		Burger{
	// 			Toppings: []string{"bacon", "peanut butter"},
	// 			Price:    50,
	// 			Calories: 500,
	// 		},
	// 	},
	// }
	//
	// w.Header().Set("Content-Type", "application/json")
	//
	// json.NewEncoder(w).Encode(burgers)
}

// MarshalBinary encodes a user to binary format.
func (b *Burger) MarshalBinary() ([]byte, error) {
	pb := internal.Burger{
		ID:       proto.Int64(int64(b.ID)),
		Price:    proto.Int64(int64(b.Price)),
		Calories: proto.Int64(int64(b.Calories)),
	}

	return proto.Marshal(&pb)
}

func (b *Burger) UnmarshalBinary(data []byte) error {
	var pb internal.Burger
	if err := proto.Unmarshal(data, &pb); err != nil {
		return err
	}

	b.ID = int(pb.GetID())
	b.Price = int(pb.GetPrice())
	b.Calories = int(pb.GetCalories())

	return nil
}

// Store represents the data storage layer.
type Store struct {
	// Filepath to the data file.
	Path string

	db *bolt.DB
}

func (s *Store) Open() error {
	db, err := bolt.Open(s.Path, 0666, nil)
	if err != nil {
		return err
	}
	s.db = db

	// Start a writable transaction.
	tx, err := s.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Initialize buckets to guarantee that they exist.
	tx.CreateBucketIfNotExists([]byte("Burgers"))

	// Commit the transaction.
	return tx.Commit()
}

func (s *Store) CreateBurger(b *Burger) error {
	tx, err := s.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Retrieve bucket and create new ID.
	bkt := tx.Bucket([]byte("Burgers"))
	seq, _ := bkt.NextSequence()
	b.ID = int(seq)

	// Marshal our user into bytes.
	buf, err := b.MarshalBinary()
	if err != nil {
		return err
	}

	// Save user to the bucket.
	if err := bkt.Put(itob(b.ID), buf); err != nil {
		return err
	}
	return tx.Commit()
}

// itob encodes v as a big endian integer.
func itob(v int) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(v))
	return buf
}

func (s *Store) Burger(id int) (*Burger, error) {
	// Start a read-only transaction.
	tx, err := s.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Read encoded user bytes.
	v := tx.Bucket([]byte("Burgers")).Get(itob(id))
	if v == nil {
		return nil, nil
	}

	// Unmarshal bytes into a user.
	var b Burger
	if err := b.UnmarshalBinary(v); err != nil {
		return nil, err
	}

	return &b, nil
}

func main() {
	http.HandleFunc("/", indexHandler)

	fmt.Println("Listening on port :8080")
	http.ListenAndServe(":8080", nil)
}
