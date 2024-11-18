from sqlalchemy import DateTime, Index, Integer, Text, text, inspect
from sqlalchemy.orm import mapped_column

from .base import Base


class Urls(Base):
    __tablename__ = "urls"
    __table_args__ = (
        Index("idx_urls_status", "status"),  # 상태별 URL 조회 인덱스
        Index("idx_urls_priority", "priority"),  # 우선순위별 조회 인덱스
    )

    urls_id = mapped_column(Integer, primary_key=True, autoincrement=True)  # 고유 ID
    url = mapped_column(Text, nullable=False, unique=True)  # 크롤링할 URL
    domain = mapped_column(Text, nullable=False)  # URL의 도메인
    path = mapped_column(Text)  # URL의 경로
    status = mapped_column(Text, nullable=False, server_default=text("'pending'"))  # URL 상태
    priority = mapped_column(Integer, nullable=False, server_default=text("0"))  # 크롤링 우선순위
    reg_date = mapped_column(DateTime(True), server_default=text("now()"))  # URL이 추가된 시간
    update_date = mapped_column(DateTime(True), server_default=text("now()"))  # URL이 업데이트된 시간

    def to_dict(self):
        """테이블 데이터를 딕셔너리로 반환"""
        return {c.key: getattr(self, c.key) for c in inspect(self).mapper.column_attrs}
