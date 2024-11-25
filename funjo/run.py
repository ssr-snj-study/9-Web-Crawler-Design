from flask import Flask, request, jsonify
from bs4 import BeautifulSoup
import time
from hashlib import sha256
from urllib.parse import urlparse, urlunparse
import psycopg2
from concurrent.futures import ThreadPoolExecutor
import os
from selenium import webdriver  # selenium.webdriver 임포트 추가
from selenium.webdriver.chrome.service import Service  # Service 임포트 추가
from webdriver_manager.chrome import ChromeDriverManager  # ChromeDriverManager 임포트 추가
import heapq
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.by import By

app = Flask(__name__)

# PostgreSQL 연결 설정
DB_CONFIG = {
    "dbname": "funjodb",
    "user": "funjo",
    "password": "funjopass",
    "host": "192.168.169.2",
    "port": "5432"
}
# DB 연결 함수
def get_db_connection():
    conn = psycopg2.connect(
        dbname=DB_CONFIG['dbname'],
        user=DB_CONFIG['user'],
        password=DB_CONFIG['password'],
        host=DB_CONFIG['host'],
        port=DB_CONFIG['port']
    )
    return conn


@app.route('/start_crawl', methods=['POST'])
def start_crawl():
    data = request.get_json()
    start_url = data.get('url')  # URL을 받음
    storage_directory = "crawled_contents"  # 콘텐츠 저장 디렉토리
    depth = data.get('depth', 1)  # 깊이 설정 (기본값 1)

    # 크롤링을 시작
    crawl([start_url], storage_directory, depth)

    return jsonify({"message": "Crawling started.", "url": start_url}), 200

#우선순위
def classify_url(url):
    if "finance" in url or "series" in url:
        return 1
    elif "comic" in url or "m." in url:
        return 2
    elif "weather" in url or "news" in url:
        return 3
    elif "novel" in url or "map" in url:
        return 4
    else:
        return 5

# 미수집 URL 저장소
class URLQueue:
    def __init__(self):
        self.queue = []  # 우선순위 큐로 URL 저장
        self.counter = 0  # 중복 URL을 방지하는 카운터

    def add_url(self, url, depth, priority):
        # 우선순위와 URL을 튜플로 저장
        # 우선순위가 낮을수록 먼저 처리되므로, 우선순위 값이 낮을수록 큐의 앞에 배치됨 classify_url의 리턴값 priority
        heapq.heappush(self.queue, (priority, self.counter, url, depth))
        self.counter += 1

    def get_next_url(self):
        if self.queue:
            # 가장 우선순위가 높은 URL을 반환 (우선순위가 낮은 값이 먼저 나옴)
            return heapq.heappop(self.queue)[2:]
        return None

    def is_empty(self):
        return len(self.queue) == 0


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

            # 페이지 로딩 완료 대기 (동적 콘텐츠가 로드될 때까지)
            WebDriverWait(driver, 10).until(EC.presence_of_element_located((By.TAG_NAME, 'body')))  # 예시로 body 태그가 로드되었을 때까지 대기

            html_content = driver.page_source  # 페이지 소스 가져오기
            driver.quit()  # 드라이버 종료

            # 정규화 처리 추가
            normalized_content = ContentNormalizer.normalize(html_content)
            return normalized_content

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

class ContentNormalizer:
    @staticmethod
    def normalize(html_content):
        soup = BeautifulSoup(html_content, 'html.parser')

        # 동적 태그 제거
        for tag in soup(['script', 'style', 'meta']):
            tag.decompose()

        # 난수 기반 속성 제거 (예: 동적인 ID, 클래스)
        for tag in soup.find_all(True):  # 모든 태그 순회
            for attr in ['id', 'class']:
                if attr in tag.attrs and any(char.isdigit() for char in tag[attr]):
                    del tag.attrs[attr]

        # 정규화된 HTML 반환
        return soup.prettify()


# 중복 확인
class DuplicateChecker:
    def __init__(self):
        self.conn = get_db_connection()

    def url_is_duplicate(self, url):

        domain = DomainNameConverter().get_domain(url)
        # DB에서 해당 해시 값이 있는지 확인
        cursor = self.conn.cursor()
        cursor.execute("SELECT id FROM downloaded_html WHERE domain_name = %s", ( domain, ))
        result = cursor.fetchone()

        if result:
            return True
        else:
            return False

    def html_is_duplicate(self, content, url):
        # HTML 콘텐츠 해시 생성
        content_hash = sha256(content.encode('utf-8')).hexdigest()
        # DB에서 해당 해시 값이 있는지 확인
        cursor = self.conn.cursor()
        cursor.execute("SELECT id FROM downloaded_html WHERE file_hash = %s ", (content_hash, ))
        result = cursor.fetchone()

        if result:  # 중복된 콘텐츠가 있다면
            return True
        else:
            # 중복되지 않으면 DB에 저장
            self.save_to_db(content, content_hash, url)
            return False

    def save_to_db(self, content, content_hash, url):
        # 파일 크기와 데이터 유형을 파악
        file_size = len(content)
        data_type = "HTML"  # 예시로 HTML만

        cursor = self.conn.cursor()
        cursor.execute("""
            INSERT INTO downloaded_html (domain_name, file_hash, file_size, data_type)
            VALUES (%s, %s, %s, %s) RETURNING id
        """, (urlparse(url).netloc, content_hash, file_size, data_type))

        downloaded_html_id = cursor.fetchone()[0]
        self.conn.commit()

        # parsed_urls 테이블에 URL 추가
        cursor.execute("""
            INSERT INTO parsed_urls (downloaded_html_id, url)
            VALUES (%s, %s)
        """, (downloaded_html_id, url))
        self.conn.commit()

        print(f"도메인 저장 db 명  : {url}")


#콘텐츠 저장소
class ContentStorage:
    def __init__(self, directory):
        self.directory = directory
        os.makedirs(directory, exist_ok=True)  # 저장소 디렉토리 생성

    def save_content(self, domain, content): # 저장
        filename = os.path.join(self.directory, f"{domain}.html")
        with open(filename, 'w', encoding='utf-8') as f:
            f.write(content)
        print(f"도메인 생성: {domain}")

    def delete_content(self, domain): #삭제
        filename = os.path.join(self.directory, f"{domain}.html")
        if os.path.exists(filename):
            os.remove(filename)
            print(f"도메인 삭제: {domain}")


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

    conn = get_db_connection()
    url_queue = URLQueue()

    for url in start_urls:
        priority = classify_url(url)  # URL 분류
        url_queue.add_url(url, 0, priority)  # depth 0으로 시작, 우선순위

    downloader = HTMLDownloader()
    domain_converter = DomainNameConverter()
    parser = ContentParser()
    checker = DuplicateChecker()
    storage = ContentStorage(storage_directory)

    with ThreadPoolExecutor(max_workers=5) as executor:
        while not url_queue.is_empty():
            url, depth = url_queue.get_next_url()
            print(f"Crawling: {url} at depth: {depth}")

            if depth >= max_depth:
                continue

            if checker.url_is_duplicate(url):
                print(f" 중복 url : {url}")
                continue

            # HTML 다운로드
            html_content = downloader.download(url)
            if html_content is None:
                continue

            # 도메인 변환
            domain = domain_converter.get_domain(url)

            # 콘텐츠 파싱
            soup = parser.parse(html_content)

            # 중복 확인
            if checker.html_is_duplicate(html_content, url):
                print(f" 중복 url : {url}")
                continue

            # 콘텐츠 저장
            storage.save_content(domain, html_content)

            # URL 추출
            cursor = conn.cursor()
            cursor.execute("SELECT id FROM downloaded_html WHERE domain_name = %s ORDER BY download_date DESC LIMIT 1",
                           (domain,))
            downloaded_html_id = cursor.fetchone()
            if downloaded_html_id:
                downloaded_html_id = downloaded_html_id[0]

                # URL 파싱 후 저장
                extracted_urls = URLExtractor.extract_urls(soup, url)
                for extracted_url in extracted_urls:
                    cursor.execute("SELECT id FROM parsed_urls WHERE url = %s", (extracted_url,))
                    if cursor.fetchone():
                        print(f"중복 url 탐지 : {extracted_url}")
                        continue

                    # 우선순위에 맞게 URL을 큐에 추가
                    priority = classify_url(extracted_url)
                    url_queue.add_url(extracted_url, depth + 1, priority)

                    cursor.execute("""
                        INSERT INTO parsed_urls (url, downloaded_html_id)
                        VALUES (%s, %s)
                    """, (extracted_url, downloaded_html_id))
                    conn.commit()

            cursor.close()

            # HTML 파일 삭제
            storage.delete_content(domain)

    # 크롤링 후 parsed_urls에 있는 URL들을 크롤링
    process_urls()

    return url_queue


# URL 처리 함수 (parsed_urls에 들어간 URL을 다시 크롤링 큐에 넣어서 반복적으로 크롤링)
def process_urls():
    conn = get_db_connection()
    cursor = conn.cursor()

    # parsed_urls 테이블에 있는 모든 URL을 처리
    cursor.execute("SELECT url FROM parsed_urls")
    rows = cursor.fetchall()
    for row in rows:
        url = row[0]
        print(f"Processing URL: {url}")
        crawl([url], "crawled_contents", 1)  # 해당 URL을 재크롤링 하도록 함

    conn.commit()
    cursor.close()
    conn.close()

if __name__ == "__main__":
    app.run(debug=True, host="0.0.0.0", port=5002)