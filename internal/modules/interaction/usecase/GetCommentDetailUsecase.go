package usecase

import (
	"context"
	"net/http"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/IRepository/IRepositoryMongoDB"
)

type GetCommentDetailUsecase struct {
	commentRepo  IRepositoryMongoDB.ICommentRepository
	editLogRepo  IRepositoryMongoDB.ICommentEditLogsRepository
	reactionRepo IRepositoryCassandra.IReactionsRepository
	pool         IRepositoryShare.IWorkerPool
}

func NewGetCommentDetailUsecase(commentRepo IRepositoryMongoDB.ICommentRepository, editLogRepo IRepositoryMongoDB.ICommentEditLogsRepository, reactionRepo IRepositoryCassandra.IReactionsRepository, pool IRepositoryShare.IWorkerPool) *GetCommentDetailUsecase {
	return &GetCommentDetailUsecase{
		commentRepo:  commentRepo,
		editLogRepo:  editLogRepo,
		reactionRepo: reactionRepo,
		pool:         pool,
	}
}

func (u *GetCommentDetailUsecase) Execute(ctx context.Context, req *req.GetCommentDetailRequest) (*response.Response, error) {
	// Implement the logic for getting comment details here
	type taskResult struct {
		dataComment   *res.CommentRes
		dataEditLog   []*res.CommentEditLogRes
		dataReactions []*res.EntityReactionRes
		err           error
	}
	workerCount := 3
	resultChan := make(chan taskResult, workerCount)

	for i := 0; i < workerCount; i++ {
		err := u.pool.Run(ctx, func() {
			switch i {
			case 0:
				comment, err := u.commentRepo.GetCommentByID(ctx, req.CommentID)
				res := mapper.ToResponseComment(comment)
				resultChan <- taskResult{dataComment: res, err: err}
			case 1:
				editLogs, err := u.editLogRepo.GetVersionBulkEditLogsByTargetID(ctx, req.CommentID)
				var editLogResList []*res.CommentEditLogRes
				for _, log := range editLogs {
					logRes := mapper.ToResCommentEditLogs(log)
					editLogResList = append(editLogResList, logRes)
				}
				resultChan <- taskResult{dataEditLog: editLogResList, err: err}
			case 2:
				reactions, err := u.reactionRepo.GetReactionsByTargetID(ctx, req.CommentID)
				var reactionResList []*res.EntityReactionRes
				for _, reaction := range reactions {
					reactionRes := mapper.FromEntityToResEntityReactions(reaction)
					reactionResList = append(reactionResList, reactionRes)
				}
				resultChan <- taskResult{dataReactions: reactionResList, err: err}
			}
		})
		if err != nil {
			return response.NewResponse(response.WithData(""), response.WithMessage("Failed to start worker"), response.WithStatus(http.StatusInternalServerError)), err
		}

	}
	u.pool.Wait() // Đợi tất cả goroutine hoàn thành
	finalRes := &res.CommentDetailResponse{}
	for i := 0; i < workerCount; i++ {
		result := <-resultChan
		if result.err != nil {
			return nil, result.err // Or handle partial failures
		}

		// Assign based on which field is populated
		if result.dataComment != nil {
			finalRes.Comment = result.dataComment
		}
		if result.dataEditLog != nil {
			finalRes.Edit = result.dataEditLog
		}
		if result.dataReactions != nil {
			finalRes.Reactions = result.dataReactions
		}
	}
	return response.NewResponse(response.WithData(finalRes), response.WithMessage("Comment details retrieved successfully"), response.WithStatus(http.StatusOK)), nil
}
