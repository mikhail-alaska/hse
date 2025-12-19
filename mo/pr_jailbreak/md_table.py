import csv
from collections import defaultdict

FILES = [
    "results_pc.csv",
   # "results_laptop.csv"
]

def load_rows(files):
    rows = []
    for file in files:
        with open(file, newline="", encoding="utf-8") as f:
            reader = csv.DictReader(f)
            rows.extend(reader)
    return rows

def main():
    rows = load_rows(FILES)

    # агрегаты
    model_total = defaultdict(int)
    model_success = defaultdict(int)

    model_attack_total = defaultdict(lambda: defaultdict(int))
    model_attack_success = defaultdict(lambda: defaultdict(int))

    attack_types = set()
    models = set()

    for row in rows:
        model = row["Модель"]
        attack = row["Тип атаки"]
        success = row["Успешен"].strip() == "1"

        models.add(model)
        attack_types.add(attack)

        model_total[model] += 1
        model_attack_total[model][attack] += 1

        if success:
            model_success[model] += 1
            model_attack_success[model][attack] += 1

    models = sorted(models)
    attack_types = sorted(attack_types)

    # ===== ОБЩАЯ СТАТИСТИКА =====
    print("\n## 📊 Общая статистика по моделям\n")
    print("| Модель | Успешных | Всего | % успеха |")
    print("|-------|----------|-------|----------|")

    for model in models:
        total = model_total[model]
        success = model_success[model]
        percent = (success / total * 100) if total else 0
        print(f"| `{model}` | {success} | {total} | {percent:.1f}% |")

    # ===== ТАБЛИЦА ПО ТИПАМ АТАК =====
    print("\n## 🧪 Процент успешных джейлбрейков по типам атак\n")

    header = "| Модель | " + " | ".join(attack_types) + " |"
    separator = "|---|" + "|".join(["---"] * len(attack_types)) + "|"

    print(header)
    print(separator)

    for model in models:
        row = [f"`{model}`"]
        for attack in attack_types:
            total = model_attack_total[model][attack]
            success = model_attack_success[model][attack]
            if total == 0:
                cell = "—"
            else:
                cell = f"{(success / total * 100):.1f}%"
            row.append(cell)
        print("| " + " | ".join(row) + " |")

if __name__ == "__main__":
    main()
