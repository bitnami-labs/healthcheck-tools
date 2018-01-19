all: ssl-checker

lint-all: ssl-checker-lint

test-all: ssl-checker-test

ssl-checker:
	$(MAKE) -C ssl-checker/

ssl-checker-lint:
	$(MAKE) -C ssl-checker/ lint


ssl-checker-test:
	$(MAKE) -C ssl-checker/ test
