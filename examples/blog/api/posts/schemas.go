package posts

type PutPostSchema struct {
	UUID string `braza:"in=path"`
	Text string
}
