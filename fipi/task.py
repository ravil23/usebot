from typing import Dict, List, NamedTuple, Optional

from bs4 import BeautifulSoup


class FIPITask(NamedTuple):
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
    def from_response(data: Dict) -> 'FIPITask':
        return FIPITask(
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
            text=BeautifulSoup(data['taskText'], 'html.parser').get_text().strip(),
            title=data['taskTitle'].strip(),
            version=data['taskVersion'],
            options={
                tag['number']: tag.get_text().strip()
                for tag in BeautifulSoup(data['html'], 'html.parser').find_all(attrs={'class': 'answer'})
            },
            doc=data['docHtml'] and BeautifulSoup(data['docHtml'], 'html.parser').get_text().strip(),
        )

    def to_dict(self) -> Dict:
        return {
            'id': self.id,
            'type': f'{self.type_name} [{self.type_id}]',
            'text': self.text,
            'answer': self.answer,
            'options': self.options,
        }
