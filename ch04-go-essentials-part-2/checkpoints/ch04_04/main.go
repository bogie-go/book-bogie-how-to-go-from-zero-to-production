package main

import (
	"fmt"
	"time"
)

type Timestamps struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t Timestamps) Age() time.Duration {
	return time.Since(t.CreatedAt)
}

type Order struct {
	Timestamps
	Number int
}

func main() {
	order := Order{Number: 1001}
	order.CreatedAt = time.Now()

	fmt.Println(order.Number, order.Age() < time.Second)
}
