init:
	python -c"from minitwit import init_db; init_db()"

build:
	gcc tools/flag_tool.c -l sqlite3 -o bin/flag_tool

clean:
	rm flag_tool
