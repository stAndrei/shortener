Укорачиватель ссылок


Необходимо реализовать сервис, который должен предоставлять API по созданию сокращённых ссылок следующего формата:
Ссылка должна быть уникальной и на один оригинальный URL должна ссылаться только одна сокращенная ссылка
Ссылка должна быть длинной 10 символов
Ссылка должна состоять из символов латинского алфавита в нижнем и верхнем регистре, цифр и символа _ (подчеркивание)

Сервис должен содержать следующие endpoint-ы:
Запрос GET /{shortUrl}, который будет перенаправлять на оригинальный URL
Запрос POST /, который будет принимать оригинальный URL и возвращать сокращённый

Решение должно быть предоставлено в «конечном виде», а именно:
Сервис должен быть распространён в виде Docker-образа
В качестве хранилища можно использовать in-memory решение
Для API должно быть описание в формате OpenAPI в виде endpoint-a (GET /docs/openapi)

Покрыть реализованный функционал Unit-тестами


Получилось:
Вызыв сервера
$ curl -X POST -i -F 'url=https://yandex.ru?q=126' 0.0.0.0:8000
$ curl -i http://0.0.0.0:8000/UrzIFqGc2h

Прогон тестов:
$ make test
