from flask import Flask
from app.api.crawl import main  # main 블루프린트 가져오기

def crawl_app():
    app = Flask(__name__)
    app.config.from_object('app.config.Config')

    # Blueprint 등록
    app.register_blueprint(main)

    return app
