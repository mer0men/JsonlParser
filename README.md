#JsonlParser

Тестовое задание Mail.ru.

В качестве задания требовалось обработать семпл выгрузки из хадупа.\
Строка хронит в себе json объект формата:
```
{
	Url             string   `json:"url"`
	State           string   `json:"state"`
	Categories      []string `json:"categories"`
	CategoryAnother string   `json:"category_another"`
	ForMainPage     bool     `json:"for_main_page"`
	Ctime           int64    `json:"ctime"`
}
```
Каждую строку требуется распарсить и скачать html по Url, для того, чтобы достать от туда title и description,
после чего сохранить данные в `Categories[n].tsv` в формате `Url\ttitle\tdescription`

## Запуск

`go run main.go`

## Процесс выполнения программы

 На вход мы получаем файл `500.jsonl`, который попадает в `JsonlSiteReader`, где считывается построчно.\
 Каждая строка проходит через `json.Unmarshal` и попадает в канал(`chan Site`).\
 `JsonlSiteReader` возвращает `chan Site`, который далее попадает в `SiteReceiver.Recive(ctx, chan Site)` запущенную в 
 10 горутинах.\
 В `Recive` берется `Site` из канала, после чего делается `GET` запрос по `Site.Url`, через `http.Client` который берется
 из `SiteReciver.HttpClientPool`(`sync.Pool`)\
 Ответ запроса парсится через `goquery` для получения `title` и `description`.\
 После формируется строка формата `Site.Url\ttitle\tdescription`, которая записывается в буфер(`bytes.Buffer`) через `SiteWriter`
 с использованием `sync.Mutex` для синхронизации потоков.\
 `SiteWriter` хранится в `SiteReciver.SiteWriters`(`sync.Map`) для каждой категории(`Site.Categories[n]`).\
 Если в мапе нет `SiteWriter` для полученной категории, то создается новый `SiteWriter` в который складывается буфер
 из `SiteReciver.BufferPool`(`sync.Pool`).\
 После получения `SiteWriter`, проверятся длинна буфера(`SiteWriter.Length()`),
 если длина меньше 4096(байт), то производится запись(`SiteWriter.Write(str)`),
 иначе произойдет ротация буферов, старый(заполненный) отправиться на запись в файл(`{Site.Categories[n]}.tsv`),
 а новый будет получен из `SiteReciver.BufferPool`(`sync.Pool`).\
 После того как все сообщения из канала будут обработаны, произойдет запись данных в файлы, что остались в буферах.

