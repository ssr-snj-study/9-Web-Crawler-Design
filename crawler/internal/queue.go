package internal

import "fmt"

// StringQueue 구조체 정의
type StringQueue struct {
	items []string
}

// Enqueue: 큐에 문자열 추가
func (q *StringQueue) Enqueue(item string) {
	q.items = append(q.items, item)
	fmt.Println("insert queue: ", item)
}

// Dequeue: 큐에서 문자열 제거 및 반환
func (q *StringQueue) Dequeue() (string, bool) {
	if len(q.items) == 0 {
		return "", false // 큐가 비어있음
	}
	dequeued := q.items[0] // 첫 번째 요소 가져오기
	q.items = q.items[1:]  // 첫 번째 요소 제거
	fmt.Println("Dequeue queue: ", dequeued)
	return dequeued, true
}

// Peek: 큐의 첫 번째 문자열 확인
func (q *StringQueue) Peek() (string, bool) {
	if len(q.items) == 0 {
		return "", false // 큐가 비어있음
	}
	return q.items[0], true
}

// IsEmpty: 큐가 비어있는지 확인
func (q *StringQueue) IsEmpty() bool {
	return len(q.items) == 0
}
