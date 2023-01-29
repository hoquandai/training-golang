package main

import (
  "fmt"
)

type House struct {
  Material string
  HasFireplace bool
  Floors int
}

func NewHouse(opts ...HouseOption) *House {
  const (
    defaultFloors = 2
    defaultHasFireplace = false
    defaultMaterial = "wood"
  )

  h := &House {
    Material: defaultMaterial,
    HasFireplace: defaultHasFireplace,
    Floors: defaultFloors,
  }

  for _, opt := range opts {
    opt(h)
  }

  return h
}

type HouseOption func(*House)

func WithConcrete() HouseOption {
  return func(h *House) {
    h.Material = "concrete"
  }
}

func WithFireplace() HouseOption {
  return func(h *House) {
    h.HasFireplace = true
  }
}

func WithFloors(floors int) HouseOption {
  return func(h *House) {
    h.Floors = floors
  }
}

func main() {
  h := NewHouse(WithFloors(3), WithFireplace(), WithConcrete())

  fmt.Println("floors: ", h.Floors)
  fmt.Println("fireplace: ", h.HasFireplace)
  fmt.Println("material: ", h.Material)

  h = NewHouse()

  fmt.Println("floors: ", h.Floors)
  fmt.Println("fireplace: ", h.HasFireplace)
  fmt.Println("material: ", h.Material)
}
