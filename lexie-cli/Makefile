
all:
	go build

# ./lexie-cli -l ../in/django3.lex -D match,db_DumpDFAPool,db_DumpPool,db_Matcher_02,db_NFA_LnNo,DumpParseNodes2,db_FlushTokenBeforeAfter ./bob.in

test1:
	go build
	./lexie-cli -l ../in/django3.lex -i ./bob.in

test2:
	go build
	./lexie-cli -l ../in/django3.lex ./bob.in

test3db:
	go build
	./lexie-cli -l ../in/django3.lex -X match,db_DumpDFAPool,db_DumpPool,db_Matcher_02,db_NFA_LnNo,DumpParseNodes2,db_FlushTokenBeforeAfter -i ./autoescape.tpl

test3:
	go build
	./lexie-cli -l ../in/django3.lex -i ./autoescape.tpl

test4:
	go build
	./lexie-cli -l ../in/django3.lex -i ./filters.tpl

test5:
	go build
	./lexie-cli -l ../in/django3.lex -X match,db_DumpDFAPool,db_DumpPool,db_Matcher_02,db_NFA_LnNo,DumpParseNodes2,db_FlushTokenBeforeAfter,in-echo-machine -i ./dig.tpl

test6:
	go build
	./lexie-cli -l ../in/django3.lex -i ./comment.06.tpl


