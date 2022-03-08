# mediastorage backend  

## /media?[cursor=\<string>]

получение списка всех изображений и видео постранично. одна страница - 50 элементов
ответ:  

```JSON
{
  "media": [
    {
      "thumb_url": "http://localhost:8080/media/LnRtcC9kaXIxL2kxLmpwZw",
      "detail_url": "http://localhost:8080/media/LnRtcC9kaXIxL2kxLmpwZw",
      "original_url": "http://localhost:8080/media/LnRtcC9kaXIxL2kxLmpwZw"
    },
    ...
  ],
  "cursor": ""
}
```

## /media/{id}

получение конкретного файла
