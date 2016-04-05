.PHONY: generate

gschemas.compiled: *.xml
	glib-compile-schemas .

schemas.go: gschemas.compiled
	ruby ./generate.rb

generate: schemas.go

touch:
	ls *.xml | xargs -n1 touch

doctor: touch generate
	git diff --exit-code .
