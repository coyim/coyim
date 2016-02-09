package gui

import (
	"fmt"
	"log"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

const (
	ulIndexId          = 0
	ulIndexDisplayName = 1
	ulIndexJid         = 2
	ulIndexColor       = 3
	ulIndexBackground  = 4
	ulIndexWeight      = 5
	ulIndexTooltip     = 6
	ulStatusIcon       = 7
)

var ulAllIndexValues = []int{0, 1, 2, 3, 4, 5, 6, 7}

type unifiedLayout struct {
	ui           *gtkUI
	cl           *conversationList
	leftPane     *gtk.Box
	revealer     *gtk.Revealer
	notebook     *gtk.Notebook
	header       *gtk.Label
	headerBox    *gtk.Box
	close        *gtk.Button
	convsVisible bool
	inPageSet    bool
	itemMap      map[int]*conversationStackItem
}

type conversationList struct {
	layout *unifiedLayout
	view   *gtk.TreeView
	model  *gtk.ListStore
}

type conversationStackItem struct {
	*conversationPane
	pageIndex int
	iter      *gtk.TreeIter
	layout    *unifiedLayout
}

func newUnifiedLayout(ui *gtkUI, left, parent *gtk.Box) *unifiedLayout {
	ul := &unifiedLayout{
		ui:       ui,
		cl:       &conversationList{},
		leftPane: left,
		itemMap:  make(map[int]*conversationStackItem),
	}
	ul.cl.layout = ul

	builder := newBuilder("UnifiedLayout")
	builder.getItems(
		"treeview", &ul.cl.view,
		"liststore", &ul.cl.model,

		"revealer", &ul.revealer,
		"notebook", &ul.notebook,
		"header_label", &ul.header,
		"header_box", &ul.headerBox,
		"close_button", &ul.close,
	)
	builder.ConnectSignals(map[string]interface{}{
		"on_activate":    ul.cl.onActivate,
		"on_clicked":     ul.onCloseClicked,
		"on_switch_page": ul.onSwitchPage,
	})
	parent.PackStart(ul.revealer, false, true, 0)
	parent.SetChildPacking(left, false, true, 0, gtk.PACK_START)
	ul.notebook.SetSizeRequest(500, -1)
	ul.revealer.Hide()
	left.SetHAlign(gtk.ALIGN_FILL)
	left.SetHExpand(true)
	return ul
}

func (ul *unifiedLayout) createConversation(account *account, uid string) conversationView {
	cp := createConversationPane(account, uid, ul.ui, &ul.ui.window.Window)
	cp.menubar.Hide()
	idx := ul.notebook.AppendPage(cp.widget, nil)
	if idx < 0 {
		panic("Failed to append page to notebook")
	}

	csi := &conversationStackItem{
		conversationPane: cp,
		pageIndex:        idx,
		layout:           ul,
	}

	ul.notebook.SetTabLabelText(cp.widget, csi.shortName())
	ul.itemMap[idx] = csi
	buffer, _ := csi.history.GetBuffer()
	buffer.Connect("changed", func() {
		ul.onConversationChanged(csi)
	})
	return csi
}

func (ul *unifiedLayout) onConversationChanged(csi *conversationStackItem) {
	if ul.notebook.GetCurrentPage() != csi.pageIndex {
		csi.setBold(true)
	}
}

func (cl *conversationList) add(csi *conversationStackItem) {
	if csi.iter == nil {
		csi.iter = cl.model.Append()
		cl.updateItem(csi)
	}
}

func (cl *conversationList) remove(csi *conversationStackItem) {
	if csi.iter == nil {
		return
	}
	cl.model.Remove(csi.iter)
	csi.iter = nil
}

func (cl *conversationList) updateItem(csi *conversationStackItem) {
	peer, ok := cl.layout.ui.getPeer(csi.account, csi.to)
	if !ok {
		log.Printf("No peer found for %s", csi.to)
		return
	}
	cl.model.Set(csi.iter, ulAllIndexValues, []interface{}{
		csi.pageIndex,
		csi.shortName(),
		peer.Jid,
		decideColorFor(peer),
		"#ffffff",
		500,
		createTooltipFor(peer),
		statusIcons[decideStatusFor(peer)].getPixbuf(),
	},
	)
}

func (ul *unifiedLayout) showConversations() {
	if ul.convsVisible {
		return
	}
	ul.leftPane.SetHExpand(false)
	ul.revealer.Show()
	ul.revealer.SetHExpand(true)
	ul.revealer.SetRevealChild(true)

	ul.convsVisible = true
}

func (ul *unifiedLayout) hideConversations() {
	if !ul.convsVisible {
		return
	}
	width := ul.leftPane.GetAllocatedWidth()
	height := ul.ui.window.GetAllocatedHeight()
	ul.revealer.SetRevealChild(false)
	ul.revealer.SetHExpand(false)
	ul.revealer.Hide()
	ul.leftPane.SetHExpand(true)
	ul.ui.window.Resize(width, height)
	ul.convsVisible = false
}

func (csi *conversationStackItem) SetEnabled(enabled bool) {
	log.Printf("csi.SetEnabled(%v)", enabled)
}

func (csi *conversationStackItem) shortName() string {
	ss := strings.Split(csi.to, "@")
	return ss[0]
}

func (csi *conversationStackItem) setBold(bold bool) {
	if csi.iter == nil {
		return
	}
	weight := 500
	if bold {
		weight = 700
	}
	csi.layout.cl.model.SetValue(csi.iter, ulIndexWeight, weight)
}

func (csi *conversationStackItem) Show(userInitiated bool) {
	csi.layout.showConversations()
	csi.updateSecurityWarning()
	csi.layout.cl.add(csi)
	csi.widget.Show()
	if userInitiated {
		csi.bringToFront()
		return
	}
	if csi.layout.notebook.GetCurrentPage() != csi.pageIndex {
		csi.setBold(true)
	}
}

func (csi *conversationStackItem) bringToFront() {
	csi.layout.showConversations()
	csi.setBold(false)
	csi.layout.setCurrentPage(csi)
	csi.layout.header.SetText(csi.to)
	csi.entry.GrabFocus()
}

func (csi *conversationStackItem) remove() {
	csi.layout.cl.remove(csi)
	csi.widget.Hide()
}

func (cl *conversationList) getItemForIter(iter *gtk.TreeIter) *conversationStackItem {
	val, err := cl.model.GetValue(iter, ulIndexId)
	if err != nil {
		log.Printf("Error getting ulIndexId value: %v", err)
		return nil
	}
	gv, err := val.GoValue()
	if err != nil {
		fmt.Printf("Error getting GoValue for ulIndexId: %v", err)
		return nil
	}
	return cl.layout.itemMap[gv.(int)]
}

func (cl *conversationList) onActivate(v *gtk.TreeView, path *gtk.TreePath) {
	iter, err := cl.model.GetIter(path)
	if err != nil {
		log.Printf("Error converting path to iter: %v", err)
		return
	}
	csi := cl.getItemForIter(iter)
	if csi != nil {
		csi.bringToFront()
	}
}

func (cl *conversationList) removeSelection() {
	selection, _ := cl.view.GetSelection()
	selection.GetSelectedRows(cl.model).Foreach(func(item interface{}) {
		selection.UnselectPath(item.(*gtk.TreePath))
	})
}

func (ul *unifiedLayout) setCurrentPage(csi *conversationStackItem) {
	ul.inPageSet = true
	ul.notebook.SetCurrentPage(csi.pageIndex)
	ul.inPageSet = false
}

func (ul *unifiedLayout) onCloseClicked() {
	page := ul.notebook.GetCurrentPage()
	if page < 0 {
		return
	}
	item := ul.itemMap[page]
	if item != nil {
		item.remove()
	}

	if !ul.displayFirstConvo() {
		ul.header.SetText("")
		ul.hideConversations()
	}
}

func (ul *unifiedLayout) onSwitchPage(notebook *gtk.Notebook, page *gtk.Widget, idx int) {
	if ul.inPageSet {
		return
	}
	if csi := ul.itemMap[idx]; csi != nil {
		csi.bringToFront()
	}
	ul.cl.removeSelection()
}

func (ul *unifiedLayout) displayFirstConvo() bool {
	if iter, ok := ul.cl.model.GetIterFirst(); ok {
		if csi := ul.cl.getItemForIter(iter); csi != nil {
			csi.bringToFront()
			return true
		}
	}
	return false
}

func (ul *unifiedLayout) update() {
	for it, ok := ul.cl.model.GetIterFirst(); ok; ok = ul.cl.model.IterNext(it) {
		csi := ul.cl.getItemForIter(it)
		if csi != nil {
			ul.cl.updateItem(csi)
		}
	}
}
