BINDIR=../bin

GIT_VERSION=$(shell git rev-parse HEAD)
TAG_VERSION=$(shell git tag -l --contains $$GIT_VERSION | tail -1)
KEYID=$(shell gpg2 --keyid-format 0xlong -K | grep '^sec' | head -1 | cut -d\  -f4 | cut -d\/ -f2)

sign-reproducible:
	./sign_build_info_with_key.sh $(KEYID) $(BINDIR)

send-reproducible-signature:
	./mail_build_info.sh $(BINDIR)/build_info.$(KEYID).rasc $(TAG_VERSION)

check-reproducible-signatures:
	./check_build_info_signatures_gh.rb $(TAG_VERSION)
