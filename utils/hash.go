package utils

import (
	"errors"
	"fmt"
	"hash/fnv"
	"time"
)

func HashValue(val interface{}) (uint32, error) {
	h := fnv.New32a()

	switch v := val.(type) {
	case int:
		_, err := fmt.Fprintf(h, "%d", v)
		if err != nil {
			return 0, fmt.Errorf("failed to hash int: %w", err)
		}
	case string:
		_, err := h.Write([]byte(v))
		if err != nil {
			return 0, fmt.Errorf("failed to hash string: %w", err)
		}
	case bool:
		_, err := fmt.Fprintf(h, "%t", v)
		if err != nil {
			return 0, fmt.Errorf("failed to hash bool: %w", err)
		}
	case time.Time:
		_, err := h.Write([]byte(v.String()))
		if err != nil {
			return 0, fmt.Errorf("failed to hash time.Time: %w", err)
		}
	case nil:
		return 0, errors.New("cannot hash nil value")
	default:
		return 0, fmt.Errorf("unsupported type: %T", v)
	}

	return h.Sum32(), nil
}