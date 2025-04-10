
all:
	go build

# ./lexie-cli -l ../in/django3.lex -D match,db_DumpDFAPool,db_DumpPool,db_Matcher_02,db_NFA_LnNo,DumpParseNodes2,db_FlushTokenBeforeAfter ./bob.in

setup:
	go build

tests: setup test1 test2 test3db test3 test4 test5 test6

test1:
	./lexie-cli -l ../in/django3.lex -i ./bob.in > out/test1.out
	diff ./out/test1.out ./ref/test1.out

test2:
	./lexie-cli -l ../in/django3.lex ./bob.in

test3db:
	./lexie-cli -l ../in/django3.lex -X match,db_DumpDFAPool,db_DumpPool,db_Matcher_02,db_NFA_LnNo,DumpParseNodes2,db_FlushTokenBeforeAfter -i ./autoescape.tpl

test3:
	./lexie-cli -l ../in/django3.lex -i ./autoescape.tpl

test4:
	./lexie-cli -l ../in/django3.lex -i ./filters.tpl

test5:
	./lexie-cli -l ../in/django3.lex -X match,db_DumpDFAPool,db_DumpPool,db_Matcher_02,db_NFA_LnNo,DumpParseNodes2,db_FlushTokenBeforeAfter,in-echo-machine -i ./dig.tpl

test6:
	./lexie-cli -l ../in/django3.lex -i ./comment.06.tpl


