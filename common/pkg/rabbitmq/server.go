package rabbitmq

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	rabbitmq "github.com/wagslane/go-rabbitmq"
)

// ConsumerConfig describes a single consumer subscription.
// It mirrors the "routes register handlers" style used in the gateway,
// but for RabbitMQ consumers.
type ConsumerConfig struct {
	// Name is used for logs and troubleshooting.
	Name string

	Queue    string
	Exchange string
	// Kind must be one of: "fanout", "direct", "topic", "headers".
	Kind string

	// RoutingKeys is optional (fanout typically uses none).
	RoutingKeys []string
}

// Server owns a single RabbitMQ connection and one-or-more consumers.
// It provides a DI-friendly surface for registering handlers and then
// starting everything until shutdown.
type Server struct {
	conn   *rabbitmq.Conn
	logger *slog.Logger

	mu        sync.Mutex
	consumers []*rabbitmq.Consumer
	closed    bool
}

func NewServer(
	url string,
	logger *slog.Logger,
	connOpts ...func(*rabbitmq.ConnectionOptions),
) (*Server, error) {
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}
	if url == "" {
		return nil, errors.New("rabbitmq url is empty")
	}

	conn, err := rabbitmq.NewConn(url, connOpts...)
	if err != nil {
		return nil, err
	}

	return &Server{
		conn:   conn,
		logger: logger,
	}, nil
}

func (s *Server) AddConsumer(cfg ConsumerConfig, handler rabbitmq.Handler) error {
	if handler == nil {
		return errors.New("handler is nil")
	}

	if cfg.Name == "" {
		cfg.Name = cfg.Queue
	}
	if cfg.Queue == "" {
		return errors.New("consumer queue is empty")
	}
	if cfg.Exchange == "" {
		return errors.New("consumer exchange is empty")
	}
	if cfg.Kind == "" {
		return errors.New("consumer exchange kind is empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return errors.New("server is closed")
	}

	opts := []func(*rabbitmq.ConsumerOptions){
		rabbitmq.WithConsumerOptionsExchangeName(cfg.Exchange),
		rabbitmq.WithConsumerOptionsExchangeKind(cfg.Kind),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	}

	// For fanout exchanges, routing keys are ignored by RabbitMQ, but bindings
	// still exist. If you supply none, we bind to "" by default.
	if len(cfg.RoutingKeys) == 0 {
		opts = append(opts, rabbitmq.WithConsumerOptionsRoutingKey(""))
	} else {
		for _, rk := range cfg.RoutingKeys {
			opts = append(opts, rabbitmq.WithConsumerOptionsRoutingKey(rk))
		}
	}

	consumer, err := rabbitmq.NewConsumer(s.conn, cfg.Queue, opts...)
	if err != nil {
		return fmt.Errorf("[%s] new consumer: %w", cfg.Name, err)
	}

	s.consumers = append(s.consumers, consumer)

	// Start the consumer loop in its own goroutine. Any fatal behavior should be
	// decided by the caller; here we just log errors and stop that consumer.
	go func() {
		s.logger.Info("rabbitmq consumer starting",
			"name", cfg.Name,
			"queue", cfg.Queue,
			"exchange", cfg.Exchange,
			"kind", cfg.Kind,
		)

		if err := consumer.Run(handler); err != nil {
			s.logger.Error("rabbitmq consumer stopped",
				"name", cfg.Name,
				"error", err,
			)
		}
	}()

	return nil
}

// GetPublisher returns a publisher that uses the server's connection.
// This allows handlers to publish messages to other exchanges.
func (s *Server) GetPublisher() (*rabbitmq.Publisher, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil, errors.New("server is closed")
	}

	return rabbitmq.NewPublisher(s.conn)
}

// StartAndBlock blocks until SIGINT/SIGTERM, then closes all consumers and the
// underlying connection.
func (s *Server) StartAndBlock() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	s.logger.Info("rabbitmq server started; awaiting shutdown signal")
	<-sigs

	s.logger.Info("rabbitmq server shutting down")
	s.Close()
}

func (s *Server) Close() {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return
	}
	s.closed = true
	consumers := append([]*rabbitmq.Consumer(nil), s.consumers...)
	s.consumers = nil
	s.mu.Unlock()

	// Close consumers first, then the connection.
	for _, c := range consumers {
		c.Close()
	}
	_ = s.conn.Close()
}

