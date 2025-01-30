package core

import (
	"gitub.com/sriramr98/go_kvdb/store"
)

func DeleteExpiredKeys(store *store.InMemoryStore) {
	deleted_perc := store.DeleteExpiredKeysFromXRandom(20)
	if deleted_perc > 0.25 {
		DeleteExpiredKeys(store)
	}
}
