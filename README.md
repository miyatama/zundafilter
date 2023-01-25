![LOGO](./zundafilter_logo.png)
<h2 align="center">全てのテキストをずんだへ</h2>

<p align="center">
<a alt="code: Golang" href="https://go.dev/">
  <img src="https://img.shields.io/badge/code-go-ff69b4.svg">
</a>
<a alt="MIT License" href="https://kawakawaritsuki.mit-license.org/">
  <img src="https://img.shields.io/badge/license-MIT-blue.svg">
</a>
</p>

# build and run

1.build

```shell
cd ./data
./setup.sh
cd ..
make clean
make build
```

2.run

```shell
echo "継ぎます" | ./bin/zundaFilter
```

# attention

file for test

+ `./data/convert.sh`
+ `./data/source.txt`
