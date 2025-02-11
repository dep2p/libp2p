package multihash

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/dep2p/libp2p/multiformats/varint"
)

// maybe silly, but makes it so changing
// the table accidentally has to happen twice.
var tCodes = map[uint64]string{
	0x00:   "identity",
	0x11:   "sha1",
	0x12:   "sha2-256",
	0x13:   "sha2-512",
	0x14:   "sha3-512",
	0x15:   "sha3-384",
	0x16:   "sha3-256",
	0x17:   "sha3-224",
	0x56:   "dbl-sha2-256",
	0x22:   "murmur3-x64-64",
	0x1A:   "keccak-224",
	0x1B:   "keccak-256",
	0x1C:   "keccak-384",
	0x1D:   "keccak-512",
	0x1E:   "blake3",
	0x18:   "shake-128",
	0x19:   "shake-256",
	0x1100: "x11",
	0xd5:   "md5",
	0x1012: "sha2-256-trunc254-padded",
	0xb401: "poseidon-bls12_381-a2-fc1",
}

type TestCase struct {
	hex  string
	code uint64
	name string
}

var testCases = []TestCase{
	{"2c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae", 0x00, "identity"},
	{"", 0x00, "identity"},
	{"0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33", 0x11, "sha1"},
	{"0beec7b5", 0x11, "sha1"},
	{"2c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae", 0x12, "sha2-256"},
	{"2c26b46b", 0x12, "sha2-256"},
	{"2c26b46b68ffc68ff99b453c1d30413413", 0xb240, "blake2b-512"},
	{"243ddb9e", 0x22, "murmur3-x64-64"},
	{"f00ba4", 0x1b, "keccak-256"},
	{"f84e95cb5fbd2038863ab27d3cdeac295ad2d4ab96ad1f4b070c0bf36078ef08", 0x18, "shake-128"},
	{"1af97f7818a28edfdfce5ec66dbdc7e871813816d7d585fe1f12475ded5b6502b7723b74e2ee36f2651a10a8eaca72aa9148c3c761aaceac8f6d6cc64381ed39", 0x19, "shake-256"},
	{"4bca2b137edc580fe50a88983ef860ebaca36c857b1f492839d6d7392452a63c82cbebc68e3b70a2a1480b4bb5d437a7cba6ecf9d89f9ff3ccd14cd6146ea7e7", 0x14, "sha3-512"},
	{"d41d8cd98f00b204e9800998ecf8427e", 0xd5, "md5"},
	{"14fcb37dc45fa9a3c492557121bd4d461c0db40e5dcfcaa98498bd238486c307", 0x1012, "sha2-256-trunc254-padded"},
	{"14fcb37dc45fa9a3c492557121bd4d461c0db40e5dcfcaa98498bd238486c307", 0xb401, "poseidon-bls12_381-a2-fc1"},
	{"04e0bb39f30b1a3feb89f536c93be15055482df748674b00d26e5a75777702e9", 0x1e, "blake3"},
}

func (tc TestCase) Multihash() (Multihash, error) {
	ob, err := hex.DecodeString(tc.hex)
	if err != nil {
		return nil, err
	}

	pre := make([]byte, 2*binary.MaxVarintLen64)
	spot := pre
	n := binary.PutUvarint(spot, tc.code)
	spot = pre[n:]
	n += binary.PutUvarint(spot, uint64(len(ob)))

	nb := append(pre[:n], ob...)
	return Cast(nb)
}

func TestEncode(t *testing.T) {
	for _, tc := range testCases {
		ob, err := hex.DecodeString(tc.hex)
		if err != nil {
			t.Error(err)
			continue
		}

		pre := make([]byte, 2*binary.MaxVarintLen64)
		spot := pre
		n := binary.PutUvarint(spot, tc.code)
		spot = pre[n:]
		n += binary.PutUvarint(spot, uint64(len(ob)))

		nb := append(pre[:n], ob...)

		encC, err := Encode(ob, tc.code)
		if err != nil {
			t.Error(err)
			continue
		}

		if !bytes.Equal(encC, nb) {
			t.Error("encoded byte mismatch: ", encC, nb)
			t.Error(hex.Dump(nb))
		}

		encN, err := EncodeName(ob, tc.name)
		if err != nil {
			t.Error(err)
			continue
		}

		if !bytes.Equal(encN, nb) {
			t.Error("encoded byte mismatch: ", encN, nb)
		}

		h, err := tc.Multihash()
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(h, nb) {
			t.Error("Multihash func mismatch.")
		}
	}
}

func ExampleEncodeName() {
	// ignores errors for simplicity - don't do that at home.
	buf, _ := hex.DecodeString("0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33")
	mhbuf, _ := EncodeName(buf, "sha1")
	mhhex := hex.EncodeToString(mhbuf)
	fmt.Printf("hex: %v\n", mhhex)

	// Output:
	// hex: 11140beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33
}

func TestDecode(t *testing.T) {
	for _, tc := range testCases {
		ob, err := hex.DecodeString(tc.hex)
		if err != nil {
			t.Error(err)
			continue
		}
		pre := make([]byte, 2*binary.MaxVarintLen64)
		spot := pre
		n := binary.PutUvarint(spot, tc.code)
		spot = pre[n:]
		n += binary.PutUvarint(spot, uint64(len(ob)))

		nb := append(pre[:n], ob...)

		mustNotAllocateMore(t, 0, func() {
			dec, err := Decode(nb)
			if err != nil {
				t.Error(err)
				return
			}

			if dec.Code != tc.code {
				t.Error("decoded code mismatch: ", dec.Code, tc.code)
			}

			if dec.Name != tc.name {
				t.Error("decoded name mismatch: ", dec.Name, tc.name)
			}

			if dec.Length != len(ob) {
				t.Error("decoded length mismatch: ", dec.Length, len(ob))
			}

			if !bytes.Equal(dec.Digest, ob) {
				t.Error("decoded byte mismatch: ", dec.Digest, ob)
			}
		})
	}
}

func TestTable(t *testing.T) {
	for k, v := range tCodes {
		if Codes[k] != v {
			t.Error("Table mismatch: ", Codes[k], v)
		}
		if Names[v] != k {
			t.Error("Table mismatch: ", Names[v], k)
		}
	}

	for name, code := range Names {
		if tCodes[code] != name {
			if strings.HasPrefix(name, "blake") {
				// skip these
				continue
			}
			if name == "sha3" {
				// tested as "sha3-512"
				continue
			}
			t.Error("missing test case for: ", name)
		}
	}

	for code, name := range Codes {
		if tCodes[code] != name {
			if strings.HasPrefix(name, "blake") {
				// skip these
				continue
			}
			if name == "sha3" {
				// tested as "sha3-512"
				continue
			}
			t.Error("missing test case for: ", name)
		}
	}
}

func ExampleDecode() {
	// ignores errors for simplicity - don't do that at home.
	buf, _ := hex.DecodeString("0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33")
	mhbuf, _ := EncodeName(buf, "sha1")
	o, _ := Decode(mhbuf)
	mhhex := hex.EncodeToString(o.Digest)
	fmt.Printf("obj: %v 0x%x %d %s\n", o.Name, o.Code, o.Length, mhhex)

	// Output:
	// obj: sha1 0x11 20 0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33
}

func TestCast(t *testing.T) {
	for _, tc := range testCases {
		ob, err := hex.DecodeString(tc.hex)
		if err != nil {
			t.Error(err)
			continue
		}

		pre := make([]byte, 2*binary.MaxVarintLen64)
		spot := pre
		n := binary.PutUvarint(spot, tc.code)
		spot = pre[n:]
		n += binary.PutUvarint(spot, uint64(len(ob)))

		nb := append(pre[:n], ob...)

		mustNotAllocateMore(t, 0, func() {
			if _, err := Cast(nb); err != nil {
				t.Error(err)
				return
			}
		})

		// 1 for the error object.
		mustNotAllocateMore(t, 1, func() {
			if _, err = Cast(ob); err == nil {
				t.Error("cast failed to detect non-multihash")
				return
			}
		})
	}
}

func TestHex(t *testing.T) {
	for _, tc := range testCases {
		ob, err := hex.DecodeString(tc.hex)
		if err != nil {
			t.Error(err)
			continue
		}

		pre := make([]byte, 2*binary.MaxVarintLen64)
		spot := pre
		n := binary.PutUvarint(spot, tc.code)
		spot = pre[n:]
		n += binary.PutUvarint(spot, uint64(len(ob)))

		nb := append(pre[:n], ob...)

		hs := hex.EncodeToString(nb)
		mh, err := FromHexString(hs)
		if err != nil {
			t.Error(err)
			continue
		}

		if !bytes.Equal(mh, nb) {
			t.Error("FromHexString failed", nb, mh)
			continue
		}

		if mh.HexString() != hs {
			t.Error("Multihash.HexString failed", hs, mh.HexString())
			continue
		}
	}
}
func TestDecodeErrorInvalid(t *testing.T) {
	_, err := FromB58String("/ipfs/QmQTw94j68Dgakgtfd45bG3TZG6CAfc427UVRH4mugg4q4")
	if err != ErrInvalidMultihash {
		t.Fatalf("expected: %s, got %s\n", ErrInvalidMultihash, err)
	}
}

func TestBadVarint(t *testing.T) {
	_, err := Cast([]byte{129, 128, 128, 128, 128, 128, 128, 128, 128, 128, 129, 1})
	if err != varint.ErrOverflow {
		t.Error("expected error from varint longer than 64bits, got: ", err)
	}
	_, err = Cast([]byte{129, 128, 128})
	if err != varint.ErrUnderflow {
		t.Error("expected error from cut-off varint, got: ", err)
	}

	_, err = Cast([]byte{129, 0})
	if err != varint.ErrNotMinimal {
		t.Error("expected error non-minimal varint, got: ", err)
	}

	_, err = Cast([]byte{128, 0})
	if err != varint.ErrNotMinimal {
		t.Error("expected error non-minimal varint, got: ", err)
	}
}

func BenchmarkEncode(b *testing.B) {
	tc := testCases[0]
	ob, err := hex.DecodeString(tc.hex)
	if err != nil {
		b.Error(err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Encode(ob, tc.code)
	}
}

func BenchmarkDecode(b *testing.B) {
	tc := testCases[0]
	ob, err := hex.DecodeString(tc.hex)
	if err != nil {
		b.Error(err)
		return
	}

	pre := make([]byte, 2)
	pre[0] = byte(uint8(tc.code))
	pre[1] = byte(uint8(len(ob)))
	nb := append(pre, ob...)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decode(nb)
	}
}

func BenchmarkCast(b *testing.B) {
	tc := testCases[0]
	ob, err := hex.DecodeString(tc.hex)
	if err != nil {
		b.Error(err)
		return
	}

	pre := make([]byte, 2)
	pre[0] = byte(uint8(tc.code))
	pre[1] = byte(uint8(len(ob)))
	nb := append(pre, ob...)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Cast(nb)
	}
}
