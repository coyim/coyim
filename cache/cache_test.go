package cache

import (
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type CacheSuite struct{}

var _ = Suite(&CacheSuite{})

func (cs *CacheSuite) Test_New_ReturnsNotNil(c *C) {
	r := New()
	c.Assert(r, Not(IsNil))
}

func (cs *CacheSuite) Test_NewWithExpiry_ReturnsNotNil(c *C) {
	r := NewWithExpiry()
	c.Assert(r, Not(IsNil))
}

func (cs *CacheSuite) Test_EmptyCache_ReturnsCorrectResults(c *C) {
	r := New()

	_, found := r.Get("foo")
	c.Assert(found, Equals, false)

	found = r.Put("bar", "flux")
	c.Assert(found, Equals, false)

	found = r.Has("quux")
	c.Assert(found, Equals, false)

	found = r.Remove("blarg")
	c.Assert(found, Equals, false)
}

func (cs *CacheSuite) Test_CacheWithExistingValue_ReturnsCorrectResults(c *C) {
	r := New()
	r.Clear()
	r.Put("foo", "bar")

	v, found := r.Get("foo")
	c.Assert(found, Equals, true)
	c.Assert(v.(string), Equals, "bar")

	found = r.Put("foo", "flux")
	c.Assert(found, Equals, true)

	found = r.Has("foo")
	c.Assert(found, Equals, true)

	found = r.Remove("foo")
	c.Assert(found, Equals, true)

	r.Clear()

	found = r.Has("foo")
	c.Assert(found, Equals, false)
}

func (cs *CacheSuite) Test_PutIfAbsent_CallsFunction(c *C) {
	r := New()

	called := false
	found := r.PutIfAbsent("foo", func(key string) interface{} {
		called = true
		return "bla"
	})

	c.Assert(found, Equals, true)
	c.Assert(called, Equals, true)
	v, _ := r.Get("foo")
	c.Assert(v.(string), Equals, "bla")

	called = false
	found = r.PutIfAbsent("foo", func(key string) interface{} {
		called = true
		return "bla2"
	})

	c.Assert(found, Equals, false)
	c.Assert(called, Equals, false)
	v, _ = r.Get("foo")
	c.Assert(v.(string), Equals, "bla")
}

func (cs *CacheSuite) Test_Expiry_WillRemoveEntry(c *C) {
	r := NewWithExpiry()

	found := r.PutTimed("foo", 200*time.Millisecond, "bar")
	c.Assert(found, Equals, false)
	c.Assert(r.Has("foo"), Equals, true)
	time.Sleep(500 * time.Millisecond)
	c.Assert(r.Has("foo"), Equals, false)
}

func (cs *CacheSuite) Test_Expiry_DealsOKWithOverridingATimedEntryWithARegularOne(c *C) {
	r := NewWithExpiry()

	r.PutTimed("foo", 200*time.Millisecond, "bar")
	r.Put("foo", "blarg")
	time.Sleep(500 * time.Millisecond)
	c.Assert(r.Has("foo"), Equals, true)
	v, _ := r.Get("foo")
	c.Assert(v.(string), Equals, "blarg")
}

func (cs *CacheSuite) Test_PutIfAbsent_withExpiry_WillCallTheFunction(c *C) {
	r := NewWithExpiry()

	called := false
	added := r.PutTimedIfAbsent("foo", 200*time.Millisecond, func(key string) interface{} {
		called = true
		return "bla"
	})
	c.Assert(added, Equals, true)
	c.Assert(called, Equals, true)

	called = false
	added = r.PutTimedIfAbsent("foo", 200*time.Millisecond, func(key string) interface{} {
		called = true
		return "bla2"
	})
	c.Assert(added, Equals, false)
	c.Assert(called, Equals, false)

	c.Assert(r.Has("foo"), Equals, true)
	time.Sleep(500 * time.Millisecond)
	c.Assert(r.Has("foo"), Equals, false)
}
