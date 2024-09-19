package envcfg

import (
	"os"
	"testing"
)

func TestParseJWT(t *testing.T) {
	in := "cd3efUoLdh8TNTWFG8DAqGjgUVFPZB12554iLS9Mmy5Vc8bXbjXxGccDuWLUWaB8ix0SR7mc2B1Uh3TuqmX/qg=="
	os.Setenv("JWT_SECRET_KEY", in)

	jwt := parseJWT()

	if len(jwt.Secret) != 64 {
		t.Fatalf("length of the []byte should be 64, but got %d\n", len(jwt.Secret))
	}
}
