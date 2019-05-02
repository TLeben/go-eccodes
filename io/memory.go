package io

import "C"
import (
	"runtime"
	"strings"

	"github.com/tleben/go-eccodes/debug"
	"github.com/tleben/go-eccodes/native"
)

type Memory interface {
	isOpen() bool
	Native() native.Cbytes
	Close() error
	GetSize() native.Clong
}

type memory struct {
	debugID string
	memory  native.Cbytes
	size    native.Clong
}

func OpenMemory(content []byte, size int) (Memory, error) {
	m := &memory{memory: C.CBytes(content), debugID: strings.Join([]string{"memory=", string(len(content)), "', mode='", "byte_array", "'"}, ""), size: native.Clong(size)}
	runtime.SetFinalizer(m, memoryFinalizer)
	return m, nil
}

func (m *memory) isOpen() bool {
	return m.memory != nil
}

func (m *memory) Native() native.Cbytes {
	return m.memory
}

func (m *memory) GetSize() native.Clong {
	return m.size
}

func (m *memory) Close() error {
	defer func() { m.memory = nil }()
	native.Cmclose(m.Native())
	return nil
}

func memoryFinalizer(m *memory) {
	if m.isOpen() {
		debug.MemoryLeakLogger.Printf("'%s' is not closed", m.debugID)
		m.Close()
	}
}
