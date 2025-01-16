package cache

import (
	"hash/fnv"
	"sync"
	"sync/atomic"
	"time"
)

// CacheItem 缓存项结构
type CacheItem struct {
	Value      interface{}
	ExpireTime time.Time
}

// Segment 缓存分段
type Segment struct {
	data sync.Map
}

// Cache 缓存主结构
type Cache struct {
	segments      [32]Segment   // 分段锁
	enabled       bool          // 是否启用缓存
	ttl           time.Duration // 默认过期时间
	cleanInterval time.Duration // 清理间隔
	itemCount     int64         // 总缓存项数
	stopClear     chan struct{} // 停止清理信号
	metrics       *CacheMetrics // 缓存指标
	maxSize       int64         // 最大缓存项数
}

// CacheMetrics 缓存指标
type CacheMetrics struct {
	hits   int64 // 命中次数
	misses int64 // 未命中次数
}

// NewCache 创建新的缓存实例
func NewCache(enabled bool, defaultTTL time.Duration, clearup time.Duration, maxSize int64) *Cache {
	cache := &Cache{
		enabled:       enabled,
		ttl:           defaultTTL,
		cleanInterval: clearup,
		stopClear:     make(chan struct{}),
		metrics:       &CacheMetrics{},
		maxSize:       maxSize,
	}

	// 初始化分段
	for i := range cache.segments {
		cache.segments[i] = Segment{}
	}

	// 启动清理器
	if enabled {
		go cache.startCleaner()
	}

	return cache
}

// getSegmentIndex 获取分段索引
func (c *Cache) getSegmentIndex(key string) uint {
	hash := fnv.New32a()
	hash.Write([]byte(key))
	return uint(hash.Sum32()) % uint(len(c.segments))
}

// Get 获取缓存值
func (c *Cache) Get(key string) (interface{}, bool) {
	if !c.enabled {
		return nil, false
	}

	segmentIndex := c.getSegmentIndex(key)
	segment := &c.segments[segmentIndex]

	value, found := segment.data.Load(key)
	if !found {
		atomic.AddInt64(&c.metrics.misses, 1)
		return nil, false
	}

	item, ok := value.(*CacheItem)
	if !ok || time.Now().After(item.ExpireTime) {
		segment.data.Delete(key)
		atomic.AddInt64(&c.itemCount, -1)
		atomic.AddInt64(&c.metrics.misses, 1)
		return nil, false
	}

	atomic.AddInt64(&c.metrics.hits, 1)
	return item.Value, true
}

// Set 设置缓存值
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	if !c.enabled {
		return
	}

	if ttl == 0 {
		ttl = c.ttl
	}

	if atomic.LoadInt64(&c.itemCount) >= c.maxSize {
		return // 超过最大缓存项数限制
	}

	segmentIndex := c.getSegmentIndex(key)
	segment := &c.segments[segmentIndex]

	segment.data.Store(key, &CacheItem{
		Value:      value,
		ExpireTime: time.Now().Add(ttl),
	})
	atomic.AddInt64(&c.itemCount, 1)
}

// Clear 清空缓存
func (c *Cache) Clear() {
	if !c.enabled {
		return
	}

	for i := range c.segments {
		segment := &c.segments[i]
		segment.data = sync.Map{}
	}
	atomic.StoreInt64(&c.itemCount, 0)
}

// startCleaner 启动清理器
func (c *Cache) startCleaner() {
	ticker := time.NewTicker(c.cleanInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.clearExpired()
		case <-c.stopClear:
			return
		}
	}
}

// clearExpired 清理过期项
func (c *Cache) clearExpired() {
	now := time.Now()
	removed := int64(0)

	for i := range c.segments {
		segment := &c.segments[i]
		segment.data.Range(func(key, value interface{}) bool {
			if item, ok := value.(*CacheItem); ok && now.After(item.ExpireTime) {
				segment.data.Delete(key)
				removed++
			}
			return true
		})
	}

	if removed > 0 {
		atomic.AddInt64(&c.itemCount, -removed)
	}
}

// Close 关闭缓存
func (c *Cache) Close() {
	if c.enabled {
		close(c.stopClear)
	}
}

// GetMetrics 获取缓存指标
func (c *Cache) GetMetrics() (hits, misses, count int64) {
	hits = atomic.LoadInt64(&c.metrics.hits)
	misses = atomic.LoadInt64(&c.metrics.misses)
	count = atomic.LoadInt64(&c.itemCount)
	return
}
