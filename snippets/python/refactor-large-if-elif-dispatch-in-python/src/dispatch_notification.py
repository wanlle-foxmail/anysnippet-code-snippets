from typing import Callable


def send_email(payload: dict[str, str]) -> dict[str, object]:
    return {
        "recipient": payload["recipient"],
        "subject": payload.get("subject", "Notification"),
        "body": payload["message"],
    }


def send_sms(payload: dict[str, str]) -> dict[str, object]:
    return {
        "recipient": payload["recipient"],
        "text": payload["message"],
    }


def send_webhook(payload: dict[str, str]) -> dict[str, object]:
    return {
        "url": payload["url"],
        "method": "POST",
        "json": {"message": payload["message"]},
    }


NOTIFICATION_HANDLERS: dict[str, Callable[[dict[str, str]], dict[str, object]]] = {
    "email": send_email,
    "sms": send_sms,
    "webhook": send_webhook,
}


def dispatch_notification(channel: str, payload: dict[str, str]) -> dict[str, object]:
    """Replace if/elif dispatch with a handler mapping."""
    handler = NOTIFICATION_HANDLERS.get(channel)
    if handler is None:
        raise ValueError(f"Unsupported channel: {channel}")

    return handler(payload)


if __name__ == "__main__":
    result = dispatch_notification(
        "email",
        {
            "recipient": "user@example.com",
            "message": "Your message",
        },
    )
    print(result)