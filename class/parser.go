package class

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"golang.org/x/exp/constraints"
)

type errReader struct {
	r   io.Reader
	err error

	name string
}

func item(e *errReader, name string, f func(e *errReader) bool) bool {
	if e.err != nil {
		return false
	}
	e.name = name
	return f(e)
}

func parent(e *errReader, parent string) *errReader {
	return &errReader{r: e.r, err: e.err, name: parent}
}

func child(e *errReader, fieldName string, f func(e *errReader) bool) bool {
	if e.err != nil {
		return false
	}
	return f(&errReader{r: e.r, err: e.err, name: fmt.Sprintf("%s.%s", e.name, fieldName)})
}

func integer[T constraints.Integer](data *T, vs ...validator[T]) func(e *errReader) bool {
	return func(e *errReader) bool {
		return readInteger(e, data, vs...)
	}
}

func readInteger[T constraints.Integer](e *errReader, data *T, vs ...validator[T]) bool {
	if e.err != nil {
		return false
	}
	if err := binary.Read(e.r, binary.BigEndian, data); err != nil {
		e.err = fmt.Errorf("fail to parse %s: %w", e.name, err)
		return false
	}

	for _, v := range vs {
		if err := v.validate(*data, e.name); err != nil {
			e.err = err
			return false
		}
	}
	return true
}

func bytes(bytes []byte, vs ...validator[[]byte]) func(e *errReader) bool {
	return func(e *errReader) bool {
		return readBytes(e, bytes, vs...)
	}
}

func readBytes(e *errReader, bytes []byte, vs ...validator[[]byte]) bool {
	if e.err != nil {
		return false
	}
	if n, err := e.r.Read(bytes); err != nil {
		e.err = fmt.Errorf("fail to parse %s: %w", e.name, err)
		return false
	} else if n == 0 {
		e.err = errors.New("fail to parse %s")
		return false
	}

	for _, v := range vs {
		if err := v.validate(bytes, e.name); err != nil {
			e.err = err
			return false
		}
	}
	return true
}

func entries[T any](es []T, f func(e *errReader) T, vs ...validator[T]) func(e *errReader) bool {
	return func(e *errReader) bool {
		return readEntries(e, es, f, vs...)
	}
}

func readEntries[T any](e *errReader, es []T, f func(e *errReader) T, vs ...validator[T]) bool {
	if e.err != nil {
		return false
	}
	for i := range es {
		entry := f(e)
		if e.err != nil {
			e.err = fmt.Errorf("fail to parse %s[%d]: %w", e.name, i, e.err)
			return false
		}
		for _, v := range vs {
			if err := v.validate(entry, e.name); err != nil {
				e.err = err
				return false
			}
		}
		es[i] = entry
	}
	return true
}
