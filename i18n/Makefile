DEFINITIONS=$(shell ls ../gui/definitions/*.xml)
IN_FILES=$(shell find ../gui -name "*.go" ! -path "gui/definitions/*")

# https://www.gnu.org/software/gettext/manual/gettext.html#Overview
# translatables from source -> .po
# .po -> translate to each locale -> .po (translated)
# .po (translated) -> .mo (binary)

# For now, a single translation domain. We could use multiple domains if
# needed.
DOMAIN=coy

TEMPLATE_UI=$(DOMAIN)-ui.pot
TEMPLATE_CODE=$(DOMAIN)-code.pot

TRANSLATIONS=\
	en_US/LC_MESSAGES/$(DOMAIN).mo \
	pt_BR/LC_MESSAGES/$(DOMAIN).mo \
	sv_SE/LC_MESSAGES/$(DOMAIN).mo

gettext: $(TEMPLATE_UI) $(TEMPLATE_CODE) $(TRANSLATIONS)

$(TEMPLATE_UI): $(DEFINITIONS)
	xgettext --package-name=$(DOMAIN) --default-domain=$(DOMAIN) --from-code=utf-8 -L glade $(DEFINITIONS) -o $(TEMPLATE_UI)

$(TEMPLATE_CODE): $(IN_FILES)
	xgettext --package-name=$(DOMAIN) --default-domain=$(DOMAIN) --from-code=utf-8 -L python $(IN_FILES) --keyword=Local -o $(TEMPLATE_CODE)

%.po: $(TEMPLATE_UI) $(TEMPLATE_CODE)
	./update-locales $(TEMPLATE_UI) $@
	./update-locales $(TEMPLATE_CODE) $@

%.mo: %.po
	msgfmt -c -v -o $@ $<

update-locales: $(TEMPLATE_UI) $(TEMPLATE_CODE)
	./update-locales $(TEMPLATE_UI) $(TRANSLATIONS)
	./update-locales $(TEMPLATE_CODE) $(TRANSLATIONS)

clean:
	$(RM) $(TRANSLATIONS)

# We need this to prevent make from removing the intermediate .po
.PRECIOUS: %.po
