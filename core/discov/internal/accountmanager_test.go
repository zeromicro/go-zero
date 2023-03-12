package internal

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stringx"
)

const (
	certContent = `-----BEGIN CERTIFICATE-----
MIIDazCCAlOgAwIBAgIUEg9GVO2oaPn+YSmiqmFIuAo10WIwDQYJKoZIhvcNAQEM
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAgFw0yMzAzMTExMzIxMjNaGA8yMTIz
MDIxNTEzMjEyM1owRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUx
ITAfBgNVBAoMGEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDCCASIwDQYJKoZIhvcN
AQEBBQADggEPADCCAQoCggEBALplXlWsIf0O/IgnIplmiZHKGnxyfyufyE2FBRNk
OofRqbKuPH8GNqbkvZm7N29fwTDAQ+mViAggCkDht4hOzoWJMA7KYJt8JnTSWL48
M1lcrpc9DL2gszC/JF/FGvyANbBtLklkZPFBGdHUX14pjrT937wqPtm+SqUHSvRT
B7bmwmm2drRcmhpVm98LSlV7uQ2EgnJgsLjBPITKUejLmVLHfgX0RwQ2xIpX9pS4
FCe1BTacwl2gGp7Mje7y4Mfv3o0ArJW6Tuwbjx59ZXwb1KIP71b7bT04AVS8ZeYO
UMLKKuB5UR9x9Rn6cLXOTWBpcMVyzDgrAFLZjnE9LPUolZMCAwEAAaNRME8wHwYD
VR0jBBgwFoAUeW8w8pmhncbRgTsl48k4/7wnfx8wCQYDVR0TBAIwADALBgNVHQ8E
BAMCBPAwFAYDVR0RBA0wC4IJbG9jYWxob3N0MA0GCSqGSIb3DQEBDAUAA4IBAQAI
y9xaoS88CLPBsX6mxfcTAFVfGNTRW9VN9Ng1cCnUR+YGoXGM/l+qP4f7p8ocdGwK
iYZErVTzXYIn+D27//wpY3klJk3gAnEUBT3QRkStBw7XnpbeZ2oPBK+cmDnCnZPS
BIF1wxPX7vIgaxs5Zsdqwk3qvZ4Djr2wP7LabNWTLSBKgQoUY45Liw6pffLwcGF9
UKlu54bvGze2SufISCR3ib+I+FLvqpvJhXToZWYb/pfI/HccuCL1oot1x8vx6DQy
U+TYxlZsKS5mdNxAX3dqEkEMsgEi+g/tzDPXJImfeCGGBhIOXLm8SRypiuGdEbc9
xkWYxRPegajuEZGvCqVs
-----END CERTIFICATE-----`
	keyContent = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAumVeVawh/Q78iCcimWaJkcoafHJ/K5/ITYUFE2Q6h9Gpsq48
fwY2puS9mbs3b1/BMMBD6ZWICCAKQOG3iE7OhYkwDspgm3wmdNJYvjwzWVyulz0M
vaCzML8kX8Ua/IA1sG0uSWRk8UEZ0dRfXimOtP3fvCo+2b5KpQdK9FMHtubCabZ2
tFyaGlWb3wtKVXu5DYSCcmCwuME8hMpR6MuZUsd+BfRHBDbEilf2lLgUJ7UFNpzC
XaAansyN7vLgx+/ejQCslbpO7BuPHn1lfBvUog/vVvttPTgBVLxl5g5Qwsoq4HlR
H3H1Gfpwtc5NYGlwxXLMOCsAUtmOcT0s9SiVkwIDAQABAoIBAD5meTJNMgO55Kjg
ESExxpRcCIno+tHr5+6rvYtEXqPheOIsmmwb9Gfi4+Z3WpOaht5/Pz0Ppj6yGzyl
U//6AgGKb+BDuBvVcDpjwPnOxZIBCSHwejdxeQu0scSuA97MPS0XIAvJ5FEv7ijk
5Bht6SyGYURpECltHygoTNuGgGqmO+McCJRLE9L09lTBI6UQ/JQwWJqSr7wx6iPU
M1Ze/srIV+7cyEPu6i0DGjS1gSQKkX68Lqn1w6oE290O+OZvleO0gZ02fLDWCZke
aeD9+EU/Pw+rqm3H6o0szOFIpzhRp41FUdW9sybB3Yp3u7c/574E+04Z/e30LMKs
TCtE1QECgYEA3K7KIpw0NH2HXL5C3RHcLmr204xeBfS70riBQQuVUgYdmxak2ima
80RInskY8hRhSGTg0l+VYIH8cmjcUyqMSOELS5XfRH99r4QPiK8AguXg80T4VumY
W3Pf+zEC2ssgP/gYthV0g0Xj5m2QxktOF9tRw5nkg739ZR4dI9lm/iECgYEA2Dnf
uwEDGqHiQRF6/fh5BG/nGVMvrefkqx6WvTJQ3k/M/9WhxB+lr/8yH46TuS8N2b29
FoTf3Mr9T7pr/PWkOPzoY3P56nYbKU8xSwCim9xMzhBMzj8/N9ukJvXy27/VOz56
eQaKqnvdXNGtPJrIMDGHps2KKWlKLyAlapzjVTMCgYAA/W++tACv85g13EykfT4F
n0k4LbsGP9DP4zABQLIMyiY72eAncmRVjwrcW36XJ2xATOONTgx3gF3HjZzfaqNy
eD/6uNNllUTVEryXGmHgNHPL45VRnn6memCY2eFvZdXhM5W4y2PYaunY0MkDercA
+GTngbs6tBF88KOk04bYwQKBgFl68cRgsdkmnwwQYNaTKfmVGYzYaQXNzkqmWPko
xmCJo6tHzC7ubdG8iRCYHzfmahPuuj6EdGPZuSRyYFgJi5Ftz/nAN+84OxtIQ3zn
YWOgskQgaLh9YfsKsQ7Sf1NDOsnOnD5TX7UXl07fEpLe9vNCvAFiU8e5Y9LGudU5
4bYTAoGBAMdX3a3bXp4cZvXNBJ/QLVyxC6fP1Q4haCR1Od3m+T00Jth2IX2dk/fl
p6xiJT1av5JtYabv1dFKaXOS5s1kLGGuCCSKpkvFZm826aQ2AFm0XGqEQDLeei5b
A52Kpy/YJ+RkG4BTFtAooFq6DmA0cnoP6oPvG2h6XtDJwDTPInJb
-----END RSA PRIVATE KEY-----`
	caContent = `-----BEGIN CERTIFICATE-----
MIIDbTCCAlWgAwIBAgIUBJvFoCowKich7MMfseJ+DYzzirowDQYJKoZIhvcNAQEM
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAgFw0yMzAzMTExMzIxMDNaGA8yMTIz
MDIxNTEzMjEwM1owRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUx
ITAfBgNVBAoMGEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDCCASIwDQYJKoZIhvcN
AQEBBQADggEPADCCAQoCggEBAO4to2YMYj0bxgr2FCiweSTSFuPx33zSw2x/s9Wf
OR41bm2DFsyYT5f3sOIKlXZEdLmOKty2e3ho3yC0EyNpVHdykkkHT3aDI17quZax
kYi/URqqtl1Z08A22txolc04hAZisg2BypGi3vql81UW1t3zyloGnJoIAeXR9uca
ljP6Bk3bwsxoVBLi1JtHrO0hHLQaeHmKhAyrys06X0LRdn7Px48yRZlt6FaLSa8X
YiRM0G44bVy/h6BkoQjMYGwVmCVk6zjJ9U7ZPFqdnDMNxAfR+hjDnYodqdLDMTTR
1NPVrnEnNwFx0AMLvgt/ba/45vZCEAmSZnFXFAJJcM7ai9ECAwEAAaNTMFEwHQYD
VR0OBBYEFHlvMPKZoZ3G0YE7JePJOP+8J38fMB8GA1UdIwQYMBaAFHlvMPKZoZ3G
0YE7JePJOP+8J38fMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQEMBQADggEB
AMX8dNulADOo9uQgBMyFb9TVra7iY0zZjzv4GY5XY7scd52n6CnfAPvYBBDnTr/O
BgNp5jaujb4+9u/2qhV3f9n+/3WOb2CmPehBgVSzlXqHeQ9lshmgwZPeem2T+8Tm
Nnc/xQnsUfCFszUDxpkr55+aLVM22j02RWqcZ4q7TAaVYL+kdFVMc8FoqG/0ro6A
BjE/Qn0Nn7ciX1VUjDt8l+k7ummPJTmzdi6i6E4AwO9dzrGNgGJ4aWL8cC6xYcIX
goVIRTFeONXSDno/oPjWHpIPt7L15heMpKBHNuzPkKx2YVqPHE5QZxWfS+Lzgx+Q
E2oTTM0rYKOZ8p6000mhvKI=
-----END CERTIFICATE-----`
)

func TestAccount(t *testing.T) {
	endpoints := []string{
		"192.168.0.2:2379",
		"192.168.0.3:2379",
		"192.168.0.4:2379",
	}
	username := "foo" + stringx.Rand()
	password := "bar"
	anotherPassword := "any"

	_, ok := GetAccount(endpoints)
	assert.False(t, ok)

	AddAccount(endpoints, username, password)
	account, ok := GetAccount(endpoints)
	assert.True(t, ok)
	assert.Equal(t, username, account.User)
	assert.Equal(t, password, account.Pass)

	AddAccount(endpoints, username, anotherPassword)
	account, ok = GetAccount(endpoints)
	assert.True(t, ok)
	assert.Equal(t, username, account.User)
	assert.Equal(t, anotherPassword, account.Pass)
}

func TestTLSMethods(t *testing.T) {
	certFile := createTempFile(t, []byte(certContent))
	defer os.Remove(certFile)
	keyFile := createTempFile(t, []byte(keyContent))
	defer os.Remove(keyFile)
	caFile := createTempFile(t, []byte(caContent))
	defer os.Remove(caFile)

	assert.NoError(t, AddTLS([]string{"foo"}, certFile, keyFile, caFile, false))
	cfg, ok := GetTLS([]string{"foo"})
	assert.True(t, ok)
	assert.NotNil(t, cfg)

	assert.Error(t, AddTLS([]string{"bar"}, "bad-file", keyFile, caFile, false))
	assert.Error(t, AddTLS([]string{"bar"}, certFile, keyFile, "bad-file", false))
}

func createTempFile(t *testing.T, body []byte) string {
	tmpFile, err := os.CreateTemp(os.TempDir(), "go-unit-*.tmp")
	if err != nil {
		t.Fatal(err)
	}

	tmpFile.Close()
	if err = os.WriteFile(tmpFile.Name(), body, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	return tmpFile.Name()
}
