package core

import (
	"container/list"
	"errors"
	"log"
	"sync"
)

const initCursor = 2000

var _ = log.Printf

var q = newQueue()

type Team struct {
	Size int `json:"size"`
	ID   int `json:"id"`
}

type Queue struct {
	li   *list.List
	cur  int
	lock *sync.RWMutex
}

func newQueue() *Queue {
	q := Queue{}
	q.li = list.New()
	q.cur = initCursor
	q.lock = new(sync.RWMutex)
	return &q
}

func AddTeamToQueue(teamSize int) (*Team, error) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.cur += 1
	t := Team{Size: teamSize, ID: q.cur}
	q.li.PushBack(&t)
	return &t, nil
}

func ResetQueue() error {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.li.Init()
	q.cur = initCursor
	return nil
}

// 将拉去previousID之后count个team，previousID可以是0，表示从头开始拉取
func GetTeamsFromQueue(previousID int, count int) ([]*Team, error) {
	if previousID < 0 {
		return nil, errors.New("previousID must >= 0")
	}
	q.lock.RLock()
	defer q.lock.RUnlock()
	if q.li.Len() == 0 {
		return nil, nil
	}
	var firstElement *list.Element = nil
	var remainTeamCount = q.li.Len()
	if previousID == 0 {
		firstElement = q.li.Front()
	} else {
		for e := q.li.Front(); e != nil; e = e.Next() {
			remainTeamCount -= 1
			team := e.Value.(*Team)
			if team.ID == previousID {
				firstElement = e.Next()
				break
			}
		}
	}
	if firstElement == nil {
		return nil, nil
	}
	resultCount := MinInt(remainTeamCount, count)
	result := make([]*Team, resultCount)
	for e, i := firstElement, 0; e != nil; e, i = e.Next(), i+1 {
		result[i] = e.Value.(*Team)
	}
	return result, nil
}
