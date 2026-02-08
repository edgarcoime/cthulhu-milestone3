module github.com/cthulhu-platform/gateway

go 1.25.6

require (
	github.com/cthulhu-platform/auth v0.0.0
	github.com/cthulhu-platform/common v0.0.0
	github.com/cthulhu-platform/filemanager v0.0.0
	github.com/cthulhu-platform/lifecycle v0.0.0
	github.com/cthulhu-platform/proto v0.0.0
	github.com/gofiber/fiber/v2 v2.52.10
	github.com/spf13/viper v1.21.0
)

require (
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	go.opentelemetry.io/otel v1.38.0 // indirect
	go.opentelemetry.io/otel/trace v1.38.0 // indirect
	golang.org/x/net v0.48.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251029180050-ab9386a59fda // indirect
	google.golang.org/grpc v1.78.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

require (
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/google/uuid v1.6.0
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/sagikazarmark/locafero v0.11.0 // indirect
	github.com/samber/slog-fiber v1.20.1
	github.com/sourcegraph/conc v0.3.1-0.20240121214520-5f936abd7ae8 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.59.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
)

replace github.com/cthulhu-platform/lifecycle => ../lifecycle

replace github.com/cthulhu-platform/auth => ../auth

replace github.com/cthulhu-platform/common => ../common

replace github.com/cthulhu-platform/filemanager => ../filemanager

replace github.com/cthulhu-platform/proto => ../proto
