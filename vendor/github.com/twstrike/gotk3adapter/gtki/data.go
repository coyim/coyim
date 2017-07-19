package gtki

// AccelFlags is a representation of GTK's GtkAccelFlags
type AccelFlags int

var (
	ACCEL_VISIBLE AccelFlags
	ACCEL_LOCKED  AccelFlags
	ACCEL_MASK    AccelFlags
)

// Align is a representation of GTK's GtkAlign.
type Align int

var (
	ALIGN_FILL   Align
	ALIGN_START  Align
	ALIGN_END    Align
	ALIGN_CENTER Align
)

// AssistantPageType is a representation of GTK's GtkAssistantPageType.
type AssistantPageType int

var (
	ASSISTANT_PAGE_CONTENT  AssistantPageType
	ASSISTANT_PAGE_INTRO    AssistantPageType
	ASSISTANT_PAGE_CONFIRM  AssistantPageType
	ASSISTANT_PAGE_SUMMARY  AssistantPageType
	ASSISTANT_PAGE_PROGRESS AssistantPageType
	ASSISTANT_PAGE_CUSTOM   AssistantPageType
)

// FileChooserAction is a representation of GTK's GtkFileChooserAction.
type FileChooserAction int

var (
	FILE_CHOOSER_ACTION_OPEN          FileChooserAction
	FILE_CHOOSER_ACTION_SAVE          FileChooserAction
	FILE_CHOOSER_ACTION_SELECT_FOLDER FileChooserAction
	FILE_CHOOSER_ACTION_CREATE_FOLDER FileChooserAction
)

type IconSize int

var (
	ICON_SIZE_INVALID       IconSize
	ICON_SIZE_MENU          IconSize
	ICON_SIZE_SMALL_TOOLBAR IconSize
	ICON_SIZE_LARGE_TOOLBAR IconSize
	ICON_SIZE_BUTTON        IconSize
	ICON_SIZE_DND           IconSize
	ICON_SIZE_DIALOG        IconSize
)

// PackType is a representation of GTK's GtkPackType.
type PackType int

var (
	PACK_START PackType
	PACK_END   PackType
)

// ResponseType is a representation of GTK's GtkResponseType.
type ResponseType int

var (
	RESPONSE_NONE         ResponseType
	RESPONSE_REJECT       ResponseType
	RESPONSE_ACCEPT       ResponseType
	RESPONSE_DELETE_EVENT ResponseType
	RESPONSE_OK           ResponseType
	RESPONSE_CANCEL       ResponseType
	RESPONSE_CLOSE        ResponseType
	RESPONSE_YES          ResponseType
	RESPONSE_NO           ResponseType
	RESPONSE_APPLY        ResponseType
	RESPONSE_HELP         ResponseType
)

// StateFlags is a representation of GTK's GtkStateFlags.
type StateFlags int

var (
	STATE_FLAG_NORMAL       StateFlags
	STATE_FLAG_ACTIVE       StateFlags
	STATE_FLAG_PRELIGHT     StateFlags
	STATE_FLAG_SELECTED     StateFlags
	STATE_FLAG_INSENSITIVE  StateFlags
	STATE_FLAG_INCONSISTENT StateFlags
	STATE_FLAG_FOCUSED      StateFlags
	STATE_FLAG_BACKDROP     StateFlags
)

// StyleProviderPriority is a representation of GTK's GtkStyleProviderPriority.
type StyleProviderPriority int

var (
	STYLE_PROVIDER_PRIORITY_FALLBACK    StyleProviderPriority
	STYLE_PROVIDER_PRIORITY_THEME       StyleProviderPriority
	STYLE_PROVIDER_PRIORITY_SETTINGS    StyleProviderPriority
	STYLE_PROVIDER_PRIORITY_APPLICATION StyleProviderPriority
	STYLE_PROVIDER_PRIORITY_USER        StyleProviderPriority
)
