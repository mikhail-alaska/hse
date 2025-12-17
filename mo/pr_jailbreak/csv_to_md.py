import csv

INPUT = "results.csv"
OUTPUT = "results.md"

def render(text):
    if not text:
        return ""
    return text.replace("\\n", "\n")

with open(INPUT, newline="", encoding="utf-8") as f, \
     open(OUTPUT, "w", encoding="utf-8") as out:

    reader = csv.DictReader(f)

    for row in reader:
        out.write(f"## 🧪 Тест {row['Номер теста']}\n\n")
        out.write(f"**Модель:** `{row['Модель']}`  \n")
        out.write(f"**Тип атаки:** `{row['Тип атаки']}`  \n")
        out.write(f"**Успех:** `{row['Успешен']}`  \n")
        out.write(f"**Комментарий:** {row['Комментарий']}\n\n")

        out.write("### 🔹 Промпт\n")
        out.write(render(row["Промпт"]))

        out.write("### 🔹 Ответ модели\n")
        out.write(render(row["Ответ"]))

        out.write("---\n\n")

print("Готово: results.md")
