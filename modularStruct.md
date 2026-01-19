
chi tiết từng folder
Đây là một yêu cầu cực kỳ quan trọng. Rất nhiều dự án có cấu trúc thư mục đẹp nhưng nội dung bên trong lại để lộn xộn (ví dụ: viết logic SQL trong Controller, hay viết logic validate trong Model).
Để đảm bảo project của bạn "Đúng 100%" theo chuẩn Clean Architecture kết hợp Modular Monolith, tôi sẽ đi vào chi tiết từng folder chứa cái gì, file gì và code gì.
Chúng ta sẽ đi từ Ngoài vào Trong (từ tầng giao tiếp đến tầng cốt lõi), và đi từ Chung đến Riêng.
1. cmd/ (Entry Point - Cửa ngõ)
Đây là nơi ứng dụng bắt đầu. Nó không chứa business logic nào cả.
  cmd/api/main.go:
  Chứa: Hàm main().
  Nhiệm vụ:
 Load Config (env, yaml).
 Khởi tạo Logger (Zap, Logrus).
 Kết nối Database (Postgres, Redis).
 Wiring (Nối dây): Gọi InitializeApp() (của Wire) hoặc tự khởi tạo các module.
 Chạy HTTP Server (Gin/Echo/Fiber run).
2. internal/common/ (Shared Kernel - Tài nguyên dùng chung)
Nơi chứa các code kỹ thuật (technical code) mà mọi module đều cần dùng.
  server/: Cấu hình Gin/Echo engine, setup middleware global (CORS, Recovery).
  logs/: Wrapper cho logger. Ví dụ hàm logs.Info(), logs.Error() để cả project log theo 1 chuẩn format JSON.
  errors/: Định nghĩa các lỗi chung. Ví dụ: errors.ErrNotFound, errors.ErrUnauthorized.
  database/: Hàm khởi tạo kết nối DB (NewPostgresDB), cấu hình connection pool.
  utils/: Các hàm tiện ích thuần túy. Ví dụ: HashPassword(), GenerateUUID(), ValidateEmailRegex().
  middleware/: Các middleware dùng chung như RequestID, LoggerMiddleware.
3. internal/modules/<module_name>/ (Trái tim của Modular)
Đây là phần quan trọng nhất. Hãy lấy ví dụ module identity.
A. domain/ (Enterprise Business Rules - Cốt lõi)
Đây là lớp trong cùng, không phụ thuộc vào ai cả (No dependencies).
  entity.go (hoặc user.go):
  Chứa struct User.
  Lưu ý: Chỉ chứa các field dữ liệu thuần túy và các method logic nội tại (ví dụ: User.IsActive()). Hạn chế dùng tag json hay gorm ở đây nếu muốn tuân thủ Clean Arch triệt để (nhưng có thể du di tag json cho tiện).
  interfaces.go (hoặc repository.go):
  Chứa: Interface định nghĩa các hành động lưu trữ. Ví dụ: type UserRepository interface { FindByID(...) ... }.
  Tại sao: Để tầng Usecase code theo interface này mà không cần biết bên dưới là SQL hay Mongo (Dependency Inversion).
  errors.go: Các lỗi nghiệp vụ riêng của module. Ví dụ: ErrUserAlreadyExists, ErrPasswordWeak.
B. usecase/ (Application Business Rules - Logic ứng dụng)
Đây là nơi xử lý các yêu cầu của người dùng ("Người dùng muốn làm gì?").
  service.go (Interface):
  Định nghĩa các chức năng mà module này cung cấp cho bên ngoài (Handler hoặc Module khác).
  Ví dụ: Register(ctx, req), Login(ctx, req).
  register_user.go / login_user.go (Implementation):
  Chứa: Logic thực thi.
  Luồng đi:
 Nhận dữ liệu đầu vào.
 Validate nghiệp vụ (Email đã tồn tại chưa?).
 Gọi domain.User để xử lý logic nội tại.
 Gọi UserRepository (Interface) để lưu xuống DB.
 Bắn Event (nếu cần).
 Trả về kết quả.
  Tuyệt đối không: Không chứa code xử lý HTTP (như gin.Context), không chứa câu lệnh SQL.
C. infrastructure/ (Frameworks & Drivers - Cơ sở hạ tầng)
Nơi duy nhất biết về Database, External API.
  postgres/user_repository.go:
  Implement UserRepository interface đã định nghĩa ở domain.
  Chứa code Gorm/SQLx: db.Create(&user), SELECT * FROM users.
  redis/: Code lưu cache, session.
  external_api/: Nếu module Identity cần gọi Google API hay Facebook API, code gọi HTTP Client nằm ở đây.
D. delivery/ (Interface Adapters - Giao tiếp)
Nơi tiếp nhận input từ thế giới bên ngoài (HTTP, gRPC, CLI).
  http/v1/dto.go (Data Transfer Objects):
  Chứa các struct Request/Response. Ví dụ: RegisterRequest, LoginResponse.
  Chứa tag json:"email" binding:"required". Đây là nơi validate format dữ liệu (email đúng dạng không, password đủ dài không).
  http/v1/handler.go:
  Chứa các hàm xử lý HTTP (func (h \*Handler) Register(c \*gin.Context)).
  Nhiệm vụ:
 Parse request body vào DTO (c.ShouldBindJSON).
 Gọi Usecase.Register().
 Map kết quả trả về JSON response (c.JSON(200, result)).
  http/v1/routes.go: Đăng ký đường dẫn URL (POST /register, GET /me).
E. module.go (Module Configuration)
  Nhiệm vụ: Đóng gói module.
  Chứa struct Module và hàm NewModule.
  Nơi duy nhất import cả infrastructure, usecase, delivery để lắp ráp chúng lại (trừ khi dùng Wire ở ngoài).
4. migrations/ (Database Changes)
  Chứa các file .sql để tạo bảng, sửa bảng.
  Ví dụ:
  000001\_create\_users_table.up.sql (Lệnh CREATE TABLE)
  000001\_create\_users_table.down.sql (Lệnh DROP TABLE - để rollback)
5. configs/ (Configuration)
  config.yaml: File chứa giá trị config mặc định.
  config.go: Struct Go để map file yaml vào code (sử dụng Viper).
6. deployments/ (Vận hành)
  Dockerfile: Cách đóng gói app thành Docker Image.
  docker-compose.yaml: Chạy toàn bộ hệ thống (App + Postgres + Redis) ở local.
Bảng Tóm tắt: Code gì nằm ở đâu?
Loại Code
Nằm ở Folder nào?
Ví dụ cụ thể
Struct User
internal/modules/identity/domain
type User struct {...}
Interface Repo
internal/modules/identity/domain
type UserRepo interface {...}
Logic đăng ký
internal/modules/identity/usecase
func (s *service) Register(...)
Câu lệnh SQL
internal/modules/identity/infrastructure
db.Where("email = ?", email).First(&user)
Parse JSON
internal/modules/identity/delivery
c.ShouldBindJSON(&req)
Route URL
internal/modules/identity/delivery
r.POST("/login", h.Login)
Connect DB
internal/common/database
gorm.Open(...)
Logger
internal/common/logs
zap.NewProduction()
Tại sao chia kỹ thế này lại là "Đúng"?
 Dễ thay đổi DB: Bạn muốn đổi từ Postgres sang MongoDB? Chỉ cần viết lại folder infrastructure, các folder domain, usecase, delivery giữ nguyên 100%.
 Dễ thay đổi Framework: Bạn chán Gin muốn qua Fiber? Chỉ cần viết lại folder delivery, logic nghiệp vụ trong usecase giữ nguyên.
 Dễ Test: Bạn có thể viết Unit Test cho usecase mà không cần Database thật (bằng cách Mock cái Interface Repo trong domain).
Next Step: Bạn có muốn tôi viết mẫu chi tiết nội dung của dto.go (trong delivery) và service.go (trong usecase) để xem cách dữ liệu chuyển đổi từ JSON -> DTO -> Entity như thế nào cho đúng chuẩn không?
Chào bạn, một câu hỏi rất hay và đầy tham vọng. Trong kỹ thuật phần mềm, khái niệm "đúng 100%" là tương đối vì nó phụ thuộc vào quy mô và nghiệp vụ. Tuy nhiên, để đạt được tiêu chuẩn Modular Monolith chuẩn mực (Industry Standard) trong Go, cấu trúc cần đảm bảo các tiêu chí:
High Cohesion (Kết dính cao): Code liên quan đến một nghiệp vụ nằm chung một chỗ.
Low Coupling (Liên kết lỏng): Các module không được import trực tiếp implementation của nhau (tránh Spaghetti code).
Separation of Concerns: Tách biệt rõ logic business, database và transport (HTTP/gRPC).
Dưới đây là cấu trúc Project Backend Modular Monolith chuẩn và đầy đủ nhất áp dụng các nguyên lý DDD (Domain-Driven Design) và Clean Architecture, được các công ty công nghệ lớn sử dụng.
1. Tổng quan Cấu trúc Thư mục (Directory Tree)
Đây là bản đồ toàn diện cho project của bạn:
my-project/
├── cmd/                        # Entry points của ứng dụng
│   └── api/
│       └── main.go             # Hàm main: Khởi tạo config, logger, DB, wire các modules
├── configs/                    # File cấu hình (config.yaml, .env)
├── deployments/                # Dockerfile, docker-compose, k8s charts
├── docs/                       # Swagger, Architecture diagrams
├── internal/                   # Code riêng tư, library bên ngoài không thể import (BẮT BUỘC)
│   ├── common/                 # Các tiện ích dùng chung cho TẤT CẢ modules (Shared Kernel)
│   │   ├── server/             # HTTP Server config (Gin/Echo/Fiber)
│   │   ├── logs/               # Logger (Zap/Logrus)
│   │   ├── errors/             # Custom error types
│   │   ├── events/             # Event Bus interface (cho giao tiếp bất đồng bộ)
│   │   └── database/           # Gorm/Sqlx wrapper
│   │
│   └── modules/                # TRÁI TIM CỦA MODULAR MONOLITH
│       ├── payment/            # Module Payment (Ví dụ)
│       ├── order/              # Module Order (Ví dụ)
│       └── identity/           # Module User/Auth (Ví dụ chi tiết bên dưới)
│           ├── delivery/       # Tầng giao tiếp bên ngoài (HTTP/gRPC)
│           ├── domain/         # Entity, Interface Repository (Core logic)
│           ├── usecase/        # Business Logic (Application layer)
│           ├── infrastructure/ # Implement Repository, External API calls
│           └── module.go       # File định nghĩa Public API của module này
│
├── migrations/                 # Database migrations (SQL files)
├── pkg/                        # Code public có thể share cho dự án khác (ít dùng trong Monolith)
├── go.mod
├── go.sum
└── Makefile                    # Các lệnh build, run, test, migrate
2. Chi tiết cấu trúc bên trong 1 Module (Ví dụ: internal/modules/identity)
Để đảm bảo tính "Modular", mỗi module phải hoạt động như một "mini-service". Các module khác TUYỆT ĐỐI KHÔNG được trọc vào folder usecase hay infrastructure của module này.
Cấu trúc chuẩn của 1 module:
internal/modules/identity/
├── domain/                     # 1. CORE: Không phụ thuộc vào bất kỳ lib nào bên ngoài
│   ├── user.go                 # Struct User (Entity)
│   ├── errors.go               # Lỗi riêng của module này
│   └── repository.go           # Interface định nghĩa việc lưu trữ (UserRepo interface)
│
├── usecase/                    # 2. APPLICATION: Business Logic
│   ├── register_user.go        # Logic đăng ký
│   ├── login_user.go           # Logic đăng nhập
│   └── service.go              # Interface Service
│
├── infrastructure/             # 3. INFRA: Kết nối DB, Redis, 3rd Party
│   ├── postgres/
│   │   └── user_repository.go  # Implement UserRepo interface bằng Gorm/SQL
│   └── jwt/                    # JWT Token generator
│
├── delivery/                   # 4. PRESENTATION: Nhận request
│   ├── http/
│   │   ├── handlers.go         # Xử lý HTTP Request/Response
│   │   └── routes.go           # Đăng ký route (POST /register)
│   └── grpc/                   # Nếu có dùng gRPC
│
└── module.go                   # 5. PUBLIC CONTRACT (Quan trọng nhất)
                                # Nơi khởi tạo module và expose interface cho module khác dùng
3. Quy tắc "Vàng" để giữ cấu trúc đúng 100%
Để không biến Monolith thành "Spaghetti Monolith", bạn phải tuân thủ nghiêm ngặt các quy tắc giao tiếp:
Quy tắc 1: Database Isolation (Cách ly dữ liệu)
Mặc dù dùng chung 1 Database vật lý, nhưng về mặt logic (Schema/Table):
Module Order không được JOIN bảng của Module Identity.
Muốn lấy thông tin User trong Order? Hãy gọi qua API (Go Interface) của module Identity.
Quy tắc 2: Giao tiếp giữa các Module
Các module giao tiếp với nhau như thế nào?
Giao tiếp đồng bộ (Synchronous): Dùng Go Interface.
Ví dụ: Order module cần check user tồn tại.
Order sẽ định nghĩa 1 interface UserChecker.
Main.go sẽ tiêm (inject) Identity service vào Order service.
Giao tiếp bất đồng bộ (Asynchronous): Dùng Event Bus (In-memory hoặc Kafka/RabbitMQ).
Ví dụ: User đăng ký xong -> Bắn event UserRegistered.
Module Email lắng nghe event này và gửi mail welcome.
Quy tắc 3: Dependency Injection (DI)
Tất cả việc kết nối các module diễn ra tại cmd/api/main.go. Các module không tự khởi tạo lẫn nhau bên trong code của chúng.
4. Ví dụ Code Wiring (Tại cmd/api/main.go)
Đây là nơi ma thuật xảy ra để ghép nối các mảnh ghép:
Tại sao cấu trúc này là "Đúng và Đủ"?
Tính mở rộng (Scalability): Khi module Order phình to quá mức, bạn có thể bê nguyên folder internal/modules/order ra làm một Microservice riêng biệt rất dễ dàng vì nó đã tách biệt về Domain và Database logic.
Dễ bảo trì (Maintainability): Sửa logic Auth thì chỉ vào folder identity, không sợ làm hỏng logic của payment.
Teamwork: Team A làm module Payment, Team B làm module Order độc lập, ít conflict code.
Testability: Dễ dàng viết Unit Test vì các tầng giao tiếp qua Interface.
Next Step
Bạn có muốn tôi viết mẫu code chi tiết cho file internal/modules/identity/module.go để xem cách đóng gói (encapsulation) một module chuẩn và cách define interface giao tiếp với bên ngoài không?