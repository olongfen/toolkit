package db_data

import (
	"fmt"
	"testing"
)

func TestProcessDBWhere(t *testing.T) {
	p := ProcessDBWhere("id", []int{23423, 32423}, "in")
	fmt.Println(p)
}
