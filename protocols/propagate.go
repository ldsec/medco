package protocols

import (
	"golang.org/x/xerrors"
	"reflect"
	"strings"
	"sync"
	"time"

	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/log"
	"go.dedis.ch/onet/v3/network"
)

func init() {
	network.RegisterMessage(PropagateSendData{})
	network.RegisterMessage(PropagateReply{})
}

// How long to wait before timing out on waiting for the time-out.
const initialWait = 100000 * time.Millisecond

// Propagate is a protocol that sends some data to all attached nodes
// and waits for confirmation before returning.
type Propagate struct {
	*onet.TreeNodeInstance
	onDataToChildren PropagationOneMsg
	onDataToRoot     PropagationOneMsgSend
	onDoneCb         PropagationMultiMsg
	sd               *PropagateSendData
	ChannelSD        chan struct {
		*onet.TreeNode
		PropagateSendData
	}
	ChannelReply chan struct {
		*onet.TreeNode
		PropagateReply
	}

	allowedFailures int
	sync.Mutex
	closing chan bool
}

// PropagateSendData is the message to pass the data to the children
type PropagateSendData struct {
	// Data is the data to transmit to the children
	Data []byte
	// How long the root will wait for the children before
	// timing out.
	Timeout time.Duration
}

// PropagateReply is sent from the children back to the root
type PropagateReply struct {
	// Data is the data to transmit to the root
	Data []byte
	// Level is how many children replied
	Level int
}

// PropagationFunc starts the propagation protocol and blocks until all children
// minus the exception stored the new value or the timeout has been reached.
// The return value is the number of nodes that acknowledged having
// stored the new value or an error if the protocol couldn't start.
type PropagationFunc func(el *onet.Roster, msg network.Message,
	timeout time.Duration) ([]network.Message, error)

// PropagationOneMsg is the function that will store the new data.
type PropagationOneMsg func(network.Message) error

// PropagationOneMsgSend is the function that will store the new data.
type PropagationOneMsgSend func() network.Message

// PropagationMultiMsg is the function that will store the new data.
type PropagationMultiMsg func([]network.Message)

// propagationContext is used for testing.
type propagationContext interface {
	ProtocolRegister(name string, protocol onet.NewProtocol) (onet.ProtocolID, error)
	ServerIdentity() *network.ServerIdentity
	CreateProtocol(name string, t *onet.Tree) (onet.ProtocolInstance, error)
}

// NewPropagationProtocol creates a new protocl for propagation.
func NewPropagationProtocol(n *onet.TreeNodeInstance) (onet.ProtocolInstance, error) {
	p := &Propagate{
		sd:               &PropagateSendData{[]byte{}, initialWait},
		TreeNodeInstance: n,
		closing:          make(chan bool),
		allowedFailures:  (len(n.Roster().List) - 1) / 3,
	}
	for _, h := range []interface{}{&p.ChannelSD, &p.ChannelReply} {
		if err := p.RegisterChannel(h); err != nil {
			return nil, err
		}
	}
	return p, nil
}

// NewPropagationFunc registers a new protocol name with the context c and will
// set f as handler for every new instance of that protocol.
// The protocol will fail if more than thresh nodes per subtree fail to respond.
// If thresh == -1, the threshold defaults to len(n.Roster().List-1)/3. Thus, for a roster of
// 5, t = int(4/3) = 1, e.g. 1 node out of the 5 can fail.
func NewPropagationFunc(c propagationContext, name string,
	thresh int) (PropagationFunc, error) {
	return NewPropagationFuncTest(c, name, thresh, nil, nil)
}

// NewPropagationFuncTest takes two callbacks for easier testing without
// a `service.NewProtocl`
func NewPropagationFuncTest(c propagationContext, name string, thresh int,
	onDataToChildren PropagationOneMsg,
	onDataToRoot PropagationOneMsgSend) (PropagationFunc, error) {
	pid, err := c.ProtocolRegister(name, func(n *onet.TreeNodeInstance) (onet.
		ProtocolInstance, error) {
		pi, err := NewPropagationProtocol(n)
		if err != nil {
			return nil, xerrors.Errorf("couldn't create protocol: %+v", err)
		}
		proto := pi.(*Propagate)
		proto.onDataToChildren = onDataToChildren
		proto.onDataToRoot = onDataToRoot
		return pi, err
	})

	log.Lvl3("Registering new propagation for", c.ServerIdentity(),
		name, pid)
	return func(el *onet.Roster, msg network.Message,
		to time.Duration) ([]network.Message, error) {
		rooted := el.NewRosterWithRoot(c.ServerIdentity())
		if rooted == nil {
			return nil, xerrors.New("we're not in the roster")
		}
		tree := rooted.GenerateNaryTree(len(el.List))
		if tree == nil {
			return nil, xerrors.New("Didn't find root in tree")
		}
		log.Lvl3(el.List[0].Address, "Starting to propagate", reflect.TypeOf(msg))
		pi, err := c.CreateProtocol(name, tree)
		if err != nil {
			return nil, err
		}
		proto := pi.(*Propagate)
		proto.Lock()
		t := thresh
		if t == -1 {
			t = (len(el.List) - 1) / 3
		}
		proto.allowedFailures = t

		if msg != nil {
			d, err := network.Marshal(msg)
			if err != nil {
				proto.Unlock()
				return nil, err
			}
			proto.sd.Data = d
		}
		proto.sd.Timeout = to

		done := make(chan []network.Message)
		proto.onDoneCb = func(msg []network.Message) {
			done <- msg
		}
		proto.Unlock()

		if err := proto.Start(); err != nil {
			return nil, err
		}
		select {
		case replies := <-done:
			return replies, nil
		case <-proto.closing:
			return nil, nil
		}
	}, err
}

// Start will contact everyone and make the connections
func (p *Propagate) Start() error {
	log.Lvl4("going to contact", p.Root().ServerIdentity)
	return p.SendTo(p.Root(), p.sd)
}

// Dispatch can handle timeouts
func (p *Propagate) Dispatch() error {
	process := true
	var received int
	var rcvMsgs [][]byte
	log.Lvl4(p.ServerIdentity(), "Start dispatch")
	defer p.Done()
	defer func() {
		if p.IsRoot() {
			if p.onDoneCb != nil {
				var rcvNetMsgs []network.Message
				for _, data := range rcvMsgs {
					if data != nil && len(data) > 0 {
						_, netMsg, err := network.Unmarshal(data, p.Suite())
						if err != nil {
							log.Warnf("Got error while unmarshaling: %+v", err)
						} else {
							rcvNetMsgs = append(rcvNetMsgs, netMsg)
						}
					}
				}
				p.onDoneCb(rcvNetMsgs)
			}
		}
	}()

	var gotSendData bool
	var errs []error
	subtreeCount := p.TreeNode().SubtreeCount()

	for process {
		p.Lock()
		timeout := p.sd.Timeout
		log.Lvl4("Got timeout", timeout, "from SendData")
		p.Unlock()
		select {
		case msg := <-p.ChannelSD:
			if gotSendData {
				log.Error("already got msg")
				continue
			}
			gotSendData = true
			log.Lvl3(p.ServerIdentity(), "Got data from", msg.ServerIdentity, "and setting timeout to", msg.Timeout)
			p.sd.Timeout = msg.Timeout
			if p.onDataToChildren != nil {
				_, netMsg, err := network.Unmarshal(msg.Data, p.Suite())
				if err != nil {
					log.Lvlf2("Unmarshal failed with %v", err)
				} else {
					err := p.onDataToChildren(netMsg)
					if err != nil {
						log.Lvlf2("Propagation callback failed: %v", err)
					}
				}
			}
			if !p.IsRoot() {
				log.Lvl3(p.ServerIdentity(), "Sending to parent")
				var data []byte
				if p.onDataToRoot != nil {
					var err error
					data, err = network.Marshal(p.onDataToRoot())
					if err != nil {
						return xerrors.Errorf("couldn't marshal message: %+v",
							err)
					}
				}
				if err := p.SendToParent(
					&PropagateReply{Data: data}); err != nil {
					return err
				}
			}
			if p.IsLeaf() {
				process = false
			} else {
				log.Lvl3(p.ServerIdentity(), "Sending to children")
				if errs = p.SendToChildrenInParallel(&msg.PropagateSendData); len(errs) != 0 {
					var errsStr []string
					for _, e := range errs {
						errsStr = append(errsStr, e.Error())
					}
					if len(errs) > p.allowedFailures {
						return xerrors.New(strings.Join(errsStr, "\n"))
					}
					log.Lvl2("Error while sending to children:", errsStr)
				}
			}
		case rep := <-p.ChannelReply:
			if !gotSendData {
				log.Error("got response before send")
				continue
			}
			received++
			log.Lvl4(p.ServerIdentity(), "received:", received, subtreeCount)
			if !p.IsRoot() {
				if err := p.SendToParent(&PropagateReply{
					Data: rep.Data}); err != nil {
					return err
				}
			} else {
				rcvMsgs = append(rcvMsgs, rep.Data)
			}
			// Only wait for the number of children that successfully received our message.
			if received == subtreeCount-len(errs) && received >= subtreeCount-p.allowedFailures {
				process = false
			}
		case <-time.After(timeout):
			if received < subtreeCount-p.allowedFailures {
				_, _, err := network.Unmarshal(p.sd.Data, p.Suite())
				return xerrors.Errorf("Timeout of %s reached, "+
					"got %v but need %v, err: %+v",
					timeout, received, subtreeCount-p.allowedFailures, err)
			}
			process = false
		case <-p.closing:
			process = false
			p.onDoneCb = nil
		}
	}
	log.Lvl3(p.ServerIdentity(), "done, isroot:", p.IsRoot())
	return nil
}

// RegisterOnDone takes a function that will be called once the data has been
// sent to the whole tree. It receives the number of nodes that replied
// successfully to the propagation.
func (p *Propagate) RegisterOnDone(fn PropagationMultiMsg) {
	p.onDoneCb = fn
}

// RegisterOnDataToChildren takes a function that will be called for that node if it
// needs to update its data.
func (p *Propagate) RegisterOnDataToChildren(fn PropagationOneMsg) {
	p.onDataToChildren = fn
}

// RegisterOnDataToRoot takes a function that will be called for that node if it
// needs to update its data.
func (p *Propagate) RegisterOnDataToRoot(fn PropagationOneMsgSend) {
	p.onDataToRoot = fn
}

// Config stores the basic configuration for that protocol.
func (p *Propagate) Config(d []byte, timeout time.Duration) {
	p.sd.Data = d
	p.sd.Timeout = timeout
}

// Shutdown informs the Dispatch method to stop
// waiting.
func (p *Propagate) Shutdown() error {
	close(p.closing)
	return nil
}
