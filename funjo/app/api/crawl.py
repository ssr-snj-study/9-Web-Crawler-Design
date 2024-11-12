from flask import request, jsonify
from app.main import main
import requests
from bs4 import BeautifulSoup
from urllib.parse import urlparse
import os
from selenium import webdriver
from selenium.webdriver.chrome.service import Service
from webdriver_manager.chrome import ChromeDriverManager
import time
from concurrent.futures import ThreadPoolExecutor

@main.route('/start_crawl', methods=['POST'])
def start_crawl():
    data = request.get_json()
    start_url = data.get('url')  # 단일 URL을 받음
    storage_directory = "crawled_contents"  # 콘텐츠 저장 디렉토리
    depth = data.get('depth', 1)  # 깊이 설정 (기본값 1)
    crawl([start_url], storage_directory, depth)  # 시작 URL을 리스트로 변환하여 크롤링 함수 호출
    return jsonify({"message": "Crawling started.", "url": start_url}), 200


# 미수집 URL 저장소
class URLQueue:
    def __init__(self):
        self.queue = set()  # 중복을 피하기 위해 set 사용

    def add_url(self, url, depth):
        self.queue.add((url, depth))  # URL과 깊이를 함께 저장

    def get_next_url(self):
        return self.queue.pop() if self.queue else None


# HTML 다운로더
class HTMLDownloader:
    @staticmethod
    def download(url):
        try:
            # Selenium WebDriver 설정
            options = webdriver.ChromeOptions()
            options.add_argument('--headless')  # 헤드리스 모드로 실행
            driver = webdriver.Chrome(service=Service(ChromeDriverManager().install()), options=options)

            driver.get(url)  # URL 열기
            time.sleep(3)  # 페이지가 완전히 로드될 때까지 대기

            html_content = driver.page_source  # 페이지 소스 가져오기
            driver.quit()  # 드라이버 종료

            return html_content
        except Exception as e:
            print(f"Error downloading {url}: {e}")
            return None


# 도메인 이름 변환기
class DomainNameConverter:
    @staticmethod
    def get_domain(url):
        parsed_url = urlparse(url)
        return parsed_url.netloc


# 콘텐츠 파서
class ContentParser:
    @staticmethod
    def parse(html_content):
        soup = BeautifulSoup(html_content, 'html.parser')
        return soup


# 중복 확인
class DuplicateChecker:
    def __init__(self):
        self.seen_content = set()

    def is_duplicate(self, content):
        content_hash = hash(content)  # 콘텐츠의 해시를 사용하여 중복 확인
        if content_hash in self.seen_content:
            return True
        self.seen_content.add(content_hash)
        return False


# 콘텐츠 저장소
class ContentStorage:
    def __init__(self, directory):
        self.directory = directory
        os.makedirs(directory, exist_ok=True)  # 저장소 디렉토리 생성

    def save_content(self, domain, content):
        filename = os.path.join(self.directory, f"{domain}.html")
        with open(filename, 'w', encoding='utf-8') as f:
            f.write(content)
        print(f"Content saved for domain: {domain}")


# URL 추출기
class URLExtractor:
    @staticmethod
    def extract_urls(soup, base_url):
        urls = []
        for a_tag in soup.find_all('a', href=True):
            url = a_tag['href']
            # 절대 URL로 변환
            if not url.startswith(('http://', 'https://')):
                url = urlparse(base_url)._replace(path=url).geturl()
            urls.append(url)
        return urls


# 크롤링 프로세스
def crawl(start_urls, storage_directory, max_depth):
    url_queue = URLQueue()
    for url in start_urls:
        url_queue.add_url(url, 0)  # 깊이 0으로 시작

    downloader = HTMLDownloader()
    domain_converter = DomainNameConverter()
    parser = ContentParser()
    checker = DuplicateChecker()
    storage = ContentStorage(storage_directory)

    with ThreadPoolExecutor(max_workers=5) as executor:
        while url_queue.queue:
            url, depth = url_queue.get_next_url()
            print(f"Crawling: {url} at depth: {depth}")

            # 현재 깊이가 최대 깊이에 도달했으면 크롤링 중단
            if depth >= max_depth:
                continue

            # HTML 다운로드를 비동기적으로 실행
            future = executor.submit(downloader.download, url)
            html_content = future.result()  # 결과를 기다림

            if html_content is None:
                continue  # 다운로드 실패 시 다음 URL로 진행

            # 도메인 변환
            domain = domain_converter.get_domain(url)

            # 콘텐츠 파싱
            soup = parser.parse(html_content)

            # 중복 확인
            if checker.is_duplicate(html_content):
                print(f"Duplicate content found for URL: {url}")
                continue

            # 콘텐츠 저장
            storage.save_content(domain, html_content)

            # URL 추출
            extracted_urls = URLExtractor.extract_urls(soup, url)  # base_url로 현재 URL 전달
            for extracted_url in extracted_urls:
                # 중복된 URL은 추가하지 않음
                url_queue.add_url(extracted_url, depth + 1)  # 깊이를 1 증가시켜 추가
