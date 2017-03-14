generate:
	go generate -x .

touch:
	ls definitions/*.xml | xargs -n1 touch

doctor: touch generate
	git diff --exit-code .