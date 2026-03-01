package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/IRepository/IRepositoryMongodb"
)

type DeleteGroup struct {
	groupRepo       IRepositoryMongodb.IGroupRepository
	groupMemberRepo IRepositoryMongodb.IGroupMembersRepository
	groupfileRepo   IRepositoryMongodb.IGroupFilesRepository
	groupqaRepo     IRepositoryMongodb.IGroupjoinQuestionsRepository
	groupEventrepo  IRepositoryMongodb.IGroupEventsRepository
	pool            IRepositoryShare.IWorkerPool
}

func NewDeleteGroup(groupRepo IRepositoryMongodb.IGroupRepository, groupMemberRepo IRepositoryMongodb.IGroupMembersRepository, groupfileRepo IRepositoryMongodb.IGroupFilesRepository, groupqaRepo IRepositoryMongodb.IGroupjoinQuestionsRepository, groupEventrepo IRepositoryMongodb.IGroupEventsRepository, pool IRepositoryShare.IWorkerPool) *DeleteGroup {
	return &DeleteGroup{
		groupRepo:       groupRepo,
		groupMemberRepo: groupMemberRepo,
		groupfileRepo:   groupfileRepo,
		groupqaRepo:     groupqaRepo,
		groupEventrepo:  groupEventrepo,
		pool:            pool,
	}
}

func (uc *DeleteGroup) Execute(ctx context.Context, req *req.DeleteGroupRequest) ([]*res.FailedGroup, error) {
	workercount := 5
	results := make(chan *res.FailedGroup, workercount)
	for i := 0; i < workercount; i++ {
		err := uc.pool.Run(ctx, func() {
			switch i {
			case 0:
				err := uc.groupRepo.DeleteGroup(ctx, req.GroupID)
				if err != nil {
					results <- &res.FailedGroup{
						GroupID:      req.GroupID,
						ErrorMessage: err,
					}
					return
				}
				results <- &res.FailedGroup{
					GroupID:      req.GroupID,
					ErrorMessage: nil,
				}
			case 1:
				err := uc.groupMemberRepo.DeleteGroupMember(ctx, req.GroupID)
				if err != nil {
					results <- &res.FailedGroup{
						GroupID:      req.GroupID,
						ErrorMessage: err,
					}
					return
				}
				results <- &res.FailedGroup{
					GroupID:      req.GroupID,
					ErrorMessage: nil,
				}

			case 2:
				err := uc.groupfileRepo.DeleteGroupFilesByGroupID(ctx, req.GroupID)
				if err != nil {
					results <- &res.FailedGroup{
						GroupID:      req.GroupID,
						ErrorMessage: err,
					}
					return
				}
				results <- &res.FailedGroup{
					GroupID:      req.GroupID,
					ErrorMessage: nil,
				}

			case 3:
				_, err := uc.groupqaRepo.DeleteGroupJoinQuestionsByGroupID(ctx, req.GroupID)
				if err != nil {
					results <- &res.FailedGroup{
						GroupID:      req.GroupID,
						ErrorMessage: err,
					}
					return
				}
				results <- &res.FailedGroup{
					GroupID:      req.GroupID,
					ErrorMessage: nil,
				}
				return

			case 4:
				err := uc.groupEventrepo.DeleteGroupEventsByGroupID(ctx, req.GroupID)
				if err != nil {
					results <- &res.FailedGroup{
						GroupID:      req.GroupID,
						ErrorMessage: err,
					}
					return
				}
				results <- &res.FailedGroup{
					GroupID:      req.GroupID,
					ErrorMessage: nil,
				}
			}
		})
		if err != nil {
			return nil, err
		}

	}
	uc.pool.Wait()
	close(results)
	var failedGroups []*res.FailedGroup
	for res := range results {
		if res.ErrorMessage != nil {
			failedGroups = append(failedGroups, res)
		}
	}
	if len(failedGroups) > 0 {
		return failedGroups, nil
	}
	return failedGroups, nil
}
