_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go generate ./...` from the project root._

# Lazygit 鍵盤快捷鍵

## 全域快捷鍵

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+r> `` | 切換到最近使用的版本庫 |  |
| `` <pgup> (fn+up/shift+k) `` | 向上捲動主面板 |  |
| `` <pgdown> (fn+down/shift+j) `` | 向下捲動主面板 |  |
| `` @ `` | 開啟命令記錄選單 | View options for the command log e.g. show/hide the command log and focus the command log. |
| `` P `` | 推送 | 推送到遠端。如果沒有設定遠端，會開啟設定視窗。 |
| `` p `` | 拉取 | 從遠端同步當前分支。如果沒有設定遠端，會開啟設定視窗。 |
| `` ) `` | Increase rename similarity threshold | Increase the similarity threshold for a deletion and addition pair to be treated as a rename.<br><br>The default can be changed in the config file with the key 'git.renameSimilarityThreshold'. |
| `` ( `` | Decrease rename similarity threshold | Decrease the similarity threshold for a deletion and addition pair to be treated as a rename.<br><br>The default can be changed in the config file with the key 'git.renameSimilarityThreshold'. |
| `` } `` | 增加差異檢視中顯示變更周圍上下文的大小 | Increase the amount of the context shown around changes in the diff view.<br><br>The default can be changed in the config file with the key 'git.diffContextSize'. |
| `` { `` | 減小差異檢視中顯示變更周圍上下文的大小 | Decrease the amount of the context shown around changes in the diff view.<br><br>The default can be changed in the config file with the key 'git.diffContextSize'. |
| `` : `` | Execute shell command | Bring up a prompt where you can enter a shell command to execute. |
| `` <ctrl+p> `` | 檢視自訂補丁選項 |  |
| `` m `` | 查看合併/變基選項 | View options to abort/continue/skip the current merge/rebase. |
| `` mc `` | Continue rebase / merge | View options to abort/continue/skip the current merge/rebase. |
| `` ma `` | Abort rebase / merge | View options to abort/continue/skip the current merge/rebase. |
| `` ms `` | Skip current rebase commit | View options to abort/continue/skip the current merge/rebase. |
| `` R `` | 重新整理 | Refresh the git state (i.e. run `git status`, `git branch`, etc in background to update the contents of panels). This does not run `git fetch`. |
| `` + `` | 下一個螢幕模式（常規/半螢幕/全螢幕） |  |
| `` _ `` | 上一個螢幕模式 |  |
| `` \| `` | Cycle pagers | Choose the next pager in the list of configured pagers |
| `` <esc> `` | 取消 |  |
| `` ? `` | 開啟選單 |  |
| `` <ctrl+s> `` | 檢視篩選路徑選項 | View options for filtering the commit log, so that only commits matching the filter are shown. |
| `` W `` | 開啟差異比較選單 | View options relating to diffing two refs e.g. diffing against selected ref, entering ref to diff against, and reversing the diff direction. |
| `` <ctrl+e> `` | 開啟差異比較選單 | View options relating to diffing two refs e.g. diffing against selected ref, entering ref to diff against, and reversing the diff direction. |
| `` q `` | 結束 |  |
| `` <ctrl+z> `` | Suspend the application |  |
| `` <ctrl+w> `` | 切換是否在差異檢視中顯示空格變更 | Toggle whether or not whitespace changes are shown in the diff view.<br><br>The default can be changed in the config file with the key 'git.ignoreWhitespaceInDiffView'. |
| `` z `` | 復原 | 將使用 reflog 確任 git 指令以復原。這不包括工作區更改；只考慮提交。 |
| `` Z `` | 取消復原 | 將使用 reflog 確任 git 指令以重作。這不包括工作區更改；只考慮提交。 |

## 移動

| Key | Action | Info |
|-----|--------|-------------|
| `` , `` | 上一頁 |  |
| `` . `` | 下一頁 |  |
| `` < (<home>) `` | 捲動到頂部 |  |
| `` > (<end>) `` | 捲動到底部 |  |
| `` v `` | 切換拖曳選擇 |  |
| `` <shift+down> `` | Range select down |  |
| `` <shift+up> `` | Range select up |  |
| `` / `` | 搜尋 |  |
| `` H `` | 向左捲動 |  |
| `` L `` | 向右捲動 |  |
| `` ] `` | 下一個索引標籤 |  |
| `` [ `` | 上一個索引標籤 |  |

## Input prompt

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | 確認 |  |
| `` <esc> `` | 關閉/取消 |  |

## 主面板 (補丁生成)

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | 選擇上一段 |  |
| `` <right> `` | 選擇下一段 |  |
| `` v `` | 切換拖曳選擇 |  |
| `` a `` | Toggle hunk selection | Toggle line-by-line vs. hunk selection mode. |
| `` <ctrl+o> `` | 複製所選文本至剪貼簿 |  |
| `` o `` | 開啟檔案 | 使用預設軟體開啟 |
| `` e `` | 編輯檔案 | 使用外部編輯器開啟 |
| `` <space> `` | 向 (或從) 補丁中添加/刪除行 |  |
| `` d `` | Remove lines from commit | Remove the selected lines from this commit. This runs an interactive rebase in the background, so you may get a merge conflict if a later commit also changes these lines. |
| `` <esc> `` | 退出自訂補丁建立器 |  |
| `` / `` | 搜尋 |  |

## 主面板（一般）

| Key | Action | Info |
|-----|--------|-------------|
| `` <mouse wheel down> (fn+up) `` | 向下捲動 |  |
| `` <mouse wheel up> (fn+down) `` | 向上捲動 |  |
| `` <tab> `` | 切換至另一個面板 (已預存/未預存更改) | Switch to other view (staged/unstaged changes). |
| `` <esc> `` | Exit back to side panel |  |
| `` / `` | 搜尋 |  |

## 主面板（合併）

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | 挑選程式碼片段 |  |
| `` b `` | 挑選所有程式碼片段 |  |
| `` <up> `` | 選擇上一段 |  |
| `` <down> `` | 選擇下一段 |  |
| `` <left> `` | 選擇上一個衝突 |  |
| `` <right> `` | 選擇下一個衝突 |  |
| `` z `` | 復原 | Undo last merge conflict resolution. |
| `` e `` | 編輯檔案 | 使用外部編輯器開啟 |
| `` o `` | 開啟檔案 | 使用預設軟體開啟 |
| `` M `` | View merge conflict options | View options for resolving merge conflicts. |
| `` <esc> `` | 返回檔案面板 |  |

## 主面板（預存）

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | 選擇上一段 |  |
| `` <right> `` | 選擇下一段 |  |
| `` v `` | 切換拖曳選擇 |  |
| `` a `` | Toggle hunk selection | Toggle line-by-line vs. hunk selection mode. |
| `` <ctrl+o> `` | 複製所選文本至剪貼簿 |  |
| `` <space> `` | 切換預存 | 切換現有行的狀態 (已預存/未預存) |
| `` d `` | 刪除變更 (git reset) | When unstaged change is selected, discard the change using `git reset`. When staged change is selected, unstage the change. |
| `` o `` | 開啟檔案 | 使用預設軟體開啟 |
| `` e `` | 編輯檔案 | 使用外部編輯器開啟 |
| `` <esc> `` | 返回檔案面板 |  |
| `` <tab> `` | 切換至另一個面板 (已預存/未預存更改) | Switch to other view (staged/unstaged changes). |
| `` E `` | 編輯程式碼塊 | Edit selected hunk in external editor. |
| `` c `` | 提交變更 | 提交暫存區變更 |
| `` w `` | 沒有預提交 hook 就提交更改 |  |
| `` C `` | 使用 git 編輯器提交變更 |  |
| `` <ctrl+f> `` | Find base commit for fixup | Find the commit that your current changes are building upon, for the sake of amending/fixing up the commit. This spares you from having to look through your branch's commits one-by-one to see which commit should be amended/fixed up. See docs: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` / `` | 搜尋 |  |

## 功能表

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | 執行 |  |
| `` <esc> `` | 關閉/取消 |  |
| `` / `` | 搜尋 |  |

## 子提交

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | Copy abbreviated commit hash to clipboard |  |
| `` <space> `` | 檢出 | Checkout the selected commit as a detached HEAD. |
| `` y `` | 複製提交屬性 | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | 在瀏覽器中開啟提交 |  |
| `` n `` | 從提交建立新分支 |  |
| `` N `` | Move commits to new branch | Create a new branch and move the unpushed commits of the current branch to it. Useful if you meant to start new work and forgot to create a new branch first.<br><br>Note that this disregards the selection, the new branch is always created either from the main branch or stacked on top of the current branch (you get to choose which). |
| `` g `` | Reset to ref |  |
| `` gm `` | Mixed reset | Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as unstaged changes. |
| `` gs `` | 軟重設 | Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as staged changes. |
| `` gh `` | 強制重設 | Reset HEAD to the chosen commit, and discard all changes between the current and chosen commit, as well as all current modifications in the working tree. |
| `` C `` | 複製提交 (揀選) | Mark commit as copied. Then, within the local commits view, you can press `V` to paste (cherry-pick) the copied commit(s) into your checked out branch. At any time you can press `<esc>` to cancel the selection. |
| `` <ctrl+r> `` | 重設選定的揀選 (複製) 提交 |  |
| `` <ctrl+t> `` | 開啟外部差異工具 (git difftool) |  |
| `` * `` | Select commits of current branch |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | 檢視所選項目的檔案 |  |
| `` w `` | 檢視工作目錄選項 |  |
| `` / `` | 搜尋 |  |

## 子模組

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | 複製子模組名稱到剪貼簿 |  |
| `` <enter> `` | Enter | 進入子模組 |
| `` d `` | Remove | Remove the selected submodule and its corresponding directory. |
| `` u `` | Update | 更新子模組 |
| `` n `` | 新增子模組 |  |
| `` e `` | 更新子模組 URL |  |
| `` i `` | Initialize | 初始化子模組 |
| `` b `` | 查看批量子模組選項 |  |
| `` bi `` | 批量初始化子模組 |  |
| `` bu `` | 批量更新子模組 |  |
| `` br `` | Bulk init and update submodules recursively |  |
| `` bd `` | 批量解除子模組初始化 |  |
| `` / `` | 搜尋 |  |

## 工作目錄

| Key | Action | Info |
|-----|--------|-------------|
| `` n `` | New worktree |  |
| `` <space> `` | Switch | Switch to the selected worktree. |
| `` o `` | 在編輯器中開啟 |  |
| `` d `` | Remove | Remove the selected worktree. This will both delete the worktree's directory, as well as metadata about the worktree in the .git directory. |
| `` / `` | 搜尋 |  |

## 提交

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | Copy abbreviated commit hash to clipboard |  |
| `` <ctrl+r> `` | 重設選定的揀選 (複製) 提交 |  |
| `` b `` | 查看二分選項 |  |
| `` bb `` | Mark commit as bad |  |
| `` bg `` | Mark commit as good |  |
| `` bs `` | Skip current bisect commit |  |
| `` bS `` | Skip selected commit |  |
| `` br `` | 重設二分查找 |  |
| `` bb `` | Mark commit as bad (start bisect) |  |
| `` bg `` | Mark commit as good (start bisect) |  |
| `` bt `` | Choose bisect terms |  |
| `` s `` | 壓縮 (Squash) | Squash the selected commit into the commit below it. The selected commit's message will be appended to the commit below it. |
| `` f `` | 修復 (Fixup) | Meld the selected commit into the commit below it. Similar to squash, but the selected commit's message will be discarded. |
| `` ff `` | 修復 (Fixup) | Meld the selected commit into the commit below it. Similar to squash, but the selected commit's message will be discarded. |
| `` fc `` | Fixup and use this commit's message | Squash the selected commit into the commit below, using this commit's message, discarding the message of the commit below. |
| `` c `` | Set fixup message | Set the message option for the fixup commit. The -C option means to use this commit's message instead of the target commit's message. |
| `` r `` | 改寫提交 | 改寫選中的提交訊息 |
| `` R `` | 使用編輯器改寫提交 |  |
| `` d `` | 刪除提交 | Drop the selected commit. This will remove the commit from the branch via a rebase. If the commit makes changes that later commits depend on, you may need to resolve merge conflicts. |
| `` e `` | 編輯(開始互動變基) | 編輯提交 |
| `` i `` | 開始互動變基 | Start an interactive rebase for the commits on your branch. This will include all commits from the HEAD commit down to the first merge commit or main branch commit.<br>If you would instead like to start an interactive rebase from the selected commit, press `e`. |
| `` p `` | 挑選 | 挑選提交 (於變基過程中) |
| `` F `` | 建立修復提交 | 為此提交建立修復提交 |
| `` S `` | 壓縮上方所有「fixup」提交（自動壓縮） | 是否壓縮上方 {{.commit}} 所有「fixup」提交？ |
| `` <alt+down> `` | 向下移動提交 |  |
| `` <alt+up> `` | 向上移動提交 |  |
| `` V `` | 貼上提交 (揀選) |  |
| `` B `` | 為了變基已標注提交為基準提交 | 請為了下一次變基選擇一項基準提交；此將執行 `git rebase --onto`。 |
| `` A `` | 修改 | 使用已預存的更改修正提交 |
| `` a `` | 設定/重設提交作者 | Set/Reset commit author or set co-author. |
| `` t `` | 還原 | Create a revert commit for the selected commit, which applies the selected commit's changes in reverse. |
| `` T `` | 打標籤到提交 | Create a new tag pointing at the selected commit. You'll be prompted to enter a tag name and optional description. |
| `` <ctrl+l> `` | 開啟記錄選單 | View options for commit log e.g. changing sort order, hiding the git graph, showing the whole git graph. |
| `` G `` | Open pull request in browser |  |
| `` <space> `` | 檢出 | Checkout the selected commit as a detached HEAD. |
| `` y `` | 複製提交屬性 | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | 在瀏覽器中開啟提交 |  |
| `` n `` | 從提交建立新分支 |  |
| `` N `` | Move commits to new branch | Create a new branch and move the unpushed commits of the current branch to it. Useful if you meant to start new work and forgot to create a new branch first.<br><br>Note that this disregards the selection, the new branch is always created either from the main branch or stacked on top of the current branch (you get to choose which). |
| `` g `` | 檢視重設選項 | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` gm `` | Mixed reset | Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as unstaged changes. |
| `` gs `` | 軟重設 | Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as staged changes. |
| `` gh `` | 強制重設 | Reset HEAD to the chosen commit, and discard all changes between the current and chosen commit, as well as all current modifications in the working tree. |
| `` C `` | 複製提交 (揀選) | Mark commit as copied. Then, within the local commits view, you can press `V` to paste (cherry-pick) the copied commit(s) into your checked out branch. At any time you can press `<esc>` to cancel the selection. |
| `` <ctrl+t> `` | 開啟外部差異工具 (git difftool) |  |
| `` * `` | Select commits of current branch |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | 檢視所選項目的檔案 |  |
| `` w `` | 檢視工作目錄選項 |  |
| `` / `` | 搜尋 |  |

## 提交摘要

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | 確認 |  |
| `` <esc> `` | 關閉 |  |

## 提交檔案

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | 複製檔案名稱到剪貼簿 |  |
| `` y `` | Copy to clipboard |  |
| `` yn `` | 檔案名稱 |  |
| `` yp `` | Relative path |  |
| `` yP `` | Absolute path |  |
| `` ys `` | 所選檔案的差異 |  |
| `` ya `` | 所有檔案的差異 |  |
| `` yc `` | Content of selected file |  |
| `` c `` | 檢出 | 檢出檔案 |
| `` d `` | 捨棄 | Discard this commit's changes to this file. This runs an interactive rebase in the background, so you may get a merge conflict if a later commit also changes this file. |
| `` o `` | 開啟檔案 | 使用預設軟體開啟 |
| `` e `` | 編輯 | 使用外部編輯器開啟 |
| `` <ctrl+t> `` | 開啟外部差異工具 (git difftool) |  |
| `` <space> `` | 切換檔案是否包含在補丁中 | Toggle whether the file is included in the custom patch. See https://github.com/jesseduffield/lazygit#rebase-magic-custom-patches. |
| `` a `` | 切換所有檔案是否包含在補丁中 | Add/remove all commit's files to custom patch. See https://github.com/jesseduffield/lazygit#rebase-magic-custom-patches. |
| `` <enter> `` | 輸入檔案以將選定的行添加至補丁（或切換目錄折疊） | If a file is selected, enter the file so that you can add/remove individual lines to the custom patch. If a directory is selected, toggle the directory. |
| `` ` `` | 顯示檔案樹狀視圖 | Toggle file view between flat and tree layout. Flat layout shows all file paths in a single list, tree layout groups files by directory.<br><br>The default can be changed in the config file with the key 'gui.showFileTree'. |
| `` - `` | Collapse all files | Collapse all directories in the files tree |
| `` = `` | Expand all files | Expand all directories in the file tree |
| `` 0 `` | Focus main view |  |
| `` / `` | 搜尋 |  |

## 收藏 (Stash)

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | 套用 | Apply the stash entry to your working directory. |
| `` g `` | 還原 | Apply the stash entry to your working directory and remove the stash entry. |
| `` d `` | 捨棄 | Remove the stash entry from the stash list. |
| `` n `` | 新分支 | Create a new branch from the selected stash entry. This works by git checking out the commit that the stash entry was created from, creating a new branch from that commit, then applying the stash entry to the new branch as an additional commit. |
| `` r `` | 重新命名收藏 |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | 檢視所選項目的檔案 |  |
| `` w `` | 檢視工作目錄選項 |  |
| `` / `` | 搜尋 |  |

## 日誌

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | Copy abbreviated commit hash to clipboard |  |
| `` <space> `` | 檢出 | Checkout the selected commit as a detached HEAD. |
| `` y `` | 複製提交屬性 | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | 在瀏覽器中開啟提交 |  |
| `` n `` | 從提交建立新分支 |  |
| `` N `` | Move commits to new branch | Create a new branch and move the unpushed commits of the current branch to it. Useful if you meant to start new work and forgot to create a new branch first.<br><br>Note that this disregards the selection, the new branch is always created either from the main branch or stacked on top of the current branch (you get to choose which). |
| `` g `` | Reset to ref |  |
| `` gm `` | Mixed reset | Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as unstaged changes. |
| `` gs `` | 軟重設 | Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as staged changes. |
| `` gh `` | 強制重設 | Reset HEAD to the chosen commit, and discard all changes between the current and chosen commit, as well as all current modifications in the working tree. |
| `` C `` | 複製提交 (揀選) | Mark commit as copied. Then, within the local commits view, you can press `V` to paste (cherry-pick) the copied commit(s) into your checked out branch. At any time you can press `<esc>` to cancel the selection. |
| `` <ctrl+r> `` | 重設選定的揀選 (複製) 提交 |  |
| `` <ctrl+t> `` | 開啟外部差異工具 (git difftool) |  |
| `` * `` | Select commits of current branch |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | 檢視提交 |  |
| `` w `` | 檢視工作目錄選項 |  |
| `` / `` | 搜尋 |  |

## 本地分支

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | 複製分支名稱到剪貼簿 |  |
| `` i `` | 顯示 git-flow 選項 |  |
| `` iF `` | Finish git-flow branch |  |
| `` if `` | Start git-flow feature |  |
| `` ih `` | Start git-flow hotfix |  |
| `` ib `` | Start git-flow bugfix |  |
| `` ir `` | Start git-flow release |  |
| `` <space> `` | 檢出 | 檢出選定的項目。 |
| `` n `` | 新分支 |  |
| `` N `` | Move commits to new branch | Create a new branch and move the unpushed commits of the current branch to it. Useful if you meant to start new work and forgot to create a new branch first.<br><br>Note that this disregards the selection, the new branch is always created either from the main branch or stacked on top of the current branch (you get to choose which). |
| `` o `` | 建立拉取請求 |  |
| `` O `` | 建立拉取請求選項 |  |
| `` G `` | Open pull request in browser |  |
| `` <ctrl+y> `` | 複製拉取請求的 URL 到剪貼板 |  |
| `` c `` | 根據名稱檢出 | Checkout by name. In the input box you can enter '-' to switch to the previous branch. |
| `` - `` | Checkout previous branch |  |
| `` F `` | 強制檢出 | Force checkout selected branch. This will discard all local changes in your working directory before checking out the selected branch. |
| `` d `` | 刪除 | View delete options for local/remote branch. |
| `` dc `` | Delete local branch | View delete options for local/remote branch. |
| `` dr `` | 刪除遠端分支 | Delete the remote branch from the remote. |
| `` db `` | Delete local and remote branch | View delete options for local/remote branch. |
| `` r `` | 將已檢出的分支變基至此分支 | Rebase the checked-out branch onto the selected branch. |
| `` rs `` | Simple rebase | Rebase the checked-out branch onto the selected branch. |
| `` ri `` | Interactive rebase | 開始一個互動變基，以中斷開始，這樣你可以在繼續之前更新TODO提交 |
| `` rb `` | Rebase onto base branch | Rebase the checked out branch onto its base branch (i.e. the closest main branch). |
| `` M `` | 合併到當前檢出的分支 | View options for merging the selected item into the current branch (regular merge, squash merge) |
| `` Mm `` | 合併到當前檢出的分支 | View options for merging the selected item into the current branch (regular merge, squash merge) |
| `` Mn `` | Regular merge (with merge commit) | Merge '{{.selectedBranch}}' into '{{.checkedOutBranch}}', creating a merge commit. |
| `` Mf `` | Regular merge (fast-forward) | Fast-forward '{{.checkedOutBranch}}' to '{{.selectedBranch}}' without creating a merge commit. |
| `` Ms `` | Squash merge (uncommitted) | Squash merge '{{.selectedBranch}}' into the working tree. |
| `` MS `` | Squash merge (committed) | Squash merge '{{.selectedBranch}}' into '{{.checkedOutBranch}}' as a single commit. |
| `` f `` | 從上游快進此分支 | 從遠端快進所選的分支 |
| `` T `` | 建立標籤 |  |
| `` s `` | 排序規則 |  |
| `` g `` | Reset to ref |  |
| `` gm `` | Mixed reset | Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as unstaged changes. |
| `` gs `` | 軟重設 | Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as staged changes. |
| `` gh `` | 強制重設 | Reset HEAD to the chosen commit, and discard all changes between the current and chosen commit, as well as all current modifications in the working tree. |
| `` R `` | 重新命名分支 |  |
| `` u `` | 檢視遠端設定 | 檢視有關遠端分支的設定（例如重設至遠端） |
| `` ud `` | 檢視與遠端的差異 | 檢視有關遠端分支的設定（例如重設至遠端） |
| `` uD `` | View divergence from base branch |  |
| `` us `` | 設定選定分支的遠端分支 |  |
| `` uu `` | 重置選定分支的遠端 |  |
| `` ugm `` | Mixed reset to upstream | Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as unstaged changes. |
| `` ugs `` | Soft reset to upstream | Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as staged changes. |
| `` ugh `` | Hard reset to upstream | Reset HEAD to the chosen commit, and discard all changes between the current and chosen commit, as well as all current modifications in the working tree. |
| `` urs `` | Simple rebase onto upstream |  |
| `` uri `` | Interactive rebase onto upstream | 開始一個互動變基，以中斷開始，這樣你可以在繼續之前更新TODO提交 |
| `` urb `` | Rebase onto base branch | Rebase the checked out branch onto its base branch (i.e. the closest main branch). |
| `` <ctrl+t> `` | 開啟外部差異工具 (git difftool) |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | 檢視提交 |  |
| `` w `` | 檢視工作目錄選項 |  |
| `` / `` | 搜尋 |  |

## 標籤

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | Copy tag to clipboard |  |
| `` <space> `` | 檢出 | Checkout the selected tag as a detached HEAD. |
| `` n `` | 建立標籤 | Create new tag from current commit. You'll be prompted to enter a tag name and optional description. |
| `` d `` | 刪除 |  |
| `` dc `` | 刪除本地標籤 | View delete options for local/remote tag. |
| `` dr `` | Delete remote tag | View delete options for local/remote tag. |
| `` db `` | Delete local and remote tag | View delete options for local/remote tag. |
| `` P `` | 推送標籤 | Push the selected tag to a remote. You'll be prompted to select a remote. |
| `` g `` | Reset to ref |  |
| `` gm `` | Mixed reset | Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as unstaged changes. |
| `` gs `` | 軟重設 | Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as staged changes. |
| `` gh `` | 強制重設 | Reset HEAD to the chosen commit, and discard all changes between the current and chosen commit, as well as all current modifications in the working tree. |
| `` <ctrl+t> `` | 開啟外部差異工具 (git difftool) |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | 檢視提交 |  |
| `` w `` | 檢視工作目錄選項 |  |
| `` / `` | 搜尋 |  |

## 檔案

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | 複製檔案名稱到剪貼簿 |  |
| `` <space> `` | 切換預存 | Toggle staged for selected file. |
| `` <ctrl+b> `` | 篩選檔案 (預存/未預存) |  |
| `` <ctrl+b>s `` | 僅顯示預存的檔案 |  |
| `` <ctrl+b>u `` | 僅顯示未預存的檔案 |  |
| `` <ctrl+b>t `` | Show only tracked files |  |
| `` <ctrl+b>T `` | Show only untracked files |  |
| `` <ctrl+b>r `` | No filter |  |
| `` y `` | 複製到剪貼簿 |  |
| `` yn `` | 檔案名稱 |  |
| `` yp `` | Relative path |  |
| `` yP `` | Absolute path |  |
| `` ys `` | 所選檔案的差異 | 如果有已預存的項目，此指令只考慮它們。否則，它將考慮所有未暫存的項目。 |
| `` ya `` | 所有檔案的差異 | 如果有已預存的項目，此指令只考慮它們。否則，它將考慮所有未暫存的項目。 |
| `` c `` | 提交變更 | 提交暫存區變更 |
| `` w `` | 沒有預提交 hook 就提交更改 |  |
| `` A `` | 修改上次提交 |  |
| `` C `` | 使用 git 編輯器提交變更 |  |
| `` <ctrl+f> `` | Find base commit for fixup | Find the commit that your current changes are building upon, for the sake of amending/fixing up the commit. This spares you from having to look through your branch's commits one-by-one to see which commit should be amended/fixed up. See docs: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` e `` | 編輯 | 使用外部編輯器開啟 |
| `` o `` | 開啟檔案 | 使用預設軟體開啟 |
| `` i `` | 忽略或排除檔案 |  |
| `` ii `` | 添加到 .gitignore |  |
| `` ie `` | 添加到 .git/info/exclude |  |
| `` r `` | 重新整理檔案 |  |
| `` s `` | 收藏 | Stash all changes. Press capital S for variations (keep index, include untracked, staged only, unstaged only). |
| `` S `` | 檢視收藏選項 | View stash options (e.g. stash all, stash staged, stash unstaged). |
| `` Si `` | 收藏所有變更並保留預存區 |  |
| `` SU `` | 收藏所有變更，包括未追蹤檔案 |  |
| `` Ss `` | 收藏已預存變更 |  |
| `` Su `` | 收藏未預存變更 |  |
| `` a `` | 全部預存/取消預存 | Toggle staged/unstaged for all files in working tree. |
| `` <enter> `` | 選擇檔案中的單個程式碼塊/行，或展開/折疊目錄 | If the selected item is a file, focus the staging view so you can stage individual hunks/lines. If the selected item is a directory, collapse/expand it. |
| `` d `` | Discard changes |  |
| `` dc `` | 捨棄 | 檢視選中變動進行捨棄復原 |
| `` du `` | 刪除未預存變更 |  |
| `` g `` | Reset to upstream |  |
| `` gm `` | Mixed reset to upstream | Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as unstaged changes. |
| `` gs `` | Soft reset to upstream | Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as staged changes. |
| `` gh `` | Hard reset to upstream | Reset HEAD to the chosen commit, and discard all changes between the current and chosen commit, as well as all current modifications in the working tree. |
| `` D `` | 重設 | View reset options for working tree (e.g. nuking the working tree). |
| `` Dx `` | 刪除工作目錄 | 如果你想讓所有工作樹上的變更消失，這就是正確的選項。如果有未提交的子模組變更，它們將被收藏在子模組中。 |
| `` Du `` | 刪除未預存變更 |  |
| `` Dc `` | 刪除未追蹤檔案 |  |
| `` DS `` | 刪除已預存變更 | 這將創建一個新的存儲條目，其中只包含預存檔案，然後如果存儲條目不需要，將其刪除，因此工作樹僅保留未預存的變更。 |
| `` Ds `` | 軟重設 |  |
| `` Dm `` | mixed reset |  |
| `` Dh `` | 強制重設 |  |
| `` ` `` | 顯示檔案樹狀視圖 | Toggle file view between flat and tree layout. Flat layout shows all file paths in a single list, tree layout groups files by directory.<br><br>The default can be changed in the config file with the key 'gui.showFileTree'. |
| `` <ctrl+t> `` | 開啟外部差異工具 (git difftool) |  |
| `` M `` | View merge conflict options | View options for resolving merge conflicts. |
| `` f `` | 擷取 | 同步遠端異動 |
| `` - `` | Collapse all files | Collapse all directories in the files tree |
| `` = `` | Expand all files | Expand all directories in the file tree |
| `` 0 `` | Focus main view |  |
| `` / `` | 搜尋 |  |

## 次要

| Key | Action | Info |
|-----|--------|-------------|
| `` <tab> `` | 切換至另一個面板 (已預存/未預存更改) | Switch to other view (staged/unstaged changes). |
| `` <esc> `` | Exit back to side panel |  |
| `` / `` | 搜尋 |  |

## 狀態

| Key | Action | Info |
|-----|--------|-------------|
| `` o `` | 開啟設定檔案 | 使用預設軟體開啟 |
| `` e `` | 編輯設定檔案 | 使用外部編輯器開啟 |
| `` u `` | 檢查更新 |  |
| `` <enter> `` | 切換到最近使用的版本庫 |  |
| `` a `` | Show/cycle all branch logs |  |
| `` A `` | Show/cycle all branch logs (reverse) |  |
| `` 0 `` | Focus main view |  |

## 確認面板

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | 確認 |  |
| `` <esc> `` | 關閉/取消 |  |
| `` <ctrl+o> `` | 複製到剪貼簿 |  |

## 遠端

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | View branches |  |
| `` n `` | 新增遠端 |  |
| `` d `` | Remove | Remove the selected remote. Any local branches tracking a remote branch from the remote will be unaffected. |
| `` e `` | 編輯 | 編輯遠端 |
| `` f `` | 擷取 | 擷取遠端 |
| `` F `` | Add fork remote | Quickly add a fork remote by replacing the owner in the origin URL and optionally check out a branch from new remote. |
| `` / `` | 搜尋 |  |

## 遠端分支

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | 複製分支名稱到剪貼簿 |  |
| `` <space> `` | 檢出 | Checkout a new local branch based on the selected remote branch, or the remote branch as a detached head. |
| `` n `` | 新分支 |  |
| `` M `` | Merge |  |
| `` Mm `` | 合併到當前檢出的分支 | View options for merging the selected item into the current branch (regular merge, squash merge) |
| `` Mn `` | Non-fast-forward merge | Merge '{{.selectedBranch}}' into '{{.checkedOutBranch}}', creating a merge commit. |
| `` Mf `` | Fast-forward only merge | Fast-forward '{{.checkedOutBranch}}' to '{{.selectedBranch}}' without creating a merge commit. |
| `` Ms `` | Squash merge (uncommitted) | Squash merge '{{.selectedBranch}}' into the working tree. |
| `` MS `` | Squash merge (committed) | Squash merge '{{.selectedBranch}}' into '{{.checkedOutBranch}}' as a single commit. |
| `` r `` | Rebase options |  |
| `` rs `` | 將已檢出的分支變基至此分支 | Rebase the checked-out branch onto the selected branch. |
| `` ri `` | Interactive rebase | 開始一個互動變基，以中斷開始，這樣你可以在繼續之前更新TODO提交 |
| `` rb `` | Rebase onto base branch | Rebase the checked out branch onto its base branch (i.e. the closest main branch). |
| `` d `` | 刪除 | Delete the remote branch from the remote. |
| `` us `` | 設置為遠端 | 將此分支設為當前分支之遠端 |
| `` s `` | 排序規則 |  |
| `` g `` | Reset to ref |  |
| `` gm `` | Mixed reset | Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as unstaged changes. |
| `` gs `` | 軟重設 | Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as staged changes. |
| `` gh `` | 強制重設 | Reset HEAD to the chosen commit, and discard all changes between the current and chosen commit, as well as all current modifications in the working tree. |
| `` <ctrl+t> `` | 開啟外部差異工具 (git difftool) |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | 檢視提交 |  |
| `` w `` | 檢視工作目錄選項 |  |
| `` / `` | 搜尋 |  |
