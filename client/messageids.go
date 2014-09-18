package gnatt

import (
	"sync"
)

type mids struct {
	sync.RWMutex
	index map[uint16]Token
}

func (m *mids) freeId(id uint16) {
	m.Lock()
	defer m.Unlock()
	delete(m.index, id)
}

func (m *mids) getId(t Token) uint16 {
	m.Lock()
	defer m.Unlock()
	for i := uint16(1); i < uint16(65535); i++ {
		if _, ok := m.index[i]; !ok {
			m.index[i] = t
			return i
		}
	}
	return 0
}

func (m *mids) getToken(id uint16) Token {
	m.RLock()
	defer m.RUnlock()
	if token, ok := m.index[id]; ok {
		return token
	}
	return nil
}
