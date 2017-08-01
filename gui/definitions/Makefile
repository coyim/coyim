generate:
	ruby ./generate.rb
	gofmt -w .

touch:
	ls *.xml | xargs -n1 touch

doctor: touch generate
	git diff --exit-code .