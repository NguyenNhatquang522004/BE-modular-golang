package dataloader

// import (
// 	"context"
// 	"fmt"
// 	"time"
// 	"your-project/internal/modules/identity/domain"
// 	"your-project/internal/modules/identity/usecase"
//     "github.com/graph-gophers/dataloader" // Library phổ biến cho Go
// )

// // Key context để lưu loader
// type ctxKey string
// const KeyLoaders = ctxKey("dataloaders")

// // Loaders struct chứa tất cả loader
// type Loaders struct {
// 	UserLoader *dataloader.Loader
// }

// // NewLoaders khởi tạo logic batching
// func NewLoaders(userUC usecase.UserUseCase) *Loaders {
// 	// Hàm batch function: Nhận vào 1 mảng keys, trả về 1 mảng results
// 	batchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
// 		var ids []string
// 		for _, key := range keys {
// 			ids = append(ids, key.String())
// 		}

// 		// Gọi sang Module Identity để lấy data 1 lần duy nhất (WHERE IN)
// 		userMap, err := userUC.GetUsersByIDs(ctx, ids)
		
// 		results := make([]*dataloader.Result, len(keys))
// 		for i, key := range keys {
// 			user, ok := userMap[key.String()]
// 			if err != nil {
// 				results[i] = &dataloader.Result{Error: err}
// 			} else if !ok {
// 				results[i] = &dataloader.Result{Error: fmt.Errorf("user not found")}
// 			} else {
// 				results[i] = &dataloader.Result{Data: user} // Mapping domain user to result
// 			}
// 		}
// 		return results
// 	}

// 	return &Loaders{
// 		UserLoader: dataloader.NewBatchedLoader(batchFn, dataloader.WithWait(2*time.Millisecond)),
// 	}
// }

// // Middleware để inject Loader vào Context
// func Middleware(loaders *Loaders) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			ctx := context.WithValue(r.Context(), KeyLoaders, loaders)
// 			r = r.WithContext(ctx)
// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }

// // Helper để lấy Loader từ Context ra dùng
// func GetUser(ctx context.Context, userID string) (*domain.User, error) {
// 	loaders, ok := ctx.Value(KeyLoaders).(*Loaders)
// 	if !ok {
// 		return nil, fmt.Errorf("loaders not found in context")
// 	}
// 	// Thunk: Chờ data trả về
// 	thunk := loaders.UserLoader.Load(ctx, dataloader.StringKey(userID))
// 	result, err := thunk()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return result.(*domain.User), nil
// }