package store

// NOTE: We provide dummy methods just so that we can have a shared interface between these types
// You can think of it as a weird way of achieving union types in Go since we don't have generics (yet)

// Mappable denotes all the value types that can be mapped to our store (global and local)
type Mappable interface {
	isMappable()
}

// Number is a Mappable with underlying type int
type Number int

func (n Number) isMappable()

// String is a Mappable with underlying type string
type String string

func (s String) isMappable()

type Key = string
type Value = Mappable

type KVPair struct {
	First  Key
	Second Value
}
