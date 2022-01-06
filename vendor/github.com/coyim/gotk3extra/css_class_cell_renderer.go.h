#ifndef CSS_CLASS_CELL_RENDERER_H
#define CSS_CLASS_CELL_RENDERER_H

#include <gtk/gtk.h>

#define TYPE_CSS_CLASS_CELL_RENDERER             (css_class_cell_renderer_get_type())
#define CSS_CLASS_CELL_RENDERER(obj)             (G_TYPE_CHECK_INSTANCE_CAST((obj),  TYPE_CSS_CLASS_CELL_RENDERER, CSSClassCellRenderer))
#define CSS_CLASS_CELL_RENDERER_CLASS(klass)     (G_TYPE_CHECK_CLASS_CAST ((klass),  TYPE_CSS_CLASS_CELL_RENDERER, CSSClassCellRendererClass))
#define IS_CSS_CLASS_CELL_RENDERER(obj)          (G_TYPE_CHECK_INSTANCE_TYPE ((obj), TYPE_CSS_CLASS_CELL_RENDERER))
#define IS_CSS_CLASS_CELL_RENDERER_CLASS(klass)  (G_TYPE_CHECK_CLASS_TYPE ((klass),  TYPE_CSS_CLASS_CELL_RENDERER))
#define CSS_CLASS_CELL_RENDERER_GET_CLASS(obj)   (G_TYPE_INSTANCE_GET_CLASS ((obj),  TYPE_CSS_CLASS_CELL_RENDERER, CSSClassCellRendererClass))


typedef struct _CSSClassCellRenderer CSSClassCellRenderer;
typedef struct _CSSClassCellRendererClass CSSClassCellRendererClass;

struct _CSSClassCellRenderer
{
    GtkCellRenderer   parent;
    gchar *css;
    GtkCellRenderer *real;
};


struct _CSSClassCellRendererClass
{
    GtkCellRendererClass  parent_class;
};


GType                css_class_cell_renderer_get_type (void);

GtkCellRenderer     *css_class_cell_renderer_new (void);

void                 css_class_cell_renderer_set_real (CSSClassCellRenderer *cr, GtkCellRenderer *real);

static void     css_class_cell_renderer_init       (CSSClassCellRenderer      *cr);
static void     css_class_cell_renderer_class_init (CSSClassCellRendererClass *klass);
static void     css_class_cell_renderer_get_property  (GObject                    *object,
                                                       guint                       param_id,
                                                       GValue                     *value,
                                                       GParamSpec                 *pspec);
static void     css_class_cell_renderer_set_property  (GObject                    *object,
                                                       guint                       param_id,
                                                       const GValue               *value,
                                                       GParamSpec                 *pspec);
static void     css_class_cell_renderer_finalize (GObject *gobject);

static void     css_class_cell_renderer_get_size   (GtkCellRenderer            *cell,
                                                    GtkWidget                  *widget,
                                                    const GdkRectangle         *cell_area,
                                                    gint                       *x_offset,
                                                    gint                       *y_offset,
                                                    gint                       *width,
                                                    gint                       *height);

static void     css_class_cell_renderer_render     (GtkCellRenderer            *cell,
                                                    cairo_t         *cr,
                                                    GtkWidget       *widget,
                                                    const GdkRectangle   *background_area,
                                                    const GdkRectangle   *cell_area,
                                                    GtkCellRendererState  flags);

static GtkCellEditable *css_class_cell_renderer_start_editing (GtkCellRenderer      *cell,
                                                               GdkEvent             *event,
                                                               GtkWidget            *widget,
                                                               const gchar          *path,
                                                               const GdkRectangle   *background_area,
                                                               const GdkRectangle   *cell_area,
                                                               GtkCellRendererState  flags);

static void       css_class_cell_renderer_get_preferred_width            (GtkCellRenderer       *cell,
                                                                          GtkWidget             *widget,
                                                                          gint                  *minimal_size,
                                                                          gint                  *natural_size);
static void       css_class_cell_renderer_get_preferred_height           (GtkCellRenderer       *cell,
                                                                          GtkWidget             *widget,
                                                                          gint                  *minimal_size,
                                                                          gint                  *natural_size);
static void       css_class_cell_renderer_get_preferred_height_for_width (GtkCellRenderer       *cell,
                                                                          GtkWidget             *widget,
                                                                          gint                   width,
                                                                          gint                  *minimum_height,
                                                                          gint                  *natural_height);
static void       css_class_cell_renderer_get_aligned_area               (GtkCellRenderer       *cell,
                                                                          GtkWidget             *widget,
                                                                          GtkCellRendererState   flags,
                                                                          const GdkRectangle    *cell_area,
                                                                          GdkRectangle          *aligned_area);

enum
{
  PROP_0,

  // Properties for the CSS Class Renderer
  PROP_CSS,

  // Properties for the Text Cell Renderer
  PROP_TEXT,

  // TODO: implement the rest of these
  PROP_MARKUP,
  PROP_BACKGROUND,
  PROP_FOREGROUND,
  PROP_UNDERLINE,
  PROP_WEIGHT,

  // Properties for the Pixbuf Cell Renderer
  PROP_PIXBUF,

  // Properties for the Spinner Cell Renderer
  PROP_ACTIVE,
};

static   gpointer parent_class;

GType
css_class_cell_renderer_get_type (void)
{
  static GType css_class_cell_renderer_type = 0;

  if (css_class_cell_renderer_type == 0)
  {
    static const GTypeInfo css_class_cell_renderer_info =
    {
      sizeof (CSSClassCellRendererClass),
      NULL,                                                     /* base_init */
      NULL,                                                     /* base_finalize */
      (GClassInitFunc) css_class_cell_renderer_class_init,
      NULL,                                                     /* class_finalize */
      NULL,                                                     /* class_data */
      sizeof (CSSClassCellRenderer),
      0,                                                        /* n_preallocs */
      (GInstanceInitFunc) css_class_cell_renderer_init,
    };

    css_class_cell_renderer_type = g_type_register_static (GTK_TYPE_CELL_RENDERER,
                                                           "CSSClassCellRenderer",
                                                           &css_class_cell_renderer_info,
                                                           0);
  }

  return css_class_cell_renderer_type;
}

static void
css_class_cell_renderer_init (CSSClassCellRenderer *cr)
{
}


static void
css_class_cell_renderer_class_init (CSSClassCellRendererClass *klass)
{
  GtkCellRendererClass *cell_class   = GTK_CELL_RENDERER_CLASS(klass);
  GObjectClass         *object_class = G_OBJECT_CLASS(klass);

  parent_class           = g_type_class_peek_parent (klass);
  object_class->finalize = css_class_cell_renderer_finalize;

  object_class->get_property = css_class_cell_renderer_get_property;
  object_class->set_property = css_class_cell_renderer_set_property;

  cell_class->get_size = css_class_cell_renderer_get_size;
  cell_class->render   = css_class_cell_renderer_render;
  cell_class->start_editing = css_class_cell_renderer_start_editing;
  cell_class->get_preferred_width = css_class_cell_renderer_get_preferred_width;
  cell_class->get_preferred_height = css_class_cell_renderer_get_preferred_height;
  cell_class->get_preferred_height_for_width = css_class_cell_renderer_get_preferred_height_for_width;
  cell_class->get_aligned_area = css_class_cell_renderer_get_aligned_area;

  g_object_class_install_property (object_class,
                                   PROP_CSS,
                                   g_param_spec_string ("css",
                                                        "CSS Class",
                                                        "The CSS class to set",
                                                         NULL,
                                                         G_PARAM_READWRITE));

  g_object_class_install_property (object_class,
                                   PROP_TEXT,
                                   g_param_spec_string ("text",
                                                        "Text",
                                                        "Text to render (for GtkCellRendererText)",
                                                        NULL,
                                                        G_PARAM_READWRITE));
}

static void
css_class_cell_renderer_finalize (GObject *object)
{
  CSSClassCellRenderer  *cr = CSS_CLASS_CELL_RENDERER(object);
  g_free (cr->css);
  if (cr->real)
      g_object_unref (cr->real);

  (* G_OBJECT_CLASS (parent_class)->finalize) (object);
}

static void
css_class_cell_renderer_get_property (GObject    *object,
                                      guint       param_id,
                                      GValue     *value,
                                      GParamSpec *psec) {
    CSSClassCellRenderer *cr;

    cr = CSS_CLASS_CELL_RENDERER(object);

    switch (param_id) {
    case PROP_CSS:
        g_value_set_string (value, cr->css);
        break;

    case PROP_TEXT:
        if (cr->real)
            g_object_get_property (G_OBJECT (cr->real), "text", value);
        break;

    default:
        G_OBJECT_WARN_INVALID_PROPERTY_ID (object, param_id, psec);
        break;
    }
}

static void
css_class_cell_renderer_set_property (GObject      *object,
                                      guint         param_id,
                                      const GValue *value,
                                      GParamSpec   *pspec) {
    CSSClassCellRenderer *cr;

    cr = CSS_CLASS_CELL_RENDERER (object);

    switch (param_id) {
    case PROP_CSS:
        g_free (cr->css);
        cr->css = g_value_dup_string (value);
        break;

    case PROP_TEXT:
        if (cr->real)
            g_object_set_property (G_OBJECT (cr->real), "text", value);
        break;

    default:
        G_OBJECT_WARN_INVALID_PROPERTY_ID(object, param_id, pspec);
        break;
    }
}

GtkCellRenderer *
css_class_cell_renderer_new (void)
{
  return g_object_new(TYPE_CSS_CLASS_CELL_RENDERER, NULL);
}

void
css_class_cell_renderer_set_real (CSSClassCellRenderer *cr, GtkCellRenderer *real) {
  if (cr->real)
    g_object_unref (cr->real);

  cr->real = real;

  if (cr->real)
    g_object_ref_sink (cr->real);
}

static void css_class_cell_renderer_apply_style_from_css (GtkCellRenderer *cell,
                                                          GtkWidget *widget) {
    GtkStyleContext *context;
    CSSClassCellRenderer *cellr;
    GdkRGBA background;
    gchar *printed_style_context;
    GValue value = G_VALUE_INIT;
    PangoFontDescription *font_desc;

    cellr = CSS_CLASS_CELL_RENDERER(cell);

    context = gtk_widget_get_style_context (widget);

    if (!cellr->real)
        return;

G_GNUC_BEGIN_IGNORE_DEPRECATIONS
    gtk_style_context_get_background_color (context, gtk_style_context_get_state (context), &background);
G_GNUC_END_IGNORE_DEPRECATIONS

    g_value_init (&value, GDK_TYPE_RGBA);
    g_value_set_boxed (&value, &background);
    g_object_set_property (G_OBJECT (cellr->real), "background-rgba", &value);
    g_value_unset (&value);

    gtk_style_context_get (context, gtk_style_context_get_state (context), "font", &font_desc, NULL);
    g_value_init (&value, PANGO_TYPE_FONT_DESCRIPTION);
    g_value_set_boxed (&value, font_desc);
    g_object_set_property (G_OBJECT (cellr->real), "font-desc", &value);
    g_value_unset (&value);

    pango_font_description_free (font_desc);

    printed_style_context = gtk_style_context_to_string (context, GTK_STYLE_CONTEXT_PRINT_SHOW_STYLE | GTK_STYLE_CONTEXT_PRINT_RECURSE);

    g_value_init (&value, G_TYPE_BOOLEAN);
    if (strstr(printed_style_context, "text-decoration-line: underline") != NULL) {
        g_value_set_boolean (&value, TRUE);
    } else {
        g_value_set_boolean (&value, FALSE);
    }
    g_object_set_property (G_OBJECT (cellr->real), "underline", &value);
    g_value_unset (&value);

    g_free(printed_style_context);
}

static void css_class_cell_renderer_set_class (GtkCellRenderer *cell,
                                               GtkWidget *widget) {
    GtkStyleContext *context;
    CSSClassCellRenderer *cellr;

    cellr = CSS_CLASS_CELL_RENDERER(cell);

    context = gtk_widget_get_style_context (widget);

    if (cellr->css)
        gtk_style_context_add_class (context, cellr->css);
}

static void
css_class_cell_renderer_get_size (GtkCellRenderer *cell,
                                  GtkWidget       *widget,
                                  const GdkRectangle    *cell_area,
                                  gint            *x_offset,
                                  gint            *y_offset,
                                  gint            *width,
                                  gint            *height)
{
    CSSClassCellRenderer *cellr;

    css_class_cell_renderer_set_class(cell, widget);
    css_class_cell_renderer_apply_style_from_css (cell, widget);

    cellr = CSS_CLASS_CELL_RENDERER(cell);

    if (cellr->real && GTK_CELL_RENDERER_GET_CLASS (cellr->real)->get_size) {
        GTK_CELL_RENDERER_GET_CLASS (cellr->real)->get_size (GTK_CELL_RENDERER (cellr->real), widget, cell_area, x_offset, y_offset, width, height);
    }
}


static void
css_class_cell_renderer_render (GtkCellRenderer *cell,
                                cairo_t         *cr,
                                GtkWidget       *widget,
                                const GdkRectangle   *background_area,
                                const GdkRectangle   *cell_area,
                                GtkCellRendererState  flags)
{
    CSSClassCellRenderer *cellr;

    css_class_cell_renderer_set_class (cell, widget);
    css_class_cell_renderer_apply_style_from_css (cell, widget);

    cellr = CSS_CLASS_CELL_RENDERER(cell);

    if (cellr->real && GTK_CELL_RENDERER_GET_CLASS (cellr->real)->render) {
        GTK_CELL_RENDERER_GET_CLASS (cellr->real)->render (GTK_CELL_RENDERER (cellr->real), cr, widget, background_area, cell_area, flags);
    }
}



static GtkCellEditable *
css_class_cell_renderer_start_editing (GtkCellRenderer      *cell,
                                       GdkEvent             *event,
                                       GtkWidget            *widget,
                                       const gchar          *path,
                                       const GdkRectangle   *background_area,
                                       const GdkRectangle   *cell_area,
                                       GtkCellRendererState  flags) {
    CSSClassCellRenderer *cellr;

    cellr = CSS_CLASS_CELL_RENDERER(cell);

    if (cellr->real && GTK_CELL_RENDERER_GET_CLASS (cellr->real)->start_editing) {
        return GTK_CELL_RENDERER_GET_CLASS (cellr->real)->start_editing (GTK_CELL_RENDERER (cellr->real), event, widget, path, background_area, cell_area, flags);
    }

    return NULL;
}

static void
css_class_cell_renderer_get_preferred_width            (GtkCellRenderer       *cell,
                                                        GtkWidget             *widget,
                                                        gint                  *minimal_size,
                                                        gint                  *natural_size) {
    CSSClassCellRenderer *cellr;

    css_class_cell_renderer_set_class(cell, widget);
    css_class_cell_renderer_apply_style_from_css (cell, widget);

    cellr = CSS_CLASS_CELL_RENDERER(cell);

    if (cellr->real && GTK_CELL_RENDERER_GET_CLASS (cellr->real)->get_preferred_width) {
        GTK_CELL_RENDERER_GET_CLASS (cellr->real)->get_preferred_width (GTK_CELL_RENDERER (cellr->real), widget, minimal_size, natural_size);
    }
}

static void
css_class_cell_renderer_get_preferred_height           (GtkCellRenderer       *cell,
                                                        GtkWidget             *widget,
                                                        gint                  *minimal_size,
                                                        gint                  *natural_size) {
    CSSClassCellRenderer *cellr;

    css_class_cell_renderer_set_class(cell, widget);
    css_class_cell_renderer_apply_style_from_css (cell, widget);

    cellr = CSS_CLASS_CELL_RENDERER(cell);

    if (cellr->real && GTK_CELL_RENDERER_GET_CLASS (cellr->real)->get_preferred_height) {
        GTK_CELL_RENDERER_GET_CLASS (cellr->real)->get_preferred_height (GTK_CELL_RENDERER (cellr->real), widget, minimal_size, natural_size);
    }
}

static void
css_class_cell_renderer_get_preferred_height_for_width (GtkCellRenderer       *cell,
                                                        GtkWidget             *widget,
                                                        gint                   width,
                                                        gint                  *minimum_height,
                                                        gint                  *natural_height) {
    CSSClassCellRenderer *cellr;

    css_class_cell_renderer_set_class(cell, widget);
    css_class_cell_renderer_apply_style_from_css (cell, widget);

    cellr = CSS_CLASS_CELL_RENDERER(cell);

    if (cellr->real && GTK_CELL_RENDERER_GET_CLASS (cellr->real)->get_preferred_height_for_width) {
        GTK_CELL_RENDERER_GET_CLASS (cellr->real)->get_preferred_height_for_width (GTK_CELL_RENDERER (cellr->real), widget, width, minimum_height, natural_height);
    }
}

static void
css_class_cell_renderer_get_aligned_area               (GtkCellRenderer       *cell,
                                                        GtkWidget             *widget,
                                                        GtkCellRendererState   flags,
                                                        const GdkRectangle    *cell_area,
                                                        GdkRectangle          *aligned_area) {
    CSSClassCellRenderer *cellr;

    css_class_cell_renderer_set_class(cell, widget);
    css_class_cell_renderer_apply_style_from_css (cell, widget);

    cellr = CSS_CLASS_CELL_RENDERER(cell);

    if (cellr->real && GTK_CELL_RENDERER_GET_CLASS (cellr->real)->get_aligned_area) {
        GTK_CELL_RENDERER_GET_CLASS (cellr->real)->get_aligned_area (GTK_CELL_RENDERER (cellr->real), widget, flags, cell_area, aligned_area);
    }
}

static CSSClassCellRenderer*
toCSSClassCellRenderer(void *p)
{
	return (CSS_CLASS_CELL_RENDERER(p));
}

static GtkCellRenderer*
toGtkCellRenderer(void *p) {
  return (GTK_CELL_RENDERER(p));
}

#endif
