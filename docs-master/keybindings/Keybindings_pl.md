_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go generate ./...` from the project root._

# Lazygit Skróty klawiszowe

## Globalne skróty klawiszowe

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+r> `` | Przełącz na ostatnie repozytorium |  |
| `` <pgup> (fn+up/shift+k) `` | Przewiń główne okno w górę |  |
| `` <pgdown> (fn+down/shift+j) `` | Przewiń główne okno w dół |  |
| `` @ `` | Pokaż opcje dziennika poleceń | Pokaż opcje dla dziennika poleceń, np. pokazywanie/ukrywanie dziennika poleceń i skupienie na dzienniku poleceń. |
| `` P `` | Wypchnij | Wypchnij bieżącą gałąź do jej gałęzi nadrzędnej. Jeśli nie skonfigurowano gałęzi nadrzędnej, zostaniesz poproszony o skonfigurowanie gałęzi nadrzędnej. |
| `` p `` | Pociągnij | Pociągnij zmiany z zdalnego dla bieżącej gałęzi. Jeśli nie skonfigurowano gałęzi nadrzędnej, zostaniesz poproszony o skonfigurowanie gałęzi nadrzędnej. |
| `` ) `` | Increase rename similarity threshold | Increase the similarity threshold for a deletion and addition pair to be treated as a rename.<br><br>The default can be changed in the config file with the key 'git.renameSimilarityThreshold'. |
| `` ( `` | Decrease rename similarity threshold | Decrease the similarity threshold for a deletion and addition pair to be treated as a rename.<br><br>The default can be changed in the config file with the key 'git.renameSimilarityThreshold'. |
| `` } `` | Zwiększ rozmiar kontekstu w widoku różnic | Increase the amount of the context shown around changes in the diff view.<br><br>The default can be changed in the config file with the key 'git.diffContextSize'. |
| `` { `` | Zmniejsz rozmiar kontekstu w widoku różnic | Decrease the amount of the context shown around changes in the diff view.<br><br>The default can be changed in the config file with the key 'git.diffContextSize'. |
| `` : `` | Execute shell command | Bring up a prompt where you can enter a shell command to execute. |
| `` <ctrl+p> `` | Wyświetl opcje niestandardowej łatki |  |
| `` m `` | Pokaż opcje scalania/rebase | Pokaż opcje do przerwania/kontynuowania/pominięcia bieżącego scalania/rebase. |
| `` mc `` | Continue rebase / merge | Pokaż opcje do przerwania/kontynuowania/pominięcia bieżącego scalania/rebase. |
| `` ma `` | Abort rebase / merge | Pokaż opcje do przerwania/kontynuowania/pominięcia bieżącego scalania/rebase. |
| `` ms `` | Skip current rebase commit | Pokaż opcje do przerwania/kontynuowania/pominięcia bieżącego scalania/rebase. |
| `` R `` | Odśwież | Odśwież stan git (tj. uruchom `git status`, `git branch`, itp. w tle, aby zaktualizować zawartość paneli). To nie uruchamia `git fetch`. |
| `` + `` | Następny tryb ekranu (normalny/półpełny/pełnoekranowy) |  |
| `` _ `` | Poprzedni tryb ekranu |  |
| `` \| `` | Cycle pagers | Choose the next pager in the list of configured pagers |
| `` <esc> `` | Anuluj |  |
| `` ? `` | Otwórz menu przypisań klawiszy |  |
| `` <ctrl+s> `` | Pokaż opcje filtrowania | Pokaż opcje filtrowania dziennika commitów, tak aby pokazywane były tylko commity pasujące do filtra. |
| `` W `` | Pokaż opcje różnicowania | Pokaż opcje dotyczące różnicowania dwóch refów, np. różnicowanie względem wybranego refa, wprowadzanie refa do różnicowania i odwracanie kierunku różnic. |
| `` <ctrl+e> `` | Pokaż opcje różnicowania | Pokaż opcje dotyczące różnicowania dwóch refów, np. różnicowanie względem wybranego refa, wprowadzanie refa do różnicowania i odwracanie kierunku różnic. |
| `` q `` | Wyjdź |  |
| `` <ctrl+z> `` | Suspend the application |  |
| `` <ctrl+w> `` | Przełącz białe znaki | Toggle whether or not whitespace changes are shown in the diff view.<br><br>The default can be changed in the config file with the key 'git.ignoreWhitespaceInDiffView'. |
| `` z `` | Cofnij | Dziennik reflog zostanie użyty do określenia, jakie polecenie git należy uruchomić, aby cofnąć ostatnie polecenie git. Nie obejmuje to zmian w drzewie roboczym; brane są pod uwagę tylko commity. |
| `` Z `` | Ponów | Dziennik reflog zostanie użyty do określenia, jakie polecenie git należy uruchomić, aby ponowić ostatnie polecenie git. Nie obejmuje to zmian w drzewie roboczym; brane są pod uwagę tylko commity. |

## Nawigacja panelu listy

| Key | Action | Info |
|-----|--------|-------------|
| `` , `` | Poprzednia strona |  |
| `` . `` | Następna strona |  |
| `` < (<home>) `` | Przewiń do góry |  |
| `` > (<end>) `` | Przewiń do dołu |  |
| `` v `` | Przełącz zaznaczenie zakresu |  |
| `` <shift+down> `` | Zaznacz zakres w dół |  |
| `` <shift+up> `` | Zaznacz zakres w górę |  |
| `` / `` | Szukaj w bieżącym widoku po tekście |  |
| `` H `` | Przewiń w lewo |  |
| `` L `` | Przewiń w prawo |  |
| `` ] `` | Następna zakładka |  |
| `` [ `` | Poprzednia zakładka |  |

## Commity

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | Copy abbreviated commit hash to clipboard |  |
| `` <ctrl+r> `` | Resetuj wybrane (cherry-picked) commity |  |
| `` b `` | Zobacz opcje bisect |  |
| `` bb `` | Mark commit as bad |  |
| `` bg `` | Mark commit as good |  |
| `` bs `` | Skip current bisect commit |  |
| `` bS `` | Skip selected commit |  |
| `` br `` | Resetuj bisect |  |
| `` bb `` | Mark commit as bad (start bisect) |  |
| `` bg `` | Mark commit as good (start bisect) |  |
| `` bt `` | Wybierz terminy bisect |  |
| `` s `` | Scal | Scal wybrany commit z commitami poniżej. Wiadomość wybranego commita zostanie dołączona do commita poniżej. |
| `` f `` | Poprawka | Włącz wybrany commit do commita poniżej. Podobnie do fixup, ale wiadomość wybranego commita zostanie odrzucona. |
| `` ff `` | Poprawka | Włącz wybrany commit do commita poniżej. Podobnie do fixup, ale wiadomość wybranego commita zostanie odrzucona. |
| `` fc `` | Fixup and use this commit's message | Squash the selected commit into the commit below, using this commit's message, discarding the message of the commit below. |
| `` c `` | Set fixup message | Set the message option for the fixup commit. The -C option means to use this commit's message instead of the target commit's message. |
| `` r `` | Przeformułuj | Przeformułuj wiadomość wybranego commita. |
| `` R `` | Przeformułuj za pomocą edytora |  |
| `` d `` | Usuń | Usuń wybrany commit. To usunie commit z gałęzi za pomocą rebazowania. Jeśli commit wprowadza zmiany, od których zależą późniejsze commity, być może będziesz musiał rozwiązać konflikty scalania. |
| `` e `` | Edytuj (rozpocznij interaktywne rebazowanie) | Edytuj wybrany commit. Użyj tego, aby rozpocząć interaktywne rebazowanie od wybranego commita. Podczas trwania rebazowania, to oznaczy wybrany commit do edycji, co oznacza, że po kontynuacji rebazowania, rebazowanie zostanie wstrzymane na wybranym commicie, aby umożliwić wprowadzenie zmian. |
| `` i `` | Rozpocznij interaktywny rebase | Rozpocznij interaktywny rebase dla commitów na twoim branchu. To będzie zawierać wszystkie commity od HEAD do pierwszego commita scalenia lub commita głównego brancha.<br>Jeśli chcesz zamiast tego rozpocząć interaktywny rebase od wybranego commita, naciśnij `e`. |
| `` p `` | Wybierz | Oznacz wybrany commit do wybrania (podczas rebazowania). Oznacza to, że commit zostanie zachowany po kontynuacji rebazowania. |
| `` F `` | Utwórz commit fixup | Utwórz commit 'fixup!' dla wybranego commita. Później możesz nacisnąć `S` na tym samym commicie, aby zastosować wszystkie powyższe commity fixup. |
| `` S `` | Zastosuj commity fixup | Scal wszystkie commity 'fixup!', albo powyżej wybranego commita, albo wszystkie w bieżącej gałęzi (autosquash). |
| `` <alt+down> `` | Przesuń commit w dół |  |
| `` <alt+up> `` | Przesuń commit w górę |  |
| `` V `` | Wklej (cherry-pick) |  |
| `` B `` | Oznacz jako bazowy commit dla rebase | Wybierz bazowy commit dla następnego rebase. Kiedy robisz rebase na branch, tylko commity powyżej bazowego commita zostaną przeniesione. Używa to polecenia `git rebase --onto`. |
| `` A `` | Popraw | Popraw commit ze zmianami zatwierdzonymi. Jeśli wybrany commit jest commit HEAD, to wykona `git commit --amend`. W przeciwnym razie commit zostanie poprawiony za pomocą rebazowania. |
| `` a `` | Popraw atrybut commita | Ustaw/Resetuj autora commita lub ustaw współautora. |
| `` t `` | Cofnij | Utwórz commit cofający dla wybranego commita, który stosuje zmiany wybranego commita w odwrotnej kolejności. |
| `` T `` | Otaguj commit | Utwórz nowy tag wskazujący na wybrany commit. Zostaniesz poproszony o wprowadzenie nazwy tagu i opcjonalnego opisu. |
| `` <ctrl+l> `` | Zobacz opcje logów | Zobacz opcje dla logów commitów, np. zmiana kolejności sortowania, ukrywanie grafu gita, pokazywanie całego grafu gita. |
| `` G `` | Open pull request in browser |  |
| `` <space> `` | Przełącz | Przełącz wybrany commit jako odłączoną HEAD. |
| `` y `` | Kopiuj atrybut commita do schowka | Kopiuj atrybut commita do schowka (np. hash, URL, różnice, wiadomość, autor). |
| `` o `` | Otwórz commit w przeglądarce |  |
| `` n `` | Utwórz nową gałąź z commita |  |
| `` N `` | Move commits to new branch | Create a new branch and move the unpushed commits of the current branch to it. Useful if you meant to start new work and forgot to create a new branch first.<br><br>Note that this disregards the selection, the new branch is always created either from the main branch or stacked on top of the current branch (you get to choose which). |
| `` g `` | Reset | Wyświetl opcje resetu (miękki/mieszany/twardy) do wybranego elementu. |
| `` gm `` | Mixed reset | Resetuj HEAD do wybranego commita, zachowując zmiany między bieżącym a wybranym commit jako zmiany niezatwierdzone. |
| `` gs `` | Miękki reset | Resetuj HEAD do wybranego commita, zachowując zmiany między bieżącym a wybranym commit jako zmiany zatwierdzone. |
| `` gh `` | Twardy reset | Resetuj HEAD do wybranego commita, odrzucając wszystkie zmiany między bieżącym a wybranym commit, jak również wszystkie bieżące modyfikacje w drzewie roboczym. |
| `` C `` | Kopiuj (cherry-pick) | Oznacz commit jako skopiowany. Następnie, w widoku lokalnych commitów, możesz nacisnąć `V`, aby wkleić (cherry-pick) skopiowane commity do sprawdzonej gałęzi. W dowolnym momencie możesz nacisnąć `<esc>`, aby anulować zaznaczenie. |
| `` <ctrl+t> `` | Otwórz zewnętrzne narzędzie różnic (git difftool) |  |
| `` * `` | Select commits of current branch |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | Wyświetl pliki |  |
| `` w `` | Zobacz opcje drzewa pracy |  |
| `` / `` | Szukaj w bieżącym widoku po tekście |  |

## Dodatkowy

| Key | Action | Info |
|-----|--------|-------------|
| `` <tab> `` | Przełącz widok | Przełącz na inny widok (zatwierdzone/niezatwierdzone zmiany). |
| `` <esc> `` | Exit back to side panel |  |
| `` / `` | Szukaj w bieżącym widoku po tekście |  |

## Drzewa pracy

| Key | Action | Info |
|-----|--------|-------------|
| `` n `` | Nowe drzewo pracy |  |
| `` <space> `` | Przełącz | Przełącz do wybranego drzewa pracy. |
| `` o `` | Otwórz w edytorze |  |
| `` d `` | Usuń | Usuń wybrane drzewo pracy. To usunie zarówno katalog drzewa pracy, jak i metadane o drzewie pracy w katalogu .git. |
| `` / `` | Filtruj bieżący widok po tekście |  |

## Główny panel (budowanie łatki)

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | Idź do poprzedniego fragmentu |  |
| `` <right> `` | Idź do następnego fragmentu |  |
| `` v `` | Przełącz zaznaczenie zakresu |  |
| `` a `` | Toggle hunk selection | Toggle line-by-line vs. hunk selection mode. |
| `` <ctrl+o> `` | Kopiuj zaznaczony tekst do schowka |  |
| `` o `` | Otwórz plik | Otwórz plik w domyślnej aplikacji. |
| `` e `` | Edytuj plik | Otwórz plik w zewnętrznym edytorze. |
| `` <space> `` | Przełącz linie w łatce |  |
| `` d `` | Remove lines from commit | Remove the selected lines from this commit. This runs an interactive rebase in the background, so you may get a merge conflict if a later commit also changes these lines. |
| `` <esc> `` | Wyjdź z budowniczego niestandardowej łatki |  |
| `` / `` | Szukaj w bieżącym widoku po tekście |  |

## Input prompt

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Potwierdź |  |
| `` <esc> `` | Zamknij/Anuluj |  |

## Lokalne gałęzie

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | Kopiuj nazwę gałęzi do schowka |  |
| `` i `` | Pokaż opcje git-flow |  |
| `` iF `` | Finish git-flow branch |  |
| `` if `` | Start git-flow feature |  |
| `` ih `` | Start git-flow hotfix |  |
| `` ib `` | Start git-flow bugfix |  |
| `` ir `` | Start git-flow release |  |
| `` <space> `` | Przełącz | Przełącz wybrany element. |
| `` n `` | Nowa gałąź |  |
| `` N `` | Move commits to new branch | Create a new branch and move the unpushed commits of the current branch to it. Useful if you meant to start new work and forgot to create a new branch first.<br><br>Note that this disregards the selection, the new branch is always created either from the main branch or stacked on top of the current branch (you get to choose which). |
| `` o `` | Utwórz żądanie ściągnięcia |  |
| `` O `` | Zobacz opcje tworzenia pull requesta |  |
| `` G `` | Open pull request in browser |  |
| `` <ctrl+y> `` | Kopiuj adres URL żądania ściągnięcia do schowka |  |
| `` c `` | Przełącz według nazwy | Przełącz według nazwy. W polu wprowadzania możesz wpisać '-' aby przełączyć się na ostatnią gałąź. |
| `` - `` | Checkout previous branch |  |
| `` F `` | Wymuś przełączenie | Wymuś przełączenie wybranej gałęzi. To spowoduje odrzucenie wszystkich lokalnych zmian w drzewie roboczym przed przełączeniem na wybraną gałąź. |
| `` d `` | Usuń | Wyświetl opcje usuwania lokalnej/odległej gałęzi. |
| `` dc `` | Usuń lokalną gałąź | Wyświetl opcje usuwania lokalnej/odległej gałęzi. |
| `` dr `` | Usuń gałąź zdalną | Usuń gałąź zdalną ze zdalnego. |
| `` db `` | Delete local and remote branch | Wyświetl opcje usuwania lokalnej/odległej gałęzi. |
| `` r `` | Przebazuj | Przebazuj przełączoną gałąź na wybraną gałąź. |
| `` rs `` | Simple rebase | Przebazuj przełączoną gałąź na wybraną gałąź. |
| `` ri `` | Interactive rebase | Rozpocznij interaktywny rebase z przerwaniem na początku, abyś mógł zaktualizować commity TODO przed kontynuacją. |
| `` rb `` | Rebase onto base branch | Rebase the checked out branch onto its base branch (i.e. the closest main branch). |
| `` M `` | Scal | Scal wybraną gałąź z aktualnie sprawdzoną gałęzią. |
| `` Mm `` | Scal | Scal wybraną gałąź z aktualnie sprawdzoną gałęzią. |
| `` Mn `` | Regular merge (with merge commit) | Merge '{{.selectedBranch}}' into '{{.checkedOutBranch}}', creating a merge commit. |
| `` Mf `` | Regular merge (fast-forward) | Fast-forward '{{.checkedOutBranch}}' to '{{.selectedBranch}}' without creating a merge commit. |
| `` Ms `` | Squash merge (uncommitted) | Squash merge '{{.selectedBranch}}' into the working tree. |
| `` MS `` | Squash merge (committed) | Squash merge '{{.selectedBranch}}' into '{{.checkedOutBranch}}' as a single commit. |
| `` f `` | Szybkie przewijanie | Szybkie przewijanie wybranej gałęzi z jej źródła. |
| `` T `` | Nowy tag |  |
| `` s `` | Kolejność sortowania |  |
| `` g `` | Reset to ref |  |
| `` gm `` | Mixed reset | Resetuj HEAD do wybranego commita, zachowując zmiany między bieżącym a wybranym commit jako zmiany niezatwierdzone. |
| `` gs `` | Miękki reset | Resetuj HEAD do wybranego commita, zachowując zmiany między bieżącym a wybranym commit jako zmiany zatwierdzone. |
| `` gh `` | Twardy reset | Resetuj HEAD do wybranego commita, odrzucając wszystkie zmiany między bieżącym a wybranym commit, jak również wszystkie bieżące modyfikacje w drzewie roboczym. |
| `` R `` | Zmień nazwę gałęzi |  |
| `` u `` | Pokaż opcje upstream | Pokaż opcje dotyczące upstream gałęzi, np. ustawianie/usuwanie upstream i resetowanie do upstream. |
| `` ud `` | Wyświetl rozbieżność od upstream | Pokaż opcje dotyczące upstream gałęzi, np. ustawianie/usuwanie upstream i resetowanie do upstream. |
| `` uD `` | View divergence from base branch |  |
| `` us `` | Ustaw upstream wybranej gałęzi |  |
| `` uu `` | Usuń upstream wybranej gałęzi |  |
| `` ugm `` | Mixed reset to upstream | Resetuj HEAD do wybranego commita, zachowując zmiany między bieżącym a wybranym commit jako zmiany niezatwierdzone. |
| `` ugs `` | Soft reset to upstream | Resetuj HEAD do wybranego commita, zachowując zmiany między bieżącym a wybranym commit jako zmiany zatwierdzone. |
| `` ugh `` | Hard reset to upstream | Resetuj HEAD do wybranego commita, odrzucając wszystkie zmiany między bieżącym a wybranym commit, jak również wszystkie bieżące modyfikacje w drzewie roboczym. |
| `` urs `` | Simple rebase onto upstream |  |
| `` uri `` | Interactive rebase onto upstream | Rozpocznij interaktywny rebase z przerwaniem na początku, abyś mógł zaktualizować commity TODO przed kontynuacją. |
| `` urb `` | Rebase onto base branch | Rebase the checked out branch onto its base branch (i.e. the closest main branch). |
| `` <ctrl+t> `` | Otwórz zewnętrzne narzędzie różnic (git difftool) |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | Pokaż commity |  |
| `` w `` | Zobacz opcje drzewa pracy |  |
| `` / `` | Filtruj bieżący widok po tekście |  |

## Menu

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Wykonaj |  |
| `` <esc> `` | Zamknij/Anuluj |  |
| `` / `` | Filtruj bieżący widok po tekście |  |

## Panel główny (normalny)

| Key | Action | Info |
|-----|--------|-------------|
| `` <mouse wheel down> (fn+up) `` | Przewiń w dół |  |
| `` <mouse wheel up> (fn+down) `` | Przewiń w górę |  |
| `` <tab> `` | Przełącz widok | Przełącz na inny widok (zatwierdzone/niezatwierdzone zmiany). |
| `` <esc> `` | Exit back to side panel |  |
| `` / `` | Szukaj w bieżącym widoku po tekście |  |

## Panel główny (scalanie)

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | Wybierz fragment |  |
| `` b `` | Wybierz wszystkie fragmenty |  |
| `` <up> `` | Poprzedni fragment |  |
| `` <down> `` | Następny fragment |  |
| `` <left> `` | Poprzedni konflikt |  |
| `` <right> `` | Następny konflikt |  |
| `` z `` | Cofnij | Cofnij ostatnie rozwiązanie konfliktu scalania. |
| `` e `` | Edytuj plik | Otwórz plik w zewnętrznym edytorze. |
| `` o `` | Otwórz plik | Otwórz plik w domyślnej aplikacji. |
| `` M `` | View merge conflict options | View options for resolving merge conflicts. |
| `` <esc> `` | Wróć do panelu plików |  |

## Panel główny (zatwierdzanie)

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | Idź do poprzedniego fragmentu |  |
| `` <right> `` | Idź do następnego fragmentu |  |
| `` v `` | Przełącz zaznaczenie zakresu |  |
| `` a `` | Toggle hunk selection | Toggle line-by-line vs. hunk selection mode. |
| `` <ctrl+o> `` | Kopiuj zaznaczony tekst do schowka |  |
| `` <space> `` | Zatwierdź | Przełącz zaznaczenie zatwierdzone/niezatwierdzone. |
| `` d `` | Odrzuć | Gdy zaznaczona jest niezatwierdzona zmiana, odrzuć ją używając `git reset`. Gdy zaznaczona jest zatwierdzona zmiana, cofnij zatwierdzenie. |
| `` o `` | Otwórz plik | Otwórz plik w domyślnej aplikacji. |
| `` e `` | Edytuj plik | Otwórz plik w zewnętrznym edytorze. |
| `` <esc> `` | Wróć do panelu plików |  |
| `` <tab> `` | Przełącz widok | Przełącz na inny widok (zatwierdzone/niezatwierdzone zmiany). |
| `` E `` | Edytuj fragment | Edytuj wybrany fragment w zewnętrznym edytorze. |
| `` c `` | Commit | Zatwierdź zmiany zatwierdzone. |
| `` w `` | Zatwierdź zmiany bez hooka pre-commit |  |
| `` C `` | Zatwierdź zmiany używając edytora git |  |
| `` <ctrl+f> `` | Znajdź bazowy commit do poprawki | Znajdź commit, na którym opierają się Twoje obecne zmiany, w celu poprawienia/zmiany commita. To pozwala Ci uniknąć przeglądania commitów w Twojej gałęzi jeden po drugim, aby zobaczyć, który commit powinien być poprawiony/zmieniony. Zobacz dokumentację: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` / `` | Szukaj w bieżącym widoku po tekście |  |

## Panel potwierdzenia

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Potwierdź |  |
| `` <esc> `` | Zamknij/Anuluj |  |
| `` <ctrl+o> `` | Kopiuj do schowka |  |

## Pliki

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | Kopiuj ścieżkę do schowka |  |
| `` <space> `` | Zatwierdź | Przełącz zatwierdzenie dla wybranego pliku. |
| `` <ctrl+b> `` | Filtruj pliki według statusu |  |
| `` <ctrl+b>s `` | Pokaż tylko zatwierdzone pliki |  |
| `` <ctrl+b>u `` | Pokaż tylko niezatwierdzone pliki |  |
| `` <ctrl+b>t `` | Show only tracked files |  |
| `` <ctrl+b>T `` | Show only untracked files |  |
| `` <ctrl+b>r `` | No filter |  |
| `` y `` | Kopiuj do schowka |  |
| `` yn `` | Nazwa pliku |  |
| `` yp `` | Relative path |  |
| `` yP `` | Absolute path |  |
| `` ys `` | Różnice wybranego pliku | Jeśli istnieją zatwierdzone elementy, ta komenda bierze pod uwagę tylko je. W przeciwnym razie bierze pod uwagę wszystkie niezatwierdzone. |
| `` ya `` | Różnice wszystkich plików | Jeśli istnieją zatwierdzone elementy, ta komenda bierze pod uwagę tylko je. W przeciwnym razie bierze pod uwagę wszystkie niezatwierdzone. |
| `` c `` | Commit | Zatwierdź zmiany zatwierdzone. |
| `` w `` | Zatwierdź zmiany bez hooka pre-commit |  |
| `` A `` | Popraw ostatni commit |  |
| `` C `` | Zatwierdź zmiany używając edytora git |  |
| `` <ctrl+f> `` | Znajdź bazowy commit do poprawki | Znajdź commit, na którym opierają się Twoje obecne zmiany, w celu poprawienia/zmiany commita. To pozwala Ci uniknąć przeglądania commitów w Twojej gałęzi jeden po drugim, aby zobaczyć, który commit powinien być poprawiony/zmieniony. Zobacz dokumentację: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` e `` | Edytuj | Otwórz plik w zewnętrznym edytorze. |
| `` o `` | Otwórz plik | Otwórz plik w domyślnej aplikacji. |
| `` i `` | Ignoruj lub wyklucz plik |  |
| `` ii `` | Dodaj do .gitignore |  |
| `` ie `` | Dodaj do .git/info/exclude |  |
| `` r `` | Odśwież pliki |  |
| `` s `` | Schowaj | Schowaj wszystkie zmiany. Dla innych wariantów schowania, użyj klawisza wyświetlania opcji schowka. |
| `` S `` | Wyświetl opcje schowka | Wyświetl opcje schowka (np. schowaj wszystko, schowaj zatwierdzone, schowaj niezatwierdzone). |
| `` Si `` | Schowaj wszystkie zmiany i zachowaj indeks |  |
| `` SU `` | Schowaj wszystkie zmiany włącznie z nieśledzonymi plikami |  |
| `` Ss `` | Schowaj zatwierdzone zmiany |  |
| `` Su `` | Schowaj niezatwierdzone zmiany |  |
| `` a `` | Zatwierdź wszystko | Przełącz zatwierdzenie/odznaczenie dla wszystkich plików w drzewie roboczym. |
| `` <enter> `` | Zatwierdź linie / Zwiń katalog | Jeśli wybrany element jest plikiem, skup się na widoku zatwierdzania, aby móc zatwierdzać poszczególne fragmenty/linie. Jeśli wybrany element jest katalogiem, zwiń/rozwiń go. |
| `` d `` | Discard changes |  |
| `` dc `` | Odrzuć | Wyświetl opcje odrzucania zmian w wybranym pliku. |
| `` du `` | Odrzuć niezatwierdzone zmiany |  |
| `` g `` | Reset to upstream |  |
| `` gm `` | Mixed reset to upstream | Resetuj HEAD do wybranego commita, zachowując zmiany między bieżącym a wybranym commit jako zmiany niezatwierdzone. |
| `` gs `` | Soft reset to upstream | Resetuj HEAD do wybranego commita, zachowując zmiany między bieżącym a wybranym commit jako zmiany zatwierdzone. |
| `` gh `` | Hard reset to upstream | Resetuj HEAD do wybranego commita, odrzucając wszystkie zmiany między bieżącym a wybranym commit, jak również wszystkie bieżące modyfikacje w drzewie roboczym. |
| `` D `` | Reset | Wyświetl opcje resetu dla drzewa roboczego (np. zniszczenie drzewa roboczego). |
| `` Dx `` | Zniszcz drzewo robocze | Jeśli chcesz, aby wszystkie zmiany w drzewie pracy zniknęły, to jest sposób na to. Jeśli są brudne zmiany w submodule, to zostaną one zapisane w submodule(s). |
| `` Du `` | Odrzuć niezatwierdzone zmiany |  |
| `` Dc `` | Odrzuć nieśledzone pliki |  |
| `` DS `` | Odrzuć zatwierdzone zmiany | To stworzy nowy wpis stash zawierający tylko pliki w stanie staged, a następnie go usunie, tak że drzewo pracy zostanie tylko ze zmianami niezatwierdzonymi |
| `` Ds `` | Miękki reset |  |
| `` Dm `` | mixed reset |  |
| `` Dh `` | Twardy reset |  |
| `` ` `` | Przełącz widok drzewa plików | Toggle file view between flat and tree layout. Flat layout shows all file paths in a single list, tree layout groups files by directory.<br><br>The default can be changed in the config file with the key 'gui.showFileTree'. |
| `` <ctrl+t> `` | Otwórz zewnętrzne narzędzie różnic (git difftool) |  |
| `` M `` | View merge conflict options | View options for resolving merge conflicts. |
| `` f `` | Pobierz | Pobierz zmiany ze zdalnego serwera. |
| `` - `` | Collapse all files | Collapse all directories in the files tree |
| `` = `` | Expand all files | Expand all directories in the file tree |
| `` 0 `` | Focus main view |  |
| `` / `` | Filtruj bieżący widok po tekście |  |

## Pliki commita

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | Kopiuj ścieżkę do schowka |  |
| `` y `` | Copy to clipboard |  |
| `` yn `` | Nazwa pliku |  |
| `` yp `` | Relative path |  |
| `` yP `` | Absolute path |  |
| `` ys `` | Różnice wybranego pliku |  |
| `` ya `` | Różnice wszystkich plików |  |
| `` yc `` | Content of selected file |  |
| `` c `` | Przełącz | Przełącz plik. Zastępuje plik w twoim drzewie roboczym wersją z wybranego commita. |
| `` d `` | Odrzuć | Odrzuć zmiany w tym pliku z tego commita. Uruchamia interaktywny rebase w tle, więc możesz otrzymać konflikt scalania, jeśli późniejszy commit również zmienia ten plik. |
| `` o `` | Otwórz plik | Otwórz plik w domyślnej aplikacji. |
| `` e `` | Edytuj | Otwórz plik w zewnętrznym edytorze. |
| `` <ctrl+t> `` | Otwórz zewnętrzne narzędzie różnic (git difftool) |  |
| `` <space> `` | Przełącz plik włączony w łatkę | Przełącz, czy plik jest włączony w niestandardową łatkę. Zobacz https://github.com/jesseduffield/lazygit#rebase-magic-custom-patches. |
| `` a `` | Przełącz wszystkie pliki | Dodaj/usuń wszystkie pliki commita do niestandardowej łatki. Zobacz https://github.com/jesseduffield/lazygit#rebase-magic-custom-patches. |
| `` <enter> `` | Wejdź do pliku / Przełącz zwiń katalog | Jeśli plik jest wybrany, wejdź do pliku, aby móc dodawać/usuwać poszczególne linie do niestandardowej łatki. Jeśli wybrany jest katalog, przełącz katalog. |
| `` ` `` | Przełącz widok drzewa plików | Toggle file view between flat and tree layout. Flat layout shows all file paths in a single list, tree layout groups files by directory.<br><br>The default can be changed in the config file with the key 'gui.showFileTree'. |
| `` - `` | Collapse all files | Collapse all directories in the files tree |
| `` = `` | Expand all files | Expand all directories in the file tree |
| `` 0 `` | Focus main view |  |
| `` / `` | Filtruj bieżący widok po tekście |  |

## Podsumowanie commita

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Potwierdź |  |
| `` <esc> `` | Zamknij |  |

## Reflog

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | Copy abbreviated commit hash to clipboard |  |
| `` <space> `` | Przełącz | Przełącz wybrany commit jako odłączoną HEAD. |
| `` y `` | Kopiuj atrybut commita do schowka | Kopiuj atrybut commita do schowka (np. hash, URL, różnice, wiadomość, autor). |
| `` o `` | Otwórz commit w przeglądarce |  |
| `` n `` | Utwórz nową gałąź z commita |  |
| `` N `` | Move commits to new branch | Create a new branch and move the unpushed commits of the current branch to it. Useful if you meant to start new work and forgot to create a new branch first.<br><br>Note that this disregards the selection, the new branch is always created either from the main branch or stacked on top of the current branch (you get to choose which). |
| `` g `` | Reset to ref |  |
| `` gm `` | Mixed reset | Resetuj HEAD do wybranego commita, zachowując zmiany między bieżącym a wybranym commit jako zmiany niezatwierdzone. |
| `` gs `` | Miękki reset | Resetuj HEAD do wybranego commita, zachowując zmiany między bieżącym a wybranym commit jako zmiany zatwierdzone. |
| `` gh `` | Twardy reset | Resetuj HEAD do wybranego commita, odrzucając wszystkie zmiany między bieżącym a wybranym commit, jak również wszystkie bieżące modyfikacje w drzewie roboczym. |
| `` C `` | Kopiuj (cherry-pick) | Oznacz commit jako skopiowany. Następnie, w widoku lokalnych commitów, możesz nacisnąć `V`, aby wkleić (cherry-pick) skopiowane commity do sprawdzonej gałęzi. W dowolnym momencie możesz nacisnąć `<esc>`, aby anulować zaznaczenie. |
| `` <ctrl+r> `` | Resetuj wybrane (cherry-picked) commity |  |
| `` <ctrl+t> `` | Otwórz zewnętrzne narzędzie różnic (git difftool) |  |
| `` * `` | Select commits of current branch |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | Pokaż commity |  |
| `` w `` | Zobacz opcje drzewa pracy |  |
| `` / `` | Filtruj bieżący widok po tekście |  |

## Schowek

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | Zastosuj | Zastosuj wpis schowka do katalogu roboczego. |
| `` g `` | Wyciągnij | Zastosuj wpis schowka do katalogu roboczego i usuń wpis schowka. |
| `` d `` | Usuń | Usuń wpis schowka z listy schowka. |
| `` n `` | Nowa gałąź | Utwórz nową gałąź z wybranego wpisu schowka. Działa poprzez przełączenie git na commit, na którym wpis schowka został utworzony, tworzenie nowej gałęzi z tego commita, a następnie zastosowanie wpisu schowka do nowej gałęzi jako dodatkowego commita. |
| `` r `` | Zmień nazwę schowka |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | Wyświetl pliki |  |
| `` w `` | Zobacz opcje drzewa pracy |  |
| `` / `` | Filtruj bieżący widok po tekście |  |

## Status

| Key | Action | Info |
|-----|--------|-------------|
| `` o `` | Otwórz plik konfiguracyjny | Otwórz plik w domyślnej aplikacji. |
| `` e `` | Edytuj plik konfiguracyjny | Otwórz plik w zewnętrznym edytorze. |
| `` u `` | Sprawdź aktualizacje |  |
| `` <enter> `` | Przełącz na ostatnie repozytorium |  |
| `` a `` | Show/cycle all branch logs |  |
| `` A `` | Show/cycle all branch logs (reverse) |  |
| `` 0 `` | Focus main view |  |

## Sub-commity

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | Copy abbreviated commit hash to clipboard |  |
| `` <space> `` | Przełącz | Przełącz wybrany commit jako odłączoną HEAD. |
| `` y `` | Kopiuj atrybut commita do schowka | Kopiuj atrybut commita do schowka (np. hash, URL, różnice, wiadomość, autor). |
| `` o `` | Otwórz commit w przeglądarce |  |
| `` n `` | Utwórz nową gałąź z commita |  |
| `` N `` | Move commits to new branch | Create a new branch and move the unpushed commits of the current branch to it. Useful if you meant to start new work and forgot to create a new branch first.<br><br>Note that this disregards the selection, the new branch is always created either from the main branch or stacked on top of the current branch (you get to choose which). |
| `` g `` | Reset to ref |  |
| `` gm `` | Mixed reset | Resetuj HEAD do wybranego commita, zachowując zmiany między bieżącym a wybranym commit jako zmiany niezatwierdzone. |
| `` gs `` | Miękki reset | Resetuj HEAD do wybranego commita, zachowując zmiany między bieżącym a wybranym commit jako zmiany zatwierdzone. |
| `` gh `` | Twardy reset | Resetuj HEAD do wybranego commita, odrzucając wszystkie zmiany między bieżącym a wybranym commit, jak również wszystkie bieżące modyfikacje w drzewie roboczym. |
| `` C `` | Kopiuj (cherry-pick) | Oznacz commit jako skopiowany. Następnie, w widoku lokalnych commitów, możesz nacisnąć `V`, aby wkleić (cherry-pick) skopiowane commity do sprawdzonej gałęzi. W dowolnym momencie możesz nacisnąć `<esc>`, aby anulować zaznaczenie. |
| `` <ctrl+r> `` | Resetuj wybrane (cherry-picked) commity |  |
| `` <ctrl+t> `` | Otwórz zewnętrzne narzędzie różnic (git difftool) |  |
| `` * `` | Select commits of current branch |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | Wyświetl pliki |  |
| `` w `` | Zobacz opcje drzewa pracy |  |
| `` / `` | Szukaj w bieżącym widoku po tekście |  |

## Submoduły

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | Kopiuj nazwę submodułu do schowka |  |
| `` <enter> `` | Wejdź | Wejdź do submodułu. Po wejściu do submodułu możesz nacisnąć `<esc>`, aby wrócić do repozytorium nadrzędnego. |
| `` d `` | Usuń | Usuń wybrany submoduł i odpowiadający mu katalog. |
| `` u `` | Aktualizuj | Aktualizuj wybrany submoduł. |
| `` n `` | Nowy submoduł |  |
| `` e `` | Zaktualizuj URL submodułu |  |
| `` i `` | Zainicjuj | Zainicjuj wybrany submoduł, aby przygotować do pobrania. Prawdopodobnie chcesz to kontynuować, wywołując akcję 'update', aby pobrać submoduł. |
| `` b `` | Pokaż opcje masowych operacji na submodułach |  |
| `` bi `` | Masowe inicjowanie submodułów |  |
| `` bu `` | Masowa aktualizacja submodułów |  |
| `` br `` | Bulk init and update submodules recursively |  |
| `` bd `` | Masowe wyłączanie submodułów |  |
| `` / `` | Filtruj bieżący widok po tekście |  |

## Tagi

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | Copy tag to clipboard |  |
| `` <space> `` | Przełącz | Przełącz wybrany tag jako odłączoną głowę (detached HEAD). |
| `` n `` | Nowy tag | Utwórz nowy tag z bieżącego commita. Zostaniesz poproszony o wprowadzenie nazwy tagu i opcjonalnego opisu. |
| `` d `` | Usuń |  |
| `` dc `` | Usuń lokalny tag | Wyświetl opcje usuwania lokalnego/odległego tagu. |
| `` dr `` | Usuń zdalny tag | Wyświetl opcje usuwania lokalnego/odległego tagu. |
| `` db `` | Delete local and remote tag | Wyświetl opcje usuwania lokalnego/odległego tagu. |
| `` P `` | Wyślij tag | Wyślij wybrany tag do zdalnego. Zostaniesz poproszony o wybranie zdalnego. |
| `` g `` | Reset to ref |  |
| `` gm `` | Mixed reset | Resetuj HEAD do wybranego commita, zachowując zmiany między bieżącym a wybranym commit jako zmiany niezatwierdzone. |
| `` gs `` | Miękki reset | Resetuj HEAD do wybranego commita, zachowując zmiany między bieżącym a wybranym commit jako zmiany zatwierdzone. |
| `` gh `` | Twardy reset | Resetuj HEAD do wybranego commita, odrzucając wszystkie zmiany między bieżącym a wybranym commit, jak również wszystkie bieżące modyfikacje w drzewie roboczym. |
| `` <ctrl+t> `` | Otwórz zewnętrzne narzędzie różnic (git difftool) |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | Pokaż commity |  |
| `` w `` | Zobacz opcje drzewa pracy |  |
| `` / `` | Filtruj bieżący widok po tekście |  |

## Zdalne

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Wyświetl gałęzie |  |
| `` n `` | Nowy zdalny |  |
| `` d `` | Usuń | Usuń wybrany zdalny. Wszelkie lokalne gałęzie śledzące gałąź zdalną z tego zdalnego nie zostaną dotknięte. |
| `` e `` | Edytuj | Edytuj nazwę lub URL wybranego zdalnego. |
| `` f `` | Pobierz | Pobierz aktualizacje z zdalnego repozytorium. Pobiera nowe commity i gałęzie bez scalania ich z lokalnymi gałęziami. |
| `` F `` | Add fork remote | Quickly add a fork remote by replacing the owner in the origin URL and optionally check out a branch from new remote. |
| `` / `` | Filtruj bieżący widok po tekście |  |

## Zdalne gałęzie

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | Kopiuj nazwę gałęzi do schowka |  |
| `` <space> `` | Przełącz | Przełącz na nową lokalną gałąź na podstawie wybranej gałęzi zdalnej. Nowa gałąź będzie śledzić gałąź zdalną. |
| `` n `` | Nowa gałąź |  |
| `` M `` | Merge |  |
| `` Mm `` | Scal | Scal wybraną gałąź z aktualnie sprawdzoną gałęzią. |
| `` Mn `` | Non-fast-forward merge | Merge '{{.selectedBranch}}' into '{{.checkedOutBranch}}', creating a merge commit. |
| `` Mf `` | Fast-forward only merge | Fast-forward '{{.checkedOutBranch}}' to '{{.selectedBranch}}' without creating a merge commit. |
| `` Ms `` | Squash merge (uncommitted) | Squash merge '{{.selectedBranch}}' into the working tree. |
| `` MS `` | Squash merge (committed) | Squash merge '{{.selectedBranch}}' into '{{.checkedOutBranch}}' as a single commit. |
| `` r `` | Rebase options |  |
| `` rs `` | Przebazuj | Przebazuj przełączoną gałąź na wybraną gałąź. |
| `` ri `` | Interactive rebase | Rozpocznij interaktywny rebase z przerwaniem na początku, abyś mógł zaktualizować commity TODO przed kontynuacją. |
| `` rb `` | Rebase onto base branch | Rebase the checked out branch onto its base branch (i.e. the closest main branch). |
| `` d `` | Usuń | Usuń gałąź zdalną ze zdalnego. |
| `` us `` | Ustaw jako upstream | Ustaw wybraną gałąź zdalną jako upstream sprawdzonej gałęzi. |
| `` s `` | Kolejność sortowania |  |
| `` g `` | Reset to ref |  |
| `` gm `` | Mixed reset | Resetuj HEAD do wybranego commita, zachowując zmiany między bieżącym a wybranym commit jako zmiany niezatwierdzone. |
| `` gs `` | Miękki reset | Resetuj HEAD do wybranego commita, zachowując zmiany między bieżącym a wybranym commit jako zmiany zatwierdzone. |
| `` gh `` | Twardy reset | Resetuj HEAD do wybranego commita, odrzucając wszystkie zmiany między bieżącym a wybranym commit, jak również wszystkie bieżące modyfikacje w drzewie roboczym. |
| `` <ctrl+t> `` | Otwórz zewnętrzne narzędzie różnic (git difftool) |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | Pokaż commity |  |
| `` w `` | Zobacz opcje drzewa pracy |  |
| `` / `` | Filtruj bieżący widok po tekście |  |
