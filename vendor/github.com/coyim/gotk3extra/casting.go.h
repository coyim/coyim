
#pragma once

static const gchar *
object_get_class_name(GObject *object)
{
	return G_OBJECT_CLASS_NAME(G_OBJECT_GET_CLASS(object));
}

static GObject *
toGObject(void *p)
{
	return (G_OBJECT(p));
}
