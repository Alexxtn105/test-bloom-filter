# Тестовый проект с использованием [Bloom Filters](https://github.com/bits-and-blooms/bloom)
Особенности bloom-фильтров:
 - Дают гарантированный ответ о том, что искомая величина **отсутствует** в наборе (false negative);
 - Не гарантируют то, что искомая величина **присутствует** в наборе (false positive);
 - высокая скорость и минимум расхода памяти.

Идеально подходит для:
 - дедупликация;
 - кэширование;
 - предварительная фильтрация большого набора данных.

Например, при использовании API часто требуется обращаться к БД на предмет наличия у пользователя прав доступа. Здесь на выручку приходит bloom filter, для сокращения количества запросов к БД.

Установка в Go:
```bash
go get -github github.com/bits-and-blooms/bloom/v3
```


Проверка:

http://localhost:8080/feature_access?user_id=99&feature_id=3

http://localhost:8080/estimate_fp?n=10&desired_fp_rate=0.001
