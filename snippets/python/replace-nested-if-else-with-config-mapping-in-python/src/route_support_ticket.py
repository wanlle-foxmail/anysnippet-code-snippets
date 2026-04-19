from typing import TypedDict


class RoutingDecision(TypedDict):
    normalized_channel: str
    normalized_customer_tier: str
    normalized_issue_type: str
    queue: str
    priority: str
    sla_hours: int


ROUTING_RULES = {
    ("chat", "enterprise", "billing"): {
        "queue": "priority-billing",
        "priority": "urgent",
        "sla_hours": 1,
    },
    ("email", "pro", "*"): {
        "queue": "pro-email",
        "priority": "high",
        "sla_hours": 4,
    },
    ("*", "enterprise", "security"): {
        "queue": "security-incident",
        "priority": "urgent",
        "sla_hours": 1,
    },
    ("*", "*", "*"): {
        "queue": "general-support",
        "priority": "normal",
        "sla_hours": 24,
    },
}


def route_support_ticket(
    channel: str,
    customer_tier: str,
    issue_type: str,
) -> RoutingDecision:
    """Route support tickets with configuration rules instead of nested if/else blocks.

    Args:
        channel: Ticket source such as chat, email, or phone.
        customer_tier: Customer tier such as starter, pro, or enterprise.
        issue_type: Issue category such as billing or security.

    Returns:
        A dictionary containing normalized inputs and the selected routing rule.

    Raises:
        ValueError: If any input value is blank.
    """
    normalized_channel = _normalize_value(channel, "channel")
    normalized_customer_tier = _normalize_value(customer_tier, "customer_tier")
    normalized_issue_type = _normalize_value(issue_type, "issue_type")

    for key in _candidate_keys(normalized_channel, normalized_customer_tier, normalized_issue_type):
        rule = ROUTING_RULES.get(key)
        if rule is not None:
            return {
                "normalized_channel": normalized_channel,
                "normalized_customer_tier": normalized_customer_tier,
                "normalized_issue_type": normalized_issue_type,
                "queue": rule["queue"],
                "priority": rule["priority"],
                "sla_hours": rule["sla_hours"],
            }

    raise ValueError("No routing rule matched the provided inputs")


def _normalize_value(value: str, field_name: str) -> str:
    if not isinstance(value, str) or not value.strip():
        raise ValueError(f"{field_name} must be a non-empty string")
    return value.strip().lower()


def _candidate_keys(
    channel: str,
    customer_tier: str,
    issue_type: str,
) -> list[tuple[str, str, str]]:
    return [
        (channel, customer_tier, issue_type),
        (channel, customer_tier, "*"),
        (channel, "*", issue_type),
        (channel, "*", "*"),
        ("*", customer_tier, issue_type),
        ("*", customer_tier, "*"),
        ("*", "*", issue_type),
        ("*", "*", "*"),
    ]