package helpers

import "hash/fnv"

func GetInt16Hash(s string) int16 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int16(h.Sum32())
}
