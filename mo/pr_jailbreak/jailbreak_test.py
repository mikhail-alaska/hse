import json
import requests
import csv
import time
import sys

OLLAMA_URL = "http://localhost:11434/api/generate"
PROMPTS_FILE = "prompts.json"

def load_prompts(path):
    with open(path, "r", encoding="utf-8") as f:
        return json.load(f)

def ask(model, prompt, attack_type):
    if attack_type!="Multi-turn":
        r = requests.post(
            OLLAMA_URL,
            json={
                "model": model,
                "prompt": prompt,
                "stream": False
            },
            timeout=180
        )
        r.raise_for_status()
    else:
        for elem in prompt:

            r = requests.post(
                OLLAMA_URL,
                json={
                    "model": model,
                    "prompt": prompt,
                    "stream": False
                },
                timeout=180
            )
            r.raise_for_status()

    return r.json()["response"]



def main():
    prompts_data = load_prompts(PROMPTS_FILE)

    with open("results.csv", "w", newline="", encoding="utf-8") as f:
        writer = csv.writer(f)
        writer.writerow([
            "Модель",
            "Тип атаки",
            "Номер теста",
            "Промпт",
            "Ответ",
            "Успешен",
            "Комментарий"
        ])

        for model, attacks in prompts_data.items():
            for attack_type, prompts in attacks.items():
                for test_id, prompt in enumerate(prompts, start=1):

                    print(f"[{model}] {attack_type} — тест {test_id}")

                    response = ask(model, prompt, attack_type)
                    success = 0

                    comment = (
                        "Полностью выполнила"
                        if success else
                        "Отказала / уклонилась"
                    )

                    writer.writerow([
                        model,
                        attack_type,
                        test_id,
                        prompt,
                        response,
                        success,
                        comment
                    ])

                    time.sleep(1)

    print("\nГотово: results.csv создан.")

if __name__ == "__main__":
    try:
        main()
    except Exception as e:
        print(f"Ошибка: {e}", file=sys.stderr)
        sys.exit(1)
