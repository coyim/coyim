package gtki

type Notebook interface {
	Container

	AppendPage(Widget, Widget) int
	GetCurrentPage() int
	GetNPages() int
	GetNthPage(int) (Widget, error)
	NextPage()
	PrevPage()
	SetCurrentPage(int)
	SetShowTabs(bool)
	SetTabLabelText(Widget, string)
}

func AssertNotebook(_ Notebook) {}
