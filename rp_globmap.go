package replacer

type gmap struct {
	m map[string][]byte
}

func newGMap() *gmap {
	return &gmap{
		m: make(map[string][]byte),
	}
}

func (m *gmap) get(key string) []byte {
	if val, ok := m.m[key]; ok {
		return val
	}

	return nil
}

func (m *gmap) set(key string, value []byte) {
	m.m[key] = value
}
