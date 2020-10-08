package protocols

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"go.dedis.ch/kyber/v3/suites"
	"golang.org/x/xerrors"
	"reflect"
	"sync"
	"testing"
	"time"

	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/log"
	"go.dedis.ch/onet/v3/network"
)

type propagateMsg struct {
	Data []byte
}

func init() {
	network.RegisterMessage(propagateMsg{})
}

func TestPropagation(t *testing.T) {
	propagate(t,
		[]int{3, 10, 14, 4, 8, 8},
		[]int{0, 0, 0, 1, 3, 6})
}

var tSuite = suites.MustFind("Ed25519")

// Tests an n-node system
func propagate(t *testing.T, nbrNodes, nbrFailures []int) {
	for i, n := range nbrNodes {
		local := onet.NewLocalTest(tSuite)
		servers, el, _ := local.GenTree(n, true)
		var recvCount int
		var iMut sync.Mutex
		msg := &propagateMsg{[]byte("propagate")}
		propFuncs := make([]PropagationFunc, n)

		// setup the servers
		var err error
		for n, server := range servers {
			pc := &PC{server, local.Overlays[server.ServerIdentity.ID]}
			propFuncs[n], err = NewPropagationFuncTest(pc, "Propagate",
				nbrFailures[i],
				func(m network.Message) error {
					if bytes.Equal(msg.Data, m.(*propagateMsg).Data) {
						iMut.Lock()
						recvCount++
						iMut.Unlock()
						return nil
					}

					t.Error("Didn't receive correct data")
					return xerrors.New("Didn't receive correct data")
				},
				func() network.Message {
					return &propagateMsg{Data: []byte{1, 2, 3}}
				})
			require.NoError(t, err)
		}

		// shut down some servers to simulate failure
		for k := 0; k < nbrFailures[i]; k++ {
			err = servers[len(servers)-1-k].Close()
			require.NoError(t, err)
		}

		// start the propagation
		log.Lvl2("Starting to propagate", reflect.TypeOf(msg))
		datas, err := propFuncs[0](el, msg,
			1*time.Second)
		require.NoError(t, err)
		require.Equal(t, n, recvCount+nbrFailures[i], "Didn't get data-request")
		require.Equal(t, n-1, len(datas)+nbrFailures[i], "Not all nodes replied")

		local.CloseAll()
		log.AfterTest(t)
	}
}

type PC struct {
	C *onet.Server
	O *onet.Overlay
}

func (pc *PC) ProtocolRegister(name string, protocol onet.NewProtocol) (onet.ProtocolID, error) {
	return pc.C.ProtocolRegister(name, protocol)
}
func (pc *PC) ServerIdentity() *network.ServerIdentity {
	return pc.C.ServerIdentity

}
func (pc *PC) CreateProtocol(name string, t *onet.Tree) (onet.ProtocolInstance, error) {
	return pc.O.CreateProtocol(name, t, onet.NilServiceID)
}
