package discov

import (
	"context"
	"errors"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/discov/internal"
	"github.com/zeromicro/go-zero/core/lang"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stringx"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"
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

func init() {
	logx.Disable()
}

func TestPublisher_register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	const id = 1
	cli := internal.NewMockEtcdClient(ctrl)
	restore := setMockClient(cli)
	defer restore()
	cli.EXPECT().Ctx().AnyTimes()
	cli.EXPECT().Grant(gomock.Any(), timeToLive).Return(&clientv3.LeaseGrantResponse{
		ID: id,
	}, nil)
	cli.EXPECT().Put(gomock.Any(), makeEtcdKey("thekey", id), "thevalue", gomock.Any())
	pub := NewPublisher(nil, "thekey", "thevalue",
		WithPubEtcdAccount(stringx.Rand(), "bar"))
	_, err := pub.register(cli)
	assert.Nil(t, err)
}

func TestPublisher_registerWithOptions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	const id = 2
	cli := internal.NewMockEtcdClient(ctrl)
	restore := setMockClient(cli)
	defer restore()
	cli.EXPECT().Ctx().AnyTimes()
	cli.EXPECT().Grant(gomock.Any(), timeToLive).Return(&clientv3.LeaseGrantResponse{
		ID: 1,
	}, nil)
	cli.EXPECT().Put(gomock.Any(), makeEtcdKey("thekey", id), "thevalue", gomock.Any())

	certFile := createTempFile(t, []byte(certContent))
	defer os.Remove(certFile)
	keyFile := createTempFile(t, []byte(keyContent))
	defer os.Remove(keyFile)
	caFile := createTempFile(t, []byte(caContent))
	defer os.Remove(caFile)
	pub := NewPublisher(nil, "thekey", "thevalue", WithId(id),
		WithPubEtcdTLS(certFile, keyFile, caFile, true))
	_, err := pub.register(cli)
	assert.Nil(t, err)
}

func TestPublisher_registerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cli := internal.NewMockEtcdClient(ctrl)
	restore := setMockClient(cli)
	defer restore()
	cli.EXPECT().Ctx().AnyTimes()
	cli.EXPECT().Grant(gomock.Any(), timeToLive).Return(nil, errors.New("error"))
	pub := NewPublisher(nil, "thekey", "thevalue")
	val, err := pub.register(cli)
	assert.NotNil(t, err)
	assert.Equal(t, clientv3.NoLease, val)
}

func TestPublisher_revoke(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	const id clientv3.LeaseID = 1
	cli := internal.NewMockEtcdClient(ctrl)
	restore := setMockClient(cli)
	defer restore()
	cli.EXPECT().Ctx().AnyTimes()
	cli.EXPECT().Revoke(gomock.Any(), id)
	pub := NewPublisher(nil, "thekey", "thevalue")
	pub.lease = id
	pub.revoke(cli)
}

func TestPublisher_revokeError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	const id clientv3.LeaseID = 1
	cli := internal.NewMockEtcdClient(ctrl)
	restore := setMockClient(cli)
	defer restore()
	cli.EXPECT().Ctx().AnyTimes()
	cli.EXPECT().Revoke(gomock.Any(), id).Return(nil, errors.New("error"))
	pub := NewPublisher(nil, "thekey", "thevalue")
	pub.lease = id
	pub.revoke(cli)
}

func TestPublisher_keepAliveAsyncError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	const id clientv3.LeaseID = 1
	cli := internal.NewMockEtcdClient(ctrl)
	restore := setMockClient(cli)
	defer restore()
	cli.EXPECT().Ctx().AnyTimes()
	cli.EXPECT().KeepAlive(gomock.Any(), id).Return(nil, errors.New("error"))
	pub := NewPublisher(nil, "thekey", "thevalue")
	pub.lease = id
	assert.NotNil(t, pub.keepAliveAsync(cli))
}

func TestPublisher_keepAliveAsyncQuit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	const id clientv3.LeaseID = 1
	cli := internal.NewMockEtcdClient(ctrl)
	cli.EXPECT().ActiveConnection()
	cli.EXPECT().Close()
	defer cli.Close()
	cli.ActiveConnection()
	restore := setMockClient(cli)
	defer restore()
	cli.EXPECT().Ctx().AnyTimes()
	cli.EXPECT().KeepAlive(gomock.Any(), id)
	var wg sync.WaitGroup
	wg.Add(1)
	cli.EXPECT().Revoke(gomock.Any(), id).Do(func(_, _ any) {
		wg.Done()
	})
	pub := NewPublisher(nil, "thekey", "thevalue")
	pub.lease = id
	pub.Stop()
	assert.Nil(t, pub.keepAliveAsync(cli))
	wg.Wait()
}

func TestPublisher_keepAliveAsyncPause(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	const id clientv3.LeaseID = 1
	cli := internal.NewMockEtcdClient(ctrl)
	restore := setMockClient(cli)
	defer restore()
	cli.EXPECT().Ctx().AnyTimes()
	cli.EXPECT().KeepAlive(gomock.Any(), id)
	pub := NewPublisher(nil, "thekey", "thevalue")
	var wg sync.WaitGroup
	wg.Add(1)
	cli.EXPECT().Revoke(gomock.Any(), id).Do(func(_, _ any) {
		pub.Stop()
		wg.Done()
	})
	pub.lease = id
	assert.Nil(t, pub.keepAliveAsync(cli))
	pub.Pause()
	wg.Wait()
}

func TestPublisher_Resume(t *testing.T) {
	publisher := new(Publisher)
	publisher.resumeChan = make(chan lang.PlaceholderType)
	go func() {
		publisher.Resume()
	}()
	go func() {
		time.Sleep(time.Minute)
		t.Fail()
	}()
	<-publisher.resumeChan
}

func TestPublisher_keepAliveAsync(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	const id clientv3.LeaseID = 1
	conn := createMockConn(t)
	defer conn.Close()
	cli := internal.NewMockEtcdClient(ctrl)
	cli.EXPECT().ActiveConnection().Return(conn).AnyTimes()
	cli.EXPECT().Close()
	defer cli.Close()
	cli.ActiveConnection()
	restore := setMockClient(cli)
	defer restore()
	cli.EXPECT().Ctx().AnyTimes()
	cli.EXPECT().KeepAlive(gomock.Any(), id)
	cli.EXPECT().Grant(gomock.Any(), timeToLive).Return(&clientv3.LeaseGrantResponse{
		ID: 1,
	}, nil)
	cli.EXPECT().Put(gomock.Any(), makeEtcdKey("thekey", int64(id)), "thevalue", gomock.Any())
	var wg sync.WaitGroup
	wg.Add(1)
	cli.EXPECT().Revoke(gomock.Any(), id).Do(func(_, _ any) {
		wg.Done()
	})
	pub := NewPublisher([]string{"the-endpoint"}, "thekey", "thevalue")
	pub.lease = id
	assert.Nil(t, pub.KeepAlive())
	pub.Stop()
	wg.Wait()
}

func createMockConn(t *testing.T) *grpc.ClientConn {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Error while listening. Err: %v", err)
	}
	defer lis.Close()
	lisAddr := resolver.Address{Addr: lis.Addr().String()}
	lisDone := make(chan struct{})
	dialDone := make(chan struct{})
	// 1st listener accepts the connection and then does nothing
	go func() {
		defer close(lisDone)
		conn, err := lis.Accept()
		if err != nil {
			t.Errorf("Error while accepting. Err: %v", err)
			return
		}
		framer := http2.NewFramer(conn, conn)
		if err := framer.WriteSettings(http2.Setting{}); err != nil {
			t.Errorf("Error while writing settings. Err: %v", err)
			return
		}
		<-dialDone // Close conn only after dial returns.
	}()

	r := manual.NewBuilderWithScheme("whatever")
	r.InitialState(resolver.State{Addresses: []resolver.Address{lisAddr}})
	client, err := grpc.DialContext(context.Background(), r.Scheme()+":///test.server",
		grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithResolvers(r))
	close(dialDone)
	if err != nil {
		t.Fatalf("Dial failed. Err: %v", err)
	}

	timeout := time.After(1 * time.Second)
	select {
	case <-timeout:
		t.Fatal("timed out waiting for server to finish")
	case <-lisDone:
	}

	return client
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
