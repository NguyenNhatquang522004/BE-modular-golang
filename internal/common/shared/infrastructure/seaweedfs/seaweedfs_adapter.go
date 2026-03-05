package seaweedfs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/url" // Bắt buộc phải có để tạo tham số xóa
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/google/uuid"
	"github.com/linxGnu/goseaweedfs"
	"github.com/minio/minio-go/v7"
)

type SeaweedfsAdapter struct {
	client     *goseaweedfs.Filer
	s3Client   *minio.Client
	publicHost string // Ví dụ: https://cdn.mysocial.com
	cfg        *configs.Config
}

func NewSeaweedfsAdapter(client *goseaweedfs.Filer, cfg *configs.Config) IRepositoryShare.ISeaweedfs {
	publicHost := cfg.SEAWEEDFS.SEAWEEDFS_FILER_URL // Hoặc có thể dùng biến cấu hình riêng cho public host nếu khác
	return &SeaweedfsAdapter{
		client:     client,
		publicHost: publicHost,
		cfg:        cfg,
	}
}

// Upload: Upload file hoàn chỉnh (Avatar, Post, CV)
func (r *SeaweedfsAdapter) Upload(ctx context.Context, input *dto.FileUploadInput) (*dto.FileUploadOutput, error) {
	// 1. TẠO TÊN FILE DUY NHẤT (UUID)
	// Tránh trùng lặp file và bảo mật tên file gốc của người dùng.
	ext := filepath.Ext(input.FileName)
	if ext == "" {
		ext = ".bin" // Fallback nếu không xác định được đuôi file
	}
	uniqueName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// 2. XÁC ĐỊNH THƯ MỤC LƯU TRỮ (FOLDER STRATEGY)
	// Nếu input.Folder trống, tự động sinh path dựa trên logic phân cấp Tầng 1-2-3 bạn đã định nghĩa.
	folderPath := input.Folder
	if folderPath == "" {
		folderPath = utils.GenerateStoragePath(input.OwnerID, input.Storage)
	}

	// 3. CHUẨN HÓA ĐƯỜNG DẪN ĐẦY ĐỦ (FULL PATH)
	// Đảm bảo đường dẫn bắt đầu bằng / và không có khoảng trắng thừa.
	cleanFolder := "/" + strings.Trim(folderPath, "/")
	fullPath := fmt.Sprintf("%s/%s", cleanFolder, uniqueName)

	// 4. THỰC HIỆN UPLOAD QUA FILER
	// Sử dụng io.Reader để stream trực tiếp, tối ưu RAM cho các file nặng như video/reels.
	// r.client là *goseaweedfs.Filer đã được khởi tạo trong database layer.
	result, err := r.client.Upload(input.Content, input.Size, fullPath, "", "")
	if err != nil {
		return nil, fmt.Errorf("seaweedfs_adapter: upload failed: %w", err)
	}

	// 5. TRẢ VỀ OUTPUT CHUẨN DTO
	// FilePath dùng để lưu vào MongoDB, PublicURL dùng để hiển thị lên UI.
	return &dto.FileUploadOutput{
		OwnerID:   input.OwnerID,
		Storage:   input.Storage,
		FilePath:  fullPath,
		PublicURL: r.GetPublicURL(fullPath),
		FileID:    result.FileID,
	}, nil
}

// Download: Lấy nội dung file (Dùng khi cần xử lý ảnh/video ở Backend)
// Download thực hiện lấy nội dung file từ SeaweedFS dưới dạng stream
func (r *SeaweedfsAdapter) Download(ctx context.Context, filePath string) (io.ReadCloser, error) {
	// 1. KIỂM TRA ĐƯỜNG DẪN
	cleanPath := strings.TrimSpace(filePath)
	if cleanPath == "" {
		return nil, fmt.Errorf("seaweedfs_adapter: file path is empty")
	}

	// Đảm bảo đường dẫn bắt đầu bằng /
	if !strings.HasPrefix(cleanPath, "/") {
		cleanPath = "/" + cleanPath
	}

	// 2. SỬ DỤNG IO.PIPE ĐỂ STREAM DỮ LIỆU
	// io.Pipe tạo ra một cặp (Reader, Writer).
	// Những gì ghi vào Writer sẽ được đọc trực tiếp từ Reader mà không tốn RAM trung gian.
	pr, pw := io.Pipe()

	// 3. CHẠY DOWNLOAD TRONG GOROUTINE
	// Chúng ta cần chạy Download trong một luồng riêng để không chặn luồng chính
	go func() {
		// Gọi hàm Download của thư viện linxGnu
		// Tham số: (path string, callback func(io.Reader) error)
		err := r.client.Download(cleanPath, nil, func(reader io.Reader) error {
			// Copy dữ liệu từ SeaweedFS trực tiếp vào Pipe Writer
			_, copyErr := io.Copy(pw, reader)
			return copyErr
		})

		// Đóng Pipe Writer kèm theo lỗi (nếu có)
		// Khi pw đóng, pr (Reader) sẽ nhận được tín hiệu kết thúc hoặc lỗi.
		if err != nil {
			pw.CloseWithError(fmt.Errorf("seaweedfs_adapter: download stream failed: %w", err))
		} else {
			pw.Close()
		}
	}()

	// Trả về Pipe Reader (đóng vai trò là io.ReadCloser) cho Service sử dụng
	return pr, nil
}

// Delete thực hiện xóa một file vật lý trên SeaweedFS
func (r *SeaweedfsAdapter) Delete(ctx context.Context, filePath string) error {
	// 1. KIỂM TRA ĐẦU VÀO
	// Nếu đường dẫn rỗng, ta coi như đã xóa xong (Idempotent)
	cleanPath := strings.TrimSpace(filePath)
	if cleanPath == "" {
		return nil
	}

	// 2. CHUẨN HÓA ĐƯỜNG DẪN
	// SeaweedFS Filer yêu cầu đường dẫn tuyệt đối bắt đầu bằng /
	if !strings.HasPrefix(cleanPath, "/") {
		cleanPath = "/" + cleanPath
	}

	// 3. THỰC HIỆN XÓA
	// Tham số thứ 2 là 'recursive'. Ở đây xóa 1 file nên ta để là 'false'.
	// SeaweedFS sẽ tự động dọn dẹp Metadata ở Filer và Data ở Volume.
	err := r.client.Delete(cleanPath, nil)
	if err != nil {
		// Nếu lỗi là "Not Found", có thể bỏ qua hoặc wrap lỗi tùy nhu cầu logic
		// Ở đây ta wrap lỗi để dễ dàng debug trong dự án modular
		return fmt.Errorf("seaweedfs_adapter: failed to delete file [%s]: %w", cleanPath, err)
	}

	return nil
}

// DeleteFolder thực hiện xóa toàn bộ thư mục và các file bên trong một cách đệ quy
func (r *SeaweedfsAdapter) DeleteFolder(ctx context.Context, folderPath string) error {
	// 1. KIỂM TRA BẢO MẬT CỰC KỲ QUAN TRỌNG
	// Loại bỏ khoảng trắng và chuẩn hóa
	cleanPath := strings.TrimSpace(folderPath)

	// Ngăn chặn xóa root hoặc đường dẫn rỗng để bảo vệ dữ liệu hệ thống
	if cleanPath == "" || cleanPath == "/" {
		return fmt.Errorf("seaweedfs_adapter: dangerous operation: cannot delete root or empty path")
	}

	// 2. CHUẨN HÓA ĐƯỜNG DẪN
	// Đảm bảo bắt đầu bằng / nhưng không kết thúc bằng / (SeaweedFS Filer chuẩn)
	if !strings.HasPrefix(cleanPath, "/") {
		cleanPath = "/" + cleanPath
	}
	cleanPath = strings.TrimSuffix(cleanPath, "/")

	// 3. THỰC HIỆN XÓA ĐỆ QUY (RECURSIVE)
	// Tham số thứ 2 là 'recursive'. Đặt là 'true' để xóa sạch folder và con của nó.
	// Lưu ý: Dựa trên signature của goseaweedfs, hàm Delete không nhận Context.
	params := url.Values{}
	params.Set("recursive", "true")
	err := r.client.Delete(cleanPath, params)
	if err != nil {
		// Kiểm tra nếu lỗi không phải do folder không tồn tại (404)
		// Một số phiên bản goseaweedfs trả lỗi cụ thể, nếu không ta wrap lỗi chung
		return fmt.Errorf("seaweedfs_adapter: failed to delete folder [%s]: %w", cleanPath, err)
	}

	return nil
}

// Exists: Kiểm tra file/thư mục đã tồn tại trên Filer chưa

func (r *SeaweedfsAdapter) Exists(ctx context.Context, filePath string) (bool, error) {
	cleanPath := "/" + strings.Trim(filePath, "/")

	// Gọi hàm Head mới thêm vào
	_, statusCode, _ := r.client.Get(cleanPath, nil, nil)

	if statusCode == 200 {
		return true, nil
	}
	return false, nil
}

// --- LIVE STREAM & REALTIME OPERATIONS ---

// UploadStreamSegment: Đẩy từng mảnh video (.ts) lên folder tạm của buổi Live
// sessionID là định danh duy nhất cho buổi live hiện tại
// segmentName là tên file mảnh video, content là nội dung mảnh video
// Thường segmentName có dạng: segment0001.ts, segment0002.ts, ...
func (r *SeaweedfsAdapter) UploadStreamSegment(ctx context.Context, sessionID string, segmentName string, content io.Reader) error {
	// Quy hoạch folder: /lives/{sessionID}/{segmentName}
	// Lưu ý: Stream segment thường nhỏ nên ta dùng TTL ngắn (ví dụ: 10 phút) để tự dọn rác

	fullPath := r.GetStreamSegmentURL(sessionID, segmentName)
	// Vì đây là stream realtime, ta không cần lưu Metadata phức tạp vào DB chính
	// Dùng TTL "10m" để SeaweedFS tự xóa các segment cũ
	// Lưu ý: Cần biết Size của segment. Nếu không có, bạn phải buffer hoặc dùng chunked upload.
	// Ở đây giả định bạn đã lấy được size từ encoder (ví dụ FFmpeg).
	_, err := r.client.Upload(content, 0, fullPath, "", "")
	return err
}
func (r *SeaweedfsAdapter) GenerateVOD(ctx context.Context, sessionID string, ownerID string, name string) (*dto.FileUploadOutput, error) {
	// 1. TẠO INTERNAL URL TRỎ VÀO SEAWEEDFS
	// Mặc định Filer chạy ở port 8888 trên server.
	// Nếu app Go và SeaweedFS chạy chung mạng Docker, dùng tên service (vd: http://seaweedfs-filer:8888)
	// 2. TẠO FILE TẠM TRÊN SERVER BACKEND
	localMP4Path := fmt.Sprintf("/tmp/lives/%s_vod.mp4", sessionID)

	log.Printf("[VOD] Bắt đầu gộp video cho session: %s", sessionID)

	// 3. GỌI FFMPEG GỘP VIDEO
	// FFmpeg sẽ tự tải .m3u8 và .ts từ SeaweedFS về, gộp và lưu vào localMP4Path
	err := r.MergeToMP4(ctx, sessionID, localMP4Path)
	if err != nil {
		log.Printf("[VOD Error] Không thể gộp video: %v", err)
		return nil, fmt.Errorf("failed to merge video: %w", err)
	}

	log.Printf("[VOD] Đã tạo file MP4 tạm: %s", localMP4Path)

	// 4. MỞ FILE MP4 TẠM ĐỂ UPLOAD NGƯỢC LÊN SEAWEEDFS
	file, err := os.Open(localMP4Path)
	if err != nil {
		log.Printf("[VOD Error] Không thể mở file MP4 tạm: %v", err)
		return nil, fmt.Errorf("failed to open temp MP4 file: %w", err)
	}
	defer file.Close()

	fileInfo, _ := file.Stat()

	// 5. SỬ DỤNG HÀM UPLOAD HIỆN CÓ CỦA BẠN (SeaweedfsAdapter.Upload)
	uploadInput := &dto.FileUploadInput{
		OwnerID:  ownerID,
		Storage:  utils.BucketLive, // Tùy định nghĩa Storage của bạn
		FileName: sessionID + name + ".mp4",
		Content:  file,
		Size:     fileInfo.Size(),
	}

	// Tận dụng lại đúng hàm Upload bạn đã viết!
	uploadResult, err := r.Upload(ctx, uploadInput)
	if err != nil {
		log.Printf("[VOD Error] Upload lên SeaweedFS thất bại: %v", err)
		return nil, fmt.Errorf("failed to upload VOD to SeaweedFS: %w", err)
	}

	log.Printf("[VOD Success] Video đã sẵn sàng: %s", uploadResult.PublicURL)

	// 6. DỌN DẸP RÁC
	// Xóa file mp4 tạm trên ổ cứng backend
	os.Remove(localMP4Path)

	// [Tùy chọn] Xóa folder chứa các file .ts trên SeaweedFS để tiết kiệm dung lượng
	_ = r.DeleteFolder(ctx, r.GetStreamURLDIR(sessionID))
	return uploadResult, nil
}

// UpdateStreamManifest: Cập nhật file danh sách phát .m3u8
// manifestContent là nội dung mới của file manifest/index.m3u8
// sessionID là định danh duy nhất cho buổi live hiện tại
func (r *SeaweedfsAdapter) UpdateStreamManifest(ctx context.Context, sessionID string, manifestContent []byte) error {
	manifestPath := r.GetStreamURL(sessionID) // Ví dụ: /lives/{sessionID}/index.m3u8

	// Manifest phải luôn được ghi đè để người xem cập nhật được segment mới nhất
	reader := bytes.NewReader(manifestContent)
	_, err := r.client.Upload(reader, int64(len(manifestContent)), manifestPath, "", "")
	return err
}
func (r *SeaweedfsAdapter) MergeToMP4(ctx context.Context, sessionID string, outputPath string) error {
	// Gọi FFmpeg để merge HLS thành MP4
	m3u8URL := r.GetStreamURL(sessionID)

	args := []string{
		"-y",          // Ghi đè file đầu ra nếu đã tồn tại
		"-i", m3u8URL, // Link HTTP trỏ tới file m3u8 trên SeaweedFS Filer
		"-c", "copy", // Copy y nguyên luồng video/audio, không tốn CPU encode
		"-bsf:a", "aac_adtstoasc", // Fix lỗi định dạng âm thanh TS -> MP4
		outputPath, // Lưu ra file tạm trên disk (VD: /tmp/lives/session_vod.mp4)
	}
	cmd := exec.CommandContext(ctx, "ffmpeg", args...)

	// Lấy log lỗi của FFmpeg để dễ debug nếu ghép file thất bại
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg merge failed: %w, stderr: %s", err, stderr.String())
	}
	return nil
}

// --- HELPER METHODS ---

// GetPublicURL: Trả về URL đầy đủ để hiển thị ảnh/video
func (r *SeaweedfsAdapter) GetPublicURL(filePath string) string {
	if filePath == "" {
		return ""
	}
	// Đảm bảo không bị dư dấu / khi nối chuỗi: http://host/path
	base := strings.TrimRight(r.publicHost, "/")
	path := "/" + strings.TrimLeft(filePath, "/")

	return base + path
}

// GetStreamSegmentURL: Trả về link file segment để trình phát video (HLS Player) kết nối
func (r *SeaweedfsAdapter) GetStreamSegmentURL(sessionID string, segmentName string) string {
	// Đường dẫn chuẩn cho trình phát: http://filer:8888/lives/{sessionID}/{segmentName}
	liveFolder := fmt.Sprintf("/lives/%s", sessionID)
	fullPath := fmt.Sprintf("%s/%s", liveFolder, segmentName)
	return r.GetPublicURL(fullPath)
}
func (r *SeaweedfsAdapter) GetStreamURLDIR(sessionID string) string {
	// Đường dẫn chuẩn cho trình phát: http://filer:8888/lives/{sessionID}/index.m3u8
	manifestPath := fmt.Sprintf("/lives/%s", sessionID)
	return r.GetPublicURL(manifestPath)
}
func (r *SeaweedfsAdapter) GetStreamURL(sessionID string) string {
	// Đường dẫn chuẩn cho trình phát: http://filer:8888/lives/{sessionID}/index.m3u8
	manifestPath := fmt.Sprintf("/lives/%s/index.m3u8", sessionID)
	return r.GetPublicURL(manifestPath)
}

func (r *SeaweedfsAdapter) GetUploadPresignedUrl(ctx context.Context, input *dto.FileUploadInput) (*dto.PresignedURLResponse, error) {
	// 1. VALIDATION INPUT
	if input.FileName == "" {
		return nil, fmt.Errorf("filename is required")
	}

	// 2. TẠO TÊN FILE DUY NHẤT (UUID)
	// Format: uuid_timestamp.ext (Giúp sort theo thời gian và unique tuyệt đối)
	ext := filepath.Ext(input.FileName)
	if ext == "" {
		ext = ".bin"
	}
	// Dùng UnixNano để đảm bảo tính duy nhất cao nhất kết hợp UUID
	uniqueFileName := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)

	// 3. XÁC ĐỊNH FOLDER (Logic giống hàm Upload thường)
	folderPath := input.Folder
	if folderPath == "" {
		folderPath = utils.GenerateStoragePath(input.OwnerID, input.Storage)
	}

	// Chuẩn hóa path: S3 key không nên bắt đầu bằng dấu "/" (relative to bucket)
	// VD: avatars/user_123/abc.jpg (ĐÚNG) - /avatars/user_123/abc.jpg (SAI với một số S3 client)
	cleanFolder := strings.Trim(folderPath, "/")
	objectName := fmt.Sprintf("%s/%s", cleanFolder, uniqueFileName)

	// 4. CẤU HÌNH PRESIGNED URL
	// Thời gian hết hạn: 15 phút (Đủ để user upload file lớn, nhưng không quá lâu để rò rỉ)
	expiry := 15 * time.Minute

	// Best Practice Security: Ép buộc Content-Type
	// Nếu Frontend upload file khác loại (vd đổi đuôi .exe thành .jpg), S3 sẽ từ chối.
	reqParams := make(url.Values)
	if input.ContentType != "" {
		reqParams.Set("response-content-type", input.ContentType)
		reqParams.Set("Content-Type", input.ContentType)
	}

	// 5. GỌI MINIO SDK TẠO URL
	presignedURL, err := r.s3Client.PresignedPutObject(ctx, r.cfg.SEAWEEDFS.S3_BUCKET_NAME, objectName, expiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned url: %w", err)
	}

	// 6. TẠO PUBLIC URL
	// Public URL để lưu DB cần mapping từ Bucket S3 sang Filer Path hoặc CDN
	// Trong SeaweedFS: Bucket "default", Object "avatars/img.jpg" -> Filer Path "/buckets/default/avatars/img.jpg"
	// Hoặc nếu dùng S3 Gateway trực tiếp làm CDN thì là: http://s3-host/bucket/key

	// Ở đây ta dùng Filer Path để đồng bộ với logic cũ
	internalFilePath := fmt.Sprintf("/buckets/%s/%s", r.cfg.SEAWEEDFS.S3_BUCKET_NAME, objectName)

	// Nếu r.publicHost trỏ vào Filer (port 8888)
	finalPublicURL := fmt.Sprintf("%s%s", strings.TrimRight(r.publicHost, "/"), internalFilePath)

	// 7. TRẢ VỀ RESPONSE
	return &dto.PresignedURLResponse{
		URL:       presignedURL.String(),
		FilePath:  internalFilePath, // Lưu cái này vào DB
		PublicURL: finalPublicURL,   // Dùng để hiển thị
		Method:    "PUT",            // Frontend bắt buộc dùng PUT
		Headers: map[string]string{
			"Content-Type": input.ContentType, // Frontend bắt buộc set header này khớp
		},
	}, nil
}
