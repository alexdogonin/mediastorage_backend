# mediastorage backend  

## How to run

```ROOT_PATH=/home/user/Images PORT=3000 server```

## Envirenment

**ROOT_PATH** - directory to scan. Required  
**PORT** - port to listen. Required  
**ADDRESS** - IP address to listen. Default: *0.0.0.0*  
**SCHEME** - scheme. Default: *http*  

## Endpoints
### __/media?[cursor=\<string>]__

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

### __/media/{id}__

получение конкретного файла
