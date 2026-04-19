from typing import Any, Callable, Mapping, Optional, TypedDict


NotificationPayload = Mapping[str, Any]
NotificationHandler = Callable[[NotificationPayload], dict[str, Any]]


class DispatchResult(TypedDict):
    channel: str
    handled_by: str
    output: dict[str, Any]


def _require_string(payload: NotificationPayload, field_name: str) -> str:
    value = payload.get(field_name)
    if not isinstance(value, str) or not value.strip():
        raise ValueError(f"payload field {field_name!r} must be a non-empty string")
    return value.strip()


def _normalize_channel_name(channel: Any) -> str:
    if not isinstance(channel, str) or not channel.strip():
        raise ValueError("channel must be a non-empty string")
    return channel.strip().lower()


def _email_handler(payload: NotificationPayload) -> dict[str, Any]:
    recipient = _require_string(payload, "recipient")
    message = _require_string(payload, "message")
    subject = str(payload.get("subject", "Notification")).strip() or "Notification"
    return {
        "recipient": recipient,
        "subject": subject,
        "body": message,
    }


def _sms_handler(payload: NotificationPayload) -> dict[str, Any]:
    recipient = _require_string(payload, "recipient")
    message = _require_string(payload, "message")
    return {
        "recipient": recipient,
        "text": message,
    }


def _webhook_handler(payload: NotificationPayload) -> dict[str, Any]:
    url = _require_string(payload, "url")
    message = _require_string(payload, "message")
    return {
        "url": url,
        "method": "POST",
        "json": {"message": message},
    }


def _normalize_handler_registry(
    handlers: Mapping[str, NotificationHandler],
) -> dict[str, NotificationHandler]:
    normalized_handlers = {}
    for channel_name, handler in handlers.items():
        normalized_handlers[_normalize_channel_name(channel_name)] = handler
    return normalized_handlers


DEFAULT_NOTIFICATION_HANDLERS: dict[str, NotificationHandler] = {
    "email": _email_handler,
    "sms": _sms_handler,
    "webhook": _webhook_handler,
}


def dispatch_notification(
    channel: str,
    payload: NotificationPayload,
    handlers: Optional[Mapping[str, NotificationHandler]] = None,
) -> DispatchResult:
    """Dispatch a notification through a handler registry instead of if/elif branches.

    Args:
        channel: Notification channel name such as email, sms, or webhook.
        payload: Input payload for the selected handler.
        handlers: Optional custom registry that overrides the default handlers.

    Returns:
        A dictionary describing the normalized channel and handler output.

    Raises:
        ValueError: If the channel is blank or unsupported.
    """
    normalized_channel = _normalize_channel_name(channel)

    handler_registry = dict(DEFAULT_NOTIFICATION_HANDLERS)
    if handlers is not None:
        handler_registry.update(_normalize_handler_registry(handlers))

    handler = handler_registry.get(normalized_channel)
    if handler is None:
        raise ValueError(f"Unsupported channel: {channel}")

    return {
        "channel": normalized_channel,
        "handled_by": normalized_channel,
        "output": handler(dict(payload)),
    }