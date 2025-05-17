package store

var Store = DBStore{
	Data: make(map[string]Item),
}
