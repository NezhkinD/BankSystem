package db

import (
	"BankSystem/internal/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

// BuildDSN строит строку подключения к PostgreSQL
func BuildDSN(cfg config.DBConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)
}

// New создает новый пул соединений с базой данных
func New(ctx context.Context, cfg config.DBConfig) (*pgxpool.Pool, error) {
	dsn := BuildDSN(cfg)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка при парсинге конфигурации пула: %w", err)
	}

	// Настройка параметров пула соединений
	poolConfig.MaxConns = 10
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = 1 * time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	// Создание пула с таймаутом
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании пула соединений: %w", err)
	}

	// Проверка соединения
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ошибка при проверке соединения с БД: %w", err)
	}

	return pool, nil
}
