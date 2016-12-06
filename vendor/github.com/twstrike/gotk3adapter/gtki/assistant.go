package gtki

type Assistant interface {
	Window
	GetCurrentPage() int
	// SetCurrentPage(pageNum int)
	// GetNPages() int
	GetNthPage(pageNum int) (Widget, error)
	// PrependPage(page Widget) int
	AppendPage(page Widget) int
	// InsertPage(page Widget, position int) int
	// RemovePage(pageNum int)
	// SetPageType(page Widget, ptype gtk.AssistantPageType)
	// GetPageType(page Widget) gtk.AssistantPageType
	// SetPageTitle(page Widget, title string)
	// GetPageTitle(page Widget) string
	SetPageComplete(page Widget, complete bool)
	GetPageComplete(page Widget) bool
	// AddActionWidget(child Widget)
	// RemoveActionWidget(child Widget)
	// UpdateButtonsState()
	Commit()
	NextPage()
	PreviousPage()
}

func AssertAssistant(_ Assistant) {}
