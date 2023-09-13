package stores

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/discentem/cavorite/internal/stores/pluginproto"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CommandNotFoundError struct {
	Err error
}

func (p *CommandNotFoundError) Error() string {
	return fmt.Sprintf("command not found: %s", p.Err.Error())
}

func (p *CommandNotFoundError) Unwrap() error {
	return p.Err
}

var (
	PluginSet = plugin.PluginSet{
		"store": &grpcPlugin{},
	}

	HandshakeConfig = plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "BASIC_PLUGIN",
		MagicCookieValue: "cavorite",
	}

	// FIXME: make configurable?
	HLog = hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})
)

// ConfigData is hypothetical data you'd want to pass to the plugin upon initialization
type ConfigData struct {
	A string
	B int64
	C bool
}

// pluginStore wraps Store. pluginStore is the actual interface used by the plugin so that ConfigData can be passed to the plugin on init,
// but plugins only need to implement Store
type pluginStore interface {
	Store
	Init(*ConfigData) error
}

// pluginGRPCClient is the client cavorite uses to communicate with the plugin, translating pluginStore calls to GRPC calls
type pluginGRPCClient struct {
	pluginproto.PluginClient
}

func (p *pluginGRPCClient) Upload(ctx context.Context, objects ...string) error {
	_, err := p.PluginClient.Upload(ctx, &pluginproto.Objects{Objects: objects})
	return err
}

func (p *pluginGRPCClient) Retrieve(ctx context.Context, objects ...string) error {
	_, err := p.PluginClient.Retrieve(ctx, &pluginproto.Objects{Objects: objects})
	return err
}

func (p *pluginGRPCClient) GetOptions() (Options, error) {
	opts, err := p.PluginClient.GetOptions(context.Background(), &emptypb.Empty{})
	if err != nil {
		return Options{}, err
	}

	return Options{
		BackendAddress:        opts.BackendAddress,
		MetadataFileExtension: opts.MetadataFileExtension,
		Region:                opts.Region,
	}, nil
}

func (p *pluginGRPCClient) Init(config *ConfigData) error {
	if _, err := p.PluginClient.Init(context.Background(), &pluginproto.ConfigData{
		A: config.A,
		B: config.B,
		C: config.C,
	}); err != nil {
		return err
	}

	return nil
}

func (p *pluginGRPCClient) Close() error {
	return nil
}

// pluginGRPCServer is the server the plugin uses to communicate with cavorite, translating GRPC calls to pluginStore calls
type pluginGRPCServer struct {
	pluginStore
	pluginproto.UnimplementedPluginServer
}

func (p *pluginGRPCServer) Upload(ctx context.Context, objects *pluginproto.Objects) (*emptypb.Empty, error) {
	err := p.pluginStore.Upload(ctx, objects.Objects...)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (p *pluginGRPCServer) Retrieve(ctx context.Context, objects *pluginproto.Objects) (*emptypb.Empty, error) {
	err := p.pluginStore.Retrieve(ctx, objects.Objects...)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (p *pluginGRPCServer) GetOptions(_ context.Context, _ *emptypb.Empty) (*pluginproto.Options, error) {
	opts, err := p.pluginStore.GetOptions()
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &pluginproto.Options{
		BackendAddress:        opts.BackendAddress,
		MetadataFileExtension: opts.MetadataFileExtension,
		Region:                opts.Region,
	}, nil
}

func (p *pluginGRPCServer) Init(_ context.Context, config *pluginproto.ConfigData) (*emptypb.Empty, error) {
	if err := p.pluginStore.Init(&ConfigData{
		A: config.A,
		B: config.B,
		C: config.C,
	}); err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &emptypb.Empty{}, nil
}

// grpcPlugin is the go-plugin wrapper around pluginStore that implements plugin.GRPCPlugin using pluginGRPCClient and pluginGRPCServer
type grpcPlugin struct {
	plugin.Plugin
	pluginStore
}

func (p *grpcPlugin) GRPCServer(_ *plugin.GRPCBroker, server *grpc.Server) error {
	pluginproto.RegisterPluginServer(server, &pluginGRPCServer{pluginStore: p.pluginStore})
	return nil
}

func (p *grpcPlugin) GRPCClient(ctx context.Context, _ *plugin.GRPCBroker, client *grpc.ClientConn) (interface{}, error) {
	return &pluginGRPCClient{PluginClient: pluginproto.NewPluginClient(client)}, nil
}

// PluggableStore is the Store used by cavorite that wraps go-plugin
type PluggableStore struct {
	client *plugin.Client
	pluginStore
}

func NewPluggableStore(config *ConfigData, cmd *exec.Cmd) (*PluggableStore, error) {
	if _, err := exec.LookPath(cmd.Path); err != nil {
		return nil, &CommandNotFoundError{Err: err}
	}
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  HandshakeConfig,
		Plugins:          PluginSet,
		Cmd:              cmd,
		Logger:           HLog,
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	})

	// connect to RPC
	rpcClient, err := client.Client()
	if err != nil {
		client.Kill()
		return nil, fmt.Errorf("could not create rpc client: %w", err)
	}

	// get plugin as <any>
	raw, err := rpcClient.Dispense("store")
	if err != nil {
		client.Kill()
		return nil, fmt.Errorf("could not dispense plugin: %w", err)
	}

	// assert to Store
	store := raw.(pluginStore)

	// pass config data to plugin
	if err := store.Init(config); err != nil {
		client.Kill()
		return nil, fmt.Errorf("could not initialize plugin: %w", err)
	}

	return &PluggableStore{
		client:      client,
		pluginStore: store,
	}, nil
}

func (p *PluggableStore) Close() error {
	p.client.Kill()
	return nil
}

// PluginServer is the interface used by the plugin
type PluginServer interface {
	// Config returns the ConfigData passed to the plugin during initialization
	Config() *ConfigData
	// Wait blocks until the plugin server stops
	Wait()
}

// pluginServer is a wrapper around Store that implements pluginStore and PluginServer
type pluginServer struct {
	Store
	config     *ConfigData
	configWait chan struct{}
	wait       chan struct{}
}

// NewPluginServer creates and starts a new PluginServer
func NewPluginServer(store Store, logger hclog.Logger) PluginServer {
	s := &pluginServer{
		Store:      store,
		configWait: make(chan struct{}),
		wait:       make(chan struct{}),
	}

	PluginSet["store"] = &grpcPlugin{pluginStore: s}

	go func() {
		plugin.Serve(&plugin.ServeConfig{
			HandshakeConfig: HandshakeConfig,
			Plugins:         PluginSet,
			Logger:          logger,
			GRPCServer:      plugin.DefaultGRPCServer,
		})
		close(s.wait)
	}()

	// wait for config from Init
	<-s.configWait

	return s
}

func (s *pluginServer) Init(config *ConfigData) error {
	s.config = config
	close(s.configWait)
	return nil
}

func (s *pluginServer) Config() *ConfigData {
	return s.config
}

func (s *pluginServer) Wait() {
	<-s.wait
}
