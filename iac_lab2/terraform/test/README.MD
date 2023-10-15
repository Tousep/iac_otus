# Домашняя работа №3
Домашняя работа №3 курса Infrastructure As Code от OTUS. Выполняется по 
[методическим указаниям](https://hackmd.io/@otus/B1TK3dUMK).

## Цель работы
Создать тесты, проверяющие корректно ли созданы ресурсы при помощи терраформ-манифестов.

Проверки:

* наличия IP Load balancer-а в state-е терраформа
* возможность подключиться по ssh к одной из виртуальных машин
* возможность подключения к базе данных (это задание “со звездочкой”)

## Запуск тестов.
```
go test -v ./ -timeout 30m -folder 'b1geb6a2k7nj5robi0h6' -ssh-key-pass 'C:/Users/User/ssh/srt' -db-ssh-key-pass 'C:/Users/User/mysql/crt'
```

## Задание со звездочкой
Дополните тесты проверкой доступа к базе данных.

Параметры подключения в БД.
```
// Читаем список хостов из output-переменных терраформа(value = tolist(local.dbhosts))
databaseHostFQDNs := terraform.OutputList(t, terraformOptions, "database_host_fqdn")
// Переменные окружения
dbuser := os.Getenv("DBUSER")
dbpass := os.Getenv("DBPASS")
// Константы
net := "tcp"
dbname := "db"
port := 3306
```

Проверка заведения переменных окружения.
```
if dbuser == "" || dbpass == "" {
    t.Fatal("env DBUSER or DBPASS not found!!")
}
```

Регистрация TLS сертификата для подключения к БД.
```
rootCertPool := x509.NewCertPool()
dbkey, err := ioutil.ReadFile(*dbSshKeyPath)
if err != nil {
    t.Fatalf("Unable to read private key -db-ssh-key-pass: %v", err)
}

if !rootCertPool.AppendCertsFromPEM(dbkey) {
    t.Fatalf("AppendCertsFromPEM err")
}

mysql.RegisterTLSConfig("custom", &tls.Config{
    RootCAs: rootCertPool,
})
```

Цикл тестов по хостам кластера БД(подключение + тестовый запрос)
```
for _, host := range databaseHostFQDNs {
mysqlInfo := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?tls=custom",
    dbuser, dbpass, net, host, port, dbname)

conn, err := sql.Open("mysql", mysqlInfo)
if err != nil {
    t.Fatalf("ERROR! %s: %v", host, err)
}

fmt.Println(fmt.Sprintf("Connected! %s", host))
defer conn.Close()

q, err := conn.Query("SELECT version()")
if err != nil {
    t.Fatalf("ERROR! %s: %v", host, err)
}

var res string

for q.Next() {
    q.Scan(&res)
    fmt.Printf(fmt.Sprintf("MySql version %s: %s\n", host, res))
}
```

