package fungi_test

import (
	"fmt"

	"github.com/getkalido/fungi"
)

func ExampleConcat() {
	err := fungi.Loop(func(in int64) {
		fmt.Println(in)
	})(
		fungi.Concat(
			fungi.SliceStream([]int64{1, 2, 3}),
			fungi.SliceStream([]int64{4, 5}),
		),
	)
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
