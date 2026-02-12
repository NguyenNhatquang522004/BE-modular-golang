### PHẦN 3: CHỈ MỤC CHO MỐI QUAN HỆ (RELATIONSHIP INDEXES) - CỰC QUAN TRỌNG

**Mục đích:** Đây là **Best Practice nâng cao** của Neo4j (từ bản 4.3+). Vì thuật toán Ranking của bạn dựa hoàn toàn vào điểm số trên Cạnh (Edge), nên **BẮT BUỘC** phải đánh index cho Cạnh.

```cypher
// 1. Chỉ mục cho AFFINITY_SCORE (Biến vàng)
// Giúp câu lệnh: ORDER BY r.affinity_score DESC chạy nhanh gấp 100 lần.
CREATE INDEX rel_interacted_affinity_idx IF NOT EXISTS 
FOR ()-[r:INTERACTED_WITH]-() ON (r.affinity_score);

// 2. Chỉ mục cho Time Decay
// Giúp lọc nhanh những tương tác đã quá cũ
CREATE INDEX rel_interacted_time_idx IF NOT EXISTS 
FOR ()-[r:INTERACTED_WITH]-() ON (r.last_interaction_at);

// 3. Chỉ mục cho Sở thích (Interest Score)
// Giúp tìm nhanh User thích Topic nào nhất
CREATE INDEX rel_interested_score_idx IF NOT EXISTS 
FOR ()-[r:INTERESTED_IN]-() ON (r.score);

// 4. Chỉ mục cho Thời gian kết bạn (Friend Since)
// Hiển thị list bạn bè: "Mới kết bạn gần đây"
CREATE INDEX rel_friend_since_idx IF NOT EXISTS 
FOR ()-[r:FRIEND]-() ON (r.since);

// 5. Chỉ mục cho Chống Spam (Used Device)
// Tìm nhanh các User cùng dùng 1 thiết bị trong khoảng thời gian X
CREATE INDEX rel_used_device_time_idx IF NOT EXISTS 
FOR ()-[r:USED_DEVICE]-() ON (r.last_used_at);

```

---

### PHẦN 4: VECTOR INDEX (CHO AI / MACHINE LEARNING)

**Mục đích:** Để thực hiện chức năng tìm kiếm ngữ nghĩa và tìm người giống nhau (Similarity Search) bằng vector embedding 128 chiều.

```cypher
// Tạo Vector Index cho User Embedding
// dimension: 128 (Khớp với model bạn dùng)
// similarityFunction: 'cosine' (Chuẩn nhất cho so sánh hành vi người dùng)
CREATE VECTOR INDEX user_embedding_idx IF NOT EXISTS
FOR (u:User) ON (u.embedding)
OPTIONS {indexConfig: {
 `vector.dimensions`: 128,
 `vector.similarity_function`: 'cosine'
}};

```

---

### PHẦN 5: RÀNG BUỘC TỒN TẠI (EXISTENCE CONSTRAINTS) - (Enterprise Edition Only)

*Lưu ý: Nếu bạn dùng bản Neo4j Enterprise, hãy chạy thêm đoạn này. Nó bắt buộc các trường quan trọng không được phép Null (tránh lỗi dữ liệu rác).*

```cypher
// Bắt buộc User phải có ID và ngày tạo
CREATE CONSTRAINT user_id_exists IF NOT EXISTS FOR (u:User) REQUIRE u.user_id IS NOT NULL;
CREATE CONSTRAINT user_created_exists IF NOT EXISTS FOR (u:User) REQUIRE u.created_at IS NOT NULL;

// Bắt buộc Group phải có ID
CREATE CONSTRAINT group_id_exists IF NOT EXISTS FOR (g:Group) REQUIRE g.group_id IS NOT NULL;

```

### TỔNG KẾT CHIẾN LƯỢC

1. **Chạy Identity Constraints trước tiên:** Để đảm bảo tính nhất quán của dữ liệu (Consistency).
2. **Vector Index là bắt buộc:** Nếu không có nó, tính năng "Tìm người có hành vi giống nhau" sẽ chạy rất chậm vì phải quét toàn bộ Database (Full Scan).
3. **Relationship Index là vũ khí bí mật:** Hầu hết các hướng dẫn cơ bản bỏ qua cái này, nhưng với hệ thống Ranking dựa trên điểm số (`affinity_score`), nó là yếu tố quyết định app của bạn load Newfeed trong 50ms hay 5s.


// Quan hệ này giúp AI hiểu ngữ nghĩa sâu hơn
// VD: Tìm người thích 'Technology' sẽ ra cả người thích 'Golang'
CREATE INDEX rel_topic_child_of_idx IF NOT EXISTS 
FOR ()-[r:CHILD_OF]-() ON (r.weight);

// Truy vết nguồn gốc user
CREATE INDEX rel_invited_timestamp_idx IF NOT EXISTS 
FOR ()-[r:INVITED]-() ON (r.timestamp);

CREATE CONSTRAINT country_code_unique IF NOT EXISTS 
FOR (c:Country) REQUIRE c.code IS UNIQUE;