from bs4 import BeautifulSoup
from urllib.parse import urljoin


async def extract_urls_from_soup(soup: BeautifulSoup, base_url: str) -> set[str]:
    """
    BeautifulSoup 객체에서 URL을 추출하는 함수
    :param soup: BeautifulSoup 객체
    :return: 추출된 URL 목록
    """
    urls = set()

    for link in soup.find_all("a", href=True):
        href = link["href"]

        # 앵커, JavaScript 무시
        if href.startswith("#") or href.startswith("javascript") or not href.strip():
            continue

        full_url = urljoin(f"https://{base_url}", href)

        urls.add(full_url)

    for link in soup.find_all("img"):
        href = link["src"]

        # 앵커, JavaScript 무시
        if href.startswith("#") or href.startswith("javascript") or not href.strip():
            continue

        full_url = urljoin(f"https://{base_url}", href)

        urls.add(full_url)

    return urls
