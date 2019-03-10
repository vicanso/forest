package service

import (
	"testing"
	"time"

	"github.com/vicanso/forest/util"
)

func TestGetRedisClient(t *testing.T) {
	if GetRedisClient() == nil {
		t.Fatalf("get client fail")
	}
}

func TestLock(t *testing.T) {
	key := util.RandomString(8)
	ttl := 10 * time.Millisecond
	success, err := Lock(key, ttl)
	if err != nil || !success {
		t.Fatalf("the first time lock fail, %v", err)
	}

	success, err = Lock(key, ttl)
	if err != nil || success {
		t.Fatalf("the second time lock should be fail, %v", err)
	}
	time.Sleep(ttl)
	success, err = Lock(key, ttl)
	if err != nil || !success {
		t.Fatalf("wait for ttl, lock fail, %v", err)
	}
}

func TestLockWithDone(t *testing.T) {
	key := util.RandomString(8)
	ttl := 10 * time.Second
	success, done, err := LockWithDone(key, ttl)
	if err != nil || !success {
		t.Fatalf("the first time lock with done fail, %v", err)
	}

	success, _, err = LockWithDone(key, ttl)
	if err != nil || success {
		t.Fatalf("the second time lock with done should be fail, %v", err)
	}

	done()
	success, _, err = LockWithDone(key, ttl)
	if err != nil || !success {
		t.Fatalf("after done lock with done fail, %v", err)
	}
}

func TestIncWithTTL(t *testing.T) {
	key := util.RandomString(8)
	ttl := 10 * time.Millisecond
	count, err := IncWithTTL(key, ttl)
	if err != nil || count != 1 {
		t.Fatalf("inc fail, %v", err)
	}
	count, err = IncWithTTL(key, ttl)
	if err != nil || count != 2 {
		t.Fatalf("inc fail, %v", err)
	}

	time.Sleep(ttl)
	count, err = IncWithTTL(key, ttl)
	if err != nil || count != 1 {
		t.Fatalf("inc fail, %v", err)
	}
}
