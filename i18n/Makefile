DEFINITIONS=$(shell ls ../gui/definitions/*.xml)

# https://www.gnu.org/software/gettext/manual/gettext.html#Overview
# translatables from source -> .po
# .po -> translate to each locale -> .po (translated)
# .po (translated) -> .mo (binary)

# For now, a single translation domain. We could use multiple domains if
# needed.
DOMAIN=coy
TEMPLATE=$(DOMAIN).pot
TRANSLATIONS=\
	en_US.utf8/LC_MESSAGES/$(DOMAIN).mo \
	pt_BR.utf8/LC_MESSAGES/$(DOMAIN).mo

gettext: $(TEMPLATE) $(TRANSLATIONS)

$(TEMPLATE): $(DEFINITIONS)
	xgettext --package-name=$(DOMAIN) --default-domain=$(DOMAIN) -L glade $(DEFINITIONS) -o $(TEMPLATE)

%.po: $(TEMPLATE)
	./update-locales $(TEMPLATE) $@

%.mo: %.po
	msgfmt -c -v -o $@ $<

update-locales: $(TEMPLATE)
	./update-locales $(TEMPLATE) $(TRANSLATIONS)

clean:
	$(RM) $(TRANSLATIONS)

# We need this to prevent make from removing the intermediate .po
.PRECIOUS: %.po

