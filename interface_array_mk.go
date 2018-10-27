package main

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"math"
)
type Hashable interface {
    GetHash() string
}

type MerkleTreeI interface {
    ComputeTree(hashes []Hashable)
    GetRoot() string
    GetTree() []string
    SetTree(leavesCount int, tree []string) error
    GetPath(hash Hashable) []MTPathNode
    VerifyPath(hash Hashable, path []MTPathNode) bool
    GetPathByIndex(idx int) []MTPathNode
}

const (
    Left = 0
    Right = 1
)

type MTPathNode struct {
    Hash string
    Side byte
}

func Hash(text string) string {
	hash := sha1.New()
	io.WriteString(hash, text)
	return hex.EncodeToString(hash.Sum(nil));

}

func MHash(h1 string, h2 string) string {
    return Hash(h1 + h2)
}

type Txn struct {
	data string
}

func (a Txn) GetHash() string {
	return Hash(a.data)
}

func HashObject(h Hashable) string {
	return h.GetHash()
}

func CalculatenumNodes(curr int) int {
	total := curr
	for curr != 1 {
		curr = curr / 2
		if curr % 2 != 0 {
			if curr != 1 {
				curr = curr + 1
			}
		}
		total += curr
	}
	return total
}
func getLeftChild(a int) int{
	return (2*a) + 1
}

func getRightChild(a int) int{
	return 2*(a+1)
}

func buildLevel(t *MerkleArray, level int) {
	nodeIndex := int(math.Pow(2, float64(level))) - 1
	nodeLast :=  int(math.Pow(2, float64(level+1))) - 2
	nodeBuilt := 0
	for ; nodeIndex <= nodeLast; nodeBuilt++ {
		l := getLeftChild(nodeIndex)
		r := getRightChild(nodeIndex)
		lChildHash := t.tree[l].Hash
		rChildHash := t.tree[r].Hash
		if lChildHash == "" {
			if nodeBuilt % 2 != 0 {
				t.tree[nodeIndex].Hash = t.tree[nodeIndex-1].Hash
				nodeBuilt++
			}
			return
		}
		t.tree[nodeIndex].Hash = MHash(lChildHash, rChildHash)
		t.tree[l].path.Side = Left
		t.tree[l].path.Hash = t.tree[nodeIndex].Hash
		t.tree[r].path.Side = Right
		t.tree[r].path.Hash = t.tree[nodeIndex].Hash
		nodeIndex++
	}
}

func (t *MerkleArray) VerifyPath(hash Hashable, path []MTPathNode) bool {
	result := false
	actualPath := t.GetPath(hash)
	countA := len(actualPath)
	countB := len(path)
	if countA != countB {return result}
	for i:=0; i < countA; i++ {
		if(actualPath[i] != path[i]) {
			return result
		}
	}
	result = true
	return result
}

func (t *MerkleArray) GetPathByIndex(i int) []MTPathNode {
	var path []MTPathNode
	if i >= t.numLeaves {return path}
	if t.numLeaves == 1 {
		path = append(path, t.tree[0].path)
		return path
	}
	index := int(math.Pow(2, float64(t.height))) - 1
	index = index + i
	for {
		path = append(path, t.tree[index].path)
		index = (index-1) / 2
		if index == 0 {
			return path
		}
	}
	return path
}
func (t *MerkleArray) GetPath(hash Hashable) []MTPathNode {
	var path []MTPathNode
	i := t.ToPath[hash.GetHash()]
	if i > t.numLeaves {return path}
	if t.numLeaves == 1 {
		path = append(path, t.tree[0].path)
		return path
	}
	i--
	index := int(math.Pow(2, float64(t.height))) - 1
	index = index + i
	for {
		path = append(path, t.tree[index].path)
		index = (index-1) / 2
		if index == 0 {
			return path
		}
	}
	return path
}

func (t *MerkleArray) GetTree() []string {
	var tree_str []string
        for _, c := range t.tree {
		if c.Hash != "" {
			tree_str = append(tree_str, c.Hash)
		}
	}
	return tree_str
}

func(t *MerkleArray) GetRoot() string{
		return t.tree[0].Hash
}

func (t *MerkleArray) ComputeTree(hashes []Hashable) {
	leaf_len := len(hashes)
	i:=0
	j:=0
	index:=0

	if leaf_len == 0 {return}

	t.numLeaves = leaf_len
	if leaf_len % 2 == 1 {
		t.numLeaves = leaf_len + 1
		hashes = append(hashes, hashes[leaf_len-1])
	}
	t.height = int(math.Ceil(math.Log2(float64(t.numLeaves))))
	t.numNodes = CalculatenumNodes(t.numLeaves)
	treelen := int(math.Pow(2, float64(t.height+1))) - 1
	t.tree = make([]TreeElement, treelen)
	i=0
	t.ToPath = make(map[string]int)
	for index = int(math.Pow(2, float64(t.height))) -1 ; j < t.numLeaves; index++ {
		t.tree[index].Hash = hashes[i].GetHash()
		if j < leaf_len {
			t.ToPath[t.tree[index].Hash] = j + 1
		}
		i++
		j++
	}
	level:= t.height-1
	for ; level >= 0; level-- {
		buildLevel(t, level);
	}

}
type TreeElement struct {
        Hash	string
	path	MTPathNode
}

type MerkleArray struct {
	tree []TreeElement
	height int
	numLeaves int
	numNodes int
	ToPath	map[string]int
}
