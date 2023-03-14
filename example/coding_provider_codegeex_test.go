package script

import (
	"fmt"
	"testing"
)

func TestComplete(t *testing.T) {
	content := `
// Write a function that returns the sum of the numbers from 1 to n,
// but using a for loop instead of a while loop.

func sumN(n int) int {
	var sum int	
`
	retList := Complete("go", content)
	fmt.Println(content)
	fmt.Println(retList)
}

func TestConvert(t *testing.T) {
	content := `
def pairs_sum_to_zero(l):
    for i, l1 in enumerate(l):
        for j in range(i + 1, len(l)):
            if l1 + l[j] == 0:
                return True
    return False
`
	retList := Convert("python", "go", content)
	fmt.Println(retList)
}

func TestExplain(t *testing.T) {
	content := `
def sum(a,b):
    return a+b	
`
	retList := Explain("python", content)
	fmt.Println(retList)
}
