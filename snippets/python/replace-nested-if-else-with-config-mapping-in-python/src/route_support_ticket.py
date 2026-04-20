DEFAULT_ROUTE = {
    "queue": "general-support",
    "priority": "normal",
    "sla_hours": 24,
}


ROUTING_RULES = {
    ("chat", "enterprise", "billing"): {
        "queue": "priority-billing",
        "priority": "urgent",
        "sla_hours": 1,
    },
    ("email", "pro", "login"): {
        "queue": "pro-login",
        "priority": "high",
        "sla_hours": 4,
    },
    ("phone", "enterprise", "security"): {
        "queue": "security-incident",
        "priority": "urgent",
        "sla_hours": 1,
    },
}


def route_support_ticket(channel: str, customer_tier: str, issue_type: str) -> dict[str, object]:
    """Replace nested if/else logic with a config mapping."""
    route = ROUTING_RULES.get((channel, customer_tier, issue_type), DEFAULT_ROUTE)
    return dict(route)


if __name__ == "__main__":
    result = route_support_ticket("email", "pro", "login")
    print(result)