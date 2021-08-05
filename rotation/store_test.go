package rotation

import "testing"

func buildKeystoreDecorator() KeyStore {
	return &KeyStoreDecorator{Wrapped: nil}
}

//really this is an enforcement of the expectation of meeting a contract.  intended to easy the number of message when
//the contract changes
func TestKeyStoreDecorator_castableToKeyStore(t *testing.T) {
	buildKeystoreDecorator()
}
