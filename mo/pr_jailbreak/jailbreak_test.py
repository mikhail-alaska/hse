import json
import requests
import csv
import time
import sys

GENERATE_URL = "http://localhost:11434/api/generate"
CHAT_URL = "http://localhost:11434/api/chat"
PROMPTS_FILE = "prompts.json"


def load_prompts(path):
    with open(path, "r", encoding="utf-8") as f:
        return json.load(f)


def normalize_csv_field(text):
    return (
        text
        .replace("\\", "\\\\")
        .replace("\r\n", "\\n")
        .replace("\n", "\\n")
        .replace("\r", "\\n")
    )


def ask_single(model, prompt):
    r = requests.post(
        GENERATE_URL,
        json={
            "model": model,
            "prompt": prompt,
            "stream": False,
            "options": {
                "num_predict": 256
            }
        },
        timeout=900
    )
    r.raise_for_status()
    return r.json()["response"]


def ask_multiturn(model, messages):
    chat = []
    last_response = ""

    for msg in messages:
        chat.append({
            "role": "user",
            "content": msg
        })

        r = requests.post(
            CHAT_URL,
            json={
                "model": model,
                "messages": chat,
                "stream": False,
                "options": {
                    "num_predict": 256
                }
            },
            timeout=900
        )
        r.raise_for_status()

        last_response = r.json()["message"]["content"]

        chat.append({
            "role": "assistant",
            "content": last_response
        })

    return last_response


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

                print(f"[{model}] {attack_type} - тест {test_id}")

                if attack_type == "Multi-turn":
                    response = ask_multiturn(model, prompt)
                    prompt_for_csv = " → ".join(prompt)
                else:
                    response = ask_single(model, prompt)
                    prompt_for_csv = prompt

                # временная заглушка, позже вручную менял значения
                success = 0
                comment = "Отказала / уклонилась"

                writer.writerow([
                    model,
                    attack_type,
                    test_id,
                    normalize_csv_field(prompt_for_csv),
                    normalize_csv_field(response),
                    success,
                    comment
                ])

                time.sleep(1)



