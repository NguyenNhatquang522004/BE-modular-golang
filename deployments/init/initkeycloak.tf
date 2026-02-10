# deployments/init/initkeycloak.tf

terraform {
  required_providers {
    keycloak = {
      source  = "mrparkers/keycloak"
      version = ">= 4.4.0" # Nên dùng bản mới nhất
    }
  }
}

# 1. Cấu hình Provider
provider "keycloak" {
  client_id     = "admin-cli"
  username      = "admin"
  password      = "admin"
  url           = "http://keycloak_BE-modular-golang:8080" # URL nội bộ Docker
}

# 2. Tạo Realm
resource "keycloak_realm" "realm" {
  realm        = "SocialNetworkRealm"
  enabled      = true
  display_name = "Social Network"
}

# 3. Tạo Role User (Viết thường để khớp với Go Enum)
resource "keycloak_role" "user_role" {
  realm_id    = keycloak_realm.realm.id
  name        = "user" # Khớp với RoleTypeUser.String()
  description = "Role mặc định cho người dùng"
}

# 4. Tạo Role Admin (Viết thường để khớp với Go Enum)
resource "keycloak_role" "admin_role" {
  realm_id    = keycloak_realm.realm.id
  name        = "admin" # Khớp với RoleTypeAdmin.String()
  description = "Role quản trị hệ thống"
}

# 5. Set "user" làm Default Role
# Lưu ý: Resource này sẽ định nghĩa lại toàn bộ default roles.
resource "keycloak_default_roles" "default_roles" {
  realm_id = keycloak_realm.realm.id
  default_roles = [
    "offline_access",
    "uma_authorization",
    keycloak_role.user_role.name # Tự động gán role 'user' khi đăng ký
  ]
}

# 6. Tạo Client cho Backend Golang
resource "keycloak_openid_client" "backend_client" {
  realm_id                 = keycloak_realm.realm.id
  client_id                = "backend-api"
  name                     = "Backend API Service"
  enabled                  = true
  access_type              = "CONFIDENTIAL" # Bắt buộc CONFIDENTIAL để có Service Account
  standard_flow_enabled    = false          # Backend không login bằng trình duyệt
  service_accounts_enabled = true           # Cho phép lấy token dạng Client Credentials
}

# =============================================================================
# 7. QUAN TRỌNG: CẤP QUYỀN "ADMIN" CHO BACKEND CLIENT
# Để code Go có thể: Lấy danh sách role, Gán role cho user, Xóa role...
# =============================================================================

# 7.1. Lấy thông tin role "realm-admin" có sẵn của Keycloak (thuộc client realm-management)
data "keycloak_openid_client" "realm_management" {
  realm_id  = keycloak_realm.realm.id
  client_id = "realm-management"
}

data "keycloak_role" "realm_admin" {
  realm_id  = keycloak_realm.realm.id
  client_id = data.keycloak_openid_client.realm_management.id
  name      = "realm-admin" # Role siêu quyền lực quản lý realm
}

# 7.2. Gán role "realm-admin" vào Service Account của "backend-api"
resource "keycloak_openid_client_service_account_role" "backend_service_account_role" {
  realm_id                = keycloak_realm.realm.id
  service_account_user_id = keycloak_openid_client.backend_client.service_account_user_id
  client_id               = data.keycloak_openid_client.realm_management.id
  role                    = data.keycloak_role.realm_admin.name
}

# =============================================================================
# 8. OUTPUTS (In ra màn hình console sau khi chạy xong)
# =============================================================================

output "backend_client_id" {
  value = keycloak_openid_client.backend_client.client_id
}

output "backend_client_secret" {
  value     = keycloak_openid_client.backend_client.client_secret
  sensitive = true # Terraform sẽ che lại, muốn xem chạy lệnh: terraform output -raw backend_client_secret
}