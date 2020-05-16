import re
from typing import Dict, List, NamedTuple, Optional

from bs4 import BeautifulSoup


class Task(NamedTuple):
    subject_id: int
    level_id: int
    level_name: str
    answer: str
    theme_names: List[str]
    requirement_names: List[str]
    img_url: Optional[str]

    id: int
    type_id: int
    type_name: str
    text: str
    title: str
    version: int
    options: Dict
    doc: Optional[str]

    @staticmethod
    def from_response(data: Dict) -> 'Task':
        def clean_text(text: str) -> str:
            return re.sub(r" +", " ", re.sub(r"MathType.* ", "", text)).strip()

        return Task(
            subject_id=data['subjectId'],
            level_id=data['levelId'],
            level_name=data['levelName'].strip(),
            answer=data['answer'].strip(),
            theme_names=data['themeNames'],
            requirement_names=data['requirementNames'],
            img_url=data['imgUrl'],
            id=data['id'],
            type_id=data['taskTypeId'],
            type_name=data['taskTypeName'].strip(),
            text=clean_text(BeautifulSoup(data['taskText'], 'html.parser').get_text()),
            title=data['taskTitle'].strip(),
            version=data['taskVersion'],
            options={
                tag['number']: clean_text(tag.get_text())
                for tag in BeautifulSoup(data['html'], 'html.parser').find_all(attrs={'class': 'answer'})
                if len(clean_text(tag.get_text())) != 0
            },
            doc=data['docHtml'] and clean_text(BeautifulSoup(data['docHtml'], 'html.parser').get_text()),
        )

    def to_dict(self) -> Dict:
        return {
            'id': self.id,
            'type': self.type_id,
            'level': self.level_id,
            'text': self.text,
            'answer': self.answer,
            'options': self.options,
            'themes': self.theme_names,
            'doc': self.doc,
            'sendAsPoll': (
                len(self.text) <= 255 and
                self.doc is not None and
                all(len(option) <= 100
                    for option in self.options.values())
            )
        }
