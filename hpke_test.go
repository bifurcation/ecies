package hpke

import (
	"bytes"
	"crypto/rand"
	"testing"
)

var (
	psk      = []byte("mellon")
	pskID    = []byte("Ennyn Durin aran Moria")
	original = []byte("Beauty is truth, truth beauty")
	aad      = []byte("that is all // Ye know on earth, and all ye need to know")
	info     = []byte("Ode on a Grecian Urn")
	rtts     = 10
)

func roundTrip(t *testing.T, id uint16, enc *EncryptContext, dec *DecryptContext) {
	for range make([]struct{}, rtts) {
		encrypted := enc.Seal(aad, original)
		decrypted, err := dec.Open(aad, encrypted)
		if err != nil {
			t.Fatalf("[%d] Error in Open: %s", id, err)
		}

		if !bytes.Equal(decrypted, original) {
			t.Fatalf("[%d] Incorrect decryption: [%x] != [%x]", id, decrypted, original)
		}
	}
}

func TestBase(t *testing.T) {
	for id, suite := range ciphersuites {
		skR, pkR, err := suite.KEM.GenerateKeyPair(rand.Reader)
		if err != nil {
			t.Fatalf("[%d] Error generating DH key pair: %s", id, err)
		}

		enc, ctxI, err := SetupBaseI(suite, rand.Reader, pkR, info)
		if err != nil {
			t.Fatalf("[%d] Error in SetupBaseI: %s", id, err)
		}

		ctxR, err := SetupBaseR(suite, skR, enc, info)
		if err != nil {
			t.Fatalf("[%d] Error in SetupBaseI: %s", id, err)
		}

		roundTrip(t, id, ctxI, ctxR)
	}
}

func TestPSK(t *testing.T) {
	for id, suite := range ciphersuites {
		skR, pkR, err := suite.KEM.GenerateKeyPair(rand.Reader)
		if err != nil {
			t.Fatalf("[%d] Error generating DH key pair: %s", id, err)
		}

		enc, ctxI, err := SetupPSKI(suite, rand.Reader, pkR, psk, pskID, info)
		if err != nil {
			t.Fatalf("[%d] Error in SetupPSKI: %s", id, err)
		}

		ctxR, err := SetupPSKR(suite, skR, enc, psk, pskID, info)
		if err != nil {
			t.Fatalf("[%d] Error in SetupBaseI: %s", id, err)
		}

		roundTrip(t, id, ctxI, ctxR)
	}
}

func TestAuth(t *testing.T) {
	for id, suite := range ciphersuites {
		skI, pkI, err := suite.KEM.GenerateKeyPair(rand.Reader)
		if err != nil {
			t.Fatalf("[%d] Error generating initiator DH key pair: %s", id, err)
		}

		skR, pkR, err := suite.KEM.GenerateKeyPair(rand.Reader)
		if err != nil {
			t.Fatalf("[%d] Error generating responder DH key pair: %s", id, err)
		}

		enc, ctxI, err := SetupAuthI(suite, rand.Reader, pkR, skI, info)
		if err != nil {
			t.Fatalf("[%d] Error in SetupAuthI: %s", id, err)
		}

		ctxR, err := SetupAuthR(suite, skR, pkI, enc, info)
		if err != nil {
			t.Fatalf("[%d] Error in SetupBaseI: %s", id, err)
		}

		roundTrip(t, id, ctxI, ctxR)
	}
}

func TestPSKAuth(t *testing.T) {
	for id, suite := range ciphersuites {
		skI, pkI, err := suite.KEM.GenerateKeyPair(rand.Reader)
		if err != nil {
			t.Fatalf("[%d] Error generating initiator DH key pair: %s", id, err)
		}

		skR, pkR, err := suite.KEM.GenerateKeyPair(rand.Reader)
		if err != nil {
			t.Fatalf("[%d] Error generating responder DH key pair: %s", id, err)
		}

		enc, ctxI, err := SetupPSKAuthI(suite, rand.Reader, pkR, skI, psk, pskID, info)
		if err != nil {
			t.Fatalf("[%d] Error in SetupPSKAuthI: %s", id, err)
		}

		ctxR, err := SetupPSKAuthR(suite, skR, pkI, enc, psk, pskID, info)
		if err != nil {
			t.Fatalf("[%d] Error in SetupBaseI: %s", id, err)
		}

		roundTrip(t, id, ctxI, ctxR)
	}
}
