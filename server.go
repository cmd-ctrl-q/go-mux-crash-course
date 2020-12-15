package main

import (
	"os"

	"github.com/cmd-ctrl-q/go-mux-crash-course/cache"
	"github.com/cmd-ctrl-q/go-mux-crash-course/repository"

	"github.com/cmd-ctrl-q/go-mux-crash-course/controller"
	router "github.com/cmd-ctrl-q/go-mux-crash-course/http"
	"github.com/cmd-ctrl-q/go-mux-crash-course/service"
)

var (
	// postRepository repository.PostRepository = repository.NewSQLiteRepository("posts")		// sqlite
	// postRepository repository.PostRepository = repository.NewPostgresRepository() 			// postgres
	// postRepository repository.PostRepository = repository.NewFirestoreRepository()    		// firestore
	postRepository repository.PostRepository = repository.NewDynamoDBRepository()     // dynamodb
	postService    service.PostService       = service.NewPostService(postRepository) // service
	postCache      cache.PostCache           = cache.NewRedisCache("localhost:6379", 1, 10)
	postController controller.PostController = controller.NewPostController(postService, postCache) // controller
	httpRouter     router.Router             = router.NewMuxRouter()                                // router / server / mux
	// httpRouter     router.Router             = router.NewChiRouter()                  			// router / server / mux
)

func main() {
	// const PORT string = ":8080"

	// httpRouter.GET("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintln(w, "Up and running...")
	// })

	httpRouter.GET("/posts", postController.GetAllPosts)
	httpRouter.GET("/posts/{id}", postController.GetPostByID)
	httpRouter.POST("/posts", postController.AddPost)

	// httpRouter.SERVE(PORT)
	httpRouter.SERVE(":" + os.Getenv("PORT"))

}
