package stalecucumber

type globalSentinel struct{
	Package string 
	Name string
}

type instanceSentinel struct{
	Package string 
	Name string
	Args []interface{}
}

