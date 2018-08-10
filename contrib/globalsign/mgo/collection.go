package mgo

import (
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/globalsign/mgo"
)

// Collection provides a mgo.Collection along with
// data used for APM Tracing.
type Collection struct {
	*mgo.Collection
	cfg mongoConfig
}

// Create invokes and traces Collection.Create
func (c *Collection) Create(info *mgo.CollectionInfo) error {
	span := newChildSpanFromContext(c.cfg)
	err := c.Collection.Create(info)
	span.Finish(tracer.WithError(err))
	return err
}

// DropCollection invokes and traces Collection.DropCollection
func (c *Collection) DropCollection() error {
	span := newChildSpanFromContext(c.cfg)
	err := c.Collection.DropCollection()
	span.Finish(tracer.WithError(err))
	return err
}

// EnsureIndexKey invokes and traces Collection.EnsureIndexKey
func (c *Collection) EnsureIndexKey(key ...string) error {
	span := newChildSpanFromContext(c.cfg)
	err := c.Collection.EnsureIndexKey(key...)
	span.Finish(tracer.WithError(err))
	return err
}

// EnsureIndex invokes and traces Collection.EnsureIndex
func (c *Collection) EnsureIndex(index mgo.Index) error {
	span := newChildSpanFromContext(c.cfg)
	err := c.Collection.EnsureIndex(index)
	span.Finish(tracer.WithError(err))
	return err
}

// DropIndex invokes and traces Collection.DropIndex
func (c *Collection) DropIndex(key ...string) error {
	span := newChildSpanFromContext(c.cfg)
	err := c.Collection.DropIndex(key...)
	span.Finish(tracer.WithError(err))
	return err
}

// DropIndexName invokes and traces Collection.DropIndexName
func (c *Collection) DropIndexName(name string) error {
	span := newChildSpanFromContext(c.cfg)
	err := c.Collection.DropIndexName(name)
	span.Finish(tracer.WithError(err))
	return err
}

// DropAllIndexes invokes and traces Collection.DropAllIndexes
/*func (c *Collection) DropAllIndexes() error {
	span := newChildSpanFromContext(c.cfg)
	err := c.Collection.DropAllIndexes()
	span.Finish(tracer.WithError(err))
	return err
}*/

// Indexes invokes and traces Collection.Indexes
func (c *Collection) Indexes() (indexes []mgo.Index, err error) {
	span := newChildSpanFromContext(c.cfg)
	indexes, err = c.Collection.Indexes()
	span.Finish(tracer.WithError(err))
	return indexes, err
}

// Insert invokes and traces Collectin.Insert
func (c *Collection) Insert(docs ...interface{}) error {
	span := newChildSpanFromContext(c.cfg)
	err := c.Collection.Insert(docs...)
	span.Finish(tracer.WithError(err))
	return err
}

// Find invokes and traces Collection.Find
func (c *Collection) Find(query interface{}) *Query {
	return &Query{
		Query: c.Collection.Find(query),
		cfg:   c.cfg,
	}
}

// FindId invokes and traces Collection.FindId
func (c *Collection) FindId(id interface{}) *Query { // nolint
	return &Query{
		Query: c.Collection.FindId(id),
		cfg:   c.cfg,
	}
}

// Count invokes and traces Collection.Count
func (c *Collection) Count() (n int, err error) {
	span := newChildSpanFromContext(c.cfg)
	n, err = c.Collection.Count()
	span.Finish(tracer.WithError(err))
	return n, err
}

// Bulk creates a trace ready wrapper around Collection.Bulk
func (c *Collection) Bulk() *Bulk {
	return &Bulk{
		Bulk: c.Collection.Bulk(),
		cfg:  c.cfg,
	}
}

// Update invokes and traces Collection.Update
func (c *Collection) Update(selector interface{}, update interface{}) error {
	span := newChildSpanFromContext(c.cfg)
	err := c.Collection.Update(selector, update)
	span.Finish(tracer.WithError(err))
	return err
}

// UpdateId invokes and traces Collection.UpdateId
func (c *Collection) UpdateId(id interface{}, update interface{}) error { // nolint
	span := newChildSpanFromContext(c.cfg)
	err := c.Collection.UpdateId(id, update)
	span.Finish(tracer.WithError(err))
	return err
}

// UpdateAll invokes and traces Collection.UpdateAll
func (c *Collection) UpdateAll(selector interface{}, update interface{}) (info *mgo.ChangeInfo, err error) {
	span := newChildSpanFromContext(c.cfg)
	info, err = c.Collection.UpdateAll(selector, update)
	span.Finish(tracer.WithError(err))
	return info, err
}

// Upsert invokes and traces Collection.Upsert
func (c *Collection) Upsert(selector interface{}, update interface{}) (info *mgo.ChangeInfo, err error) {
	span := newChildSpanFromContext(c.cfg)
	info, err = c.Collection.Upsert(selector, update)
	span.Finish(tracer.WithError(err))
	return info, err
}

// UpsertId invokes and traces Collection.UpsertId
func (c *Collection) UpsertId(id interface{}, update interface{}) (info *mgo.ChangeInfo, err error) { // nolint
	span := newChildSpanFromContext(c.cfg)
	info, err = c.Collection.UpsertId(id, update)
	span.Finish(tracer.WithError(err))
	return info, err
}

// Remove invokes and traces Collection.Remove
func (c *Collection) Remove(selector interface{}) error {
	span := newChildSpanFromContext(c.cfg)
	err := c.Collection.Remove(selector)
	span.Finish(tracer.WithError(err))
	return err
}

// RemoveId invokes and traces Collection.RemoveId
func (c *Collection) RemoveId(id interface{}) error { // nolint
	span := newChildSpanFromContext(c.cfg)
	err := c.Collection.RemoveId(id)
	span.Finish(tracer.WithError(err))
	return err
}

// RemoveAll invokes and traces Collection.RemoveAll
func (c *Collection) RemoveAll(selector interface{}) (info *mgo.ChangeInfo, err error) {
	span := newChildSpanFromContext(c.cfg)
	info, err = c.Collection.RemoveAll(selector)
	span.Finish(tracer.WithError(err))
	return info, err
}

// Repair invokes and traces Collection.Repair
func (c *Collection) Repair() *Iter {
	span := newChildSpanFromContext(c.cfg)
	iter := c.Collection.Repair()
	span.Finish()
	return &Iter{
		Iter: iter,
		cfg:  c.cfg,
	}
}
