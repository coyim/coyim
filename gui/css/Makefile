generate:
	go generate -x .

touch:
	ls definitions/*.css | xargs -n1 touch

doctor: touch generate
	git diff --exit-code .
