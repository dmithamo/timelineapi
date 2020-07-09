// package authservice manages authentication/authorization in the application
package dbservice

import "github.com/gomodule/redigo/redis"

// ConnectCache establishes a new conn to the redis db powering auth
func ConnectCache(dsn *string) (redis.Conn, error) {
	conn, err := redis.Dial("tcp", *dsn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
