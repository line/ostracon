package v2

import (
	"fmt"
	"net"
	"os"
	"sort"
	"sync"
	"testing"
	"time"

	abci "github.com/line/ostracon/abci/types"
	"github.com/line/ostracon/behaviour"
	bc "github.com/line/ostracon/blockchain"
	cfg "github.com/line/ostracon/config"
	"github.com/line/ostracon/libs/log"
	"github.com/line/ostracon/libs/service"
	"github.com/line/ostracon/mempool/mock"
	"github.com/line/ostracon/p2p"
	"github.com/line/ostracon/p2p/conn"
	bcproto "github.com/line/ostracon/proto/ostracon/blockchain"
	"github.com/line/ostracon/proxy"
	sm "github.com/line/ostracon/state"
	"github.com/line/ostracon/store"
	"github.com/line/ostracon/types"
	tmtime "github.com/line/ostracon/types/time"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	dbm "github.com/tendermint/tm-db"
)

type mockPeer struct {
	service.Service
	id p2p.ID
}

func (mp mockPeer) FlushStop()           {}
func (mp mockPeer) ID() p2p.ID           { return mp.id }
func (mp mockPeer) RemoteIP() net.IP     { return net.IP{} }
func (mp mockPeer) RemoteAddr() net.Addr { return &net.TCPAddr{IP: mp.RemoteIP(), Port: 8800} }

func (mp mockPeer) IsOutbound() bool   { return true }
func (mp mockPeer) IsPersistent() bool { return true }
func (mp mockPeer) CloseConn() error   { return nil }

func (mp mockPeer) NodeInfo() p2p.NodeInfo {
	return p2p.DefaultNodeInfo{
		DefaultNodeID: "",
		ListenAddr:    "",
	}
}
func (mp mockPeer) Status() conn.ConnectionStatus { return conn.ConnectionStatus{} }
func (mp mockPeer) SocketAddr() *p2p.NetAddress   { return &p2p.NetAddress{} }

func (mp mockPeer) Send(byte, []byte) bool    { return true }
func (mp mockPeer) TrySend(byte, []byte) bool { return true }

func (mp mockPeer) Set(string, interface{}) {}
func (mp mockPeer) Get(string) interface{}  { return struct{}{} }

func (mp mockPeer) String() string { return fmt.Sprintf("%v", mp.id) }

// nolint:unused // ignore
type mockBlockStore struct {
	blocks map[int64]*types.Block
}

// nolint:unused // ignore
func (ml *mockBlockStore) Height() int64 {
	return int64(len(ml.blocks))
}

// nolint:unused // ignore
func (ml *mockBlockStore) LoadBlock(height int64) *types.Block {
	return ml.blocks[height]
}

// nolint:unused // ignore
func (ml *mockBlockStore) SaveBlock(block *types.Block, part *types.PartSet, commit *types.Commit) {
	ml.blocks[block.Height] = block
}

type mockBlockApplier struct {
}

// XXX: Add whitelist/blacklist?
func (mba *mockBlockApplier) ApplyBlock(
	state sm.State, blockID types.BlockID, block *types.Block, times *sm.CommitStepTimes,
) (sm.State, int64, error) {
	state.LastBlockHeight++
	return state, 0, nil
}

type mockSwitchIo struct {
	mtx                 sync.Mutex
	switchedToConsensus bool
	numStatusResponse   int
	numBlockResponse    int
	numNoBlockResponse  int
}

func (sio *mockSwitchIo) sendBlockRequest(peerID p2p.ID, height int64) error {
	return nil
}

func (sio *mockSwitchIo) sendStatusResponse(base, height int64, peerID p2p.ID) error {
	sio.mtx.Lock()
	defer sio.mtx.Unlock()
	sio.numStatusResponse++
	return nil
}

func (sio *mockSwitchIo) sendBlockToPeer(block *types.Block, peerID p2p.ID) error {
	sio.mtx.Lock()
	defer sio.mtx.Unlock()
	sio.numBlockResponse++
	return nil
}

func (sio *mockSwitchIo) sendBlockNotFound(height int64, peerID p2p.ID) error {
	sio.mtx.Lock()
	defer sio.mtx.Unlock()
	sio.numNoBlockResponse++
	return nil
}

func (sio *mockSwitchIo) trySwitchToConsensus(state sm.State, skipWAL bool) bool {
	sio.mtx.Lock()
	defer sio.mtx.Unlock()
	sio.switchedToConsensus = true
	return true
}

func (sio *mockSwitchIo) broadcastStatusRequest() error {
	return nil
}

type testReactorParams struct {
	logger      log.Logger
	genDoc      *types.GenesisDoc
	privVals    []types.PrivValidator
	startHeight int64
	mockA       bool
}

func newTestReactor(p testReactorParams) *BlockchainReactor {
	store, state, _ := newReactorStore(p.genDoc, p.privVals, p.startHeight)
	reporter := behaviour.NewMockReporter()

	var appl blockApplier

	if p.mockA {
		appl = &mockBlockApplier{}
	} else {
		app := &testApp{}
		cc := proxy.NewLocalClientCreator(app)
		proxyApp := proxy.NewAppConns(cc)
		err := proxyApp.Start()
		if err != nil {
			panic(fmt.Errorf("error start app: %w", err))
		}
		db := dbm.NewMemDB()
		stateStore := sm.NewStore(db)
		appl = sm.NewBlockExecutor(stateStore, p.logger, proxyApp.Consensus(), mock.Mempool{}, sm.EmptyEvidencePool{})
		if err = stateStore.Save(state); err != nil {
			panic(err)
		}
	}

	r := newReactor(state, store, reporter, appl, true)
	logger := log.TestingLogger()
	r.SetLogger(logger.With("module", "blockchain"))

	return r
}

// This test is left here and not deleted to retain the termination cases for
// future improvement in [#4482](https://github.com/tendermint/tendermint/issues/4482).
// func TestReactorTerminationScenarios(t *testing.T) {

// 	config := cfg.ResetTestRoot("blockchain_reactor_v2_test")
// 	defer os.RemoveAll(config.RootDir)
// 	genDoc, privVals := randGenesisDoc(config.ChainID(), 1, false, 30)
// 	refStore, _, _ := newReactorStore(genDoc, privVals, 20)

// 	params := testReactorParams{
// 		logger:      log.TestingLogger(),
// 		genDoc:      genDoc,
// 		privVals:    privVals,
// 		startHeight: 10,
// 		bufferSize:  100,
// 		mockA:       true,
// 	}

// 	type testEvent struct {
// 		evType string
// 		peer   string
// 		height int64
// 	}

// 	tests := []struct {
// 		name   string
// 		params testReactorParams
// 		msgs   []testEvent
// 	}{
// 		{
// 			name:   "simple termination on max peer height - one peer",
// 			params: params,
// 			msgs: []testEvent{
// 				{evType: "AddPeer", peer: "P1"},
// 				{evType: "ReceiveS", peer: "P1", height: 13},
// 				{evType: "BlockReq"},
// 				{evType: "ReceiveB", peer: "P1", height: 11},
// 				{evType: "BlockReq"},
// 				{evType: "BlockReq"},
// 				{evType: "ReceiveB", peer: "P1", height: 12},
// 				{evType: "Process"},
// 				{evType: "ReceiveB", peer: "P1", height: 13},
// 				{evType: "Process"},
// 			},
// 		},
// 		{
// 			name:   "simple termination on max peer height - two peers",
// 			params: params,
// 			msgs: []testEvent{
// 				{evType: "AddPeer", peer: "P1"},
// 				{evType: "AddPeer", peer: "P2"},
// 				{evType: "ReceiveS", peer: "P1", height: 13},
// 				{evType: "ReceiveS", peer: "P2", height: 15},
// 				{evType: "BlockReq"},
// 				{evType: "BlockReq"},
// 				{evType: "ReceiveB", peer: "P1", height: 11},
// 				{evType: "ReceiveB", peer: "P2", height: 12},
// 				{evType: "Process"},
// 				{evType: "BlockReq"},
// 				{evType: "BlockReq"},
// 				{evType: "ReceiveB", peer: "P1", height: 13},
// 				{evType: "Process"},
// 				{evType: "ReceiveB", peer: "P2", height: 14},
// 				{evType: "Process"},
// 				{evType: "BlockReq"},
// 				{evType: "ReceiveB", peer: "P2", height: 15},
// 				{evType: "Process"},
// 			},
// 		},
// 		{
// 			name:   "termination on max peer height - two peers, noBlock error",
// 			params: params,
// 			msgs: []testEvent{
// 				{evType: "AddPeer", peer: "P1"},
// 				{evType: "AddPeer", peer: "P2"},
// 				{evType: "ReceiveS", peer: "P1", height: 13},
// 				{evType: "ReceiveS", peer: "P2", height: 15},
// 				{evType: "BlockReq"},
// 				{evType: "BlockReq"},
// 				{evType: "ReceiveNB", peer: "P1", height: 11},
// 				{evType: "BlockReq"},
// 				{evType: "ReceiveB", peer: "P2", height: 12},
// 				{evType: "ReceiveB", peer: "P2", height: 11},
// 				{evType: "Process"},
// 				{evType: "BlockReq"},
// 				{evType: "BlockReq"},
// 				{evType: "ReceiveB", peer: "P2", height: 13},
// 				{evType: "Process"},
// 				{evType: "ReceiveB", peer: "P2", height: 14},
// 				{evType: "Process"},
// 				{evType: "BlockReq"},
// 				{evType: "ReceiveB", peer: "P2", height: 15},
// 				{evType: "Process"},
// 			},
// 		},
// 		{
// 			name:   "termination on max peer height - two peers, remove one peer",
// 			params: params,
// 			msgs: []testEvent{
// 				{evType: "AddPeer", peer: "P1"},
// 				{evType: "AddPeer", peer: "P2"},
// 				{evType: "ReceiveS", peer: "P1", height: 13},
// 				{evType: "ReceiveS", peer: "P2", height: 15},
// 				{evType: "BlockReq"},
// 				{evType: "BlockReq"},
// 				{evType: "RemovePeer", peer: "P1"},
// 				{evType: "BlockReq"},
// 				{evType: "ReceiveB", peer: "P2", height: 12},
// 				{evType: "ReceiveB", peer: "P2", height: 11},
// 				{evType: "Process"},
// 				{evType: "BlockReq"},
// 				{evType: "BlockReq"},
// 				{evType: "ReceiveB", peer: "P2", height: 13},
// 				{evType: "Process"},
// 				{evType: "ReceiveB", peer: "P2", height: 14},
// 				{evType: "Process"},
// 				{evType: "BlockReq"},
// 				{evType: "ReceiveB", peer: "P2", height: 15},
// 				{evType: "Process"},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			reactor := newTestReactor(params)
// 			reactor.Start()
// 			reactor.reporter = behaviour.NewMockReporter()
// 			mockSwitch := &mockSwitchIo{switchedToConsensus: false}
// 			reactor.io = mockSwitch
// 			// time for go routines to start
// 			time.Sleep(time.Millisecond)

// 			for _, step := range tt.msgs {
// 				switch step.evType {
// 				case "AddPeer":
// 					reactor.scheduler.send(bcAddNewPeer{peerID: p2p.ID(step.peer)})
// 				case "RemovePeer":
// 					reactor.scheduler.send(bcRemovePeer{peerID: p2p.ID(step.peer)})
// 				case "ReceiveS":
// 					reactor.scheduler.send(bcStatusResponse{
// 						peerID: p2p.ID(step.peer),
// 						height: step.height,
// 						time:   time.Now(),
// 					})
// 				case "ReceiveB":
// 					reactor.scheduler.send(bcBlockResponse{
// 						peerID: p2p.ID(step.peer),
// 						block:  refStore.LoadBlock(step.height),
// 						size:   10,
// 						time:   time.Now(),
// 					})
// 				case "ReceiveNB":
// 					reactor.scheduler.send(bcNoBlockResponse{
// 						peerID: p2p.ID(step.peer),
// 						height: step.height,
// 						time:   time.Now(),
// 					})
// 				case "BlockReq":
// 					reactor.scheduler.send(rTrySchedule{time: time.Now()})
// 				case "Process":
// 					reactor.processor.send(rProcessBlock{})
// 				}
// 				// give time for messages to propagate between routines
// 				time.Sleep(time.Millisecond)
// 			}

// 			// time for processor to finish and reactor to switch to consensus
// 			time.Sleep(20 * time.Millisecond)
// 			assert.True(t, mockSwitch.hasSwitchedToConsensus())
// 			reactor.Stop()
// 		})
// 	}
// }

func TestReactorHelperMode(t *testing.T) {
	var (
		channelID = byte(0x40)
	)
	config := cfg.ResetTestRoot("blockchain_reactor_v2_test")
	defer os.RemoveAll(config.RootDir)
	genDoc, privVals := randGenesisDoc(config.ChainID(), 1, false, 30)

	params := testReactorParams{
		logger:      log.TestingLogger(),
		genDoc:      genDoc,
		privVals:    privVals,
		startHeight: 20,
		mockA:       true,
	}

	type testEvent struct {
		peer  string
		event interface{}
	}

	tests := []struct {
		name   string
		params testReactorParams
		msgs   []testEvent
	}{
		{
			name:   "status request",
			params: params,
			msgs: []testEvent{
				{"P1", bcproto.StatusRequest{}},
				{"P1", bcproto.BlockRequest{Height: 13}},
				{"P1", bcproto.BlockRequest{Height: 20}},
				{"P1", bcproto.BlockRequest{Height: 22}},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			reactor := newTestReactor(params)
			mockSwitch := &mockSwitchIo{switchedToConsensus: false}
			reactor.io = mockSwitch
			err := reactor.Start()
			require.NoError(t, err)

			for i := 0; i < len(tt.msgs); i++ {
				step := tt.msgs[i]
				switch ev := step.event.(type) {
				case bcproto.StatusRequest:
					old := mockSwitch.numStatusResponse
					msg, err := bc.EncodeMsg(&ev)
					assert.NoError(t, err)
					reactor.Receive(channelID, mockPeer{id: p2p.ID(step.peer)}, msg)
					assert.Equal(t, old+1, mockSwitch.numStatusResponse)
				case bcproto.BlockRequest:
					if ev.Height > params.startHeight {
						old := mockSwitch.numNoBlockResponse
						msg, err := bc.EncodeMsg(&ev)
						assert.NoError(t, err)
						reactor.Receive(channelID, mockPeer{id: p2p.ID(step.peer)}, msg)
						assert.Equal(t, old+1, mockSwitch.numNoBlockResponse)
					} else {
						old := mockSwitch.numBlockResponse
						msg, err := bc.EncodeMsg(&ev)
						assert.NoError(t, err)
						assert.NoError(t, err)
						reactor.Receive(channelID, mockPeer{id: p2p.ID(step.peer)}, msg)
						assert.Equal(t, old+1, mockSwitch.numBlockResponse)
					}
				}
			}
			err = reactor.Stop()
			require.NoError(t, err)
		})
	}
}

func TestReactorSetSwitchNil(t *testing.T) {
	config := cfg.ResetTestRoot("blockchain_reactor_v2_test")
	defer os.RemoveAll(config.RootDir)
	genDoc, privVals := randGenesisDoc(config.ChainID(), 1, false, 30)

	reactor := newTestReactor(testReactorParams{
		logger:   log.TestingLogger(),
		genDoc:   genDoc,
		privVals: privVals,
	})
	reactor.SetSwitch(nil)

	assert.Nil(t, reactor.Switch)
	assert.Nil(t, reactor.io)
}

//----------------------------------------------
// utility funcs

func makeTxs(height int64) (txs []types.Tx) {
	for i := 0; i < 10; i++ {
		txs = append(txs, types.Tx([]byte{byte(height), byte(i)}))
	}
	return txs
}

func makeBlock(privVal types.PrivValidator, height int64, state sm.State, lastCommit *types.Commit) *types.Block {
	message := state.MakeHashMessage(0)
	proof, _ := privVal.GenerateVRFProof(message)
	proposerAddr := state.Validators.SelectProposer(state.LastProofHash, height, 0).Address
	block, _ := state.MakeBlock(height, makeTxs(height), lastCommit, nil, proposerAddr, 0, proof)
	return block
}

type testApp struct {
	abci.BaseApplication
}

func randGenesisDoc(chainID string, numValidators int, randPower bool, minPower int64) (
	*types.GenesisDoc, []types.PrivValidator) {
	validators := make([]types.GenesisValidator, numValidators)
	privValidators := make([]types.PrivValidator, numValidators)
	for i := 0; i < numValidators; i++ {
		val, privVal := types.RandValidator(randPower, minPower)
		validators[i] = types.GenesisValidator{
			PubKey: val.PubKey,
			Power:  val.VotingPower,
		}
		privValidators[i] = privVal
	}
	sort.Sort(types.PrivValidatorsByAddress(privValidators))

	return &types.GenesisDoc{
		GenesisTime: tmtime.Now(),
		ChainID:     chainID,
		Validators:  validators,
	}, privValidators
}

// Why are we importing the entire blockExecutor dependency graph here
// when we have the facilities to
func newReactorStore(
	genDoc *types.GenesisDoc,
	privVals []types.PrivValidator,
	maxBlockHeight int64) (*store.BlockStore, sm.State, *sm.BlockExecutor) {
	if len(privVals) != 1 {
		panic("only support one validator")
	}
	app := &testApp{}
	cc := proxy.NewLocalClientCreator(app)
	proxyApp := proxy.NewAppConns(cc)
	err := proxyApp.Start()
	if err != nil {
		panic(fmt.Errorf("error start app: %w", err))
	}

	stateDB := dbm.NewMemDB()
	blockStore := store.NewBlockStore(dbm.NewMemDB())
	stateStore := sm.NewStore(stateDB)
	state, err := stateStore.LoadFromDBOrGenesisDoc(genDoc)
	if err != nil {
		panic(fmt.Errorf("error constructing state from genesis file: %w", err))
	}

	db := dbm.NewMemDB()
	stateStore = sm.NewStore(db)
	blockExec := sm.NewBlockExecutor(stateStore, log.TestingLogger(), proxyApp.Consensus(),
		mock.Mempool{}, sm.EmptyEvidencePool{})
	if err = stateStore.Save(state); err != nil {
		panic(err)
	}

	// add blocks in
	for blockHeight := int64(1); blockHeight <= maxBlockHeight; blockHeight++ {
		lastCommit := types.NewCommit(blockHeight-1, 0, types.BlockID{}, nil)
		if blockHeight > 1 {
			lastBlockMeta, err := blockStore.LoadBlockMeta(blockHeight - 1)
			if err != nil {
				panic(err)
			}
			lastBlock, err := blockStore.LoadBlock(blockHeight - 1)
			if err != nil {
				panic(err)
			}
			vote, err := types.MakeVote(
				lastBlock.Header.Height,
				lastBlockMeta.BlockID,
				state.Validators,
				privVals[0],
				lastBlock.Header.ChainID,
				time.Now(),
			)
			if err != nil {
				panic(err)
			}
			lastCommit = types.NewCommit(vote.Height, vote.Round,
				lastBlockMeta.BlockID, []types.CommitSig{vote.CommitSig()})
		}

		thisBlock := makeBlock(privVals[0], blockHeight, state, lastCommit)

		thisParts := thisBlock.MakePartSet(types.BlockPartSizeBytes)
		blockID := types.BlockID{Hash: thisBlock.Hash(), PartSetHeader: thisParts.Header()}

		state, _, err = blockExec.ApplyBlock(state, blockID, thisBlock, nil)
		if err != nil {
			panic(fmt.Errorf("error apply block: %w", err))
		}

		blockStore.SaveBlock(thisBlock, thisParts, lastCommit)
	}
	return blockStore, state, blockExec
}
