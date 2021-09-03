# GTK Tips and Tricks

The goal of this file is to provide a few guidelines on how to do things that might not be completely obvious. Please
add things here that you think could be helpful for other people working in the UI. This document assumes that you are
using `gotk3adapter`. Basically, if you spend time figuring out how to do something weird, this is the place to document
the solution - especially if you think it's something other people will want to do.

You can also add entries here and ask for a solution, if you don't have one yourself. If someone comes up with a good
idea for how to do it, they can add it here!


## Store arbitrary data in a GTK ListStore or TreeStore

Sometimes you want to store arbitrary data along side an entry in a ListStore or TreeStore. This is not completely
obvious. The actual C structures does allow you to store a raw pointer, but this doesn't seem to work well with our
wrappers. There is also the problem of garbage collection wit hthis approach. A simple method that requires just a
little bit extra work is to maintain an independent map from integer to the value you want to store, and then store the
integer in a column in the map. You should generally randomly select the integer used to store, and check that it isn't
already used as a key. The only tricky part about this pattern is that you need to make sure that the integer is not too
large, since there seems to be weird interactions and overflow if you store really large numbers. For this reason, I
recommend using a 31 bit number, to be on the safe side. However, you can't actually use the `int32` data type. These
snippes from `gui/muc_public_rooms.go` show how I put the data in, and how I get it back out:

```golang
func (prv *mucPublicRoomsView) getFreshRoomInfoIdentifierAndSet(rl *muc.RoomListing) int {
	for {
		v := int(rand.Int31())
		_, ok := prv.roomInfos[v]
		if !ok {
			prv.roomInfos[v] = rl
			return v
		}
	}
}
```

```golang
	roomInfoValue, e1 := prv.roomsModel.GetValue(iter, mucListRoomsIndexRoomInfo)
	roomInfoRef, e2 := roomInfoValue.GoValue()
	if e1 != nil || e2 != nil {
		return nil, nil, errNoRoomSelected
	}

	roomInfoRealVal := roomInfoRef.(int)
	rl, ok := prv.roomInfos[roomInfoRealVal]
	if !ok || rl == nil {
		return nil, nil, errNoPossibleSelection
	}
```

It does require a little bit of book keeping, but nothing too complicated.


## Keep something centered, even when you're adding or removing things above, below or to the sides

Very often you want to have some kind of main content, and this content should be centered in the view. This is fairly
simple to achieve, but say that you want to add a side bar - now the content will end up being centered in the remaining
space. Another example is if you have a form and you want to display one or several validation warnings. But when you
add the warnings, the form will be pushed down more and more. From a usability perspective, both of these are unwanted
behavior.

The simple solution to this is to use the Center Widget functionality that exists on boxes. There are two steps to do
it - first, you have to separate out the object in the XML definition. The object can _not_ be part of the Box
already. Just put it outside. Then, once you are loading and creating the object, after everything else is set up, but
before you show the object, you call `SetCenterWidget` with the object in question. By doing this, all other things you
already have added or will add to any side of the center widget will not impact the position of the center
widget. Remember, a Box is only working in one dimension, either horizontally or vertically. So the center widget will
be dependent on what type of box it is, whether it will be centered vertically or horizontally. You can take a look at
`gui/muc_room_lobby.go` and how we add the box with ID `mainContent` as a center widget, even though we might add one or
several warnings at the top of the box.


## Don't allow a text label to expand the box and window it's inside - imperfect solution

If you add a text label with a lot of text, you will generally end up in a situation where the text expands vertically,
and this will expand the box its inside, all the way out to the window. I have not found a good way to stop this
behavior and let it expand vertically instead of horizontally. However, you can use the `max-width-chars` property to
make the label wrap at a specific point. This makes it possible to stop the expansion of the box, to some degree. If you
do this, and the left margin becomes unpredictable, sometimes more indented and sometimes less, you can solve this by
setting the `halign` and `xalign` properties on the label. They should be set to `start` and `0`, respectively.


## Don't allow a text label to expand the box and window it's inside - BETTER SOLUTION WANTED

There are several problems with the above solution. Measuring in terms of characters is not great, since it will change
with the font and font size. Also, weird things happen with the left side of the box, making the indentation different,
depending on the text. This makes my eyes hurt. This solution also makes it weird when using different languages and
scripts.

I would like to have a better solution to this, where I can basically say that a box or label should NOT expand to the
sides, only vertically.

PLEASE FILL IN YOUR SOLUTION HERE.


## Get faster feedback about GTK warnings

When you get warnings in the style of 

```
(CoyIM:19511): Gtk-CRITICAL **: 21:40:31.567: gtk_widget_set_visible: assertion 'GTK_IS_WIDGET (widget)' failed

(CoyIM:19511): Gtk-CRITICAL **: 21:40:31.567: gtk_widget_set_visible: assertion 'GTK_IS_WIDGET (widget)' failed
```

It can sometimes be very hard to find where and why they are happening. By setting the environment variable `G_DEBUG` to
`fatal_warnings`, the first such warning will crash the Go program and give you a helpful stack trace. You can invoke
Coy in this way to make that happen easily:

```
G_DEBUG=fatal_warnings ./bin/coyim
```


## Get information about garbage collection problems

Sometimes garbage collection will happen and cause problems, and sometimes Golang objects are freed too early, causing
reference counting problems. One way of looking at these kinds of problems is to modify the Finalizer for GObject, to
print information during finalization. In the `Take` method in `vendor/github.com/gotk3/gotk3/glib/glib.go`, I sometimes
add this kind of code:

```
		refc := v.native().ref_count
		if refc < 2 {
			nm := v.TypeFromInstance().Name()
			extra := ""
			if nm == "GtkLabel" {
				extra = HelperBla(v)
			}
			fmt.Printf("1. Finalizing object of type: %v (ref count: %d) extra: %s\n", nm, refc, extra)
		}
``

This will print some useful information if the ref count is 1 or 0, indicating that the object will be freed.


## Do Garbage Collection more often to provoke issues wit reference counting

Sometimes you need to push the Garbage Collection a bit harder to see crashes in a consistent way. You can do this by
adding this snippet of code to `main.go`:

```
	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			runtime.GC()
		}
	}()

```


## To debug GTK objects, use this kind of print code

Sometimes you need a deeper understanding of GTK objects. These snippets of code can be added to Gotk3 temporarily to
print a bunch of information. Add more evolution to this code when you use it!

```
func genIn(indent int) string {
	res := ""
	for i := 0; i < indent; i++ {
		res = res + " "
	}
	return res
}

func printCollection(indent int, v Container) {
	v.GetChildren().Foreach(func(item interface{}) {
		printAllNative(indent, item)
	})
}

func (w *Widget) Print() {
	printAllNative(0, w)
}

func printAllNative(indent int, v interface{}) {
	switch vv := v.(type) {
	case *Widget:
		va, er := vv.TestGoValue()
		if er != nil {
			fmt.Printf(" GOT ERROR CONVERTING: %#v\n", va)
		} else {
			printAllNative(indent, va)
		}
	case *ScrolledWindow:
		fmt.Printf("%sScrolledWindow(%#v) {\n", genIn(indent), vv.native())
		printCollection(indent+1, vv.Bin.Container)
		fmt.Printf("%s}\n", genIn(indent))
	case *Viewport:
		fmt.Printf("%sViewport(%#v) {\n", genIn(indent), vv.native())
		printCollection(indent+1, vv.Bin.Container)
		fmt.Printf("%s}\n", genIn(indent))
	case *TextBuffer:
		fmt.Printf("%sTextBuffer(%#v) {\n", genIn(indent), vv.native())
	case *ListBox:
		fmt.Printf("%sListBox(%#v) {\n", genIn(indent), vv.native())
		printCollection(indent+1, vv.Container)
		fmt.Printf("%s}\n", genIn(indent))
	case *ListBoxRow:
		fmt.Printf("%sListBoxRow(%#v) {\n", genIn(indent), vv.native())
		printCollection(indent+1, vv.Bin.Container)
		fmt.Printf("%s}\n", genIn(indent))
	case *Box:
		fmt.Printf("%sBox(%#v) {\n", genIn(indent), vv.native())
		printCollection(indent+1, vv.Container)
		fmt.Printf("%s}\n", genIn(indent))
	case *Menu:
		fmt.Printf("%sMenu(%#v)\n", genIn(indent), vv.native())
	case *MenuItem:
		fmt.Printf("%sMenuItem(%#v)\n", genIn(indent), vv.native())
	case *SeparatorMenuItem:
		fmt.Printf("%sSeparatorMenuItem(%#v)\n", genIn(indent), vv.native())
	case *CheckMenuItem:
		fmt.Printf("%sCheckMenuItem(%#v)\n", genIn(indent), vv.native())
	case *HeaderBar:
		fmt.Printf("%sHeaderBar(%#v) {\n", genIn(indent), vv.native())
		printCollection(indent+1, vv.Container)
		fmt.Printf("%s}\n", genIn(indent))
	case *Label:
		ttt, _ := vv.GetText()
		fmt.Printf("%sLabel(%#v) Text=%s\n", genIn(indent), vv.native(), ttt)
	case *Notebook:
		fmt.Printf("%sNotebook(%#v) {\n", genIn(indent), vv.native())
		printCollection(indent+1, vv.Container)
		fmt.Printf("%s}\n", genIn(indent))
	case *Overlay:
		fmt.Printf("%sOverlay(%#v) {\n", genIn(indent), vv.native())
		printCollection(indent+1, vv.Bin.Container)
		fmt.Printf("%s}\n", genIn(indent))
	case *Separator:
		fmt.Printf("%sSeparator(%#v)\n", genIn(indent), vv.native())
	case *SizeGroup:
		fmt.Printf("%sSizeGroup(%#v)\n", genIn(indent), vv.native())
	case *Spinner:
		fmt.Printf("%sSpinner(%#v)\n", genIn(indent), vv.native())
	case *CheckButton:
		fmt.Printf("%sCheckButton(%#v)\n", genIn(indent), vv.native())
	case *Grid:
		fmt.Printf("%sGrid(%#v) {\n", genIn(indent), vv.native())
		printCollection(indent+1, vv.Container)
		fmt.Printf("%s}\n", genIn(indent))
	default:
		fmt.Printf("NO TYPE MATCH: %#v\n", vv)
	}
}
```

For this to work, you also need to add this snippet to the `glib` part of `gotk3`:

```
func (v *Object) TestGoValue() (interface{}, error) {
       return v.goValue()
}
```
