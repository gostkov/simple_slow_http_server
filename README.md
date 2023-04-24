## simple_slow_http_server

### Description
It is simple HTTP server, which can help to you test a some functions, from your projects.

You can use 3 methods (slow, fast, error):
- http://127.0.0.1:8080/slow/timeout=10 this method emulates the delay of server, add parameter `timeout`
- http://127.0.0.1:8080/fast/ this method gives a quick answer
- http://127.0.0.1:8080/error/code=500 this method returns the custom error, you need pass parameter `code`

### Usage

You can pass environment variables `BASIC_AUTH_LOGIN` and `BASIC_AUTH_PASSWORD` for enable basic_auth. By default, basic auth is disabled.

You can build local
```shell
git clone https://github.com/gostkov/simple_slow_http_server
cd simple_slow_http_server
go run simple_slow_http_server
```

or you can use docker image from docker hub

```shell
docker run -p 8080:8080 gostkov/simple_slow_http_server
```

### Описание
Это простой HTTP сервер, с помощью которого, вы можете протестить какие-то функции из вашего проекта.

Вы можете использовать 3 метода:
- http://127.0.0.1:8080/slow/timeout=10 данный метод эмулирует задержку, параметр `timeout`
- http://127.0.0.1:8080/fast/ данный метод дает быстрый ответ
- http://127.0.0.1:8080/error/code=500 данный метод возвращает нужную вам ошибку, через параметр `code`

А также включить basic авторизацию

### Использование
Если вы хотите использовать basic авторизацию, то необходимо задать две переменные: `BASIC_AUTH_LOGIN` and `BASIC_AUTH_PASSWORD`. По-умолчанию выключена.

Запустить можно скачав с github'а
```shell
git clone https://github.com/gostkov/simple_slow_http_server
cd simple_slow_http_server
go run simple_slow_http_server
```

Или используя образ из docker hub

```shell
docker run -p 8080:8080 gostkov/simple_slow_http_server
```
