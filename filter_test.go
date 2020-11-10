package rxgo

import (
	"fmt"
	"os"

	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)


func TestDebounce(t *testing.T){
	Just(0,3,4,5,6,7,8, 88).Map(func(x int)int{
		if x != 0{
			time.Sleep(1*time.Millisecond)
		}
		return x
	}).Debounce(2*time.Millisecond).Subscribe(func(x int){
		if x != 88 {
			fmt.Printf("error Debounce with %d\n", x)
			os.Exit(-1)
		}
		fmt.Printf("Debunce %d\n", x)
	})
}

func TestDistinct(t *testing.T){
	res := []int{}
	Just(1, 8, 9, 10, 1, 8, 8, 8, 8, 8).Map(func(x int) int {
		return x
	}).Distinct().Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{1, 8, 9, 10}, res, "TestDistinct Error!")
}

func TestElementAt(t *testing.T){
	Just("hello", "h","aa").ElementAt(0).Subscribe(func(x string){
		if x!= "hello" {
			fmt.Printf("TestElementAt Error with value = %s", x)
			os.Exit(1)
		}
	})
}

func TestFirst(t *testing.T) {
	Just(9,"b","c").First().Subscribe(func(x int) {
		if x != 9 {
			fmt.Printf("TestFirst Error with value: %d\n", x)
			os.Exit(1)
		}
	})
}

func TestIgnoreElements(t *testing.T){
	Just("aaa", "b", "100", "0", "mary", "peter").IgnoreElements().Subscribe(func(x string) {
		fmt.Printf("TestIgnoreElements Error with value: %s\n", x)
		os.Exit(1)
	})
}

func TestLast(t *testing.T){
	Just("aaa", "b", "100", "0", "mary", "peter").Last().Subscribe(func(x string) {
		if x != "peter" {
			fmt.Errorf("TestLast Error with value: %s\n", x)
			os.Exit(1)
		}
	})
}

func TestSample(t *testing.T){
	result := []int{2, 5, 8, 9}
	actual := []int{}
	Just(0,1,2,3,4, 5, 6, 7, 8, 9).Map(func(x int)int{
		time.Sleep(18*time.Millisecond)
		return x
	}).Sample(3*time.Millisecond).Subscribe(func(x int){
		actual = append(actual, x)
	})
	assert.Equal(t, result, actual, "--same---")
}

func TestSkip(t *testing.T) {
	res := []int{}
	Just(10, 20, 30, 40, 50).Map(func(x int) int {
		return x
	}).Skip(4).Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{50}, res, "Skip Test Error!")
}