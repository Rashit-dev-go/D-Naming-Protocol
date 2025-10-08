# 🛰 D Naming Protocol v3

Полностью автономный CLI-инструмент, который:

* создаёт проект в заданной папке,
* генерирует `go.mod` , `Makefile` , `README` ,
* инициализирует `git` ,
* ведёт централизированный журнал всех ваших миссий,
* **автоматически создаёт репозиторий на GitHub и пушит первый коммит**.

---

## ⚙️ Возможности:

* `--prefix`  — задать имя вручную
* `--desc`  — описание проекта
* `--dir`  — временно переопределить корневую папку
* автогенерация структуры проекта
* автоинициализация `git`  и `go.mod`
* **автоматическое создание GitHub репозитория** (с OAuth-токеном)
* конфигурация через `~/.dnp/config.yaml`

---

## 📦 Установка

1. Создай папку для исходника:

   ```bash
   mkdir -p ~/dev/dnp && cd ~/dev/dnp
   ```
2. Собери и установи:

   ```bash
   go build -o dnp main.go
   sudo mv dnp /usr/local/bin/
   ```
3. Создай конфиг:

   ```bash
   mkdir -p ~/.dnp
   cat > ~/.dnp/config.yaml <<EOF
   root_dir: /home/$USER/Projects/D
   default_type: LAB
   default_domain: CORE
   git_init: true
   github_token: "your_github_personal_access_token"
   EOF
   ```

   **Получите токен GitHub:**
   - Перейдите в Settings → Developer settings → Personal access tokens
   - Создайте токен с областью `repo`
   - Вставьте его в `github_token`

---

## 💻 Пример использования

```bash
dnp create proto billing --prefix=ARGUS --desc="SaaS биллинг с интеграцией Stripe"
```

→ результат:

```
Создан проект: ARGUS-PROTO-BILLING
Расположение: ~/Projects/D/argus-proto-billing/
Git инициализирован: ✅
GitHub репозиторий создан: https://github.com/username/argus-proto-billing
Первый коммит запушен на GitHub: ✅
```

---

## 🧱 Структура созданного проекта

```
~/Projects/D/argus-proto-billing/
├── cmd/
│   └── main.go
├── internal/
│   └── core/
├── go.mod
├── Makefile
├── README.md
├── .gitignore
└── .git/ (с remote origin на GitHub)
```

---

## 🚀 Версия 3: GitHub Integration

DNP v3 теперь поддерживает автоматическое развертывание на GitHub:

- Создание публичного репозитория под вашей учётной записью
- Установка origin и пуш первого коммита
- Интеграция через Personal Access Token

**Без токена:** Работает как v2 (только локально)

**С токеном:** Полная автоматизация от создания до деплоя!

---

## 📜 Команды

- `dnp create [type] [domain] [--prefix=NAME] [--desc='описание'] [--dir=/путь]` — создать проект
- `dnp list` — показать журнал всех созданных проектов

---

## 🔧 Конфигурация

Файл `~/.dnp/config.yaml`:

```yaml
root_dir: /home/user/Projects/D      # Корневая папка для проектов
default_type: LAB                        # Тип по умолчанию
default_domain: CORE                     # Домен по умолчанию
git_init: true                          # Инициализировать Git
github_token: "ghp_..."                 # Токен GitHub (опционально)
```

---

## 📊 Журнал миссий

Все созданные проекты логируются в `~/.dnp/projects.log`:

```
2025-10-08 13:04 — ARGUS-PROTO-BILLING (/home/user/Projects/D/argus-proto-billing)
```

---

## 🛠 Сборка

```bash
go mod tidy
go build -o dnp main.go
sudo mv dnp /usr/local/bin/
```
