package main

type NoImage struct{}

func (e NoImage) Error() string {
	return "no image represented"
}
