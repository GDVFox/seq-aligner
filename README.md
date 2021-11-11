# 🧪 Sequence Aligner

## Дисклеймер

Sequence Aligner был написан мной, Гавриловским Даниилом, на курсе «Алгоритмы Биоинформатики» в 7-м семестре на кафедре ИУ9 МГТУ им. Н.Э. Баумана в 2020 году. Использовал этот Sequence Aligner для решения задач на ROSALIND тоже я.

## Описание

Производит оптимальное глобальное (или локальное с флагом `--local`) парное выравнивание двух последовательностей, получает для него оценку в скоринговой системе. Использует алгоритм [Нидлмана—Вунша](https://en.wikipedia.org/wiki/Needleman%E2%80%93Wunsch_algorithm).

## Build

```bash
mkdir _build && go build -o _build/seq-aligner *.go
```

## Run

```bash
cd _build && ./seq-aligner <flag_options> <your_fasta_file> [<your_second_fasta_file>]
```

### Входные данные

1. Один файл `your_fasta_file` с двумя последовательностями в формате fasta. Из него для выравнивания будут загружены две первые последовательности.
2. Два файла `your_fasta_file` и `our_second_fasta_file`. Первая последовательность будет взята из первого файла, вторая из второго.

### Доступные опции

Опции передаются как флаги командной строки _перед_ входными файлами.

| Опция   |      Тип      |  По-умолчанию | Описание |
|----------|---------------|-------|-------|
| `--gap` | int | -2 | цена установки `-` в скоринговой системе |
| `--gap-open` | int | -2 | цена установки первого (следующего 1-м символом строки или после буквы) `-` в скоринговой системе |
| `--gap-extend` | int | 0 | цена установки новых `-` следующих за существующими `-` в скоринговой системе. Если флаг не передан, то всегда используется значение параметра `--gap` |
| `--mode` | dna\|protein_b62\|protein_p250\|default | default | выбор [алфавита и скоринга](#алфавиты) |
| `--pretty` | bool | false | вывод в `🦄🌈⭐красивом режиме⭐🌈🦄` |
| `--mem-save` | bool | false | эффективный по памяти режим работы с незначительными ограничениями |
| `--line` | int | 100 | количество символов последовательности в одной строке |
| `--out` | string |  | имя файла для вывода, если не указано, вывод в консоль |
| `--spen` | bool | false | штрафовать за `-` в _начале_ последовательности |
| `--epen` | bool | false | штрафовать за `-` в _конце_ последовательности |
| `--local` | bool | false | работать в режиме локального выравнивания |

### Алфавиты

На данные момент поддерживаются:

* DNA (`--mode=dna`): последовательности нуклеотидов. Алфавит состоит из символов `{A,T,G,C}`. Для скоринга используется матрица [DNAFull](http://rosalind.info/glossary/dnafull/).
* Protein (`--mode=protein_b62` и `--mode=protein_p250`): последовательности аминокислот. Алфавит состоит из символов `{A,R,N,D,C,Q,E,G,H,I,L,K,M,F,P,S,T,W,Y,V}`. Для скоринга используется матрица [BLOSUM62](https://www.ncbi.nlm.nih.gov/Class/BLAST/BLOSUM62.txt) или [PAM250](https://www.ncbi.nlm.nih.gov/IEB/ToolBox/C_DOC/lxr/source/data/PAM250) в зависимости от указанного режима.
* Произвольный (`--mode=default`): произвольные последовательности. Алфавит состоит из всеъ символов, кроме `-`. Для скоринга используется правило: совпадение символов — `+1`, несовпадение символов — `-1`.
