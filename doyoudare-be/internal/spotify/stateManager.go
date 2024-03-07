// internal/spotify/stateManager.go
package spotify

var stateVerifierMap = make(map[string]string)

// SaveStateVerifier saves a state and its associated code verifier
func SaveStateVerifier(state, verifier string) {
	stateVerifierMap[state] = verifier
}

// GetVerifierByState retrieves the code verifier associated with a state
func GetVerifierByState(state string) (verifier string, ok bool) {
	verifier, ok = stateVerifierMap[state]
	return verifier, ok
}

// DeleteStateVerifier removes the state and its associated verifier from the map
func DeleteStateVerifier(state string) {
	delete(stateVerifierMap, state)
}
