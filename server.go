package main

import (
	"fmt"
	"net/http"

	"github.com/cmd-ctrl-q/go-mux-crash-course/repository"

	"github.com/cmd-ctrl-q/go-mux-crash-course/controller"
	router "github.com/cmd-ctrl-q/go-mux-crash-course/http"
	"github.com/cmd-ctrl-q/go-mux-crash-course/service"
)

var (
	// postRepository repository.PostRepository = repository.NewSQLiteRepository("posts")          // database
	// postRepository repository.PostRepository = repository.NewFirestoreRepository()       // database
	postRepository repository.PostRepository = repository.NewPostgresRepository()        // database
	postService    service.PostService       = service.NewPostService(postRepository)    // service
	postController controller.PostController = controller.NewPostController(postService) // controller
	httpRouter     router.Router             = router.NewMuxRouter()                     // router / server / mux
	// httpRouter     router.Router             = router.NewChiRouter()                  	// router / server / mux
)

func main() {
	const PORT string = ":8080"

	httpRouter.GET("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Up and running...")
	})

	httpRouter.GET("/posts", postController.GetAllPosts)
	httpRouter.POST("/posts", postController.AddPost)

	httpRouter.SERVE(PORT)

}
