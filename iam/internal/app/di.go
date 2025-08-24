package app

import (
	"context"
	"fmt"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v5"

	authApi "github.com/dfg007star/go_rocket/iam/internal/api/auth/v1"
	userApi "github.com/dfg007star/go_rocket/iam/internal/api/user/v1"
	"github.com/dfg007star/go_rocket/iam/internal/config"
	"github.com/dfg007star/go_rocket/iam/internal/repository"
	authRepository "github.com/dfg007star/go_rocket/iam/internal/repository/session"
	userRepository "github.com/dfg007star/go_rocket/iam/internal/repository/user"
	"github.com/dfg007star/go_rocket/iam/internal/service"
	authService "github.com/dfg007star/go_rocket/iam/internal/service/auth"
	userService "github.com/dfg007star/go_rocket/iam/internal/service/user"
	"github.com/dfg007star/go_rocket/platform/pkg/cache"
	"github.com/dfg007star/go_rocket/platform/pkg/cache/redis"
	"github.com/dfg007star/go_rocket/platform/pkg/closer"
	"github.com/dfg007star/go_rocket/platform/pkg/logger"
	authV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/auth/v1"
	userV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/user/v1"
)

type diContainer struct {
	userV1API      userV1.UserServiceServer
	userService    service.UserService
	userRepository repository.UserRepository

	authV1API      authV1.AuthServiceServer
	authService    service.AuthService
	authRepository repository.SessionRepository

	redisPool   *redigo.Pool
	redisClient cache.RedisClient

	postgresClient *pgx.Conn
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

// AuthAPI
func (d *diContainer) AuthV1API(ctx context.Context) authV1.AuthServiceServer {
	if d.authV1API == nil {
		d.authV1API = authApi.NewAuthAPI(d.AuthService(ctx))
	}

	return d.authV1API
}

func (d *diContainer) AuthService(ctx context.Context) service.AuthService {
	if d.authService == nil {
		d.authService = authService.NewAuthService(d.AuthRepository(), d.UserRepository(ctx))
	}

	return d.authService
}

func (d *diContainer) AuthRepository() repository.SessionRepository {
	if d.authRepository == nil {
		d.authRepository = authRepository.NewSessionRepository(d.RedisClient())
	}

	return d.authRepository
}

// UserAPI
func (d *diContainer) UserV1API(ctx context.Context) userV1.UserServiceServer {
	if d.userV1API == nil {
		d.userV1API = userApi.NewUserAPI(d.UserService(ctx))
	}

	return d.userV1API
}

func (d *diContainer) UserService(ctx context.Context) service.UserService {
	if d.userService == nil {
		d.userService = userService.NewUserService(d.UserRepository(ctx))
	}

	return d.userService
}

func (d *diContainer) UserRepository(ctx context.Context) repository.UserRepository {
	if d.userRepository == nil {
		d.userRepository = userRepository.NewUserRepository(d.PostgresClient(ctx))
	}

	return d.userRepository
}

// Redis
func (d *diContainer) RedisPool() *redigo.Pool {
	if d.redisPool == nil {
		d.redisPool = &redigo.Pool{
			MaxIdle:     config.AppConfig().Redis.MaxIdle(),
			IdleTimeout: config.AppConfig().Redis.IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", config.AppConfig().Redis.Address())
			},
		}
	}

	return d.redisPool
}

func (d *diContainer) RedisClient() cache.RedisClient {
	if d.redisClient == nil {
		d.redisClient = redis.NewClient(d.RedisPool(), logger.Logger(), config.AppConfig().Redis.ConnectionTimeout())
	}

	return d.redisClient
}

// Postgres
func (d *diContainer) PostgresClient(ctx context.Context) *pgx.Conn {
	if d.postgresClient == nil {
		con, err := pgx.Connect(ctx, config.AppConfig().Postgres.URI())
		if err != nil {
			panic(fmt.Errorf("failed to connect to database: %w", err))
		}

		closer.AddNamed("Postgres client", func(ctx context.Context) error {
			return con.Close(ctx)
		})

		err = con.Ping(ctx)
		if err != nil {
			panic(fmt.Errorf("database is unavailable: %w", err))
		}

		d.postgresClient = con
	}

	return d.postgresClient
}
