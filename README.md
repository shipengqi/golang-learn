# golang-learn

Golang language learning ...

## Usage

development:

```sh
# install hugo-book as git submodule
git submodule add git@github.com:alex-shpak/hugo-book.git themes/book

hugo server --minify --theme book
```

deploy:

```sh
./deploy.sh
```

## Menu

By default, the hugo-book theme will render pages from the `content/docs` section as a menu in a tree structure.
You can set title and weight in the front matter of pages to adjust the order and titles in the menu.