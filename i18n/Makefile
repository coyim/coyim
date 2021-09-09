DEFINITIONS=$(shell ls ../gui/definitions/*.xml)
IN_FILES=$(shell find ../gui -name "*.go" ! -path "../gui/definitions/*" ! -path "../gui/definitions.go" ! -path "../gui/*_test.go")

# https://www.gnu.org/software/gettext/manual/gettext.html#Overview
# translatables from source -> .po
# .po -> translate to each locale -> .po (translated)
# .po (translated) -> .mo (binary)

# For now, a single translation domain. We could use multiple domains if
# needed.
DOMAIN=coy

TRANSLATIONS=\
	ar/LC_MESSAGES/$(DOMAIN).mo \
	pt/LC_MESSAGES/$(DOMAIN).mo \
	sv/LC_MESSAGES/$(DOMAIN).mo \
	zh/LC_MESSAGES/$(DOMAIN).mo \
	nb_NO/LC_MESSAGES/$(DOMAIN).mo \
	fr/LC_MESSAGES/$(DOMAIN).mo \
	es/LC_MESSAGES/$(DOMAIN).mo

gettext: $(TRANSLATIONS)

$(DOMAIN).pot: $(DEFINITIONS) $(IN_FILES)
	xgettext --package-name=$(DOMAIN) --default-domain=$(DOMAIN) --from-code=utf-8 -L glade $(DEFINITIONS) -o $(DOMAIN).pot
	xgettext --package-name=$(DOMAIN) --default-domain=$(DOMAIN) --from-code=utf-8 -L python $(IN_FILES) --keyword=Local --keyword=Localf -j -o $(DOMAIN).pot

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

generate:
	go generate -x .

touch:
	ls $(TRANSLATIONS:.mo=.po) | xargs -n1 touch

doctor: touch generate
	git diff --exit-code .
