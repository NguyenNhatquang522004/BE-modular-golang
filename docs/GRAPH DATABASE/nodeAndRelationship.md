# NEO4J GRAPH SCHEMA DESIGN (FINAL - 100% HYPERSCALE READY)

Phiên bản này được thiết kế dựa trên kiến trúc của Facebook (TAO), Twitter (GraphJet) và Pinterest (Pixie), tối ưu hóa cho:
1.  **Ranking (Sắp xếp Newfeed)** dựa trên điểm số tương tác và thời gian.
2.  **Recommendation (Gợi ý)** sử dụng AI Vector Embedding và Graph Traversal.
3.  **Security (An toàn)** với cơ chế chống Spam/Bot Farm bằng Device Fingerprint và Invite Graph.
4.  **Growth (Tăng trưởng)** với khả năng truy vết nguồn gốc người dùng (Provenance).
5.  **Real-time (Thời gian thực)** với Short-term Context (Viral Posts).

---

## PHẦN 1: CÁC LOẠI NODE (NODES) & BIẾN (PROPERTIES)

### 1. Node `User` (Trái tim của Graph)
* `user_id` (String/UUID - Unique): ID đồng bộ từ Postgres.
* `created_at` (Long/Timestamp): Thời gian tạo tài khoản.
* `last_active_at` (Long/Timestamp): Lọc User "ma" (inactive).
* `is_verified` (Boolean): Tích xanh (Boost độ uy tín).
* `embedding` (List<Float>): Vector 128D biểu diễn hành vi người dùng cho GNN.
* `risk_score` (Float): Điểm rủi ro (0.0 - 1.0). Cao = Bot/Spammer.
* `page_rank_score` (Float): [MỚI] Điểm uy tín trong mạng lưới (Tính bằng GDS PageRank).
* `community_id` (Integer): [MỚI] ID cụm cộng đồng (Tính bằng GDS Louvain).

### 2. Node `Topic` (Interest Graph - Có phân cấp)
* `name` (String - Unique): Tên chủ đề (VD: "golang", "backend").
* `trending_score` (Float): Điểm xu hướng hiện tại (Update mỗi giờ).
* *(Lưu ý: `category` cũ đã được chuyển thành quan hệ `CHILD_OF` để hỗ trợ suy luận bắc cầu).*

### 3. Node `Group` (Community)
* `group_id` (String - Unique).
* `privacy` (String): 'public', 'closed', 'secret'.
* `member_count` (Integer): Số thành viên.

### 4. Node `Page` (Brand/Fanpage)
* `page_id` (String - Unique).
* `category_id` (String).
* `rating` (Float): Điểm đánh giá.

### 5. Node `Post` (Short-term / Viral Content - [MỚI])
* *Chỉ lưu các bài đang Trending hoặc Viral trong 24-48h.*
* `post_id` (String - Unique).
* `created_at` (Long).
* `ttl` (Long): Thời điểm tự hủy (Time-To-Live).

### 6. Node `Device` (Anti-Spam)
* `device_id` (String - Unique): Hash Fingerprint.
* `trust_level` (Float): Độ tin cậy thiết bị.

### 7. Node `PhoneContact` (Identity)
* `phone_hash` (String - Unique): Hash SHA256 số điện thoại.

### 8. Node `Location` (Geo - Có phân cấp)
* `city_id` (String).
* `geo_hash` (String).
* `country_code` (String).

---

## PHẦN 2: CÁC MỐI QUAN HỆ (RELATIONSHIPS) & BIẾN (PROPERTIES)

### 1. Nhóm Xã hội & Tăng trưởng (Social & Growth Graph)
* **`(:User)-[:FRIEND]->(:User)`**
    * `since` (Long).
    * `type` (String): 'normal', 'close_friend', 'family'.
    * `interaction_frequency` (Float): [Pre-calc] Tần suất tương tác.
* **`(:User)-[:FOLLOWS]->(:User)`**
    * `since` (Long).
    * `source` (String): 'profile', 'search', 'suggested'.
* **`(:User)-[:INVITED]->(:User)` [MỚI - Growth]**
    * `timestamp` (Long).
    * `code` (String): Mã mời/Method.
    * *Tác dụng: Truy vết nguồn gốc (User Provenance) để diệt cụm Bot.*

### 2. Nhóm Tương tác (Interaction Graph - Ranking Core)
* **`(:User)-[:INTERACTED_WITH]->(:User)` [Long-term Memory]**
    * *Cạnh tổng hợp (Aggregated Edge).*
    * `last_interaction_at` (Long): Tính Time Decay.
    * `like_count` (Integer).
    * `comment_count` (Integer).
    * `message_count` (Integer).
    * `share_count` (Integer).
    * `profile_view_count` (Integer).
    * `affinity_score` (Float): Điểm thân thiết (Calculated by Worker).
* **`(:User)-[:INTERACTED_RECENTLY]->(:Post)` [MỚI - Short-term Memory]**
    * `timestamp` (Long).
    * `type` (String): 'view', 'like', 'share'.
    * `weight` (Float): Trọng số tức thời.
    * *Tác dụng: Real-time Recommendation (TikTok style).*

### 3. Nhóm Sở thích & Nội dung (Interest Graph)
* **`(:User)-[:INTERESTED_IN]->(:Topic)`**
    * `score` (Float): 0.0 - 1.0 (Độ thích).
    * `last_engaged_at` (Long).
* **`(:Topic)-[:CHILD_OF]->(:Topic)` [MỚI - Hierarchy]**
    * `weight` (Float): Độ mạnh quan hệ cha-con.
    * *Tác dụng: Suy luận bắc cầu (Thích "Golang" -> Thích "Backend").*
* **`(:User)-[:MEMBER_OF]->(:Group)`**
    * `role` (String): 'admin', 'member'.
    * `joined_at` (Long).
* **`(:User)-[:LIKES_PAGE]->(:Page)`**
    * `since` (Long).

### 4. Nhóm Tín hiệu Tiêu cực (Negative Signals - Safety)
* **`(:User)-[:BLOCKS]->(:User)`**
    * `since` (Long).
* **`(:User)-[:MUTES]->(:User)`**
    * `since` (Long).
* **`(:User)-[:HIDDEN_POST_FROM]->(:User)`**
    * `count` (Integer).
    * `last_hidden_at` (Long).
* **`(:User)-[:REPORTED]->(:User)`**
    * `reason` (String).
    * `timestamp` (Long).

### 5. Nhóm Định danh & Thiết bị (Identity Graph)
* **`(:User)-[:USED_DEVICE]->(:Device)`**
    * `last_used_at` (Long).
    * `login_count` (Integer).
* **`(:User)-[:HAS_CONTACT]->(:PhoneContact)`**
    * `uploaded_at` (Long).

---

## PHẦN 3: SCHEMA CONSTRAINTS & INDEXES (CYPHER SCRIPT)

Copy và chạy toàn bộ script này để khởi tạo Database chuẩn.

```cypher
// ==========================================
// 1. UNIQUENESS CONSTRAINTS (Định danh)
// ==========================================
CREATE CONSTRAINT user_id_unique IF NOT EXISTS FOR (u:User) REQUIRE u.user_id IS UNIQUE;
CREATE CONSTRAINT topic_name_unique IF NOT EXISTS FOR (t:Topic) REQUIRE t.name IS UNIQUE;
CREATE CONSTRAINT group_id_unique IF NOT EXISTS FOR (g:Group) REQUIRE g.group_id IS UNIQUE;
CREATE CONSTRAINT page_id_unique IF NOT EXISTS FOR (p:Page) REQUIRE p.page_id IS UNIQUE;
CREATE CONSTRAINT post_id_unique IF NOT EXISTS FOR (p:Post) REQUIRE p.post_id IS UNIQUE; // [MỚI]
CREATE CONSTRAINT device_id_unique IF NOT EXISTS FOR (d:Device) REQUIRE d.device_id IS UNIQUE;
CREATE CONSTRAINT phone_hash_unique IF NOT EXISTS FOR (c:PhoneContact) REQUIRE c.phone_hash IS UNIQUE;

// ==========================================
// 2. NODE INDEXES (Hiệu năng tìm kiếm)
// ==========================================
// User Metadata
CREATE INDEX user_verified_idx IF NOT EXISTS FOR (u:User) ON (u.is_verified);
CREATE INDEX user_created_at_idx IF NOT EXISTS FOR (u:User) ON (u.created_at);
CREATE INDEX user_last_active_idx IF NOT EXISTS FOR (u:User) ON (u.last_active_at);
CREATE INDEX user_risk_score_idx IF NOT EXISTS FOR (u:User) ON (u.risk_score);
CREATE INDEX user_pagerank_idx IF NOT EXISTS FOR (u:User) ON (u.page_rank_score); // [MỚI]

// Content & Topics
CREATE INDEX topic_trending_idx IF NOT EXISTS FOR (t:Topic) ON (t.trending_score);
CREATE INDEX post_ttl_idx IF NOT EXISTS FOR (p:Post) ON (p.ttl); // [MỚI - Để xóa bài cũ]

// Group & Geo
CREATE INDEX group_member_count_idx IF NOT EXISTS FOR (g:Group) ON (g.member_count);
CREATE INDEX location_geohash_idx IF NOT EXISTS FOR (l:Location) ON (l.geo_hash);

// ==========================================
// 3. RELATIONSHIP INDEXES (Ranking & Traversal)
// ==========================================
// Interaction & Ranking (Quan trọng nhất)
CREATE INDEX rel_interacted_affinity_idx IF NOT EXISTS FOR ()-[r:INTERACTED_WITH]-() ON (r.affinity_score);
CREATE INDEX rel_interacted_time_idx IF NOT EXISTS FOR ()-[r:INTERACTED_WITH]-() ON (r.last_interaction_at);
CREATE INDEX rel_interacted_recent_time_idx IF NOT EXISTS FOR ()-[r:INTERACTED_RECENTLY]-() ON (r.timestamp); // [MỚI]

// Interest & Growth
CREATE INDEX rel_interested_score_idx IF NOT EXISTS FOR ()-[r:INTERESTED_IN]-() ON (r.score);
CREATE INDEX rel_topic_child_of_weight_idx IF NOT EXISTS FOR ()-[r:CHILD_OF]-() ON (r.weight); // [MỚI]
CREATE INDEX rel_invited_timestamp_idx IF NOT EXISTS FOR ()-[r:INVITED]-() ON (r.timestamp); // [MỚI]

// Social & Security
CREATE INDEX rel_friend_since_idx IF NOT EXISTS FOR ()-[r:FRIEND]-() ON (r.since);
CREATE INDEX rel_used_device_time_idx IF NOT EXISTS FOR ()-[r:USED_DEVICE]-() ON (r.last_used_at);

// ==========================================
// 4. VECTOR INDEX (AI Semantic Search)
// ==========================================
CREATE VECTOR INDEX user_embedding_idx IF NOT EXISTS
FOR (u:User) ON (u.embedding)
OPTIONS {indexConfig: {
 `vector.dimensions`: 128,
 `vector.similarity_function`: 'cosine'
}};

// ==========================================
// 5. EXISTENCE CONSTRAINTS (Enterprise Only)
// ==========================================
// CREATE CONSTRAINT user_id_exists IF NOT EXISTS FOR (u:User) REQUIRE u.user_id IS NOT NULL;
// CREATE CONSTRAINT user_created_exists IF NOT EXISTS FOR (u:User) REQUIRE u.created_at IS NOT NULL;