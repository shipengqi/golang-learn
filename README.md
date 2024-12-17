# golang-learn

Go language learning ...

## Usage

Initialize project as a Hugo Module using the `hugo mod init` command:

```sh
hugo mod init github.com/shipengqi/golang-learn
```

Adding the [Hextra](https://github.com/imfing/hextra) theme:

```yaml
baseURL: "https://shipengqi.github.io/golang-learn"
title: "Golang Learning"

module:
  imports:
    - path: github.com/imfing/hextra
```

Start server:

```sh
hugo server -D
```

Manually deploy:

```sh
./deploy.sh
```

> Any changes in the `content` directory will automatically trigger a deployment.

## Update Theme

To update all Hugo modules in your project to their latest versions, run the following command:

```
hugo mod get -u
```

To update Hextra to the latest released version, run the following command:

```
hugo mod get -u github.com/imfing/hextra
```
