package services

import "github.com/garyburd/redigo/redis"

type Storage interface {
	// Closes the connection with the storage server
	Close()
	// Retrieves the timestamp of the last time the channel was updated
	GetLastChannelRefresh(string) uint64
	SetLastChannelRefresh(channel string, lastRefresh uint64)
	GetChannelURL(channel string) (string, error)
	SetChannelURL(channel, url string)
	GetNewsByChannel(channel string) (string, error)
	SetNewsByChannel(channel, news string)
	EraseChannel(channel string)
}

type storage struct {
	redis.Conn
}

var (
	LAST_REFRESH = "em:last_refresh"
	NEWS_CACHE   = "em:news_cache"
	CHANNEL_URLS = "em:channel_urls"
	TWEETS_CACHE = "em:tweets_cache:"
)

// NewStorage returns a new connection to the Storage service
func NewStorage(address string) (Storage, error) {
	var (
		conn redis.Conn
		err  error
	)

	if conn, err = redis.Dial("tcp", address); err != nil {
		return nil, err
	}

	return &storage{conn}, nil
}

func (s *storage) Close() {
	s.Conn.Close()
}

func (s *storage) getStringFromHash(key, item string) (string, error) {
	buff, err := redis.Bytes(s.Conn.Do("HGET", key, item))
	if err != nil {
		return "", err
	}

	return string(buff), nil
}

func (s *storage) GetLastChannelRefresh(channel string) uint64 {
	lastRefresh, err := redis.Uint64(s.Conn.Do("HGET", LAST_REFRESH, channel))
	if err != nil {
		return uint64(0)
	}

	return lastRefresh
}

func (s *storage) SetLastChannelRefresh(channel string, lastRefresh uint64) {
	s.Conn.Do("HSET", LAST_REFRESH, channel, lastRefresh)
}

func (s *storage) GetChannelURL(channel string) (string, error) {
	return s.getStringFromHash(CHANNEL_URLS, channel)
}

func (s *storage) SetChannelURL(channel, url string) {
	s.Conn.Do("HSET", CHANNEL_URLS, channel, url)
}

func (s *storage) GetNewsByChannel(channel string) (string, error) {
	return s.getStringFromHash(NEWS_CACHE, channel)
}

func (s *storage) SetNewsByChannel(channel, news string) {
	s.Conn.Do("HSET", NEWS_CACHE, channel, news)
}

func (s *storage) EraseChannel(channel string) {
	s.Conn.Do("HDEL", NEWS_CACHE, channel)
	s.Conn.Do("HDEL", LAST_REFRESH, channel)
	s.Conn.Do("DEL", TWEETS_CACHE+channel)
}
