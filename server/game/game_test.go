package game

import "testing"

func TestCollidesWith(t *testing.T) {
	// Defining the columns of the table
	var tests = []struct {
		name       string
		inputHB    Hitbox
		inputOther Hitbox
		want       bool
	}{
		// the table itself
		{"Same HB should collide",
			Hitbox{
				point:  &Point{x: 1, y: 2},
				height: 5,
				width:  5,
			},
			Hitbox{
				point:  &Point{x: 1, y: 2},
				height: 5,
				width:  5,
			},
			true},
		{"Non-overlapping HB should not collide",
			Hitbox{
				point:  &Point{x: 1, y: 2},
				height: 5,
				width:  5,
			},
			Hitbox{
				point:  &Point{x: 100, y: 200},
				height: 5,
				width:  5,
			},
			false},
		{"HB1 top right vertex & HB2 bottom left vertex should collide",
			Hitbox{
				point:  &Point{x: 20, y: 20},
				height: 5,
				width:  5,
			},
			Hitbox{
				point:  &Point{x: 30, y: 10},
				height: 10,
				width:  10,
			},
			true},
		{"HB1 bottom left vertex & HB2 top right vertex should collide",
			Hitbox{
				point:  &Point{x: 30, y: 10},
				height: 5,
				width:  5,
			},
			Hitbox{
				point:  &Point{x: 20, y: 20},
				height: 10,
				width:  10,
			},
			true},
		{"HB1 bottom right vertex & HB2 top left vertex should collide",
			Hitbox{
				point:  &Point{x: 20, y: 20},
				height: 5,
				width:  5,
			},
			Hitbox{
				point:  &Point{x: 30, y: 30},
				height: 5,
				width:  5,
			},
			true},
		{"HB1 top left vertex & HB2 bottom right vertex should collide",
			Hitbox{
				point:  &Point{x: 30, y: 30},
				height: 5,
				width:  5,
			},
			Hitbox{
				point:  &Point{x: 20, y: 20},
				height: 5,
				width:  5,
			},
			true},
	}
	// The execution loop
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ans := tc.inputHB.collidesWith(&tc.inputOther)
			if ans != tc.want {
				t.Errorf("got %t, want %t", ans, tc.want)
			}
		})
	}
}
