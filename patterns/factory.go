package main

import (
	"fmt"
)

type Animal struct {
	species string
	age int
}

type AnimalHouse struct {
	name string
	sizeInMeters int
}

type AnimalFactory struct {
	species string
	houseName string
}

func (af *AnimalFactory) NewAnimal(age int) {
	return Animal{
		species: af.species,
		age: age,
	}
}

func (af *AnimalFactory) NewHouse(sizeInMeters int) {
	return AnimalHouse{
		name: af.houseName,
		sizeInMeters: sizeInMeters,
	}
}

func main() {

}