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
func (r *SeaweedfsAdapter) UploadStreamSegment(ctx context.Context, OwnerID string, sessionID string, segmentName string, content io.Reader) error {
	// Quy hoạch folder: /lives/{sessionID}/{segmentName}
	// Lưu ý: Stream segment thường nhỏ nên ta dùng TTL ngắn (ví dụ: 10 phút) để tự dọn rác

	fullPath := r.GetStreamSegmentURL(OwnerID, sessionID, segmentName , utils.BucketGroupStream)
	// Vì đây là stream realtime, ta không cần lưu Metadata phức tạp vào DB chính
	// Dùng TTL "10m" để SeaweedFS tự xóa các segment cũ
	// Lưu ý: Cần biết Size của segment. Nếu không có, bạn phải buffer hoặc dùng chunked upload.
	// Ở đây giả định bạn đã lấy được size từ encoder (ví dụ FFmpeg).
	_, err := r.client.Upload(content, 0, fullPath, "", "")
	return err
}
func (r *SeaweedfsAdapter) GenerateVOD(ctx context.Context, OwnerID string, sessionID string, name string) (*dto.FileUploadOutput, error) {
	// 1. TẠO INTERNAL URL TRỎ VÀO SEAWEEDFS
	// Mặc định Filer chạy ở port 8888 trên server.
	// Nếu app Go và SeaweedFS chạy chung mạng Docker, dùng tên service (vd: http://seaweedfs-filer:8888)
	// 2. TẠO FILE TẠM TRÊN SERVER BACKEND
	localMP4Path := fmt.Sprintf("/tmp/lives/%s_vod.mp4", sessionID)

	log.Printf("[VOD] Bắt đầu gộp video cho session: %s", sessionID)

	// 3. GỌI FFMPEG GỘP VIDEO
	// FFmpeg sẽ tự tải .m3u8 và .ts từ SeaweedFS về, gộp và lưu vào localMP4Path
	err := r.MergeToMP4(ctx, OwnerID, sessionID, localMP4Path)
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
		OwnerID:  OwnerID,
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
	_ = r.DeleteFolder(ctx, r.GetStreamURLDIR(OwnerID, sessionID, utils.BucketGroupStream))
	return uploadResult, nil
}

// UpdateStreamManifest: Cập nhật file danh sách phát .m3u8
// manifestContent là nội dung mới của file manifest/index.m3u8
// sessionID là định danh duy nhất cho buổi live hiện tại
func (r *SeaweedfsAdapter) UpdateStreamManifest(ctx context.Context, OwnerID string, sessionID string, manifestContent []byte) error {
	manifestPath := r.GetStreamURL(OwnerID, sessionID, utils.BucketGroupStream) // Ví dụ: /lives/{sessionID}/index.m3u8

	// Manifest phải luôn được ghi đè để người xem cập nhật được segment mới nhất
	reader := bytes.NewReader(manifestContent)
	_, err := r.client.Upload(reader, int64(len(manifestContent)), manifestPath, "", "")
	return err
}
func (r *SeaweedfsAdapter) MergeToMP4(ctx context.Context, OwnerID string, sessionID string, outputPath string) error {
	// Gọi FFmpeg để merge HLS thành MP4
	m3u8URL := r.GetStreamURL(OwnerID, sessionID, utils.BucketGroupStream)

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
func (r *SeaweedfsAdapter) getInternalLivePath(ownerID string, sessionID string, Storage utils.StorageType) string {
	// Kết quả: users/{ownerID}/lives/YYYY/MM
	basePartition := utils.GenerateStoragePath(ownerID, Storage)

	// Nối thêm sessionID để nhóm các file .ts và .m3u8 của cùng 1 phiên live
	// Kết quả cuối: /users/{ownerID}/lives/YYYY/MM/{sessionID}
	cleanPath := fmt.Sprintf("/%s/%s", strings.Trim(basePartition, "/"), sessionID)
	return cleanPath
}

// GetStreamSegmentURL: Trả về link file segment để trình phát video (HLS Player) kết nối
func (r *SeaweedfsAdapter) GetStreamSegmentURL(ownerID string, sessionID string, segmentName string, Storage utils.StorageType) string {
	liveFolder := r.getInternalLivePath(ownerID, sessionID, Storage)
	fullPath := fmt.Sprintf("%s/%s", liveFolder, segmentName)
	return r.GetPublicURL(fullPath)
}
func (r *SeaweedfsAdapter) GetStreamURLDIR(ownerID string, sessionID string, Storage utils.StorageType) string {
	manifestDir := r.getInternalLivePath(ownerID, sessionID, Storage)
	return r.GetPublicURL(manifestDir)
}
func (r *SeaweedfsAdapter) GetStreamURL(ownerID string, sessionID string, Storage utils.StorageType) string {
	liveFolder := r.getInternalLivePath(ownerID, sessionID, Storage)
	manifestPath := fmt.Sprintf("%s/index.m3u8", liveFolder)
	return r.GetPublicURL(manifestPath)
}

func (r *SeaweedfsAdapter) GetUploadPresignedUrl(ctx context.Context, input *dto.FileUploadInput) (*dto.PresignedURLResponse, error) {
	// 1. VALIDATION CƠ BẢN
	if input.FileName == "" {
		return nil, fmt.Errorf("filename is required")
	}
	if input.OwnerID == "" {
		return nil, fmt.Errorf("owner_id is required")
	}

	// 2. CHUẨN HÓA ĐUÔI FILE & SINH TÊN DUY NHẤT
	// Áp dụng thêm Timestamp (Unix) để đảm bảo không bao giờ trùng lặp ngay cả khi tạo cùng 1 mili-giây
	ext := filepath.Ext(input.FileName)
	if ext == "" {
		ext = ".bin"
	}
	uniqueFileName := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)

	// 3. XÁC ĐỊNH ĐƯỜNG DẪN THƯ MỤC (FOLDER PARTITIONING)
	// Sử dụng hàm utils.GenerateStoragePath để tự động chia thư mục theo thời gian cho các StorageType nặng
	folderPath := input.Folder
	if folderPath == "" {
		folderPath = utils.GenerateStoragePath(input.OwnerID, input.Storage)
	}

	// 4. TẠO S3 OBJECT KEY CHUẨN MỰC
	// S3 Object Key tuyệt đối KHÔNG ĐƯỢC bắt đầu bằng dấu "/", phải là dạng: "users/123/posts/2026/03/abc.jpg"
	cleanFolder := strings.Trim(folderPath, "/")
	objectName := fmt.Sprintf("%s/%s", cleanFolder, uniqueFileName)

	// 5. CẤU HÌNH THỜI GIAN SỐNG & PRESIGNED URL
	expiry := 15 * time.Minute // 15 phút là con số lý tưởng để chống rò rỉ link

	// Gọi SDK của MinIO để cấp link
	presignedURL, err := r.s3Client.PresignedPutObject(ctx, r.cfg.SEAWEEDFS.S3_BUCKET_NAME, objectName, expiry)
	if err != nil {
		return nil, fmt.Errorf("seaweedfs_adapter: failed to generate upload presigned url: %w", err)
	}

	// 6. REWRITE HOST (Xử lý cho môi trường Docker/Microservices)
	// Biến đổi http://seaweedfs-s3:8333 thành https://cdn.mysocial.com
	if parsedPublic, err := url.Parse(r.publicHost); err == nil && parsedPublic.Host != "" {
		presignedURL.Scheme = parsedPublic.Scheme
		presignedURL.Host = parsedPublic.Host
	}

	// 7. MAP TỪ S3 KEY SANG FILER PATH
	// Khi upload qua Gateway S3 của SeaweedFS, file vật lý sẽ nằm ở: /buckets/{bucket_name}/{object_key}
	internalFilePath := fmt.Sprintf("/buckets/%s/%s", r.cfg.SEAWEEDFS.S3_BUCKET_NAME, objectName)
	finalPublicURL := r.GetPublicURL(internalFilePath)

	// 8. TRẢ VỀ PAYLOAD CHO FRONTEND
	return &dto.PresignedURLResponse{
		URL:       presignedURL.String(),
		FilePath:  internalFilePath, // Lưu giá trị này vào CSDL (Postgres/MongoDB)
		PublicURL: finalPublicURL,   // Dùng để Frontend biết sau khi up xong thì truy cập link nào
		Method:    "PUT",            // Báo Frontend phải dùng HTTP PUT
		Headers: map[string]string{
			"Content-Type": input.ContentType, // Ép Frontend phải gắn đúng Content-Type khi PUT lên S3
		},
	}, nil
}

// GetDownloadPresignedUrl: Sinh link tải file trực tiếp với tốc độ cực cao, bypass Backend
func (r *SeaweedfsAdapter) GetDownloadPresignedUrl(ctx context.Context, filePath string, forceDownload bool) (string, error) {
	// 1. KIỂM TRA ĐẦU VÀO
	cleanPath := strings.TrimSpace(filePath)
	if cleanPath == "" {
		return "", fmt.Errorf("invalid empty file path")
	}

	// 2. TRÍCH XUẤT S3 OBJECT KEY TỪ FILER PATH
	// FilePath đang lưu trong DB là dạng: /buckets/{bucket_name}/users/123/posts/2026/03/abc.jpg
	// Ta cần bóc phần "/buckets/{bucket_name}/" đi để trả lại ObjectKey nguyên thủy
	bucketPrefix := fmt.Sprintf("/buckets/%s/", r.cfg.SEAWEEDFS.S3_BUCKET_NAME)
	objectName := strings.TrimPrefix(cleanPath, bucketPrefix)
	objectName = strings.TrimPrefix(objectName, "/")

	if objectName == "" || objectName == cleanPath {
		// Log cảnh báo nếu đường dẫn trong DB không khớp chuẩn S3 Bucket
		log.Printf("[Warning] FilePath không có bucket prefix chuẩn: %s", cleanPath)
		// Fallback: Lấy toàn bộ filepath làm objectKey (cắt "/" ở đầu)
		objectName = strings.TrimPrefix(cleanPath, "/")
	}

	// 3. CẤU HÌNH THỜI GIAN SỐNG
	expiry := 1 * time.Hour // Thời gian cho phép tải (1 tiếng)

	// 4. XỬ LÝ CONTENT-DISPOSITION (ÉP TẢI XUỐNG HAY XEM TRỰC TIẾP)
	reqParams := make(url.Values)
	if forceDownload {
		// Dùng filepath.Base để trích xuất tên file thật (VD: lấy "abc.jpg" từ "users/123/posts/abc.jpg")
		fileName := filepath.Base(objectName)

		// Set header ép trình duyệt mở hộp thoại tải file (Save As...) thay vì mở tab xem ảnh
		disposition := fmt.Sprintf("attachment; filename=\"%s\"", fileName)
		reqParams.Set("response-content-disposition", disposition)
	}

	// 5. GỌI SDK TẠO LINK KÝ MÃ HÓA
	presignedURL, err := r.s3Client.PresignedGetObject(ctx, r.cfg.SEAWEEDFS.S3_BUCKET_NAME, objectName, expiry, reqParams)
	if err != nil {
		return "", fmt.Errorf("seaweedfs_adapter: failed to generate download presigned url: %w", err)
	}

	// 6. REWRITE HOST CHO PUBLIC NETWORK
	// Biến URL nội bộ thành Public URL để Client tải được
	if parsedPublic, err := url.Parse(r.publicHost); err == nil && parsedPublic.Host != "" {
		presignedURL.Scheme = parsedPublic.Scheme
		presignedURL.Host = parsedPublic.Host
	}

	return presignedURL.String(), nil
}
