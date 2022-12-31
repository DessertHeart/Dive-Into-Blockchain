package main

// SPV验证

import "crypto/sha256"

// Merkle树结构（二叉树）
type MerkleTree struct {
	RootNode *MerkleNode
}

// Merkle树节点
type MerkleNode struct {
	Left *MerkleNode
	Right *MerkleNode
	HashData []byte
}

// 新建节点
func NewMerkleNode(left, right *MerkleNode, dataHash []byte) *MerkleNode {
	newNode := &MerkleNode{}

	// 处理HashData参数：左右节点都为空， 则说明是初始叶子节点，打包输入的data(初始为序列化后的transactions)，否则打包左右节点的hash和
	if left == nil && right == nil {
		hash := sha256.Sum256(dataHash)
		newNode.HashData = hash[:]
	} else {
		data := append(left.HashData, right.HashData...) 
		hash := sha256.Sum256(data)
		newNode.HashData = hash[:]
	}
	// 处理左右子节点
	newNode.Left = left
	newNode.Right = right

	return newNode
}

// 新建Merkle Tree
func NewMerkleTree(data [][]byte) *MerkleTree {
	var nodes []MerkleNode

	// 初始叶子结点数必须为双数，否则将复制最后一个tx的Hash
	if len(data) % 2 != 0 {
		data = append(data, data[len(data) - 1])
	}

	// 将所以的交易打包为最底层叶子结点
	for _, datum := range data {
		node := NewMerkleNode(nil, nil, datum)
		nodes = append(nodes, *node)
	}

	// 一层一层hash，两两打包故次数为len(data)/2，直至得出最终rootnode
	for i := 0; i < len(data)/2; i++ {
		var newLevel []MerkleNode

		for j := 0; j < len(nodes); j += 2 {
			node := NewMerkleNode(&nodes[j], &nodes[j+1], nil)
			newLevel = append(newLevel, *node)
		}

		// 每新一层merkle tree更新一次该层的节点
		nodes = newLevel
	}

	// 最后就是根节点哈希
	newTree := MerkleTree{&nodes[0]}
	return &newTree
}

