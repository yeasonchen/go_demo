// Package creational 单例模式
package creational

import "sync"

type Singleton struct{}

// 饿汉模式

var singleton *Singleton

func init() {
	singleton = &Singleton{}
}

func GetInstance() *Singleton {
	return singleton
}

// 懒汉模式

var (
	lazySingleton *Singleton
	once          = &sync.Once{}
)

func GetLazySingleton() *Singleton {
	if lazySingleton == nil {
		once.Do(func() {
			lazySingleton = &Singleton{}
		})
	}
	return lazySingleton
}
