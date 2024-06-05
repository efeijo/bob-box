package metadata

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type MetadataFile struct {
	mu           sync.Mutex
	states       State
	file         *os.File
	currentState map[string]int64
}
type State struct {
	AllStates []map[string]int64 `json:"states"`
}

func NewMetadataFile(configFilePath, hiddenFileName string) (*MetadataFile, error) {
	// create hidden file with metadata for folder
	file, err := os.Create(configFilePath + hiddenFileName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &MetadataFile{
		currentState: make(map[string]int64),
		states: State{
			AllStates: make([]map[string]int64, 0),
		},
		file: file,
	}, nil
}

// close mnetdata file
func (m *MetadataFile) Close() error {
	return m.file.Close()
}

func (m *MetadataFile) PersistFile() error {
	m.file.Truncate(0)
	jsonBytes, err := json.Marshal(m.states)
	if err != nil {
		return err
	}
	_, err = m.file.Write(jsonBytes)
	return err
}

func (m *MetadataFile) AddFile(path string, size int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.currentState[path] = size
	m.states.AllStates = append(m.states.AllStates, m.currentState)

	m.PersistFile()
}

func (m *MetadataFile) RemoveFile(path string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.currentState, path)
	m.states.AllStates = append(m.states.AllStates, m.currentState)
	m.PersistFile()
}
