generate:
	ruby ./generate.rb

touch:
	ls *.xml | xargs -n1 touch

doctor: touch generate
	git diff --exit-code .