## blockchain-chat

Приложение позволяет подключаться к локальной p2p-сети и проводить обмен зашифрованными анонимными сообщениями. Под анонимностью понимается:
- отсутствие информации об IP-адресах отправителя и получателя в сообщении
- отсутствие публичных ключей отправителя и получателя в сообщении
При этом корректность диалога (отсутствие подмены сообщений) всегда может быть подтверждено благодаря использованию технологии, аналогичной blockchain.

## Установка и запуск
Установка:
```
go get github.com/poslegm/blockchain-chat
make install
```
Запуск:
```
make run
```

После этого при переходе на localhost:8080 в браузере должен открываться графический пользовательский интерфейс.
Для начала работы необходимо добавить в файл ips.txt в корневой директории приложения известные IP-адреса участников сети. Если этот файл пустой, то приложение будет работать в режиме ожидания запроса на подключение к сети.
Кроме того, для отправки сообщений необходимо добавить в локальную базу данных пару GPG ключей. Сделать это можно через графический интерфейс, перейдя по ссылке «Добавить пару ключей»
