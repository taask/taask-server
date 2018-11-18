package crypto

import "encoding/json"

// KeyTypePair and others represent the type of key used to encrypt or sign a message
const (
	KeyTypePair      = "astro.key.pair"
	KeyTypeSymmetric = "astro.key.sym"
)

// ToJSON serializes a message to JSON
func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// MessageFromJSON deserializes a message from JSON
func MessageFromJSON(src []byte) (*Message, error) {
	m := &Message{}
	if err := json.Unmarshal(src, m); err != nil {
		return nil, err
	}

	return m, nil
}

// ToJSON serializes a signature to JSON
func (s *Signature) ToJSON() ([]byte, error) {
	return json.Marshal(s)
}

// SignatureFromJSON deserializes a signature from JSON
func SignatureFromJSON(src []byte) (*Signature, error) {
	s := &Signature{}
	if err := json.Unmarshal(src, s); err != nil {
		return nil, err
	}

	return s, nil
}
