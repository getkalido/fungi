package fungi_test

import (
	"fmt"

	"github.com/getkalido/fungi"
)

func ExampleFlatten() {
	someStream := fungi.SliceStream([][]int64{{1, 2, 3}, {4, 5}})
	err := fungi.Loop(func(in int64) {
		fmt.Println(in)
	})(fungi.Flatten[int64](someStream))
	if err != nil {
		panic(err)
	}

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
}
