package mongodb

import (
	"context"
	"golang.org/x/sync/semaphore"
	"gopkg.in/mgo.v2"
	"log"
	"ranking/config"
)

var sema = semaphore.NewWeighted(10)

// GetMongodbSessionBySemaphore 获取MongoDb连接会话，使用信号量限制最多同时10个协程访问连接
func GetMongodbSessionBySemaphore() *mgo.Session {
	err := sema.Acquire(context.TODO(), 1)
	if err != nil {
		log.Panicln("acquire semaphore failed", err)
	}
	session, err := mgo.Dial(config.AppConfig.MongoDb.Address)
	if err != nil {
		log.Panicln("connect to mongodb failed", err)
	}
	sema.Release(1)
	return session
}

var bufChannel = make(chan struct{}, 10)

// GetMongodbSession 获取MongoDb连接会话，使用管道限制最多同时10个协程访问连接
func GetMongodbSession() *mgo.Session {
	bufChannel <- struct{}{}
	session, err := mgo.Dial(config.AppConfig.MongoDb.Address)
	if err != nil {
		log.Panicln("connect to mongodb failed", err)
	}
	<-bufChannel
	return session
}

// CloseMongodbSession 关闭MongoDb连接会话
func CloseMongodbSession(session *mgo.Session) {
	session.Close()
}
