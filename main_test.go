package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheck(t *testing.T) {
	for i := 1; i <= 64; i++ {
		assert.Equal(t, i, len(salt(i)))
	}

	// http://www.cisco.com/c/en/us/td/docs/ios-xml/ios/security/d1/sec-d1-cr-book/sec-cr-e1.html
	// username demo8 algorithm-type sha256 secret cisco
	// username demo8 secret 8 $8$dsYGNam3K1SIJO$7nv/35M/qr6t.dVc7UY9zrJDWRVqncHub1PE9UlMQFs
	assert.Equal(t, type8("cisco", "dsYGNam3K1SIJO"), "$8$dsYGNam3K1SIJO$7nv/35M/qr6t.dVc7UY9zrJDWRVqncHub1PE9UlMQFs")

	// http://www.cisco.com/c/en/us/td/docs/ios-xml/ios/security/d1/sec-d1-cr-book/sec-cr-e1.html
	// username demo9 algorithm-type scrypt secret cisco
	// username demo9 secret 9 $9$nhEmQVczB7dqsO$X.HsgL6x1il0RxkOSSvyQYwucySCt7qFm4v7pqCxkKM
	assert.Equal(t, type9("cisco", "nhEmQVczB7dqsO"), "$9$nhEmQVczB7dqsO$X.HsgL6x1il0RxkOSSvyQYwucySCt7qFm4v7pqCxkKM")
}
