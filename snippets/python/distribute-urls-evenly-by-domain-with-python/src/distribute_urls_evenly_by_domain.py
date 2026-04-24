from collections import deque
from typing import Optional


def distribute_urls_evenly_by_domain(url_to_domain: dict[str, str]) -> list[str]:
    """Return URLs ordered to reduce same-domain clustering in crawl queues."""
    # Flow:
    #   URL map -> validate items and group URLs by domain
    #              |
    #              +-> track remaining counts and previous domain
    #                   |
    #                   +-> pick the best next domain each round
    #                        |
    #                        +-> emit one URL and repeat until all queues are drained
    if not isinstance(url_to_domain, dict):
        raise TypeError("url_to_domain must be a dictionary of URL-to-domain mappings")

    grouped_urls: dict[str, deque[str]] = {}
    domain_order: dict[str, int] = {}

    for url, domain in url_to_domain.items():
        if not isinstance(url, str):
            raise TypeError("each URL must be a string")
        if not isinstance(domain, str):
            raise TypeError("each domain must be a string")

        if domain not in grouped_urls:
            grouped_urls[domain] = deque()
            domain_order[domain] = len(domain_order)

        grouped_urls[domain].append(url)

    remaining_by_domain = {
        domain: len(urls)
        for domain, urls in grouped_urls.items()
    }
    ordered_urls: list[str] = []
    previous_domain: Optional[str] = None

    while remaining_by_domain:
        available_domains = [
            domain for domain in remaining_by_domain if domain != previous_domain
        ]
        if not available_domains:
            available_domains = list(remaining_by_domain)

        next_domain = min(
            available_domains,
            key=lambda domain: (-remaining_by_domain[domain], domain_order[domain]),
        )

        ordered_urls.append(grouped_urls[next_domain].popleft())
        remaining_by_domain[next_domain] -= 1
        if remaining_by_domain[next_domain] == 0:
            del remaining_by_domain[next_domain]

        previous_domain = next_domain

    return ordered_urls


if __name__ == "__main__":
    result = distribute_urls_evenly_by_domain(
        {
            "https://www.reuters.com/world/article-001/": "www.reuters.com",
            "https://www.reuters.com/business/article-002/": "www.reuters.com",
            "https://www.reuters.com/technology/article-003/": "www.reuters.com",
            "https://www.bbc.com/news/articles/article-101": "www.bbc.com",
            "https://www.bbc.com/news/articles/article-102": "www.bbc.com",
            "https://www.bbc.com/news/articles/article-103": "www.bbc.com",
            "https://www.bbc.com/news/articles/article-104": "www.bbc.com",
            "https://www.afp.com/en/news/article-201": "www.afp.com",
            "https://www.afp.com/en/news/article-202": "www.afp.com",
        }
    )
    print(result)