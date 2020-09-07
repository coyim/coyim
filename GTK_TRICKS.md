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
