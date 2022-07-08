package storage

import (
	"context"

	"github.com/dgraph-io/badger"
	"github.com/dgraph-io/badger/pb"
	"github.com/google/uuid"
)

func (s *Storage) WalkAndPruneQueue(f func(uuid.UUID) error) error {
	stream := s.s.NewStream()
	stream.Prefix = queuePrefix
	stream.NumGo = 2

	stream.Send = func(k *pb.KVList) error {
		for _, kv := range k.Kv {
			UUID, err := uuid.FromBytes(kv.Value)
			if err != nil {
				return err
			}

			err = f(UUID)
			if err != nil {
				return err
			}

			err = s.s.Update(func(txn *badger.Txn) error {
				return txn.Delete(s.queuedItemKey(UUID))
			})
			if err != nil {
				return err
			}
		}

		return nil
	}

	return stream.Orchestrate(context.Background())
}
