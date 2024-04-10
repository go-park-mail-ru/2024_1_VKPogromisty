package hash_test

import (
	"socio/pkg/hash"
	"testing"
)

type HashTestCase struct {
	Password string
	Salt     []byte
	Expected string
	Match    bool
}

var HashTestCases = map[string]HashTestCase{
	"match password": {
		Password: "admin",
		Salt:     []byte("salt"),
		Expected: "3c4a79782143337be4492be072abcfe979dd703c00541a8c39a0f3df4bab2029c050cf46fddc47090b5b04ac537b3e78189e3de16e601e859f95c51ac9f6dafb",
		Match:    true,
	},
}

func TestHashPassword(t *testing.T) {
	for name, tc := range HashTestCases {
		t.Run(name, func(t *testing.T) {
			hashedPassword := hash.HashPassword(tc.Password, tc.Salt)

			if hashedPassword != tc.Expected || !hash.MatchPasswords(hashedPassword, tc.Password, tc.Salt) {
				t.Errorf("wrong hashedPassword: got %s, expected %s", hashedPassword, tc.Expected)
				return
			}
		})
	}
}
