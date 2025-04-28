package fungi

// Unique returns a stream of unique items. The incoming stream must be a stream
// of comparable items.
func Unique[T comparable](items Stream[T]) Stream[T] {
	return UniqueBy(identity[T])(items)
}

// UniqueBy accepts a stream of any items and returns a stream of unique items.
// To be able to do that efficiently, it requires a function that can convert
// a stream item into something comparable. It's like a hashing function, but it
// doesn't have to be.
//
// For example, let's say there's a struct called User:
//
//	type User struct {
//	    ID int // unique database identifier
//	    Name string
//	}
//
// You can create a stream of unique users like so:
//
//	var usersStream fungi.Stream[*User] = service.GetUsers()
//	//  usersStream might contain duplicates
//	uniqueUsers := fungi.UniqueBy(func(u *User) int { return u.ID })
//	uniqueUsersStream := uniqueUsers(usersStream)
func UniqueBy[T any, K comparable](id func(T) K) StreamIdentity[T] {
	if id == nil {
		panic("nil id function in unique by")
	}
	return func(s Stream[T]) Stream[T] {
		return &uniqueBy[T, K]{
			id:     id,
			source: s,
		}
	}
}

type uniqueBy[T any, K comparable] struct {
	id      func(T) K
	source  Stream[T]
	seenIDs map[K]struct{}
}

func (u *uniqueBy[T, K]) Next() (item T, err error) {
	if u.seenIDs == nil {
		u.seenIDs = make(map[K]struct{})
	}
	for {
		item, err = u.source.Next()
		if err != nil {
			return
		}
		id := u.id(item)
		if _, seen := u.seenIDs[id]; !seen {

			u.seenIDs[id] = struct{}{}
			return
		}
	}
}

func identity[T any](item T) T { return item }
