package codec

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/fs"
)

const (
	priKey = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQC4TJk3onpqb2RYE3wwt23J9SHLFstHGSkUYFLe+nl1dEKHbD+/
Zt95L757J3xGTrwoTc7KCTxbrgn+stn0w52BNjj/kIE2ko4lbh/v8Fl14AyVR9ms
fKtKOnhe5FCT72mdtApr+qvzcC3q9hfXwkyQU32pv7q5UimZ205iKSBmgQIDAQAB
AoGAM5mWqGIAXj5z3MkP01/4CDxuyrrGDVD5FHBno3CDgyQa4Gmpa4B0/ywj671B
aTnwKmSmiiCN2qleuQYASixes2zY5fgTzt+7KNkl9JHsy7i606eH2eCKzsUa/s6u
WD8V3w/hGCQ9zYI18ihwyXlGHIgcRz/eeRh+nWcWVJzGOPUCQQD5nr6It/1yHb1p
C6l4fC4xXF19l4KxJjGu1xv/sOpSx0pOqBDEX3Mh//FU954392rUWDXV1/I65BPt
TLphdsu3AkEAvQJ2Qay/lffFj9FaUrvXuftJZ/Ypn0FpaSiUh3Ak3obBT6UvSZS0
bcYdCJCNHDtBOsWHnIN1x+BcWAPrdU7PhwJBAIQ0dUlH2S3VXnoCOTGc44I1Hzbj
Rc65IdsuBqA3fQN2lX5vOOIog3vgaFrOArg1jBkG1wx5IMvb/EnUN2pjVqUCQCza
KLXtCInOAlPemlCHwumfeAvznmzsWNdbieOZ+SXVVIpR6KbNYwOpv7oIk3Pfm9sW
hNffWlPUKhW42Gc+DIECQQDmk20YgBXwXWRM5DRPbhisIV088N5Z58K9DtFWkZsd
OBDT3dFcgZONtlmR1MqZO0pTh30lA4qovYj3Bx7A8i36
-----END RSA PRIVATE KEY-----`
	pubKey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC4TJk3onpqb2RYE3wwt23J9SHL
FstHGSkUYFLe+nl1dEKHbD+/Zt95L757J3xGTrwoTc7KCTxbrgn+stn0w52BNjj/
kIE2ko4lbh/v8Fl14AyVR9msfKtKOnhe5FCT72mdtApr+qvzcC3q9hfXwkyQU32p
v7q5UimZ205iKSBmgQIDAQAB
-----END PUBLIC KEY-----`
	testBody = `this is the content`
)

func TestCryption(t *testing.T) {
	enc, err := NewRsaEncrypter([]byte(pubKey))
	assert.Nil(t, err)
	ret, err := enc.Encrypt([]byte(testBody))
	assert.Nil(t, err)

	file, err := fs.TempFilenameWithText(priKey)
	assert.Nil(t, err)
	defer os.Remove(file)
	dec, err := NewRsaDecrypter(file)
	assert.Nil(t, err)
	actual, err := dec.Decrypt(ret)
	assert.Nil(t, err)
	assert.Equal(t, testBody, string(actual))

	actual, err = dec.DecryptBase64(base64.StdEncoding.EncodeToString(ret))
	assert.Nil(t, err)
	assert.Equal(t, testBody, string(actual))
}

func TestBadPubKey(t *testing.T) {
	_, err := NewRsaEncrypter([]byte("foo"))
	assert.Equal(t, ErrPublicKey, err)
}
