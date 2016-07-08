PREFIX       ?= /usr
INSTALL_DIR  ?= install -d -m 755
INSTALL_PROG ?= install -m 755
INSTALL_FILE ?= install -m 644
RM           ?= rm -f

all:
	@echo Run \'make install\' to install punf.

install:
	@echo "Installing binaries."
	$(INSTALL_DIR) $(DESTDIR)$(PREFIX)/bin
	$(INSTALL_PROG) punf $(DESTDIR)$(PREFIX)/bin/punf
	@echo "Installing configs."
	$(INSTALL_DIR) $(DESTDIR)$(PREFIX)/share/punf
	$(INSTALL_FILE) configs/config $(DESTDIR)$(PREFIX)/share/punf/config
	@echo "Installing completions."
	$(INSTALL_DIR) $(DESTDIR)$(PREFIX)/share/fish/completions
	$(INSTALL_FILE) completions/punf.fish $(DESTDIR)$(PREFIX)/share/fish/completions/punf.fish

uninstall:
	@echo "Uninstalling binaries."
	$(RM) $(DESTDIR)$(PREFIX)/bin/punf
	@echo "Uninstalling configs."
	$(RM) -r $(DESTDIR)$(PREFIX)/share/punf
	@echo "Uninstalling completions."
	$(RM) $(DESTDIR)$(PREFIX)/share/fish/completions/punf.fish
