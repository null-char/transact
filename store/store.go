package store

// Store consists of a map which maps to Mappable values including methods to manipulate the map data
type Store struct {
	data map[Key]Value
}

// MakeNewStore constructs a new key value store
func MakeNewStore() *Store {
	s := &Store{}
	return s
}

// Set sets a new key value pair or updates it if the key already exists
func (s *Store) Set(key Key, value Value) {
	s.data[key] = value
}

// Get fetches the Mappable value if the key mapping exists
func (s Store) Get(key Key) (Value, bool) {
	value, ok := s.data[key]
	return value, ok
}

// Delete will remove the key value mapping associated with the provided key. Returns a status indicating a no-op if false (key mapping does not exist)
func (s *Store) Delete(key Key) bool {
	// delete is a no-op if there is no such key so we need to manually check if key exists
	_, ok := s.Get(key)
	delete(s.data, key)
	return ok
}

// Count returns an uint representing the number of times the provided value mapping exists in the store
func (s Store) Count(value Value) uint {
	var count uint = 0

	for _, v := range s.data {
		if v == value {
			count++
		}
	}

	return count
}

// GetKVPairs returns a slice of key value pairs with first field indicating the key and second field indicating the value
func (s Store) GetKVPairs() []KVPair {
	var pairs []KVPair = make([]KVPair, 0, len(s.data))

	for k, v := range s.data {
		pairs = append(pairs, KVPair{First: k, Second: v})
	}

	return pairs
}
