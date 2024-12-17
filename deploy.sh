#!/bin/bash
# exit on errors
set -e

# Build the project.
hugo --minify # if using a theme, replace with `hugo -t <YOURTHEME>`

cd public

git init
git config user.name 'shipengqi'
git config user.email 'pooky.shi@gmail.com'
git add -A

# Commit changes.
msg="rebuilding site `date`"
if [ $# -eq 1 ]
  then msg="$1"
fi
git commit -m "$msg"

# Push source and build repos.
git push -f git@github.com:shipengqi/golang-learn.git master:gh-pages

# 如果发布到 https://<USERNAME>.github.io
# git push -f git@github.com:<USERNAME>/<USERNAME>.github.io.git master

# 如果发布到 https://<USERNAME>.github.io/<REPO>
# git push -f git@github.com:shipengqi/golang-learn.git master:gh-pages
