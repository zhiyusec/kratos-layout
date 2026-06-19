package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/go-kratos/kratos/contrib/otel/v3/tracing"
	"github.com/go-kratos/kratos/v3"
	"github.com/go-kratos/kratos/v3/log"
	"github.com/go-kratos/kratos/v3/transport/grpc"
	"github.com/zhiyusec/kratos-layout/internal/conf"
	"github.com/zhiyusec/kratos-layout/internal/logic"
	"github.com/zhiyusec/kratos-layout/internal/repo"
	"github.com/zhiyusec/kratos-layout/internal/server"
	"github.com/zhiyusec/kratos-layout/internal/service"

	_ "go.uber.org/automaxprocs"
)

var (
	Name     string
	Version  string
	flagconf string
	id, _    = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf ../../configs")
}

func newApp(logger *slog.Logger, gs *grpc.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(gs),
	)
}

func main() {
	flag.Parse()
	logger := log.NewLogger(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		}),
		log.WithExtractor(tracing.TraceAttrs),
	).With(
		slog.String("service.id", id),
		slog.String("service.name", Name),
		slog.String("service.version", Version),
	)
	log.SetDefault(logger)

	bc, err := conf.Load(flagconf)
	if err != nil {
		panic(err)
	}

	// 3 layers: repo -> logic -> service
	todoRepo := repo.NewTodoRepo()
	todoLogic := logic.NewTodoLogic(todoRepo)
	todoSvc := service.NewTodoService(todoLogic)

	grpcSrv := server.NewGRPCServer(bc.Server, todoSvc)
	app := newApp(logger, grpcSrv)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
