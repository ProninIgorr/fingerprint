# fingerprint 
- [Кейс 2 - Распознавание отпечатков пальцев (Go)](https://www.zavodit.ru/ru/profile/hackathons-participant/case/16)

# Описание кейса
- [__Памятка участникам__](https://drive.google.com/file/d/13-EpP_g8ppJ6VYVP0g4AoQaDiEJU6R1S/view?usp=share_link)
- __Презентация кейса №2__
  - [pdf](https://drive.google.com/file/d/1EXuFBBhypP_CjIvAzKugvxSdcDJqs5am/view?usp=share_link)
  - [mkv](https://drive.google.com/file/d/1zJUwMx-o5AiOcnEUD62aqpx3DPCD2EOA/view?usp=share_link)
- __Целевая аудитория__: 
  - Разработчики Go
- __Задача__: 
  - Реализовать алгоритм сравнения отпечатков пальцев на языке Go. 
  - На вход передается отпечаток с разрешением 500 dpi размером 440x500. 
  - Отпечаток сравнивается с базой отпечатков, составленной из датасета SOCOFing.
    - Данный датасет состоит из 6000 изображений от 600 личностей, по 1 отпечатку для каждого пальца. 
    - Также в датасете имеются аугментации изображений разной силы.
  - По входному изображению (не из базы) требуется найти данный отпечаток в базе. 
  - Для тестирования точности работы алгоритма будет выдан тестовый сет изображений отпечатков с метками. 
  - Необходимо посчитать метрикy качества работы алгоритма, а именно точность (accuracy).
  - Алгоритм должен масштабироваться горизонтально (увеличение числа ядер процессора), 
    - максимальное время сравнения с образцом - 2 секунды. 
  - Необходимо предоставить результаты времени работы алгоритма в зависимости от используемого числа ядер (потоков).

- __Dataset__:
  - source: 
    - working [Датасет.zip](https://drive.google.com/file/d/1RoPUZ9WivSKc6cvVDCp3z_iLEkGX4ObM/view?usp=share_link)
    - original [SOCOFing](https://www.kaggle.com/datasets/ruizgara/socofing)
  - archive: 
    - [zip full original](https://drive.google.com/file/d/1FNjNfDlFAdQn2gM_w5XBBvslYWiR0Mev/view?usp=sharing)
    - [zip real head 100](https://drive.google.com/file/d/1-0MdvDUjwWeYoQ7uH7xceyYvSGIoe6tU/view?usp=share_link)
