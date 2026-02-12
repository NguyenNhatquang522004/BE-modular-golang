--USER 
CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    phone_number VARCHAR(20) UNIQUE, -- Thêm để login bằng SĐT
    password_hash VARCHAR(255) NOT NULL,
    username VARCHAR(50) UNIQUE, -- Quan trọng để tag tên (@hieu)
    type_login VARCHAR(20), -- 'email', 'google', 'facebook', 'apple'
    is_active BOOLEAN DEFAULT TRUE,
    settings_active_status BOOLEAN DEFAULT TRUE, -- Trạng thái online tổng
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);


Profiles Collection
{
  "_id": ObjectId,
  "user_id": UUID, // Index Unique - Khóa chính logic của hệ thống
  
  // 1. Thông tin định danh
  "first_name": String,
  "last_name": String,
  "full_name": String, // (Denormalization) Lưu gộp "Họ + Tên"f để search nhanh hơn
  "slug": String, // (VD: nguyen-van-a) Dùng cho SEO URL
  
  "bio": String,
  "date_of_birth": Date,
  "gender": String, // ENUM('male', 'female', 'other', 'hidden')
  // 2. LIÊN KẾT VỚI GRIDFS (Avatar & Cover)
  // Thay vì chỉ string url, ta dùng object để quản lý file chặt chẽ
  "avatar": {
    "url": String,       
    "updated_at": Date   // Để browser cache busting (avatar.jpg?t=123456)
  },
  
  "cover_photo": {       // Ảnh bìa trang cá nhân
    "url": String,
    "position_y": Number // Căn chỉnh vị trí ảnh bìa (0-100%)
  },

  // 3. Thông tin liên hệ & Địa lý
  "address": {
    "street": String,
    "city": String,
    "country": String,
    "coordinates": [Number, Number] // [Long, Lat] GeoJSON nếu cần map
  },
  "phone_number": String,
  
  "social_links": {
    "facebook": String,
    "twitter": String,
    "linkedin": String,
    "instagram": String,
    "github": String,
    "website": String
  },

  // 4. Hồ sơ năng lực (Professional Info)
  // Có thể lưu CV dạng file PDF từ GridFS
  "cv_document": {
      "file_id": ObjectId, // Trỏ tới fs_files._id (file PDF/Docx)
      "filename": String,  // Tên file gốc (VD: CV_NguyenVanA.pdf)
      "uploaded_at": Date
  },

  "work_experience": [
    {
      "_id": ObjectId, // Để dễ dàng edit/delete từng dòng
      "company": String,
      "position": String,
      "start_date": Date,
      "end_date": Date, // Null nếu đang làm việc
      "is_current": Boolean,
      "description": String
    }
  ],

  "education": [
    {
      "_id": ObjectId,
      "institution": String, // Trường học
      "degree": String,      // Bằng cấp
      "field_of_study": String, // Chuyên ngành
      "start_date": Date,
      "end_date": Date
    }
  ],

  // 5. System Meta
  "settings": {
      "is_private": Boolean, // Hồ sơ công khai hay riêng tư
      "allow_search_engine": Boolean // Cho phép Google index không?
  },
  "notifications": {
      "email_frequency": "weekly",
      "push_types": ["comment", "friend_request"] // Các loại nhận
  }
  "created_at": Date,
  "updated_at": Date,
  "deleted_at": Date // Soft delete
}
-- friendship
CREATE TABLE friendships (
    friendship_id UUID PRIMARY KEY,
    requester_id UUID NOT NULL REFERENCES users(user_id),
    recipient_id UUID NOT NULL REFERENCES users(user_id),
    status VARCHAR(20), -- 'pending', 'accepted', 'blocked'
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(requester_id, recipient_id) -- Tránh duplicate request
);
-- followers
Create table Followers (
follower_id

is_muted -- ENUM('yes', 'no')
-- NGƯỜI ĐI THEO DÕI (Fan / Follower) - Khóa ngoại
follower_user_id UUID NOT NULL REFERENCES users(user_id),
-- NGƯỜI ĐƯỢC THEO DÕI (Idol / Followee) - Khóa ngoại
followed_user_id UUID NOT NULL REFERENCES users(user_id),
CreatedAt
UpdatedAt
DeletedAt
PRIMARY KEY (follower_id, followed_id)
)
-- user blocks
CREATE TABLE User_Blocks (
    block_id UUID PRIMARY KEY,
    blocker_user_id UUID NOT NULL REFERENCES users(user_id), -- Người chặn - Khóa ngoại
    blocked_user_id UUID NOT NULL REFERENCES users(user_id),  -- Người bị chặn - Khóa ngoại
    reason TEXT,                    -- Lý do chặn (tùy chọn)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP               -- Xoá mềm
    type VARCHAR(20) DEFAULT 'full', -- 'full' (chặn hoàn toàn), 'partial' (chặn một phần)
    PRIMARY KEY (blocker_id, blocked_id)
)


-- chức năng  Bộ sưu tập ảnh/video: Album ảnh, Ảnh có mặt bạn, Video đã tải lên.


-- chức năng  bài  POST
--Mongodb media
Posts Collection {
 "_id": ObjectId,
  "user_id": ObjectId, // Index
  
  // 1. Phân loại & Ngữ cảnh
  "type": String, // 'text', 'media', 'share', 'background', 'qna', 'live'
  "context": {
      "type": String, // 'user_wall', 'group', 'page'
      "target_id": ObjectId // ID của Group/Page. Null nếu là tường nhà.
  },

  // 2. Nội dung & Tóm tắt (Quan trọng để render nhanh)
  "content": String, 
  "slug": String,
  "summary": {
      "feeling_icon": String, // "😄"
      "feeling_name": String, // "đang cảm thấy hạnh phúc"
      "location_name": String, // "tại Starbucks"
      "has_media": Boolean,
      "media_count": Number, // "5 ảnh"
      "thumbnail_url": String, // URL ảnh đầu tiên (để hiện preview)
      "background_theme_id": String // ID màu nền (nếu có)
  },

  // 3. Logic hiển thị
  "privacy": {
      "scope": String, // 'public', 'friends', 'only_me', 'custom'
      "allow_comment": Boolean,
      "allow_share": Boolean
  },
  "status": String, // 'published', 'draft', 'archived', 'hidden'
  "is_pinned": Boolean,
  "is_edited": Boolean, 

  // 4. Counters
  "stats": {
      "total_reactions": Number, 
      "comments": Number,
      "shares": Number,
      "views": Number,
      "top_reaction_types": [String] // ['haha', 'love'] - Cache 2 icon nhiều nhất
  },

  // 5. Search & Tagging
  "hashtags": [String], // Index
  "mentions": [ObjectId], // Index

  // 6. Timestamps
  "published_at": Date, // Index (Sort Newfeed)
  "created_at": Date,
  "updated_at": Date,
  "deleted_at": Date
}
Post_Extensions collection{
"_id": ObjectId,
  "post_id": ObjectId, // Index Unique
  
  // A. Nếu là bài Share
  "share_data": {
      "original_post_id": ObjectId, // Bài gốc
      "parent_post_id": ObjectId,   // Bài cha (Direct parent)
      "snapshot": { // Cache hiển thị
          "author_id": ObjectId,
          "author_name": String,
          "author_avatar": String,
          "content_excerpt": String,
          "media_thumb": String,
          "created_at": Date
      }
  },

  // B. Nếu là bài Background
  "background_data": {
      "theme_id": String, 
      "text_color": String 
  },

  // C. Nếu là bài Q&A
  "qna_data": {
      "question": String,
      "button_text": String
  },

  // D. Chi tiết Feeling/Activity
  "activity_data": {
      "type": String, // 'watching', 'traveling'
      "object_id": String, // ID phim/sách
      "object_name": String // "Phim Mai"
  },
  
  // E. Chi tiết Location (GeoJSON)
  "location_detail": {
      "type": { type: String, default: "Point" },
      "coordinates": [Number, Number], // [Long, Lat]
      "address": String,
      "map_url": String
  }
}

Post_Media collection{
"_id": ObjectId,
  "post_id": ObjectId, // Index Unique
  
  "items": [
    {
       "_id": ObjectId, 
       "media_type": String, // 'image', 'video', 'gif'
       "url": String, // Full quality URL (SeaweedFS)
       "thumbnail_url": String, // Low quality URL
       
       "metadata": { 
           "width": Number, 
           "height": Number, 
           "duration": Number, 
           "size_bytes": Number,
           "mime_type": String
       },
       
       "order": Number, // Thứ tự hiển thị 1, 2, 3
       
       // Tag mặt người (Social Feature)
       "tagged_users": [ 
           { 
               "user_id": ObjectId, 
               "name": String, // Cache tên
               "x": Number, // Tọa độ 0.0 - 1.0
               "y": Number 
           } 
       ]
    }
  ]
}
Post_Settings collection {
"_id": ObjectId,
  "post_id": ObjectId, // Index Unique
  
  // 1. Quản lý đăng bài
  "schedule": {
      "is_scheduled": Boolean,
      "publish_time": Date,
      "publisher_user_id": UUID, // Admin bấm đăng
      "author_role_snapshot": String // 'editor'
  },

  // 2. Quảng cáo (Liên kết Postgres)
  "ads_info": {
      "campaign_id": ObjectId, // Link sang bảng Ads
      "is_sponsored": Boolean,
      "cta_link": String // "Learn More" link
  },

  // 3. Target Audience (Ai nhìn thấy bài này?)
  "targeting": {
      "locations": [String],
      "age_min": Number,
      "age_max": Number,
      "genders": [String],
      "languages": [String],
      "interests": [String]
  }
}
Entity_Edit_Logs collection{
 "_id": ObjectId,
  "target_collection": String, // 'posts' hoặc 'comments'
  "target_id": ObjectId, // ID đối tượng
  
  "version": Number, 
  "editor_id": ObjectId, 
  "edited_at": Date,
  
  "diff": {
      "old_content": String,
      "new_content": String,
      "changed_fields": [String] // ['content', 'privacy', 'media']
  },
  
  "ip_address": String, 
  "user_agent": String
}
-- chức năng  Bộ sưu tập ảnh/video: Album ảnh, Ảnh có mặt bạn, Video đã tải lên.
Albums Collection
{
  "_id": ObjectId,
  "user_id": ObjectId, // Reference User
  "group_id": ObjectId, // Null nếu là album cá nhân
  // Thông tin cơ bản
  "title": String, 
  "description": String,
  "type": String, // ENUM('normal', 'profile', 'cover', 'mobile', 'trash') - Để phân loại album hệ thống
  
  // Ảnh bìa album (Quan trọng để hiển thị gallery)
  "cover_asset_id": ObjectId, // Reference tới collection MediaAssets
  
  // Denormalization (Lưu thừa để hiển thị nhanh UI mà không cần count)
  "asset_count": Number, // Tổng số ảnh/video
  "last_asset_added_at": Date, // Để sort album nào mới cập nhật lên đầu
  
  // Tương tác trên chính Album (Facebook cho like/comment cả Album)
  "reactions": {
      "total": Number,
      "like": Number,
      "love": Number
      // ... các icon khác
  },
  "comment_count": Number,

  // Cài đặt
  "privacy": {
    "level": String, // ENUM('public', 'friends', 'only_me', 'custom')
    "allow_list": [ObjectId],
    "block_list": [ObjectId]
  },
  
  "created_at": Date,
  "updated_at": Date,
  "deleted_at": Date // Soft delete
}
User_MediaAssets Collection
{
  "_id": ObjectId,
  "user_id": ObjectId, // Người upload
  "album_id": ObjectId, // Reference Albums. (Index trường này để load ảnh trong album nhanh)
  "post_id": ObjectId, // Reference user_Post
  "group_id": ObjectId, // Quan trọng: Để load tab "Media" trong nhóm
  // --- KẾT NỐI GRIDFS (Cốt lõi) ---
  "url":  "" , // ID file trong SeaweedFS
  
  // URLs (Cần thiết kế 2 loại link)
  "original_url": String, // Link stream file gốc (gọi qua API GridFS)
  "thumbnail_url": String, // Link ảnh nhỏ (Thumb) để load list nhanh (thường lưu ở server static hoặc CDN, không dùng GridFS cho thumb)
  
  "asset_type": String, // ENUM('image', 'video' , 'gif', 'sticker') - Loại media
  
  // Metadata hiển thị (Tách ra để render UI khung ảnh trước khi load nội dung)
  "metadata": {
    "width": Number, 
    "height": Number,
    "duration": Number, // Nếu là video (seconds)
    "size_bytes": Number,
    "mime_type": String // 'image/jpeg', 'video/mp4'
  },

  // Tính năng xã hội (Social Features)
  "caption": String, // Mô tả riêng cho từng ảnh
  "hashtags": [String], // Index trường này
  
  "tagged_users": [
    {
      "user_id": ObjectId,
      "name": String, // Lưu thừa tên để hiển thị nhanh
      "position": { "x": Number, "y": Number }, // Tọa độ mặt (0.0 đến 1.0)
      "status": String // ENUM('pending', 'approved')
    }
  ],
  
  // Tương tác (Reaction/Comment cho từng ảnh)
  "reactions_count": { "total": Number, "like": Number, "love": Number },
  "comment_count": Number,

  // Cài đặt
  "order": Number, // Thứ tự sắp xếp trong album (nếu user muốn custom)
  "privacy": { // Thường sẽ inherit (kế thừa) từ Album, nhưng có thể override
      "level": String, 
      "inherit_from_album": Boolean // Mặc định là true
  },

  "created_at": Date,
  "updated_at": Date,
  "deleted_at": Date
}
-- chức năng comment cho post 
Comments collection{
  "_id": ObjectId,
  "post_id": ObjectId, // Reference bài viết gốc
  "user_id": ObjectId, // Người comment
  
  // 1. NGỮ CẢNH (CONTEXT) - Giữ nguyên 100%
  "asset_id": ObjectId, // Null hoặc ID ảnh/video (nếu comment vào ảnh cụ thể)

  // 2. NỘI DUNG & MEDIA - Giữ nguyên 100%
  "content": String, // Text nội dung HIỆN TẠI (Mới nhất)
  
  "media": { 
      "type": String, // 'image', 'gif', 'sticker', 'video'
      "url": String,  
      "display_meta": { "width": Number, "height": Number }
  },
  
  "mentions": [ObjectId], // Danh sách user_id được tag

  // 3. CẤU TRÚC PHẢN HỒI - Giữ nguyên 100%
  "parent_comment_id": ObjectId, // Null nếu là level 1
  "root_comment_id": ObjectId,   // Null nếu là level 1
  
  // 4. TRẠNG THÁI & MODERATION - Giữ nguyên 100%
  "status": String, // 'active', 'hidden', 'deleted', 'pending'
  
  "hidden_metadata": {
      "is_hidden": Boolean, 
      "hidden_at": Date,
      "hidden_by_user_id": ObjectId, 
      "reason": String, 
      "is_ghost_banned": Boolean 
  },
  
  "deleted_metadata": {
      "deleted_at": Date, 
      "deleted_by_user_id": ObjectId 
  },
  
  "report_count": Number,

  // 5. METRICS - Giữ nguyên 100%
  "reactions": {
      "total": Number,
      "like": Number,
      "love": Number,
      "haha": Number,
      "wow": Number,
      "sad": Number,
      "angry": Number
  },
  "reply_count": Number, 
  "mention_count": Number,

  // 6. LỊCH SỬ CHỈNH SỬA (Đã tách mảng array -> Dùng cờ flag)
  "is_edited": Boolean, // True nếu đã từng sửa
  "last_edited_at": Date, // [MỚI] Thời điểm sửa gần nhất (để hiển thị UI "Đã sửa 5p trước")

  "created_at": Date,
  "updated_at": Date,
  "deleted_at": Date    
}
Entity_Edit_Logs collection{
    "_id": ObjectId,
  
  // Định danh (Mapping về comment gốc)
  "target_collection": String, // "comments"
  "target_id": ObjectId, // ID của Comment gốc
  
  // Meta data của lần sửa
  "version": Number, // Lần 1, 2, 3...
  "editor_id": ObjectId, // Người thực hiện sửa
  "edited_at": Date, // (Mapping từ trường edited_at trong mảng cũ)
  
  // Nội dung ĐÃ BỊ THAY THẾ (Lưu lại cái cũ để restore nếu cần)
  "diff": {
      "old_content": String, // (Mapping từ trường content trong mảng cũ)
      "new_content": String,
      
      // Nếu sửa cả ảnh đính kèm
      "old_media": { "type": String, "url": String },
      "new_media": { "type": String, "url": String }
  },
  
  "ip_address": String, 
  "user_agent": String
}

--Truyền thông đa phương tiện (Multimedia & Short Form)
-- Stories
Stories Collection{
  "_id": ObjectId,
  "user_id": UUID, // Ref sang Postgres
  
  // Media (SeaweedFS)
  "media": {
    "url": String, 
    "type": String, // 'image', 'video'
    "duration": Number,
    "thumbnail_url": String
    "size_bytes": Number // [MỚI] Để quản lý dung lượng user upload
  },
// [CẬP NHẬT] Logic Privacy nâng cao
  "privacy": {
      "type": String, // 'public', 'friends', 'close_friends', 'custom'
      "allow_list": [UUID], // Nếu là 'custom'
      "block_list": [UUID]  // Ẩn story với người yêu cũ
  },
  "settings": {
      "allow_reply": Boolean,
      "allow_share": Boolean // Cho phép người khác share story của mình không?
  },
  // [MỚI] Cache danh sách người xem gần nhất (Để hiển thị UI nhanh mà không cần gọi Cassandra)
  "preview_viewers": [
      { "user_id": UUID, "avatar": String, "name": String } // Lưu tối đa 3 người mới nhất
  ],
  // Interactive Overlays (Cấu trúc phức tạp -> Thế mạnh Mongo)
  // Lưu tọa độ, nội dung của Sticker, Poll, Question
  "overlays": [
    {
       "type": "poll", // 'text', 'music', 'mention', 'location', 'poll'
       "position": { "x": 0.5, "y": 0.5, "rotation": 0, "scale": 1 },
       "data": {
           "question": "Ăn tối chưa?",
           "options": ["Rồi", "Chưa"] // Poll
       }
    },
    {
       "type": "music",
       "data": { "track_id": ObjectId, "start_time": 15, "duration": 15 }
    }
  ],

  // Logic 24h
  "created_at": Date,
  "expires_at": Date, // Index TTL (nếu muốn xóa thật) hoặc Index thường để filter
  "is_archived": Boolean, // [MỚI] Sau 24h, set true thay vì xóa (nếu user bật lưu trữ)
  // Settings
  "privacy": String, // 'public', 'friends', 'close_friends'
  "allow_reply": Boolean,
  
 "stats": {
      "views_count": Number,
      "likes_count": Number,
      "reply_count": Number // [MỚI] Số tin nhắn phản hồi
  }
}
-- Cassandra
CREATE TABLE story_views (
    story_id UUID,
    viewed_at TIMESTAMP,
    viewer_id UUID,
    
    -- [MỚI] Denormalization: Lưu luôn tên/avatar người xem để load list cho nhanh, 
    -- tránh việc query lại Postgres 1000 lần cho 1000 views.
    viewer_name TEXT,
    viewer_avatar_url TEXT,

    interaction_type TEXT, -- 'view', 'reaction', 'poll_vote'
    reaction_code TEXT, -- [MỚI] Lưu '❤️', '😂', '😡' hoặc ID của sticker
    poll_option_index INT, -- [MỚI] Nếu vote, thì vote cho phương án 0 hay 1?
    
    PRIMARY KEY (story_id, viewed_at, viewer_id)
) WITH CLUSTERING ORDER BY (viewed_at DESC);
--REELS
Reels Collection{
  "_id": ObjectId,
  "user_id": UUID,
  // [MỚI] Quản lý trạng thái xử lý Video (CỰC KỲ QUAN TRỌNG)
  "processing_status": String, // 'pending', 'processing', 'active', 'failed'
  "video": {
      "url": String, // Stream URL (HLS/DASH từ SeaweedFS)
      "width": Number, // Luôn ưu tiên dọc (9:16)
      "height": Number,
      "duration": Number
      "thumbnail_url": String, // Ảnh tĩnh (JPG)
      "preview_gif_url": String // [MỚI] Ảnh động khi rê chuột vào video
  },
  
  "caption": String,
  "hashtags": [String], // Index để tìm kiếm/trend
  "mentions": [UUID], // [MỚI] Danh sách người được tag trong caption

  // Tính năng Âm nhạc (Quan trọng của Reels)
  "audio_meta": {
      "track_id": ObjectId, // Ref tới Music Library
      "is_original_audio": Boolean, // True nếu dùng mic thu trực tiếp
      "volumn_adjust": Number
      "audio_start_time": Number // [MỚI] User có thể chọn đoạn điệp khúc bắt đầu từ giây thứ 30
  },

  // Tính năng Remix/Duet
  "remix_info": {
      "parent_reel_id": ObjectId, // Video gốc được remix
      "type": String // 'duet', 'remix'
      "is_remixable": Boolean // [MỚI] Chủ video có cho phép người khác Remix không?
  },

  // Metrics (Cập nhật định kỳ từ Redis hoặc dùng $inc)
  "stats": {
      "views": Number,
      "likes": Number,
      "shares": Number,
      "saves": Number, // Lưu lại xem sau
      "comments": Number
  },

  "privacy": String,
  "created_at": Date
  "deleted_at": Date // [MỚI] Soft delete
}
Music_Library collection{
  "_id": ObjectId,
  "title": String,
  "artist": String,
  "album": String, // [MỚI] Tên album
  "cover_url": String,
  "stream_url": String, // File mp3/aac cắt sẵn 15s, 30s, 60s trên SeaweedFS
  "duration": Number, // [MỚI] Độ dài bài nhạc
  // [MỚI] Phục vụ tìm kiếm lời bài hát
  "lyrics_snippet": String, // Đoạn lời nổi bật để search
  // [MỚI] Quản lý bản quyền (Tránh bị report)
  "copyright_info": {
      "provider": String, // Sony, Universal...
      "allowed_regions": [String] // ['VN', 'US']
  },
  "genre": [String], // Pop, Rock...
  "usage_count": Number // Để tính Trending Music
  
}
--LIVESTREAM
Live_Sessions collection{
 "_id": ObjectId,
  "host_user_id": UUID,
  
  "title": String,
  "description": String, // [MỚI] Mô tả buổi live
  "category_id": String, // [MỚI] Game, Talkshow, Music (Để filter)
  
  "status": String, 
  
  // Thông tin kỹ thuật
  "stream_key": String, 
  "playback_url": String, 
  "recording_setting": { // [MỚI] Cài đặt có lưu lại video sau khi live không
      "is_recorded": Boolean,
      "archive_url": String
  },

  // [MỚI] Quản lý phòng Live
  "banned_users": [UUID], // Những người bị kick khỏi phòng
  "pinned_comment_id": UUID, // ID của comment đang được ghim lên màn hình

  "started_at": Date,
  "ended_at": Date,
  
  "stats": {
     "peak_viewers": Number, // Mắt xem cao nhất
     "total_views": Number,
     "total_likes": Number,
     "total_comments": Number
  }
}
-- Cassandra
CREATE TABLE live_comments (
    stream_id UUID, 
    created_at TIMESTAMP, 
    comment_id UUID,
    
    user_id UUID,
    -- [MỚI - QUAN TRỌNG] Lưu cứng thông tin User tại thời điểm comment
    user_nickname TEXT, 
    user_avatar_url TEXT,
    user_badges SET<TEXT>, -- Ví dụ: {'top_fan', 'moderator'}
    
    content TEXT,
    is_pinned BOOLEAN, 
    
    PRIMARY KEY (stream_id, created_at, comment_id)
) WITH CLUSTERING ORDER BY (created_at DESC);


--Nhóm (Groups)
Groups Collection{
 "_id": ObjectId,
  "creator_id": UUID,
  
  // 1. Định danh & SEO
  "name": String,
  "slug": String, // Index Unique
  "description": String,
  "tags": [String], 
  
  // 2. Media (Giữ nguyên)
  "cover": { "url": String, "position_y": Number },
  "avatar": { "url": String },
  
  // 3. Phân loại & Quyền (Giữ nguyên)
  "privacy": String, // 'public', 'private', 'secret'
  "category_id": String, 
  
  // 4. Nội quy (Giữ lại vì text này thường không quá lớn và ít thay đổi)
  // Nếu nội quy quá dài, cũng có thể tách ra, nhưng thường giữ ở đây là OK.
  "rules": [ 
      { "title": String, "content": String }
  ],

  // 5. Cài đặt quản trị (Giữ nguyên)
  "settings": {
      "require_approval_to_join": Boolean,
      "require_approval_to_post": Boolean,
      "allow_member_posting": Boolean,
      "who_can_approve_member": String
  },

  // 6. Cấu hình tính năng (Feature Flags)
  // [QUAN TRỌNG] Đã XÓA mảng 'channels'. Chỉ giữ cờ bật/tắt.
  "community_chats": {
      "is_enabled": Boolean 
  },
  
  // [QUAN TRỌNG] Đã XÓA mảng 'join_questions'. Chỉ giữ cờ.
  "membership_questions": {
      "is_enabled": Boolean
  },

  // 7. Metrics (Giữ nguyên)
  "stats": {
      "member_count": Number,
      "post_count": Number,
      "pending_member_count": Number,
      "pending_post_count": Number,
      "reported_post_count": Number
  },

  "created_at": Date,
  "updated_at": Date,
  "deleted_at": Date
}
Group_Join_Questions {
    "_id": ObjectId,
  "group_id": ObjectId, // Reference tới Groups
  
  "content": String, // "Bạn sinh năm bao nhiêu?"
  "type": String, // 'text', 'multiple_choice', 'checkbox'
  
  // Nếu là trắc nghiệm thì có thêm options
  "options": [
      { "text": "Dưới 18", "value": "under_18" },
      { "text": "Trên 18", "value": "over_18" }
  ],
  
  "is_required": Boolean, // Bắt buộc trả lời không?
  "order": Number, // Thứ tự hiển thị 1, 2, 3
  
  "created_at": Date,
  "updated_at": Date
}
Group_Members collection{
  "_id": ObjectId,
  "group_id": ObjectId, // Ref Groups
  "user_id": UUID,      // Ref Users
  // [CẬP NHẬT] Thêm status 'invited'
  "status": String, // 'active', 'pending' (xin vào), 'invited' (được mời), 'banned', 'muted'
  // 1. Vai trò & Trạng thái
  "role": String,   // ENUM: 'admin', 'moderator', 'member'
  "status": String, // ENUM: 'active', 'pending' (chờ duyệt), 'banned', 'muted' (bị cấm chat)
  // Thông tin mời/xin vào
  "inviter_id": UUID, // Người mời (nếu status = invited)
  // 2. Nếu status = 'pending', lưu câu trả lời của họ
  "join_answers": [
      {
          "question_id": ObjectId,
          "answer": "Mình sinh năm 1995"
      }
  ],

  // 3. Nếu status = 'banned' hoặc 'muted'
  "discipline_info": {
      "reason": String,
      "banned_by": UUID,
      "until_date": Date // Null nếu ban vĩnh viễn
  },

  // 4. Huy hiệu & Gamification
  "badges": [
      "founding_member", // Thành viên sáng lập
      "top_contributor", // Fan cứng (Tính dựa trên số like/comment)
      "conversation_starter" // Người bắt chuyện
  ],
  
  // 5. Metadata cá nhân trong nhóm
  "joined_at": Date, // Ngày tham gia
  "last_active_at": Date, // Lần cuối vào nhóm (để lọc thành viên ảo)
  "inviter_id": UUID // Ai là người mời vào nhóm?
}
Group_Events Collection{
  "_id": ObjectId,
  "group_id": ObjectId,
  "creator_id": UUID,

  "title": String,
  "description": String,
  "cover_url": String,
  
  // Thời gian
  "start_time": Date,
  "end_time": Date,
  
  // Địa điểm
  "location": {
      "type": "online", // hoặc 'offline'
      "address": String, // "Link Zoom" hoặc "Phố đi bộ"
      "coordinates": [Number, Number] // GeoJSON
  },
  
  // Người tham gia (Chỉ lưu số lượng, list chi tiết lưu collection riêng nếu sự kiện lớn)
  "attendees_count": {
      "going": Number,
      "interested": Number
  },
  
  "created_at": Date
}
-- 4. Redis: Caching & Performance (Best Practice)
-- Đây là phần quan trọng để hệ thống không bị chậm khi check quyền hạn liên tục.

-- Key 1: Kiểm tra quyền (Permission Cache)

-- Key: group_perm:{group_id}:{user_id}

-- Value: "admin" | "mod" | "member" | "banned"

-- TTL: 1 giờ (Xóa key khi Admin thay đổi quyền của user).

-- Tác dụng: Khi User đăng bài, Backend check Redis thay vì query MongoDB.

-- Key 2: Counter Cache

-- Key: group_stats:{group_id}

-- Value: { members: 10500, online: 120 }

-- Tác dụng: Hiển thị số thành viên ngay tiêu đề nhóm mà không cần count() trong DB.
Group_Files collection{
"_id": ObjectId,
   "group_id": ObjectId,
   "uploader_id": UUID,
   
   "file_name": String, // "Bao_cao_tai_chinh.xlsx"
   "file_type": String, // 'pdf', 'docx', 'xlsx'
   "file_size": Number,
   "seaweedfs_file_id": String, // ID file trong SeaweedFS
   
   "download_count": Number,
   "created_at": Date
}

--Module Trang (Fanpages)
--Collection: Pages (Thông tin lõi)
Pages Collection{
  "_id": ObjectId,
  "creator_user_id": UUID, // Người tạo ra page (Owner gốc)
  
  // 1. Định danh
  "name": String,
  "slug": String, // Unique Index (VD: @cafebiz.vn)
  "category_id": String, // Ref tới danh mục (F&B, Media, Blog...)
  "is_verified": Boolean, // Tích xanh
  "status": String, // 'published', 'unpublished', 'banned'

  // 2. Branding
  "avatar": { "url": String },
  "cover": { "url": String, "position_y": Number },
  
  // 3. Thông tin Doanh nghiệp (Business Info)
  "bio": String,
  "website": String,
  "email": String,
  "phone_number": String,
  "address": {
      "street": String,
      "city": String,
      "zipcode": String,
      "coordinates": [Number, Number] // GeoJSON để tìm quán ăn gần đây
  },
  
  // 4. Giờ mở cửa (Cấu trúc linh động của Mongo)
  "business_hours": [
      { "day": "Monday", "open": "08:00", "close": "22:00" },
      { "day": "Tuesday", "open": "08:00", "close": "22:00" }
      // ...
  ],

  // 5. Nút hành động (Call To Action - CTA)
  "cta_button": {
      "type": "send_message", // 'call_now', 'visit_website', 'shop_now'
      "value": "https://m.me/..."
  },

  // 6. Cài đặt chuyên sâu
  "settings": {
      "allow_visitor_post": Boolean, // Cho khách đăng bài lên tường không?
      "profanity_filter": String, // 'off', 'medium', 'strong' (Lọc từ tục tĩu)
      "messaging_status": String // 'online', 'away'
  },

  // 7. Metrics (Denormalization)
  "stats": {
      "followers_count": Number,
      "likes_count": Number,
      "rating_score": Number, // 4.8/5.0
      "review_count": Number
  },

  "created_at": Date,
  "updated_at": Date
}
Page_Roles Collection{
"_id": ObjectId,
  "page_id": ObjectId,
  "user_id": UUID,
  
  // Vai trò (Best Practice của Facebook)
  "role": String, 
  // ENUM: 
  // 'admin' (Toàn quyền), 
  // 'editor' (Đăng bài, sửa bài), 
  // 'moderator' (Trả lời comment, ban user), 
  // 'advertiser' (Chỉ chạy ads), 
  // 'analyst' (Chỉ xem thống kê)

  "custom_permissions": [String], // Nếu muốn phân quyền chi tiết hơn (VD: ['manage_jobs', 'manage_events'])
  
  "created_at": Date,
  "assigned_by": UUID // Ai là người cấp quyền này?
}
Page_Followers collection{
  "_id": ObjectId,
  "page_id": ObjectId, // Index
  "user_id": UUID,     // Index
  
  "type": String, // 'follow' hoặc 'like' (Facebook phân biệt 2 cái này)
  
  "settings": {
      "notification_level": String, // 'all', 'highlight', 'off'
      "is_favorite": Boolean // "See First" (Xem trước trên Newfeed)
  },
  
  "followed_at": Date
}
--Cassandra: Insights & Analytics (Thống kê)
-- Cassandra
CREATE TABLE page_daily_metrics (
    page_id UUID,
    metric_date DATE, -- 2024-01-01
    
    -- Các chỉ số quan trọng
    reach_total BIGINT, -- Số người tiếp cận
    reach_paid BIGINT, -- Tiếp cận từ quảng cáo
    reach_organic BIGINT, -- Tiếp cận tự nhiên
    
    impressions_total BIGINT, -- Số lần hiển thị
    
    new_followers INT,
    unfollows INT,
    
    profile_views INT,
    website_clicks INT,
    cta_clicks INT,
    
    PRIMARY KEY (page_id, metric_date)
) WITH CLUSTERING ORDER BY (metric_date DESC);
-- Cassandra
CREATE TABLE post_insights (
    post_id UUID, -- Partition Key
    
    -- Snapshot trọn đời (Lifetime)
    reach BIGINT,
    impressions BIGINT,
    engagement_rate FLOAT, -- Tỷ lệ tương tác
    
    -- Chi tiết tương tác
    reactions_total INT,
    comments_total INT,
    shares_total INT,
    clicks_total INT,
    video_views_3s INT, -- View video trên 3 giây
    
    updated_at TIMESTAMP,
    PRIMARY KEY (post_id)
);
--Postgres: Ads & Billing (Quảng cáo)
--Postgres
CREATE TABLE Ad_Accounts (
    account_id UUID PRIMARY KEY,
    owner_user_id UUID NOT NULL REFERENCES users(user_id), -- Người sở hữu tài khoản
    currency VARCHAR(3) DEFAULT 'VND',
    timezone VARCHAR(50),
    balance DECIMAL(15, 2), -- Số dư (nếu trả trước)
    credit_limit DECIMAL(15, 2), -- Hạn mức tín dụng (nếu trả sau)
    status VARCHAR(20), -- 'active', 'disabled', 'settled'
    
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE TABLE Ad_Campaigns (
    campaign_id UUID PRIMARY KEY,
    account_id UUID REFERENCES Ad_Accounts(account_id),
    
    name VARCHAR(255),
    objective VARCHAR(50), -- 'reach', 'traffic', 'messages', 'conversions'
    buying_type VARCHAR(20), -- 'auction', 'fixed_price'
    
    daily_budget DECIMAL(15, 2),
    lifetime_budget DECIMAL(15, 2),
    
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    status VARCHAR(20) -- 'active', 'paused', 'completed', 'archived'
);
CREATE TABLE Ads (
    ad_id UUID PRIMARY KEY,
    campaign_id UUID REFERENCES Ad_Campaigns(campaign_id),
    
    -- Liên kết với bài viết bên MongoDB
    target_post_id VARCHAR(50), -- Lưu ObjectId dạng String
    
    bid_amount DECIMAL(10, 2), -- Giá thầu
    status VARCHAR(20), -- 'reviewing' (đang xét duyệt), 'active', 'rejected'
    rejection_reason TEXT -- Lý do từ chối nếu có
);
-- 5. Redis: Caching & Performance
-- Page Role Cache: page_role:{page_id}:{user_id} -> Value: 'admin'.

-- Tác dụng: Khi User vào quản lý Page, không cần query Mongo Page_Roles liên tục.

-- Post Scheduler Queue: Sử dụng Redis ZSET (Sorted Set).

-- Key: scheduled_posts

-- Score: Timestamp (thời gian đăng).

-- Value: post_id.

-- Worker: Quét Redis mỗi phút, nếu score <= now() thì đổi status bài viết từ draft -> published.


-- chat 

Conversations collection{
"_id": ObjectId,
  "type": String, // 'private', 'group'
  "scope": String, // 'messenger', 'community_channel'
  "status": String, // 'active', 'pending', 'spam'
  
  // 1. Group Info
  "name": String, 
  "avatar": { "url": String },
  
  // 2. Ownership & Linking
  "creator_id": UUID,
  "owner_id": UUID, 
  "related_group_id": ObjectId, // Link tới Module Groups (nếu là kênh cộng đồng)

  // 3. Permissions (Giữ nguyên)
  "permissions": {
      "send_message": String, // 'everyone', 'admin_only'
      "add_member": String    // 'everyone', 'admin_only'
  },

  // 4. Interface (Giữ nguyên)
  "theme": { 
      "color": String, 
      "emoji": String, 
      "background_url": String 
  },

  // 5. Caching (Giữ nguyên để hiển thị list bên ngoài)
  "last_message": {
      "message_id": TimeUUID, 
      "content": String, 
      "sender_id": UUID,
      "type": String, 
      "created_at": Date
  },
  
  // [BỔ SUNG] Counter để biết nhóm đông thế nào mà không cần count collection kia
  "participant_count": Number, 

  "created_at": Date,
  "updated_at": Date
}
Conversation_Participants collection{

    "_id": ObjectId,
  
  // --- KHÓA NGOẠI (LINKING) ---
  "conversation_id": ObjectId, // Thuộc về nhóm chat nào
  "user_id": UUID,             // Là thành viên nào
  
  // --- CÁC TRƯỜNG TỪ MẢNG 'MEMBERS' CŨ (GIỮ ĐỦ 100%) ---
  
  // 1. Vai trò & Định danh trong nhóm
  "role": String, // 'admin', 'moderator', 'member'
  "nickname": String, // Biệt danh (nếu đặt riêng trong chat này)
  
  // 2. Trạng thái Đọc (Read Status)
  "last_seen_at": Date, // Thời điểm cuối user mở box chat
  "last_seen_message_id": String, // [BỔ SUNG] ID tin nhắn cuối cùng đã đọc (để tính unread count)
  
  // 3. Cài đặt cá nhân (Personal Settings)
  "is_muted": Boolean,      // Tắt thông báo
  "is_archived": Boolean,   // Lưu trữ (ẩn khỏi list)
  
  // 4. Tính năng Xóa lịch sử
  "clear_history_at": Date, // Mốc thời gian xóa (Chỉ hiện tin nhắn sau mốc này)
  
  // 5. Metadata tham gia
  "joined_at": Date,
  "added_by_user_id": UUID // [BỔ SUNG] Ai là người add vào nhóm?
}
Call_Logs collection{
    "_id": ObjectId,
  "conversation_id": ObjectId,
  "caller_id": UUID,
  
  "participants": [UUID], // Những người tham gia
  
  "type": String, // 'voice', 'video'
  "status": String, // 'missed', 'ended', 'rejected', 'busy'
  
  "started_at": Date,
  "ended_at": Date,
  "duration_seconds": Number,
  
  "is_group_call": Boolean
}

-- Cassandra
-- Partition Key: conversation_id (Gom tin nhắn cùng 1 box chat)
-- Clustering Key: bucket (Phân mảnh theo tháng) + message_id (TimeUUID - Sắp xếp thời gian)

CREATE TABLE messages (
conversation_id TEXT, 
    bucket INT, 
    message_id TIMEUUID, 
    
    sender_id UUID,
    type TEXT, -- 'text', 'image', 'video', 'audio', 'file', 'location', 'story_reply', 'sticker'
    content TEXT, 
    attachments LIST<TEXT>, 
    
    reply_to_message_id UUID,
    story_ref_id UUID, -- [MỚI] Nếu reply story
    
    is_revoked BOOLEAN, 
    
    created_at TIMESTAMP,
    
    PRIMARY KEY ((conversation_id, bucket), message_id)
) WITH CLUSTERING ORDER BY (message_id DESC);
-- Cassandra
-- User A đã đọc đến tin nhắn nào trong cuộc hội thoại X?
CREATE TABLE conversation_read_state (
    conversation_id TEXT,
    user_id UUID,
    last_read_message_id TIMEUUID, -- User đã đọc đến tin này
    last_read_at TIMESTAMP,
    PRIMARY KEY (conversation_id, user_id)
);
-- [MỚI] Cassandra: Message Reactions
-- Partition Key: message_id (Để load tất cả reaction của 1 tin nhắn)
-- Clustering Key: user_id (Để biết ai thả)

CREATE TABLE message_reactions (
    conversation_id TEXT, -- Partition Key phụ (Optional - để dễ dọn dẹp)
    message_id TIMEUUID,
    user_id UUID,
    
    reaction_code TEXT, -- '❤️', '😆', '😢'
    created_at TIMESTAMP,
    
    PRIMARY KEY ((conversation_id, message_id), user_id)
);
3. Redis: Real-time & Presence (Trạng th-- [MỚI] Cassandra: Message Reactions
-- Partition Key: message_id (Để load tất cả reaction của 1 tin nhắn)
-- Clustering Key: user_id (Để biết ai thả)

CREATE TABLE message_reactions (
    conversation_id TEXT, -- Partition Key phụ (Optional - để dễ dọn dẹp)
    message_id TIMEUUID,
    user_id UUID,
    
    reaction_code TEXT, -- '❤️', '😆', '😢'
    created_at TIMESTAMP,
    
    PRIMARY KEY ((conversation_id, message_id), user_id)
);
3. Redis: Real-time & Presence (Trạng thái)
Sử dụng Redis để xử lý các tính năng tức thời.

A. Trạng thái Online (Presence)
Key: user:presence:{user_id}

Value: {"status": "online", "last_active": 1700000000, "device": "mobile"}

TTL: 5 phút (Client phải gửi heartbeat ping lên server mỗi 3 phút để gia hạn key. Nếu không ping -> Hết hạn -> Coi như Offline).

B. Trạng thái đang nhập (Typing)
Key: chat:typing:{conversation_id}

Data Structure: Set (Tập hợp các user_id đang gõ).

Value: [user_id_1, user_id_2]

TTL: 5 giây (Nếu user ngừng gõ, key tự mất).

 ES Index: search_users
{
  "_id": "user_id_uuid",
  "full_name": "Nguyen Van A", // Text (Analyzer: Vietnamese)
  "username": "nguyenvana", // Keyword
  "email": "a@gmail.com",
  "bio": "Yêu màu hồng",
  "avatar": "url...",
  
  // Quan trọng cho Ranking (Điểm số hiển thị)
  "follower_count": 10500,
  "is_verified": true, // Tích xanh -> Boost điểm cao hơn
  "location": { "lat": 10.7, "lon": 106.6 }, // Geo-point

  "type": "user" // hoặc 'page'
}

 ES Index: search_posts
{
  "_id": "post_id_objectid",
  "content": "Hôm nay trời đẹp quá #dulich #dalat", // Text
  "hashtags": ["dulich", "dalat"], // Keyword
  
  "author_id": "user_uuid",
  "group_id": "group_objectid", // Null nếu post tường nhà
  "page_id": "page_objectid",

  // Media (Để lọc tìm riêng Ảnh/Video)
  "media_types": ["image", "video"], 
  
  "privacy": "public", // Chỉ index bài Public/Group Public
  
  "created_at": "2024-01-22T10:00:00Z",
  
  // Metrics để Ranking (Bài nhiều like lên trước)
  "likes_count": 500,
  "comments_count": 20
}


// ES Index: search_groups
{
  "_id": "group_id_objectid",
  "name": "Hội Yêu Mèo",
  "description": "Chia sẻ kinh nghiệm nuôi mèo",
  "tags": ["pet", "cat", "hcm"],
  "privacy": "public", // Không index nhóm Secret
  "member_count": 50000,
  "location": { "lat": ..., "lon": ... }
}

Notification_Templates collection{
    "_id": ObjectId,
  "type": "POST_LIKE",
  "template": {
      "vi": "<strong>{{actor_name}}</strong> đã thích bài viết của bạn.",
      "en": "<strong>{{actor_name}}</strong> liked your post."
  },
  "icon_url": "like_icon.png",
  "action_link": "/post/{{target_id}}" // Deeplink
}

User_Notification_Settings collection{
    "_id": ObjectId,
  "user_id": UUID,
  "settings": {
      "push_enabled": Boolean, // Tổng
      "email_frequency": String, // 'daily', 'weekly', 'never'
      
      // Chi tiết
      "on_post_like": { "push": true, "email": false },
      "on_comment": { "push": true, "email": true },
      "on_birthday": { "push": true, "email": false }
  },
  "fcm_tokens": [ // Token để gửi xuống Mobile App (Firebase Cloud Messaging)
      { "token": "fcm_xyz...", "device_id": "iphone_13", "updated_at": Date }
  ]
}

-- Cassandra
-- Partition Key: user_id (Gom thông báo của 1 người vào 1 chỗ)
-- Clustering Key: created_at (Sắp xếp mới nhất lên đầu)

CREATE TABLE notifications (
    user_id UUID, -- Người nhận (Recipient)
    created_at TIMESTAMP,
    notification_id TIMEUUID,
    
    type TEXT, -- 'POST_LIKE', 'COMMENT_REPLY', 'FRIEND_REQUEST'
    
    -- Actor (Người gây ra hành động - Denormalized để hiển thị nhanh)
    actor_id UUID,
    actor_name TEXT,
    actor_avatar TEXT,
    
    -- Target (Đối tượng bị tác động)
    target_id TEXT, -- ID bài viết, ID comment, ID group...
    target_preview TEXT, -- "Bài viết về con mèo..." (Snapshot ngắn)
    
    -- Status
    is_read BOOLEAN,
    is_clicked BOOLEAN,
    
    -- Grouping (Dùng để gộp thông báo: "A, B và 5 người khác đã like")
    group_key TEXT, -- VD: "POST_LIKE:post_id_123"
    
    PRIMARY KEY (user_id, created_at, notification_id)
) WITH CLUSTERING ORDER BY (created_at DESC);


III. GRAPH DATABASE (NEO4J) - RECOMMENDATION ENGINE
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

User_Settings collection{
 "_id": ObjectId,
  "user_id": UUID, // Index Unique
  
  // --- 1. APP CONFIGURATION (Phẳng hóa) ---
  // Thay vì appearance.theme -> theme_mode
  "theme_mode": String, // 'system', 'light', 'dark', 'high_contrast'
  "font_size": String, // 'medium', 'large'
  "compact_mode": Boolean, 
  
  "lang_code": String, // 'vi', 'en', 'jp'
  "timezone": String, // 'Asia/Ho_Chi_Minh'
  "auto_translate": Boolean, 

  // --- 2. PRIVACY DEFAULTS (Phẳng hóa) ---
  // Thay vì privacy_defaults.post_audience -> default_post_audience
  "default_post_audience": String, // 'public', 'friends', 'only_me'
  "default_story_audience": String, // 'friends', 'close_friends'

  // --- 3. ACCESS CONTROL (Phẳng hóa & Rename cho rõ nghĩa) ---
  "allow_friend_request_from": String, // 'everyone', 'friends_of_friends'
  "allow_friend_list_view_from": String, // 'public', 'friends', 'only_me'
  "allow_email_lookup_from": String, // 'everyone', 'friends'
  "allow_phone_lookup_from": String, // 'everyone', 'friends'
  "allow_search_engine_indexing": Boolean, 

  // --- 4. TIMELINE & TAGGING (Phẳng hóa) ---
  "allow_timeline_posting_from": String, // 'friends', 'only_me'
  "review_tags_enabled": Boolean, // Tên ngắn gọn hơn review_tags_before_appearing
  "review_timeline_posts_enabled": Boolean, // Tên ngắn gọn hơn review_posts_before_appearing

  // --- 5. NOTIFICATIONS (Giữ nguyên hoặc tách ra nếu quá nhiều) ---
  // Vì notification settings thường đi theo cụm, giữ object này cũng được, 
  // nhưng nếu quá nhiều loại noti thì nên tách ra collection riêng.
  "notifications": {
      "email_frequency": String, // 'weekly'
      "push_interactions": Boolean, // Like, comment
      "push_friends": Boolean, // Friend request
      "push_groups": Boolean,
      "push_events": Boolean,
      "push_birthdays": Boolean
  },

  "updated_at": Date
 }

postgres
 CREATE TABLE User_Sessions (
    session_id UUID PRIMARY KEY,
    user_id UUID NOT NULL  REFERENCES users(user_id),
    
    -- Thông tin thiết bị
    device_name VARCHAR(255), -- "iPhone 14 Pro"
    os_version VARCHAR(50), -- "iOS 17.2"
    browser VARCHAR(50), -- "Chrome Mobile"
    ip_address VARCHAR(45),
    
    -- Location (GeoIP)
    location_city VARCHAR(100),
    location_country VARCHAR(100),
    
    -- Trạng thái
    refresh_token TEXT, -- Token để gia hạn
    is_active BOOLEAN DEFAULT TRUE,
    last_active_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);
-- Index user_id để load danh sách nhanh
CREATE INDEX idx_sessions_user ON User_Sessions(user_id);

3. Redis: Caching Logic Chặn (Critical)
Hệ thống phải kiểm tra "A có chặn B không?" hàng triệu lần (mỗi khi load comment, load chat, search). Query Postgres liên tục sẽ sập DB.

Design Redis Key:

Key: user:blocks:{user_id}

Type: Set (Tập hợp các ID bị user này chặn).

Value: [blocked_user_id_1, blocked_user_id_2, ...]

Logic:

User A vào xem Profile B.

Code check Redis: SISMEMBER user:blocks:B A (B có chặn A không?).

Nếu True -> Trả về "Content Unavailable".

-- Cassandra
-- Partition Key: target_id (Gom tất cả like của 1 bài viết/comment vào 1 chỗ)
-- Clustering Key: user_id (Để check nhanh user A có like chưa)

CREATE TABLE entity_reactions (
    target_id TEXT, -- ID của Post hoặc Comment
    target_type TEXT, -- 'post' hoặc 'comment'
    user_id UUID,
    
    reaction_code TEXT, -- 'like', 'love', 'haha', 'sad', 'angry'
    created_at TIMESTAMP,
    
    PRIMARY KEY (target_id, user_id)
);
-- Bảng phụ: Để phục vụ "Nhật ký hoạt động" (Activity Log)
-- Xem lại lịch sử mình đã like những gì
CREATE TABLE user_reaction_history (
    user_id UUID,
    created_at TIMESTAMP,
    target_id TEXT,
    target_type TEXT,
    reaction_code TEXT,
    PRIMARY KEY (user_id, created_at)
) WITH CLUSTERING ORDER BY (created_at DESC);

// MongoDB: User_Saved_Items
// Index: { user_id: 1, created_at: -1 }
User_Saved_Items collection
{
  "_id": ObjectId,
  "user_id": UUID,
  
  "target_id": ObjectId, // ID bài viết / video / reel
  "target_type": String, // 'post', 'reel', 'video'
  
  // Lưu snapshot nhẹ để hiển thị danh sách đã lưu mà không cần join
  "snapshot": {
      "author_name": String,
      "content_preview": String,
      "thumbnail_url": String
  },    
  
  "collection_name": String, // (Nâng cao) Lưu vào bộ sưu tập nào? VD: "Món ăn ngon"
  "created_at": Date
}
Search_History collection{
    "_id": ObjectId,
  "user_id": UUID,
  "keyword": String, // "quán cafe đẹp"
  
  "target_id": String, // Nếu user click vào 1 User/Group cụ thể từ gợi ý
  "target_type": String, // 'user', 'group'
  
  "created_at": Date // TTL Index: Tự xóa sau 30 ngày
}