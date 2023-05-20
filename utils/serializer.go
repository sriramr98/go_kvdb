package utils

type Serializer[T any] interface {
	Serialize(T) ([]byte, error)
	Deserialize([]byte) (T, error)
}
