package graphql

// // Setup Resolver (cần inject các UseCase vào Resolver)
// func InitializeResolver() (*resolver.Resolver, error) {
// 	wire.Build(
// 		identitySet,
// 		socialSet,
//         // Struct Resolver chứa các usecases
// 		wire.Struct(new(resolver.Resolver), "IdentityUserUseCase", "SocialPostUseCase"), 
// 	)
// 	return &resolver.Resolver{}, nil
// }

// // Setup DataLoader (Cần UserUseCase)
// func InitializeLoaders() (*dataloader.Loaders, error) {
//     wire.Build(
//         identitySet,
//         dataloader.NewLoaders,
//     )
//     return &dataloader.Loaders{}, nil
// }

// func main() {
//     // 1. Init Dependencies qua Wire
//     resolvers, err := InitializeResolver() 
//     if err != nil { log.Fatal(err) }
    
//     loaders, err := InitializeLoaders()
//     if err != nil { log.Fatal(err) }

//     // 2. Setup GQL Server
// 	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolvers}))

//     // 3. Setup Router (dùng Chi, Gin, hoặc http standard)
//     // QUAN TRỌNG: Wrap Middleware DataLoader
// 	http.Handle("/query", dataloader.Middleware(loaders)(srv))
	
// 	http.Handle("/", playground.Handler("GraphQL playground", "/query"))

// 	log.Printf("connect to http://localhost:8080/ for GraphQL playground")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }