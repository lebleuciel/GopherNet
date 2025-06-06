// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"

	"gophernet/pkg/db/ent/migrate"

	"gophernet/pkg/db/ent/burrow"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
	// Burrow is the client for interacting with the Burrow builders.
	Burrow *BurrowClient
}

// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
	client := &Client{config: newConfig(opts...)}
	client.init()
	return client
}

func (c *Client) init() {
	c.Schema = migrate.NewSchema(c.driver)
	c.Burrow = NewBurrowClient(c.config)
}

type (
	// config is the configuration for the client and its builder.
	config struct {
		// driver used for executing database requests.
		driver dialect.Driver
		// debug enable a debug logging.
		debug bool
		// log used for logging on debug mode.
		log func(...any)
		// hooks to execute on mutations.
		hooks *hooks
		// interceptors to execute on queries.
		inters *inters
	}
	// Option function to configure the client.
	Option func(*config)
)

// newConfig creates a new config for the client.
func newConfig(opts ...Option) config {
	cfg := config{log: log.Println, hooks: &hooks{}, inters: &inters{}}
	cfg.options(opts...)
	return cfg
}

// options applies the options on the config object.
func (c *config) options(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
	if c.debug {
		c.driver = dialect.Debug(c.driver, c.log)
	}
}

// Debug enables debug logging on the ent.Driver.
func Debug() Option {
	return func(c *config) {
		c.debug = true
	}
}

// Log sets the logging function for debug mode.
func Log(fn func(...any)) Option {
	return func(c *config) {
		c.log = fn
	}
}

// Driver configures the client driver.
func Driver(driver dialect.Driver) Option {
	return func(c *config) {
		c.driver = driver
	}
}

// Open opens a database/sql.DB specified by the driver name and
// the data source name, and returns a new client attached to it.
// Optional parameters can be added for configuring the client.
func Open(driverName, dataSourceName string, options ...Option) (*Client, error) {
	switch driverName {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		drv, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}
		return NewClient(append(options, Driver(drv))...), nil
	default:
		return nil, fmt.Errorf("unsupported driver: %q", driverName)
	}
}

// ErrTxStarted is returned when trying to start a new transaction from a transactional client.
var ErrTxStarted = errors.New("ent: cannot start a transaction within a transaction")

// Tx returns a new transactional client. The provided context
// is used until the transaction is committed or rolled back.
func (c *Client) Tx(ctx context.Context) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, ErrTxStarted
	}
	tx, err := newTx(ctx, c.driver)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = tx
	return &Tx{
		ctx:    ctx,
		config: cfg,
		Burrow: NewBurrowClient(cfg),
	}, nil
}

// BeginTx returns a transactional client with specified options.
func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, errors.New("ent: cannot start a transaction within a transaction")
	}
	tx, err := c.driver.(interface {
		BeginTx(context.Context, *sql.TxOptions) (dialect.Tx, error)
	}).BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = &txDriver{tx: tx, drv: c.driver}
	return &Tx{
		ctx:    ctx,
		config: cfg,
		Burrow: NewBurrowClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		Burrow.
//		Query().
//		Count(ctx)
func (c *Client) Debug() *Client {
	if c.debug {
		return c
	}
	cfg := c.config
	cfg.driver = dialect.Debug(c.driver, c.log)
	client := &Client{config: cfg}
	client.init()
	return client
}

// Close closes the database connection and prevents new queries from starting.
func (c *Client) Close() error {
	return c.driver.Close()
}

// Use adds the mutation hooks to all the entity clients.
// In order to add hooks to a specific client, call: `client.Node.Use(...)`.
func (c *Client) Use(hooks ...Hook) {
	c.Burrow.Use(hooks...)
}

// Intercept adds the query interceptors to all the entity clients.
// In order to add interceptors to a specific client, call: `client.Node.Intercept(...)`.
func (c *Client) Intercept(interceptors ...Interceptor) {
	c.Burrow.Intercept(interceptors...)
}

// Mutate implements the ent.Mutator interface.
func (c *Client) Mutate(ctx context.Context, m Mutation) (Value, error) {
	switch m := m.(type) {
	case *BurrowMutation:
		return c.Burrow.mutate(ctx, m)
	default:
		return nil, fmt.Errorf("ent: unknown mutation type %T", m)
	}
}

// BurrowClient is a client for the Burrow schema.
type BurrowClient struct {
	config
}

// NewBurrowClient returns a client for the Burrow from the given config.
func NewBurrowClient(c config) *BurrowClient {
	return &BurrowClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `burrow.Hooks(f(g(h())))`.
func (c *BurrowClient) Use(hooks ...Hook) {
	c.hooks.Burrow = append(c.hooks.Burrow, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `burrow.Intercept(f(g(h())))`.
func (c *BurrowClient) Intercept(interceptors ...Interceptor) {
	c.inters.Burrow = append(c.inters.Burrow, interceptors...)
}

// Create returns a builder for creating a Burrow entity.
func (c *BurrowClient) Create() *BurrowCreate {
	mutation := newBurrowMutation(c.config, OpCreate)
	return &BurrowCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Burrow entities.
func (c *BurrowClient) CreateBulk(builders ...*BurrowCreate) *BurrowCreateBulk {
	return &BurrowCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *BurrowClient) MapCreateBulk(slice any, setFunc func(*BurrowCreate, int)) *BurrowCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &BurrowCreateBulk{err: fmt.Errorf("calling to BurrowClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*BurrowCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &BurrowCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Burrow.
func (c *BurrowClient) Update() *BurrowUpdate {
	mutation := newBurrowMutation(c.config, OpUpdate)
	return &BurrowUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *BurrowClient) UpdateOne(b *Burrow) *BurrowUpdateOne {
	mutation := newBurrowMutation(c.config, OpUpdateOne, withBurrow(b))
	return &BurrowUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *BurrowClient) UpdateOneID(id int) *BurrowUpdateOne {
	mutation := newBurrowMutation(c.config, OpUpdateOne, withBurrowID(id))
	return &BurrowUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Burrow.
func (c *BurrowClient) Delete() *BurrowDelete {
	mutation := newBurrowMutation(c.config, OpDelete)
	return &BurrowDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *BurrowClient) DeleteOne(b *Burrow) *BurrowDeleteOne {
	return c.DeleteOneID(b.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *BurrowClient) DeleteOneID(id int) *BurrowDeleteOne {
	builder := c.Delete().Where(burrow.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &BurrowDeleteOne{builder}
}

// Query returns a query builder for Burrow.
func (c *BurrowClient) Query() *BurrowQuery {
	return &BurrowQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeBurrow},
		inters: c.Interceptors(),
	}
}

// Get returns a Burrow entity by its id.
func (c *BurrowClient) Get(ctx context.Context, id int) (*Burrow, error) {
	return c.Query().Where(burrow.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *BurrowClient) GetX(ctx context.Context, id int) *Burrow {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *BurrowClient) Hooks() []Hook {
	return c.hooks.Burrow
}

// Interceptors returns the client interceptors.
func (c *BurrowClient) Interceptors() []Interceptor {
	return c.inters.Burrow
}

func (c *BurrowClient) mutate(ctx context.Context, m *BurrowMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&BurrowCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&BurrowUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&BurrowUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&BurrowDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Burrow mutation op: %q", m.Op())
	}
}

// hooks and interceptors per client, for fast access.
type (
	hooks struct {
		Burrow []ent.Hook
	}
	inters struct {
		Burrow []ent.Interceptor
	}
)
