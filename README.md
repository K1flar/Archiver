# Archiver

Утилита для Ubuntu, позволяющая архивировать и разархивировать файлы.

## Описание

1. Архиватор принимает папку с файлами, которую сжимает в один файл-архив, объединяющий все файлы из входной директории, а также содержащий заголовок с информацией об именах файлов и их размерах.

2. Разархиватор принимает архив, разделяющий архивированный файл обратно в файлы.

## Поддержка флагов

``-o [path]`` - имя (или путь) выходного архива или директория для разархивации 

``-bwt`` - преобразование Барроуза — Уилера (Burrows-Wheeler transform)

``-mtf`` - преобразование MTF (move-to-front)

``-huff`` - сжатие алгоритмом Хаффмана

``-mirror`` - зеркальное отражение всех символов в архиве

``-caesar [shift]`` - закодировать архив шифром Цезаря с указанным сдвигом (если сдвиг не указан - по умолчанию 0) 