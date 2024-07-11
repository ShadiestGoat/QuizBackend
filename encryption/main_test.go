package encryption_test

import (
	"strconv"
	"testing"

	"whotfislucy.com/encryption"
)

const SECRET = "1234567890123456"

func Test(t *testing.T) {
	testAmt := encryption.BlockSize * 2 + 2 + 1

	inp := make([]string, testAmt)
	out := make([]string, testAmt)

	for i := 0; i < testAmt; i++ {
		dSize := i + 1
		d := ""

		for j := 0; j < dSize; j++ {
			d += strconv.Itoa(j)
		}

		inp = append(inp, d[:dSize])
	}

	err := encryption.Init(SECRET)
	if err != nil {
		panic(err)
	}

	for i := 0; i < testAmt; i++ {
		out[i] = encryption.Encrypt(inp[i])
	}

	for i := 0; i < testAmt; i++ {
		if encryption.Decrypt(out[i]) != inp[i] {
			t.Errorf("Failed to decrypt #%v. inp: '%v', dec: '%v'", i, inp[i], out[i])
		}

		j := testAmt - i - 1

		if encryption.Decrypt(out[j]) != inp[j] {
			t.Errorf("Failed to decrypt #%v (opposite dir). inp: '%v', dec: '%v'", j, inp[j], out[j])
		}
	}

	if err := encryption.Init(SECRET); err != nil {
		panic(err)
	}

	for i := 0; i < testAmt; i++ {
		enc := encryption.Encrypt(inp[i])

		if enc != out[i] {
			t.Errorf("Inconsistent encryption (#%d), enc1: '%v', enc2: '%v'", i, inp[i], enc)
		}
	}

	for i := 0; i < testAmt; i++ {
		if encryption.Decrypt(out[i]) != inp[i] {
			t.Errorf("Failed to decrypt #%v (pass #2). inp: '%v', dec: '%v'", i, inp[i], out[i])
		}

		j := testAmt - i - 1

		if encryption.Decrypt(out[j]) != inp[j] {
			t.Errorf("Failed to decrypt #%v (opposite dir, pass #2). inp: '%v', dec: '%v'", j, inp[j], out[j])
		}
	}
}
