# LayoutConfig.api - REST API основной WEB API для взаимодействия со схемой размещения (layout-ом)

[![pipeline status](https://git.countmax.ru/countmax/layoutconfig.api/badges/master/pipeline.svg)](https://git.countmax.ru/countmax/layoutconfig.api/-/commits/master) [![coverage report](https://git.countmax.ru/countmax/layoutconfig.api/badges/master/coverage.svg)](https://git.countmax.ru/countmax/layoutconfig.api/-/commits/master)

## назначение компонента

Создавался для решения [general#73](https://git.countmax.ru/countmax/general/-/issues/73)  
WEB API доступ к схеме размещения и ее элементам для выполнения операций CRUD

### [Техническое решение](https://git.countmax.ru/countmax/layoutconfig.api/wikis/%D0%A2%D0%B5%D1%85%D0%BD%D0%B8%D1%87%D0%B5%D1%81%D0%BA%D0%BE%D0%B5%20%D1%80%D0%B5%D1%88%D0%B5%D0%BD%D0%B8%D0%B5)  

### [Контракт.v2](https://git.countmax.ru/countmax/layoutconfig.api/-/wikis/%D0%9A%D0%BE%D0%BD%D1%82%D1%80%D0%B0%D0%BA%D1%82.v2)

>Определяется в [swagger спеке](https://git.countmax.ru/countmax/layoutconfig.api/-/blob/master/docs/swagger.yaml)

## перед первым началом работы с исходным кодом компонента

* Установить go >=1.13
* Установить git
* Клонировать репозиторий `git clone git@git.countmax.ru:countmax/layoutconfig.api.git`
* Перейти в папку проекта и установить go зависимости `cd layoutconfig.api && go mod download`

## для построения компонента

### makefile

```bash
cd path/to/layoutconfig.api
make build
make run # start layoutconfig.api bind to 9001 and 9002 tcp ports
```

### В ручном режиме

```bash
cd path/to/layoutconfig.api
go build
./layoutconfig.api -p=":9001" -sp=":9002" # start layoutconfig.api bind to 8080 tcp port
```

### docker

```bash
cd path/to/layoutconfig.api
make docker
```

## CI

В качестве CI используется [gitlab-ci](https://docs.gitlab.com/ee/ci/)  
Детали отображены в файле .gitlab-ci.yml в корне проекта  
При чери-пике в ветку `production` / `pre-production` создается docker образ и загружается на приватный [docker-hub](https://hub.watcom.ru)

## шаги необходимые выполнить для получения результатов построения компонента

* Выполнить билд (см пред пункт)
* Сконфигурировать приложение
* Запустить его

### Запуск в docker-е

```bash
$ docker run --name layout-api -p 8080:8000 -p 8081:8001 \
  -e LAYOUT_COUNTMAX_TOKEN=$(Token-commonapi-value) \
  -d hub.watcom.ru/layoutconfig.api
```

Или настроить параметры в файле `docker-compose.yml`

```bash
$ docker-compose up -d
```

### Запуск windows service

> в системе должна быть установлена утилита [nssm](http://nssm.cc/usage)

```powershell
cd path/to/layoutconfig.api.exe
nssm install layoutconfig.api
```

### Запуск linux demon

[Linux service howto](https://jonathanmh.com/deploying-go-apps-systemd-10-minutes-without-docker/)

* Потребуется предварительная сборка приложения под linux (`make linuxbuild`)
* Загрузка исполняемого и конфигурационого файла не описывается (используйте `scp/winscp`)
* Добавить свойства исполняемого файла `chmod +x layoutconfig.api`  
* Создать service file /lib/systemd/system/layoutconfig.api.service:

```bash
sudo systemctl edit /lib/systemd/system/layoutconfig.api.service
```

* Открыть файл в редакторе (требуются права root) **nano/vim**

```bash
[Unit]
Description=layoutconfig.api
Wants=network-online.target
After=network-online.target

[Service]
User=userForExecute
Group=userForExecute
Type=simple
Restart=always
RestartSec=5s
ExecStart=/usr/local/bin/layoutconfig.api -p=":8000" -sp=":8001"

[Install]
WantedBy=multi-user.target
```

* Сохранить и закрыть файл layoutconfig.api.service.  
* Для запука сервиса перезагрузите systemd

```bash
sudo systemctl daemon-reload
```

* Запуск сервиса

```bash
service layoutconfig.api start
```

* Проверка статуса

```bash
service layoutconfig.api status
```

* Для запуска сервиса при запске сервера:

```bash
service layoutconfig.api enable
```

## требования к окружению для работы компонента

* Требуется наличие конфигурационного файла `config.yaml` в той же папке, где и сам исполняемый файл или запуск с флагом `-c=/path/to/config.yaml`
* Сетевой доступ к хранилищу (СУБД) той БД, путь к котрой прописан в конфиге или который получен от commonapi
* Сетевой доступ к CONSUL серверу, сервис пытается зарегистрироваться при запуске
* Если параметры доступа к хранилищам требуется получать из commonapi. тогда требуется иметь к нему доступ: url и token

В рамках проекта с Х5 были разработаны методы запроса данных из бд событий и скриншотов, в обычных вариантах они не нужны и их следут отсавить отключенными

```yaml
devicemanager:
  isuse: true # флаг, использовтаь или нет данную БД в работе
  ...
events:
  isuse: true # флаг, использовтаь или нет данную БД в работе
  ...
```

## описание параметров конфигурации компонента

> конфигурация стандартно файл > переменные окружения > флаги

Префикс для переменных окружения `LAYOUT`, тогда если задана переменная окружения `LAYOUT_COUNTMAX_SOURCE=api`, она будет замещать значение из файла конфигурации

```yaml
countmax:
    source: "api"
```

Флаги, значения которых используются:

* level - уровень логирования
* logfile - путь в лог файлу
* c - путь к файлу конфигурации
* consul - адрес:TCPport CONSUL сервера
* p - host:port на котором будет отвечать основное API
* sp - host:port на котором будет отвечать сервисное API (/health, /metrics, /pprof)

Конфигурационный файл содержит комментарии, объясняющие смысл каждого поля

параметр countmax/ids - перечисление кодов 1С клиентов, по которым будут запращиваться строки подклчения у commonapi  
в связи с тем, что возможно наличе нескольких баз данных по одному клиенту, после кода 1С через двоеточие можно указать тип сервера БД (см [CM_Info].[dbo].[DBTypes])

```yaml
app: # метаданные приложения
  name: "layoutconfig.api" # наименование приложения
httpd:
  main:
      host_port: ":8001" # хост:порт, который будет пытаться открыть приложения и принимать на него http запросы
      static: "assets" # тестовая статика, на работу не влияет
      allow_origin: # разрешение cors запросов
        - "*"
  service:
      host_port: ":8002" # хост:порт, который будет пытаться открыть приложения и принимать на него http запросы для метрик и проверки здоровья
countmax:
  source: "api" # источник коннектов к базам данных countmax, может быть config - конфигурационный файл или api - url  к ендпоинту откуда получить коннекты по списку
  url:  "http://dworker-01.watcom.local:7001" #"sqlserver://root:master@study-app:1433?database=CM_Karpov523&connection_timeout=0&encrypt=disable" #"postgres://retail:retail@localhost:5432/retail?sslmode=disable&pool_max_conns=2" #"sqlserver://root:master@study-app:1433?database=CM_GribMall523&connection_timeout=0&encrypt=disable"
  token: "eyJpc3MiOiJ0b3B0YWwuY29" # токен авторизации к АПИ commonapi для получения коннектов
  version: countmax523 # верися БД countMax: countmax523/countmax600
  timeout: 30s # timeout с которым будут работать запросы к БД
  ids: 1000001:2,1000002 #123,456,789 # коды 1С клиентов для поиска по ним строк подключения через commonapi
permissions:
  policy: allow # allow|deny политика по умолчанию, если не передаются права пользователя в запросе
  cache:
    url: memory # memory|redis url TODO: redis implementation
    expire: 59m
devicemanager:
  isuse: true # флаг, использовтаь или нет данную БД в работе
  url: "postgres://devmgr:devmgr@elk-01.watcom.local:35432/device.manager?sslmode=disable&pool_max_conns=2" #"sqlserver://root:master@study-app:1433?database=CM_Karpov523&connection_timeout=0&encrypt=disable" #"postgres://retail:retail@localhost:5432/retail?sslmode=disable&pool_max_conns=2" #"sqlserver://root:master@study-app:1433?database=CM_GribMall523&connection_timeout=0&encrypt=disable"
  timeout: 30s # timeout с которым будут работать запросы к БД
  aliases:
    src: http://cdn.countmax.ru:9001
    dest:
      - https://s3.watcom.ru
      - https://s3.countmax.ru
events:
  isuse: true # флаг, использовать или нет данную БД в работе
  url: "postgres://events:events@elk-01.watcom.local:45432/events?sslmode=disable&pool_max_conns=2" #"sqlserver://root:master@study-app:1433?database=CM_Karpov523&connection_timeout=0&encrypt=disable" #"postgres://retail:retail@localhost:5432/retail?sslmode=disable&pool_max_conns=2" #"sqlserver://root:master@study-app:1433?database=CM_GribMall523&connection_timeout=0&encrypt=disable"
  timeout: 30s # timeout с которым будут работать запросы к БД
env: production # тип окружения в котором запускается сервис, production - логи в json формате, все отсальное обычный logrus формат, котрый лучше выводить в текстовый файл и смотреть VSCode-ом
log:
  level: debug # уровень логирования сервиса: debug, info, warn, error
  file: "" # имя файла лога, если пусто или stdout - будет выводить в stdout, если указано имя фацйла, будет писать в него
consul:
  url: "elk-01:8500" # адрес consul сервера
  serviceid: "layoutconfig.api-dev" # уникальный идентификатор сервиса, соответсвует имени контейнера (имена контейнеров во всей системе не должны совпадать)!
  address: "elk-01.watcom.local" # адрес/fqnd имя сервера по которому будет видент данный сервис, host docker машины
  port: 8000 # порт по которому доступны метрики и проверка здоровья сервиса снаружи, service.host_port
tags:
  "develop,countmax,layoutconfig.api,office" # теги сервиса по которым будет осущестялться поиск и разметка в мониторинге, количетсво и порядок строго определенные
  #- develop # №1 окружение: develop, stage, production
  #- countmax # №2 проект откуда сервис: countmax, grib, focus etc...
  #- layoutconfig.api # №3 семейство сервисов: commonapi, dbscanner, incidentmaker, transport.webui etc...
  #- office # №4 локация/датацентр где работает сервис
```

## особенности публикации и эксплуатации компонента

имеет стандартный набор метрик для Prometheus-a `/metrics`  
при запуске регистрируется в consul-e для service discovering-a  
