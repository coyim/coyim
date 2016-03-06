package glibi

type Application interface {
	Object

	Quit()
	Run([]string) int
}

func AssertApplication(_ Application) {}
