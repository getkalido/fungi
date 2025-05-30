package fungi_test

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/getkalido/fungi"
)

func ExampleChunk() {
	someStream := fungi.SliceStream([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	err := fungi.Loop(func(in []int64) {
		fmt.Println(in)
	})(fungi.Chunk[int64](3)(someStream))
	if err != nil {
		panic(err)
	}

	// Output:
	// [1 2 3]
	// [4 5 6]
	// [7 8 9]
	// [10]
}

func ExampleChunk_ChunkAndReflatten() {
	someStream := fungi.SliceStream([]int64{1, 2, 3, 4, 5})

	err := fungi.Loop(func(in int32) {
		fmt.Println(in)
	})(
		fungi.Flatten[int32](
			fungi.Map[[]int64, []int32](
				func(in []int64) []int32 {
					// This transformer is just for demonstration purposes,
					// but can be replaced with, for example, API calls
					fmt.Println(in)
					result := make([]int32, 0, len(in))
					for _, v := range in {
						result = append(result, int32(v))
					}
					return result
				})(
				fungi.Chunk[int64](3)(
					someStream,
				),
			),
		),
	)
	if err != nil {
		panic(err)
	}

	// Output:
	// [1 2 3]
	// 1
	// 2
	// 3
	// [4 5]
	// 4
	// 5
}

func TestChunkZero(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for chunk size 0, but did not panic")
		}
	}()
	fungi.Chunk[int64](0)
}

type SequentialStream struct {
	MaxCounter int64
	FailFunc   func(int64) (int64, error)
	counter    int64
}

func (s *SequentialStream) Next() (int64, error) {
	if s.MaxCounter <= 0 {
		s.MaxCounter = 10
	}
	if s.counter < s.MaxCounter {
		s.counter++
		return s.counter, nil
	}
	if s.FailFunc == nil {
		return 0, io.EOF
	}
	return s.FailFunc(s.counter)
}

func SliceEqual(a []int64, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestChunkLengthSubdividesExactly(t *testing.T) {
	stream := &SequentialStream{
		MaxCounter: 9,
	}
	chunker := fungi.Chunk[int64](3)(stream)

	var (
		chunk []int64
		err   error
	)
	for i, expectedChunk := range [][]int64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	} {
		chunk, err := chunker.Next()
		if err != nil {
			t.Fatalf("Expected no error on chunk %d, got: %v", i+1, err)
		}
		if !SliceEqual(chunk, expectedChunk) {
			t.Fatalf("Expected chunk %d to be %v, got: %v", i+1, expectedChunk, chunk)
		}
	}

	chunk, err = chunker.Next()
	if len(chunk) > 0 {
		t.Fatalf("Expected exhausted chunk to be empty, got: %v", chunk)
	}
	if err != io.EOF {
		t.Fatalf("Expected io.EOF on fourth chunk, but got %v", err)
	}
}

func TestChunkError(t *testing.T) {
	stream := &SequentialStream{
		MaxCounter: 10,
		FailFunc: func(counter int64) (int64, error) {
			return 0, fmt.Errorf("error at counter %d", counter)
		},
	}
	chunker := fungi.Chunk[int64](3)(stream)

	var (
		chunk []int64
		err   error
	)
	for i, expectedChunk := range [][]int64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	} {
		chunk, err := chunker.Next()
		if err != nil {
			t.Fatalf("Expected no error on chunk %d, got: %v", i+1, err)
		}
		if !SliceEqual(chunk, expectedChunk) {
			t.Fatalf("Expected chunk %d to be %v, got: %v", i+1, expectedChunk, chunk)
		}
	}

	chunk, err = chunker.Next()
	if len(chunk) > 0 {
		t.Fatalf("Expected errored chunk to be empty, got: %v", chunk)
	}
	if err == nil {
		t.Fatal("Expected error on fourth chunk, but got nil")
	}

	chunk, finalErr := chunker.Next()
	if len(chunk) > 0 {
		t.Fatalf("Expected errored chunk to be empty, got: %v", chunk)
	}
	if err == nil {
		t.Fatal("Expected error on chunk after error, but got nil")
	}
	errors.Is(finalErr, err)
	if err == nil {
		t.Fatalf("To receive error again, but got %v", finalErr)
	}
}

func TestChunkPullAfterExhaustion(t *testing.T) {
	stream := &SequentialStream{
		MaxCounter: 10,
	}
	chunker := fungi.Chunk[int64](3)(stream)

	var (
		chunk []int64
		err   error
	)
	for i, expectedChunk := range [][]int64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
		{10},
	} {
		chunk, err := chunker.Next()
		if err != nil {
			t.Fatalf("Expected no error on chunk %d, got: %v", i+1, err)
		}
		if !SliceEqual(chunk, expectedChunk) {
			t.Fatalf("Expected chunk %d to be %v, got: %v", i+1, expectedChunk, chunk)
		}
	}

	chunk, err = chunker.Next()
	if len(chunk) > 0 {
		t.Fatalf("Expected exhausted chunk to be empty, got: %v", chunk)
	}
	if err != io.EOF {
		t.Fatalf("Expected io.EOF on fourth chunk, but got %v", err)
	}

	chunk, err = chunker.Next()
	if len(chunk) > 0 {
		t.Fatalf("Expected exhausted chunk to be empty again, got: %v", chunk)
	}
	if err != io.EOF {
		t.Fatalf("Expected io.EOF on fifth chunk, but got %v", err)
	}
}
