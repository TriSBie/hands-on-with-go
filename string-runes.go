package main

import (
	"errors"
	"fmt"
	"math"
	"unicode/utf8"
)

type person struct {
	name string
	age  int
}

/* Using Method */
// type rect1 struct {
// 	width, length int
// }

type rect struct {
	width, height float64
}

type circle struct {
	radius float64
}

// func (r *rect) area() int {
// 	return r.width * r.length
// }

// func (r *rect) perimeter() int {
// 	return 2*r.width + 2*r.length
// }

// create a new person and return a pointer to it
func newPerson(name string) *person {
	p := person{name: name}
	p.age = 42
	return &p
}

type Vertex struct {
	X, Y float64
}

func (v *Vertex) abs() float64 {
	return math.Abs(v.X) + math.Abs(v.Y)
}

type geometry interface {
	area() float64
	petrim() float64
}

func (r rect) area() float64 {
	return float64(r.height) * float64(r.width)
}

func (r rect) petrim() float64 {
	return 2*r.width + 2*r.height
}

func (c circle) area() float64 {
	return math.Pi * c.radius * c.radius
}

func (c circle) petrim() float64 {
	return 2 * math.Pi * c.radius
}

func measure(g geometry) {
	fmt.Println(g)
	fmt.Println(g.area())
	fmt.Println(g.petrim())
}

type Dog struct{}

type Cat struct{}

type Animal interface {
	Speak() string
}

func (d Dog) Speak() string {
	return "Woof"
}

func (c *Cat) Speak() string {
	return "Meoww"
}

func printAll(vals []interface{}) {
	for _, val := range vals {
		fmt.Println(val)
	}
}

/* Using Enum Type*/
type ServerState int

const (
	StateIdle = iota
	StateConnected
	StateError
	StateRetrying
)

var stateName = map[ServerState]string{
	StateIdle:      "Idle",
	StateConnected: "Connected",
	StateError:     "Error",
	StateRetrying:  "Retrying",
}

func (ss ServerState) String() string {
	return stateName[ss]
}

type base struct {
	num int
}

func (b base) describe() string {
	return fmt.Sprintf("base with num= %v", b.num)
}

type container struct {
	base
	str string
}

// Using Errors
func f(arg int) (int, error) {
	if arg == 32 {
		return -1, errors.New("Can't work with 32")
	}
	return arg + 3, nil // nil indicated that there was no error
}

// using with sentiel error
var ErrOutOfTea = fmt.Errorf("No more tea available")
var ErrPower = fmt.Errorf("Can't boil water")

func makeTea(arg int) error {
	if arg == 2 {
		return ErrOutOfTea
	} else if arg == 4 {
		// wrap errors with higher-level errors to add context
		return fmt.Errorf("Making tea: %w", ErrPower)
	}
	return nil
}

// Using custom errors

type argErr struct {
	arg     int
	message string
}

func (arg *argErr) Error() string {
	return fmt.Sprintf("%d - %s", arg.arg, arg.message)
}

// declare a function which return int and error
func f1(arg int) (int, error) {
	if arg == 42 {
		return -1, &argErr{arg, "can't work with it"} // return a reference argErr
	}
	return arg + 3, nil
}

func main() {
	const nihongo = "日本語"
	fmt.Println(nihongo)
	fmt.Printf("%x\n", nihongo)  // space between each byte
	fmt.Printf("% x\n", nihongo) // space between each byte
	fmt.Printf("%q\n", nihongo)
	fmt.Printf("%+q\n")
	fmt.Println(len(nihongo)) // since string holds the arbitrary byte
	for i, w := 0, 0; i < len(nihongo); i += w {
		runeValue, width := utf8.DecodeRuneInString(nihongo[i:])
		fmt.Printf("%#U starts at %d\n", runeValue, i)
		w = width
	}

	fmt.Println(person{name: "Bob"})
	fmt.Println(&person{name: "Bob", age: 40})
	fmt.Println(newPerson("Bob"))

	s := person{name: "Anna", age: 20}
	fmt.Println(s)

	sp := &s // asign pointer to variable sp
	sp.age = 29
	fmt.Println(s, sp)

	v := Vertex{3, 4}
	fmt.Println(v.abs())

	sc := make([]int, 5, 10)

	// Print the length and capacity
	fmt.Printf("Slice: %v\n", sc)
	fmt.Printf("Length: %d\n", len(sc))
	fmt.Printf("Capacity: %d\n", cap(sc))

	// Append elements to the slice
	sc = append(sc, 1, 2, 3)

	// Print the length and capacity after appending
	fmt.Printf("\nAfter appending elements:\n")
	fmt.Printf("Slice: %v\n", sc)
	fmt.Printf("Length: %d\n", len(sc))
	fmt.Printf("Capacity: %d\n", cap(sc))

	// Further append to exceed the initial capacity
	sc = append(sc, 4, 5, 6, 7, 8, 9)

	// Print the length and capacity after exceeding initial capacity
	fmt.Printf("\nAfter exceeding initial capacity:\n")
	fmt.Printf("Slice: %v\n", sc)
	fmt.Printf("Length: %d\n", len(sc))
	fmt.Printf("Capacity: %d\n", cap(sc))

	// Interface implement
	rectVar := rect{3, 4}
	circleVar := circle{5}

	measure(rectVar)
	measure(circleVar)

	// using slice interface

	geometryVar := []geometry{
		rect{3, 4},
		circle{5},
	}

	for _, g := range geometryVar {
		fmt.Println(g.area())
		fmt.Println(g.petrim())

	}

	names := []string{"john", "standley", "mickey"}

	// printAll(names) - caused error for string[] cannot cast as interface[]

	// make(typeT[], size, capacity)
	vals := make([]interface{}, len(names))
	for i, val := range names {
		vals[i] = val
	}
	printAll(vals)

	animals := []Animal{&Dog{}, new(Cat)}
	for _, animal := range animals {
		fmt.Println(animal.Speak())
	}

	// var input = `
	// 	{"created_at" : "Thu May 31 00:00:01 +0000 2012"}
	// `

	// var val map[string]time.Time

	// if err := json.Unmarshal([]byte(input), &val); err != nil {
	// 	panic(err)
	// }

	// fmt.Println(val)
	// for k, v := range val {
	// 	fmt.Println(k, reflect.TypeOf(v))
	// }

	ns := transition(StateIdle)
	fmt.Println(ns)

	// Using embedding struct
	co := container{
		base: base{
			num: 1,
		},
		str: "some name",
	}

	/**
	With embedding struct
	When nested item has methods -> method will be called by nesting fields
	for example : co fields holds base
	-> base has methods
	-> methods of base became methods of co
	*/

	fmt.Println(co.base.num)
	fmt.Println(co.str)
	fmt.Println(co.describe())

	type describer interface {
		describe() string
	}

	var d describer = co
	fmt.Println(d.describe())

	// Using with Generics Types
	var maps = map[int]string{1: "2", 2: "3", 3: "4"}

	fmt.Println(MapKeys(maps))

	_abc := MapKeys[int, string](maps)
	fmt.Println(_abc)

	lst := List[int]{}

	lst.Push(1)
	lst.Push(2)
	lst.Push(3)

	fmt.Println(lst.GetAll())

	// ERRORS in GO
	for _, i := range []int{7, 32} {
		if r, e := f(i); e != nil {
			fmt.Println("f failed:", e)
		} else {
			fmt.Println("f worked:", r)
		}
	}

	for i := range []int{1, 2, 3} {
		if err := makeTea(i); err != nil {
			if errors.Is(err, ErrOutOfTea) {
				fmt.Println("We should buy new tea")
			} else if errors.Is(err, ErrPower) {
				fmt.Println("Not it is dark")
			} else {
				fmt.Println("Unknown error : %s\n", err)
			}
			continue
		}
		fmt.Println("Tea is ready")
	}

	/*
		Difference between errors.Is vs errors.As
		- errors.Is for sentinel errors (simple - predefined)
		- errors.As for a customize errors (with struct)
			+ errors.As required target should be a pointer by a errors type
	*/
	_, err := f1(42)
	var ae *argErr
	if errors.As(err, &ae) {
		fmt.Println(ae.arg)
		fmt.Println(ae.message)
	} else {
		fmt.Println("Error doesn't match argError")
	}
}

func transition(s ServerState) ServerState {
	switch s {
	case StateIdle:
		return StateConnected
	case StateConnected, StateRetrying:
		return StateIdle
	case StateError:
		return StateError
	default:
		panic(fmt.Errorf("unknown state %s", s))
	}
}

// map[K]V -> K : Key types ? > V : value types
// Generic types
func MapKeys[K comparable, V any](m map[K]V) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k) // return key of map
	}

	return r
}

type List[T any] struct {
	head, tail *element[T] // head and tail of list stored as reference
}

type element[T any] struct {
	next *element[T]
	val  T // any
}

func (lst *List[T]) Push(v T) {
	if lst.tail == nil {
		lst.head = &element[T]{
			val: v,
		}
		lst.tail = lst.head
	} else {
		lst.tail.next = &element[T]{
			val: v,
		}
		lst.tail = lst.tail.next
	}
}

func (lst *List[T]) GetAll() []T {
	var elems []T
	for e := lst.head; e != nil; e = e.next {
		elems = append(elems, e.val)
	}
	return elems
}
