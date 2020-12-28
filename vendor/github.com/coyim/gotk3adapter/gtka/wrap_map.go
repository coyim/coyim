package gtka

import (
	"fmt"

	"github.com/coyim/gotk3adapter/gliba"
	"github.com/coyim/gotk3extra"
	"github.com/gotk3/gotk3/gtk"
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
		val := WrapAboutDialogSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.AccelGroup:
		val := WrapAccelGroupSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Adjustment:
		val := WrapAdjustmentSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.ApplicationWindow:
		val := WrapApplicationWindowSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Assistant:
		val := WrapAssistantSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Bin:
		val := WrapBinSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Box:
		val := WrapBoxSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Button:
		val := WrapButtonSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.ModelButton:
		val := WrapModelButtonSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.LinkButton:
		val := WrapLinkButtonSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.CellRenderer:
		val := WrapCellRendererSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.CellRendererText:
		val := WrapCellRendererTextSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.CellRendererToggle:
		val := WrapCellRendererToggleSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.CheckButton:
		val := WrapCheckButtonSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.CheckMenuItem:
		val := WrapCheckMenuItemSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.ComboBox:
		val := WrapComboBoxSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.ComboBoxText:
		val := WrapComboBoxTextSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Container:
		val := WrapContainerSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Dialog:
		val := WrapDialogSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Entry:
		val := WrapEntrySimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.EventBox:
		val := WrapEventBoxSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.ButtonBox:
		val := WrapButtonBoxSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.FileChooserDialog:
		val := WrapFileChooserDialogSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Grid:
		val := WrapGridSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.HeaderBar:
		val := WrapHeaderBarSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Image:
		val := WrapImageSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.InfoBar:
		val := WrapInfoBarSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Label:
		val := WrapLabelSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.ListStore:
		val := WrapListStoreSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Menu:
		val := WrapMenuSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.MenuBar:
		val := WrapMenuBarSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.MenuItem:
		val := WrapMenuItemSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.MenuButton:
		val := WrapMenuButtonSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.MenuShell:
		val := WrapMenuShellSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.MessageDialog:
		val := WrapMessageDialogSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Notebook:
		val := WrapNotebookSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.ProgressBar:
		val := WrapProgressBarSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Revealer:
		val := WrapRevealerSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.ScrolledWindow:
		val := WrapScrolledWindowSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.SearchBar:
		val := WrapSearchBarSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.SearchEntry:
		val := WrapSearchEntrySimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.SeparatorMenuItem:
		val := WrapSeparatorMenuItemSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.SpinButton:
		val := WrapSpinButtonSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Spinner:
		val := WrapSpinnerSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.TextBuffer:
		val := WrapTextBufferSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.TextTag:
		val := WrapTextTagSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.TextTagTable:
		val := WrapTextTagTableSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.TextView:
		val := WrapTextViewSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.ToggleButton:
		val := WrapToggleButtonSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.TreePath:
		val := WrapTreePathSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.TreeSelection:
		val := WrapTreeSelectionSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.TreeStore:
		val := WrapTreeStoreSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.TreeView:
		val := WrapTreeViewSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.TreeViewColumn:
		val := WrapTreeViewColumnSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Widget:
		val := WrapWidgetSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Window:
		val := WrapWindowSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.ToolItem:
		val := WrapToolItemSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.ToolButton:
		val := WrapToolButtonSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gotk3extra.MenuToolButton:
		val := WrapMenuToolButtonSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Overlay:
		val := WrapOverlaySimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Popover:
		val := WrapPopoverSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.RadioButton:
		val := WrapRadioButtonSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gtk.Switch:
		val := WrapSwitchSimple(oo)
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
		val := UnwrapAboutDialog(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *accelGroup:
		val := UnwrapAccelGroup(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *adjustment:
		val := UnwrapAdjustment(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *applicationWindow:
		val := UnwrapApplicationWindow(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *assistant:
		val := UnwrapAssistant(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *bin:
		val := UnwrapBin(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *box:
		val := UnwrapBox(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *button:
		val := UnwrapButton(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *modelButton:
		val := UnwrapModelButton(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *linkButton:
		val := UnwrapLinkButton(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *cellRenderer:
		val := UnwrapCellRenderer(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *cellRendererText:
		val := UnwrapCellRendererText(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *cellRendererToggle:
		val := UnwrapCellRendererToggle(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *checkButton:
		val := UnwrapCheckButton(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *checkMenuItem:
		val := UnwrapCheckMenuItem(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *comboBox:
		val := UnwrapComboBox(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *comboBoxText:
		val := UnwrapComboBoxText(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *container:
		val := UnwrapContainer(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *dialog:
		val := UnwrapDialog(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *entry:
		val := UnwrapEntry(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *eventBox:
		val := UnwrapEventBox(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *buttonBox:
		val := UnwrapButtonBox(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *fileChooserDialog:
		val := UnwrapFileChooserDialog(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *grid:
		val := UnwrapGrid(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *headerBar:
		val := UnwrapHeaderBar(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *image:
		val := UnwrapImage(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *infoBar:
		val := UnwrapInfoBar(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *label:
		val := UnwrapLabel(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *listStore:
		val := UnwrapListStore(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *menu:
		val := UnwrapMenu(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *menuBar:
		val := UnwrapMenuBar(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *menuItem:
		val := UnwrapMenuItem(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *menuButton:
		val := UnwrapMenuButton(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *menuShell:
		val := UnwrapMenuShell(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *messageDialog:
		val := UnwrapMessageDialog(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *notebook:
		val := UnwrapNotebook(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *progressBar:
		val := UnwrapProgressBar(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *revealer:
		val := UnwrapRevealer(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *scrolledWindow:
		val := UnwrapScrolledWindow(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *searchBar:
		val := UnwrapSearchBar(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *searchEntry:
		val := UnwrapSearchEntry(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *separatorMenuItem:
		val := UnwrapSeparatorMenuItem(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *spinButton:
		val := UnwrapSpinButton(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *spinner:
		val := UnwrapSpinner(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *textBuffer:
		val := UnwrapTextBuffer(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *textTag:
		val := UnwrapTextTag(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *textTagTable:
		val := UnwrapTextTagTable(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *textView:
		val := UnwrapTextView(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *toggleButton:
		val := UnwrapToggleButton(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *treePath:
		val := UnwrapTreePath(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *treeSelection:
		val := UnwrapTreeSelection(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *treeStore:
		val := UnwrapTreeStore(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *treeView:
		val := UnwrapTreeView(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *treeViewColumn:
		val := UnwrapTreeViewColumn(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *widget:
		val := UnwrapWidget(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *window:
		val := UnwrapWindow(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *toolItem:
		val := UnwrapToolItem(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *toolButton:
		val := UnwrapToolButton(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *menuToolButton:
		val := UnwrapMenuToolButton(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *overlay:
		val := UnwrapOverlay(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *popover:
		val := UnwrapPopover(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *radioButton:
		val := UnwrapRadioButton(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *zwitch:
		val := UnwrapSwitch(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	default:
		return nil, false
	}
}
