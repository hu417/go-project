package test

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type Node struct {
	Text     string `json:"text"`
	Children []Node `json:"children"`
}

var (
	stRootDirTree   string
	stSeparatorTree string
	iRootNode       Node
)

const stJsonFileName = "dir.json"

func loadJsonTree() {
	stSeparatorTree = string(filepath.Separator)
	stWorkDir, _ := os.Getwd()
	stRootDirTree = stWorkDir[:strings.LastIndex(stWorkDir, stSeparatorTree)]

	gnJsonFileBytes, _ := os.ReadFile(stWorkDir + stSeparatorTree + stJsonFileName)
	err := json.Unmarshal(gnJsonFileBytes, &iRootNode)
	if err != nil {
		panic("Load Json Data Error: " + err.Error())
	}
}

func parseNode(iNode Node, stParentDir string) {
	if iNode.Text != "" {
		createDir02(iNode, stParentDir)
	}

	if stParentDir != "" {
		stParentDir = stParentDir + stSeparatorTree
	}

	if iNode.Text != "" {
		stParentDir = stParentDir + iNode.Text
	}

	for _, v := range iNode.Children {
		parseNode(v, stParentDir)
	}
}

func createDir02(iNode Node, stParentDir string) {
	stDirPath := stRootDirTree + stSeparatorTree
	if stParentDir != "" {
		stDirPath = stDirPath + stParentDir + stSeparatorTree
	}
	stDirPath = stDirPath + iNode.Text

	err := os.MkdirAll(stDirPath, fs.ModePerm)
	if err != nil {
		panic("Create Dir Error: " + err.Error())
	}
}

func TestGenerateDirTree(t *testing.T) {
	loadJsonTree()
	parseNode(iRootNode, "")
}
