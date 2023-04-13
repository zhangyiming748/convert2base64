# convert2base64
图片转换为企业微信支持的base64编码

查找同级目录下jpg和png文件

支持多个文件

分别生成

+ base64字符串
+ MD5字符串
+ json请求体
# 使用方法

以Linux为例

同级目录下需要放置配置文件`setting.ini`
格式如下

```ini
[log]
level = Debug
# level = Info
# level = Warn
# level = Error

```

```bash
chmod a+x convertBase64ForLinux
./convertBase64ForLinux
```

会自动查找同级目录下的jpg和png文件,生成对应的base64文件、md5文件和json结构体