DEFINITIONS=$(shell ls ../gui/definitions/*.xml)
IN_FILES=$(shell find ../gui -name "*.go" ! -path "gui/definitions/*")

# https://www.gnu.org/software/gettext/manual/gettext.html#Overview
# translatables from source -> .po
# .po -> translate to each locale -> .po (translated)
# .po (translated) -> .mo (binary)

# For now, a single translation domain. We could use multiple domains if
# needed.
DOMAIN=coy

TRANSLATIONS=\
	ar/LC_MESSAGES/$(DOMAIN).mo \
	en_US/LC_MESSAGES/$(DOMAIN).mo \
	pt_BR/LC_MESSAGES/$(DOMAIN).mo \
	sv_SE/LC_MESSAGES/$(DOMAIN).mo \
	zh_CN/LC_MESSAGES/$(DOMAIN).mo \
	es_EC/LC_MESSAGES/$(DOMAIN).mo

gettext: $(TRANSLATIONS)

$(DOMAIN).pot: $(DEFINITIONS) $(IN_FILES)
	xgettext --package-name=$(DOMAIN) --default-domain=$(DOMAIN) --from-code=utf-8 -L glade $(DEFINITIONS) -o $(DOMAIN).pot
	xgettext --package-name=$(DOMAIN) --default-domain=$(DOMAIN) --from-code=utf-8 -L python $(IN_FILES) --keyword=Local -j -o $(DOMAIN).pot

%.po: $(DOMAIN).pot
	./update-locales $(DOMAIN).pot $@

%.mo: %.po
	msgfmt -c -v -o $@ $<

update-locales: $(DOMAIN).pot
	./update-locales $(DOMAIN).pot $(TRANSLATIONS)

clean:
	$(RM) $(TRANSLATIONS)

# We need this to prevent make from removing the intermediate .po
.PRECIOUS: %.po
