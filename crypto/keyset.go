package crypto

// KeySet represents all the keys a node needs to operate
type KeySet struct {
	GlobalKey *SymKey
	KeyPair   *KeyPair
	pairs     map[string]*KeyPair
}

// AddKeyPair adds a keyPair to the keySet
func (aks *KeySet) AddKeyPair(pair *KeyPair) {
	if aks.pairs == nil {
		aks.pairs = make(map[string]*KeyPair)
	}

	aks.pairs[pair.KID] = pair
}

// KeyPairWithKID checks if a particular keyPair exists in the keySet
func (aks *KeySet) KeyPairWithKID(kid string) *KeyPair {
	if aks.KeyPair.KID == kid {
		return aks.KeyPair
	}

	pair, ok := aks.pairs[kid]
	if !ok {
		return nil
	}

	return pair
}
