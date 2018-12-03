package simplcrypto

// KeySet represents all the keys a node needs to operate
type KeySet struct {
	KeyPair *KeyPair // the server's keypair
	pairs   map[string]*KeyPair
	syms    map[string]*SymKey
}

// AddKeyPair adds a keyPair to the keySet
func (aks *KeySet) AddKeyPair(pair *KeyPair) {
	if aks.pairs == nil {
		aks.pairs = make(map[string]*KeyPair)
	}

	aks.pairs[pair.KID] = pair
}

// AddSymKey adds a SymKey to the keyset
func (aks *KeySet) AddSymKey(sym *SymKey) {
	if aks.syms == nil {
		aks.syms = make(map[string]*SymKey)
	}

	aks.syms[sym.KID] = sym
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

// SymKeyWithKID checks if a particular symkey exists in the keySet
func (aks *KeySet) SymKeyWithKID(kid string) *SymKey {
	sym, ok := aks.syms[kid]
	if !ok {
		return nil
	}

	return sym
}
