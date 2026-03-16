

---

### PHẦN 1: CÁC SCORE NẰM TRÊN ĐỈNH (NODE)

#### 1. `PageRankScore` (Trong `UserNode`)

* **Ý nghĩa:** Độ uy tín, tầm ảnh hưởng (Influence) của một người dùng trên toàn mạng lưới. Người có điểm này cao đăng bài sẽ dễ lên xu hướng hơn.
* **Cách tính (Batch Job ban đêm):**
Sử dụng thuật toán PageRank kinh điển thông qua thư viện Graph Data Science (GDS). Nó không đếm số lượng Follower đơn thuần, mà đánh giá "chất lượng" của Follower.
* *Công thức đơn giản hóa:* Điểm của User A bằng tổng điểm của những người có cạnh (`FollowRel`, `FriendRel`, `InteractedWithRel`) trỏ về A, chia cho số liên kết họ có.
* Nếu một KOL (PageRank cao) follow bạn, điểm của bạn sẽ tăng vọt so với việc 100 nick ảo follow bạn.



#### 2. `TrendingScore` (Trong `TopicNode`)

* **Ý nghĩa:** Thể hiện chủ đề này (Hashtag) đang "Nóng" đến mức nào tại thời điểm hiện tại.
* **Cách tính (Real-time Stream):**
Sử dụng hàm suy giảm thời gian (Time-Decay Function), thường dựa trên thuật toán của Hacker News hoặc Reddit.

$$\text{TrendingScore} = \frac{ \sum (\text{Weight từ InteractedRecentlyRel}) }{ (\text{Tuổi của Topic tính bằng giờ} + 2)^G }$$


* *Trong đó:* $G$ là hằng số trọng lực (Ví dụ: $1.8$). Điểm này sẽ tự động tụt rất nhanh theo thời gian nếu không có người tương tác mới, giúp Bảng xếp hạng Trending luôn tươi mới.



---

### PHẦN 2: CÁC SCORE NẰM TRÊN CẠNH (RELATIONSHIP)

#### 3. `ConfidenceScore` (Trong `HasTopicRel` nối Post $\to$ Topic)

* **Ý nghĩa:** Độ tự tin của hệ thống khi dán nhãn bài viết này thuộc về chủ đề kia.
* **Cách tính (Lúc tạo Post):**
* **Tĩnh (Heuristic):** Nếu User tự gõ hashtag `#golang` vào bài viết $\to$ Điểm cố định là $1.0$ (Tin tưởng 100%).
* **AI dự đoán:** Nếu User đăng 1 tấm ảnh, bạn gọi qua Ollama (Vision Model). Model trả về xác suất bài này liên quan đến "ẩm thực" là 85% $\to$ Ghi vào DB là $0.85$.



#### 4. `Weight` (Trong `InteractedRecentlyRel` nối User $\to$ Post)

* **Ý nghĩa:** Trọng số tức thời của một hành động tương tác. Dùng để "bắt trend" ngắn hạn.
* **Cách tính (Gán tĩnh bằng Rule-based Engine):**
Gán cứng dựa trên giá trị hành động đối với nền tảng:
* Lướt qua (Impression) = $0.0$
* Dừng lại xem > 3s (View) = $1.0$
* Like = $3.0$
* Comment = $7.0$
* Share (Rất quan trọng cho Viral) = $15.0$



#### 5. `AffinityScore` (Trong `InteractedWithRel` nối User $\to$ User/Page)

* **Ý nghĩa:** Điểm thân thiết (EdgeRank). Quyết định 80% việc bài viết của User B có xuất hiện trên News Feed của User A hay không.
* **Cách tính (Event-driven với Exponential Decay):**
Tính tổng các tương tác trong quá khứ, nhưng tương tác càng cũ thì giá trị càng giảm.

$$\text{AffinityScore} = \sum_{i} ( W_i \times e^{-\lambda \Delta t_i} )$$


* $W_i$: Trọng số hành động (Chat = 10, Cmt = 7, Like = 3).
* $\Delta t_i$: Thời gian trôi qua từ lúc thực hiện hành động.
* $\lambda$: Tốc độ quên (Decay rate).



#### 6. `Score` (Trong `InterestedInRel` nối User $\to$ Topic)

* **Ý nghĩa:** Mức độ yêu thích của User với một chủ đề (Dùng để gợi ý nội dung từ người lạ).
* **Cách tính (Cộng dồn & Chuẩn hóa):**
Mỗi khi User tương tác với một bài Post, hệ thống tra ngược về Topic của bài Post đó và cộng điểm:
$\text{Điểm cộng thêm} = \text{Weight của Tương tác} \times \text{ConfidenceScore của Bài viết}$
Sau đó, điểm tổng có thể được đưa qua hàm Sigmoid để ép về khoảng từ $0.0$ đến $1.0$.

$$S(x) = \frac{1}{1 + e^{-x}}$$



#### 7. `InteractionFrequency` (Trong `FriendRel`)

* **Ý nghĩa:** Tần suất tương tác trung bình giữa 2 người bạn. Dùng để phân loại "Bạn thân" hay "Bạn xã giao".
* **Cách tính (Batch Job chạy định kỳ):**
Tính bằng: $\frac{\text{Tổng số tương tác 2 chiều}}{\text{Số ngày từ lúc kết bạn (Since)}}$. Nếu chỉ số này lớn hơn một ngưỡng (Threshold) nhất định, hệ thống tự động đổi `Type` của cạnh thành `close_friend`.

#### 8. `Weight` (Trong `ChildOfRel` nối Topic $\to$ Topic)

* **Ý nghĩa:** Độ liên quan giữa 2 chủ đề trong Cây tri thức (Knowledge Graph).
* **Cách tính:**
* **Manual:** Đội ngũ Admin tự set (Ví dụ: `Spring Boot` là con của `Java` với Weight = $0.9$).
* **Machine Learning (Co-occurrence):** Đếm tần suất 2 hashtag xuất hiện cùng nhau trong các bài viết. Tính toán bằng công thức tương đồng Cosine (Cosine Similarity) hoặc Jaccard Index.



---

### Tóm tắt lại luồng chảy của các Score

Khi một cú **Click (Like)** xảy ra, hiệu ứng domino của các Score chạy như sau:

1. Sinh ra `InteractedRecentlyRel` với **Weight** tĩnh $\to$ Đẩy qua Kafka.
2. Cộng dồn vào **TrendingScore** của Topic $\to$ Giúp bài viết lên Xu hướng.
3. Cộng dồn vào **AffinityScore** của tác giả $\to$ Giúp 2 người xích lại gần nhau hơn trên Feed.
4. Cộng dồn vào **Score** của Topic $\to$ Hệ thống học được sở thích của người Like.

Bạn đã nắm rõ bức tranh toàn cảnh về cách các con số này "nhảy múa" rồi chứ? Bạn có muốn tôi viết thử một câu lệnh truy vấn **Cypher (Neo4j)** kết hợp `AffinityScore` và `TrendingScore` để lấy ra danh sách 10 bài viết cho Bảng tin của một user không?





Trong 8 loại Score (và thêm 1 trường đặc biệt) của bạn, chúng ta chia làm 2 nhóm rõ rệt:

---

### NHÓM 1: TÍNH BẰNG THỦ CÔNG & TOÁN HỌC (RULE-BASED / MATH)

*Nhóm này hệ thống tự tính toán bằng code Backend (Golang) hoặc chạy thuật toán trực tiếp trên Graph DB. Rất nhanh, Real-time và không cần Model AI.*

**1. `Weight` (Trọng số tương tác tức thời)**

* **Cách tính:** Thủ công (Hard-code bằng lệnh `if-else` hoặc file config).
* **Lý do:** View = 1, Like = 3, Share = 10. Đây là logic kinh doanh (Business Logic) do bạn tự quyết định để điều hướng hành vi người dùng, không cần AI phải học.

**2. `TrendingScore` (Độ hot của Topic)**

* **Cách tính:** Bằng công thức Toán học (Hàm suy giảm thời gian - Time-decay function).
* **Lý do:** Chỉ là phép tính cộng trừ nhân chia dựa trên tổng lượng Like/Share và số giờ trôi qua. Code Backend tự tính mỗi khi có Event.

**3. `AffinityScore` (Điểm thân thiết giữa 2 User)**

* **Cách tính:** Bằng công thức Toán học (Exponential Decay).
* **Lý do:** Backend hoặc Graph DB sẽ tự động cộng dồn điểm tương tác trong quá khứ và chia cho thời gian để làm điểm này "nguội" đi nếu lâu ngày không tương tác.

**4. `InteractionFrequency` (Tần suất tương tác)**

* **Cách tính:** Toán học (Tổng số tương tác chia cho Số ngày làm bạn).
* **Lý do:** Phép chia cơ bản, chạy bằng Job ban đêm.

**5. `PageRankScore` (Độ uy tín của User)**

* **Cách tính:** Thuật toán Đồ thị (Graph Algorithm).
* **Lý do:** Chạy thuật toán PageRank của Neo4j GDS (Graph Data Science). Nó là Toán học ma trận, không phải là Machine Learning dự đoán.

**6. `Score` (Điểm sở thích của User với Topic)**

* **Cách tính:** Cộng dồn tuyến tính (Tương tác bài nào thì cộng điểm bài đó vào Topic tương ứng).
2. Cách hệ thống xào nấu (Công thức tính toán)Hệ thống tính toán qua 2 bước:Bước 1: Tính $\alpha$ (Lực tác động của riêng lần tương tác này)$$\alpha = \text{actionWeight} \times \text{confidence\_score} \times \text{learningRate}$$Bước 2: Tính Điểm Sở Thích Mới ($S_{new}$)Dùng công thức tịnh tiến (Asymptotic) để cộng dồn điểm mới vào điểm cũ:$$S_{new} = S_{old} + \alpha \times (1.0 - S_{old})$$Cái hay của cụm $(1.0 - S_{old})$ là: Khi điểm của bạn càng gần 1.0 (nghĩa là bạn đã quá thích Topic này rồi), thì khoảng cách $(1.0 - S_{old})$ càng nhỏ lại. Do đó, điểm cộng thêm vào sẽ lắt nhắt từng chút một, rất khó để chạm đỉnh tuyệt đối.
---

### NHÓM 2: TÍNH BẰNG AI / MACHINE LEARNING

*Nhóm này bắt buộc phải dùng AI để "hiểu" dữ liệu phi cấu trúc (Văn bản, Hình ảnh, Hành vi).*

**1. `ConfidenceScore` (Độ tự tin khi dán nhãn Topic cho Post)**

* **Dữ liệu đầu vào (Input):** Nội dung Text của bài viết, Các hình ảnh đính kèm (dạng Base64), Video (nếu có).
* **Loại AI sử dụng:**
* **Multimodal LLM (Như bạn đang code):** LLaVA, Qwen-VL chạy qua Ollama để vừa đọc chữ vừa nhìn ảnh rồi suy luận ra Topic.
* **NLP Text Classification:** Nếu bài viết chỉ có chữ, dùng các model phân loại văn bản nhẹ nhàng hơn như PhoBERT (cho tiếng Việt) để đoán chủ đề.



**2. Trường `Embedding []float32` (Vector 128D trong UserNode)**
*(Bạn có định nghĩa trường này từ đầu, đây là thứ 100% sinh ra từ AI và là vũ khí mạnh nhất của TikTok).*

* **Dữ liệu đầu vào (Input):** Chuỗi lịch sử hành vi của User (Ví dụ: 100 bài Post gần nhất họ xem > 5s, ID của những người họ hay nhắn tin, các Topic họ Like).
* **Loại AI sử dụng:**
* **Graph Neural Networks (GNN):** Các thuật toán như GraphSAGE hoặc Node2Vec.
* **Mục đích:** AI sẽ nén toàn bộ "Sở thích và Tính cách" của User A thành 128 con số (Vector). Sau đó, bạn tính khoảng cách (Cosine Similarity) giữa Vector của User A và User B. Nếu 2 vector giống nhau, AI đoán họ cùng chung sở thích $\to$ Gợi ý kết bạn hoặc gợi ý chung một bài Post.



**3. `Weight` (Độ liên quan giữa Topic cha - con trong `ChildOfRel`)**

* **Dữ liệu đầu vào (Input):** Tần suất xuất hiện cùng nhau của các Hashtag/Topic trên toàn bộ hàng triệu bài Post của hệ thống (Co-occurrence Matrix).
* **Loại AI sử dụng:**
* **Word2Vec / NLP Embeddings:** Model AI quét qua văn bản để hiểu ngữ nghĩa. Nó sẽ tự học được rằng chữ `Golang` và chữ `Backend` thường đi chung với nhau trong một ngữ cảnh $\to$ AI tự động gán Weight = $0.92$ cho đường nối giữa 2 Topic này. Nhờ đó bạn không phải dùng sức người ngồi nối tay hàng ngàn Topic với nhau.



### Lời khuyên khi xây dựng hệ thống:

Trong giai đoạn đầu (Phase 1) của dự án:

1. Bạn hãy **tập trung code thật tốt Nhóm 1 (Nhóm Toán học/Thủ công)** vì nó dễ làm, độ chính xác 100%, và là khung xương để hệ thống chạy được ngay.
2. Đối với Nhóm 2, hiện tại bạn đang làm cái **số 1 (Dùng Ollama trích xuất Topic)** là rất đúng đắn để giải quyết bài toán "bắt ép nội dung vào khuôn khổ Sở thích". Còn Vector Embedding hay Graph ML thì để dành khi hệ thống có trên 100,000 users thì mới có đủ Data để AI học.
