package codes

import (
	"io"

	"github.com/zachaller/go-errors"

	"github.com/zachaller/go-eccodes/debug"
	cio "github.com/zachaller/go-eccodes/io"
	"github.com/zachaller/go-eccodes/native"
)

type ReaderMemory interface {
	Next() (Message, error)
}

type WriterMemory interface {
}

type Memory interface {
	ReaderMemory
	WriterMemory
	Close()
}

type memory struct {
	memory cio.Memory
}

type memoryIndexed struct {
	index native.Ccodes_index
}

//var memoryFilter = map[string]interface{}{}

func OpenMemory(m cio.Memory) (Memory, error) {
	return &memory{memory: m}, nil
}

/*
func OpenMemoryByPathWithFilter(path string, filter map[string]interface{}) (Memory, error) {
	if filter == nil {
		filter = emptyFilter
	}

	var k string
	for key, value := range filter {
		if len(k) > 0 {
			k += ","
		}
		k += key
		if value != nil {
			switch value.(type) {
			case int64, int:
				k += ":l"
			case float64, float32:
				k += ":d"
			case string:
				k += ":s"
			}
		}
	}

	i, err := native.Ccodes_index_new_from_file(native.DefaultContext, path, k)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create filtered index")
	}

	for key, value := range filter {
		if value != nil {
			err = nil
			switch value.(type) {
			case int64:
				err = native.Ccodes_index_select_long(i, key, value.(int64))
				if err != nil {
					err = errors.Wrapf(err, "failed to set filter condition '%s'=%d", key, value.(int64))
				}
			case int:
				err = native.Ccodes_index_select_long(i, key, int64(value.(int)))
				if err != nil {
					err = errors.Wrapf(err, "failed to set filter condition '%s'=%d", key, value.(int64))
				}
			case float64:
				err = native.Ccodes_index_select_double(i, key, value.(float64))
				if err != nil {
					err = errors.Wrapf(err, "failed to set filter condition '%s'=%f", key, value.(float64))
				}
			case float32:
				err = native.Ccodes_index_select_double(i, key, float64(value.(float32)))
				if err != nil {
					err = errors.Wrapf(err, "failed to set filter condition '%s'=%f", key, value.(float64))
				}
			case string:
				err = native.Ccodes_index_select_string(i, key, value.(string))
				if err != nil {
					err = errors.Wrapf(err, "failed to set filter condition '%s'='%s'", key, value.(string))
				}
			}
			if err != nil {
				native.Ccodes_index_delete(i)
				return nil, err
			}
		}
	}

	file := &fileIndexed{index: i}
	runtime.SetFinalizer(file, fileIndexedFinalizer)

	return file, nil
}
*/

func (m *memory) Next() (Message, error) {
	handle, err := native.Ccodes_handle_new_from_message_copy(native.DefaultContext, m.memory.Native(), m.memory.GetSize())
	if err != nil {
		if err == io.EOF {
			return nil, err
		}
		return nil, errors.Wrap(err, "failed create new handle from file")
	}

	return newMessage(handle), nil
}

func (m *memory) Close() {
	m.memory = nil
}

func (m *memoryIndexed) isOpen() bool {
	return m.index != nil
}

func (m *memoryIndexed) Next() (Message, error) {
	handle, err := native.Ccodes_handle_new_from_index(m.index)
	if err != nil {
		if err == io.EOF {
			return nil, err
		}
		return nil, errors.Wrap(err, "failed to create handle from index")
	}

	return newMessage(handle), nil
}

func (f *memoryIndexed) Close() {
	if f.isOpen() {
		defer func() { f.index = nil }()
		native.Ccodes_index_delete(f.index)
	}
}

func memoryIndexedFinalizer(f *fileIndexed) {
	if f.isOpen() {
		debug.MemoryLeakLogger.Print("file is not closed")
		f.Close()
	}
}
