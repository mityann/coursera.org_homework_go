# coursera.org_homework_go
hw1_tree_v.0.5

Утилита tree.
Запуск go run main.go . -f

Выводит дерево каталогов и файлов (если указана опция -f).

Запускать тесты через `go test -v` находясь в папке c заданием.

TODO:
* Если вы пользуетесь windows - помните, что там и в linux разделители директорий различаются - используйте лучше `string(os.PathSeparator)`
* Результаты ( список папок-файлов ) должны быть отсортированы по алфавиту. Т.е. у вас должен быть код который отсортирует уровень. Смотрите для этого пакет sort. Это самая частая причина непрохождения тестов. Тесты запускаются в среде linux. В задании есть докер-файл для тестов ровно в тех же условиях, он сразу выявит все проблемы.
