# golang-loki-adapter 
Сервис для сбора логов из журнала событий CMS Bitrix в Loki Grafana 

# Установка зависимостей на CentOS
Шаги установки ядра Go из под root:
1.	`wget https://dl.google.com/go/go1.21.3.linux-amd64.tar.gz`
2.	`sha256sum go1.21.3.linux-amd64.tar.gz`
3.	`sudo tar -C /usr/local -xzf go1.21.3.linux-amd64.tar.gz`
4.	`export PATH=$PATH:/usr/local/go/bin export GOROOT=/usr/local/go export GOPATH=$HOME/Documents/go`
5.  `source /etc/profile`

Должно работать `go version`.

# Создание очереди в mysql
Любыми средствами добавить очередь:
```sql
CREATE TABLE IF NOT EXISTS `bpm_event_log_queue` (
`ID` int unsigned NOT NULL AUTO_INCREMENT,
`DATA` text COLLATE utf8mb4_general_ci DEFAULT NULL,
`CREATED_AT` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
`ATTEMPTS` INT NOT NULL DEFAULT 0,
`LAST_ATTEMPT_AT` datetime
PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
```
И триггер для заполнения очереди из журнала событий:
```sql
CREATE TRIGGER after_insert_b_event_log
    AFTER INSERT ON b_event_log
    FOR EACH ROW
    INSERT INTO bpm_event_log_queue (DATA)
    VALUES (
               CONCAT_WS("",
                       '{"ID":"', NEW.ID,
                       '","SEVERITY":"', NEW.SEVERITY,
                       '","MODULE_ID":"', NEW.MODULE_ID,
                       '","ITEM_ID":"', NEW.ITEM_ID,
                       '","REQUEST_URI":"', NEW.REQUEST_URI,
                       '","AUDIT_TYPE_ID":"', NEW.AUDIT_TYPE_ID,
                       '","DESCRIPTION":"', NEW.DESCRIPTION,
                       '","TIMESTAMP_X":"', NEW.TIMESTAMP_X,
                       '","DATE":"', NOW(), '"}\n'
               )
    );
```

Дампы находятся в проекте `install/dump.sql`.

# Настройка конфигураций
В `internal/config` нужно создать файл `config.yaml` с настройками среды:
1. Подключение к БД
2. Настройки для Loki

# Настройка systemd unit файла
1. touch /etc/systemd/system/loki-adapter.service
2. echo '[Unit]\nDescription=Event log loki log adapter\nAfter=default.target\n\n[Service]\nUser=root \nRestart=always \nEnvironment="CONFIG_PATH=/home/bitrix/www/local/services/golang-loki-adapter/internal/config/config.yaml"\nExecStart=/home/bitrix/www/local/services/golang-loki-adapter/golang-loki-adapter.local\n\n[Install]\nAlias=golang-loki-adapter.service \nWantedBy=nginx.service' > /etc/systemd/system/loki-adapter.service
3. ln -s loki-adapter.service golang-loki-adapter.service

Результат
```
[Unit]
Description=Event log loki log adapter
After=default.target

[Service]
User=root
Restart=always
Environment="CONFIG_PATH=/home/bitrix/www/local/services/golang-loki-adapter/internal/config/config.yaml"
ExecStart=/home/bitrix/www/local/services/golang-loki-adapter/golang-loki-adapter.local

[Install]
Alias=golang-loki-adapter.service
WantedBy=nginx.service
```


`CONFIG_PATH` настроить под расположение config.yaml.

Полезные команды:
1. systemctl start loki-adapter.service
2. systemctl restart loki-adapter.service
3. systemctl status loki-adapter.service