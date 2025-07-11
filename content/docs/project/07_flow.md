---
title: Git 工作流程
weight: 7
---

涉及到多人协作的项目，多个开发者向同一个仓库提交代码，如果处理不好会出现代码丢失，冲突等问题。所以一个规范的工作流程，可以让开发者更有效地合作，使项目更好地发展下去。

最常用的工作流程有三种：

- Git Flow
- GitHub Flow
- Forking Flow

## Git Flow

Git Flow 是最早出现的一种工作流程。

Git Flow 存在两种长期分支：

- `master`：这个分支永远是稳定的发布版本，不能直接在该分支上开发。每次合并一个 hotfix/release 分支，都在 `master` 上打一个版本标签。
- `develop`：日常开发的分支，存放最新的开发版。同样不能在这个分支上直接开发，这个分支只做合并操作。

三种短期分支：

- feature branch：用于功能开发，基于 `develop` 创建新的 feature 分支，可以命名为 `feat/xxx-xx`。开发完成之后，合并到 `develop` 并删除。
- hotfix branch：补丁分支，在维护阶段用于紧急的 bug 修复。基于 `master` 创建，可以命名为 `hotfix/xxx-xx`。完成后合并到 `master` 分支并，然后在 `master` 打上标签删除并删除 hotfix 
分支。一般 `develop` 也需要合并 hotfix 分支。
- release branch：预发布分支，在发布阶段，基于 `develop` 创建，可以命名为 `release/xxx-xx`。 例如 `v1.0.0` 版本开发完成后，代码已经全部合并到 `develop` 分支。发布之前，基于 
`develop` 创建`release/1.0.0` 分支，基于 `release/1.0.0` 进行测试，如果发现 bug，就在 `release/1.0.0` 分支上修复。测试完成后，合并到 `master` 和 `develop` 分支。然后在 `master`
打上标签，并删除 `release/1.0.0` 分支。

> 这三种短期分支会在开发完成后合并到 develop 或者 master，然后删除。

![git-flow](https://raw.gitcode.com/shipengqi/illustrations/files/main/go/git-flow.png)

Git flow 的优点是每个分支分工明确，可以最大程度减少它们之间的相互影响。但是需要同时维护两个长期分支，相对比较复杂，需要经常在 `master` 分支 `develop` 分支进行切换。

## GitHub Flow

GitHub Flow 是 Git flow 的简化版。只要一个长期分支 `master`。

流程：

1. 基于 `master` 创建新的 feature/hotfix 分支。
2. 开发完成后，向 `master` 分支发起一个 pull request（PR）。
3. PR 需要 review，review 过程中可以不断的提交代码进行修改。
4. PR 被 approve ，然后合并到 `master` 并删除 feature/hotfix 分支。

GitHub Flow 非常简单，适合持续发布的产品，`master` 分支就是当前的线上代码。

但是有时候代码合并到 `master` 并不代表就可以发布了，比如，苹果商店的 APP 提交审核以后，等一段时间才能上架。这时，如果还有新的代码提交，`master` 分支就会与刚发布的版本不一致。这种情况，
只有 `master` 一个主分支就不够用了。通常，不得不在 `master` 分支以外，另外新建一个 `production` 分支跟踪线上版本。

### Forking Flow

开源项目中常用的是 Forking Flow，比如 Kubernetes。Forking Flow 在 Git Flow 的基础上充分利用了 Git 的 Fork 和 pull request 的功能以达到代码审核的目的。可以安全可靠地管理大团队的开发者，并能接受不信任贡献者的提交。

Forking Flow 和 GitHub Flow 是差不多的：

1. Fork 项目到自己的仓库。
2. 开发完成后，推送到自己的仓库。
3. 向上游仓库发起 PR。
4. PR 需要 review，review 过程中可以不断的提交代码进行修改。 
5. PR 被 approve，然后合并到上游仓库。

和 Github Flow 的区别就是没有创建新分支，而是创建了一个新的 fork。
