package gtki

type Assistant interface {
	Window
	Commit()
	NextPage()
	PreviousPage()
	SetCurrentPage(pageNum int)
	GetCurrentPage() int
	// GetNPages() int
	GetNthPage(pageNum int) (Widget, error)
	// PrependPage(page Widget) int
	AppendPage(page Widget) int
	// InsertPage(page Widget, position int) int
	// RemovePage(pageNum int)
	SetPageType(page Widget, ptype AssistantPageType)
	GetPageType(page Widget) AssistantPageType
	SetPageTitle(page Widget, title string)
	// GetPageTitle(page Widget) string
	SetPageComplete(page Widget, complete bool)
	GetPageComplete(page Widget) bool
	AddActionWidget(child Widget)
	RemoveActionWidget(child Widget)
	UpdateButtonsState()

	// The following are suppossed to be helper methods to work with the assistant.
	// They are not part of the main GTK api.
	GetButtons() []Button
	GetButtonSizeGroup() (SizeGroup, error)
	GetHeaderBar() (HeaderBar, error)
	GetSidebar() (Box, error)
	GetNotebook() (Notebook, error)
	HideBottomActionArea()
}

func AssertAssistant(_ Assistant) {}
