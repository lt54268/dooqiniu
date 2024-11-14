# 七牛云 Kodo
运行：go run main.go

API文档：http://127.0.0.1:9090/swagger/index.html

## 一、上传接口（POST）
http://127.0.0.1:9090/api/v1/upload

参数：filePath、objectName

返回示例：
```
{
    "code": 200,
    "data": {
        "content-length": 10089,
        "etag": "FiVXXXXXXXXXXXXXXXXXXXXXXXXXXX",
        "last-modified": "2024-11-14T06:25:50Z"
    },
    "msg": "上传成功"
}
```

## 二、下载接口（GET）
http://127.0.0.1:9090/api/v1/download

参数：objectName

返回示例：返回文件下载链接
```
{
    "code": 200,
    "data": {
        "downloadURL": "http://xxxxxxx.clouddn.com/xxxxxxxxxxxxxx.xxxx?e=xxxxxxxxx&token=xxxxxxxxxxxxxxxxxxxx="
    },
    "msg": "生成下载链接成功"
}
```

## 三、删除接口（DELETE）
http://127.0.0.1:9090/api/v1/delete

参数：objectName

返回示例：
```
{
    "code": 200,
    "msg": "文件删除成功"
}
```

## 四、获取文件列表接口（GET）
http://127.0.0.1:9090/api/v1/list

参数（可选）：
prefix（返回的文件前缀，留空默认全部返回）
marker（游标，列举时继续读取上次的marker）
limit（每次返回的文件数量，默认一次返回1000条数据）

返回示例：
```
{
    "code": 200,
    "files": [
        {
            "key": "10.14会议纪要.docx",
            "content-length": 13769,
            "etag": "Fg9XXXXXXXXXXXXXXXXXXXXXXXXX",
            "last_modified": "2024-11-07T03:47:15Z"
        },
        {
            "key": "11.05会议纪要.docx",
            "content-length": 14567,
            "etag": "FnGhUXXXXXXXXXXXXXXXXXXXXXXX",
            "last_modified": "2024-11-07T02:09:10Z"
        },
        {
            "key": "AnythingLLMdocx",
            "content-length": 61521,
            "etag": "FkXXXXXXXXXXXXXXXXXXXXXXXXXX",
            "last_modified": "2024-11-12T06:55:58Z"
        }
    ],
    "msg": "文件列表获取成功",
    "next_marker": ""
}
```

## 五、拷贝接口（POST）
http://127.0.0.1:9090/api/v1/copy

参数：srcObject、destObject

返回示例：只能在同一个桶内进行操作
```
{
  "code": 200,
  "msg": "文件拷贝成功"
}
```

## 六、移动接口（POST）
http://127.0.0.1:9090/api/v1/move

参数：rcObject、destObject

返回示例：只能在同一个桶内进行操作
```
{
  "code": 200,
  "msg": "文件移动成功"
}
```

### 说明：
七牛云Kodo对象存储，上传同名文件会无法覆盖，需要先删除再上传。
