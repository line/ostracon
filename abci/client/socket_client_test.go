package abcicli_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	abcicli "github.com/line/ostracon/abci/client"
	"github.com/line/ostracon/abci/server"
	"github.com/line/ostracon/abci/types"
	tmrand "github.com/line/ostracon/libs/rand"
	"github.com/line/ostracon/libs/service"
)

func TestProperSyncCalls(t *testing.T) {
	app := slowApp{}

	s, c := setupClientServer(t, app)
	t.Cleanup(func() {
		if err := s.Stop(); err != nil {
			t.Error(err)
		}
	})
	t.Cleanup(func() {
		if err := c.Stop(); err != nil {
			t.Error(err)
		}
	})

	resp := make(chan error, 1)
	go func() {
		// This is BeginBlockSync unrolled....
		reqres := c.BeginBlockAsync(types.RequestBeginBlock{}, nil)
		_, err := c.FlushSync()
		require.NoError(t, err)
		res := reqres.Response.GetBeginBlock()
		require.NotNil(t, res)
		resp <- c.Error()
	}()

	select {
	case <-time.After(time.Second):
		require.Fail(t, "No response arrived")
	case err, ok := <-resp:
		require.True(t, ok, "Must not close channel")
		assert.NoError(t, err, "This should return success")
	}
}

func TestHangingSyncCalls(t *testing.T) {
	app := slowApp{}

	s, c := setupClientServer(t, app)
	t.Cleanup(func() {
		if err := s.Stop(); err != nil {
			t.Log(err)
		}
	})
	t.Cleanup(func() {
		if err := c.Stop(); err != nil {
			t.Log(err)
		}
	})

	resp := make(chan error, 1)
	go func() {
		// Start BeginBlock and flush it
		reqres := c.BeginBlockAsync(types.RequestBeginBlock{}, nil)
		flush := c.FlushAsync(nil)
		// wait 20 ms for all events to travel socket, but
		// no response yet from server
		time.Sleep(20 * time.Millisecond)
		// kill the server, so the connections break
		err := s.Stop()
		require.NoError(t, err)

		// wait for the response from BeginBlock
		reqres.Wait()
		flush.Wait()
		resp <- c.Error()
	}()

	select {
	case <-time.After(time.Second):
		require.Fail(t, "No response arrived")
	case err, ok := <-resp:
		require.True(t, ok, "Must not close channel")
		assert.Error(t, err, "We should get EOF error")
	}
}

func setupClientServer(t *testing.T, app types.Application) (
	service.Service, abcicli.Client) {
	// some port between 20k and 30k
	port := 20000 + tmrand.Int32()%10000
	addr := fmt.Sprintf("localhost:%d", port)

	s, err := server.NewServer(addr, "socket", app)
	require.NoError(t, err)
	err = s.Start()
	require.NoError(t, err)

	c := abcicli.NewSocketClient(addr, true)
	err = c.Start()
	require.NoError(t, err)

	return s, c
}

type slowApp struct {
	types.BaseApplication
}

func (slowApp) BeginBlock(req types.RequestBeginBlock) types.ResponseBeginBlock {
	time.Sleep(200 * time.Millisecond)
	return types.ResponseBeginBlock{}
}

func TestSockerClientCalls(t *testing.T) {
	app := sampleApp{}

	s, c := setupClientServer(t, app)
	t.Cleanup(func() {
		if err := s.Stop(); err != nil {
			t.Error(err)
		}
	})
	t.Cleanup(func() {
		if err := c.Stop(); err != nil {
			t.Error(err)
		}
	})

	gbCalled := false
	c.SetGlobalCallback(func(*types.Request, *types.Response) {
		gbCalled = true
	})

	cbCalled := false
	res := c.EchoAsync("msg", getResponseCallback(t, &cbCalled))
	res.Wait()
	require.True(t, cbCalled)
	require.True(t, gbCalled)

	cbCalled = false
	res = c.FlushAsync(getResponseCallback(t, &cbCalled))
	res.Wait()
	require.True(t, cbCalled)

	cbCalled = false
	res = c.InfoAsync(types.RequestInfo{}, getResponseCallback(t, &cbCalled))
	res.Wait()
	require.True(t, cbCalled)

	cbCalled = false
	res = c.SetOptionAsync(types.RequestSetOption{}, getResponseCallback(t, &cbCalled))
	res.Wait()
	require.True(t, cbCalled)

	cbCalled = false
	res = c.DeliverTxAsync(types.RequestDeliverTx{}, getResponseCallback(t, &cbCalled))
	res.Wait()
	require.True(t, cbCalled)

	cbCalled = false
	res = c.CheckTxAsync(types.RequestCheckTx{}, getResponseCallback(t, &cbCalled))
	res.Wait()
	require.True(t, cbCalled)

	cbCalled = false
	res = c.QueryAsync(types.RequestQuery{}, getResponseCallback(t, &cbCalled))
	res.Wait()
	require.True(t, cbCalled)

	cbCalled = false
	res = c.CommitAsync(getResponseCallback(t, &cbCalled))
	res.Wait()
	require.True(t, cbCalled)

	cbCalled = false
	res = c.InitChainAsync(types.RequestInitChain{}, getResponseCallback(t, &cbCalled))
	res.Wait()
	require.True(t, cbCalled)

	cbCalled = false
	res = c.BeginBlockAsync(types.RequestBeginBlock{}, getResponseCallback(t, &cbCalled))
	res.Wait()
	require.True(t, cbCalled)

	cbCalled = false
	res = c.EndBlockAsync(types.RequestEndBlock{}, getResponseCallback(t, &cbCalled))
	res.Wait()
	require.True(t, cbCalled)

	cbCalled = false
	res = c.BeginRecheckTxAsync(types.RequestBeginRecheckTx{}, getResponseCallback(t, &cbCalled))
	res.Wait()
	require.True(t, cbCalled)

	cbCalled = false
	res = c.EndRecheckTxAsync(types.RequestEndRecheckTx{}, getResponseCallback(t, &cbCalled))
	res.Wait()
	require.True(t, cbCalled)

	cbCalled = false
	res = c.ListSnapshotsAsync(types.RequestListSnapshots{}, getResponseCallback(t, &cbCalled))
	res.Wait()
	require.True(t, cbCalled)

	cbCalled = false
	res = c.OfferSnapshotAsync(types.RequestOfferSnapshot{}, getResponseCallback(t, &cbCalled))
	res.Wait()
	require.True(t, cbCalled)

	cbCalled = false
	res = c.LoadSnapshotChunkAsync(types.RequestLoadSnapshotChunk{}, getResponseCallback(t, &cbCalled))
	res.Wait()
	require.True(t, cbCalled)

	cbCalled = false
	res = c.ApplySnapshotChunkAsync(types.RequestApplySnapshotChunk{}, getResponseCallback(t, &cbCalled))
	res.Wait()
	require.True(t, cbCalled)

	_, err := c.EchoSync("msg")
	require.NoError(t, err)

	_, err = c.FlushSync()
	require.NoError(t, err)

	_, err = c.InfoSync(types.RequestInfo{})
	require.NoError(t, err)

	_, err = c.SetOptionSync(types.RequestSetOption{})
	require.NoError(t, err)

	_, err = c.DeliverTxSync(types.RequestDeliverTx{})
	require.NoError(t, err)

	_, err = c.CheckTxSync(types.RequestCheckTx{})
	require.NoError(t, err)

	_, err = c.QuerySync(types.RequestQuery{})
	require.NoError(t, err)

	_, err = c.CommitSync()
	require.NoError(t, err)

	_, err = c.InitChainSync(types.RequestInitChain{})
	require.NoError(t, err)

	_, err = c.BeginBlockSync(types.RequestBeginBlock{})
	require.NoError(t, err)

	_, err = c.EndBlockSync(types.RequestEndBlock{})
	require.NoError(t, err)

	_, err = c.BeginRecheckTxSync(types.RequestBeginRecheckTx{})
	require.NoError(t, err)

	_, err = c.EndRecheckTxSync(types.RequestEndRecheckTx{})
	require.NoError(t, err)

	_, err = c.ListSnapshotsSync(types.RequestListSnapshots{})
	require.NoError(t, err)

	_, err = c.OfferSnapshotSync(types.RequestOfferSnapshot{})
	require.NoError(t, err)

	_, err = c.LoadSnapshotChunkSync(types.RequestLoadSnapshotChunk{})
	require.NoError(t, err)

	_, err = c.ApplySnapshotChunkSync(types.RequestApplySnapshotChunk{})
	require.NoError(t, err)
}

type sampleApp struct {
	types.BaseApplication
}

func getResponseCallback(t *testing.T, called *bool) abcicli.ResponseCallback {
	return func (res *types.Response) {
		require.NotNil(t, res)
		*called = true
	}
}
