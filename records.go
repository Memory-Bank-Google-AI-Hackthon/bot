package main

import "sync"

type Record struct {
	UserId   string `json:"userId"`
	UserName string `json:"userName"`
	Message  string `json:"message"`
}

type Records struct {
	sync.Mutex
	Messages []Record `json:"messages"`
}

func NewRecords() *Records {
	return &Records{}
}

func (r *Records) AddRecord(record Record) {
	r.Lock()
	defer r.Unlock()
	r.Messages = append(r.Messages, record)
}

func (r *Records) GetRecords() []Record {
	r.Lock()
	defer r.Unlock()
	return r.Messages
}

func (r *Records) Clear() {
	r.Lock()
	defer r.Unlock()
	r.Messages = []Record{}
}
