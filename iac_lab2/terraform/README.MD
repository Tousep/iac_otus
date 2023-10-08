# Домашняя работа №2
Домашняя работа №2 курса Infrastructure As Code от OTUS. Выполняется по 
[методическим указаниям](https://hackmd.io/@otus/rkCER8MfK#%D0%94%D0%BE%D0%BC%D0%B0%D1%88%D0%BD%D0%B5%D0%B5-%D0%B7%D0%B0%D0%B4%D0%B0%D0%BD%D0%B8%D0%B5-%D0%BA-%D0%B7%D0%B0%D0%BD%D1%8F%D1%82%D0%B8%D1%8E-%E2%80%9CTerraform-%D0%BA%D0%B0%D0%BA-%D0%B8%D0%BD%D1%81%D1%82%D1%80%D1%83%D0%BC%D0%B5%D0%BD%D1%82-%D0%B4%D0%BB%D1%8F-%D0%B4%D0%B5%D0%BA%D0%BB%D0%B0%D1%80%D0%B0%D1%82%D0%B8%D0%B2%D0%BD%D0%BE%D0%B3%D0%BE-%D0%BE%D0%BF%D0%B8%D1%81%D0%B0%D0%BD%D0%B8%D1%8F-%D0%B8%D0%BD%D1%84%D1%80%D0%B0%D1%81%D1%82%D1%80%D1%83%D0%BA%D1%82%D1%83%D1%80%D1%8B%E2%80%9D).

## Генерация ключа
Для доступа к ВМ генерируем ключи командой. 
```
ssh-keygen -t ed25519
```

## Задание со звездочкой
Общение с YC, используя токен не самый надежный и удобный способ, поэтому попробуйте перейти на использования доступа через сервисный аккаунт.

Порядок действий:
* Создаем аккаунт.
```sh
yc iam service-account get --name iac-sacc --folder-id b1gceigr50kq5nudfjfb
```
* Добавляем права.
```sh
yc resource-manager folder add-access-binding --id b1gceigr50kq5nudfjfb --role editor --service-account-id ajeca8ve3noid3a9mebd
```
* Достаем key.json.
```sh
yc iam key create --service-account-id ajeca8ve3noid3a9mebd --output key.json
```

## Задание с двумя звездочками.
** 
Обратите внимание, что в манифесте группы хостов для балансировщика мы явно указывали отдельные блоки target с именами виртуальных машин

address = yandex_compute_instance.wp-app-1.network_interface.0.ip_address

Понятно, что это может быть не очень удобное, если мы захотим уменьшить или увеличить кол-во хостов, где будет развернут WordPress.

Попробуйте изменить манифесты создания виртуальных машин и балансировщика, чтобы кол-вом хостов можно было управлять при помощи переменной и при этом не требовалось бы вносить изменения в манифесты терраформа. 
**

Я пытался...
Но так и не нашел способа реализовать подобную динамику средствами Terraform.