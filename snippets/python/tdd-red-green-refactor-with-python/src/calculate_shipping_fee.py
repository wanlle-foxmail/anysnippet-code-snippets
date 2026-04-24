FREE_SHIPPING_THRESHOLD_CENTS = 5000
STANDARD_SHIPPING_FEE_CENTS = 500


def calculate_shipping_fee(subtotal_cents: int, customer_tier: str = "standard") -> int:
    """Return the shipping fee in cents for one order."""
    # Flow: validate subtotal -> validate tier -> premium stays free -> threshold decides standard fee
    if isinstance(subtotal_cents, bool):
        raise TypeError("subtotal_cents must be an integer")
    if not isinstance(subtotal_cents, int):
        raise TypeError("subtotal_cents must be an integer")
    if subtotal_cents < 0:
        raise ValueError("subtotal_cents must be greater than or equal to 0")
    if customer_tier not in {"standard", "premium"}:
        raise ValueError("customer_tier must be 'standard' or 'premium'")

    if customer_tier == "premium":
        return 0
    if subtotal_cents >= FREE_SHIPPING_THRESHOLD_CENTS:
        return 0
    return STANDARD_SHIPPING_FEE_CENTS


if __name__ == "__main__":
    print(calculate_shipping_fee(3200))