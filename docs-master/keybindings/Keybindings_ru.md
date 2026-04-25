_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go generate ./...` from the project root._

# Lazygit Связки клавиш

## Глобальные сочетания клавиш

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+r> `` | Переключиться на последний репозиторий |  |
| `` <pgup> (fn+up/shift+k) `` | Прокрутить вверх главную панель |  |
| `` <pgdown> (fn+down/shift+j) `` | Прокрутить вниз главную панель |  |
| `` @ `` | Открыть меню журнала команд | View options for the command log e.g. show/hide the command log and focus the command log. |
| `` P `` | Отправить изменения | Push the current branch to its upstream branch. If no upstream is configured, you will be prompted to configure an upstream branch. |
| `` p `` | Получить и слить изменения | Pull changes from the remote for the current branch. If no upstream is configured, you will be prompted to configure an upstream branch. |
| `` ) `` | Increase rename similarity threshold | Increase the similarity threshold for a deletion and addition pair to be treated as a rename.<br><br>The default can be changed in the config file with the key 'git.renameSimilarityThreshold'. |
| `` ( `` | Decrease rename similarity threshold | Decrease the similarity threshold for a deletion and addition pair to be treated as a rename.<br><br>The default can be changed in the config file with the key 'git.renameSimilarityThreshold'. |
| `` } `` | Увеличить размер контекста, отображаемого вокруг изменений в просмотрщике сравнении | Increase the amount of the context shown around changes in the diff view.<br><br>The default can be changed in the config file with the key 'git.diffContextSize'. |
| `` { `` | Уменьшите размер контекста, отображаемого вокруг изменений в просмотрщике сравнении | Decrease the amount of the context shown around changes in the diff view.<br><br>The default can be changed in the config file with the key 'git.diffContextSize'. |
| `` : `` | Execute shell command | Bring up a prompt where you can enter a shell command to execute. |
| `` <ctrl+p> `` | Просмотреть пользовательские параметры патча |  |
| `` m `` | Просмотреть параметры слияния/перебазирования | View options to abort/continue/skip the current merge/rebase. |
| `` mc `` | Continue rebase / merge | View options to abort/continue/skip the current merge/rebase. |
| `` ma `` | Abort rebase / merge | View options to abort/continue/skip the current merge/rebase. |
| `` ms `` | Skip current rebase commit | View options to abort/continue/skip the current merge/rebase. |
| `` R `` | Обновить | Refresh the git state (i.e. run `git status`, `git branch`, etc in background to update the contents of panels). This does not run `git fetch`. |
| `` + `` | Следующий режим экрана (нормальный/полуэкранный/полноэкранный) |  |
| `` _ `` | Предыдущий режим экрана |  |
| `` \| `` | Cycle pagers | Choose the next pager in the list of configured pagers |
| `` <esc> `` | Отменить |  |
| `` ? `` | Открыть меню |  |
| `` <ctrl+s> `` | Просмотреть параметры фильтрации по пути | View options for filtering the commit log, so that only commits matching the filter are shown. |
| `` W `` | Открыть меню сравнении | View options relating to diffing two refs e.g. diffing against selected ref, entering ref to diff against, and reversing the diff direction. |
| `` <ctrl+e> `` | Открыть меню сравнении | View options relating to diffing two refs e.g. diffing against selected ref, entering ref to diff against, and reversing the diff direction. |
| `` q `` | Выйти |  |
| `` <ctrl+z> `` | Suspend the application |  |
| `` <ctrl+w> `` | Переключить отображение изменении пробелов в просмотрщике сравнении | Toggle whether or not whitespace changes are shown in the diff view.<br><br>The default can be changed in the config file with the key 'git.ignoreWhitespaceInDiffView'. |
| `` z `` | Отменить (через reflog) (экспериментальный) | Журнал ссылок (reflog) будет использоваться для определения того, какую команду git запустить, чтобы отменить последнюю команду git. Сюда не входят изменения в рабочем дереве; учитываются только коммиты. |
| `` Z `` | Повторить (через reflog) (экспериментальный) | Журнал ссылок (reflog) будет использоваться для определения того, какую команду git нужно запустить, чтобы повторить последнюю команду git. Сюда не входят изменения в рабочем дереве; учитываются только коммиты. |

## Навигация по панели списка

| Key | Action | Info |
|-----|--------|-------------|
| `` , `` | Предыдущая страница |  |
| `` . `` | Следующая страница |  |
| `` < (<home>) `` | Пролистать наверх |  |
| `` > (<end>) `` | Прокрутить вниз |  |
| `` v `` | Переключить выборку перетаскивания |  |
| `` <shift+down> `` | Range select down |  |
| `` <shift+up> `` | Range select up |  |
| `` / `` | Найти |  |
| `` H `` | Прокрутить влево |  |
| `` L `` | Прокрутить вправо |  |
| `` ] `` | Следующая вкладка |  |
| `` [ `` | Предыдущая вкладка |  |

## Input prompt

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Подтвердить |  |
| `` <esc> `` | Закрыть/отменить |  |

## Worktrees

| Key | Action | Info |
|-----|--------|-------------|
| `` n `` | New worktree |  |
| `` <space> `` | Switch | Switch to the selected worktree. |
| `` o `` | Open in editor |  |
| `` d `` | Remove | Remove the selected worktree. This will both delete the worktree's directory, as well as metadata about the worktree in the .git directory. |
| `` / `` | Filter the current view by text |  |

## Вторичный

| Key | Action | Info |
|-----|--------|-------------|
| `` <tab> `` | Переключиться на другую панель (проиндексированные/непроиндексированные изменения) | Switch to other view (staged/unstaged changes). |
| `` <esc> `` | Exit back to side panel |  |
| `` / `` | Найти |  |

## Главная панель (Индексирование)

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | Выбрать предыдущую часть |  |
| `` <right> `` | Выбрать следующую часть |  |
| `` v `` | Переключить выборку перетаскивания |  |
| `` a `` | Toggle hunk selection | Toggle line-by-line vs. hunk selection mode. |
| `` <ctrl+o> `` | Скопировать выделенный текст в буфер обмена |  |
| `` <space> `` | Переключить индекс | Переключить строку в проиндексированные / непроиндексированные |
| `` d `` | Отменить изменение (git reset) | When unstaged change is selected, discard the change using `git reset`. When staged change is selected, unstage the change. |
| `` o `` | Открыть файл | Open file in default application. |
| `` e `` | Редактировать файл | Open file in external editor. |
| `` <esc> `` | Вернуться к панели файлов |  |
| `` <tab> `` | Переключиться на другую панель (проиндексированные/непроиндексированные изменения) | Switch to other view (staged/unstaged changes). |
| `` E `` | Изменить эту часть | Edit selected hunk in external editor. |
| `` c `` | Сохранить изменения | Commit staged changes. |
| `` w `` | Закоммитить изменения без предварительного хука коммита |  |
| `` C `` | Сохранить изменения с помощью редактора git |  |
| `` <ctrl+f> `` | Find base commit for fixup | Find the commit that your current changes are building upon, for the sake of amending/fixing up the commit. This spares you from having to look through your branch's commits one-by-one to see which commit should be amended/fixed up. See docs: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` / `` | Найти |  |

## Главная панель (Обычный)

| Key | Action | Info |
|-----|--------|-------------|
| `` <mouse wheel down> (fn+up) `` | Прокрутить вниз |  |
| `` <mouse wheel up> (fn+down) `` | Прокрутить вверх |  |
| `` <tab> `` | Переключиться на другую панель (проиндексированные/непроиндексированные изменения) | Switch to other view (staged/unstaged changes). |
| `` <esc> `` | Exit back to side panel |  |
| `` / `` | Найти |  |

## Главная панель (Слияние)

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | Выбрать эту часть |  |
| `` b `` | Выбрать все части |  |
| `` <up> `` | Выбрать предыдущую часть |  |
| `` <down> `` | Выбрать следующую часть |  |
| `` <left> `` | Выбрать предыдущий конфликт |  |
| `` <right> `` | Выбрать следующий конфликт |  |
| `` z `` | Отменить | Undo last merge conflict resolution. |
| `` e `` | Редактировать файл | Open file in external editor. |
| `` o `` | Открыть файл | Open file in default application. |
| `` M `` | View merge conflict options | View options for resolving merge conflicts. |
| `` <esc> `` | Вернуться к панели файлов |  |

## Главная панель (сборка патчей)

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | Выбрать предыдущую часть |  |
| `` <right> `` | Выбрать следующую часть |  |
| `` v `` | Переключить выборку перетаскивания |  |
| `` a `` | Toggle hunk selection | Toggle line-by-line vs. hunk selection mode. |
| `` <ctrl+o> `` | Скопировать выделенный текст в буфер обмена |  |
| `` o `` | Открыть файл | Open file in default application. |
| `` e `` | Редактировать файл | Open file in external editor. |
| `` <space> `` | Добавить/удалить строку(и) для патча |  |
| `` d `` | Remove lines from commit | Remove the selected lines from this commit. This runs an interactive rebase in the background, so you may get a merge conflict if a later commit also changes these lines. |
| `` <esc> `` | Выйти из сборщика пользовательских патчей |  |
| `` / `` | Найти |  |

## Журнал ссылок (Reflog)

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | Copy abbreviated commit hash to clipboard |  |
| `` <space> `` | Переключить | Checkout the selected commit as a detached HEAD. |
| `` y `` | Скопировать атрибут коммита | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | Открыть коммит в браузере |  |
| `` n `` | Создать новую ветку с этого коммита |  |
| `` N `` | Move commits to new branch | Create a new branch and move the unpushed commits of the current branch to it. Useful if you meant to start new work and forgot to create a new branch first.<br><br>Note that this disregards the selection, the new branch is always created either from the main branch or stacked on top of the current branch (you get to choose which). |
| `` g `` | Reset to ref |  |
| `` gm `` | Mixed reset | Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as unstaged changes. |
| `` gs `` | Мягкий сброс | Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as staged changes. |
| `` gh `` | Жёсткий сброс | Reset HEAD to the chosen commit, and discard all changes between the current and chosen commit, as well as all current modifications in the working tree. |
| `` C `` | Скопировать отобранные коммит (cherry-pick) | Mark commit as copied. Then, within the local commits view, you can press `V` to paste (cherry-pick) the copied commit(s) into your checked out branch. At any time you can press `<esc>` to cancel the selection. |
| `` <ctrl+r> `` | Сбросить отобранную (скопированную \| cherry-picked) выборку коммитов |  |
| `` <ctrl+t> `` | Open external diff tool (git difftool) |  |
| `` * `` | Select commits of current branch |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | Просмотреть коммиты |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## Коммиты

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | Copy abbreviated commit hash to clipboard |  |
| `` <ctrl+r> `` | Сбросить отобранную (скопированную \| cherry-picked) выборку коммитов |  |
| `` b `` | Просмотреть параметры бинарного поиска |  |
| `` bb `` | Mark commit as bad |  |
| `` bg `` | Mark commit as good |  |
| `` bs `` | Skip current bisect commit |  |
| `` bS `` | Skip selected commit |  |
| `` br `` | Сбросить бинарный поиск |  |
| `` bb `` | Mark commit as bad (start bisect) |  |
| `` bg `` | Mark commit as good (start bisect) |  |
| `` bt `` | Choose bisect terms |  |
| `` s `` | Объединить коммиты (Squash) | Squash the selected commit into the commit below it. The selected commit's message will be appended to the commit below it. |
| `` f `` | Объединить несколько коммитов в один отбросив сообщение коммита (Fixup)  | Meld the selected commit into the commit below it. Similar to squash, but the selected commit's message will be discarded. |
| `` ff `` | Объединить несколько коммитов в один отбросив сообщение коммита (Fixup)  | Meld the selected commit into the commit below it. Similar to squash, but the selected commit's message will be discarded. |
| `` fc `` | Fixup and use this commit's message | Squash the selected commit into the commit below, using this commit's message, discarding the message of the commit below. |
| `` c `` | Set fixup message | Set the message option for the fixup commit. The -C option means to use this commit's message instead of the target commit's message. |
| `` r `` | Перефразировать коммит | Reword the selected commit's message. |
| `` R `` | Переписать коммит с помощью редактора |  |
| `` d `` | Удалить коммит | Drop the selected commit. This will remove the commit from the branch via a rebase. If the commit makes changes that later commits depend on, you may need to resolve merge conflicts. |
| `` e `` | Edit (start interactive rebase) | Изменить коммит |
| `` i `` | Start interactive rebase | Start an interactive rebase for the commits on your branch. This will include all commits from the HEAD commit down to the first merge commit or main branch commit.<br>If you would instead like to start an interactive rebase from the selected commit, press `e`. |
| `` p `` | Pick | Выбрать коммит (в середине перебазирования) |
| `` F `` | Создать fixup коммит | Создать fixup коммит для этого коммита |
| `` S `` | Apply fixup commits | Объединить все 'fixup!' коммиты выше в выбранный коммит (автосохранение) |
| `` <alt+down> `` | Переместить коммит вниз на один |  |
| `` <alt+up> `` | Переместить коммит вверх на один |  |
| `` V `` | Вставить отобранные коммиты (cherry-pick) |  |
| `` B `` | Mark as base commit for rebase | Select a base commit for the next rebase. When you rebase onto a branch, only commits above the base commit will be brought across. This uses the `git rebase --onto` command. |
| `` A `` | Amend | Править последний коммит с проиндексированными изменениями |
| `` a `` | Установить/убрать автора коммита | Set/Reset commit author or set co-author. |
| `` t `` | Revert | Create a revert commit for the selected commit, which applies the selected commit's changes in reverse. |
| `` T `` | Пометить коммит тегом | Create a new tag pointing at the selected commit. You'll be prompted to enter a tag name and optional description. |
| `` <ctrl+l> `` | Открыть меню журнала | View options for commit log e.g. changing sort order, hiding the git graph, showing the whole git graph. |
| `` G `` | Open pull request in browser |  |
| `` <space> `` | Переключить | Checkout the selected commit as a detached HEAD. |
| `` y `` | Скопировать атрибут коммита | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | Открыть коммит в браузере |  |
| `` n `` | Создать новую ветку с этого коммита |  |
| `` N `` | Move commits to new branch | Create a new branch and move the unpushed commits of the current branch to it. Useful if you meant to start new work and forgot to create a new branch first.<br><br>Note that this disregards the selection, the new branch is always created either from the main branch or stacked on top of the current branch (you get to choose which). |
| `` g `` | Просмотреть параметры сброса | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` gm `` | Mixed reset | Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as unstaged changes. |
| `` gs `` | Мягкий сброс | Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as staged changes. |
| `` gh `` | Жёсткий сброс | Reset HEAD to the chosen commit, and discard all changes between the current and chosen commit, as well as all current modifications in the working tree. |
| `` C `` | Скопировать отобранные коммит (cherry-pick) | Mark commit as copied. Then, within the local commits view, you can press `V` to paste (cherry-pick) the copied commit(s) into your checked out branch. At any time you can press `<esc>` to cancel the selection. |
| `` <ctrl+t> `` | Open external diff tool (git difftool) |  |
| `` * `` | Select commits of current branch |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | Просмотреть файлы выбранного элемента |  |
| `` w `` | View worktree options |  |
| `` / `` | Найти |  |

## Локальные Ветки

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | Скопировать название ветки в буфер обмена |  |
| `` i `` | Показать параметры git-flow |  |
| `` iF `` | Finish git-flow branch |  |
| `` if `` | Start git-flow feature |  |
| `` ih `` | Start git-flow hotfix |  |
| `` ib `` | Start git-flow bugfix |  |
| `` ir `` | Start git-flow release |  |
| `` <space> `` | Переключить | Checkout selected item. |
| `` n `` | Новая ветка |  |
| `` N `` | Move commits to new branch | Create a new branch and move the unpushed commits of the current branch to it. Useful if you meant to start new work and forgot to create a new branch first.<br><br>Note that this disregards the selection, the new branch is always created either from the main branch or stacked on top of the current branch (you get to choose which). |
| `` o `` | Создать запрос на принятие изменений |  |
| `` O `` | Создать параметры запроса принятие изменений |  |
| `` G `` | Open pull request in browser |  |
| `` <ctrl+y> `` | Скопировать URL запроса на принятие изменений в буфер обмена |  |
| `` c `` | Переключить по названию | Checkout by name. In the input box you can enter '-' to switch to the previous branch. |
| `` - `` | Checkout previous branch |  |
| `` F `` | Принудительное переключение | Force checkout selected branch. This will discard all local changes in your working directory before checking out the selected branch. |
| `` d `` | Delete | View delete options for local/remote branch. |
| `` dc `` | Delete local branch | View delete options for local/remote branch. |
| `` dr `` | Удалить Удалённую Ветку | Delete the remote branch from the remote. |
| `` db `` | Delete local and remote branch | View delete options for local/remote branch. |
| `` r `` | Перебазировать переключённую ветку на эту ветку | Rebase the checked-out branch onto the selected branch. |
| `` rs `` | Simple rebase | Rebase the checked-out branch onto the selected branch. |
| `` ri `` | Interactive rebase | Начать интерактивную перебазировку с перерыва в начале, чтобы можно было обновить TODO коммиты, прежде чем продолжить. |
| `` rb `` | Rebase onto base branch | Rebase the checked out branch onto its base branch (i.e. the closest main branch). |
| `` M `` | Слияние с текущей переключённой веткой | View options for merging the selected item into the current branch (regular merge, squash merge) |
| `` Mm `` | Слияние с текущей переключённой веткой | View options for merging the selected item into the current branch (regular merge, squash merge) |
| `` Mn `` | Regular merge (with merge commit) | Merge '{{.selectedBranch}}' into '{{.checkedOutBranch}}', creating a merge commit. |
| `` Mf `` | Regular merge (fast-forward) | Fast-forward '{{.checkedOutBranch}}' to '{{.selectedBranch}}' without creating a merge commit. |
| `` Ms `` | Squash merge (uncommitted) | Squash merge '{{.selectedBranch}}' into the working tree. |
| `` MS `` | Squash merge (committed) | Squash merge '{{.selectedBranch}}' into '{{.checkedOutBranch}}' as a single commit. |
| `` f `` | Перемотать эту ветку вперёд из её upstream-ветки | Fast-forward selected branch from its upstream. |
| `` T `` | Создать тег |  |
| `` s `` | Порядок сортировки |  |
| `` g `` | Reset to ref |  |
| `` gm `` | Mixed reset | Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as unstaged changes. |
| `` gs `` | Мягкий сброс | Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as staged changes. |
| `` gh `` | Жёсткий сброс | Reset HEAD to the chosen commit, and discard all changes between the current and chosen commit, as well as all current modifications in the working tree. |
| `` R `` | Переименовать ветку |  |
| `` u `` | View upstream options | View options relating to the branch's upstream e.g. setting/unsetting the upstream and resetting to the upstream. |
| `` ud `` | View divergence from upstream | View options relating to the branch's upstream e.g. setting/unsetting the upstream and resetting to the upstream. |
| `` uD `` | View divergence from base branch |  |
| `` us `` | Установить upstream-ветку из выбранной ветки |  |
| `` uu `` | Убрать upstream-ветку из выбранной ветки |  |
| `` ugm `` | Mixed reset to upstream | Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as unstaged changes. |
| `` ugs `` | Soft reset to upstream | Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as staged changes. |
| `` ugh `` | Hard reset to upstream | Reset HEAD to the chosen commit, and discard all changes between the current and chosen commit, as well as all current modifications in the working tree. |
| `` urs `` | Simple rebase onto upstream |  |
| `` uri `` | Interactive rebase onto upstream | Начать интерактивную перебазировку с перерыва в начале, чтобы можно было обновить TODO коммиты, прежде чем продолжить. |
| `` urb `` | Rebase onto base branch | Rebase the checked out branch onto its base branch (i.e. the closest main branch). |
| `` <ctrl+t> `` | Open external diff tool (git difftool) |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | Просмотреть коммиты |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## Меню

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Выполнить |  |
| `` <esc> `` | Закрыть/отменить |  |
| `` / `` | Filter the current view by text |  |

## Панель Подтверждения

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Подтвердить |  |
| `` <esc> `` | Закрыть/отменить |  |
| `` <ctrl+o> `` | Copy to clipboard |  |

## Подкоммиты

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | Copy abbreviated commit hash to clipboard |  |
| `` <space> `` | Переключить | Checkout the selected commit as a detached HEAD. |
| `` y `` | Скопировать атрибут коммита | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | Открыть коммит в браузере |  |
| `` n `` | Создать новую ветку с этого коммита |  |
| `` N `` | Move commits to new branch | Create a new branch and move the unpushed commits of the current branch to it. Useful if you meant to start new work and forgot to create a new branch first.<br><br>Note that this disregards the selection, the new branch is always created either from the main branch or stacked on top of the current branch (you get to choose which). |
| `` g `` | Reset to ref |  |
| `` gm `` | Mixed reset | Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as unstaged changes. |
| `` gs `` | Мягкий сброс | Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as staged changes. |
| `` gh `` | Жёсткий сброс | Reset HEAD to the chosen commit, and discard all changes between the current and chosen commit, as well as all current modifications in the working tree. |
| `` C `` | Скопировать отобранные коммит (cherry-pick) | Mark commit as copied. Then, within the local commits view, you can press `V` to paste (cherry-pick) the copied commit(s) into your checked out branch. At any time you can press `<esc>` to cancel the selection. |
| `` <ctrl+r> `` | Сбросить отобранную (скопированную \| cherry-picked) выборку коммитов |  |
| `` <ctrl+t> `` | Open external diff tool (git difftool) |  |
| `` * `` | Select commits of current branch |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | Просмотреть файлы выбранного элемента |  |
| `` w `` | View worktree options |  |
| `` / `` | Найти |  |

## Подмодули

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | Скопировать название подмодуля в буфер обмена |  |
| `` <enter> `` | Enter | Ввести подмодуль |
| `` d `` | Remove | Remove the selected submodule and its corresponding directory. |
| `` u `` | Update | Обновить подмодуль |
| `` n `` | Добавить новый подмодуль |  |
| `` e `` | Обновить URL подмодуля |  |
| `` i `` | Initialize | Инициализировать подмодуль |
| `` b `` | Просмотреть параметры массового подмодуля |  |
| `` bi `` | Массовая инициализация подмодулей |  |
| `` bu `` | Массовое обновление подмодулей |  |
| `` br `` | Bulk init and update submodules recursively |  |
| `` bd `` | Массовая деинициализация подмодулей |  |
| `` / `` | Filter the current view by text |  |

## Сводка коммита

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Подтвердить |  |
| `` <esc> `` | Закрыть |  |

## Сохранить Изменения Файлов

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | Скопировать название файла в буфер обмена |  |
| `` y `` | Copy to clipboard |  |
| `` yn `` | File name |  |
| `` yp `` | Relative path |  |
| `` yP `` | Absolute path |  |
| `` ys `` | Diff of selected file |  |
| `` ya `` | Diff of all files |  |
| `` yc `` | Content of selected file |  |
| `` c `` | Переключить | Переключить файл |
| `` d `` | Просмотреть параметры «отмены изменении» | Отменить изменения коммита в этом файле |
| `` o `` | Открыть файл | Open file in default application. |
| `` e `` | Edit | Open file in external editor. |
| `` <ctrl+t> `` | Open external diff tool (git difftool) |  |
| `` <space> `` | Переключить файлы включённые в патч | Toggle whether the file is included in the custom patch. See https://github.com/jesseduffield/lazygit#rebase-magic-custom-patches. |
| `` a `` | Переключить все файлы, включённые в патч | Add/remove all commit's files to custom patch. See https://github.com/jesseduffield/lazygit#rebase-magic-custom-patches. |
| `` <enter> `` | Введите файл, чтобы добавить выбранные строки в патч (или свернуть каталог переключения) | If a file is selected, enter the file so that you can add/remove individual lines to the custom patch. If a directory is selected, toggle the directory. |
| `` ` `` | Переключить вид дерева файлов | Toggle file view between flat and tree layout. Flat layout shows all file paths in a single list, tree layout groups files by directory.<br><br>The default can be changed in the config file with the key 'gui.showFileTree'. |
| `` - `` | Collapse all files | Collapse all directories in the files tree |
| `` = `` | Expand all files | Expand all directories in the file tree |
| `` 0 `` | Focus main view |  |
| `` / `` | Filter the current view by text |  |

## Статус

| Key | Action | Info |
|-----|--------|-------------|
| `` o `` | Открыть файл конфигурации | Open file in default application. |
| `` e `` | Редактировать файл конфигурации | Open file in external editor. |
| `` u `` | Проверить обновления |  |
| `` <enter> `` | Переключиться на последний репозиторий |  |
| `` a `` | Show/cycle all branch logs |  |
| `` A `` | Show/cycle all branch logs (reverse) |  |
| `` 0 `` | Focus main view |  |

## Теги

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | Copy tag to clipboard |  |
| `` <space> `` | Переключить | Checkout the selected tag as a detached HEAD. |
| `` n `` | Создать тег | Create new tag from current commit. You'll be prompted to enter a tag name and optional description. |
| `` d `` | Delete |  |
| `` dc `` | Delete local tag | View delete options for local/remote tag. |
| `` dr `` | Delete remote tag | View delete options for local/remote tag. |
| `` db `` | Delete local and remote tag | View delete options for local/remote tag. |
| `` P `` | Отправить тег | Push the selected tag to a remote. You'll be prompted to select a remote. |
| `` g `` | Reset to ref |  |
| `` gm `` | Mixed reset | Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as unstaged changes. |
| `` gs `` | Мягкий сброс | Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as staged changes. |
| `` gh `` | Жёсткий сброс | Reset HEAD to the chosen commit, and discard all changes between the current and chosen commit, as well as all current modifications in the working tree. |
| `` <ctrl+t> `` | Open external diff tool (git difftool) |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | Просмотреть коммиты |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## Удалённые ветки

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | Скопировать название ветки в буфер обмена |  |
| `` <space> `` | Переключить | Checkout a new local branch based on the selected remote branch, or the remote branch as a detached head. |
| `` n `` | Новая ветка |  |
| `` M `` | Merge |  |
| `` Mm `` | Слияние с текущей переключённой веткой | View options for merging the selected item into the current branch (regular merge, squash merge) |
| `` Mn `` | Non-fast-forward merge | Merge '{{.selectedBranch}}' into '{{.checkedOutBranch}}', creating a merge commit. |
| `` Mf `` | Fast-forward only merge | Fast-forward '{{.checkedOutBranch}}' to '{{.selectedBranch}}' without creating a merge commit. |
| `` Ms `` | Squash merge (uncommitted) | Squash merge '{{.selectedBranch}}' into the working tree. |
| `` MS `` | Squash merge (committed) | Squash merge '{{.selectedBranch}}' into '{{.checkedOutBranch}}' as a single commit. |
| `` r `` | Rebase options |  |
| `` rs `` | Перебазировать переключённую ветку на эту ветку | Rebase the checked-out branch onto the selected branch. |
| `` ri `` | Interactive rebase | Начать интерактивную перебазировку с перерыва в начале, чтобы можно было обновить TODO коммиты, прежде чем продолжить. |
| `` rb `` | Rebase onto base branch | Rebase the checked out branch onto its base branch (i.e. the closest main branch). |
| `` d `` | Delete | Delete the remote branch from the remote. |
| `` us `` | Set as upstream | Установить как upstream-ветку переключённую ветку |
| `` s `` | Порядок сортировки |  |
| `` g `` | Reset to ref |  |
| `` gm `` | Mixed reset | Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as unstaged changes. |
| `` gs `` | Мягкий сброс | Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as staged changes. |
| `` gh `` | Жёсткий сброс | Reset HEAD to the chosen commit, and discard all changes between the current and chosen commit, as well as all current modifications in the working tree. |
| `` <ctrl+t> `` | Open external diff tool (git difftool) |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | Просмотреть коммиты |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## Удалённые репозитории

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | View branches |  |
| `` n `` | Добавить новую удалённую ветку |  |
| `` d `` | Remove | Remove the selected remote. Any local branches tracking a remote branch from the remote will be unaffected. |
| `` e `` | Edit | Редактировать удалённый репозитории |
| `` f `` | Получить изменения | Получение изменения из удалённого репозитория |
| `` F `` | Add fork remote | Quickly add a fork remote by replacing the owner in the origin URL and optionally check out a branch from new remote. |
| `` / `` | Filter the current view by text |  |

## Файлы

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | Скопировать название файла в буфер обмена |  |
| `` <space> `` | Переключить индекс | Toggle staged for selected file. |
| `` <ctrl+b> `` | Фильтровать файлы (проиндексированные/непроиндексированные) |  |
| `` <ctrl+b>s `` | Показывать только проиндексированные файлы |  |
| `` <ctrl+b>u `` | Показывать только непроиндексированные файлы |  |
| `` <ctrl+b>t `` | Show only tracked files |  |
| `` <ctrl+b>T `` | Show only untracked files |  |
| `` <ctrl+b>r `` | No filter |  |
| `` y `` | Copy to clipboard |  |
| `` yn `` | File name |  |
| `` yp `` | Relative path |  |
| `` yP `` | Absolute path |  |
| `` ys `` | Diff of selected file | If there are staged items, this command considers only them. Otherwise, it considers all the unstaged ones. |
| `` ya `` | Diff of all files | If there are staged items, this command considers only them. Otherwise, it considers all the unstaged ones. |
| `` c `` | Сохранить изменения | Commit staged changes. |
| `` w `` | Закоммитить изменения без предварительного хука коммита |  |
| `` A `` | Правка последнего коммита |  |
| `` C `` | Сохранить изменения с помощью редактора git |  |
| `` <ctrl+f> `` | Find base commit for fixup | Find the commit that your current changes are building upon, for the sake of amending/fixing up the commit. This spares you from having to look through your branch's commits one-by-one to see which commit should be amended/fixed up. See docs: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` e `` | Edit | Open file in external editor. |
| `` o `` | Открыть файл | Open file in default application. |
| `` i `` | Игнорировать или исключить файл |  |
| `` ii `` | Добавить в .gitignore |  |
| `` ie `` | Добавить в .git/info/exclude |  |
| `` r `` | Обновить файлы |  |
| `` s `` | Stash | Stash all changes. Press capital S for variations (keep index, include untracked, staged only, unstaged only). |
| `` S `` | Просмотреть параметры хранилища | View stash options (e.g. stash all, stash staged, stash unstaged). |
| `` Si `` | Припрятать все изменения и сохранить индекс |  |
| `` SU `` | Припрятать все изменения, включая неотслеживаемые файлы |  |
| `` Ss `` | Припрятать проиндексированные изменения |  |
| `` Su `` | Припрятать непроиндексированные изменения |  |
| `` a `` | Все проиндексированные/непроиндексированные | Toggle staged/unstaged for all files in working tree. |
| `` <enter> `` | Проиндексировать отдельные части/строки для файла или свернуть/развернуть для каталога | If the selected item is a file, focus the staging view so you can stage individual hunks/lines. If the selected item is a directory, collapse/expand it. |
| `` d `` | Discard changes |  |
| `` dc `` | Просмотреть параметры «отмены изменении» | View options for discarding changes to the selected file. |
| `` du `` | Отменить непроиндексированные изменения |  |
| `` g `` | Reset to upstream |  |
| `` gm `` | Mixed reset to upstream | Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as unstaged changes. |
| `` gs `` | Soft reset to upstream | Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as staged changes. |
| `` gh `` | Hard reset to upstream | Reset HEAD to the chosen commit, and discard all changes between the current and chosen commit, as well as all current modifications in the working tree. |
| `` D `` | Reset | View reset options for working tree (e.g. nuking the working tree). |
| `` Dx `` | Разбомбить рабочее дерево? | Если вы хотите, чтобы все изменения в рабочем дереве исчезли, это способ сделать это. Если есть какие-либо изменения подмодуля, эти изменения будут припрятаны в подмодуле(-ях). |
| `` Du `` | Отменить непроиндексированные изменения |  |
| `` Dc `` | Удалить неотслеживаемые файлы |  |
| `` DS `` | Отменить проиндексированные изменения | Это создаст новую запись в хранилище, содержащую только проиндексированные файлы, а затем удалит её, так что в рабочем дереве останутся только непроиндексированные изменения. |
| `` Ds `` | Мягкий сброс |  |
| `` Dm `` | mixed reset |  |
| `` Dh `` | Жёсткий сброс |  |
| `` ` `` | Переключить вид дерева файлов | Toggle file view between flat and tree layout. Flat layout shows all file paths in a single list, tree layout groups files by directory.<br><br>The default can be changed in the config file with the key 'gui.showFileTree'. |
| `` <ctrl+t> `` | Open external diff tool (git difftool) |  |
| `` M `` | View merge conflict options | View options for resolving merge conflicts. |
| `` f `` | Получить изменения | Fetch changes from remote. |
| `` - `` | Collapse all files | Collapse all directories in the files tree |
| `` = `` | Expand all files | Expand all directories in the file tree |
| `` 0 `` | Focus main view |  |
| `` / `` | Filter the current view by text |  |

## Хранилище

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | Применить припрятанные изменения | Apply the stash entry to your working directory. |
| `` g `` | Применить припрятанные изменения и тут же удалить их из хранилища | Apply the stash entry to your working directory and remove the stash entry. |
| `` d `` | Удалить припрятанные изменения из хранилища | Remove the stash entry from the stash list. |
| `` n `` | Новая ветка | Create a new branch from the selected stash entry. This works by git checking out the commit that the stash entry was created from, creating a new branch from that commit, then applying the stash entry to the new branch as an additional commit. |
| `` r `` | Переименовать хранилище |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | Просмотреть файлы выбранного элемента |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |
