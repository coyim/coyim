package gtka

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gliba"
)

func init() {
	gliba.AddWrapper(WrapLocal)

	gliba.AddUnwrapper(UnwrapLocal)
}

func Wrap(o interface{}) interface{} {
	v1, ok := WrapLocal(o)
	if !ok {
		panic(fmt.Sprintf("Unrecognized type of object: %#v", o))
	}
	return v1
}

func Unwrap(o interface{}) interface{} {
	v1, ok := UnwrapLocal(o)
	if !ok {
		panic(fmt.Sprintf("Unrecognized type of object: %#v", o))
	}
	return v1
}

func WrapLocal(o interface{}) (interface{}, bool) {
	switch oo := o.(type) {
	case *gtk.AboutDialog:
		val := wrapAboutDialogSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.AccelGroup:
		val := wrapAccelGroupSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Adjustment:
		val := wrapAdjustmentSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.ApplicationWindow:
		val := wrapApplicationWindowSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Assistant:
		val := wrapAssistantSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Bin:
		val := wrapBinSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Box:
		val := wrapBoxSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Button:
		val := wrapButtonSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.CellRenderer:
		val := wrapCellRendererSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.CellRendererText:
		val := wrapCellRendererTextSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.CellRendererToggle:
		val := wrapCellRendererToggleSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.CheckButton:
		val := wrapCheckButtonSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.CheckMenuItem:
		val := wrapCheckMenuItemSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.ComboBox:
		val := wrapComboBoxSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.ComboBoxText:
		val := wrapComboBoxTextSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Container:
		val := wrapContainerSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Dialog:
		val := wrapDialogSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Entry:
		val := wrapEntrySimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.EventBox:
		val := wrapEventBoxSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.FileChooserDialog:
		val := wrapFileChooserDialogSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Grid:
		val := wrapGridSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.HeaderBar:
		val := wrapHeaderBarSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Image:
		val := wrapImageSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.InfoBar:
		val := wrapInfoBarSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Label:
		val := wrapLabelSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.ListStore:
		val := wrapListStoreSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Menu:
		val := wrapMenuSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.MenuBar:
		val := wrapMenuBarSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.MenuItem:
		val := wrapMenuItemSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.MenuShell:
		val := wrapMenuShellSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.MessageDialog:
		val := wrapMessageDialogSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Notebook:
		val := wrapNotebookSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.ProgressBar:
		val := wrapProgressBarSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Revealer:
		val := wrapRevealerSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.ScrolledWindow:
		val := wrapScrolledWindowSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.SearchBar:
		val := wrapSearchBarSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.SearchEntry:
		val := wrapSearchEntrySimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.SeparatorMenuItem:
		val := wrapSeparatorMenuItemSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.SpinButton:
		val := wrapSpinButtonSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Spinner:
		val := wrapSpinnerSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.TextBuffer:
		val := wrapTextBufferSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.TextTag:
		val := wrapTextTagSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.TextTagTable:
		val := wrapTextTagTableSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.TextView:
		val := wrapTextViewSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.ToggleButton:
		val := wrapToggleButtonSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.TreePath:
		val := wrapTreePathSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.TreeSelection:
		val := wrapTreeSelectionSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.TreeStore:
		val := wrapTreeStoreSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.TreeView:
		val := wrapTreeViewSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.TreeViewColumn:
		val := wrapTreeViewColumnSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Widget:
		val := wrapWidgetSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Window:
		val := wrapWindowSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	default:
		return nil, false
	}
}

func UnwrapLocal(o interface{}) (interface{}, bool) {
	switch oo := o.(type) {
	case *aboutDialog:
		val := unwrapAboutDialog(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *accelGroup:
		val := unwrapAccelGroup(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *adjustment:
		val := unwrapAdjustment(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *applicationWindow:
		val := unwrapApplicationWindow(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *assistant:
		val := unwrapAssistant(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *bin:
		val := unwrapBin(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *box:
		val := unwrapBox(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *button:
		val := unwrapButton(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *cellRenderer:
		val := unwrapCellRenderer(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *cellRendererText:
		val := unwrapCellRendererText(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *cellRendererToggle:
		val := unwrapCellRendererToggle(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *checkButton:
		val := unwrapCheckButton(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *checkMenuItem:
		val := unwrapCheckMenuItem(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *comboBox:
		val := unwrapComboBox(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *comboBoxText:
		val := unwrapComboBoxText(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *container:
		val := unwrapContainer(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *dialog:
		val := unwrapDialog(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *entry:
		val := unwrapEntry(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *eventBox:
		val := unwrapEventBox(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *fileChooserDialog:
		val := unwrapFileChooserDialog(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *grid:
		val := unwrapGrid(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *headerBar:
		val := unwrapHeaderBar(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *image:
		val := unwrapImage(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *infoBar:
		val := unwrapInfoBar(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *label:
		val := unwrapLabel(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *listStore:
		val := unwrapListStore(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *menu:
		val := unwrapMenu(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *menuBar:
		val := unwrapMenuBar(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *menuItem:
		val := unwrapMenuItem(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *menuShell:
		val := unwrapMenuShell(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *messageDialog:
		val := unwrapMessageDialog(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *notebook:
		val := unwrapNotebook(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *progressBar:
		val := unwrapProgressBar(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *revealer:
		val := unwrapRevealer(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *scrolledWindow:
		val := unwrapScrolledWindow(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *searchBar:
		val := unwrapSearchBar(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *searchEntry:
		val := unwrapSearchEntry(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *separatorMenuItem:
		val := unwrapSeparatorMenuItem(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *spinButton:
		val := unwrapSpinButton(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *spinner:
		val := unwrapSpinner(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *textBuffer:
		val := unwrapTextBuffer(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *textTag:
		val := unwrapTextTag(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *textTagTable:
		val := unwrapTextTagTable(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *textView:
		val := unwrapTextView(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *toggleButton:
		val := unwrapToggleButton(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *treePath:
		val := unwrapTreePath(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *treeSelection:
		val := unwrapTreeSelection(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *treeStore:
		val := unwrapTreeStore(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *treeView:
		val := unwrapTreeView(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *treeViewColumn:
		val := unwrapTreeViewColumn(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *widget:
		val := unwrapWidget(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *window:
		val := unwrapWindow(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	default:
		return nil, false
	}
}
