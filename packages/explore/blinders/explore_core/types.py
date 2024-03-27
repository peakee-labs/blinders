import dataclasses


@dataclasses.dataclass
class MatchInfo:
    name: str
    gender: str
    major: str
    native: str
    country: str
    learnings: list[str]
    interests: list[str]
    age: int

    def __str__(self) -> str:
        return ("[BEGIN]gender: {}[SEP]age: {}[SEP]job: {}[SEP]native language: {}[SEP]learning language: {}["
                "SEP]country: {}[SEP]interests: {}[END]").format(
            self.gender,
            self.age,
            self.major,
            self.native,
            ", ".join(self.learnings),
            self.country,
            ", ".join(self.interests),
        )
