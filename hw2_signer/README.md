# coursera.org домашние задания в go
hw2_signer.v.1
попытался воссоздать вывод функций SingleHash, MultiHash, CombineResults

ни о каких каналах речи и не идет.

hw2_signer.v.2
внутри функций SingleHash, MultiHash использую каналы для ускорения вычислений
все еще не оптимально, но результат времени вычислений на лицо.
нет функции ExecutePipeline которая обеспечивает нам конвейерную обработку

hw2_signer.v.3
есть ExecutePipeline которая примерно показывает, как работает конвейерная обработка
введен левый тип tjob только для примера работы ExecutePipeline а также переделаны маленько функции

signer.go
рабочий вариант НО специально отредактирована одна функция чтобы тест не проходил.
а запуск go run signer.go common.go
отрабатывает нормально без ошибок и укладывается по времени


