package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gtki"
)

func init() {
	gtki.ACCEL_VISIBLE = gtki.AccelFlags(gtk.ACCEL_VISIBLE)
	gtki.ACCEL_LOCKED = gtki.AccelFlags(gtk.ACCEL_LOCKED)
	gtki.ACCEL_MASK = gtki.AccelFlags(gtk.ACCEL_MASK)

	gtki.ALIGN_FILL = gtki.Align(gtk.ALIGN_FILL)
	gtki.ALIGN_START = gtki.Align(gtk.ALIGN_START)
	gtki.ALIGN_END = gtki.Align(gtk.ALIGN_END)
	gtki.ALIGN_CENTER = gtki.Align(gtk.ALIGN_CENTER)

	gtki.ASSISTANT_PAGE_CONTENT = gtki.AssistantPageType(gtk.ASSISTANT_PAGE_CONTENT)
	gtki.ASSISTANT_PAGE_INTRO = gtki.AssistantPageType(gtk.ASSISTANT_PAGE_INTRO)
	gtki.ASSISTANT_PAGE_CONFIRM = gtki.AssistantPageType(gtk.ASSISTANT_PAGE_CONFIRM)
	gtki.ASSISTANT_PAGE_SUMMARY = gtki.AssistantPageType(gtk.ASSISTANT_PAGE_SUMMARY)
	gtki.ASSISTANT_PAGE_PROGRESS = gtki.AssistantPageType(gtk.ASSISTANT_PAGE_PROGRESS)
	gtki.ASSISTANT_PAGE_CUSTOM = gtki.AssistantPageType(gtk.ASSISTANT_PAGE_CUSTOM)

	gtki.ICON_SIZE_INVALID = gtki.IconSize(gtk.ICON_SIZE_INVALID)
	gtki.ICON_SIZE_MENU = gtki.IconSize(gtk.ICON_SIZE_MENU)
	gtki.ICON_SIZE_SMALL_TOOLBAR = gtki.IconSize(gtk.ICON_SIZE_SMALL_TOOLBAR)
	gtki.ICON_SIZE_LARGE_TOOLBAR = gtki.IconSize(gtk.ICON_SIZE_LARGE_TOOLBAR)
	gtki.ICON_SIZE_BUTTON = gtki.IconSize(gtk.ICON_SIZE_BUTTON)
	gtki.ICON_SIZE_DND = gtki.IconSize(gtk.ICON_SIZE_DND)
	gtki.ICON_SIZE_DIALOG = gtki.IconSize(gtk.ICON_SIZE_DIALOG)

	gtki.FILE_CHOOSER_ACTION_OPEN = gtki.FileChooserAction(gtk.FILE_CHOOSER_ACTION_OPEN)
	gtki.FILE_CHOOSER_ACTION_SAVE = gtki.FileChooserAction(gtk.FILE_CHOOSER_ACTION_SAVE)
	gtki.FILE_CHOOSER_ACTION_SELECT_FOLDER = gtki.FileChooserAction(gtk.FILE_CHOOSER_ACTION_SELECT_FOLDER)
	gtki.FILE_CHOOSER_ACTION_CREATE_FOLDER = gtki.FileChooserAction(gtk.FILE_CHOOSER_ACTION_CREATE_FOLDER)

	gtki.PACK_START = gtki.PackType(gtk.PACK_START)
	gtki.PACK_END = gtki.PackType(gtk.PACK_END)

	gtki.RESPONSE_NONE = gtki.ResponseType(gtk.RESPONSE_NONE)
	gtki.RESPONSE_REJECT = gtki.ResponseType(gtk.RESPONSE_REJECT)
	gtki.RESPONSE_ACCEPT = gtki.ResponseType(gtk.RESPONSE_ACCEPT)
	gtki.RESPONSE_DELETE_EVENT = gtki.ResponseType(gtk.RESPONSE_DELETE_EVENT)
	gtki.RESPONSE_OK = gtki.ResponseType(gtk.RESPONSE_OK)
	gtki.RESPONSE_CANCEL = gtki.ResponseType(gtk.RESPONSE_CANCEL)
	gtki.RESPONSE_CLOSE = gtki.ResponseType(gtk.RESPONSE_CLOSE)
	gtki.RESPONSE_YES = gtki.ResponseType(gtk.RESPONSE_YES)
	gtki.RESPONSE_NO = gtki.ResponseType(gtk.RESPONSE_NO)
	gtki.RESPONSE_APPLY = gtki.ResponseType(gtk.RESPONSE_APPLY)
	gtki.RESPONSE_HELP = gtki.ResponseType(gtk.RESPONSE_HELP)

	gtki.STATE_FLAG_NORMAL = gtki.StateFlags(gtk.STATE_FLAG_NORMAL)
	gtki.STATE_FLAG_ACTIVE = gtki.StateFlags(gtk.STATE_FLAG_ACTIVE)
	gtki.STATE_FLAG_PRELIGHT = gtki.StateFlags(gtk.STATE_FLAG_PRELIGHT)
	gtki.STATE_FLAG_SELECTED = gtki.StateFlags(gtk.STATE_FLAG_SELECTED)
	gtki.STATE_FLAG_INSENSITIVE = gtki.StateFlags(gtk.STATE_FLAG_INSENSITIVE)
	gtki.STATE_FLAG_INCONSISTENT = gtki.StateFlags(gtk.STATE_FLAG_INCONSISTENT)
	gtki.STATE_FLAG_FOCUSED = gtki.StateFlags(gtk.STATE_FLAG_FOCUSED)
	gtki.STATE_FLAG_BACKDROP = gtki.StateFlags(gtk.STATE_FLAG_BACKDROP)

	gtki.STYLE_PROVIDER_PRIORITY_FALLBACK = gtki.StyleProviderPriority(gtk.STYLE_PROVIDER_PRIORITY_FALLBACK)
	gtki.STYLE_PROVIDER_PRIORITY_THEME = gtki.StyleProviderPriority(gtk.STYLE_PROVIDER_PRIORITY_THEME)
	gtki.STYLE_PROVIDER_PRIORITY_SETTINGS = gtki.StyleProviderPriority(gtk.STYLE_PROVIDER_PRIORITY_SETTINGS)
	gtki.STYLE_PROVIDER_PRIORITY_APPLICATION = gtki.StyleProviderPriority(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	gtki.STYLE_PROVIDER_PRIORITY_USER = gtki.StyleProviderPriority(gtk.STYLE_PROVIDER_PRIORITY_USER)
}
