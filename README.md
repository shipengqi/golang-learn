# golang-learn

Golang language learning ...

## Usage

Development:

```sh
# install hugo-book as git submodule
git submodule add git@github.com:alex-shpak/hugo-book.git themes/book

hugo server --minify --theme book
```

Manually deploy:

```sh
./deploy.sh
```

> Any changes in the `content` directory will automatically trigger a deployment.

## Menu

By default, the [hugo-book](https://github.com/alex-shpak/hugo-book) theme will render pages from the `content/docs` section as a menu in a tree structure.
You can set `title` and `weight` in the front matter of pages to adjust the order and titles in the menu.