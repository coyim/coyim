GLIB_VERSION=$(shell pkg-config --modversion glib-2.0 | tr . _ | cut -d '_' -f 1-2)
GLIB_BUILD_TAG="glib_$(GLIB_VERSION)"

GTK_VERSION_FULL=$(shell pkg-config --modversion gtk+-3.0)
GTK_VERSION_PATCH=$(shell echo $(GTK_VERSION_FULL) | cut -f3 -d.)
GTK_VERSION=$(shell echo $(GTK_VERSION_FULL) | tr . _ | cut -d '_' -f 1-2)
GTK_BUILD_TAG="gtk_$(GTK_VERSION)"

# All this is necessary to downgrade the gtk version used to 3.22 if the
# 3.24 patch level is lower than 14. The reason for that is that
# a new variable was introduced at 3.24.14, and older patch levels
# won't compile with gotk3

GTK_VERSION_PATCH_LESS14=$(shell expr $(GTK_VERSION_PATCH) \< 14)
ifeq ($(GTK_BUILD_TAG),"gtk_3_24")
ifeq ($(GTK_VERSION_PATCH_LESS14),1)
GTK_BUILD_TAG="gtk_3_22"
endif
endif

PANGO_VERSION=$(shell pkg-config --modversion pango | tr . _ | cut -d '_' -f 1-2)
PANGO_BUILD_TAG="pango_$(PANGO_VERSION)"

CAIRO_VERSION=$(shell pkg-config --modversion cairo | tr . _ | cut -d '_' -f 1-2)
CAIRO_BUILD_TAG="cairo_$(CAIRO_VERSION)"

TAGS := -tags $(GLIB_BUILD_TAG),$(GTK_BUILD_TAG),$(PANGO_BUILD_TAG),$(CAIRO_BUILD_TAG)

GO := go
GOBUILD := $(GO) build

default: build

build:
	$(GOBUILD) $(TAGS)
