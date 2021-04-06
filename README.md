# photo-dir-date

ディレクトリ名から取得した日時を画像ファイルのオリジナル撮影日時に埋め込みます。

例えば「2021-04-06T195000」というディレクトリに格納されているJPEGファイルには「2021-04-06　19:50:00」と設定されます。

ディレクトリ名は正規表現 `(\d{4})-(\d{2})-(\d{2})[^\d]*(\d{2})(\d{2})(\d{2})` にマッチしている必要があります。

ディレクトリ内にあるファイルはファイル名の昇順でSubSecTimeOriginalが割り当てられます。

**Example**

```
$ ls ~/2021-04-06T195000
PICT0000.jpg    PICT0001.jpg    PICT0002.jpg    PICT0003.jpg    PICT0004.jpg
```

```
$ photo-dir-date -d ~/2021-04-06T195000  
```

```
$ exiftool -DateTimeOriginal -SubSecTimeOriginal ~/2021-04-06T195000
======== PICT0000.jpg
Date/Time Original              : 2021:04:06 19:50:00
Sub Sec Time Original           : 0
======== PICT0001.jpg
Date/Time Original              : 2021:04:06 19:50:00
Sub Sec Time Original           : 1
======== PICT0002.jpg
Date/Time Original              : 2021:04:06 19:50:00
Sub Sec Time Original           : 2
======== PICT0003.jpg
Date/Time Original              : 2021:04:06 19:50:00
Sub Sec Time Original           : 3
======== PICT0004.jpg
Date/Time Original              : 2021:04:06 19:50:00
Sub Sec Time Original           : 4
    5 image files read
```

## Install

実行には [exiftool](https://exiftool.org/) が必要です。 インストールしてパスが通っている場所に配置してください。

