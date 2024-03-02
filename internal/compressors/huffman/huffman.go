package huffman

import (
	cbits "archiver/internal/bits"
	"archiver/internal/btree"
	"archiver/internal/queue"
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	TempFile     = "temp.txt"
	BufSize      = 4096
	AlphabetSize = 256
)

type HuffmanCompressor struct {
	archiveName string
	codes       []cbits.Bits
}

func New(archiveName string) *HuffmanCompressor {
	codes := make([]cbits.Bits, AlphabetSize)
	for i := 0; i < AlphabetSize; i++ {
		codes[i] = cbits.NewBits(0, 0)
	}

	return &HuffmanCompressor{
		archiveName: archiveName,
		codes:       codes,
	}
}

func (c *HuffmanCompressor) Compress() error {
	f, err := os.Open(c.archiveName)
	if err != nil {
		return err
	}
	defer f.Close()

	temp, err := os.OpenFile(TempFile, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer func() {
		temp.Close()
		os.Remove(TempFile)
	}()

	var header strings.Builder

	count := make([]int, AlphabetSize)
	buf := make([]byte, BufSize)
	for {
		n, err := f.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		for i := 0; i < n; i++ {
			count[buf[i]]++
		}
	}

	root := c.buildTree(count)
	c.dfs(root, cbits.NewBits(0, 0), 0)

	// запись результата
	for i, c := range count {
		if c != 0 {
			header.WriteString(fmt.Sprintf("%d,%d;", i, c))
		}
	}

	temp.WriteString(header.String())

	bits := cbits.NewBits(0, 0)
	f.Seek(0, 0)
	for {
		n, err := f.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		for i := 0; i < n; i++ {
			code := c.codes[buf[i]]
			for j := 0; j < code.GetCount(); j++ {
				b, _ := code.GetBit(j)
				if b {
					bits.AppendBit(1)
				} else {
					bits.AppendBit(0)
				}
			}
		}
	}
	temp.WriteString(fmt.Sprintf("%d", bits.GetCount()))
	temp.WriteString("\n")
	_, err = temp.Write(bits.GetBytes())
	if err != nil {
		return err
	}

	f, err = os.OpenFile(c.archiveName, os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer f.Close()

	temp.Seek(0, 0)
	for {
		n, err := temp.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		_, err = f.Write(buf[:n])
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *HuffmanCompressor) Decompress() error {
	f, err := os.Open(c.archiveName)
	if err != nil {
		return err
	}
	defer f.Close()

	temp, err := os.OpenFile(TempFile, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer func() {
		temp.Close()
		os.Remove(TempFile)
	}()

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	countInfo := scanner.Text()
	f.Seek(int64(len(countInfo)+1), 0)

	count := make([]int, AlphabetSize)
	counts := strings.Split(countInfo, ";")
	countOfBits, err := strconv.Atoi(counts[len(counts)-1])
	if err != nil {
		return err
	}
	for _, c := range counts[:len(counts)-1] {
		ci := strings.Split(c, ",")
		ch, err := strconv.Atoi(ci[0])
		if err != nil {
			return err
		}
		freq, err := strconv.Atoi(ci[1])
		if err != nil {
			return err
		}
		count[ch] = freq
	}

	root := c.buildTree(count)

	bytes, _ := io.ReadAll(f)
	bits := cbits.NewLoadBytes(bytes, countOfBits)

	node := root
	for i := 0; i < bits.GetCount(); i++ {
		b, err := bits.GetBit(i)
		if err != nil {
			return err
		}

		if b {
			node = node.Right
		} else {
			node = node.Left
		}

		if node.Left == nil && node.Right == nil {
			_, err = temp.Write([]byte(node.Value))
			if err != nil {
				return nil
			}
			node = root
		}
	}
	if node.Left == nil && node.Right == nil {
		_, err = temp.Write([]byte(node.Value))
		if err != nil {
			return nil
		}
	}

	f, err = os.OpenFile(c.archiveName, os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := make([]byte, BufSize)
	temp.Seek(0, 0)
	for {
		n, err := temp.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		_, err = f.Write(buf[:n])
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *HuffmanCompressor) dfs(node *btree.TreeNode, bits cbits.Bits, numOfBit int) {
	if node == nil {
		return
	}

	if node.Left == nil && node.Right == nil && len(node.Value) == 1 {
		if numOfBit == 0 {
			bits.AppendBit(0)
		}
		c.codes[node.Value[0]] = cbits.Copy(bits)
		return
	}

	bitsCopy := cbits.Copy(bits)

	if node.Left != nil {
		bits.AppendBit(0)
		c.dfs(node.Left, bits, numOfBit+1)
		bits = bitsCopy
	}

	if node.Right != nil {
		bits.AppendBit(1)
		c.dfs(node.Right, bits, numOfBit+1)
	}
}

func (c *HuffmanCompressor) buildTree(count []int) *btree.TreeNode {
	priorityQueue := queue.NewPriorityQueue[*btree.TreeNode]()
	for i, c := range count {
		if c != 0 {
			node := btree.New([]byte{byte(i)})
			priorityQueue.Insert(queue.NewListNode(node, c))
		}
	}

	var root *btree.TreeNode
	for priorityQueue.GetCount() > 1 {
		first, second := priorityQueue.Extract(), priorityQueue.Extract()
		left, right := first.Value, second.Value
		root = btree.New(append(left.Value, right.Value...))
		root.Left = left
		root.Right = right
		priorityQueue.Insert(queue.NewListNode(root, first.Priority+second.Priority))
	}
	if priorityQueue.GetCount() == 1 {
		root = priorityQueue.Extract().Value
	}

	return root
}
