package mesh

import (
	"github.com/pokt-network/pocket-core/app"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/puzpuzpuz/xsync"
	"testing"
)

func init() {
	app.GlobalMeshConfig = app.DefaultMeshConfig("~/pocket")
	logger = initLogger()
}

func NewSessionStorage() *SessionStorage {
	return &SessionStorage{
		Sessions:         xsync.NewMapOf[*Session](),
		ValidationWorker: NewWorkerPool("test", "balanced", 1, 1, 10000),
	}
}

/*
Scenario 1:

RelaySessionBlockHeight = 9
servicerNodeSessionBlockHeight = 5
-> return true (node running behind)


RelaySessionBlockHeight = 10
servicerNodeSessionBlockHeight = 5
-> return false - (node running super behind) or do we want to allow for 2 blocks tolerance

RelaySessionBlockHeight = 201
servicerNodeSessionBlockHeight = 5
-> return false (malicious user case)

RelaySessionBlockHeight = 5
servicerNodeSessionBlockHeight = 9
-> should be fine, this code is not reached as we have session already cached.
*/
func TestSessionStorage_ShouldAssumeOptimisticSession(t *testing.T) {
	type fields struct {
		SessionStorage *SessionStorage
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
		{
			name: "node running 1 behind",
			fields: fields{
				SessionStorage: NewSessionStorage(),
			},
			args: args{
				relay: &pocketTypes.Relay{
					Proof: pocketTypes.RelayProof{
						SessionBlockHeight: 9,
					},
				},
				servicerNode: &fullNode{
					Status: &app.HealthResponse{
						Height: 8,
					},
					BlocksPerSession: 4,
				},
			},
			want: true,
		},
		{
			name: "node running far behind",
			fields: fields{
				SessionStorage: NewSessionStorage(),
			},
			args: args{
				relay: &pocketTypes.Relay{
					Proof: pocketTypes.RelayProof{
						SessionBlockHeight: 9,
					},
				},
				servicerNode: &fullNode{
					Status: &app.HealthResponse{
						Height: 5,
					},
					BlocksPerSession: 4,
				},
			},
			want: false,
		},
		{
			name: "node running super behind",
			fields: fields{
				SessionStorage: NewSessionStorage(),
			},
			args: args{
				relay: &pocketTypes.Relay{
					Proof: pocketTypes.RelayProof{
						SessionBlockHeight: 13,
					},
				},
				servicerNode: &fullNode{
					Status: &app.HealthResponse{
						Height: 5,
					},
					BlocksPerSession: 4,
				},
			},
			want: false,
		},
		{
			name: "malicious user case",
			fields: fields{
				SessionStorage: NewSessionStorage(),
			},
			args: args{
				relay: &pocketTypes.Relay{
					Proof: pocketTypes.RelayProof{
						SessionBlockHeight: 201,
					},
				},
				servicerNode: &fullNode{
					Status: &app.HealthResponse{
						Height: 5,
					},
					BlocksPerSession: 4,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.SessionStorage.ShouldAssumeOptimisticSession(tt.args.relay, tt.args.servicerNode); got != tt.want {
				t.Errorf("ShouldAssumeOptimisticSession() = %v, want %v", got, tt.want)
			}
		})
	}
}
