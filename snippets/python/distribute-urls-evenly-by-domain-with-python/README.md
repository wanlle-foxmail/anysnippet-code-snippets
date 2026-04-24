# Distribute URLs Evenly by Domain with Python

Reorder crawl URLs so requests stay spread across domains instead of clustering on one site.

This snippet is useful when a crawl queue contains many news detail URLs and you want consecutive requests to stay spread across domains so proxy pressure grows more slowly.

## Highlights

- Spreads requests across domains
- Keeps per-domain URL order
- Keeps exact domain strings

## Use Cases

- Spread news detail page requests across domains
- Reduce bursts against one publisher site
- Build a steadier crawl queue before workers pick URLs

## Code

```python
from src.distribute_urls_evenly_by_domain import distribute_urls_evenly_by_domain


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
```

## Notes

- The input must be a dictionary that maps each URL string to one domain string.
- Domain values are grouped exactly as provided. The snippet does not normalize case, `www.`, or ports.
- URLs from the same domain keep their original relative order.
- When more than one domain is available, the function avoids repeating the previous domain and prefers the domain with more URLs remaining.

## Verification

Run the unit tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- empty input handling
- single-domain passthrough
- balanced cross-domain alternation
- anti-clustering for skewed domains
- preference for domains with more remaining URLs
- exact-string domain grouping
- invalid input type errors

## Files

- `src/distribute_urls_evenly_by_domain.py`
- `tests/test_distribute_urls_evenly_by_domain.py`
- `snippet.json`