# photo-dir-date

ディレクトリ名から取得した日時を画像ファイルのオリジナル撮影日時に埋め込みます。

例えば「2021-04-06T195000」というディレクトリに格納されているJPEGファイルには「2021-04-06　19:50:00」と設定されます。

ディレクトリ名は正規表現 `(\d{4})-(\d{2})-(\d{2})[^\d]*(\d{2})(\d{2})(\d{2})` にマッチしている必要があります。

ディレクトリ内にあるファイルはファイル名の昇順に並べた際に先頭となるファイルのファイル作成日時との差分を足した時刻が設定されます。
差分を足したくない場合は`--set-flat-time`フラグを利用してください。

**Example**

```
$ ls ~/2021-04-08t165200
PICT0000.jpg    PICT0001.jpg    PICT0002.jpg    PICT0003.jpg    PICT0004.jpg    PICT0005.jpg    PICT0006.jpg    PICT0007.jpg    PICT0008.jpg    PICT0009.jpg    PICT0010.jpg
```

```
$ photo-dir-date -d ~/2021-04-08t165200
FilePath                                               ModifiedTime                          CreateTime                     
/Volumes/NO NAME/PHOTO/2021-04-08t165200/PICT0000.jpg  2021-04-08 16:52:00.000000000 +09:00  2016-08-26 11:31:26 +0900 JST  
/Volumes/NO NAME/PHOTO/2021-04-08t165200/PICT0001.jpg  2021-04-08 16:52:06.000000000 +09:00  2016-08-26 11:31:32 +0900 JST  
/Volumes/NO NAME/PHOTO/2021-04-08t165200/PICT0002.jpg  2021-04-08 16:52:06.000000000 +09:00  2016-08-26 11:31:32 +0900 JST  
/Volumes/NO NAME/PHOTO/2021-04-08t165200/PICT0003.jpg  2021-04-08 16:52:08.000000000 +09:00  2016-08-26 11:31:34 +0900 JST  
/Volumes/NO NAME/PHOTO/2021-04-08t165200/PICT0004.jpg  2021-04-08 16:52:20.000000000 +09:00  2016-08-26 11:31:46 +0900 JST  
/Volumes/NO NAME/PHOTO/2021-04-08t165200/PICT0005.jpg  2021-04-08 16:52:38.000000000 +09:00  2016-08-26 11:32:04 +0900 JST  
/Volumes/NO NAME/PHOTO/2021-04-08t165200/PICT0006.jpg  2021-04-08 16:52:42.000000000 +09:00  2016-08-26 11:32:08 +0900 JST  
/Volumes/NO NAME/PHOTO/2021-04-08t165200/PICT0007.jpg  2021-04-08 16:53:50.000000000 +09:00  2016-08-26 11:33:16 +0900 JST  
/Volumes/NO NAME/PHOTO/2021-04-08t165200/PICT0008.jpg  2021-04-08 16:53:56.000000000 +09:00  2016-08-26 11:33:22 +0900 JST  
/Volumes/NO NAME/PHOTO/2021-04-08t165200/PICT0009.jpg  2021-04-08 16:54:30.000000000 +09:00  2016-08-26 11:33:56 +0900 JST  
/Volumes/NO NAME/PHOTO/2021-04-08t165200/PICT0010.jpg  2021-04-08 16:55:18.000000000 +09:00  2016-08-26 11:34:44 +0900 JST  
 :
 :
```

```
$ exiftool -DateTimeOriginal ~/2021-04-08t165200
======== PICT0008.jpg
Date/Time Original              : 2021:04:08 16:53:56
======== PICT0009.jpg
Date/Time Original              : 2021:04:08 16:54:30
======== PICT0001.jpg
Date/Time Original              : 2021:04:08 16:52:06
======== PICT0000.jpg
Date/Time Original              : 2021:04:08 16:52:00
======== PICT0002.jpg
Date/Time Original              : 2021:04:08 16:52:06
======== PICT0003.jpg
Date/Time Original              : 2021:04:08 16:52:08
======== PICT0007.jpg
Date/Time Original              : 2021:04:08 16:53:50
======== PICT0006.jpg
Date/Time Original              : 2021:04:08 16:52:42
======== PICT0010.jpg
Date/Time Original              : 2021:04:08 16:55:18
======== PICT0004.jpg
Date/Time Original              : 2021:04:08 16:52:20
======== PICT0005.jpg
Date/Time Original              : 2021:04:08 16:52:38
```

## Install

実行には [exiftool](https://exiftool.org/) が必要です。 インストールしてパスが通っている場所に配置してください。

最新のバイナリは [Release](https://github.com/kasai-y/photo-dir-date/releases) よりダウンロードしてください。