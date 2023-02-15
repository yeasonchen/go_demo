package main

import (
	"fmt"
	"math"
	"net/http"
	"reflect"
)

/*
void NonRecursionTraversal(BinTree BT) {
    if (!BT) return;

    BinTree T;
    T = BT;
    StackLinked stack = CreateStackLinked();
    while (T || !IsEmptyLinked(stack)) {
        while (T) {
            PushLinked(stack, T);
            T = T->left;
        }//这里while出来的时候,T是NULL
        T = PopLinked(stack);
        printf("visit");
        T = T->right;
    }
}

BinTree T, Index;
    Index = NULL;
    T = BT;
    StackLinked stack = CreateStackLinked();
    while (T || !IsEmptyLinked(stack)) {
        while (T) {
            PushLinked(stack, T);
            T = T->left;
        }//这里while出来的时候,T是NULL
        T = GetTopLinked(stack);
        if (T->right && T->right != Index) {//右子树未被访问
            T = T->right;
        } else {//右子树已被访问或没有右子树,则可以访问当前节点了
            T = PopLinked(stack);
            printf("%d\t", T->data);
            Index = T;
            T = NULL;
        }
    }
*/

func handlerEchoCtx(w http.ResponseWriter, r *http.Request) {
	//s := r.Context().Value(http.LocalAddrContextKey)
	//w.Write([]byte("hello"))
	//fmt.Println("this is ctx ==> ", r.Context())
}

func main() {
	//fmt.Println("ahahahah")
	//
	//go func(size int32, name int64) {
	//	fmt.Println("asf")
	//}(20, 80)

	//http.HandleFunc("/", handlerEchoCtx)
	//stdoutFile, _ := os.OpenFile("./log.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModeAppend)
	//
	//syscall.Dup2(int(stdoutFile.Fd()), 1)
	//syscall.Dup2(int(stdoutFile.Fd()), 2)
	//http.ListenAndServe(":9099", nil)

	//matrix := [][]int{
	//	{1, 4, 7, 11, 15},
	//	{2, 5, 8, 12, 19},
	//	{3, 6, 9, 16, 22},
	//	{10, 13, 14, 17, 24},
	//	{18, 21, 23, 26, 30},
	//}

	//t := &TreeNode{
	//	Val: 5,
	//	Left: &TreeNode{
	//		Val: 3,
	//		Left: &TreeNode{
	//			Val:   1,
	//			Left:  nil,
	//			Right: nil,
	//		},
	//		Right: &TreeNode{
	//			Val:   4,
	//			Left:  nil,
	//			Right: nil,
	//		},
	//	},
	//	Right: &TreeNode{
	//		Val: 7,
	//		Left: &TreeNode{
	//			Val:   6,
	//			Left:  nil,
	//			Right: nil,
	//		},
	//		Right: &TreeNode{
	//			Val:   8,
	//			Left:  nil,
	//			Right: nil,
	//		},
	//	},
	//}
	//
	//fmt.Println(getK(t, 3))
	//fmt.Println(getK(t, 1))
	//fmt.Println(getK(t, 8))

	//	aaa()

	a := 12
	aType := reflect.TypeOf(a)
	fmt.Println(aType)

	aValue := reflect.ValueOf(a)
	fmt.Println(aValue.Addr())

}

func aaa() {
	// a := []int{}

	// var p interface{}
	// p = []string(nil)
	// if p == nil {
	// 	fmt.Println("sdfs")
	// }
}

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

var sK = 1

func getK(tree *TreeNode, k int) int {
	if tree == nil {
		return 0
	}
	if r := getK(tree.Left, k); r > 0 {
		return r
	}
	if k == sK {
		sK = 1
		return tree.Val
	}
	sK++
	if r := getK(tree.Right, k); r > 0 {
		return r
	}

	return 0
}

func twoEggDrop(n int) int {
	list := make([]int, n+1)
	for i := 1; i <= n; i++ {
		min := math.MaxInt
		for j := 1; j <= i; j++ {
			min = minint(maxint(list[i-j]+1, j), min)
		}
		list[i] = min
	}

	return list[n]
}

func dfs(grid [][]byte, m, n int) {
	lm, ln := len(grid), len(grid[0])
	if m < 0 || m >= lm || n < 0 || n >= ln || grid[m][n] == '0' {
		return
	}
	grid[m][n] = '0'
	dfs(grid, m, n-1)
	dfs(grid, m, n+1)
	dfs(grid, m-1, n)
	dfs(grid, m+1, n)
}

func numIslands(grid [][]byte) int {
	result := 0
	for i, g := range grid {
		for i2, b := range g {
			if b == '1' {
				result++
				// 把相同岛屿的标识为0
				dfs(grid, i, i2)
			}
		}
	}
	return result
}

func abs(num int) int {
	if num > 0 {
		return num
	}
	return -num
}

func maxint(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minint(a, b int) int {
	if a > b {
		return b
	}
	return a
}

type ListNode struct {
	Val  int
	Next *ListNode
}

type LRUCache struct {
	cap, size  int
	m          map[int]*DLinkedNode
	head, tail *DLinkedNode
}

func Constructor(capacity int) LRUCache {
	head := InitLinkedNode(0, 0)
	tail := InitLinkedNode(0, 0)
	head.Next = tail
	tail.Prev = head
	return LRUCache{
		cap:  capacity,
		size: 0,
		m:    make(map[int]*DLinkedNode, 0),
		head: head,
		tail: tail,
	}
}

func (t *LRUCache) Get(key int) int {
	if p, ok := t.m[key]; ok {
		p.Prev.Next = p.Next
		p.Next.Prev = p.Prev
		t.head.Next.Prev = p
		p.Next = t.head.Next
		t.head.Next = p
		return p.Value
	} else {
		return -1
	}
}

func (t *LRUCache) Put(key int, value int) {
	p := InitLinkedNode(key, value)
	t.m[key] = p
	t.head.Next.Prev = p
	p.Next = t.head.Next
	p.Prev = t.head
	t.head.Next = p
	if t.size == t.cap {
		last := t.tail.Prev
		last.Prev.Next = t.tail
		t.tail.Prev = last.Prev
		if len(t.m) > t.cap {
			delete(t.m, last.Key)
		}
	} else {
		t.size++
	}
}

type DLinkedNode struct {
	Key        int
	Value      int
	Prev, Next *DLinkedNode
}

func InitLinkedNode(key int, value int) *DLinkedNode {
	return &DLinkedNode{
		Key:   key,
		Value: value,
	}
}
