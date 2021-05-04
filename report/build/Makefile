# Try to subdue interactive prompts. Fail instead.
TEX=latexmk -interaction=nonstopmode -halt-on-error

# You want latexmk to always run, because make does not have all the info.
# Also, include non-file targets in .PHONY so they are run regardless of any
# file of the given name existing.
# .PHONY: phony.pdf clean watch
.PHONY: clean watch

clean:
	rm -f *.bbl *.run.xml *.xdv *.fls *.aux *.log *.fdb_latexmk

# phony.pdf:
%.pdf: %.tex
	$(TEX) -xelatex -pdf $< && $(TEX) -c $<

# Continuously make whenever there are changes.
# Works well with a good editor like Kakoune;
#   :make watch-main
# will fork. Output is sent to log buffer.
# Delete the log buffer to stop watching.
watch-%: %.tex
	$(TEX) -xelatex -pdf -pvc $<