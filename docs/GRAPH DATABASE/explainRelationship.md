
### 1. NHÓM KẾT NỐI CƠ BẢN (Basic Connections)

#### `FriendRel` (Quan hệ Bạn bè)

* **Ý nghĩa:** A là bạn của B.
* **Nguồn dữ liệu:** Bảng **`friendships`** (Postgres).
* **Mapping dữ liệu:**
* `Since`: Lấy từ cột `created_at` trong bảng `friendships` khi `status = 'accepted'`.
* `Type`: Mặc định là 'normal'. Nếu user thêm vào danh sách "Bạn thân", hệ thống sẽ update thành 'close_friend'.
* `InteractionFrequency`: Tính toán tổng hợp từ bảng `entity_reactions` và `messages` (Cassandra) chia cho thời gian.



#### `FollowRel` (Quan hệ Theo dõi)

* **Ý nghĩa:** A theo dõi B (một chiều).
* **Nguồn dữ liệu:** Bảng **`Followers`** (Postgres).
* **Mapping dữ liệu:**
* `Since`: Lấy từ cột `CreatedAt`.
* `Source`: Lấy từ ngữ cảnh UI khi user bấm nút Follow (VD: 'profile' nếu bấm tại tường nhà, 'suggested' nếu bấm từ mục gợi ý).



#### `InvitedRel` (Quan hệ Mời gọi - Growth)

* **Ý nghĩa:** A mời B tham gia mạng xã hội.
* **Nguồn dữ liệu:** Thường nằm trong logic **Sign Up** (Đăng ký). Dữ liệu có thể lấy từ bảng `users` (nếu bạn thêm cột `referrer_id`) hoặc log hệ thống Referral.
* **Mapping dữ liệu:**
* `Timestamp`: Thời điểm User B tạo tài khoản thành công (`created_at` trong bảng `users`).
* `Code`: Mã giới thiệu hoặc link mà B đã click vào.



---

### 2. NHÓM TƯƠNG TÁC (Ranking Core - Quan trọng nhất)

#### `InteractedWithRel` (Tương tác dài hạn - Giữa User với User)

* **Ý nghĩa:** Tổng hợp mức độ thân thiết giữa A và B.
* **Nguồn dữ liệu:** Tổng hợp từ nhiều nguồn:
1. **Cassandra `entity_reactions**`: Đếm số lần A thả tim bài/comment của B.
2. **MongoDB `Comments**`: Đếm số lần A bình luận vào bài của B.
3. **Cassandra `messages**`: Đếm số tin nhắn A gửi cho B.
4. **Postgres `User_Sessions**`: Đếm số lần A vào xem profile B (`profile_view_count`).


* **Mapping dữ liệu:**
* `LikeCount`, `CommentCount`, `MessageCount`: Cộng dồn từ các nguồn trên.
* `LastInteractionAt`: Lấy timestamp lớn nhất (mới nhất) từ các hành động trên.
* `AffinityScore`: `(Like * 1) + (Comment * 5) + (Message * 10)`.



#### `InteractedRecentlyRel` (Tương tác ngắn hạn - Giữa User với Post)

* **Ý nghĩa:** A vừa tương tác với bài viết P (để gợi ý real-time).
* **Nguồn dữ liệu:**
1. **Cassandra `post_insights**`: Log view/click.
2. **Cassandra `entity_reactions**`: Log thả tim mới nhất.


* **Mapping dữ liệu:**
* `Timestamp`: Thời điểm xảy ra hành động (VD: vừa xem xong).
* `Type`: 'view', 'like', 'share'.
* `Weight`: View = 1.0, Like = 5.0, Share = 10.0 (Dùng để AI biết nên gợi ý mạnh hay nhẹ).



---

### 3. NHÓM SỞ THÍCH & NỘI DUNG (Interest)

#### `InterestedInRel` (User thích Topic)

* **Ý nghĩa:** A quan tâm đến chủ đề T (VD: #Golang).
* **Nguồn dữ liệu:** Phân tích từ **MongoDB `Posts` Collection**.
* Khi User A like bài viết có `hashtags: ["golang", "backend"]`.


* **Mapping dữ liệu:**
* `Score`: Tần suất User tương tác với hashtag đó. (VD: Like 10 bài về Golang -> Score = 0.8).
* `LastEngagedAt`: Thời điểm like bài viết chứa hashtag đó gần nhất.



#### `ChildOfRel` (Phân cấp Topic)

* **Ý nghĩa:** Topic con thuộc Topic cha (VD: #Golang thuộc #Backend).
* **Nguồn dữ liệu:** **Admin Config** (Dữ liệu tĩnh) hoặc import từ Knowledge Graph bên ngoài. Không sinh ra từ hành động user.
* **Mapping dữ liệu:**
* `Weight`: Độ liên quan (Do admin set cứng, VD: 1.0).



#### `MemberOfRel` (Thành viên nhóm)

* **Ý nghĩa:** A là thành viên Group G.
* **Nguồn dữ liệu:** **MongoDB `Group_Members` Collection**.
* **Mapping dữ liệu:**
* `Role`: Lấy từ trường `role` ('admin', 'member').
* `JoinedAt`: Lấy từ trường `joined_at`.



#### `LikesPageRel` (Like Fanpage)

* **Ý nghĩa:** A like/follow Page P.
* **Nguồn dữ liệu:** **MongoDB `Page_Followers` Collection**.
* **Mapping dữ liệu:**
* `Since`: Lấy từ trường `followed_at`.



---

### 4. NHÓM ĐỊNH DANH & BẢO MẬT (Identity & Security)

#### `UsedDeviceRel` (User dùng thiết bị)

* **Ý nghĩa:** A đăng nhập trên thiết bị D.
* **Nguồn dữ liệu:** **Postgres `User_Sessions` Table**.
* **Mapping dữ liệu:**
* `LastUsedAt`: Lấy từ `last_active_at`.
* `LoginCount`: Đếm số dòng (row) trong `User_Sessions` có cùng `device_name` hoặc fingerprint.



#### `HasContactRel` (Danh bạ)

* **Ý nghĩa:** A có lưu số điện thoại P trong danh bạ.
* **Nguồn dữ liệu:** **Mobile App Upload**. (Tính năng "Sync Contacts").
* **Mapping dữ liệu:**
* `UploadedAt`: Thời điểm User bấm nút "Đồng bộ danh bạ".



---

### 5. NHÓM TÍN HIỆU TIÊU CỰC (Negative Signals)

#### `BlockRel`

* **Nguồn:** **Postgres `User_Blocks` Table**.
* **Mapping:** `Since` lấy từ `created_at`.

#### `MuteRel`

* **Nguồn:** **Postgres `Followers` Table** (trường `is_muted` = true).

#### `HiddenPostRel`

* **Nguồn:** **MongoDB `Entity_Edit_Logs**` (hoặc log ẩn bài viết). Khi user chọn "Ẩn bài viết này", hệ thống ghi nhận.
* **Mapping:** `Count` tăng lên mỗi khi user ẩn bài của người đó. `LastHiddenAt` là thời điểm ẩn gần nhất.

#### `ReportedRel`

* **Nguồn:** **MongoDB `Report` Collection** (thường sẽ có collection này cho chức năng Report).
* **Mapping:** `Reason` là lý do report ('spam', 'fake_news').

---

### TỔNG KẾT LUỒNG DỮ LIỆU (PIPELINE)

Để dữ liệu chảy từ SQL/Mongo vào các Struct này trong Neo4j, bạn sẽ cần các **Worker (Golang)**:

1. **Sync Worker:** Lắng nghe sự kiện từ Kafka (VD: `UserCreated`, `FriendRequestAccepted`, `GroupJoined`).
2. **Aggregation Worker:** Chạy định kỳ (Cron) để quét bảng `entity_reactions` (Cassandra), tính toán lại `AffinityScore` và cập nhật vào cạnh `InteractedWithRel`.
3. **Real-time Worker:** Khi có event `PostViewed`, bắn ngay lập tức cạnh `InteractedRecentlyRel` vào Neo4j (có TTL) để gợi ý tức thì.