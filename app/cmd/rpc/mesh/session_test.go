package mesh

import (
	"bytes"
	"github.com/alitto/pond"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/puzpuzpuz/xsync"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestDispatchResponse_Contains(t *testing.T) {
	type fields struct {
		BlockHeight int64
		Session     DispatchSession
	}
	type args struct {
		addr string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sn := DispatchResponse{
				BlockHeight: tt.fields.BlockHeight,
				Session:     tt.fields.Session,
			}
			if got := sn.Contains(tt.args.addr); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDispatchResponse_GetSupportedNodes(t *testing.T) {
	type fields struct {
		BlockHeight int64
		Session     DispatchSession
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sn := DispatchResponse{
				BlockHeight: tt.fields.BlockHeight,
				Session:     tt.fields.Session,
			}
			if got := sn.GetSupportedNodes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSupportedNodes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDispatchResponse_ShouldKeep(t *testing.T) {
	type fields struct {
		BlockHeight int64
		Session     DispatchSession
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sn := DispatchResponse{
				BlockHeight: tt.fields.BlockHeight,
				Session:     tt.fields.Session,
			}
			if got := sn.ShouldKeep(); got != tt.want {
				t.Errorf("ShouldKeep() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitializeSessionStorage(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitializeSessionStorage()
		})
	}
}

func TestNodeSession_CountRelay(t *testing.T) {
	type fields struct {
		PubKey          string
		RemainingRelays int64
		RelayMeta       *pocketTypes.RelayMeta
		Validated       bool
		RetryTimes      int
		IsValid         bool
		Error           *SdkErrorResponse
		Session         *Session
	}
	tests := []struct {
		name   string
		fields fields
		want   *NodeSession
		want1  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ns := &NodeSession{
				PubKey:          tt.fields.PubKey,
				RemainingRelays: tt.fields.RemainingRelays,
				RelayMeta:       tt.fields.RelayMeta,
				Validated:       tt.fields.Validated,
				RetryTimes:      tt.fields.RetryTimes,
				IsValid:         tt.fields.IsValid,
				Error:           tt.fields.Error,
				Session:         tt.fields.Session,
			}
			got, got1 := ns.CountRelay()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CountRelay() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CountRelay() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNodeSession_ReScheduleValidationTask(t *testing.T) {
	type fields struct {
		PubKey          string
		RemainingRelays int64
		RelayMeta       *pocketTypes.RelayMeta
		Validated       bool
		RetryTimes      int
		IsValid         bool
		Error           *SdkErrorResponse
		Session         *Session
	}
	type args struct {
		session        *Session
		servicerPubKey string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ns := &NodeSession{
				PubKey:          tt.fields.PubKey,
				RemainingRelays: tt.fields.RemainingRelays,
				RelayMeta:       tt.fields.RelayMeta,
				Validated:       tt.fields.Validated,
				RetryTimes:      tt.fields.RetryTimes,
				IsValid:         tt.fields.IsValid,
				Error:           tt.fields.Error,
				Session:         tt.fields.Session,
			}
			ns.ReScheduleValidationTask(tt.args.session, tt.args.servicerPubKey)
		})
	}
}

func TestSessionStorage_AddSessionToValidate(t *testing.T) {
	type fields struct {
		Sessions         *xsync.MapOf[string, *Session]
		ValidationWorker *pond.WorkerPool
	}
	type args struct {
		relay *pocketTypes.Relay
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Session
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := &SessionStorage{
				Sessions:         tt.fields.Sessions,
				ValidationWorker: tt.fields.ValidationWorker,
			}
			got, err := ss.AddSessionToValidate(tt.args.relay)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddSessionToValidate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddSessionToValidate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSessionStorage_GetNodeSession(t *testing.T) {
	type fields struct {
		Sessions         *xsync.MapOf[string, *Session]
		ValidationWorker *pond.WorkerPool
	}
	type args struct {
		relay *pocketTypes.Relay
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *NodeSession
		want1  *SdkErrorResponse
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := &SessionStorage{
				Sessions:         tt.fields.Sessions,
				ValidationWorker: tt.fields.ValidationWorker,
			}
			got, got1 := ss.GetNodeSession(tt.args.relay)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNodeSession() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetNodeSession() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSessionStorage_GetSession(t *testing.T) {
	type fields struct {
		Sessions         *xsync.MapOf[string, *Session]
		ValidationWorker *pond.WorkerPool
	}
	type args struct {
		relay *pocketTypes.Relay
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Session
		want1  *SdkErrorResponse
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := &SessionStorage{
				Sessions:         tt.fields.Sessions,
				ValidationWorker: tt.fields.ValidationWorker,
			}
			got, got1 := ss.GetSession(tt.args.relay)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSession() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetSession() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSessionStorage_NewSessionFromRelay(t *testing.T) {
	type fields struct {
		Sessions         *xsync.MapOf[string, *Session]
		ValidationWorker *pond.WorkerPool
	}
	type args struct {
		relay *pocketTypes.Relay
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := &SessionStorage{
				Sessions:         tt.fields.Sessions,
				ValidationWorker: tt.fields.ValidationWorker,
			}
			if got := ss.NewSessionFromRelay(tt.args.relay); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSessionFromRelay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSessionStorage_ShouldAssumeOptimisticSession(t *testing.T) {
	type fields struct {
		Sessions         *xsync.MapOf[string, *Session]
		ValidationWorker *pond.WorkerPool
	}
	type args struct {
		relay        *pocketTypes.Relay
		servicerNode *fullNode
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := &SessionStorage{
				Sessions:         tt.fields.Sessions,
				ValidationWorker: tt.fields.ValidationWorker,
			}
			if got := ss.ShouldAssumeOptimisticSession(tt.args.relay, tt.args.servicerNode); got != tt.want {
				t.Errorf("ShouldAssumeOptimisticSession() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_GetDispatch(t *testing.T) {
	type fields struct {
		Hash         string
		AppPublicKey string
		Chain        string
		BlockHeight  int64
		Dispatch     *DispatchResponse
		Nodes        *xsync.MapOf[string, *NodeSession]
	}
	type args struct {
		nodeSession *NodeSession
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantResult     *RPCSessionResult
		wantStatusCode int
		wantErr        bool
	}{
		// TODO: Add test cases.
	}
	servicerMap.Store("xxxx")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create a fake server and assign to servicerClient
			servicerClient = NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				assert.Equal(t, req.URL.String(), "http://example.com/some/path")
				return &http.Response{
					StatusCode: 200,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(`OK`)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			s := &Session{
				Hash:         tt.fields.Hash,
				AppPublicKey: tt.fields.AppPublicKey,
				Chain:        tt.fields.Chain,
				BlockHeight:  tt.fields.BlockHeight,
				Dispatch:     tt.fields.Dispatch,
				Nodes:        tt.fields.Nodes,
			}
			gotResult, gotStatusCode, err := s.GetDispatch(tt.args.nodeSession)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDispatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("GetDispatch() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
			if gotStatusCode != tt.wantStatusCode {
				t.Errorf("GetDispatch() gotStatusCode = %v, want %v", gotStatusCode, tt.wantStatusCode)
			}
		})
	}
}

func TestSession_GetNodeSessionByPubKey(t *testing.T) {
	type fields struct {
		Hash         string
		AppPublicKey string
		Chain        string
		BlockHeight  int64
		Dispatch     *DispatchResponse
		Nodes        *xsync.MapOf[string, *NodeSession]
	}
	type args struct {
		servicerPubKey string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *NodeSession
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				Hash:         tt.fields.Hash,
				AppPublicKey: tt.fields.AppPublicKey,
				Chain:        tt.fields.Chain,
				BlockHeight:  tt.fields.BlockHeight,
				Dispatch:     tt.fields.Dispatch,
				Nodes:        tt.fields.Nodes,
			}
			got, err := s.GetNodeSessionByPubKey(tt.args.servicerPubKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNodeSessionByPubKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNodeSessionByPubKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_InvalidateNodeSession(t *testing.T) {
	type fields struct {
		Hash         string
		AppPublicKey string
		Chain        string
		BlockHeight  int64
		Dispatch     *DispatchResponse
		Nodes        *xsync.MapOf[string, *NodeSession]
	}
	type args struct {
		servicerPubKey string
		e              *SdkErrorResponse
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *SdkErrorResponse
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				Hash:         tt.fields.Hash,
				AppPublicKey: tt.fields.AppPublicKey,
				Chain:        tt.fields.Chain,
				BlockHeight:  tt.fields.BlockHeight,
				Dispatch:     tt.fields.Dispatch,
				Nodes:        tt.fields.Nodes,
			}
			if got := s.InvalidateNodeSession(tt.args.servicerPubKey, tt.args.e); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InvalidateNodeSession() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_NewNodeFromRelay(t *testing.T) {
	type fields struct {
		Hash         string
		AppPublicKey string
		Chain        string
		BlockHeight  int64
		Dispatch     *DispatchResponse
		Nodes        *xsync.MapOf[string, *NodeSession]
	}
	type args struct {
		relay *pocketTypes.Relay
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *NodeSession
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				Hash:         tt.fields.Hash,
				AppPublicKey: tt.fields.AppPublicKey,
				Chain:        tt.fields.Chain,
				BlockHeight:  tt.fields.BlockHeight,
				Dispatch:     tt.fields.Dispatch,
				Nodes:        tt.fields.Nodes,
			}
			if got := s.NewNodeFromRelay(tt.args.relay); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNodeFromRelay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_ValidateSessionTask(t *testing.T) {
	type fields struct {
		Hash         string
		AppPublicKey string
		Chain        string
		BlockHeight  int64
		Dispatch     *DispatchResponse
		Nodes        *xsync.MapOf[string, *NodeSession]
	}
	type args struct {
		servicerPubKey string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   func()
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				Hash:         tt.fields.Hash,
				AppPublicKey: tt.fields.AppPublicKey,
				Chain:        tt.fields.Chain,
				BlockHeight:  tt.fields.BlockHeight,
				Dispatch:     tt.fields.Dispatch,
				Nodes:        tt.fields.Nodes,
			}
			if got := s.ValidateSessionTask(tt.args.servicerPubKey); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateSessionTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cleanOldSessions(t *testing.T) {
	type args struct {
		c *cron.Cron
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanOldSessions(tt.args.c)
		})
	}
}
