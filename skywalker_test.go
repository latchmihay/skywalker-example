package skywalker_test

import (
	"github.com/dixonwille/skywalker"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

var (
	rootFolder = "testingFolder"
	subFolders = []string{"the", "subfolder", "sub", "sub/folder/subfolder"}
	subFiles   = []string{"just.txt", "a.log", "few.pdf", "files", "config.hcl", "komposition.hcl"}
)

type WalkerWorker struct {
	*sync.Mutex
	found []string
}

func (ww *WalkerWorker) Work(path string) {
	//This is where the necessary work should be done.
	//This will get concurrently so make sure it is thread safe if you need info across threads.
	ww.Lock()
	defer ww.Unlock()
	ww.found = append(ww.found, path)
}

func NewWW() *WalkerWorker {
	ww := new(WalkerWorker)
	ww.Mutex = new(sync.Mutex)
	return ww
}

func TestSkywalker_test(t *testing.T) {
	require := require.New(t)

	require.NoError(standUp())
	defer tearDown()

	ww := NewWW()
	sw := skywalker.New(rootFolder, ww)
	sw.FilesOnly = true
	sw.ListType = skywalker.LTWhitelist
	sw.List = []string{"**komposition.hcl"}
	require.NoError(sw.Walk())
	t.Log(ww.found)
}

func standUp() error {
	for _, sf := range subFolders {
		if err := os.MkdirAll(filepath.Join(rootFolder, sf), 0777); err != nil {
			return err
		}
		for _, f := range subFiles {
			file, err := os.OpenFile(filepath.Join(rootFolder, sf, f), os.O_RDONLY|os.O_CREATE, 0666)
			if err != nil {
				return err
			}
			defer file.Close()
		}
	}
	return nil
}

func tearDown() error {
	return os.RemoveAll(rootFolder)
}
