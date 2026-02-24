package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"
)

type GetPostByUserIDUseCase struct {
	postRepo          IRepositoryMongodb.IPostRepository
	postMediaRepo     IRepositoryMongodb.IPostMediaRepository
	postSetting       IRepositoryMongodb.IPostSettingRepository
	postExtensionRepo IRepositoryMongodb.IPostExtensionRepository
	EditLogRepo       IRepositoryMongodb.IPostEditLogsRepository
	pool              IRepositoryShare.IWorkerPool
}

func NewGetPostByUserIDUseCase(postRepo IRepositoryMongodb.IPostRepository,
	postMediaRepo IRepositoryMongodb.IPostMediaRepository,
	postSetting IRepositoryMongodb.IPostSettingRepository,
	postExtensionRepo IRepositoryMongodb.IPostExtensionRepository,
	EditLogRepo IRepositoryMongodb.IPostEditLogsRepository,
	pool IRepositoryShare.IWorkerPool) *GetPostByUserIDUseCase {
	return &GetPostByUserIDUseCase{
		postRepo:          postRepo,
		postMediaRepo:     postMediaRepo,
		postSetting:       postSetting,
		postExtensionRepo: postExtensionRepo,
		pool:              pool,
		EditLogRepo:       EditLogRepo,
	}
}

func (uc *GetPostByUserIDUseCase) Execute(ctx context.Context, req *req.GetPostByUserIDRequest) (*response.Response, error) {
	// 1. Lấy danh sách bài viết của user từ MongoDB với pagination
	reultResponse := &res.GetPostByUserIDReponse{
		PostRes:          []*res.PostRes{},
		PostMediaRes:     []*res.PostMediaRes{},
		PostSettingRes:   []*res.PostSettingRes{},
		PostExtensionRes: []*res.PostExtensionRes{},
		PostEditlog:      []*res.PostEntityEditLogRes{},
	}
	errCh := make(chan error)
	dataPost, err := uc.postRepo.PanigationPosts(ctx, req.UserID, req.Post.Cursor, req.Post.Limit)
	captureDataPost := dataPost.Data.([]*entity.Post)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage(""),
			response.WithStatus("")), err
	}
	resultCh := make(chan *res.GetPostByUserIDReponse, len(captureDataPost))
	// dataPostMedia, err := uc.postMediaRepo.PanigationPostsByUserID(ctx, req.UserID, req.PostMedia.Cursor, req.PostMedia.Limit)
	// if err != nil {
	// 	return response.NewResponse(response.WithData(""),
	// 		response.WithMessage(""),
	// 		response.WithStatus("")), err
	// }
	for _, post := range captureDataPost {
		err := uc.pool.Run(ctx, func() {
			data := mapper.ToResPost(post)
			postMedia, err := uc.postMediaRepo.GetByPostID(ctx, post.ID.Hex())
			if err != nil {
				errCh <- err
				return
			}
			data1 := mapper.ToPostMediaResPostMedia(postMedia)
			postSetting, err := uc.postSetting.GetPostSettingByPostID(ctx, post.ID.Hex())
			if err != nil {
				errCh <- err
				return
			}
			data2 := mapper.ToPostSettingResPostSetting(postSetting)
			postExtension, err := uc.postExtensionRepo.GetByPostID(ctx, post.ID.Hex())
			if err != nil {
				errCh <- err
				return
			}
			data3 := mapper.ToPostExtensionResPostExtension(postExtension)
			postEditlog, err := uc.EditLogRepo.GetByTargetID(ctx, post.ID.Hex())
			if err != nil {
				errCh <- err
				return
			}
			data4 := mapper.ToPostEntityEditLogRes(postEditlog)
			resultCh <- &res.GetPostByUserIDReponse{
				PostRes:          []*res.PostRes{data},
				PostMediaRes:     []*res.PostMediaRes{data1},
				PostSettingRes:   []*res.PostSettingRes{data2},
				PostExtensionRes: []*res.PostExtensionRes{data3},
				PostEditlog:      []*res.PostEntityEditLogRes{data4},
			}
		})
		if err != nil {
			return response.NewResponse(response.WithData(""),
				response.WithMessage(""),
				response.WithStatus("")), err
		}
		for i := 0; i < len(captureDataPost); i++ {
			select {
			case err := <-errCh:
				if err != nil {
					// Nếu 1 trong các goroutine lỗi -> Dừng toàn bộ và trả lỗi luôn
					return nil, err
				}
			case partialRes := <-resultCh:
				if partialRes != nil {
					// Append nối các mảng lại với nhau
					reultResponse.PostRes = append(reultResponse.PostRes, partialRes.PostRes...)
					reultResponse.PostMediaRes = append(reultResponse.PostMediaRes, partialRes.PostMediaRes...)
					reultResponse.PostSettingRes = append(reultResponse.PostSettingRes, partialRes.PostSettingRes...)
					reultResponse.PostExtensionRes = append(reultResponse.PostExtensionRes, partialRes.PostExtensionRes...)
					reultResponse.PostEditlog = append(reultResponse.PostEditlog, partialRes.PostEditlog...)
				}
			}
		}
	}
	return response.NewResponse(response.WithData(reultResponse),
		response.WithMessage(""),
		response.WithStatus("")), nil
}
