from dataclasses import dataclass, field


@dataclass
class FilterRule:
    EXCLUDE_EXTENSIONS: list[str] = field(default_factory=lambda: ["jpg", "png", "pdf", "zip"])
    ALLOWED_TYPES: list[str] = field(default_factory=lambda: ["text/html"])
    EXCLUSION_PATTERNS: list[str] = field(default_factory=lambda: [r"/admin", r"/login"])
