package codes

import (
	"github.com/zachaller/go-errors/v2"

	cio "github.com/zachaller/go-eccodes/io"
	"github.com/zachaller/go-eccodes/native"
)

type ReaderMemory interface {
	GetSingleMessage() (Message, error)
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

func (m *memory) isOpen() bool {
	return m.memory != nil
}

func OpenMemory(m cio.Memory) (Memory, error) {
	return &memory{memory: m}, nil
}

func (m *memory) GetSingleMessage() (Message, error) {
	handle, err := native.Ccodes_handle_new_from_message(native.DefaultContext, m.memory.Native(), m.memory.GetSize())
	if err != nil {
		return nil, errors.Wrap(err, "failed create new handle from file")
	}
	return newMessage(handle), nil
}

func (m *memory) Close() {
	m.memory = nil
}
