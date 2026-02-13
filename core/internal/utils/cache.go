// Package utils 内存缓存工具
// 提供线程安全的内存缓存，支持 TTL 过期
package utils

import (
	"sync"
	"time"
)

// CacheItem 缓存项
type CacheItem struct {
	Value      interface{}
	ExpireTime time.Time
}

// IsExpired 判断缓存项是否过期
func (item *CacheItem) IsExpired() bool {
	if item.ExpireTime.IsZero() {
		return false // 永不过期
	}
	return time.Now().After(item.ExpireTime)
}

// Cache 内存缓存结构
type Cache struct {
	items map[string]*CacheItem
	mutex sync.RWMutex
}

// NewCache 创建新的缓存实例
func NewCache() *Cache {
	cache := &Cache{
		items: make(map[string]*CacheItem),
	}
	// 启动后台清理协程
	go cache.cleanup()
	return cache
}

// Set 设置缓存（带 TTL）
// ttl: 过期时间，0 表示永不过期
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var expireTime time.Time
	if ttl > 0 {
		expireTime = time.Now().Add(ttl)
	}

	c.items[key] = &CacheItem{
		Value:      value,
		ExpireTime: expireTime,
	}
}

// SetSeconds 设置缓存（秒为单位）
// ttlSeconds: 过期时间（秒），0 表示永不过期
func (c *Cache) SetSeconds(key string, value interface{}, ttlSeconds int) {
	c.Set(key, value, time.Duration(ttlSeconds)*time.Second)
}

// Get 获取缓存值
// 返回: 值和是否存在（过期的项也返回 false）
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	if item.IsExpired() {
		return nil, false
	}

	return item.Value, true
}

// GetString 获取字符串类型的缓存值
func (c *Cache) GetString(key string) (string, bool) {
	value, exists := c.Get(key)
	if !exists {
		return "", false
	}
	if str, ok := value.(string); ok {
		return str, true
	}
	return "", false
}

// Delete 删除缓存
func (c *Cache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.items, key)
}

// Exists 检查缓存是否存在（且未过期）
func (c *Cache) Exists(key string) bool {
	_, exists := c.Get(key)
	return exists
}

// Clear 清空所有缓存
func (c *Cache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.items = make(map[string]*CacheItem)
}

// Size 获取缓存项数量（包含已过期但未清理的）
func (c *Cache) Size() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return len(c.items)
}

// cleanup 后台清理过期缓存
func (c *Cache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		for key, item := range c.items {
			if item.IsExpired() {
				delete(c.items, key)
			}
		}
		c.mutex.Unlock()
	}
}

// GetOrSet 获取缓存，如果不存在则设置
// loader: 加载函数，返回值和错误
// ttl: 缓存过期时间
func (c *Cache) GetOrSet(key string, loader func() (interface{}, error), ttl time.Duration) (interface{}, error) {
	// 先尝试获取
	if value, exists := c.Get(key); exists {
		return value, nil
	}

	// 加载数据
	value, err := loader()
	if err != nil {
		return nil, err
	}

	// 设置缓存
	c.Set(key, value, ttl)
	return value, nil
}

// 全局缓存实例
var globalCache *Cache
var cacheOnce sync.Once

// GetGlobalCache 获取全局缓存实例（单例）
func GetGlobalCache() *Cache {
	cacheOnce.Do(func() {
		globalCache = NewCache()
	})
	return globalCache
}

// 价格缓存的专用方法

// PriceCacheKey 生成价格缓存键
func PriceCacheKey(symbol string) string {
	return "price:" + symbol
}

// RateCacheKey 生成汇率缓存键
func RateCacheKey(from, to string) string {
	return "rate:" + from + ":" + to
}

// BalanceCacheKey 生成余额缓存键
func BalanceCacheKey(sourceType string, sourceID uint) string {
	return "balance:" + sourceType + ":" + string(rune(sourceID))
}
